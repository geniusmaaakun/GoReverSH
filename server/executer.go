package server

import (
	"bufio"
	"context"
)

type Executer struct {
	Scanner  *bufio.Scanner
	Observer chan<- Notification
}

func (e Executer) WaitCommand(context.Context) {
	for e.Scanner.Scan() {
		//commands := strings.Split(e.Scanner.Text(), " ")
		//fmt.Println(commands)
		//fmt.Println(e.Scanner.Text())
		e.Observer <- Notification{Type: COMMAND, Command: e.Scanner.Text()}
	}
}
