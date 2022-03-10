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
	for _, sym := range lexer.symbols {
		// 记录各标识符出现的次数
		if sym.Id.IsIdent() {
			hash[string(sym.Value)]++
		}
	}
	for k, v := range hash {
		res := "(" + k + "," + strconv.Itoa(v) + ")"
		fmt.Println(res)
	}
}

func showLexResult(t *testing.T) {
	filepath := "../assets/a.txt"
	lexer := NewLexer(filepath)
	lexer.GetSym()
	for _, sym := range lexer.symbols {
		fmt.Println(sym.String())
	}
}

func TestExp2(t *testing.T) {
	showLexResult(t)
}
