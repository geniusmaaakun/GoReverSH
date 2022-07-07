package mock

import (
	"net"
	"time"
)

type Addr struct {
}

func (Addr) Network() string {
	return "test"
}
func (Addr) String() string {
	return "test"
}

type ConnMock struct {
}

func (ConnMock) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (ConnMock) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (ConnMock) Close() error {
	return nil
}

func (ConnMock) LocalAddr() net.Addr {
	return Addr{}
}

// RemoteAddr returns the remote network address, if known.
func (ConnMock) RemoteAddr() net.Addr {
	return Addr{}
}

func (ConnMock) SetDeadline(t time.Time) error {
	return nil
}

func (ConnMock) SetReadDeadline(t time.Time) error {
	return nil
}

func (ConnMock) SetWriteDeadline(t time.Time) error {
	return nil
}
