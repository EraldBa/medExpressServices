package wikicollector

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// WikiData is the struct object that gets returned by GetWikiData
type WikiData struct {
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

// GetWikiData returns a WikiData struct if the response from the wiki api with the provided keyword
// was successful, otherwise returns an error
func GetWikiData(keyword string) (*WikiData, error) {
	const wikiSummaryExtractURL = "https://en.wikipedia.org/api/rest_v1/page/summary/"

	finalURL := wikiSummaryExtractURL + url.PathEscape(keyword)

	response, err := http.Get(finalURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("status not accepted")
	}

	data := new(WikiData)

	err = json.NewDecoder(response.Body).Decode(data)

	return data, err
}
