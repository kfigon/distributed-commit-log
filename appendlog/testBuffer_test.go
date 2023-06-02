package appendlog

import (
	"bytes"
	"io"
)

type TestBuffer struct {
	bytes.Buffer
}

func newTestBuffer() *TestBuffer {
	return &TestBuffer{*bytes.NewBuffer(nil)}
}

func (t *TestBuffer) Close() error {
	return nil
}

func (t *TestBuffer) ReadAt(d []byte, offset int64) (int, error) {
	buf2 := bytes.NewBuffer(t.Bytes()) // make a copy so it's not consuming the data

	out, err := io.ReadAll(buf2)
	if err != nil {
		return 0, err
	}
	n := copy(d, out[offset:])

	return n, nil
}

func (t *TestBuffer) Size() int {
	return t.Len()
}