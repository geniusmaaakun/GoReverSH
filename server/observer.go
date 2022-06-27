package server

import (
	"context"
	"fmt"
	"log"
)

type Observer struct {
	Sender         Sender
	State          State
	Subject        <-chan Notification
	PromptViewFlag bool
}

//受け取った通知の種別によってメッセージの送信, あるいはメンバーの追加/削除を行います
func (observer Observer) WaitNotice(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			notice := <-observer.Subject
			fmt.Println(notice.Type)
			switch notice.Type {
			case JOIN:
				observer.State.ClientMap[notice.Client.Name] = notice.Client
				if observer.Sender.connectingClient == nil {
					observer.Sender.connectingClient = notice.Client
				}
				observer.printPrompt()

			case DEFECT:
				fmt.Println("Connection Close")
				notice.Client.Conn.Close()
				delete(observer.State.ClientMap, notice.Client.Name)
				if 0 < len(observer.State.ClientMap) {
					for _, c := range observer.State.ClientMap {
						observer.Sender.connectingClient = c
						break
					}
				}
				observer.printPrompt()

			case COMMAND:
				fmt.Println(notice.Command)
				//write
				//observer.Sender.SendMessage(notice.Command)

				observer.printPrompt()

			case MESSAGE:
				fmt.Printf("\n%s\n", notice.Message)
				observer.printPrompt()

			case CHANGE_DIR:

			case UPLOAD:

			case DOWNLOAD:

			case SCREEN_SHOT:

			case CLEAN:

			default:
				log.Println(notice)
			}
		}
	}
}

func (o Observer) printPrompt() {
	if o.Sender.connectingClient != nil {
		fmt.Printf("[GoReverSH@%s] >", o.Sender.connectingClient.Name)
	} else {
		fmt.Println("wait...")
	}
}
