package util

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeString transforms an unicode string str to ASCII form
func NormalizeString(str string) string {
	normalizer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	input := []byte(str)
	b := make([]byte, len(input))
	n, _, _ := normalizer.Transform(b, input, true)

	return string(b[:n])
}
