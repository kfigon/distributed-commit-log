package web

import (
	"commit-log/appendlog"
	"fmt"
	"net/http"
	"strings"
)

func NewServer(theLog *appendlog.AppendLog) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, newHttpError(fmt.Errorf("not found"), http.StatusNotFound))
	})

	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/append", appendToLog(theLog))
	mux.HandleFunc("/read/", readFromLog(theLog))

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJson(w, map[string]string{
		"status": "ok",
	})
}

func appendToLog(l *appendlog.AppendLog) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		offset, err := l.Append(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJson(w, map[string]int{
			"offset": offset,
		})
	}

	return panicRecoveryMiddleware(loggingMiddleware(httpMethodMiddleware(http.MethodPost, fn)))
}

func readFromLog(l *appendlog.AppendLog) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
	}

	return panicRecoveryMiddleware(loggingMiddleware(httpMethodMiddleware(http.MethodGet, fn)))
}
