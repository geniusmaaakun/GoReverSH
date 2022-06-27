package server

type NotificationType int

const (
	JOIN NotificationType = iota
	DEFECT
	CHANGE_DIR
	UPLOAD
	DOWNLOAD
	SCREEN_SHOT
	CLEAN
)

type Notification struct {
	Type     NotificationType
	Client   *Client
	Commands []string
}
