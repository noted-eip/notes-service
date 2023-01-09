package background

import (
	"notes-service/models"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	repoNote       models.NotesRepository
	processManager backGroundProcesses
}

const (
	TimeToSave time.Duration = 5 //5 seconds
	//TimeToSave time.Duration = 900 //15 minutes
)

type process struct {
	NoteId string

	debounced   func(func())
	callBackFct func()

	Quit (chan bool)
}

type backGroundProcesses struct {
	processes (chan []process)
}
