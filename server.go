package main

import (
	"GoReverSH/config"
	"GoReverSH/server"
	"GoReverSH/utils"
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
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
	return &GoReverSH{signalCh: sigCH, op: op, lock: &lock}
}

func (grsh *GoReverSH) FreeAllClientMap() bool {
	for _, c := range grsh.Observer.State.ClientMap {
		c.Conn.Close()
	}
	return true
}

func (grsh *GoReverSH) waitSignal(cancel context.CancelFunc, listner net.Listener) {
	for {
		select {
		//ここですべて終了
		case <-grsh.signalCh:
			//CleanUp
			fmt.Println("\n", "Cleanup")
			cancel()
			grsh.FreeAllClientMap()
			listner.Close()
			os.Exit(1)
		}
	}
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
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", net.JoinHostPort(grsh.op.host, grsh.op.port), &tlsConfig)
	if err != nil {
		return fmt.Errorf("Listen error: %+w\n", err)
	}
	defer listener.Close()
	fmt.Printf("Server running at %s:%s\n", grsh.op.host, grsh.op.port)

	//context
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	var channel = make(chan server.Notification)
	//通知を入れるチャネル

	//TODO NewObserver, NewExecuter
	//TODO lockを渡す
	//通知を受け取る
	state := server.State{ClientMap: make(map[string]*server.Client)}
	grsh.Observer = &server.Observer{State: state, Subject: channel, PromptViewFlag: false, Lock: grsh.lock}
	//実行コマンドを受け取る
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	grsh.Executer = &server.Executer{Scanner: scanner, Observer: channel}

	//waitSignal
	//CTRL + C
	go grsh.waitSignal(cancel, listener)

	//通知を待つ
	go grsh.Observer.WaitNotice(context)

	//コマンドを待つ
	go grsh.Executer.WaitCommand(context)

	//クライアントを待つ
	err = grsh.waitClient(context, listener, channel)
	if err != nil {
		return err
	}

	return nil
}

func (grsh *GoReverSH) waitClient(ctx context.Context, listener net.Listener, channel chan server.Notification) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		//長い名前の場合、最後まで読むようにする?
		cnameBuff := make([]byte, 1024)
		n, err := conn.Read(cnameBuff)
		if err != nil {
			return err
		}
		client := server.NewClient(conn, string(cnameBuff[:n]))

		//受信を待つ
		//read & join
		receiver := server.Receiver{Client: client, Observer: channel, Lock: grsh.lock}

		go receiver.Start(ctx)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	//自分のIPアドレスと指定
	host := flag.String("host", "127.0.0.1", "hostIP")
	port := flag.String("port", "8000", "server port")
	flag.Parse()

	fmt.Println(*host, *port)

	config.InitConfig()
	fmt.Println(config.Config)

	grsh := NewGoReverSH(*host, *port)
	err := grsh.run()
	if err != nil {
		log.Fatalln(err)
	}
}
