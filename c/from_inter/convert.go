package from_inter

import (
	"cape/c"
	"cape/interlang"
	"log"
)

func ConvertNodeFromInterLang(iNodes []*interlang.Node) ([]*c.Node, error) {
	var nodes []*c.Node
	for _, in := range iNodes {
		cn, err := toplevel(in)
		if err != nil {
			return nil, err
		}
		if cn == nil {
			continue
		}
		nodes = append(nodes, cn)
	}
	return nodes, nil
}

func convertTypeFromInterLang(tt interlang.TType) (c.TType, error) {
	switch tt {
	case interlang.Integer:
		return c.Integer, nil
	case interlang.String:
		return c.String, nil
	case interlang.Bool:
		return c.Bool, nil
	default:
		return nil, nil
	}
}

func toplevel(iNode *interlang.Node) (*c.Node, error) {
	switch iNode.GetKind() {
	case interlang.FunctionDefine:
		return functionDefine(iNode)
	default:
		return nil, nil
	}
}

func variableDeclare(iNode *interlang.Node) (*c.Node, error) {
	return nil, nil
}

func functionDeclare(iNode *interlang.Node) (*c.Node, error) {
	return nil, nil
}

func variableDefine(iNode *interlang.Node) (*c.Node, error) {
	return nil, nil
}

func functionDefine(iNode *interlang.Node) (*c.Node, error) {
	iField := iNode.GetField().(*interlang.FunctionDefineField)

	iRVType := iField.GetTType()
	RVType, err := convertTypeFromInterLang(iRVType)
	if err != nil {
		return nil, err
	}

	iIdentField := iField.Ident.GetField().(*interlang.IdentField)
	ident := iIdentField.S

	params, err := functionDefineParams(iField.Params)
	if err != nil {
		return nil, err
	}

	stmts, err := statement(iField.Block)
	if err != nil {
		return nil, err
	}

	return c.NewNode(
		c.FunctionDefine,
		&c.FunctionDefineField{
			TType:  RVType,
			Ident:  c.NewNode(c.Ident, &c.IdentField{S: ident}),
			Params: params,
			Block:  stmts,
		},
	), nil
}

func functionDefineParams(iNode *interlang.Node) (*c.Node, error) {
	_ = iNode
	return nil, nil
}

func statement(iNode *interlang.Node) (*c.Node, error) {
	switch iNode.GetKind() {
	case interlang.Block:
		var stmts []*c.Node
		iStmts := iNode.GetField().(*interlang.BlockField).Stmts
		for _, iStmt := range iStmts {
			stmt, err := statement(iStmt)
			if err != nil {
				continue
			}
			stmts = append(stmts, stmt)
		}
		return c.NewNode(c.Block, &c.BlockField{Stmts: stmts}), nil
	case interlang.Return:
		iReturnField := iNode.GetField().(*interlang.ReturnField)
		rv, err := expr(iReturnField.Value)
		if err != nil {
			return nil, err
		}
		return c.NewNode(c.Return, &c.ReturnField{Value: rv}), nil
	case interlang.IfElse:
		// TODO
	case interlang.While:
		// TODO
	case interlang.For:
		// TODO
	default:
		return expr(iNode)
	}
	return nil, nil
}

func expr(iNode *interlang.Node) (*c.Node, error) {
	return assign(iNode)
}

func assign(iNode *interlang.Node) (*c.Node, error) {
	switch iNode.GetKind() {
	case interlang.VariableDeclare:
		// TODO
		return nil, nil
	case interlang.Assign:
		iAssignField := iNode.GetField().(*interlang.AssignField)
		to, err := expr(iAssignField.To)
		if err != nil {
			return nil, err
		}
		value, err := expr(iAssignField.Value)
		if err != nil {
			return nil, err
		}
		return c.NewNode(c.Assign, &c.AssignField{
			To:    to,
			Value: value,
		}), nil
	default:
		return andor(iNode)
	}
}

