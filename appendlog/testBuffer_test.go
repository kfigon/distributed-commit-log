package appendlog

import "bytes"

type TestBuffer struct {
	bytes.Buffer
}

func newTestBuffer() *TestBuffer {
	return &TestBuffer{*bytes.NewBuffer(nil)}
}

func (t *TestBuffer) Close() error {
	return nil
}