package symbol

import (
	"gopl0/token"
	"strconv"
)

// Symbol 符号
type Symbol struct {
	Tok   token.Token // 符号枚举编号
	Value []rune      // 用户自定义的标识符值(若有)
	Num   int         // 用户自定义的数值(若有)
	Line  int         // 符号所在行
}

// String 将Symbol对象字符串化
func (s *Symbol) String() string {
	if s.Tok.IsIdent() {
		return "(indent, " + string(s.Value) + ")"
	}
	if s.Tok.IsNumber() {
		return "(number, " + strconv.Itoa(s.Num) + ")"
	}
	return "(" + s.Tok.GetDesc() + ", " + s.Tok.GetLit() + ")"
}

// GetDesc 返回符号的描述
func (s *Symbol) GetDesc() string {
	if s.Tok.IsIdent() {
		return "indent"
	}
	if s.Tok.IsNumber() {
		return "number"
	}
	return s.Tok.GetDesc()
}

// GetLit 返回符号的字面量
func (s *Symbol) GetLit() string {
	if s.Tok.IsIdent() {
		return string(s.Value)
	}
	if s.Tok.IsNumber() {
		return strconv.Itoa(s.Num)
	}
	return s.Tok.GetLit()
}
