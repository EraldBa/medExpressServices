package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"search-service/models"
)

// searchRequest is the object sent to the microservices when requesting a search
// of any kind
type searchRequest struct {
	Keyword string `json:"keyword"`
}

// RequestSearchEntry requests a SearchEntry data from the given service url
// and returns a SearchEntry or potentially an error
func RequestSearchEntry(keyword, site, searchURL string) (*models.SearchEntry, error) {
	body := searchRequest{
		Keyword: keyword,
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
