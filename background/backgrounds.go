package background

import (
	"time"

	"github.com/bep/debounce"
	"go.uber.org/zap"
)

//Faire un ProcessNote "heritant" de Process

//mettre une interface dans Process (mettre potentiellment un isTheSame() qui fait la fct du dessous dans l'interface))
//dans ProcessNote mettre un string NoteId

/*
CheckProcessIdentifier(arg interface{ LambdaArg() }, processArg) {
	switch (arg.lambda()) {
		case (ProcessNote):
			if (arg == processArg) {
				return true
			} else {
				return false
			}
	}
} return un bool
*/

// TODO : What if the note is deleted before debounced for UpdateKeyword
func (srv *service) AddProcess(lambda func() error, arg interface{ LambdaArg() }) error {

	for index := range srv.processes {
		// Prohibits the process launched before with the same noteId to execute his callback fct
		//if GetProcessIdentifier(arg, srv.processes[index]) == true { cancel la task }
		if srv.processes[index].NoteId == getNoteId(arg) { // GetProcessIdentifier() {dwitch(type) => ProcessNote {NoteID} }
			// TODO cancel the goroutine by srv.processes.task
			go srv.processes[index].debounced(func() { return })
			srv.processes = remove(srv.processes, index)
		}
	}

	// Add a process to the list launch the debounce fct
	lastIndex := len(srv.processes)
	newProcess := Process{
		NoteId:    noteId,
		debounced: debounce.New(TimeToSave * time.Second),
		callBackFct: func() {
			// lambda
			err := lambda()
			//err := srv.UpdateKeywordsByNoteId(noteId)
			// ! lambda
			if err != nil {
				srv.logger.Error("failed update keywords on noteId : "+noteId, zap.Error(err))
				return
			}
			srv.processes = remove(srv.processes, lastIndex)
		},
	}

	srv.processes = append(srv.processes, newProcess)
	go srv.processes[lastIndex].debounced(srv.processes[lastIndex].callBackFct)
	return nil
}

func getNoteId(arg struct{ LambdaArg string }) string {
	return arg.LambdaArg
}

func remove(slice []Process, idx int) []Process {
	return append(slice[:idx], slice[idx+1:]...)
}
