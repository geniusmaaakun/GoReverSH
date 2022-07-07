package server

import (
	"GoReverSH/server/mock"
	"fmt"
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

	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])

	o.defectClient(clients[0])

	if _, ok := o.State.ClientMap[clients[0].Name]; ok {
		t.Errorf("client map error: %s is defect\n", clients[0].Name)
	}
}

//upload
//後回し
func TestUpload(t *testing.T) {

}

//createfile
func TestCreatefile(t *testing.T) {
	//testdataに作成し、削除
	//outdirをが引数で指定
}

//map 出力テストは順番が変わるため難しい
/*
//clist
func TestClist(t *testing.T) {
	want := `&{{} test test}
&{{} test1 test1}
&{{} test2 test2}
&{{} test3 test3}
`

	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])

	out, _ := captureStdout(t, func() error { o.execClientlist(); return nil })

	if out != want {
		t.Errorf("\nwant: %s,\ngot: %s\n", want, out)
	}
}
*/

//cswitch
func TestCSwitch(t *testing.T) {
	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])

	beforeClient := o.Sender.connectingClient
	o.execClientSwitch(Notification{Type: CSWITCH, Client: clients[1], Command: fmt.Sprintf("cswich %s", clients[1].Name)})
	afterClient := o.Sender.connectingClient
	if beforeClient == afterClient {
		t.Errorf("after client is %v, but got %v\n", clients[1], afterClient)
	}
}

//FreeMap
func TestFreeMap(t *testing.T) {
	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])
}

//printPrompt
func TestPrintPrompt(t *testing.T) {
	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])
}

//waitnotice
func TestWaitNotice(t *testing.T) {
	ch := make(chan Notification)
	o := NewObserver(ch, &sync.Mutex{})
	clients := []*Client{
		{mock.ConnMock{}, "test", "test"},
		{mock.ConnMock{}, "test1", "test1"},
		{mock.ConnMock{}, "test2", "test2"},
		{mock.ConnMock{}, "test3", "test3"},
	}
	o.joinClient(clients[0])
	o.joinClient(clients[1])
	o.joinClient(clients[2])
	o.joinClient(clients[3])

	//senderをmock
	//test case pattern all
}
