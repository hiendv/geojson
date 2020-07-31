package util

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func NormalizeString(str string) string {
	input := []byte(str)
	b := make([]byte, len(input))
	n, _, _ := normalizer.Transform(b, input, true)

	return string(b[:n])
}
