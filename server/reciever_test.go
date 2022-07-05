package server

import (
	"GoReverSH/utils"
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
)

//通常通り
func TestNewReceiver(t *testing.T) {
	server, client := net.Pipe()
	ch := make(chan Notification)

	r := NewReceiver(client, "test", ch, &sync.Mutex{})

	if r == nil {
		t.Errorf("NewReceiver error. got %v\n", r)
	}

	server.Close()

	client.Close()
}

func TestWaitMessage(t *testing.T) {
	server, client := net.Pipe()
	ch := make(chan Notification)

	r := NewReceiver(client, "test", ch, &sync.Mutex{})

	if r == nil {
		t.Errorf("NewReceiver error. got %v\n", r)
	}

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
			client.Write(tt.writeData)
		})
	}

	server.Close()
	client.Close()
}

func genOutput(output utils.Output) []byte {
	data, _ := json.Marshal(output)
	return data
}
