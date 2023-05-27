package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	const port = 8080

	theLog := &appendLog{}

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/append", loggingMiddleware(appendToLog(theLog)))
	http.HandleFunc("/read", loggingMiddleware(readFromLog(theLog)))

	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJson(w, map[string]string{
		"status": "ok",
	})
}

type record struct {
	data   []byte
	offset int
}

type appendLog struct {
	lock    sync.Mutex
	records []record
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w,r)
		log.Printf("%v %v, took %v\n", r.Method, r.URL, time.Since(start))
	}
}

func appendToLog(l *appendLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, fmt.Errorf("method not found"), http.StatusNotFound)
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, fmt.Errorf("error reading request: %w", err), http.StatusBadRequest)
			return
		} else if len(bytes) == 0 {
			writeError(w, fmt.Errorf("empty request provided"), http.StatusBadRequest)
			return
		}

		l.lock.Lock()
		defer l.lock.Unlock()
		rec := record{
			data:   bytes,
			offset: len(l.records),
		}
		l.records = append(l.records, rec)

		writeJson(w, map[string]int{
			"offset": rec.offset,
		})
	}
}

func readFromLog(l *appendLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, fmt.Errorf("method not found"), http.StatusNotFound)
			return
		}

		rawId := strings.TrimPrefix(r.URL.Path, "/")
		offset, err := strconv.Atoi(rawId)
		if err != nil {
			writeError(w, fmt.Errorf("can't parse path paramenter: %w", err), http.StatusBadRequest)
			return
		} else if offset >= len(l.records) {
			writeError(w, fmt.Errorf("too big offset provided: %d", offset), http.StatusBadRequest)
			return
		}

		l.lock.Lock()
		defer l.lock.Unlock()
		w.Header().Set("Content-type", "application/json")
		_, err = w.Write(l.records[offset].data)
		if err != nil {
			writeError(w, fmt.Errorf("error serialising user data at offset %d: %w", offset, err), http.StatusInternalServerError)
			return
		}
	}
}

func writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func writeJson(w http.ResponseWriter, body any) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(body)
}
