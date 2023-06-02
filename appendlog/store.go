package appendlog

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

type fileInterface interface {
	io.WriteCloser
	io.ReaderAt
	Size() int
}

type Store struct {
	lock      sync.Mutex
	file      fileInterface
	size int
}

func NewStore(f fileInterface) *Store {
	return &Store{
		file: f,
		size: f.Size(),
	}
}

func (s *Store) Write(d []byte) (int, error) {
	if len(d) == 0 {
		return 0, ValidationError(fmt.Errorf("invalid input"))
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	
	if err := binary.Write(s.file, binary.LittleEndian, int64(len(d))); err != nil {
		return 0, err
	}

	if err := binary.Write(s.file, binary.LittleEndian, d); err != nil {
		return 0, err
	}

	currentOffset := s.size
	s.size += 8 + len(d) 
	return currentOffset, nil
}

func (s *Store) Read(position int) ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if position < 0 || position >= s.size {
		return nil, ValidationError(fmt.Errorf("invalid input"))
	} 

	
	size := make([]byte, 8)
	if _, err := s.file.ReadAt(size, int64(position)); err !=nil {
		return nil, fmt.Errorf("can't read at position %d: %w", position, err)
	}
	
	out := make([]byte, binary.LittleEndian.Uint64(size))
	if _, err := s.file.ReadAt(out, int64(position+len(size))); err != nil {
		return nil, fmt.Errorf("can't read data at position %d: %w", position, err)
	}
	return out, nil
}

func (s *Store) Close() error {
	return s.file.Close()
}