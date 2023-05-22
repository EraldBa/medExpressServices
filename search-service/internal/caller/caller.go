package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"search-service/internal/models"
	"search-service/internal/sites"
)

// searchRequest is the object sent to the microservices when requesting a search
// of any kind
type searchRequest struct {
	Keyword string `json:"keyword"`
	Site    string `json:"site"`
}

// RequestSearchEntry requests a SearchEntry data from the given service url
// and returns a SearchEntry or potentially an error
func RequestSearchEntry(keyword, site string) (*models.SearchEntry, error) {
	searchURL, err := getUrlForSite(site)
	if err != nil {
		return nil, err
	}

	body := searchRequest{
		Keyword: keyword,
		Site:    site,
	}

	bodyBytes, _ := json.Marshal(body)

	response, err := http.Post(searchURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return nil, errors.New("status not accepted calling service")
	}

	result := &models.SearchEntry{
		Keyword: keyword,
		Origin:  site,
	}

	err = json.NewDecoder(response.Body).Decode(&result.Data)

	return result, err
}

// RequestPDFEntry requests a pdf from the pdf service and
// returns a PDFEntry or potentially an error
func RequestPDFEntry(pmid string) (*models.PDFEntry, error) {
	const pdfCollectURL = "http://med-api-service/collect-pdf"

	body := searchRequest{
		Keyword: pmid,
	}

	bodyBytes, _ := json.Marshal(body)

	response, err := http.Post(pdfCollectURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return nil, errors.New("status not accepted calling service")
	}

	data := new(models.JsonResponse)

	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	pdf, ok := data.Data.(string)
	if !ok {
		return nil, errors.New("can't convert pdf to text from json response")
	}

	result := &models.PDFEntry{
		PMID:    pmid,
		PDFText: pdf,
	}

	return result, nil
}

// getUrlForSite gets the appropriate microservice url for the provided site
func getUrlForSite(site string) (string, error) {
	const (
		wikiURL    = "http://med-api-service/wiki-summary"
		scraperURL = "http://med-scraper-service/scrape"
	)

	switch site {
	case sites.PubMed, sites.NHS:
		return scraperURL, nil
	case sites.Wikipedia:
		return wikiURL, nil
	}

	return "", errors.New("not a valid site entry")
}
