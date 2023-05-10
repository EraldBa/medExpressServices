package main

import (
	"med-scraper-service/scrapers/pubmedscraper"
	"net/http"
)

// SearchRequest specifies the search that gets requested
type SearchRequest struct {
	Keyword string `json:"keyword"`
}

// CollectPubMed collects data from the pubmed website with the provided
// keyword using the PubMedCollector and writes the data or an errorJSON to
// the http.ResponseWriter
func CollectPubMed(w http.ResponseWriter, r *http.Request) {
	searchRequest := new(SearchRequest)

	err := readJSON(w, r, searchRequest)
	if err != nil {
		errorJSON(w, err)
		return
	}

	collector := pubmedscraper.New(searchRequest.Keyword)

	data, err := collector.GetData()
	if err != nil {
		errorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, data)
}
