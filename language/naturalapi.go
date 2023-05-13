package language

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"notes-service/models"
	"os"
	"strings"

	glanguage "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
	kgsearch "google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var protobufEnumToKeywordType = map[languagepb.Entity_Type]string{
	languagepb.Entity_UNKNOWN:       models.Unknown,
	languagepb.Entity_PERSON:        models.Person,
	languagepb.Entity_LOCATION:      models.Location,
	languagepb.Entity_ORGANIZATION:  models.Organization,
	languagepb.Entity_EVENT:         models.Event,
	languagepb.Entity_WORK_OF_ART:   models.WorkOfArt,
	languagepb.Entity_CONSUMER_GOOD: models.ConsumerGood,
	languagepb.Entity_OTHER:         models.Other,
	languagepb.Entity_PHONE_NUMBER:  models.PhoneNumber,
	languagepb.Entity_ADDRESS:       models.Address,
	languagepb.Entity_DATE:          models.Date,
	languagepb.Entity_NUMBER:        models.Number,
	languagepb.Entity_PRICE:         models.Price,
}

type KGDetailedDescription struct {
	ArticleBody string `json:"articleBody,omitempty"`
	URL         string `json:"url,omitempty"`
}

type KGImage struct {
	URL string `json:"url,omitempty"`
}

type NaturalAPIService struct {
	Service
	lClient   *glanguage.Client
	kgService *kgsearch.Service
}

func (s *NaturalAPIService) Init() error {
	jsonCredentialBase64 := os.Getenv("JSON_GOOGLE_CREDS_B64")

	if jsonCredentialBase64 == "" {
		s.lClient = nil
		return nil
	}

	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	if googleApiKey == "" {
		s.lClient = nil
		s.kgService = nil
		return nil
	}

	jsonCredential, err := base64.StdEncoding.DecodeString(jsonCredentialBase64)
	if err != nil {
		return err
	}

	client, err := glanguage.NewClient(context.Background(), option.WithCredentialsJSON(jsonCredential))
	if err != nil {
		return err
	}
	s.lClient = client

	service, err := kgsearch.NewService(context.Background(), option.WithAPIKey(googleApiKey))
	if err != nil {
		return err
	}
	s.kgService = service

	return nil
}

func (s *NaturalAPIService) doKnowledgeGraphSearch(keywords *map[string]*models.Keyword) (*kgsearch.SearchResponse, error) {
	mids := []string{}

	for mid := range *keywords {
		mids = append(mids, mid)
	}

	search := s.kgService.Entities.Search()
	search.Ids(mids...)
	search.Languages("fr")

	response, err := search.Do()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func kgInterfaceToStruct(i interface{}, s interface{}) error {
	asJson, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(asJson, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *NaturalAPIService) fillWithKnowledgeGraph(keywords *map[string]*models.Keyword) error {
	entityResult, err := s.doKnowledgeGraphSearch(keywords)
	if err != nil {
		return err
	}

	for _, element := range entityResult.ItemListElement {
		responseMap, ok := element.(map[string]interface{})
		if !ok {
			return errors.New("gkg has an invalid response")
		}

		entityResult, ok := responseMap["result"]
		if !ok {
			return errors.New("gkg response has no result for the keywords")
		}

		entityResultMap, ok := entityResult.(map[string]interface{})
		if !ok {
			return errors.New("gkg response has no result for the keywords")
		}

		mid := entityResultMap["@id"].(string)
		mid = strings.TrimPrefix(mid, "kg:")

		keyword := (*keywords)[mid]

		if detailedDescriptionInterface, ok := entityResultMap["detailedDescription"]; ok {
			detailedDescription := KGDetailedDescription{}

			err = kgInterfaceToStruct(detailedDescriptionInterface, &detailedDescription)
			if err != nil {
				// TODO : To log or not to log that is the question
			} else {
				keyword.Summary = detailedDescription.ArticleBody
				keyword.URL = detailedDescription.URL
			}

		}

		if imageInterface, ok := entityResultMap["image"]; ok {
			image := KGImage{}

			err = kgInterfaceToStruct(imageInterface, &image)
			if err != nil {
				// TODO : To log or not to log that is the question
			} else {
				keyword.ImageURL = image.URL
			}
		}

		if betterTypeInterface, ok := entityResultMap["description"]; ok {
			betterType, ok := betterTypeInterface.(string)
			if ok {
				keyword.Type = betterType
			} else {
				// TODO : To log or not to log that is the question (this one should never fail for sure)
			}
		}

	}

	return nil
}

func (s *NaturalAPIService) GetKeywordsFromTextInput(input string) ([]*models.Keyword, error) {
	if s.lClient == nil || s.kgService == nil {
		return nil, status.Error(codes.Unavailable, "credentials are not made for google's natural api or knowledge graph service")
	}

	req := &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Type:     languagepb.Document_PLAIN_TEXT,
			Source:   &languagepb.Document_Content{Content: input},
			Language: "fr", // Pass as a parameter ?
		}}

	res, err := s.lClient.AnalyzeEntities(context.Background(), req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var keywords []*models.Keyword
	keywordsWithMID := make(map[string]*models.Keyword)

	for _, entity := range res.Entities {
		newKeyword := models.Keyword{
			Keyword: entity.Name,
			Type:    protobufEnumToKeywordType[entity.Type],
		}

		if val, ok := entity.Metadata["wikipedia_url"]; ok {
			newKeyword.URL = val
		}

		if mid, ok := entity.Metadata["mid"]; ok {
			keywordsWithMID[mid] = &newKeyword
		} else if val, ok := entity.Metadata["wikipedia_url"]; ok {
			newKeyword.URL = val
		}

		keywords = append(keywords, &newKeyword)
	}

	err = s.fillWithKnowledgeGraph(&keywordsWithMID)
	if err != nil {
		// TODO: log
	}

	return keywords, nil
}
