package api

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NormalizeString(value string) string {
	value = strings.ToLower(value)
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, err := transform.String(t, value)
	if err != nil {
		return strings.TrimSpace(value)
	}
	return strings.TrimSpace(result)
}
