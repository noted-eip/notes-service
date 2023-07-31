package language

import (
	"notes-service/models"
)

type Service interface {
	Init() error
	GetKeywordsFromTextInput(input string) ([]*models.Keyword, error)
	GenerateQuizFromTextInput(input string) (*models.Quiz, error)
}
