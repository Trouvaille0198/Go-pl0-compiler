package lexer

import (
	"gopl0/fp"
	"gopl0/symbol"
	"gopl0/token"
	"gopl0/utils"
	"io/ioutil"
	"log"
	"strings"
)

const (
	MAX_TOKEN_SIZE = 10
	MAX_NUM_SIZE   = 25
)

// Lexer 词法分析器
type Lexer struct {
	file    *fp.File
	curLine int             // 当前所在行号
	Symbols []symbol.Symbol // 符号数组
}

func NewLexer(filepath string) *Lexer {
	file := fp.NewFile(filepath)
	return &Lexer{file: file, curLine: 1}
}

// getCh 读取一个字符
func (l *Lexer) getCh() (ch rune, isEnd bool) {
	ch, isEnd = l.file.Read()
	if ch == '\n' {
		l.curLine++
	}
	ch = utils.ToLower(ch)
	return
}

// GetSym DFA方式获取符号
func (l *Lexer) GetSym() {
	var num int                   // 当前识别的数字
	var numLen int                // 当前识别的数字长度
	var char [MAX_TOKEN_SIZE]rune // 当前识别的标识符或关键字
	var charIndex int             // 当前识别的标识符或关键字的索引

	curState := START // 当前状态
	ch, isEnd := l.getCh()
outerLoop:
	for !isEnd {
		switch curState {
		case START:
			if utils.IsSpace(ch) {
				// 啥都不做
			} else if ch == '{' {
				// 注释开头
				curState = COMMENT
			} else if utils.IsDigit(ch) {
				// 数字开头
				curState = INNUM
				num = num*10 + int(ch-'0')
				numLen++
			} else if utils.IsLetter(ch) {
				// 标识符开头
				if charIndex >= MAX_TOKEN_SIZE {
					panic("标识符或关键字过长")
				}
				curState = INID
				char[charIndex] = ch
				charIndex++
			} else if ch == '<' {
				curState = LES
			} else if ch == '>' {
				curState = GTR
			} else if ch == ':' {
				curState = INBECOMES
			} else {
				// 单独字符
				curState = START
				if optToken, ok := token.GetOptToken(string(ch)); ok {
					l.Symbols = append(l.Symbols, symbol.Symbol{Tok: optToken, Line: l.curLine})
				} else {
					panic("未知字符: " + string(ch))
				}
			}
		case INNUM:
			if utils.IsDigit(ch) {
				num = num*10 + int(ch-'0')
				numLen++
			} else {
				// 数字结束
				curState = START
				if numLen > MAX_NUM_SIZE {
					panic("数字过长")
				} else {
					l.Symbols = append(l.Symbols,
						symbol.Symbol{Tok: token.NUMBERSYM, Num: num, Line: l.curLine})
				}
				num, numLen = 0, 0
				continue outerLoop // 暂停对下一个字符的读取
			}
		case COMMENT:
			if ch == '}' {
				// 注释结束
				curState = START
				// 不记录注释
			}
		case INID:
			if utils.IsLetter(ch) || utils.IsDigit(ch) {
				if charIndex >= MAX_TOKEN_SIZE {
					panic("标识符或关键字过长")
				}
				char[charIndex] = ch
				charIndex++
			} else {
				// 标识符结束
				curState = START
				idToken := token.GetIdToken(string(char[:charIndex]))
				if idToken == token.IDENTSYM {
					newVal := make([]rune, charIndex)
					copy(newVal, char[:charIndex])
					l.Symbols = append(
						l.Symbols, symbol.Symbol{Tok: idToken, Value: newVal, Line: l.curLine})
				} else {
					l.Symbols = append(l.Symbols, symbol.Symbol{Tok: idToken, Line: l.curLine})
				}

				charIndex = 0
				continue outerLoop // 暂停对下一个字符的读取
			}
		case INBECOMES:
			if ch == '=' {
				curState = BECOMES
			} else {
				curState = START
				continue outerLoop
			}
		case GTR:
			if ch == '=' {
				curState = GEQ
			} else {
				curState = START
				l.Symbols = append(l.Symbols, symbol.Symbol{Tok: token.GEQSYM, Line: l.curLine})
				continue outerLoop
			}
		case LES:
			if ch == '=' {
				curState = LEQ
			} else {
				curState = START
				l.Symbols = append(l.Symbols, symbol.Symbol{Tok: token.LEQSYM, Line: l.curLine})
				continue outerLoop
			}
		case BECOMES:
			curState = START
			l.Symbols = append(l.Symbols, symbol.Symbol{Tok: token.BECOMESSYM, Line: l.curLine})
			continue outerLoop
		case GEQ:
			curState = START
			l.Symbols = append(l.Symbols, symbol.Symbol{Tok: token.GEQSYM, Line: l.curLine})
			continue outerLoop
		case LEQ:
			curState = START
			l.Symbols = append(l.Symbols, symbol.Symbol{Tok: token.LEQSYM, Line: l.curLine})
			continue outerLoop
		}
		ch, isEnd = l.getCh() // 读取下一个字符
	}
}

// Save 保存词法分析结果
func (l *Lexer) Save(path string) {
	var builder strings.Builder
	for _, sym := range l.Symbols {
		builder.WriteString(sym.String())
		builder.WriteString("\n")
	}
	err := ioutil.WriteFile(path, []byte(builder.String()), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
