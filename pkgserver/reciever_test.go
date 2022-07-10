package pkgserver

import (
	"GoReverSH/pkgserver/mock"
	"GoReverSH/utils"
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
)

//通常通り
func TestNewReceiver(t *testing.T) {
	c := NewClient(mock.ConnMock{}, "test")
	ch := make(chan Notification)

	r := NewReceiver(c, ch, &sync.Mutex{})

	if r == nil {
		t.Errorf("NewReceiver error. got %v\n", r)
	}

}

func TestWaitMessage(t *testing.T) {
	//server, client := net.Pipe()
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	clientConn, err := net.Dial("tcp", ":3000")
	if err != nil {
		t.Error(err)
	}
	defer clientConn.Close()

	tests := []struct {
		name          string
		writeData     []byte
		errorExpected bool
	}{
		{"simple out", genOutput(utils.Output{}), false},
		{"Not the specified output", []byte("aaaa"), true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				defer conn.Close()

				c := NewClient(conn, "test")

				ch := make(chan Notification)
				r := NewReceiver(c, ch, &sync.Mutex{})
				if r == nil {
					t.Errorf("NewReceiver error. got %v\n", r)
				}
				ctx := context.Background()

				go func() {
					err := r.WaitMessage(ctx)

					//エラーがある、かつ期待通り
					if err != nil && tt.errorExpected {
						t.Errorf("wantError %v, got %v\n", tt.errorExpected, err)
					}
				}()

				//test
				//output形式とそうじゃない形式のテスト
				//出力関数を作成
				_, err = clientConn.Write(tt.writeData)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func genOutput(output utils.Output) []byte {
	data, _ := json.Marshal(output)
	return data
}
