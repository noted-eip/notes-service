package background

import (
	"notes-service/language"
	"notes-service/models"

	"go.uber.org/zap"
)

type Service interface {
	AddProcess(noteId string) error
}

type service struct {
	logger    *zap.Logger
	repoNote  models.NotesRepository
	repoBlock models.BlocksRepository
	language  language.Service
	processes []Process
}

func NewService(logger *zap.Logger, repoNote models.NotesRepository, repoBlock models.BlocksRepository, language language.Service) Service {
	return &service{
		logger:    logger,
		repoNote:  repoNote,
		repoBlock: repoBlock,
		language:  language,
	}
}
