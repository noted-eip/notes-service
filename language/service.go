package language

import (
	"notes-service/models"
)

type Service interface {
	Init() error
	GetKeywordsFromTextInput(input string, lang string) ([]*models.Keyword, error)
	GenerateQuizFromTextInput(input string, lang string) (*models.Quiz, error)
	GenerateSummaryFromTextInput(input string, lang string) (*models.Summary, error)
}
