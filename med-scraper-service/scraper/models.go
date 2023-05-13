package scraper

// NHSArticle holds the nhs article data
type NHSArticle struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Summary string `json:"summary,omitempty"`
}

// PubMedArticle holds the pubmed article data
type PubMedArticle struct {
	PMID     string   `json:"pmid"`
	Title    string   `json:"title"`
	PMCID    string   `json:"pmcid,omitempty"`
	Link     string   `json:"link,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	Abstract string   `json:"abstract,omitempty"`
	Keywords string   `json:"keywords,omitempty"`
	Authors  []string `json:"authors,omitempty"`
}
