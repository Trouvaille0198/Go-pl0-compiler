package utils

import "unicode"

func IsSpace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func IsLetter(ch rune) bool {
	var lower func(ch rune) rune
	lower = func(ch rune) rune { return ('a' - 'A') | ch }

	return 'a' <= lower(ch) && lower(ch) <= 'z' && unicode.IsLetter(ch)
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
