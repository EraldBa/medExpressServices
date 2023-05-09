package scrapers

import "strings"

// sanitize sanitizes a given string by removing common clutter (spaces around string and newlines)
func Sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)

	return s
}

// sanitizeWhileRemoving sanitizes a string while removing the given unwanted word n times
func SanitizeWhileRemoving(s, unwanted string, n int) string {
	s = strings.Replace(s, unwanted, "", n)
	s = strings.ReplaceAll(s, "  ", "")

	return Sanitize(s)
}
