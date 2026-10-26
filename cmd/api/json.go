package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iliabhz/gosocial/internal/store"
)

const (
	_  = iota             // ignore first value
	KB = 1 << (10 * iota) // 1024
	MB                    // 1048576
	GB                    // 1073741824
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.MarshalWrite(w, data)
}

func writeDataJSON(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, map[string]any{"data": data})
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1 * MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := jsontext.NewDecoder(r.Body)
	return json.UnmarshalDecode(dec, data, json.RejectUnknownMembers(true))
}

func writeJSONError(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := errors.AsType[*HTTPError](err); ok {
		writeJSON(w, e.status, map[string]any{"error": e.msg})
		return
	}

	if e, ok := errors.AsType[validator.ValidationErrors](err); ok {
		writeJSON(w, http.StatusBadRequest, map[string]any{"errors": validationMessages(e)})
		return
	}

	switch {
	case errors.Is(err, store.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Resource not found"})
	//case errors.Is(err, ErrUnauthorized):
	//	writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	default:
		slog.Error("internal server error", "method", r.Method, "err", err, "path", r.URL.Path)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
}
