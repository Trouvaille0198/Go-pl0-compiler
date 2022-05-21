package lexer

import (
	"fmt"
	"strconv"
	"testing"
)

func TestExp1(t *testing.T) {
	filepath := "../assets/a.txt"
	lexer := NewLexer(filepath)
	lexer.GetSym()

	hash := make(map[string]int)
	for _, sym := range lexer.Symbols {
		// 记录各标识符出现的次数
		if sym.Tok.IsIdent() {
			hash[string(sym.Value)]++
		}
	}
	for k, v := range hash {
		res := "(" + k + "," + strconv.Itoa(v) + ")"
		fmt.Println(res)
	}
}

func showLexResult(t *testing.T) {
	filepath := "../assets/c.txt"
	lexer := NewLexer(filepath)
	lexer.GetSym()
	for _, sym := range lexer.Symbols {
		fmt.Println(sym.String())
	}
}

func TestExp2(t *testing.T) {
	showLexResult(t)
}

func TestSave(t *testing.T) {
	filepath := "../assets/b.txt"
	lexer := NewLexer(filepath)
	lexer.GetSym()

	savePath := "../assets/b-result.txt"
	lexer.Save(savePath)
}
