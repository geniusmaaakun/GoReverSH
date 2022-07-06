package server

import (
	"sync"
	"testing"
)

func TestNewObserver(t *testing.T) {
	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	if o == nil {
		t.Errorf("NewObserver error. got %v\n", o)
	}
}

//join
func TestJoinClient(t *testing.T) {

	tests := []struct {
		name    string
		clients []*Client
		want    int
	}{
		{"scusses", []*Client{{nil, "test", "test"}, {nil, "test2", "test2"}}, 2},
		{"nil check", []*Client{nil}, 0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ch := make(chan Notification)
			o := NewObserver(ch, &sync.Mutex{})

			for _, c := range tt.clients {
				t.Log(c)
				o.joinClient(c)
			}
			if len(o.State.ClientMap) != tt.want {
				t.Errorf("joinClient length error: got %v, want %v\n", len(o.State.ClientMap), tt.want)
			}
		})
	}
}

//detect
func TestDefectClient(t *testing.T) {
	//最初に数人追加しておく
	//退出通知を送る
}

//upload
func TestUpload(t *testing.T) {

}

//createfile
func TestCreatefile(t *testing.T) {
	//testdataに作成し、削除
}

//clist
func TestClist(t *testing.T) {
	//出力テスト
}

//cswitch
func TestCSwitch(t *testing.T) {

}

//FreeMap
func TestFreeMap(t *testing.T) {

}

//waitnotice
func TestWaitNotice(t *testing.T) {
	//senderをmock
}

//printPrompt
func TestPrintPrompt(t *testing.T) {

}
