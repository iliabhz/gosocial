package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"net/http"
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

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1 * MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := jsontext.NewDecoder(r.Body)
	return json.UnmarshalDecode(dec, data, json.RejectUnknownMembers(true))
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type errStruct struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &errStruct{Error: message})
}
