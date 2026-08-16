package aftercare

import "strings"

func normalizeDescription(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func limitDescription(value string, limit int) string {
	characters := []rune(value)
	if len(characters) < limit {
		return value
	}
	return string(characters[:limit-1])
}
