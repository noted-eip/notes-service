package background

import (
	"errors"
	"strconv"
	"time"

	"github.com/bep/debounce"
	"go.uber.org/zap"
)

// Add a process to the queue which is going to exec a function in X time
func (srv *service) AddProcess(process *Process) error {
	// Check for illegal processes
	if process.Identifier == nil && process.RepeatProcess == true {
		srv.logger.Error("You can't repeat a process with a nil indentifier, the process could never stop")
		return errors.New("Error : Process cannot repeat and have a nil identifier")
	}
	// Cancel the process of the same identifier
	if process.CancelProcessOnSameIdentifier {
		err := srv.CancelProcess(process)
		if err != nil {
			return err
		}
	}
	// Add a process to the list & launch the debounce fct
	lastIndex := len(srv.processes)
	process.debounced = debounce.New(time.Duration(process.SecondsToDebounce) * time.Second)
	srv.processes = append(srv.processes, *process)
	srv.debounceLogic(&srv.processes[lastIndex], lastIndex)
	return nil
}

// Cancel a process by his identifier
func (srv *service) CancelProcess(process *Process) error {
	if process.Identifier == nil {
		srv.logger.Error("Cannot cancel background process if the identifier is nil")
		return errors.New("Error : Identifier cannot be nil")
	}

	for index := 0; index < len(srv.processes); index++ {
		if srv.processes[index].Identifier == process.Identifier {
			// TODO cancel the goroutine by srv.processes.task
			go srv.processes[index].debounced(func() { return })
			srv.processes = remove(srv.processes, index)
			index--
		}
	}
	return nil
}

func (srv *service) debounceLogic(process *Process, index int) {
	logic := func() {
		err := process.CallBackFct()
		if err != nil {
			srv.logger.Error("Error in Lambda function in backgroundProcess for task : "+strconv.Itoa(int(srv.processes[index].task)), zap.Error(err))
			return
		}
		if process.RepeatProcess {
			srv.debounceLogic(process, index)
		} else {
			srv.processes = remove(srv.processes, index)
		}
	}
	go process.debounced(logic)
}

func remove(slice []Process, idx int) []Process {
	return append(slice[:idx], slice[idx+1:]...)
}
