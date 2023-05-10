package nhsscraper

import (
	"log"
	"med-scraper-service/scrapers"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type NhsScraper struct {
	keyword      string
	searchColly  *colly.Collector
	articleColly *colly.Collector
	articles     []*Article
}

type Article struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Summary string `json:"summary,omitempty"`
}

const nhsURL = "https://www.nhs.uk/search/results?q="

func New(keyword string) *NhsScraper {
	n := newDefaultNhsScraper(keyword)

	n.initScrapers()

	return n
}

func newDefaultNhsScraper(keyword string) *NhsScraper {
	const articlesPerPage = 10

	return &NhsScraper{
		keyword:  url.QueryEscape(keyword),
		articles: make([]*Article, 0, articlesPerPage),
	}
}

func (n *NhsScraper) initScrapers() {
	searchColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent(scrapers.UserAgent),
	)

	searchColly.AllowURLRevisit = false

	articleColly := n.newNhsArticleScraper()

	searchColly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting url for search:", r.URL.String())
	})

	searchColly.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s with status code: %d\n", err.Error(), r.StatusCode)
	})

	searchColly.OnResponse(func(r *colly.Response) {
		log.Printf("Connection to %s successful with status code: %d\n", r.Request.AbsoluteURL(r.Request.URL.Path), r.StatusCode)
	})

	searchColly.OnHTML(".nhsuk-list", func(h *colly.HTMLElement) {
		h.ForEach("li", func(i int, h *colly.HTMLElement) {
			link := h.DOM.Find("h2").Find("a").AttrOr("href", "")
			link = h.Request.AbsoluteURL(link)

			// If it's trying to revisit the search page
			if strings.Contains(link, nhsURL) {
				return
			}

			err := articleColly.Visit(link)
			if err != nil {
				log.Printf("can't visit link: %s with error %s\n", link, err.Error())
			}
		})

	})

	n.searchColly = searchColly
	n.articleColly = articleColly

}

func (n *NhsScraper) newNhsArticleScraper() *colly.Collector {
	articleColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.UserAgent(scrapers.UserAgent),
	)

	articleColly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting url for article collection:", r.URL.String())
	})

	articleColly.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s with status code: %d\n", err.Error(), r.StatusCode)
	})

	articleColly.OnResponse(func(r *colly.Response) {
		log.Printf("Connection to %s successful with status code: %d\n", r.Request.AbsoluteURL(r.Request.URL.Path), r.StatusCode)
	})

	articleColly.OnHTML("body", func(h *colly.HTMLElement) {
		var text string

		title := h.DOM.Find("h1").First().Text()
		title = scrapers.Sanitize(title)

		h.ForEach("p", func(i int, h *colly.HTMLElement) {
			text += scrapers.Sanitize(h.Text)
		})

		summary := scrapers.Summarize(text, 1)

		article := &Article{
			Title:   title,
			Text:    text,
			Summary: summary,
		}

		n.articles = append(n.articles, article)
	})

	return articleColly
}

func (n *NhsScraper) GetData() ([]*Article, error) {
	finalURL := nhsURL + n.keyword

	err := n.searchColly.Visit(finalURL)
	if err != nil {
		return nil, err
	}

	n.searchColly.Wait()
	n.articleColly.Wait()

	return n.articles, nil
}
