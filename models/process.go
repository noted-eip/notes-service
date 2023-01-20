// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import "time"

const (
	TimeToSave time.Duration = 900 //15 minutes
)

type Process struct {
	Task        uint32
	Identifier  interface{}
	Debounced   func(func())
	CallBackFct func()
}

type BackGroundProcesses []Process

type Action uint

const (
	NoteUpdateKeyword Action = 1
	// Put in enum the other type of actions
	//...
)
