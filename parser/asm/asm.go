package asm

import (
	"gopl0/parser/fct"
	"strconv"
)

type Asm struct {
	Fct fct.Fct // 操作符
	L   int     // 层次差
	A   int     // 位移量
}

func (a *Asm) String() string {
	return "(" + a.Fct.String() + "," + strconv.Itoa(a.L) + "," + strconv.Itoa(a.A) + ")"
}
