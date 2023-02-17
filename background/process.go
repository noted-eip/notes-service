package background

type Process struct {
	// Goroutine id of the process
	task uint32
	// Object which describe the process (any object)
	Identifier interface{}
	// Storage of debounce from the pakcage debounce
	debounced func(func())
	// Function to call
	CallBackFct func() error
	// Seconds before the CallBackFct is called
	SecondsToDebounce uint32
	// Call AddProcess with the same identifier to restart the time until the callback is called
	CancelProcessOnSameIdentifier bool
	// Repeat the process each SecondsToDebounce
	RepeatProcess bool
}

type BackGroundProcesses []Process
