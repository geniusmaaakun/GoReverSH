package pkgserver

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
)

//通常通り

func TestNewExecuter(t *testing.T) {
	ch := make(chan Notification)
	e := NewExecuter(ch)
	if e == nil {
		t.Errorf("NewExecuter error. got %v\n", e)
	}
}

func TestWaitCommand(t *testing.T) {
	ch := make(chan Notification)
	e := NewExecuter(ch)

	_, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdin := os.Stdin
	os.Stdin = w
	defer func() {
		os.Stdin = stdin
	}()

	tests := []struct {
		name    string
		command string
		want    error
	}{
		{"clist", "clist", nil},
		{"cswitch", "cswitch", nil},
		{"upload", "upload", nil},
		{"download", "download", nil},
		{"screenshot", "screenshot", nil},
		{"clean", "clean_go_reversh", nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			go func() {
				err := e.WaitCommand(ctx)
				if err != tt.want {
					t.Errorf("want %v, got %v\n", tt.want, err)
				}
			}()

			w.WriteString(tt.command)
		})
	}
}

func captureStdin(t *testing.T, f func() error, s string) error {
	t.Helper()
	_, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdin := os.Stdin
	os.Stdin = w
	defer func() {
		os.Stdin = stdin
	}()

	return f()
}

func captureStdout(t *testing.T, f func() error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	os.Stdout = w

	err = f()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), err
}
