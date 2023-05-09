package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"search-service/models"
)

// readJSON tries to read the body of a request and converts it into JSON
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes int64 = 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	var statusCode int

	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusBadRequest
	}

	payload := models.JsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return writeJSON(w, statusCode, payload)
}
