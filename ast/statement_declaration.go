package ast

type DeclarationStatement struct {
	StatementBase
	Name  *Identifier
	Type  Type
	Value Expression
}

func (*DeclarationStatement) GenerateIR(*Context) bool {
	//TO-DO // zero initialize
	//Builtin, TypeName
	return false
}
