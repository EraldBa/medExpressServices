package main

import (
	"med-scraper-service/scrapers"
	"net/http"
)

// SearchRequest specifies the search that gets requested
type SearchRequest struct {
	Keyword string `json:"keyword"`
	Site    string `json:"site"`
}

// Scrape scrapes provided site for provided keyword and responds with the collected data or an error
func Scrape(w http.ResponseWriter, r *http.Request) {
	searchRequest := new(SearchRequest)

	err := readJSON(w, r, searchRequest)
	if err != nil {
		errorJSON(w, err)
		return
	}

	scraper, err := scrapers.New(searchRequest.Keyword, searchRequest.Site)
	if err != nil {
		errorJSON(w, err)
		return
	}

	data, err := scraper.GetData()
	if err != nil {
		errorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, data)
}
