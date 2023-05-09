package models

import "time"

// RequestPayload is the only type of payload that the broker recieves from the frontend / POST request 
type RequestPayload struct {
	Action string      `json:"action"`
	Search SearchQuery `json:"search,omitempty"`
	Entry  SearchEntry `json:"log,omitempty"`
}

// SearchQuery is the type of payload that provides the search info when a search is requested
type SearchQuery struct {
	Keyword       string   `json:"keyword"`
	SitesToSearch []string `json:"sites_to_search,omitempty"`
}

// SearchEntry is the type of payload that is received from the search-service (when a search was previously requested)
// and gets returned to the requester
type SearchEntry struct {
	ID        string           `bson:"_id,omitempty" json:"id,omitempty"`
	Keyword   string           `bson:"name" json:"name"`
	Origin    string           `bson:"origin" json:"origin"`
	Data      []map[string]any `bson:"data" json:"data"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_at"`
}
