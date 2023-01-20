package background

import (
	"notes-service/models"
	"strconv"
	"time"

	"github.com/bep/debounce"
	"go.uber.org/zap"
)

// TODO : What if the note is deleted before debounced for UpdateKeyword
func (srv *service) AddProcess(lambdaFct func() error, identifier interface{}) error {

	for index := range srv.processes {
		if srv.processes[index].Identifier == identifier {
			// TODO cancel the goroutine by srv.processes.task
			go srv.processes[index].Debounced(func() { return })
			srv.processes = remove(srv.processes, index)
		}
	}

	// Add a process to the list launch the debounce fct
	lastIndex := len(srv.processes)
	newProcess := models.Process{
		Identifier: identifier,
		Debounced:  debounce.New(models.TimeToSave * time.Second),
		CallBackFct: func() {
			err := lambdaFct()
			if err != nil {
				srv.logger.Error("Error in Lambda function in backgroundProcess for task : "+strconv.Itoa(int(srv.processes[lastIndex].Task)), zap.Error(err))
				return
			}
			srv.processes = remove(srv.processes, lastIndex)
		},
	}

	srv.processes = append(srv.processes, newProcess)
	go srv.processes[lastIndex].Debounced(srv.processes[lastIndex].CallBackFct)
	return nil
}

func remove(slice []models.Process, idx int) []models.Process {
	return append(slice[:idx], slice[idx+1:]...)
}
