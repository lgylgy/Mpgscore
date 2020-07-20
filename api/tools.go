package api

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"strings"
	"unicode"
)

func NormalizeString(value string) string {
	value = strings.ToLower(value)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, value)
	if err != nil {
		return strings.TrimSpace(value)
	}
	return strings.TrimSpace(result)
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("%v variable is not present", key)
	}
	log.Printf("%v: %v\n", key, value)
	return value
}
