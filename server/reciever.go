package server

import (
	"GoReverSH/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type Receiver struct {
	Client   *Client
	Observer chan<- Notification
	Lock     *sync.Mutex
}

func NewReceiver(conn net.Conn, name string, channel chan Notification, lock *sync.Mutex) *Receiver {
	client := NewClient(conn, name)
	receiver := &Receiver{Client: client, Observer: channel, Lock: lock}
	return receiver
}

func (receiver Receiver) Start(ctx context.Context) {
	//参加通知
	receiver.Observer <- Notification{Type: JOIN, Client: receiver.Client}

	receiver.WaitMessage(ctx)
}

func (receiver Receiver) WaitMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			return ctx.Err()
		default:
			/*
				var buf = make([]byte, 1024)

				//通信切断を検知したら, メンバー退室の通知を行い, 処理を終了します.
				n, err := receiver.Client.Conn.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						//退出通知
						receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
						return err
					}
					//退出通知
					receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
					log.Println(err)
					return err
				}
			*/

			/*
				var content []byte
				buff := make([]byte, 1024)
				size := 0

				for {
					//ファイルを読み込む
					n, err := receiver.Client.Conn.Read(buff)
					if err != nil {
						if errors.Is(err, io.EOF) {
							receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
							return err
						}
						receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
						log.Println(err)
						return err
					}
					//content + bufの中身を一時的に保存。
					tmp := make([]byte, 0, size+n)
					tmp = append(content[:size], buff[:n]...)
					content = tmp
					size += n
					if n < 1024 {
						break
					}
				}
			*/

			//fmt.Println(string(content))

			//チャネルで通知　メッセージ受信
			//receiver.Observer <- Notification{Type: MESSAGE, Client: receiver.Client, Message: string(buf[:n])}

			output := utils.Output{}
			err := json.NewDecoder(receiver.Client.Conn).Decode(&output)
			if err != nil {
				if errors.Is(err, io.EOF) {
					receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
					return err
				}
				log.Println(err)
				return err
			}

			switch output.Type {
			case utils.MESSAGE:
				receiver.Observer <- Notification{Type: MESSAGE, Client: receiver.Client, Output: output}
			case utils.FILE:
				//image data decode
				base64ToImage, err := base64.StdEncoding.DecodeString(string(output.Body))
				if err != nil {
					log.Println(err)
				}
				output.FileInfo.Body = base64ToImage
				//fmt.Println(string(base64ToImage))
				receiver.Observer <- Notification{Type: CREATE_FILE, Client: receiver.Client, Output: output}

			}
		}
	}
}
