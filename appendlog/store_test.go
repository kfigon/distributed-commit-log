package appendlog

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Run("write null data", func(t *testing.T) {
		s := NewStore(newTestBuffer())
		defer s.Close()

		_, err := s.Write(nil)
		assert.Error(t, err)
	})

	t.Run("read negative data", func(t *testing.T) {
		s := NewStore(newTestBuffer())
		defer s.Close()

		_, err := s.Read(-1)
		assert.Error(t, err)
	})

	t.Run("read too big offset", func(t *testing.T) {
		s := NewStore(newTestBuffer())
		defer s.Close()

		_, err := s.Read(12345)
		assert.Error(t, err)
	})

	t.Run("read from empty store", func(t *testing.T) {
		s := NewStore(newTestBuffer())
		defer s.Close()

		got, err := s.Read(0)
		assert.Error(t, err)
		assert.Empty(t, got)
	})

	t.Run("write read to empty store", func(t *testing.T) {
		s := NewStore(newTestBuffer())
		defer s.Close()

		input := []byte(`{"name":"foo", "val": 123}`)
		pos, err := s.Write(input)
		assert.NoError(t, err)

		got, err := s.Read(pos)
		assert.NoError(t, err)
		assert.Equal(t, input, got)
	})

	t.Run("read from existing store", func(t *testing.T) {
		input := []byte(`{"data":"bar"}`)
		buf := newTestBuffer()
		binary.Write(buf, binary.LittleEndian, int64(len(input)))
		binary.Write(buf, binary.LittleEndian, input)

		s := NewStore(buf)
		defer s.Close()

		got, err := s.Read(0)
		assert.NoError(t, err)
		assert.Equal(t, input, got)
	})

	t.Run("write read from existing store", func(t *testing.T) {
		input := []byte(`{"data":"bar"}`)
		buf := newTestBuffer()
		binary.Write(buf, binary.LittleEndian, int64(len(input)))
		binary.Write(buf, binary.LittleEndian, input)

		s := NewStore(buf)
		defer s.Close()

		input2 := []byte(`{"name":"foo", "val": 123}`)
		pos, err := s.Write(input2)
		assert.NoError(t, err)

		got, err := s.Read(0)
		assert.NoError(t, err)
		assert.Equal(t, input, got)

		got2, err := s.Read(pos)
		assert.NoError(t, err)
		assert.Equal(t, input2, got2)
	})
}