package models

import "time"

// DataEntry is the interface that has to be implemented my all
// structs that are ment to be inserted in mongodb
type DataEntry interface {
	AddDefaultData()
}

// SearchEntry holds the data to insert or pull out of a 'search_logs'
type SearchEntry struct {
	ID      string           `bson:"_id,omitempty" json:"id,omitempty"`
	Keyword string           `bson:"keyword" json:"keyword"`
	Origin  string           `bson:"origin" json:"origin"`
	Data    []map[string]any `bson:"data" json:"data"`
	Times
}

// PDFEntry holds the data to insert or pull out of a 'pdf_logs' collection
type PDFEntry struct {
	ID      string `bson:"_id,omitempty" json:"id,omitempty"`
	PMID    string `bson:"pmid" json:"pmid"`
	PDFText string `bson:"pdf_text" json:"pdf_text"`
	Times
}

// SeachQuery holds the keyword to be searched as well as the site preferences
type SearchQuery struct {
	Keyword       string   `json:"keyword"`
	SitesToSearch []string `json:"sites_to_search,omitempty"`
}

// JsonResponse is the standard response object that the service writes to the http.ResponseWriter
type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Times holds the standard time data for a mongo entry
type Times struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
