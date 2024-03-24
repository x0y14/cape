package interlang

type Operation int

const (
	_ Operation = iota
	Add
	Sub
	Mul
	Div
	Mod

	And
	Or

	Eq
	Ne

	Lt
	Le
	Gt
	Ge
)
