package language

import (
	"notes-service/models"
)

type Service interface {
	Init() error
	GetKeywordsFromTextInput(input string) (*models.Keywords, error)
}
