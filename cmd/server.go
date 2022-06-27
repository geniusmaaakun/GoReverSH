package main

import (
	"GoReverSH/cmd/server"
	"GoReverSH/utils"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
)

type GoReverSH struct {
	signalCh chan os.Signal
	op       Option
	Observer *server.Observer
	Executer *server.Executer
	lock     *sync.Mutex
}

type Option struct {
	host string
	port string
}

func NewGoReverSH(host, port string) *GoReverSH {
	sigCH := make(chan os.Signal, 1)
	signal.Notify(sigCH, os.Interrupt)
	op := Option{host, port}
	lock := sync.Mutex{}
	return &GoReverSH{signalCh: sigCH, op: op, Observer: &observer, Executer: &executer, lock: &lock}
}

func (grsh *GoReverSH) FreeAllClientMap() bool {
	for _, c := range grsh.state.clientMap {
		c.conn.Close()
	}
	return true
}

func (grsh *GoReverSH) PrintPrompt() {
	fmt.Printf("[GoReverSH@%s] >", grsh.state.connectingClient.name)
}

func (grsh *GoReverSH) run() error {
	certFile, keyFile, err := utils.GenClientCerts()
	if err != nil {
		return err
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	//startserver
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", net.JoinHostPort(grsh.op.host, grsh.op.port), &config)
	if err != nil {
		return fmt.Errorf("Listen error: %+w\n", err)
	}
	defer listener.Close()
	fmt.Printf("Server running at %s:%s\n", grsh.op.host, grsh.op.port)

	//context
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	//CTRL + C
	go func() {
		for {
			select {
			//ここですべて終了
			case <-grsh.signalCh:
				//CleanUp
				fmt.Println("Cleanup")
				cancel()
				grsh.FreeAllClientMap()
				listener.Close()
				os.Exit(1)
			}
		}
	}()

	var channel = make(chan server.Notification)
	//通知を入れるチャネル

	//通知を受け取る
	state := server.State{ClientMap: make(map[string]*server.Client)}
	grsh.Observer = &server.Observer{Sender: nil, State: &state, Subject: channel, PromptViewFlag: false}
	//実行コマンドを受け取る
	grsh.Executer = &server.Executer{Observer: channel}

	//通知を待つ
	go grsh.Observer.WaitNotice(context)

	//コマンドを待つ
	go grsh.Executer.WaitCommand(context)

	//クライアントを待つ
	grsh.waitClient(context, listener, channel)

	return nil
}

func (grsh *GoReverSH) waitClient(ctx context.Context, listener net.Listener, channel chan server.Notification) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		//最後まで読むようにする?
		cnameBuff := make([]byte, 1024)
		n, err := conn.Read(cnameBuff)
		if err != nil {
			return err
		}
		client := server.NewClient(conn, string(cnameBuff[:n]))

		//受信を待つ
		//read & join
		receiver := server.Receiver{Client: client, Observer: channel}

		go receiver.Start()
	}
}

func main() {

}
