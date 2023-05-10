package main

import (
	"med-scraper-service/scrapers"
	"med-scraper-service/scrapers/nhsscraper"
	"med-scraper-service/scrapers/pubmedscraper"
	"net/http"
)

// SearchRequest specifies the search that gets requested
type SearchRequest struct {
	Keyword string `json:"keyword"`
}

// CollectPubMed collects data from the pubmed website with the provided
// keyword using the PubMedScraper and writes the data or an errorJSON to
// the http.ResponseWriter
func CollectPubMed(w http.ResponseWriter, r *http.Request) {
	collectorFor(w, r, pubmedscraper.New)
}

// CollectNHS collects data from the nhs website with the provided
// keyword using the NhsScraper and writes the data or an errorJSON to
// the http.ResponseWriter
func CollectNHS(w http.ResponseWriter, r *http.Request) {
	collectorFor(w, r, nhsscraper.New)
}

// collectFor collects data according to the generic newScraper function provided, using scrapers.Scraper.GetData method
func collectorFor[A scrapers.Article](w http.ResponseWriter, r *http.Request, newScraper func(string) scrapers.Scraper[A]) {
	searchRequest := new(SearchRequest)

	err := readJSON(w, r, searchRequest)
	if err != nil {
		errorJSON(w, err)
		return
	}

	collector := newScraper(searchRequest.Keyword)

	data, err := collector.GetData()
	if err != nil {
		errorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, data)
}
