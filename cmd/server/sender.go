package server

type Sender struct {
	connectingClient *Client
}

func (sender Sender) SendMessage(message string) {
	var buf = []byte(message)

	_, error := sender.Connection.Write(buf)
	if error != nil {
		panic(error)
	}
}
