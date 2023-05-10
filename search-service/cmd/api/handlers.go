package main

import (
	"net/http"
	"search-service/data"
	"search-service/models"
)

// LogSearchEntry inserts a SeachEntry into the mongodb
// and writes a jsonResponse to the http.ResponseWriter
// that indicates whether the insertion was successful
func LogSearchEntry(w http.ResponseWriter, r *http.Request) {
	var logPayload models.SearchEntry

	err := readJSON(w, r, &logPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	_, err = data.InsertInto("search_logs", &logPayload)
	if err != nil {
		errorJSON(w, err, http.StatusBadGateway)
		return
	}

	resp := models.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	writeJSON(w, http.StatusAccepted, &resp)
}

// SearchOneEntry searches for one SearchEntry and writes a JsonResponse
// to the http.ResponseWriter with the retrieved data or the error that occured
func SearchOneEntry(w http.ResponseWriter, r *http.Request) {
	var searchPayload models.SearchQuery

	err := readJSON(w, r, &searchPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	entry, err := data.SearchEntriesByKeyword(&searchPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	resp := models.JsonResponse{
		Error:   false,
		Message: "Search was successful!",
		Data:    entry,
	}

	writeJSON(w, http.StatusAccepted, &resp)
}

// SearchPDF searches for the PDFEntry with the provided SearchQuery 
// and writes JsonResponse with the retrieved PDFEntry or the error that occured
func SearchPDF(w http.ResponseWriter, r *http.Request) {
	var searchPayload models.SearchQuery

	err := readJSON(w, r, &searchPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	pdfEntry, err := data.SearchForPDF(&searchPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	resp := models.JsonResponse{
		Error:   false,
		Message: "PDF retrieved successfully!",
		Data:    pdfEntry,
	}

	writeJSON(w, http.StatusAccepted, &resp)

}
