package xar

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func Preprocess(input string) string {
	// Remove leading/trailing whitespace.
	input = strings.TrimSpace(input)
	// Normalize Unicode characters.
	return norm.NFC.String(input)
}
