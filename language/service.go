package language

import (
	"notes-service/models"
)

type Service interface {
	Init() error
	SetLanguage(lang string)
	GetKeywordsFromTextInput(input string) ([]*models.Keyword, error)
	GenerateQuizFromTextInput(input string) (*models.Quiz, error)
	GenerateSummaryFromTextInput(input string) (*models.Summary, error)
}
