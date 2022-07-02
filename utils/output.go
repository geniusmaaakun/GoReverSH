package utils

type OutputType int

const (
	MESSAGE OutputType = iota
	FILE
	DIR
)

type Output struct {
	Type    OutputType
	Message []byte
	FileInfo
}

type FileInfo struct {
	Name string
	Body []byte
	Size int64
}
