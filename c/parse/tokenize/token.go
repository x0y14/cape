package tokenize

type TokenKind int

const (
	_ TokenKind = iota
	Eof
	Ident
	Int
	String
	Float

	Lrb
	Rrb
	Lsb
	Rsb
	Lcb
	Rcb

	Dot
	Comma

	Colon
	Semi

	Add
	Sub
	Mul
	Div
	Mod

	Eq
	Ne
	Gt
	Lt
	Ge
	Le

	Assign

	And
	Or
	Not
)

type Token struct {
	Kind TokenKind
	S    string
	I    int
	F    float64
	Next *Token
}

func newToken(kind TokenKind, s string, i int, f float64) *Token {
	return &Token{
		Kind: kind,
		S:    s,
		I:    i,
		F:    f,
		Next: nil,
	}
}

func newLiteralToken[T int | float64 | string](v T) *Token {
	switch any(v).(type) {
	case int:
		return newToken(Int, "", any(v).(int), 0)
	case float64:
		return newToken(Float, "", 0, any(v).(float64))
	case string:
		return newToken(String, any(v).(string), 0, 0)
	}
	return nil
}

func newSymbolToken(syb string) *Token {
	switch syb {
	case "(":
		return newToken(Lrb, "", 0, 0)
	case ")":
		return newToken(Rrb, "", 0, 0)
	case "[":
		return newToken(Lsb, "", 0, 0)
	case "]":
		return newToken(Rsb, "", 0, 0)
	case "{":
		return newToken(Lcb, "", 0, 0)
	case "}":
		return newToken(Rcb, "", 0, 0)
	case ".":
		return newToken(Dot, "", 0, 0)
	case ",":
		return newToken(Comma, "", 0, 0)
	case ":":
		return newToken(Colon, "", 0, 0)
	case ";":
		return newToken(Semi, "", 0, 0)
	case "+":
		return newToken(Add, "", 0, 0)
	case "-":
		return newToken(Sub, "", 0, 0)
	case "*":
		return newToken(Mul, "", 0, 0)
	case "/":
		return newToken(Div, "", 0, 0)
	case "%":
		return newToken(Mod, "", 0, 0)
	case "<":
		return newToken(Lt, "", 0, 0)
	case ">":
		return newToken(Gt, "", 0, 0)
	case "=":
		return newToken(Assign, "", 0, 0)
	case "!":
		return newToken(Not, "", 0, 0)
	case "==":
		return newToken(Eq, "", 0, 0)
	case "!=":
		return newToken(Ne, "", 0, 0)
	case "<=":
		return newToken(Le, "", 0, 0)
	case ">=":
		return newToken(Ge, "", 0, 0)
	case "&&":
		return newToken(And, "", 0, 0)
	case "||":
		return newToken(Or, "", 0, 0)
	default:
		return nil
	}
}
