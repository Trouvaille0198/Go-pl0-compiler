package token

import (
	"strconv"
	"unicode"
)

func init() {
	initKwHash()
	initOptHash()
}

// Token 符号枚举编号
type Token int

// 符号枚举
const (
	BADTOKEN = Token(iota) // 无效字符

	literal_beg
	IDENTSYM  // 标识符
	NUMBERSYM // 数
	literal_end

	operator_beg
	// Operators
	PLUSSYM  // +
	MINUSYM  // -
	MULSYM   // *
	SLASHSYM // /

	relation_optr_beg
	EQLSYM // =
	NEQSYM // #
	LESSYM // <
	LEQSYM // <=
	GTRSYM // >
	GEQSYM // >=
	relation_optr_end

	LPARENTSYM   // (
	RPARENTSYM   // )
	COMMASYM     // ,
	SEMICOLONSYM // ;
	PERIODSYM    // .
	BECOMESSYM   // :=
	operator_end

	keyword_beg
	// Keywords
	BEGINSYM // begin
	ENDSYM   // end
	IFSYM    // if
	ELSESYM  // else
	THENSYM  // then
	WHILESYM // while
	DOSYM    // do
	CALLSYM  // call
	CONSTSYM // const
	VARSYM   // var
	PROCSYM  // procedure
	ODDSYM   // odd
	WRITESYM // write
	READSYM  // read
	keyword_end

	EOFSYM // EOF
)

// 判别字符类型

func (t Token) IsLiteral() bool     { return literal_beg < t && t < literal_end }
func (t Token) IsIdent() bool       { return t == IDENTSYM }
func (t Token) IsNumber() bool      { return t == NUMBERSYM }
func (t Token) IsBecome() bool      { return t == BECOMESSYM }
func (t Token) IsCall() bool        { return t == CALLSYM }
func (t Token) IsThen() bool        { return t == THENSYM }
func (t Token) IsSemicolon() bool   { return t == SEMICOLONSYM }
func (t Token) IsEnd() bool         { return t == ENDSYM }
func (t Token) IsDo() bool          { return t == DOSYM }
func (t Token) IsConst() bool       { return t == CONSTSYM }
func (t Token) IsComma() bool       { return t == COMMASYM }
func (t Token) IsVar() bool         { return t == VARSYM }
func (t Token) IsProc() bool        { return t == PROCSYM }
func (t Token) IsLparent() bool     { return t == LPARENTSYM }
func (t Token) IsRparent() bool     { return t == RPARENTSYM }
func (t Token) IsOperator() bool    { return operator_beg < t && t < operator_end }
func (t Token) IsKeyword() bool     { return keyword_beg < t && t < keyword_end }
func (t Token) IsRelationOpr() bool { return relation_optr_beg < t && t < relation_optr_end }

// 符号枚举在程序中的表示
var tokens = [...]string{
	CONSTSYM: "const",
	VARSYM:   "var",
	PROCSYM:  "procedure",
	CALLSYM:  "call",
	BEGINSYM: "begin",
	ENDSYM:   "end",
	IFSYM:    "if",
	THENSYM:  "then",
	ELSESYM:  "else",
	WHILESYM: "while",
	DOSYM:    "do",
	READSYM:  "read",
	WRITESYM: "write",
	ODDSYM:   "odd",

	PLUSSYM:    "+",
	MINUSYM:    "-",
	MULSYM:     "*",
	SLASHSYM:   "/",
	LPARENTSYM: "(",
	RPARENTSYM: ")",
	COMMASYM:   ",",
	PERIODSYM:  ".",

	EQLSYM: "=",
	LESSYM: "<",
	GTRSYM: ">",
	NEQSYM: "#",
	LEQSYM: "<=",
	GEQSYM: ">=",

	SEMICOLONSYM: ";",
	BECOMESSYM:   ":=",

	// 以下符号不会在程序中出现
	BADTOKEN:  "无效字符",
	NUMBERSYM: "数字",
	IDENTSYM:  "变量标识符",

	EOFSYM: "文档已结束",
}

// 符号枚举的说明
var tokenDesc = [...]string{
	CONSTSYM: "constsym",
	VARSYM:   "varsym",
	PROCSYM:  "proceduresym",
	CALLSYM:  "callsym",
	BEGINSYM: "beginsym",
	ENDSYM:   "endsym",
	IFSYM:    "ifsym",
	THENSYM:  "thensym",
	ELSESYM:  "elsesym",
	WHILESYM: "whilesym",
	DOSYM:    "dosym",
	READSYM:  "readsym",
	WRITESYM: "writesym",
	ODDSYM:   "oddsym",

	PLUSSYM:    "plus",
	MINUSYM:    "minus",
	MULSYM:     "times",
	SLASHSYM:   "slash",
	LPARENTSYM: "lparen",
	RPARENTSYM: "rparen",
	COMMASYM:   "comma",
	PERIODSYM:  "period",

	EQLSYM: "eql",
	LESSYM: "lss",
	GTRSYM: "gtr",
	NEQSYM: "neq",
	LEQSYM: "leq",
	GEQSYM: "geq",

	SEMICOLONSYM: "semicolon",
	BECOMESSYM:   "becomes",

	// 以下符号不会再程序中出现
	BADTOKEN:  "无效字符",
	NUMBERSYM: "number",
	IDENTSYM:  "ident",

	EOFSYM: "EOF",
}

// String 返回Token的字符串描述
func (t Token) String() string {
	if 0 <= t && t < Token(len(tokenDesc)) {
		return tokenDesc[t]
	} else {
		return "token(" + strconv.Itoa(int(t)) + ")"
	}
}

// GetDesc 返回Token的字符串描述
func (t Token) GetDesc() string {
	return t.String()
}

// GetLit 返回Token在程序中的字面量表示
func (t Token) GetLit() string {
	if 0 <= t && t < Token(len(tokens)) {
		return tokens[t]
	} else {
		return "token(" + strconv.Itoa(int(t)) + ")"
	}
}

var KeyWordHash map[string]Token  // 关键字哈希表
var OperatorHash map[string]Token // 操作符哈希表

// initKwHash 初始化关键字哈希表
func initKwHash() {
	KeyWordHash = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		KeyWordHash[tokens[i]] = i
	}
}

// initOptHash 初始化操作符哈希表
func initOptHash() {
	OperatorHash = make(map[string]Token)
	for i := operator_beg + 1; i < operator_end; i++ {
		if i == relation_optr_beg || i == relation_optr_end {
			continue
		}
		OperatorHash[tokens[i]] = i
	}
}

// IsIdentifier 判断name是否是一个合法的标识符
func IsIdentifier(name string) bool {
	for i, char := range name {
		if !unicode.IsLetter(char) && (i == 0 || !unicode.IsDigit(char)) {
			// 为非数字或字母 或首字母为数字
			return false
		}
	}
	if _, ok := KeyWordHash[name]; ok {
		// 不能是关键字
		return false
	}
	if name == "" {
		// 不能为空
		return false
	}
	return true
}

// GetIdToken 匹配关键字的Token 若不存在 则视为标识符Token
func GetIdToken(ident string) Token {
	if tok, isKw := KeyWordHash[ident]; isKw {
		return tok
	}
	return IDENTSYM
}

// GetOptToken 判断并匹配运算符的Token
func GetOptToken(opt string) (Token, bool) {
	if tok, isOpt := OperatorHash[opt]; isOpt {
		return tok, true
	}
	return BADTOKEN, false
}
