package web

import (
	"commit-log/appendlog"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func NewServer(theLog *appendlog.AppendLog) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, newHttpError(fmt.Errorf("not found"), http.StatusNotFound))
	})

	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/append", loggingMiddleware(appendToLog(theLog)))
	mux.HandleFunc("/read/", loggingMiddleware(readFromLog(theLog)))

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJson(w, map[string]string{
		"status": "ok",
	})
}

type httpError struct {
	error
	status int
}

func newHttpError(err error, status int) httpError {
	return httpError{err, status}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%v %v, took %v\n", r.Method, r.URL, time.Since(start))
	}
}

func withHttpMethod(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			writeError(w, newHttpError(fmt.Errorf("method not found"), http.StatusNotFound))
			return
		}
		next(w, r)
	}
}

func appendToLog(l *appendlog.AppendLog) http.HandlerFunc {
	return withHttpMethod(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		offset, err := l.Append(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJson(w, map[string]int{
			"offset": offset,
		})
	})
}

func readFromLog(l *appendlog.AppendLog) http.HandlerFunc {
	return withHttpMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		offset := strings.TrimPrefix(r.URL.Path, "/read/")
		bytes, err := l.Read(offset)
		if err != nil {
			writeError(w, err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		_, err = w.Write(bytes)
		if err != nil {
			writeError(w, newHttpError(fmt.Errorf("error serialising user data at offset %s: %w", offset, err), http.StatusInternalServerError))
			return
		}
	})
}

func writeError(w http.ResponseWriter, err error) {
	status := func() int {
		var httpErr httpError
		if errors.As(err, &httpErr) {
			return httpErr.status
		}
		var valErr appendlog.ValidationError
		if errors.As(err, &valErr) {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status())
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func writeJson(w http.ResponseWriter, body any) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(body)
}
