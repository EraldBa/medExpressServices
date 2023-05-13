package scrapers

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

const (
	PubURL          = "https://pubmed.ncbi.nlm.nih.gov/?term="
	UserAgent       = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"
	ArticlesPerPage = 10
	NhsURL          = "https://www.nhs.uk/search/results?q="
)

type scraper struct {
	url          string
	searchColly  *colly.Collector
	articleColly *colly.Collector
	articles     []any
}

// New returns a scraper and initializes it according to the site provided
func New(keyword, site string) (*scraper, error) {
	s := new(scraper)

	keyword = url.QueryEscape(keyword)

	switch site {
	case "pubmed":
		s.url = PubURL + keyword
		s.initPubMedScrapers()
	case "nhs":
		s.url = NhsURL + keyword
		s.initNhsScrapers()
	default:
		return nil, errors.New("site provided not valid")
	}

	s.articles = make([]any, 0, ArticlesPerPage)

	return s, nil
}

// GetData starts the scraper with the keyword and url
// stores the data collected in articles
func (s *scraper) GetData() ([]any, error) {
	err := s.searchColly.Visit(s.url)
	if err != nil {
		return nil, err
	}

	s.searchColly.Wait()
	s.articleColly.Wait()

	return s.articles, nil
}

// initNhsScrapers initializes scrapers for nhs
func (s *scraper) initNhsScrapers() {
	searchColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent(UserAgent),
	)

	articleColly := s.newNhsArticleScraper()

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
			if strings.Contains(link, NhsURL) {
				return
			}

			err := articleColly.Visit(link)
			if err != nil {
				log.Printf("can't visit link: %s with error %s\n", link, err.Error())
			}
		})

	})

	s.searchColly = searchColly
	s.articleColly = articleColly

}

func (s *scraper) newNhsArticleScraper() *colly.Collector {
	articleColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.UserAgent(UserAgent),
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
		title = Sanitize(title)

		h.ForEach("p", func(i int, h *colly.HTMLElement) {
			text += Sanitize(h.Text)
		})

		summary := Summarize(text, 1)

		article := &NHSArticle{
			Title:   title,
			Text:    text,
			Summary: summary,
		}

		s.articles = append(s.articles, article)
	})

	return articleColly
}

// initPubMedScrapers initializes all for pubmed
func (s *scraper) initPubMedScrapers() {
	searchColly := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent(UserAgent),
	)

	articleColly := s.newPubArticleCollector()

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

	s.searchColly = searchColly
	s.articleColly = articleColly
}

// newPubArticleCollector creates and returns the articleColly member
// for collecting the pubmed articles
func (s *scraper) newPubArticleCollector() *colly.Collector {
	articleColly := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
		colly.UserAgent(UserAgent),
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
		pmcid = Sanitize(pmcid)

		title := h.DOM.Find(".heading-title").First().Text()
		title = Sanitize(title)

		link := h.DOM.Find(".id-link").First().AttrOr("href", "")

		abstract := h.DOM.Find("[id=abstract]").Text()
		abstract = SanitizeAndRemove(abstract, "Abstract", 1)

		if i := strings.LastIndex(abstract, "Keywords:"); i > -1 {
			keywords = SanitizeAndRemove(abstract[i:], "Keywords:", 1)

			abstract = Sanitize(abstract[:i])
		}

		summary := Summarize(abstract, 1)

		h.ForEach(".expanded-authors a[class=full-name]", func(i int, h *colly.HTMLElement) {
			authors = append(authors, h.Text)
		})

		article := &PubMedArticle{
			PMID:     pmid,
			PMCID:    pmcid,
			Title:    title,
			Link:     link,
			Summary:  summary,
			Abstract: abstract,
			Keywords: keywords,
			Authors:  authors,
		}

		s.articles = append(s.articles, article)
	})

	return articleColly
}
