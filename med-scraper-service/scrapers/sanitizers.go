package scrapers

import (
	"strings"
	"unicode"
)

// Sanitize sanitizes a given string by removing common clutter (spaces around strings, newlines and invisible characters)
func Sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "  ", "")
	s = strings.TrimSpace(s)

	// Removing invisible characters
	s = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)

	return s
}

// sanitizeAndRemove sanitizes a string while removing the given unwanted word n times
func SanitizeAndRemove(s, unwanted string, n int) string {
	s = strings.Replace(s, unwanted, "", n)

	return Sanitize(s)
}
