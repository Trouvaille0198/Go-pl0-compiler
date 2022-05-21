package expparser

type sym struct {
	typ  string // const, var
	val  int    // int. constant if typ == const
	addr int    // address of symbol
}
type symtab map[string]sym // symbol table for a single scope

type procedure struct {
	addr   int
	nlocal int
	sym    symtab
}

var scopes = make(map[string]*procedure) // scopes: procedure, global
var active = ""
