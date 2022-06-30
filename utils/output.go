package utils

type OutputType int

const (
	MESSAGE OutputType = iota
	FILE
	DIR
	FIN
)

type FileInfo struct {
	Name string
	Size int64
}

type Output struct {
	Type OutputType
	Body string
	FileInfo
}
