package server

import (
	"reflect"
	"testing"
)

func TestNewReverceSH(t *testing.T) {
	grsh := NewGoReverSH("192.168.56.1", "8080")
	if grsh == nil || reflect.ValueOf(grsh).IsNil() {
		t.Log(grsh)
		t.Error("GoReverSh constructor error")
	}
}
