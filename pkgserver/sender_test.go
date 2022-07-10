package pkgserver

import (
	"GoReverSH/utils"
	"testing"
)

//net.Conn mock
func TestSendMessage(t *testing.T) {
	s := Sender{connectingClient: nil}
	err := s.SendMessage("test")
	if err == nil {
		t.Errorf("connecting client is nil but %v\n", err)
	}
}

func TestSendOutputJson(t *testing.T) {
	s := Sender{connectingClient: nil}
	err := s.SendOutputJson(utils.Output{})
	if err == nil {
		t.Errorf("connecting client is nil but %v\n", err)
	}
}
