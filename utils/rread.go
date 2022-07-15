package utils

import (
	"errors"
	"io"
	"log"
)

//繰り返し読み込む関数。効率いい
func RRead(reader io.Reader) ([]byte, error) {
	var content []byte
	var err error
	buff := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buff)
		if err != nil {
			if n == 0 || errors.Is(err, io.EOF) {
				err = nil
				break
			}
			log.Println(err)
			break
		}
		//content + bufの中身を一時的に保存。
		tmp := make([]byte, 0, size+n)
		tmp = append(content[:size], buff[:n]...)
		content = tmp
		size += n
		if n < 1024 {
			break
		}
	}

	return content[:size], err
}
