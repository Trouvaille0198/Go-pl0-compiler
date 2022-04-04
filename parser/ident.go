package parser

// IdentType 标识符类型
type IdentType int

const (
	Constant = iota // 常量
	Variable        // 变量
	Proc            // 过程
)

var identHash = [...]string{
	Constant: "常量",
	Variable: "变量",
	Proc:     "过程",
}

func (i IdentType) String() string {
	return identHash[i]
}

func (i IdentType) IsVariable() bool  { return i == Variable }
func (i IdentType) IsConstant() bool  { return i == Constant }
func (i IdentType) IsProcedure() bool { return i == Proc }
