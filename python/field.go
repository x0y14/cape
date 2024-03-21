package python

type FieldKind = NodeKind

type Field interface {
	GetKind() FieldKind
}

type FunctionDefineField struct {
	TType
	Ident  *Node
	Params *Node
	Block  *Node
}

func (f *FunctionDefineField) GetKind() FieldKind {
	return FunctionDefine
}
func (f *FunctionDefineField) GetTType() TType {
	return f.TType
}

type BlockField struct {
	Stmts []*Node
}

func (f *BlockField) GetKind() FieldKind {
	return Block
}

type MultipleField struct {
	TType
	Values []*Node
}

func (f *MultipleField) GetKind() FieldKind {
	return Multiple
}
func (f *MultipleField) GetTType() TType {
	return f.TType
}

type ReturnField struct {
	TType
	Value *Node
}

func (f *ReturnField) GetKind() FieldKind {
	return Return
}
func (f *ReturnField) GetTType() TType {
	return f.TType
}

type IfElseField struct {
	Cond      *Node
	IfBlock   *Node
	ElseBlock *Node
}

func (f *IfElseField) GetKind() FieldKind {
	return IfElse
}

type WhileField struct {
	Block *Node
}

func (f *WhileField) GetKind() FieldKind {
	return While
}

type ForField struct {
	Init  *Node
	Cond  *Node
	Loop  *Node
	Block *Node
}

func (f *ForField) GetKind() FieldKind {
	return For
}

type AssignField struct {
	To    *Node
	Value *Node
}

func (f *AssignField) GetKind() FieldKind {
	return Assign
}

type BinaryField struct {
	LHS *Node
	RHS *Node
}

func (f *BinaryField) GetKind() FieldKind {
	return Binary
}
func (f *BinaryField) GetTType() TType {
	return f.GetTType()
}

type LiteralField struct {
	TType
	I int
	F float64
	S string
}

func (f *LiteralField) GetKind() FieldKind {
	return Literal
}
func (f *LiteralField) GetTType() TType {
	return f.TType
}

type NotField struct {
	Value *Node
}

func (f *NotField) GetKind() FieldKind {
	return Not
}
func (f *NotField) GetTType() TType {
	return Bool
}

type CallField struct {
	TType
	Ident *Node
	Args  *Node
}

func (f *CallField) GetKind() FieldKind {
	return Call
}
func (f *CallField) GetTType() TType {
	return f.TType
}
