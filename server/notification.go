package server

import "GoReverSH/utils"

type NotificationType int

const (
	JOIN NotificationType = iota
	DEFECT
	MESSAGE
	COMMAND
	UPLOAD
	DOWNLOAD
	SCREEN_SHOT
	CLEAN
	CREATE_FILE
	MAKE_DIR
	CLIST
	CSWITCH
)

type Notification struct {
	Type    NotificationType
	Client  *Client
	Output  utils.Output
	Command string
}
