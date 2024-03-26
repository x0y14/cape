package python

import (
	"fmt"
	"strconv"
	"strings"
)

var nest int

type line struct {
	s string // str
	n int    // nest
}

func newLine(s string, n int) *line {
	return &line{s: s, n: n}
}

func genLine(lines []*line) string {
	var s string
	for _, l := range lines {
		s += fmt.Sprintf("%s%s\n", strings.Repeat("    ", l.n), l.s)
	}
	return s
}

func Gen(nodes []*Node) (string, error) {
	lines, err := genProgram(nodes)
	if err != nil {
		return "", err
	}
	code := genLine(lines)
	code += "if __name__ == \"__main__\":\n    main()"
	return code, nil
}

func genProgram(nodes []*Node) ([]*line, error) {
	var lines []*line
	for _, node := range nodes {
		l, err := genToplevel(node)
		if err != nil {
			return nil, err
		}
		lines = append(lines, l...)
	}
	return lines, nil
}

func genToplevel(node *Node) ([]*line, error) {
	switch node.GetKind() {
	case FunctionDefine:
		return genFunctionDefine(node)
	default:
		panic("unhandled toplevel case") // TODO
	}
	return nil, nil
}

func genFunctionDefine(node *Node) ([]*line, error) {
	field := node.GetField().(*FunctionDefineField)

	identField := field.Ident.GetField().(*IdentField)
	ident := identField.S
	params, err := genFunctionDefineParams(field.Params)
	if err != nil {
		return nil, err
	}
	block, err := genStmt(field.Block)
	if err != nil {
		return nil, err
	}

	var lines []*line
	lines = append(lines, newLine(fmt.Sprintf("def %s(%s):", ident, params), 0))
	lines = append(lines, block...)

	return lines, nil
}

func genFunctionDefineParams(node *Node) (string, error) { // TODO
	_ = node
	return "", nil
}

func genStmt(node *Node) ([]*line, error) {
	switch node.GetKind() {
	case Block:
		var lines []*line
		nest++
		blockField := node.GetField().(*BlockField)
		for _, stmtNode := range blockField.Stmts {
			statements, err := genStmt(stmtNode)
			if err != nil {
				return nil, err
			}
			lines = append(lines, statements...)
		}
		nest--
		return lines, nil
	case Return:
		rvField := node.GetField().(*ReturnField)
		rv, err := genExpr(rvField.Value)
		if err != nil {
			return nil, err
		}
		return []*line{newLine(fmt.Sprintf("return %s", rv), nest)}, nil
		//var values []string
		//rvsField := rvField.Value.GetField().(*MultipleField)
		//for _, valueNode := range rvsField.Values {
		//	v, err := genExpr(valueNode)
		//	if err != nil {
		//		return nil, err
		//	}
		//	values = append(values, v)
		//}
		//return []*line{newLine(fmt.Sprintf("return %s", strings.Join(values, ", ")), nest)}, nil
	case IfElse:
		panic("unimplemented if case")
	case While:
		panic("unimplemented while case")
	case For:
		panic("unimplemented for case")
	default:
		e, err := genExpr(node)
		if err != nil {
			return nil, err
		}
		return []*line{newLine(e, nest)}, nil
	}
	return nil, nil
}

func genExpr(node *Node) (string, error) {
	return genAssign(node)
}

