package expparser

import (
	"fmt"
	"gopl0/lexer"
	pa "gopl0/parser"
	"gopl0/parser/asm"
	"gopl0/parser/fct"
	"gopl0/symbol"
	"gopl0/token"
	"log"
	"strconv"
)

const (
	MAX_LEVEL = 3 // 嵌套的最大层数
)

// 符号表中的符号
type identItem struct {
	lit       string       // 字面量
	identType pa.IdentType // 标识符类型
	level     int          // 所在层级
	value     int          // 无符号整数的值
	addr      int          // 每层局部量的相对地址 dx
}

// fourCode 四元组
type fourCode struct {
	op     string // 操作符
	first  string // 第一个操作数
	second string // 第二个操作数
	result string // 结果
}

func (f *fourCode) String() string {
	return "( " + f.op + " " + f.first + " " + f.second + " " + f.result + " )"
}

// Parser 语法分析器
type ExpParser struct {
	lexer *lexer.Lexer // 词法分析器
	dx    int          // 当前层局部量的相对地址
	level int          // 嵌套层数

	table []identItem // 符号表
	tx    int         // 指针

	//中间代码
	codes []asm.Asm
	cidx  int

	// 当前扫描到的symbol下标
	curSymbolIndex int
	// 四元组
	FourCodes []fourCode
	fctx      int //指针
	tno       int
}

func (p *ExpParser) enterFourCode(op, first, second, result string) string {
	p.fctx++
	if result == "#" {
		result = result + strconv.Itoa(p.tno)
		p.tno += 1
	}
	fc := fourCode{op: op, first: first, second: second, result: result}
	p.FourCodes[p.fctx] = fc
	// fmt.Printf("%s\n", fc)
	return result
}

func NewExpParser(filepath string) *ExpParser {
	p := ExpParser{lexer: lexer.NewLexer(filepath),
		table: make([]identItem, 100), codes: make([]asm.Asm, 100),
		FourCodes: make([]fourCode, 100)}
	return &p
}

// getCurSymbol 获取当前扫描到的符号
func (p *ExpParser) getCurSymbol() *symbol.Symbol {
	if p.curSymbolIndex >= len(p.lexer.Symbols) {
		return &symbol.Symbol{}
	}
	return &p.lexer.Symbols[p.curSymbolIndex]
}

// goNextSymbol 扫描下一个符号
func (p *ExpParser) goNextSymbol() {
	p.curSymbolIndex++
}

// enter 添加符号至符号表
func (p *ExpParser) enter(identType pa.IdentType) {
	p.tx++
	tmp := identItem{
		lit:       p.getCurSymbol().GetLit(),
		level:     p.level,
		identType: identType,
		value:     p.getCurSymbol().Num,
		addr:      p.dx,
	}
	if identType.IsVariable() {
		// 变量 需要在数据栈里留空间
		p.dx++
	}
	p.table[p.tx] = tmp
}

// checkInTable 获取标识符在符号表的位置(若有)
func (p *ExpParser) checkInTable(lit string) (id identItem, in bool) {
	for _, v := range p.table {
		if lit == v.lit {
			return v, true
		}
	}
	return identItem{}, false
}

// gen 生成中间代码
func (p *ExpParser) gen(fct fct.Fct, y, z int) {
	p.codes = append(p.codes, asm.Asm{
		Fct: fct,
		L:   y,
		A:   z,
	})
	p.cidx++
}

// tokenIsInTokens 检查token是否在tokens中
func tokenIsInTokens(toks []token.Token, tok token.Token) bool {
	for _, v := range toks {
		if v == tok {
			return true
		}
	}
	return false
}

// test 检查当前扫描到的符号是不是在select集中
func (p *ExpParser) test(toks []token.Token, t int) {
	if !tokenIsInTokens(toks, p.getCurSymbol().Tok) {
		log.Printf("出错的当前单词: %s, %s\n", p.getCurSymbol().GetDesc(), p.getCurSymbol().GetLit())
		for !tokenIsInTokens(toks, p.getCurSymbol().Tok) && p.getCurSymbol().Tok != token.PERIODSYM {
			// 向后扫描
			p.goNextSymbol()
		}
	}
}

