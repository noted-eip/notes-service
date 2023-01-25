package background

type Process struct {
	task                          uint32
	identifier                    interface{}
	debounced                     func(func())
	callBackFct                   func()
	secondsToDebounce             uint32
	cancelProcessOnSameIdentifier bool
	//repeatProcess                 bool
}

type BackGroundProcesses []Process

type ProcessPlayLoad struct {
	// Identifier of the process
	Identifier interface{}
	// Function called by the process
	CallBackFct func() error
	// Seconds to wait before the execution of callBackFct
	SecondsToDebounce uint32
	// If the same identifier is add in process list, the seconds to debounce are going to be reset
	CancelProcessOnSameIdentifier bool
	// TODO : implem
	// If you want to shedule the reccurent execution of the process each SecondsToDebounce
	//RepeatProcess bool
}
