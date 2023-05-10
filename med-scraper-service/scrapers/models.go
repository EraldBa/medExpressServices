package scrapers

// NHSArticle holds the nhs article data
type NHSArticle struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Summary string `json:"summary,omitempty"`
}

// PubMedArticle holds the pubmed article data
type PubMedArticle struct {
	PMID     string   `json:"pmid"`
	PMCID    string   `json:"pmcid"`
	Title    string   `json:"title"`
	Link     string   `json:"link,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	Abstract string   `json:"abstract,omitempty"`
	Keywords string   `json:"keywords,omitempty"`
	Authors  []string `json:"authors,omitempty"`
}
