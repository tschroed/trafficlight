package k8090

import (
	"io"
	"os"
)

type K8090 struct {
	port io.Writer
}

func csum(n uint8) uint8 {
	return (n ^ 255) + 1
}

func (k *K8090) Set(word uint8) error {
	var w uint8
	w = 0xff
	c := csum(0x04 + 0x12 + w)
	b := []byte{0x04, 0x12, w, 0x00, 0x00, c, 0x0f}
	if _, err := k.port.Write(b); err != nil {
		return err
	}
	c = csum(0x04 + 0x11 + word)
	b = []byte{0x04, 0x11, word, 0x00, 0x00, c, 0x0f}
	if _, err := k.port.Write(b); err != nil {
		return err
	}
	return nil
}

func new(w io.Writer) *K8090 {
	return &K8090{port: w}
}

func New(tty string) (*K8090, error) {
	f, err := os.OpenFile(tty, os.O_RDWR, 0)
	return new(f), err
}
