package main

import (
	"med-api-service/collectors/pdfcollector"
	"med-api-service/collectors/wikicollector"
	"net/http"
)

// SearchRequestData is one of the types of payloads that the med-api-service receives.
// It containes the info for the search requested
type SearchRequestData struct {
	Keyword string `json:"keyword"`
}

// WikiSummary writes to the http.ResponseWriter the requested wiki summary
// received from the WikiPedia api or writes an errorJSON if a problem was encountered
func WikiSummary(w http.ResponseWriter, r *http.Request) {
	var dataSend [1]*wikicollector.WikiData

	search := new(SearchRequestData)

	err := readJSON(w, r, search)
	if err != nil {
		errorJSON(w, err)
		return
	}

	dataSend[0], err = wikicollector.GetWikiData(search.Keyword)
	if err != nil {
		errorJSON(w, err)
		return
	}

	// Sending an array because search-service expects an array of json objects
	writeJSON(w, http.StatusAccepted, dataSend)
}

// CollectPDF gets the pdf according to the provided PMCID from the SearchRequestData payload
// and writes the found pdf as a string to the htpp.ResponseWriter or an errorJSON if an error was encountered
func CollectPDF(w http.ResponseWriter, r *http.Request) {
	pmid := new(SearchRequestData)

	err := readJSON(w, r, pmid)
	if err != nil {
		errorJSON(w, err)
		return
	}

	pdf, err := pdfcollector.GetPDFByPMCID(pmid.Keyword)
	if err != nil {
		errorJSON(w, err)
		return
	}

	data := jsonResponse{
		Error:   false,
		Message: "PDF retrieval was successful",
		Data:    pdf,
	}

	writeJSON(w, http.StatusAccepted, data)
}
