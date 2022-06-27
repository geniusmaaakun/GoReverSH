package server

import "context"

type Executer struct {
	Observer chan<- Notification
}

func (e Executer) WaitCommand(context.Context) {
	for {

	}
}
