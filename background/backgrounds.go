package background

import (
	"strconv"
	"time"

	"github.com/bep/debounce"
	"go.uber.org/zap"
)

func (srv *service) AddProcess(process *ProcessPlayLoad) error {

	if process.CancelProcessOnSameIdentifier {
		err := srv.cancelProcessOnSameIdentifier(process)
		if err != nil {
			return err
		}
	}

	// Add a process to the list & launch the debounce fct
	lastIndex := len(srv.processes)
	newProcess := Process{
		identifier: process.Identifier,
		debounced:  debounce.New(time.Duration(process.SecondsToDebounce) * time.Second),
		callBackFct: func() {
			err := process.CallBackFct()
			if err != nil {
				srv.logger.Error("Error in Lambda function in backgroundProcess for task : "+strconv.Itoa(int(srv.processes[lastIndex].task)), zap.Error(err))
				return
			}
			srv.processes = remove(srv.processes, lastIndex)
		},
		secondsToDebounce:             process.SecondsToDebounce,
		cancelProcessOnSameIdentifier: process.CancelProcessOnSameIdentifier,
	}

	srv.processes = append(srv.processes, newProcess)
	go srv.processes[lastIndex].debounced(srv.processes[lastIndex].callBackFct)
	return nil
}

func (srv *service) cancelProcessOnSameIdentifier(process *ProcessPlayLoad) error {
	for index := range srv.processes {
		if !srv.processes[index].cancelProcessOnSameIdentifier {
			continue
		}
		//c pas le meme type si ?
		if srv.processes[index].identifier == process.Identifier {
			// TODO cancel the goroutine by srv.processes.task
			go srv.processes[index].debounced(func() { return })
			srv.processes = remove(srv.processes, index)
		}
	}
	return nil
}

func remove(slice []Process, idx int) []Process {
	return append(slice[:idx], slice[idx+1:]...)
}
