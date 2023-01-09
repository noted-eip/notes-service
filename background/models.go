package background

import (
	"notes-service/language"
	"notes-service/models"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	repoNote       models.NotesRepository
	repoBlock      models.BlocksRepository
	language       language.Service
	processManager backGroundProcesses
}

func NewService(logger *zap.Logger, repoNote models.NotesRepository, repoBlock models.BlocksRepository, language language.Service) Service {
	return Service{
		logger:    logger,
		repoNote:  repoNote,
		repoBlock: repoBlock,
		language:  language,
	}
}

const (
	TimeToSave time.Duration = 1 //5 seconds
	//TimeToSave time.Duration = 900 //15 minutes
)

type process struct {
	NoteId string

	debounced   func(func())
	callBackFct func()

	Quit (chan bool)
}

type backGroundProcesses struct {
	processes ( /*chan*/ []process)
}
