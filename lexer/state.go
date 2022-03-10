package lexer

type State int

// DFA自动机状态
const (
	START = State(iota)
	INNUM
	INID // 标识符或关键字
	INBECOMES
	BECOMES
	GTR
	GEQ
	LES
	LEQ
	END
	COMMENT
)
