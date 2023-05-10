package scrapers

// Article generic interface for all structs that hold article data
type Article interface {
	NHSArticle | PubMedArticle
}

// Scraper interface with generic parameter Article for all Scrapers
type Scraper[A Article] interface {
	GetData() ([]*A, error)
}
