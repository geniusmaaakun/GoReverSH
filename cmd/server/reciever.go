package server

type Receiver struct {
	Client   *Client
	Observer chan<- Notification
}

func (receiver Receiver) Start() {
	//参加通知
	receiver.Observer <- Notification{Type: Join, ClientId: receiver.Id, Connection: receiver.Connection}

	receiver.WaitMessage()
}

func (receiver Receiver) WaitMessage() {
	var buf = make([]byte, 1024)

	//通信切断を検知したら, メンバー退室の通知を行い, 処理を終了します.
	n, error := receiver.Connection.Read(buf)
	if error != nil {
		//退出通知
		receiver.Observer <- Notification{Type: Defect, ClientId: receiver.Id}
		return
	}

	//チャネルで通知　メッセージ受信
	receiver.Observer <- Notification{Type: Message, ClientId: receiver.Id, Message: string(buf[:n])}

	receiver.WaitMessage()
}