func genAssign(node *Node) (string, error) {
	switch node.GetKind() {
	case Assign:
		assignField := node.GetField().(*AssignField)
		to, err := genExpr(assignField.To)
		if err != nil {
			return "", err
		}
		val, err := genExpr(assignField.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", to, val), nil
	default:
		return genAndor(node)
	}
}

func genAndor(node *Node) (string, error) {
	if node.GetKind() != Binary {
		return genEquality(node)
	}
	binaryField := node.GetField().(*BinaryField)
	lhs, err := genAndor(binaryField.LHS)
	if err != nil {
		return "", err
	}
	rhs, err := genAndor(binaryField.RHS)
	if err != nil {
		return "", err
	}
	switch binaryField.Operation {
	case And:
		return fmt.Sprintf("%s and %s", lhs, rhs), nil
	case Or:
		return fmt.Sprintf("%s or %s", lhs, rhs), nil
	default:
		panic("unimplemented binary op")
	}
	return "", nil
}

func genEquality(node *Node) (string, error) {
	if node.GetKind() != Binary {
		return genRelational(node)
	}

	binaryField := node.GetField().(*BinaryField)
	lhs, err := genEquality(binaryField.LHS)
	if err != nil {
		return "", err
	}
	rhs, err := genEquality(binaryField.RHS)
	if err != nil {
		return "", err
	}
	switch binaryField.Operation {
	case Eq:
		return fmt.Sprintf("%s == %s", lhs, rhs), nil
	case Ne:
		return fmt.Sprintf("%s != %s", lhs, rhs), nil
	default:
		panic("unimplemented binary op")
	}

	return "", nil
}

func genRelational(node *Node) (string, error) {
	if node.GetKind() != Binary {
		return genAdd(node)
	}
	binaryField := node.GetField().(*BinaryField)
	lhs, err := genRelational(binaryField.LHS)
	if err != nil {
		return "", err
	}
	rhs, err := genRelational(binaryField.RHS)
	if err != nil {
		return "", err
	}
	switch binaryField.Operation {
	case Lt:
		return fmt.Sprintf("%s < %s", lhs, rhs), err
	case Le:
		return fmt.Sprintf("%s <= %s", lhs, rhs), err
	case Gt:
		return fmt.Sprintf("%s > %s", lhs, rhs), err
	case Ge:
		return fmt.Sprintf("%s >= %s", lhs, rhs), err
	default:
		panic("unimplemented binary op")
	}
	return "", nil
}

func genAdd(node *Node) (string, error) {
	switch node.GetKind() {
	case Binary:
		binaryField := node.GetField().(*BinaryField)
		lhs, err := genAdd(binaryField.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := genAdd(binaryField.RHS)
		if err != nil {
			return "", err
		}
		switch binaryField.Operation {
		case Add:
			return fmt.Sprintf("%s + %s", lhs, rhs), nil
		case Sub:
			return fmt.Sprintf("%s - %s", lhs, rhs), nil
		default:
			panic("unimplemented binary op")
		}
	default:
		return genMul(node)
	}
	return "", nil
}

func genMul(node *Node) (string, error) {
	switch node.GetKind() {
	case Binary:
		binaryField := node.GetField().(*BinaryField)
		lhs, err := genMul(binaryField.LHS)
		if err != nil {
			return "", err
		}
		rhs, err := genMul(binaryField.RHS)
		if err != nil {
			return "", err
		}
		switch binaryField.Operation {
		case Mul:
			return fmt.Sprintf("%s * %s", lhs, rhs), nil
		case Div:
			return fmt.Sprintf("%s / %s", lhs, rhs), nil
		case Mod:
			return fmt.Sprintf("%s %% %s", lhs, rhs), nil
		default:
			panic("unimplemented binary op")
		}
	default:
		return genUnary(node)
	}
	return "", nil
}

func genUnary(node *Node) (string, error) {
	return genPrimary(node)
}

func genPrimary(node *Node) (string, error) {
	switch node.GetKind() {
	case Ident:
		identField := node.GetField().(*IdentField)
		ident := identField.S
		return ident, nil
	case Call:
		return genCall(node)
	case Literal:
		return genLiteral(node)
	default:
		panic("unexpected primary")
	}
	return "", nil
}

func genLiteral(node *Node) (string, error) {
	literalField := node.GetField().(*LiteralField)
	switch literalField.GetTType() {
	case Integer:
		return strconv.Itoa(literalField.I), nil
	case String:
		return literalField.S, nil
	case Bool:
		panic("unimplemented literal")
	default:
		panic("unimplemented literal")
	}
	return "", nil
}

func genCall(node *Node) (string, error) {
	callField := node.GetField().(*CallField)
	identField := callField.Ident.GetField().(*IdentField)
	ident := identField.S
	argsField := callField.Args.GetField().(*MultipleField)
	args := argsField.Values

	var code string
	if ident == "printf" {
		var arguments string
		for i, arg := range args[1:] {
			a, err := genExpr(arg)
			if err != nil {
				return "", err
			}
			if i != 0 {
				arguments += ", "
			}
			arguments += a
		}

		f, err := genExpr(args[0])
		if err != nil {
			return "", err
		}
		if arguments == "" {
			code = fmt.Sprintf("print(%#v)", f)
		} else {
			code = fmt.Sprintf("print(%#v.format(%s))", formatting(f), arguments)
		}
	} else {
		var arguments string
		for i, arg := range args {
			a, err := genExpr(arg)
			if err != nil {
				return "", err
			}
			if i != 0 {
				arguments += ", "
			}
			arguments += a
		}
		code = fmt.Sprintf("%s(%s)", ident, arguments)
	}
	return code, nil
}

func formatting(s string) string {
	s = strings.ReplaceAll(s, "%d", "{}")
	s = strings.ReplaceAll(s, "%s", "{}")
	return s
}
