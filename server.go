package main

import (
	"GoReverSH/config"
	"GoReverSH/pkgserver"
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"
)

func main() {
	config.InitConfig()
	fmt.Println(config.Config)

	log.SetFlags(log.Lshortfile)
	//自分のIPアドレスと指定
	host := flag.String("host", "127.0.0.1", "hostIP")
	port := flag.String("port", "8000", "server port")
	flag.Parse()

	fmt.Println(*host, *port)

	grsh := pkgserver.NewGoReverSH(*host, *port)
	if grsh == nil || reflect.ValueOf(grsh).IsNil() {
		log.Fatalln(errors.New("GoReverSh constructor error"))
	}
	err := grsh.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
