package language

import (
	"context"
	"notes-service/models"
	"os"

	glang "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
	goption "google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var protobufEnumToKeywordType = map[languagepb.Entity_Type]models.KeywordType{
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

func (s *Service) GetKeywordsFromGoogleNaturalApi(ctx *context.Context, input string) (*models.Keywords, error) {
	// NOTE: Temporary while we get fixed on what we use as auth
	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		return nil, status.Error(codes.InvalidArgument, "Please set GOOGLE_API_KEY env variable")
	}

	client, err := glang.NewClient(*ctx, goption.WithAPIKey(apiKey))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	req := &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Type:     languagepb.Document_PLAIN_TEXT,
			Source:   &languagepb.Document_Content{Content: input},
			Language: "fr", // Pass as a parameter ?
		}}

	res, err := client.AnalyzeEntities(*ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var keywords models.Keywords

	for _, entity := range res.Entities {
		newKeyword := models.Keyword{
			Keyword: entity.Name,
			Type:    protobufEnumToKeywordType[entity.Type],
		}

		if val, ok := entity.Metadata["wikipedia_url"]; ok {
			newKeyword.URL = val
		}

		keywords = append(keywords, newKeyword)
	}

	return &keywords, nil
}
