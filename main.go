package main

import (
	"fmt"
	"gopl0/expparser"
)

func main() {
	p := expparser.NewExpParser("./assets/c.txt")
	p.Lex()
	// p.ShowLexResult()
	p.Parse()
	// savePath := "../assets/a-result.txt"
	// p.SaveLexResult(savePath)
	// fmt.Printf("%+v\n", p.FourCodes)
	fmt.Println("12321")
	p.ShowFourCodes()
}
