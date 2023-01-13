package background

import (
	"time"

	"github.com/bep/debounce"
)

func (srv *Service) AddProcess(noteId string) error {

	for _, currentProcess := range srv.processManager.processes {
		if currentProcess.NoteId == noteId {
			// suprimmer l'ancienne goroutine
			// en relancer une autre
			// suprimer le process[idx]
			currentProcess.Quit <- true
			//c sensé faire return l'ancienne goroutine et en lancer une autre
			println("Process already exist / NotedId = " + currentProcess.NoteId)
		}
	}

	println("Save new process / NoteId = " + noteId)
	lastIndex := len(srv.processManager.processes)
	newSize := lastIndex + 1
	srv.processManager.processes = make([]process, newSize)
	actualProcess := srv.processManager.processes

	actualProcess[lastIndex].NoteId = noteId
	actualProcess[lastIndex].debounced = debounce.New(TimeToSave * time.Second)
	actualProcess[lastIndex].callBackFct = func() {
		/*lock, _ := <-srv.processManager.processes[lastIndex].Quit
		if !lock {
			println("Function was called but canceled / NoteId : " + noteId)
			return
		}*/
		err := srv.UpdateKeywordsByNoteId(noteId)
		if err != nil {
			println("Failed to save new keywords / NoteId : " + noteId)
			return
		}
		println("Save new keywords / NoteId : " + noteId)
		//la variable est pas modifié dans RunProcess
		srv.processManager.processes[lastIndex].Quit <- false
		//erease the process inthe array
		srv.processManager.processes = append(srv.processManager.processes) //oposite
	}

	println("Launch process / NoteId = " + noteId)
	//go actualProcess[lastIndex].debounced(actualProcess[lastIndex].callBackFct)
	srv.RunProcess(&actualProcess[lastIndex])

	return nil
}

func (srv *Service) RunProcess(process *process) {
	process.Quit = make(chan bool)
	//process.Quit <- false
	go func() {
		for {
			value, ok := <-process.Quit

			if !value || !ok {
				println("---go return")
				return
			} else {
				//ca passe pas la
				println("---go debounce")
				process.debounced(process.callBackFct)
			}
		}
	}()
	// …
	process.Quit <- true
}

func (srv *Service) Run() error {
	return nil
}
