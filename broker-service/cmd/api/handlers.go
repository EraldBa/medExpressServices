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
	requestPayload := new(models.RequestPayload)

	err := readJSON(w, r, requestPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "log":
		logItem(w, &requestPayload.Entry)
	case "search":
		searchItem(w, &requestPayload.Search)
	case "get-pdf":
		getPdf(w, &requestPayload.Search)
	default:
		errorJSON(w, errors.New("unkown action"))
	}
}

// logItem requests to log an entry to the mongodb through the search-service
// It writes to the http.ResponseWriter a jsonResponse indicating success or failure
func logItem(w http.ResponseWriter, entry *models.SearchEntry) {
	jsonData, _ := json.Marshal(entry)

	response, err := http.Post("http://search-service/log", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("status not accepted"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Entry logged successfully",
	}

	writeJSON(w, http.StatusAccepted, payload)
}

// searchItem searches for the  requested keyword through search-service
// It writes to the http.ResponseWriter the requested data
// or an error if the search-service failed to retrieve the data
func searchItem(w http.ResponseWriter, item *models.SearchQuery) {
	jsonData, _ := json.Marshal(item)

	response, err := http.Post("http://search-service/search-entry", "application/json", bytes.NewBuffer(jsonData))
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

// getPdf gets the requested pdf with the supplied pmcid from the search-service
// It writes to the http.ResponseWriter the requested pdf
// or an error if search-service failed to retrieve it
func getPdf(w http.ResponseWriter, item *models.SearchQuery) {
	jsonData, _ := json.Marshal(item)

	response, err := http.Post("http://search-service/get-pdf", "application/json", bytes.NewBuffer(jsonData))
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
