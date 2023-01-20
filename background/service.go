package background

import (
	"notes-service/language"
	"notes-service/models"

	"go.uber.org/zap"
)

type Service interface {
	AddProcess(lambda func() error, arg interface{}) error
}

type service struct {
	logger    *zap.Logger
	repoNote  models.NotesRepository
	repoBlock models.BlocksRepository
	language  language.Service
	processes models.BackGroundProcesses
}

func NewService(logger *zap.Logger, repoNote models.NotesRepository, repoBlock models.BlocksRepository, language language.Service) Service {
	return &service{
		logger:    logger,
		repoNote:  repoNote,
		repoBlock: repoBlock,
		language:  language,
	}
}
