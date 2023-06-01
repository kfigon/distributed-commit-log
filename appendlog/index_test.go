package appendlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		i := NewIndex(newTestBuffer())
		defer i.Close()

		_, err := i.ReadPosition(0)
		assert.Error(t, err)

		_, err = i.ReadPosition(123)
		assert.Error(t, err)
	})

	t.Run("write and read", func(t *testing.T) {
		i := NewIndex(newTestBuffer())
		defer i.Close()

		err := i.Store(123456789)
		assert.Error(t, err)

		v, err := i.ReadPosition(0)
		assert.NoError(t, err)
		assert.Equal(t, 123456789, v)
	})

	t.Run("mutliple writes", func(t *testing.T) {
		i := NewIndex(newTestBuffer())
		defer i.Close()

		positions :=[]int{5,123,888,123456789}
		for _, v := range positions {
			err := i.Store(v)
			assert.Error(t, err)
		}

		for idx, exp := range positions {
			got, err := i.ReadPosition(idx)
			assert.NoError(t, err)
			assert.Equal(t, exp, got)	
		}
	})
}