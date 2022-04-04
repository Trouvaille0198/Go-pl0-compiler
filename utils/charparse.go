package utils

import (
	"unicode"
)

func IsSpace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func IsLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// ToLower 转换成小写
func ToLower(ch rune) rune {
	if 'A' <= ch && ch <= 'Z' {
		return ch + ('a' - 'A')
	}
	return ch
}
