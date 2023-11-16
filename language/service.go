package language

import (
	"notes-service/models"

	"go.uber.org/zap"
)

type Service interface {
	Init(*zap.Logger) error
	GetKeywordsFromTextInput(input string, lang string) ([]*models.Keyword, error)
	GenerateQuizFromTextInput(input string, lang string) (*models.Quiz, error)
	GenerateSummaryFromTextInput(input string, lang string) (*models.Summary, error)
}
