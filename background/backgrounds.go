package background

import (
	"fmt"
	"notes-service/models"
	"time"

	"github.com/bep/debounce"
	"go.uber.org/zap"
)

func NewService(logger *zap.Logger, repoNote models.NotesRepository) Service {
	println("Init Background Service...")
	return Service{
		logger:   logger,
		repoNote: repoNote,
	}
}

func (srv *Service) Save(noteId string, process process) {
	fmt.Println("Save for note " + noteId)
	//gen & save les keyWords
	process.Quit <- true
}

func (srv *Service) AddProcess(noteId string) error {
	for _, currentProcess := range <-srv.processManager.processes {
		//currentProcess = <- make(chan process)
		if currentProcess.NoteId == noteId {
			//cancel the currentProcess
			currentProcess.Quit <- true
			println("Stop task with notedId = " + currentProcess.NoteId)
			return nil
		}
	}
	println("Save new process / NoteId = " + noteId)
	lastIndex := len(srv.processManager.processes)
	newSize := lastIndex + 1
	srv.processManager.processes = make(chan []process, newSize)
	actualProcess := <-srv.processManager.processes
	actualProcess[lastIndex].NoteId = noteId
	actualProcess[lastIndex].debounced = debounce.New(10000000 * time.Millisecond)
	actualProcess[lastIndex].callBackFct = func() {
		println("Note saved id : " + noteId)
	}

	println("Launch process / NoteId = " + noteId)
	actualProcess[lastIndex].debounced(actualProcess[lastIndex].callBackFct)

	return nil
}

func (srv *Service) Run() {
	//println("Running Background Service")
	/*for true {
		//println("-----------loop------------")
		//println("process array size = " + strconv.Itoa(len(srv.processManager.processes)))
		for _, process := range <-srv.processManager.processes {
			select {
			//non ca va add en boucle, pas comme ca
			case <-process.Quit:
				return
			default:

				start := time.Now()
				timeElapsed := time.Since(start)

				process.Clock += timeElapsed
				if process.Clock >= TimeToSave {
					go srv.Save(process.NoteId, process)
					return
				}
			}
		}
	}*/

}

/*
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
*/
