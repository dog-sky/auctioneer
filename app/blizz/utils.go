package blizz

import (
	"unicode"
)

func isRussian(text string) bool {
	for _, r := range text {
		if !unicode.Is(unicode.Letter, r) {
			continue
		}
		if !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}
	return true
}
