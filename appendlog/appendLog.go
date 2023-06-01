package appendlog

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

type Record struct {
	data   []byte
	offset int
}

type AppendLog struct {
	lock    sync.Mutex
	records []Record
}

func NewAppendLog() *AppendLog {
	return &AppendLog{}
}

type ValidationError error

func (a *AppendLog) Append(data io.Reader) (int, error) {
	bytes, err := io.ReadAll(data)
	if err != nil {
		return 0, ValidationError(fmt.Errorf("error reading request: %w", err))
	} else if len(bytes) == 0 {
		return 0, ValidationError(fmt.Errorf("empty request provided"))
	}

	a.lock.Lock()
	defer a.lock.Unlock()
	rec := Record{
		data:   bytes,
		offset: len(a.records),
	}
	a.records = append(a.records, rec)
	return rec.offset, nil
}

func (a *AppendLog) Read(rawId string) ([]byte, error) {
	offset, err := strconv.Atoi(rawId)
	if err != nil {
		return nil, ValidationError(fmt.Errorf("can't parse path paramenter: %w", err))
	} else if offset >= len(a.records) || offset < 0 {
		return nil, ValidationError(fmt.Errorf("can't find offset: %d", offset))
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.records[offset].data, nil
}
