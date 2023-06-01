package web

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"commit-log/appendlog"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	defer req.Body.Close()

	healthCheck(rec, req)

	assertJson(t, rec, http.StatusOK, map[string]string{"status": "ok"})
}

func TestAppend(t *testing.T) {
	t.Run("correct input", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", buildInput(t))
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		appendToLog(log)(rec, req)

		assertJson(t, rec, http.StatusOK, map[string]int{"offset": 0})
	})

	t.Run("double append", func(t *testing.T) {
		log := appendlog.NewAppendLog()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", buildInput(t))
		defer req.Body.Close()

		appendToLog(log)(rec, req)

		rec2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/", buildInput(t))
		defer req2.Body.Close()
		appendToLog(log)(rec2, req2)

		assertJson(t, rec, http.StatusOK, map[string]int{"offset": 0})
		assertJson(t, rec2, http.StatusOK, map[string]int{"offset": 1})
	})

	t.Run("invalid method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		appendToLog(log)(rec, req)

		assertJson(t, rec, http.StatusNotFound, map[string]string{"error": "method not found"})
	})

	t.Run("no input", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		appendToLog(log)(rec, req)

		assertJson(t, rec, http.StatusBadRequest, map[string]string{"error": "empty request provided"})
	})
}

func TestRead(t *testing.T) {
	t.Run("empty log", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/read/0", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		readFromLog(log)(rec, req)

		assertJson(t, rec, http.StatusBadRequest, map[string]string{"error": "can't find offset: 0"})
	})

	t.Run("invalid method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/read/0", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		readFromLog(log)(rec, req)

		assertJson(t, rec, http.StatusNotFound, map[string]string{"error": "method not found"})
	})

	t.Run("too big offset", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/read/123", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		log.Append(strings.NewReader(`{"foo": "bar"}`))

		readFromLog(log)(rec, req)

		assertJson(t, rec, http.StatusBadRequest, map[string]string{"error": "can't find offset: 123"})
	})

	t.Run("negative offset", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/read/-18", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		log.Append(strings.NewReader(`{"foo": "bar"}`))

		readFromLog(log)(rec, req)

		assertJson(t, rec, http.StatusBadRequest, map[string]string{"error": "can't find offset: -18"})
	})

	t.Run("correct offset", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/read/0", nil)
		defer req.Body.Close()

		log := appendlog.NewAppendLog()
		log.Append(strings.NewReader(`{"foo": "bar"}`))
		readFromLog(log)(rec, req)

		assertJson(t, rec, http.StatusOK, map[string]string{"foo": "bar"})
	})
}

func assertJson[T any](t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody T) {
	t.Helper()

	assert.Equal(t, expectedStatus, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header().Get("Content-type"))

	var body T
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, expectedBody, body)
}

func buildInput(t *testing.T) io.Reader {
	t.Helper()
	data, err := json.Marshal(map[string]string{
		"my-data":  "123",
		"my-data2": "asdf",
	})
	assert.NoError(t, err)
	return strings.NewReader(string(data))
}
