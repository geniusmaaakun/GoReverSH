package server

import "errors"

type Sender struct {
	connectingClient *Client
}

func (sender Sender) SendMessage(message string) error {
	if sender.connectingClient == nil {
		return errors.New("connection is not exist")
	}
	var buf = []byte(message)

	_, err := sender.connectingClient.Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
