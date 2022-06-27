package server

import "net"

type Client struct {
	conn net.Conn
	name string
	addr string
}

func NewClient(conn net.Conn, name string) *Client {
	client := &Client{conn: conn, name: name, addr: conn.RemoteAddr().String()}
	return client
}
