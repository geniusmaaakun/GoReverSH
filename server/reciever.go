package server

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
)

type Receiver struct {
	Client   *Client
	Observer chan<- Notification
	Lock     *sync.Mutex
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

			//チャネルで通知　メッセージ受信
			receiver.Observer <- Notification{Type: MESSAGE, Client: receiver.Client, Message: string(buf[:n])}
		}
	}
}
