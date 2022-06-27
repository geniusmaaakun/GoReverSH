package server

import (
	"context"
)

type Observer struct {
	Sender         *Sender
	State          *State
	Subject        <-chan Notification
	PromptViewFlag bool
}

//受け取った通知の種別によってメッセージの送信, あるいはメンバーの追加/削除を行います
func (observer Observer) WaitNotice(ctx context.Context) {
	for {
		notice := <-observer.Subject

		switch notice.Type {

		case JOIN:

		case DEFECT:

		case CHANGE_DIR:

		case UPLOAD:

		case DOWNLOAD:

		case SCREEN_SHOT:

		case CLEAN:

		default:

			/*
				//メッセージの処理
				case Message:
					for i := range observer.Senders {
						observer.Senders[i].SendMessage(notice.Message)
					}

					break

					//参加処理
				case Join:
					observer.Senders = appendSender(notice.ClientId, notice.Connection, observer.Senders)

					fmt.Printf("Client %d join, now menber count is %d\n", notice.ClientId, len(observer.Senders))
					break

					//退出処理
				case Defect:
					observer.Senders = removeSender(notice.ClientId, observer.Senders)

					fmt.Printf("Client %d defect, now menber count is %d\n", notice.ClientId, len(observer.Senders))
					break

				default:
			*/

		}

	}
}

/*
//スライスに追加
func appendSender(senderId int, connection net.Conn, senders []Sender) []Sender {
	return append(senders, Sender{Id: senderId, Connection: connection})
}

func removeSender(senderId int, senders []Sender) []Sender {
	var find = -1

	for i := range senders {
		if senders[i].Id == senderId {
			find = i
			break
		}
	}

	if find == -1 {
		return senders
	}

	//スライスから削除
	return append(senders[:find], senders[find+1:]...)
}
*/

type State struct {
	ClientMap map[string]*Client
}

/*
func (o *Observer) NewClient(conn net.Conn, name string) *Client {
	client := &Client{conn: conn, name: name, addr: conn.RemoteAddr().String()}
	o.State.ClientMap[client.name] = client
	if o.Sender.connectingClient == nil {
		o.Sender.connectingClient = client
	}
	return client
}
*/
