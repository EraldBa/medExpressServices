package scraper

// NHSArticle holds the nhs article data
type NHSArticle struct {
	StandardArticleInfo
	Text string `json:"text"`
}

// PubMedArticle holds the pubmed article data
type PubMedArticle struct {
	StandardArticleInfo
	PMID     string   `json:"pmid"`
	PMCID    string   `json:"pmcid,omitempty"`
	Link     string   `json:"link,omitempty"`
	Abstract string   `json:"abstract,omitempty"`
	Authors  []string `json:"authors,omitempty"`
}

// StandardArticleInfo holds all standard article data
type StandardArticleInfo struct {
	Title    string   `json:"title"`
	Summary  string   `json:"summary,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
}
