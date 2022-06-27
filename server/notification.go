package server

type NotificationType int

const (
	JOIN NotificationType = iota
	DEFECT
	MESSAGE
	COMMAND
	CHANGE_DIR
	UPLOAD
	DOWNLOAD
	SCREEN_SHOT
	CLEAN
)

type Notification struct {
	Type    NotificationType
	Client  *Client
	Message string
	Command string
}
