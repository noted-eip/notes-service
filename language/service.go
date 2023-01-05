package language

import (
	"context"
	"notes-service/models"
)

type Service struct {
}

func (s *Service) GetKeywordsFromTextInput(ctx *context.Context, input string) (*models.Keywords, error) {
	return s.GetKeywordsFromGoogleNaturalApi(ctx, input)
}
