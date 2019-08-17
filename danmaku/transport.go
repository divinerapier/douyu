package danmaku

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Transport struct {
	readWriter io.ReadWriter
}

func NewTransport(network, address string) (*Transport, error) {
	conn, err := net.DialTimeout(network, address, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return &Transport{
		readWriter: conn,
	}, nil
}

func (t *Transport) ReadFull(buf []byte) error {
	n, err := io.ReadFull(t.readWriter, buf)
	if err == io.EOF {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("not reading enough data. expect: %d, actual: %d", len(buf), n)
	}
	return nil
}

func (t *Transport) ReadSize(size int) ([]byte, error) {
	buf := make([]byte, size)
	if err := t.ReadFull(buf); err != nil {
		buf = buf[:0:0]
		return nil, err
	}
	return buf, nil
}

func (t *Transport) Write(p []byte) error {
	n, err := t.readWriter.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return fmt.Errorf("not writing enough. expect: %d, actual: %d", len(p), n)
	}
	return nil
}
