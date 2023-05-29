package web

import (
	"commit-log/appendlog"
	"encoding/json"
	"errors"
	"net/http"
)

type httpError struct {
	error
	status int
}

func newHttpError(err error, status int) httpError {
	return httpError{err, status}
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
