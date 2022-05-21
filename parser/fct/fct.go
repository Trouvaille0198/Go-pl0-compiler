package fct

// Fct 功能码
type Fct int

const (
	Lit = iota //将常数放到栈顶
	Opr        // 关系和算术运算
	Lod        // 将变量放到栈顶
	Sto        // 将栈顶的内容送到某变量单元中
	Cal        // 调用过程
	Int        // 为被调用的过程（或主程序）在运行栈中开辟数据区
	Jmp        // 无条件转移
	Jpc        // 条件转移
)

var fctHash = [...]string{
	Lit: "LIT",
	Opr: "OPR",
	Lod: "LOD",
	Sto: "STO",
	Cal: "CAL",
	Int: "INT",
	Jmp: "JMP",
	Jpc: "JPC",
}

func (f Fct) String() string {
	return fctHash[f]
}

/*
	lit 0, a : load constant a
    opr 0, a : execute operation a
    lod l, a : load variable l, a
    sto l, a : store variable l, a
    cal l, a : call procedure a at level l
    Int 0, a : increment t-register by a
    jmp 0, a : jump to a
    jpc 0, a : jump conditional to a
*/
