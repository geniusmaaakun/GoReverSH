package config

import "testing"

func TestNewConfig(t *testing.T) {
	InitConfig()
	if Config.DownloadOutDir == "" || Config.ScreenshotDir == "" || Config.UploadDIr == "" {
		t.Errorf("configlist params not setting. got: %+v\n", Config)
	}
}