func andor(iNode *interlang.Node) (*c.Node, error) {
	if iNode.GetKind() != interlang.Binary {
		return equality(iNode)
	}

	iBinaryField := iNode.GetField().(*interlang.BinaryField)
	iType := iBinaryField.GetTType()
	pType, err := convertTypeFromInterLang(iType)
	if err != nil {
		return nil, err
	}
	lhs, err := andor(iBinaryField.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := andor(iBinaryField.RHS)
	if err != nil {
		return nil, err
	}

	switch iBinaryField.Operation {
	case interlang.And:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.And,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Or:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Or,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	default:
		return equality(iNode)
	}
}

func equality(iNode *interlang.Node) (*c.Node, error) {
	if iNode.GetKind() != interlang.Binary {
		return relational(iNode)
	}
	iBinaryField := iNode.GetField().(*interlang.BinaryField)
	iType := iBinaryField.GetTType()
	pType, err := convertTypeFromInterLang(iType)
	if err != nil {
		return nil, err
	}
	lhs, err := equality(iBinaryField.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := equality(iBinaryField.RHS)
	if err != nil {
		return nil, err
	}
	switch iBinaryField.Operation {
	case interlang.Eq:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Eq,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Ne:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Ne,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	default:
		return relational(iNode)
	}
}

func relational(iNode *interlang.Node) (*c.Node, error) {
	if iNode.GetKind() != interlang.Binary {
		return add(iNode)
	}
	iBinaryField := iNode.GetField().(*interlang.BinaryField)
	iType := iBinaryField.GetTType()
	pType, err := convertTypeFromInterLang(iType)
	if err != nil {
		return nil, err
	}
	lhs, err := relational(iBinaryField.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := relational(iBinaryField.RHS)
	if err != nil {
		return nil, err
	}
	switch iBinaryField.Operation {
	case interlang.Lt:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Lt,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Le:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Le,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Gt:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Gt,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Ge:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Ge,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	default:
		return add(iNode)
	}
}

func add(iNode *interlang.Node) (*c.Node, error) {
	if iNode.GetKind() != interlang.Binary {
		return mul(iNode)
	}
	iBinaryField := iNode.GetField().(*interlang.BinaryField)
	iType := iBinaryField.GetTType()
	pType, err := convertTypeFromInterLang(iType)
	if err != nil {
		return nil, err
	}
	lhs, err := add(iBinaryField.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := add(iBinaryField.RHS)
	if err != nil {
		return nil, err
	}
	switch iBinaryField.Operation {
	case interlang.Add:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Add,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Sub:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Sub,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	default:
		return mul(iNode)
	}
}

func mul(iNode *interlang.Node) (*c.Node, error) {
	if iNode.GetKind() != interlang.Binary {
		return unary(iNode)
	}
	iBinaryField := iNode.GetField().(*interlang.BinaryField)
	iType := iBinaryField.GetTType()
	pType, err := convertTypeFromInterLang(iType)
	if err != nil {
		return nil, err
	}
	lhs, err := mul(iBinaryField.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := mul(iBinaryField.RHS)
	if err != nil {
		return nil, err
	}
	switch iBinaryField.Operation {
	case interlang.Mul:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Mul,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Div:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Div,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	case interlang.Mod:
		return c.NewNode(c.Binary, &c.BinaryField{
			TType:     pType,
			Operation: c.Mod,
			LHS:       lhs,
			RHS:       rhs,
		}), nil
	default:
		return unary(iNode)
	}
}

func unary(iNode *interlang.Node) (*c.Node, error) {
	return primary(iNode)
}

func primary(iNode *interlang.Node) (*c.Node, error) {
	switch iNode.GetKind() {
	case interlang.Ident:
		iIdentField := iNode.GetField().(*interlang.IdentField)
		return c.NewNode(c.Ident, &c.IdentField{TType: c.String, S: iIdentField.S}), nil
	case interlang.Literal:
		return literal(iNode)
	case interlang.Call:
		log.Panicf("call unimplemented: %v", iNode)
	default:
		log.Panicf("unexpected primary node: %v", iNode)
	}
	return nil, nil
}

func literal(iNode *interlang.Node) (*c.Node, error) {
	iLitField := iNode.GetField().(*interlang.LiteralField)
	switch iLitField.GetTType() {
	case interlang.String:
		return c.NewNode(c.Literal, &c.LiteralField{TType: c.String, S: iLitField.S}), nil
	case interlang.Integer:
		return c.NewNode(c.Literal, &c.LiteralField{TType: c.Integer, I: iLitField.I}), nil
	case interlang.Bool:
		log.Panicf("unsupported literal: %v", iNode)
	default:
		log.Panicf("unsupported literal: %v", iNode)
	}
	return nil, nil
}
