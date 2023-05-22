package nlp

import (
	textrank "github.com/DavidBelicza/TextRank/v2"
)

// Summarize summarizes given text to given number of sentences
func Summarize(text string, n int) string {
	var summary string

	tr := newDefaultTextRank(text)

	sentences := textrank.FindSentencesByRelationWeight(tr, n)

	for _, sentence := range sentences {
		summary += sentence.Value
	}

	return summary
}

// KeywordsFor gets n keywords for given text
func KeywordsFor(text string, n int) []string {
	tr := newDefaultTextRank(text)

	keywords := textrank.FindSingleWords(tr)

	if len(keywords) < n {
		n -= len(keywords) - n
	}

	sanitizedKeywords := make([]string, n)

	for i := range keywords[:n] {
		sanitizedKeywords[i] = keywords[i].Word
	}

	return sanitizedKeywords
}

// newDefaultTextRank creates a new TextRank object with default settings
func newDefaultTextRank(text string) *textrank.TextRank {
	tr := textrank.NewTextRank()

	tr.Populate(text, textrank.NewDefaultLanguage(), textrank.NewDefaultRule())

	tr.Ranking(textrank.NewDefaultAlgorithm())

	return tr
}

