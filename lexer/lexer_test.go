package lexer

import (
	"testing"
)

//func TestRead(t *testing.T) {
//	filepath := "../assets/a.txt"
//	f := fp.NewFile(filepath)
//	ch, end := f.Read()
//	for end
//
//}

func TestLexer(t *testing.T) {
	filepath := "../assets/a.txt"
	lexer := NewLexer(filepath)
	lexer.GetSym()
	for _, sym := range lexer.symbols {
		t.Log(sym)
	}
}
