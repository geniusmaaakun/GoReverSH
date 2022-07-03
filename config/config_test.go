package config

import "testing"

func TestNewConfig(t *testing.T) {
	InitConfig()
	t.Log(Config)
}
