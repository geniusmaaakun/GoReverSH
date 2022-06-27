package server

type Sender struct {
	connectingClient *Client
}

func (sender Sender) SendMessage(message string) error {
	var buf = []byte(message)

	_, err := sender.connectingClient.Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
