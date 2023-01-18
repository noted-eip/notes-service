package background

import "time"

const (
	TimeToSave time.Duration = 900 //15 minutes
)

type Process struct {
	NoteId      string
	task        uint32
	debounced   func(func())
	callBackFct func()
}

type BackGroundProcesses []Process

type backGroundProcesses struct {
	processes []Process
}