// factor 因子的产生式
// <因子> → <标识符>|<无符号整数>|(<表达式>)
func (p *ExpParser) factor() string {
	first := string(p.getCurSymbol().Value)
	switch p.getCurSymbol().Tok {
	case token.IDENTSYM: // 标识符
		// 判断标识符是否已经定义 并且生成中间代码
		id, ok := p.checkInTable(p.getCurSymbol().GetLit())
		if ok {
			// 已经定义
			switch id.identType {
			// 可能是常量或者变量
			case pa.Constant:
				p.gen(fct.Lit, 0, id.value)
			case pa.Variable:
				p.gen(fct.Lod, p.level-id.level, id.addr)
			case pa.Proc:
				// 不允许接受过程
				p.Error(21)
			}
		} else {
			// 标识符未定义
			p.Error(11)
		}
		p.goNextSymbol()
	case token.NUMBERSYM: // 无符号整数
		p.gen(fct.Lit, 0, p.getCurSymbol().Num)
		p.goNextSymbol()
	case token.LPARENTSYM: // 左括号
		p.goNextSymbol()
		first = p.expression()
		// 判断表达式结束后是不是右括号
		if p.getCurSymbol().Tok == token.RPARENTSYM {
			p.goNextSymbol()
		} else {
			// 缺少右括号
			p.Error(22)
		}
	}
	return first
}

// term 项的产生式
// <项> → <因子>{<乘除运算符><因子>}
func (p *ExpParser) term() string {
	// <因子>
	first := p.factor()
	// {<乘除运算符><因子>}
	for p.getCurSymbol().Tok == token.MULSYM || p.getCurSymbol().Tok == token.SLASHSYM || p.getCurSymbol().Tok == token.MODSYM {
		opt := p.getCurSymbol().Tok
		p.goNextSymbol()
		second := p.factor()
		if opt == token.MULSYM {
			p.gen(fct.Opr, 0, 4) // 乘法
			first = p.enterFourCode("*", first, second, "#")
		} else if opt == token.SLASHSYM {
			p.gen(fct.Opr, 0, 5) // 除法
			first = p.enterFourCode("/", first, second, "#")
		} else {
			p.gen(fct.Opr, 0, 6) // 取余 肯定不是6
		}
	}
	return first
}

// expression 表达式的产生式
// <表达式> → [+|-]<项>{<加减运算符><项>}
func (p *ExpParser) expression() string {
	var opt token.Token
	var first string
	// [+|-]<项>
	if p.getCurSymbol().Tok == token.PLUSSYM || p.getCurSymbol().Tok == token.MINUSYM {
		// 可能出现正负号
		opt = p.getCurSymbol().Tok
		p.goNextSymbol()
		first = p.term()
		if opt == token.MINUSYM {
			p.gen(fct.Opr, 0, 1)
			first = p.enterFourCode("-", "0", first, "#")
		}
	} else {
		first = p.term()
	}
	// {<加减运算符><项>}
	for p.getCurSymbol().Tok == token.PLUSSYM || p.getCurSymbol().Tok == token.MINUSYM {
		opt = p.getCurSymbol().Tok
		p.goNextSymbol()
		second := p.term()
		if opt == token.PLUSSYM {
			p.gen(fct.Opr, 0, 2) // 加
			first = p.enterFourCode("+", first, second, "#")
		} else {
			p.gen(fct.Opr, 0, 3) // 减
			first = p.enterFourCode("-", first, second, "#")
		}
	}
	return first
}

// varDeclaration 变量声明 识别var后进入此过程 其实就是识别标识符
func (p *ExpParser) varDeclaration() {
	if p.getCurSymbol().Tok.IsIdent() {
		p.enter(pa.Variable)
		p.goNextSymbol()
	} else {
		p.Error(38)
	}
	for p.getCurSymbol().Tok.IsComma() {
		p.goNextSymbol()
		if p.getCurSymbol().Tok.IsIdent() {
			p.enter(pa.Variable)
			p.goNextSymbol()
		} else {
			p.Error(38)
		}
	}
}

func (p *ExpParser) Lex() {
	p.lexer.GetSym()
}

// Parse 语法分析入口
func (p *ExpParser) Parse() {
	// p.lexer.GetSym()
	p.goNextSymbol()
	p.varDeclaration()
	p.expression()
	// log.Printf("结束\n")
}

// SaveLexResult 保存词法分析结果
func (p *ExpParser) SaveLexResult(path string) {
	p.lexer.Save(path)
}

func (p *ExpParser) ShowLexResult() {
	for _, sym := range p.lexer.Symbols {
		fmt.Println(sym.String())
	}
}

func (p *ExpParser) ShowFourCodes() {
	for _, fc := range p.FourCodes {
		if fc.op != "" {
			fmt.Println(fc.String())
		}
	}
}
