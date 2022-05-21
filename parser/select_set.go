package parser

import "gopl0/token"

var (
	DeclareSelect = []token.Token{
		token.CONSTSYM,
		token.VARSYM,
		token.PROCSYM,
	}
	ExpressionSelect = []token.Token{
		token.EQLSYM,
		token.NEQSYM,
		token.LESSYM,
		token.LEQSYM,
		token.GTRSYM,
		token.GEQSYM,
	}
	FactorSelect = []token.Token{
		token.IDENTSYM,
		token.NUMBERSYM,
		token.LPARENTSYM,
	}
	StatementSelect = []token.Token{
		token.BEGINSYM,
		token.CALLSYM,
		token.IFSYM,
		token.WHILESYM,
	}
)
