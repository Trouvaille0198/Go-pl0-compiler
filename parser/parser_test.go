package parser

import (
	"fmt"
	"testing"
)

//func TestParser(t *testing.T) {
//	filepath := []string{"../assets/a.txt", "../assets/b.txt", "../assets/c.txt", "../assets/d.txt"}
//	for i := 0; i < len(filepath); i++ {
//		var wg sync.WaitGroup
//		go func() {
//			wg.Add(1)
//			p := NewParser(filepath[i])
//			p.Lex()
//			p.Parse()
//			wg.Done()
//		}()
//		wg.Wait()
//	}
//}

func TestParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewParser("../assets/c.txt")
		p.Lex()
		p.Parse()
		// savePath := "../assets/a-result.txt"
		// p.SaveLexResult(savePath)
		for _, sym := range p.codes {
			fmt.Println(sym.String())
		}
	}
}
