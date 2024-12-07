package lcus

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	tl "github.com/tschroed/trafficlight"
)

func TestSet(t *testing.T) {
	var b bytes.Buffer
	l := new(&b)
	cases := []struct {
		word  uint8
		bytes []byte
	}{{
		word:  0,
		bytes: []byte{},
	}, {
		word:  2,
		bytes: []byte{0xa0, 0x02, 0x01, 0xa3},
	}, {
		word: 7,
		bytes: []byte{
			0xa0, 0x01, 0x01, 0xa2,
			0xa0, 0x02, 0x01, 0xa3,
			0xa0, 0x03, 0x01, 0xa4,
		},
	}, {
		word: 255,
		bytes: []byte{
			0xa0, 0x01, 0x01, 0xa2,
			0xa0, 0x02, 0x01, 0xa3,
			0xa0, 0x03, 0x01, 0xa4,
		},
	}}
	for _, tc := range cases {
		want := append([]byte{
			0xa0, 0x01, 0x00, 0xa1,
			0xa0, 0x02, 0x00, 0xa2,
			0xa0, 0x03, 0x00, 0xa3,
		}, tc.bytes...)
		if err := l.Set(tc.word); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(b.Bytes(), want) {
			t.Errorf("l.Set(%v) = %v, want %v", tc.word, b.Bytes(), want)
		}
		b.Reset()
	}
}

func TestNew(t *testing.T) {
	// Do this to enforce interface implementation.
	var l tl.TrafficLight
	var err error
	fn := "/file/does/not/exist"
	if _, err = New(fn); err == nil {
		t.Errorf("New(%v) got nil, wanted error", fn)
	}
	fn = "/dev/null"
	if l, err = New(fn); err != nil {
		t.Errorf("New(%v) got %v, wanted nil", fn, err)
	}
	if _, ok := l.(tl.TrafficLight); !ok {
		t.Error("k is not a TrafficLight")
	}
}

type errWriter struct{}

func (w *errWriter) Write(p []byte) (n int, err error) {
	return -1, errors.New("This is an error.")
}

func TestSetError(t *testing.T) {
	l := new(&errWriter{})
	if err := l.Set(3); err == nil {
		t.Error("Expected error, got none.")
	}
}
