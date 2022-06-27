package server

import "net"

type Client struct {
	Conn net.Conn
	Name string
	Addr string
}

func NewClient(conn net.Conn, name string) *Client {
	client := &Client{Conn: conn, Name: name, Addr: conn.RemoteAddr().String()}
	return client
}
