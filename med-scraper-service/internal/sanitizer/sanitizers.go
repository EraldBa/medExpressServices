package sanitizer

import (
	"strings"
	"unicode"
)

const ParagraphTerminator = "|"

var (
	commonClutter = []string{
		"\n",
		"  ",
	}

	keywordClutter = []string{
		":",
		".",
		ParagraphTerminator,
	}
)

// Sanitize sanitizes a given string by removing common clutter (spaces around strings, newlines and invisible characters)
func Sanitize(s string) string {
	s = removeClutter(s, commonClutter)
	s = strings.TrimSpace(s)

	removeInvisibleRune := func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}

	// Removing invisible characters
	s = strings.Map(removeInvisibleRune, s)

	return s
}

// sanitizeAndRemove sanitizes a string while removing the given unwanted word n times
func SanitizeAndRemove(s, unwanted string, n int) string {
	s = strings.Replace(s, unwanted, "", n)

	return Sanitize(s)
}

// SanitizeKeywords sanitizes given slice of keywords from common keyword clutter
func SanitizeKeywords(keywords []string) {
	for i, keyword := range keywords {
		keyword = removeClutter(keyword, keywordClutter)
		keyword = Sanitize(keyword)

		keywords[i] = keyword
	}
}

// removeClutter removes the clutter for the provided string according to
// the provided clutterCollection
func removeClutter(s string, clutterCollection []string) string {
	for _, clutter := range clutterCollection {
		s = strings.ReplaceAll(s, clutter, "")
	}

	return s
}
