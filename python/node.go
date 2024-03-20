package python

type NodeKind int

const (
	_ NodeKind = iota
	FunctionDefine

	Block
	Multiple
	Return
	IfElse
	While
	For
	Assign
	Binary
	Literal
	Not
	Call
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

func NewNode(kind NodeKind, field Field) Node {
	return Node{kind, field}
}
