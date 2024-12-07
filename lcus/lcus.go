package lcus

import (
	"io"
	"os"
	"time"

	tl "github.com/tschroed/trafficlight"
)

type LCUS struct {
	port io.Writer
}

func (l *LCUS) Set(word uint8) error {
	// Everything off to start with.
	b := [][]byte{
		{0xa0, 0x01, 0x00, 0xa1},
		{0xa0, 0x02, 0x00, 0xa2},
		{0xa0, 0x03, 0x00, 0xa3},
		// To make the reset cycle faster, we'll assume that only the first three
		// are in use. Uncomment below if additional relays are used.
		//		{0xa0, 0x04, 0x00, 0xa4},
		//	        {0xa0, 0x05, 0x00, 0xa5},
		//		{0xa0, 0x06, 0x00, 0xa6},
		//	        {0xa0, 0x07, 0x00, 0xa7},
		//		{0xa0, 0x08, 0x00, 0xa8},
	}
	if word&tl.RED == tl.RED {
		b = append(b, []byte{0xa0, 0x01, 0x01, 0xa2})
	}
	if word&tl.GREEN == tl.GREEN {
		b = append(b, []byte{0xa0, 0x02, 0x01, 0xa3})
	}
	if word&tl.AMBER == tl.AMBER {
		b = append(b, []byte{0xa0, 0x03, 0x01, 0xa4})
	}
	for _, wb := range b {
		time.Sleep(5 * time.Millisecond)
		if _, err := l.port.Write(wb); err != nil {
			return err
		}
	}
	return nil
}

func new(w io.Writer) *LCUS {
	return &LCUS{port: w}
}

func New(tty string) (*LCUS, error) {
	f, err := os.OpenFile(tty, os.O_RDWR, 0)
	return new(f), err
}
