package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	filepath := "../assets/b.txt"
	p := NewParser(filepath)
	p.Lex()
	p.Parse()
}
