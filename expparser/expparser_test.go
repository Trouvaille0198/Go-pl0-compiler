package expparser

import (
	"fmt"
	"testing"
)

func TestExpParser(t *testing.T) {
	p := NewExpParser("../assets/c.txt")
	p.Lex()
	p.Parse()
	// savePath := "../assets/a-result.txt"
	// p.SaveLexResult(savePath)
	fmt.Println("12321")
	for _, fc := range p.FourCodes {
		fmt.Println(fc.String())
	}
}
