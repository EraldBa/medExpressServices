package pubmedscraper

import (
	"log"
	"med-scraper-service/scrapers"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// PubMedScraper holds the collectors and articles for PubMed
type PubMedScraper struct {
	keyword      string
	searchColly  *colly.Collector
	articleColly *colly.Collector
	articles     []*scrapers.PubMedArticle
}

// New creates and returns a *PubMedScraper ready for scraping
func New(keyword string) scrapers.Scraper[scrapers.PubMedArticle] {
	p := &PubMedScraper{
		keyword:  url.QueryEscape(keyword),
		articles: make([]*scrapers.PubMedArticle, 0, scrapers.ArticlesPerPage),
	}

	p.initScrapers()

	return p
}

// initScrapers initializes all scrapers of PubMedScraper
func (p *PubMedScraper) initScrapers() {
	searchColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent(scrapers.UserAgent),
	)

	articleColly := p.newPubArticleCollector()

	searchColly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting url for search:", r.URL.String())
	})

	searchColly.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s with status code: %d\n", err.Error(), r.StatusCode)
	})

	searchColly.OnResponse(func(r *colly.Response) {
		log.Printf("Connection to %s successful with status code: %d\n", r.Request.AbsoluteURL(r.Request.URL.Path), r.StatusCode)
	})

	searchColly.OnHTML(".docsum-title", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		link = h.Request.AbsoluteURL(link)

		err := articleColly.Visit(link)
		if err != nil {
			log.Printf("can't visit link: %s with error %s\n", link, err.Error())
		}
	})

	p.searchColly = searchColly
	p.articleColly = articleColly
}

// newPubArticleCollector creates and returns the articleColly member
// for collecting the pubmed articles
func (p *PubMedScraper) newPubArticleCollector() *colly.Collector {
	articleColly := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
		colly.UserAgent(scrapers.UserAgent),
	)

	articleColly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting url to find data:", r.URL.String())
	})

	articleColly.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s with status code: %d\n", err.Error(), r.StatusCode)
	})

	articleColly.OnResponse(func(r *colly.Response) {
		log.Printf("Connection to %s successful with status code: %d\n", r.Request.AbsoluteURL(r.Request.URL.Path), r.StatusCode)
	})

	articleColly.OnHTML(".article-details", func(h *colly.HTMLElement) {
		var (
			keywords string
			authors  []string
		)

		pmid := h.DOM.Find(".current-id").First().Text()

		pmcid := h.DOM.Find("[data-ga-action=PMCID]").First().Text()
		pmcid = scrapers.Sanitize(pmcid)

		title := h.DOM.Find(".heading-title").First().Text()
		title = scrapers.Sanitize(title)

		link := h.DOM.Find(".id-link").First().AttrOr("href", "")

		abstract := h.DOM.Find("[id=abstract]").Text()
		abstract = scrapers.SanitizeAndRemove(abstract, "Abstract", 1)

		if i := strings.LastIndex(abstract, "Keywords:"); i > -1 {
			keywords = scrapers.SanitizeAndRemove(abstract[i:], "Keywords:", 1)

			abstract = scrapers.Sanitize(abstract[:i])
		}

		summary := scrapers.Summarize(abstract, 1)

		h.ForEach(".expanded-authors a[class=full-name]", func(i int, h *colly.HTMLElement) {
			authors = append(authors, h.Text)
		})

		article := &scrapers.PubMedArticle{
			PMID:     pmid,
			PMCID:    pmcid,
			Title:    title,
			Link:     link,
			Summary:  summary,
			Abstract: abstract,
			Keywords: keywords,
			Authors:  authors,
		}

		p.articles = append(p.articles, article)
	})

	return articleColly
}

// GetData starts the PubMedScraper with the PubMedScraper.keyword and
// stores the data collected in PubMedScraper.articles
func (p *PubMedScraper) GetData() ([]*scrapers.PubMedArticle, error) {
	const pubURL = "https://pubmed.ncbi.nlm.nih.gov/?term="

	finalURL := pubURL + p.keyword

	err := p.searchColly.Visit(finalURL)
	if err != nil {
		return nil, err
	}

	p.searchColly.Wait()
	p.articleColly.Wait()

	return p.articles, nil
}
