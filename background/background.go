package background

import (
	"time"
)

const (
	SecondsToWait time.Duration = 900 //15 minutes
)

type result struct {
	Task int32
	Data int32
	err  error
}

type process struct {
	Task   int32
	NoteId string
	//SecondsToWait time.Duration//utile ?
	//Clock time.Duration//utile ?
}

type backGroundProcesses struct {
	Processes []process
}

//gen les keyWord & save en db
//j'ai besoin de repoNote pour save les keywords

//c pas ca qu'il faut faire
func (srv *backGroundProcesses) AddProcess(noteId string) chan bool {
	quit := make(chan bool)

	for _, process := range srv.Processes {
		if process.NoteId == noteId {
			//cancel the process
			quit <- true
			return quit
		}
	}

	//last index
	lastIndex := len(srv.Processes)
	srv.Processes[lastIndex].NoteId = noteId
	srv.Processes[lastIndex].Task = int32(lastIndex)

	quit <- false
	return quit
}

func Save(noteId string) {
	time.Sleep(SecondsToWait * 1000 * time.Millisecond)
	//gen les keyWord
}

//si on appelle 15 fois launchSaveClock il faut que ca save qu'une seule fois par note
//lancer dans le main 1 fois
func (srv *backGroundProcesses) LaunchSaveClock(noteId string) {
	for {
		select {
		//non ca va add en boucle, pas comme ca
		case <-srv.AddProcess(noteId):
			return
		default:
			return
			//time.NewTimer(1)
			/*clock += timePassed
			if clock >= timeToWait {
				fmt.Println("Save for note " + noteId)
				return
			}*/
		}
	}
}
