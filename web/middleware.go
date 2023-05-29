package web

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%v %v, took %v\n", r.Method, r.URL, time.Since(start))
	}
}

func panicRecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if out := recover(); out != nil {
				log.Printf("got panic on %v %v: %v", r.Method, r.URL, out)
			}
		}()

		next(w, r)
	}
}

func httpMethodMiddleware(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			writeError(w, newHttpError(fmt.Errorf("method not found"), http.StatusNotFound))
			return
		}
		next(w, r)
	}
}
