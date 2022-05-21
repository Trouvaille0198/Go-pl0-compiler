package parser

import (
	"gopl0/lexer"
	"gopl0/parser/asm"
	"gopl0/parser/fct"
	"gopl0/symbol"
	"gopl0/token"
	"log"
)

const (
	MAX_LEVEL = 3 // 嵌套的最大层数
)

// 符号表中的符号
type identItem struct {
	lit       string    // 字面量
	identType IdentType // 标识符类型
	level     int       // 所在层级
	value     int       // 无符号整数的值
	addr      int       // 每层局部量的相对地址 dx
}

// Parser 语法分析器
type Parser struct {
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
}

func NewParser(filepath string) *Parser {
	p := Parser{lexer: lexer.NewLexer(filepath), table: make([]identItem, 100), codes: make([]asm.Asm, 100)}
	return &p
}

// getCurSymbol 获取当前扫描到的符号
func (p *Parser) getCurSymbol() *symbol.Symbol {
	if p.curSymbolIndex > len(p.lexer.Symbols) {
		return &symbol.Symbol{}
	}
	return &p.lexer.Symbols[p.curSymbolIndex]
}

// goNextSymbol 扫描下一个符号
func (p *Parser) goNextSymbol() {
	p.curSymbolIndex++
}

// enter 添加符号至符号表
func (p *Parser) enter(identType IdentType) {
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
func (p *Parser) checkInTable(lit string) (id identItem, in bool) {
	for _, v := range p.table {
		if lit == v.lit {
			return v, true
		}
	}
	return identItem{}, false
}

// gen 生成中间代码
func (p *Parser) gen(fct fct.Fct, y, z int) {
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
func (p *Parser) test(toks []token.Token, t int) {
	if !tokenIsInTokens(toks, p.getCurSymbol().Tok) {
		p.Error(t)
		log.Printf("出错的当前单词: %s, %s\n", p.getCurSymbol().GetDesc(), p.getCurSymbol().GetLit())
		for !tokenIsInTokens(toks, p.getCurSymbol().Tok) && p.getCurSymbol().Tok != token.PERIODSYM {
			// 向后扫描
			p.goNextSymbol()
		}
	}
}

// constDeclaration 常量定义 识别const后进入此过程
// <常量定义> → <标识符>=<无符号整数>
func (p *Parser) constDeclaration() {
	if p.getCurSymbol().Tok.IsIdent() {
		p.enter(Constant)
		p.goNextSymbol() //看下一个
		if p.getCurSymbol().Tok == token.EQLSYM || p.getCurSymbol().Tok == token.BECOMESSYM {
			// 等号或赋值号
			if p.getCurSymbol().Tok == token.BECOMESSYM {
				// 容错处理，报错，但是当做等号使用
				p.Error(0)
			}

			p.goNextSymbol()
			// 数字就加入符号表
			if p.getCurSymbol().Tok.IsNumber() {
				p.goNextSymbol()
			} else {
				// 等号后不是数字
				p.Error(40)
			}
		} else {
			// 不是等号或赋值号，非常量声明部分
			p.Error(0)
		}
	} else {
		p.Error(0)
	}
}

// varDeclaration 变量声明 识别var后进入此过程 其实就是识别标识符
func (p *Parser) varDeclaration() {
	if p.getCurSymbol().Tok.IsIdent() {
		p.enter(Variable)
		p.goNextSymbol()
	} else {
		p.Error(38)
	}
	for p.getCurSymbol().Tok.IsComma() {
		p.goNextSymbol()
		if p.getCurSymbol().Tok.IsIdent() {
			p.enter(Variable)
			p.goNextSymbol()
		} else {
			p.Error(38)
		}
	}
}

// factor 因子的产生式
// <因子> → <标识符>|<无符号整数>|(<表达式>)
func (p *Parser) factor(toks []token.Token) {
	p.test(append(FactorSelect, toks...), 24)
	for tokenIsInTokens(FactorSelect, p.getCurSymbol().Tok) {
		switch p.getCurSymbol().Tok {
		case token.IDENTSYM: // 标识符
			// 判断标识符是否已经定义 并且生成中间代码
			id, ok := p.checkInTable(p.getCurSymbol().GetLit())
			if ok {
				// 已经定义
				switch id.identType {
				// 可能是常量或者变量
				case Constant:
					p.gen(fct.Lit, 0, id.value)
				case Variable:
					p.gen(fct.Lod, p.level-id.level, id.addr)
				case Proc:
					// 不允许接受过程
					p.Error(21)
				}
			} else {
				// 标识符未定义
				p.Error(11)
			}
			p.goNextSymbol()
		case token.NUMBERSYM: // 无符号整数
			// TODO 判断数字有没有溢出
			p.gen(fct.Lit, 0, p.getCurSymbol().Num)
			p.goNextSymbol()
		case token.LPARENTSYM: // 左括号
			p.goNextSymbol()
			p.expression(append(toks, token.RPARENTSYM))
			// 判断表达式结束后是不是右括号
			if p.getCurSymbol().Tok == token.RPARENTSYM {
				p.goNextSymbol()
			} else {
				// 缺少右括号
				p.Error(22)
			}
		}
		p.test(append(toks, []token.Token{token.LPARENTSYM}...), 23)
	}
}

// term 项的产生式
// <项> → <因子>{<乘除运算符><因子>}
func (p *Parser) term(toks []token.Token) {
	// <因子>
	newToks := append(toks, []token.Token{token.MULSYM, token.SLASHSYM, token.MODSYM}...)
	p.factor(newToks)
	// {<乘除运算符><因子>}
	for p.getCurSymbol().Tok == token.MULSYM || p.getCurSymbol().Tok == token.SLASHSYM || p.getCurSymbol().Tok == token.MODSYM {
		opt := p.getCurSymbol().Tok
		p.goNextSymbol()
		p.factor(newToks)
		// TODO 对取余的中间代码处理
		if opt == token.MULSYM {
			p.gen(fct.Opr, 0, 4) // 乘法
		} else if opt == token.MULSYM {
			p.gen(fct.Opr, 0, 5) // 除法
		} else {
			p.gen(fct.Opr, 0, 6) // 取余 肯定不是6
		}
	}
}

// expression 表达式的产生式
// <表达式> → [+|-]<项>{<加减运算符><项>}
func (p *Parser) expression(toks []token.Token) {
	var opt token.Token
	newToks := append(toks, []token.Token{token.PLUSSYM, token.MINUSYM}...)
	// [+|-]<项>
	if p.getCurSymbol().Tok == token.PLUSSYM || p.getCurSymbol().Tok == token.MINUSYM {
		// 可能出现正负号
		opt = p.getCurSymbol().Tok
		p.goNextSymbol()
		p.term(newToks)
		if opt == token.MINUSYM {
			p.gen(fct.Opr, 0, 1)
		}
	} else {
		p.term(newToks)
	}
	// {<加减运算符><项>}
	for p.getCurSymbol().Tok == token.PLUSSYM || p.getCurSymbol().Tok == token.MINUSYM {
		opt = p.getCurSymbol().Tok
		p.goNextSymbol()

		p.term(newToks)
		if opt == token.PLUSSYM {
			p.gen(fct.Opr, 0, 2) // 加
		} else {
			p.gen(fct.Opr, 0, 3) // 减
		}
	}
}

// condition 条件的产生式
// <条件> → <表达式><关系运算符><表达式>|odd<表达式>
func (p *Parser) condition(toks []token.Token) {
	if p.getCurSymbol().Tok == token.ODDSYM {
		// odd<表达式>
		p.goNextSymbol()
		p.expression(toks)
		p.gen(fct.Opr, 0, 6)
	} else {
		// <表达式><关系运算符><表达式>
		p.expression(append(toks, ExpressionSelect...))
		if p.getCurSymbol().Tok.IsRelationOpr() { //是关系运算符
			relopt := p.getCurSymbol().Tok
			p.goNextSymbol()

			p.expression(toks)
			switch relopt {
			case token.EQLSYM:
				p.gen(fct.Opr, 0, 8) // =
			case token.NEQSYM:
				p.gen(fct.Opr, 0, 9) // #
			case token.LESSYM:
				p.gen(fct.Opr, 0, 10) // <
			case token.GTRSYM:
				p.gen(fct.Opr, 0, 11) // >
			case token.LEQSYM:
				p.gen(fct.Opr, 0, 12) // <=
			case token.GEQSYM:
				p.gen(fct.Opr, 0, 13) // >=
			}
		}
	}
}

// statement 语句的产生式
// <语句> → <赋值语句>|<条件语句>|<当型循环语句>|<过程调用语句>|<读语句>|<写语句>|<复合语句>|<空>
// <赋值语句> → <标识符>:=<表达式>
// <复合语句> → begin<语句>{;<语句>}end
// <条件语句> → if<条件>then<语句>
// <过程调用语句> → call<标识符>
// <当循环语句> → while<条件>do<语句>
// <读语句> → read(<标识符>{,<标识符>})
// <写语句> → write(<标识符>{,<标识符>})
func (p *Parser) statement(toks []token.Token) {
	switch p.getCurSymbol().Tok {
	case token.IDENTSYM:
		// <赋值语句> → <标识符>:=<表达式>
		idName, ok := p.checkInTable(p.getCurSymbol().GetLit())
		if ok {
			if idName.identType.IsConstant() {
				// 不能改变常量的值
				p.Error(25)
				ok = false
			}
		} else {
			// 变量未定义，不能赋值
			p.Error(26)
		}
		p.goNextSymbol()
		// :=
		if p.getCurSymbol().Tok.IsBecome() {
			p.goNextSymbol()
		} else {
			p.Error(13)
		}
		// 表达式
		p.expression(toks)
		if ok {
			p.gen(fct.Sto, p.level-idName.level, idName.addr)
		}
	case token.CALLSYM:
		// <过程调用语句> → call<标识符>
		p.goNextSymbol()
		if p.getCurSymbol().Tok.IsIdent() {
			id, ok := p.checkInTable(p.getCurSymbol().GetLit())
			if ok {
				if id.identType.IsProcedure() {
					p.gen(fct.Cal, p.level-id.level, id.addr)
				} else {
					// 非过程标识符不可被调用
					p.Error(15)
				}
			} else {
				// 未找到调用的过程
				p.Error(6)
			}
			p.goNextSymbol()
		} else {
			// 不是过程调用语句
			p.Error(27)
		}
	case token.IFSYM:
		// <条件语句> → if<条件>then<语句>
		p.goNextSymbol()
		// 条件
		p.condition(append(toks, []token.Token{token.THENSYM, token.DOSYM}...))
		if p.getCurSymbol().Tok.IsThen() {
			p.goNextSymbol()
		} else {
			// 未找到then
			p.Error(29)
		}
		// TODO 看不懂
		cidx := p.cidx // 挖坑，false集的 a 需要时then后面的语句？
		p.gen(fct.Jpc, 0, 0)
		// 递归语句
		p.statement(toks)
		p.codes[cidx].A = p.cidx // false集，也就是else的部分？没有else?
	case token.BEGINSYM:
		// <复合语句> → begin<语句>{;<语句>}end
		p.goNextSymbol()
		newToks := append(toks, []token.Token{token.SEMICOLONSYM, token.ENDSYM}...)
		p.statement(newToks)
		for p.getCurSymbol().Tok.IsSemicolon() {
			p.goNextSymbol()
			p.statement(newToks)
		}
		if p.getCurSymbol().Tok.IsEnd() {
			p.goNextSymbol()
		} else {
			// 没有结束符号
			p.Error(30)
		}
	case token.WHILESYM:
		// <当循环语句> → while<条件>do<语句>
		cidx1 := p.cidx // 判断前面，循环体结束后需要跳过来
		p.goNextSymbol()
		// 条件
		p.condition(append(toks, token.DOSYM))
		cidx2 := p.cidx // 退出循环体的地址后面分配好代码后回填
		p.gen(fct.Jpc, 0, 0)
		if p.getCurSymbol().Tok.IsDo() {
			p.goNextSymbol()
		} else {
			// 缺少do
			p.Error(33)
		}
		// 语句
		p.statement(toks)
		p.gen(fct.Jmp, 0, cidx1)
		p.codes[cidx2].A = p.cidx
	case token.READSYM:
		//  <读语句> → read(<标识符>{,<标识符>})
		p.goNextSymbol()
		// 左括号
		if p.getCurSymbol().Tok.IsLparent() {
			p.goNextSymbol()
			// <标识符>
			if p.getCurSymbol().Tok.IsIdent() {
				// 检查变量表 应该是一个已经定义的变量
				id, ok := p.checkInTable(p.getCurSymbol().GetLit())
				if ok {
					if !id.identType.IsVariable() {
						// 不能改变常量的值
						p.Error(25)
						ok = false
					}
				} else {
					// 变量未定义
					p.Error(26)
				}
				if ok {
					p.gen(fct.Opr, 0, 14)                     // 读入数字放栈顶
					p.gen(fct.Sto, p.level-id.level, id.addr) // 从栈顶放到相应位置
				}
			} else {
				// 不是标识符
				p.Error(34)
			}
		} else {
			// 缺失左括号
			p.Error(35)
		}
		p.goNextSymbol()
		// {,<标识符>}
		for p.getCurSymbol().Tok.IsComma() {
			p.goNextSymbol()
			if p.getCurSymbol().Tok.IsIdent() {
				// 检查变量表 应该是一个已经定义的变量
				id, ok := p.checkInTable(p.getCurSymbol().GetLit())
				if ok {
					if id.identType.IsConstant() {
						// 不能改变常量的值
						p.Error(25)
						ok = false
					}
				} else {
					//变量未定义
					p.Error(26)
				}
				if ok {
					p.gen(fct.Opr, 0, 14)                     //读入数字放栈顶
					p.gen(fct.Sto, p.level-id.level, id.addr) //从栈顶放到相应位置
				}
			} else {
				// 不是标识符
				p.Error(34)
			}
			p.goNextSymbol()
		}
		if p.getCurSymbol().Tok.IsRparent() {
			p.goNextSymbol()
		} else {
			// 缺失右括号
			p.Error(36)
		}
	case token.WRITESYM:
		// <写语句> → write(<标识符>{，<标识符>})
		p.goNextSymbol()
		// 左括号
		if p.getCurSymbol().Tok.IsLparent() {
			p.goNextSymbol()
			// <标识符>
			if p.getCurSymbol().Tok.IsIdent() {
				// 检查变量表 应该是一个已经定义的变量
				id, ok := p.checkInTable(p.getCurSymbol().GetLit())
				if ok {
					if id.identType.IsProcedure() {
						// 不能读过程
						p.Error(28)
						ok = false
					} else if id.identType.IsConstant() {
						p.gen(fct.Lit, 0, id.value) //从相应位置读到栈顶
						p.gen(fct.Opr, 0, 15)       //从栈顶显示出来
						ok = false
					}
				} else {
					// 变量未定义
					p.Error(26)
				}
				if ok {
					p.gen(fct.Lod, p.level-id.level, id.addr) //从相应位置读到栈顶
					p.gen(fct.Opr, 0, 15)                     //从栈顶显示出来
				}
			} else {
				// 不是标识符
				p.Error(34)
			}
		} else {
			// 缺失左括号
			p.Error(35)
		}
		p.goNextSymbol()
		// {,<标识符>}
		for p.getCurSymbol().Tok.IsComma() {
			p.goNextSymbol()
			if p.getCurSymbol().Tok.IsIdent() {
				// 检查变量表 应该是一个已经定义的变量
				id, ok := p.checkInTable(p.getCurSymbol().GetLit())
				if ok {
					if id.identType.IsProcedure() {
						// 不能读过程
						p.Error(28)
						ok = false
					} else if id.identType.IsConstant() {
						p.gen(fct.Lit, 0, id.value) //从相应位置读到栈顶
						p.gen(fct.Opr, 0, 15)       //从栈顶显示出来
						ok = false
					}
				} else {
					//变量未定义
					p.Error(26)
				}
				if ok {
					p.gen(fct.Lod, p.level-id.level, id.addr) //从相应位置读到栈顶
					p.gen(fct.Opr, 0, 15)                     //从栈顶显示出来
				}
			} else {
				// 不是标识符
				p.Error(34)
			}
			p.goNextSymbol()
		}
		if p.getCurSymbol().Tok.IsRparent() {
			p.goNextSymbol()
		} else {
			// 缺失右括号
			p.Error(36)
		}
	}

	p.test(toks, 19)
}

// block
// <程序>→<分程序>.
// <分程序>→ [<常量说明部分>][<变量说明部分>][<过程说明部分>]〈语句〉
// <常量说明部分> → CONST<常量定义>{ ,<常量定义>}；
// <变量说明部分> → VAR<标识符>{ ,<标识符>}；
// <过和说明部分> → <过程首部><分程度>；{<过程说明部分>}
// <过程首部> → procedure<标识符>；
func (p *Parser) block(toks []token.Token) {
	// p.dx = 3
	tx0 := p.tx
	p.table[p.tx].addr = p.cidx
	p.gen(fct.Jmp, 0, 0)

	if p.level > MAX_LEVEL {
		// 嵌套层次太大
		p.Error(32)
	}
	//声明部分
	for {
		// 常量说明部分
		if p.getCurSymbol().Tok.IsConst() {
			// <常量说明部分> → CONST<常量定义>{,<常量定义>};
			p.goNextSymbol()
			for p.getCurSymbol().Tok.IsIdent() {
				// <常量定义>
				p.constDeclaration()
				// {,<常量定义>}
				for p.getCurSymbol().Tok.IsComma() {
					p.goNextSymbol()
					p.constDeclaration()
				}
				// 分号
				if p.getCurSymbol().Tok.IsSemicolon() {
					p.goNextSymbol()
				} else {
					// 缺少分号
					p.Error(37)
				}
			}
		}
		// 变量说明部分
		if p.getCurSymbol().Tok.IsVar() {
			// <变量说明部分> → VAR<标识符>{,<标识符>};
			p.goNextSymbol()
			for p.getCurSymbol().Tok.IsIdent() {
				p.varDeclaration()
				// 重复
				for p.getCurSymbol().Tok.IsComma() {
					p.goNextSymbol()
					p.varDeclaration()
				}
				// 分号
				if p.getCurSymbol().Tok.IsSemicolon() {
					p.goNextSymbol()
				} else {
					// 缺少分号
					p.Error(37)
				}
			}
		}
		// 过程声明部分
		// <过程说明部分> → <过程首部><分程序>;{<过程说明部分>}
		for p.getCurSymbol().Tok.IsProc() {
			// <过程首部> → procedure<标识符>;
			p.goNextSymbol()
			if p.getCurSymbol().Tok.IsIdent() {
				p.enter(Proc)
				p.goNextSymbol()
			} else {
				// 非标识符
				p.Error(39)
			}
			if p.getCurSymbol().Tok.IsSemicolon() {
				p.goNextSymbol()
			} else {
				// 过程首部里面缺少分号
				p.Error(37)
			}
			// <分程序>
			p.level++
			tx1 := p.tx
			dx1 := p.dx
			p.block(append(toks, token.SEMICOLONSYM)) // 回溯吼
			p.level--
			p.tx = tx1
			p.dx = dx1

			if p.getCurSymbol().Tok.IsSemicolon() {
				p.goNextSymbol()
				p.test(append(toks, append(StatementSelect, []token.Token{token.IDENTSYM, token.PROCSYM}...)...), 6)
			} else {
				// 缺少分号
				p.Error(37)
			}
		}

		p.test(append(StatementSelect, append(DeclareSelect, token.IDENTSYM)...), 7)
		if !tokenIsInTokens(DeclareSelect, p.getCurSymbol().Tok) {
			break
		}
	}

	p.codes[p.table[tx0].addr].A = p.cidx
	p.table[tx0].addr = p.cidx
	//cx0 := p.cidx
	p.gen(fct.Int, 0, p.dx)
	p.statement(append(toks, []token.Token{token.SEMICOLONSYM, token.ENDSYM}...))
	p.gen(fct.Opr, 0, 0)
	p.test(toks, 8)
}

func (p *Parser) Lex() {
	p.lexer.GetSym()
}

// Parse 语法分析入口
func (p *Parser) Parse() {
	// p.lexer.GetSym()
	p.block(append(append(DeclareSelect, StatementSelect...), token.PERIODSYM))
	// log.Printf("结束\n")
}

// SaveLexResult 保存词法分析结果
func (p *Parser) SaveLexResult(path string) {
	p.lexer.Save(path)
}
