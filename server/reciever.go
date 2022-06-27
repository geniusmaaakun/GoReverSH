package server

import (
	"context"
	"errors"
	"io"
	"log"
)

type Receiver struct {
	Client   *Client
	Observer chan<- Notification
}

func (receiver Receiver) Start(ctx context.Context) {
	//参加通知
	receiver.Observer <- Notification{Type: JOIN, Client: receiver.Client}

	receiver.WaitMessage(ctx)
}

func (receiver Receiver) WaitMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			return
		default:
			var buf = make([]byte, 1024)

			//通信切断を検知したら, メンバー退室の通知を行い, 処理を終了します.
			n, err := receiver.Client.Conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					//退出通知
					receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
					return
				}
				//退出通知
				receiver.Observer <- Notification{Type: DEFECT, Client: receiver.Client}
				log.Println(err)
				return
			}

			//チャネルで通知　メッセージ受信
			receiver.Observer <- Notification{Type: MESSAGE, Client: receiver.Client, Message: string(buf[:n])}
		}
	}
}
