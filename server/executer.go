package server

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
)

type Executer struct {
	Scanner  *bufio.Scanner
	Observer chan<- Notification
}

func NewExecuter(channel chan Notification) *Executer {
	//実行コマンドを受け取る
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	executer := &Executer{Scanner: scanner, Observer: channel}
	return executer
}

func (e Executer) WaitCommand(ctx context.Context) error {
	for e.Scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			return ctx.Err()
		default:
			commands := strings.Split(e.Scanner.Text(), " ")
			switch commands[0] {
			case "clist":
				e.Observer <- Notification{Type: CLIST, Command: e.Scanner.Text()}
			case "cswitch":
				e.Observer <- Notification{Type: CSWITCH, Command: e.Scanner.Text()}
			case "upload":
				e.Observer <- Notification{Type: UPLOAD, Command: e.Scanner.Text()}
			case "download":
				e.Observer <- Notification{Type: DOWNLOAD, Command: e.Scanner.Text()}
			case "screenshot":
				e.Observer <- Notification{Type: SCREEN_SHOT, Command: e.Scanner.Text()}
			case "clean_go_reversh":
				e.Observer <- Notification{Type: CLEAN, Command: e.Scanner.Text()}
			default:
				e.Observer <- Notification{Type: COMMAND, Command: e.Scanner.Text()}
			}
		}
	}
	if err := e.Scanner.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
