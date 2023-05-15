package main

import (
	"broker-service/models"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HandleSubmittion is the single point of entry to the broker service
// that handles all requests based on the action specified
func HandleSubmittion(w http.ResponseWriter, r *http.Request) {
	var (
		item    any
		service string
	)

	requestPayload := new(models.RequestPayload)

	err := readJSON(w, r, requestPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "search":
		item = &requestPayload.Search
		service = "http://search-service/search-entry"
	case "get-pdf":
		item = &requestPayload.Search
		service = "http://search-service/get-pdf"
	case "process-text":
		item = &requestPayload.NLP
		service = "http://nlp-service/process-text"
	default:
		errorJSON(w, errors.New("unkown action"))
		return
	}

	callService(service, w, item)
}

// callService calls the microservice according to the provided service with the provided item
// and writes the response to the provided http.ResponseWriter
func callService(service string, w http.ResponseWriter, item any) {
	jsonData, _ := json.Marshal(item)

	response, err := http.Post(service, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, response.Body)
	if err != nil {
		errorJSON(w, err)
	}
}
