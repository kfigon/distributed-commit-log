package main

import (
	"fmt"
	"io"
	"net/http"
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

func (a *appendLog) append(data io.ReadCloser) (int, error) {
	bytes, err := io.ReadAll(data)
	if err != nil {
		return 0, httpError{fmt.Errorf("error reading request: %w", err), http.StatusBadRequest}
	} else if len(bytes) == 0 {
		return 0, httpError{fmt.Errorf("empty request provided"), http.StatusBadRequest}
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
		return nil, httpError{fmt.Errorf("can't parse path paramenter: %w", err), http.StatusBadRequest}
	} else if offset >= len(a.records) {
		return nil, httpError{fmt.Errorf("too big offset provided: %d", offset), http.StatusBadRequest}
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.records[offset].data, nil
}
