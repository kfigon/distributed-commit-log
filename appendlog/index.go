package appendlog

import (
	"encoding/binary"
	"fmt"
	"io"
)

// bufio.Writer can reduce number of OS calls
// memory mapped file can also give better performance
type Index struct {
	file io.ReadWriteCloser
	positions []int // offset->position
}

func NewIndex(file io.ReadWriteCloser) *Index {
	return &Index{
		file: file,
		positions: nil,
	}
}

func (i *Index) Store(position int) error {
	i.positions = append(i.positions, position)
	return binary.Write(i.file, binary.LittleEndian, position)
}

func (i *Index) ReadPosition(offset int) (int, error) {
	if offset >= len(i.positions) || offset < 0 {
		return 0, ValidationError(fmt.Errorf("invalid offset %d", offset))
	}
	return i.positions[offset], nil
}

func (i *Index) Close() error {
	return i.file.Close()
}