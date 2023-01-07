package background

import (
	"fmt"
	"notes-service/models"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	repoNote       models.NotesRepository
	repoBlock      models.BlocksRepository
	processManager backGroundProcesses
}

/*
type result struct {
	Task int32
	Data int32
	err  error
}
*/

const (
	TimeToSave time.Duration = 5 //5 seconds
	//TimeToSave time.Duration = 900 //15 minutes
)

type process struct {
	Task   int32
	NoteId string
	//TimeToSave time.Duration//utile ?
	Clock time.Duration //utile ?
	Quit  (chan bool)
}

type backGroundProcesses struct {
	processes (chan []process)
}

func NewService(logger *zap.Logger, repoNote models.NotesRepository, repoBlock models.BlocksRepository) Service {
	println("Init Background Service...")
	return Service{
		logger:    logger,
		repoNote:  repoNote,
		repoBlock: repoBlock,
	}
}

func (srv *Service) AddProcess(noteId string) error {
	for _, currentProcess := range <-srv.processManager.processes {
		//currentProcess = <- make(chan process)
		if currentProcess.NoteId == noteId {
			//cancel the currentProcess
			currentProcess.Quit <- true
			print("Stop task " + strconv.Itoa(int(currentProcess.Task)))
			return nil
		}
	}

	lastIndex := len(srv.processManager.processes)
	newSize := lastIndex + 1

	srv.processManager.processes = make(chan []process, newSize)
	actualProcess := <-srv.processManager.processes
	actualProcess[lastIndex].NoteId = noteId
	actualProcess[lastIndex].Task = int32(lastIndex)
	println("Add process " + strconv.Itoa(int(actualProcess[lastIndex].Task)))
	return nil
}

func (srv *Service) Save(noteId string, process process) {
	fmt.Println("Save for note " + noteId)
	//gen & save les keyWords
	process.Quit <- true
}

func (srv *Service) Run() {
	println("Running Background Service")
	for true {
		//println("-----------loop------------")
		//println("process array size = " + strconv.Itoa(len(srv.processManager.processes)))
		for _, process := range <-srv.processManager.processes {
			select {
			//non ca va add en boucle, pas comme ca
			case <-process.Quit:
				return
			default:
				println("process " + strconv.Itoa(int(process.Task)) + " running")

				start := time.Now()
				timeElapsed := time.Since(start)

				process.Clock += timeElapsed
				if process.Clock >= TimeToSave {
					go srv.Save(process.NoteId, process)
					return
				}
			}
		}
	}
}
