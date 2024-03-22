package interlang

type NodeKind int

const (
	_ NodeKind = iota
	VariableDeclare
	FunctionDeclare
	VariableDefine
	FunctionDefine

	Block
	IfElse
	While
	For
	Assign
	Binary
	Literal
	Not
	Multiple
	Return
	Call

	Ident

	Add
	Sub
	Mul
	Div
	Mod
	AND
	OR
	Eq
	Ne
	Lt
	Le
	Gt
	Ge
)

type Node struct {
	NodeKind
	Field
}

func (n *Node) GetKind() NodeKind {
	return n.NodeKind
}
func (n *Node) GetField() Field {
	return n.Field
}

func NewNode(kind NodeKind, field Field) *Node {
	return &Node{kind, field}
}
