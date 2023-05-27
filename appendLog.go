package main

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

type record struct {
	data   []byte
	offset int
}

type appendLog struct {
	lock    sync.Mutex
	records []record
}

type validationError error

func (a *appendLog) append(data io.ReadCloser) (int, error) {
	bytes, err := io.ReadAll(data)
	if err != nil {
		return 0, validationError(fmt.Errorf("error reading request: %w", err))
	} else if len(bytes) == 0 {
		return 0, validationError(fmt.Errorf("empty request provided"))
	}

	a.lock.Lock()
	defer a.lock.Unlock()
	rec := record{
		data:   bytes,
		offset: len(a.records),
	}
	a.records = append(a.records, rec)
	return rec.offset, nil
}

func (a *appendLog) read(rawId string) ([]byte, error) {
	offset, err := strconv.Atoi(rawId)
	if err != nil {
		return nil, validationError(fmt.Errorf("can't parse path paramenter: %w", err))
	} else if offset >= len(a.records) {
		return nil, validationError(fmt.Errorf("can't find offset: %d", offset))
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.records[offset].data, nil
}
