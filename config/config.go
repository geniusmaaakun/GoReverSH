package config

import (
	_ "embed"
	"encoding/json"
	"log"
)

type ConfigList struct {
	DownloadOutDir string
	ScreenshotDir  string
	UploadDIr      string
}

//go:embed config.json
var jsonConfig []byte

var Config ConfigList

func InitConfig() {
	c := &ConfigList{}
	/*
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatalln(err)
		}
		err = json.NewDecoder(file).Decode(&c)
	*/
	err := json.Unmarshal(jsonConfig, &c)
	if err != nil {
		log.Fatalln(err)
	}
	Config = *c
}
