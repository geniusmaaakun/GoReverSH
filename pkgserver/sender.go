package pkgserver

import (
	"GoReverSH/utils"
	"encoding/json"
	"errors"
)

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

func (sender Sender) SendOutputJson(output utils.Output) error {
	if sender.connectingClient == nil {
		return errors.New("connection is not exist")
	}
	err := json.NewEncoder(sender.connectingClient.Conn).Encode(output)
	if err != nil {
		return err
	}
	return nil
}
