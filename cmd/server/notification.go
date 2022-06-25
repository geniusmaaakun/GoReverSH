package main

import "net"

type NotificationType int

const (
	NORMAL NotificationType = iota
	CHANGE_DIR
	UPLOAD
	DOWNLOAD
	SCREEN_SHOT
	CLEAN
)

type Notification struct {
	Type       NotificationType
	Commands   []string
	Connection net.Conn
}
