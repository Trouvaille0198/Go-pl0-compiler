package main

import (
	"fmt"
	"gopl0/expparser"
)

func main() {
	p := expparser.NewExpParser("./assets/c.txt")
	p.Lex()
	p.ShowLexResult()
	p.Parse()
	fmt.Print("\n")
	p.ShowFourCodes()
}
