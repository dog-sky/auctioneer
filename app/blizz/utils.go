package blizz

import (
	"strings"
	"unicode"
)

func isRussian(text string) bool {
	text = strings.ReplaceAll(text, " ", "")
	for _, r := range text {
		if !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}
	return true
}
