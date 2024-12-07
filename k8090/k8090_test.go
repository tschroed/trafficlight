package k8090

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	tl "github.com/tschroed/trafficlight"
)

func TestCsum(t *testing.T) {
	cases := []struct {
		val  uint8
		want uint8
	}{
		{
			val:  0x04 + 0x12 + 0x00,
			want: 0xea,
		},
		{
			val:  0x04 + 0x11 + 0x07,
			want: 0xe4,
		},
		{
			val:  0x04 + 0x11 + 0x03,
			want: 0xe8,
		},
		{
			val:  0x04 + 0x11 + 0x00,
			want: 0xeb,
		},
	}
	for _, tc := range cases {
		got := csum(tc.val)
		if got != tc.want {
			t.Errorf("csum(0x%x) = 0x%x, want 0x%x", tc.val, got, tc.want)
		}
	}
}

func TestSet(t *testing.T) {
	var b bytes.Buffer
	k := new(&b)
	cases := []struct {
		word uint8
	}{
		{word: 0},
		{word: 7},
		{word: 255},
	}
	all := uint8(0xff)
	for _, tc := range cases {
		want := []byte{
			0x04, 0x12, all, 0x00, 0x00, csum(0x04 + 0x12 + all), 0x0f,
			0x04, 0x11, tc.word, 0x00, 0x00, csum(0x04 + 0x11 + tc.word), 0x0f,
		}
		if err := k.Set(tc.word); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(b.Bytes(), want) {
			t.Errorf("k.Set(%v) = %v, want %v", tc.word, b.Bytes(), want)
		}
		b.Reset()
	}
}

func TestNew(t *testing.T) {
	// Do this to enforce interface implementation.
	var k tl.TrafficLight
	var err error
	fn := "/file/does/not/exist"
	if _, err = New(fn); err == nil {
		t.Errorf("New(%v) got nil, wanted error", fn)
	}
	fn = "/dev/null"
	if k, err = New(fn); err != nil {
		t.Errorf("New(%v) got %v, wanted nil", fn, err)
	}
	if _, ok := k.(tl.TrafficLight); !ok {
		t.Error("k is not a TrafficLight")
	}
}

type errWriter struct{
	failAfter int
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	if w.failAfter <= 0 {
		return -1, errors.New("This is an error.")
	}
	w.failAfter--
	return len(p), nil
}

func TestSetError(t *testing.T) {
	k := new(&errWriter{})
	if err := k.Set(3); err == nil {
		t.Error("Expected error, got none.")
	}
	k = new(&errWriter{failAfter: 1})
	if err := k.Set(3); err == nil {
		t.Error("Expected error, got none.")
	}
	k = new(&errWriter{failAfter: 2})
	if err := k.Set(3); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
