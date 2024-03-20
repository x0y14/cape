package python

type TType interface {
	IsEqual(tt2 TType) bool
}

type TPrimitive int

const (
	_ TPrimitive = iota
	Null
	Integer
	String
	Bool
)

func (tt TPrimitive) IsEqual(tt2 TType) bool { // TODO
	_ = tt2
	return false
}

type TTuple []TType

func (tt TTuple) IsEqual(tt2 TType) bool { // TODO
	_ = tt2
	return false
}
