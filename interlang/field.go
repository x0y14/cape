package interlang

type FieldKind = NodeKind

type Field interface {
	GetKind() FieldKind
	//GetTType() TType
}

type VariableDeclareField struct {
	TType
	Ident *Node
}

func (f *VariableDeclareField) GetKind() FieldKind {
	return VariableDeclare
}
func (f *VariableDeclareField) GetTType() TType {
	return f.TType
}

type FunctionDeclareField struct {
	TType
	Ident  *Node
	Params *Node
}

func (f *FunctionDeclareField) GetKind() FieldKind {
	return FunctionDeclare
}
func (f *FunctionDeclareField) GetTType() TType {
	return f.TType
}

type VariableDefineField struct {
	TType
	Ident *Node
	Value *Node
}

func (f *VariableDefineField) GetKind() FieldKind {
	return VariableDefine
}
func (f *VariableDefineField) GetTType() TType {
	return f.TType
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

type IfElseField struct {
	Cond      *Node
	IfBlock   *Node
	ElseBlock *Node
}

func (f *IfElseField) GetKind() FieldKind {
	return IfElse
}

type WhileField struct {
	Cond  *Node
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
	TType
	Operation
	LHS *Node
	RHS *Node
}

func (f *BinaryField) GetKind() FieldKind {
	return Binary
}
func (f *BinaryField) GetTType() TType {
	return f.TType
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
