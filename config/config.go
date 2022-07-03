package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigList struct {
	DownloadOutDir string
	ScreenshotDir  string
	UploadDIr      string
}

var Config ConfigList

func InitConfig() {
	c := &ConfigList{}
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	Config = *c
}
