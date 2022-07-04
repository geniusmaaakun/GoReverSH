package server

import (
	"net"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := net.Listen("tcp", ":8000")
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		t.Error(err)
	}
	c := NewClient(conn, "test")
	if c == nil {
		t.Errorf("NewClient got %v\n", c)
	}
}
