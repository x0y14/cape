package python

import (
	"cape/interlang"
	"log"
)

func ConvertNodeFromInterLang(iNodes []*interlang.Node) ([]*Node, error) {
	var nodes []*Node
	for _, in := range iNodes {
		pn, err := toplevel(in)
		if err != nil {
			return nil, err
		}
		// 中間言語にあるけどpythonにないもの
		if pn == nil {
			continue
		}
		nodes = append(nodes, pn)
	}
	return nodes, nil
}

func ConvertTypeFromInterLang(tt interlang.TType) (TType, error) {
	// TODO
	_ = tt
	return nil, nil
}

func toplevel(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.FunctionDefine:
		return functionDefine(iNode)
	default:
		return nil, nil
	}
}

func functionDefine(iNode *interlang.Node) (*Node, error) {
	iField := iNode.GetField().(*interlang.FunctionDefineField)
	// typeの変換
	iReturnValueType := iField.GetTType()
	returnValueType, err := ConvertTypeFromInterLang(iReturnValueType)

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

	return NewNode(
		FunctionDefine,
		&FunctionDefineField{
			returnValueType,
			NewNode(Ident, &IdentField{S: ident}),
			params,
			stmts,
		},
	), nil
}

func functionDefineParams(iNode *interlang.Node) (*Node, error) {
	// TODO
	return nil, nil
}

func statement(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Block:
		var stmts []*Node
		iStmts := iNode.GetField().(*interlang.BlockField).Stmts
		for _, iStmt := range iStmts {
			stmt, err := statement(iStmt)
			if err != nil {
				return nil, err
			}
			if stmt == nil {
				continue
			}
			stmts = append(stmts, stmt)
		}
		return NewNode(Block, &BlockField{stmts}), nil

	case interlang.Return:
		iReturnField := iNode.GetField().(*interlang.ReturnField)
		rv, err := expr(iReturnField.Value)
		if err != nil {
			return nil, err
		}
		return NewNode(Return, &ReturnField{Value: rv}), nil

	case interlang.IfElse:
		iIfElseField := iNode.GetField().(*interlang.IfElseField)
		cond, err := expr(iIfElseField.Cond)
		if err != nil {
			return nil, err
		}
		ifBlock, err := statement(iIfElseField.IfBlock)
		if err != nil {
			return nil, err
		}
		// ifだけ
		if iIfElseField.ElseBlock == nil {
			return NewNode(IfElse, &IfElseField{cond, ifBlock, nil}), nil
		}
		// elseあり
		elseBlock, err := statement(iIfElseField.ElseBlock)
		if err != nil {
			return nil, err
		}
		return NewNode(IfElse, &IfElseField{cond, ifBlock, elseBlock}), nil

	case interlang.While:
		// TODO
		panic("unhandled default case")
	case interlang.For:
		iForField := iNode.GetField().(*interlang.ForField)
		init, err := expr(iForField.Init)
		if err != nil {
			return nil, err
		}
		cond, err := expr(iForField.Cond)
		if err != nil {
			return nil, err
		}
		loop, err := expr(iForField.Loop)
		if err != nil {
			return nil, err
		}
		block, err := statement(iForField.Block)
		if err != nil {
			return nil, err
		}
		return NewNode(For, &ForField{init, cond, loop, block}), nil
	default:
		return expr(iNode)
	}
}

func expr(iNode *interlang.Node) (*Node, error) {
	return assign(iNode)
}

func assign(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.VariableDeclare:
		// TODO
		// pythonには存在しない概念
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
		return NewNode(Assign, &AssignField{to, value}), nil
	default:
		return andor(iNode)
	}
}

func andor(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Binary:
		iBinaryField := iNode.GetField().(*interlang.BinaryField)
		iType := iBinaryField.GetTType()
		pType, err := ConvertTypeFromInterLang(iType)
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
			return NewNode(Binary, &BinaryField{pType, And, lhs, rhs}), nil
		case interlang.Or:
			return NewNode(Binary, &BinaryField{pType, Or, lhs, rhs}), nil
		default:
			return equality(iNode)
		}
	default:
		return equality(iNode)
	}
}

func equality(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Binary:
		iBinaryField := iNode.GetField().(*interlang.BinaryField)
		iType := iBinaryField.GetTType()
		pType, err := ConvertTypeFromInterLang(iType)
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
			return NewNode(Binary, &BinaryField{pType, Eq, lhs, rhs}), nil
		case interlang.Ne:
			return NewNode(Binary, &BinaryField{pType, Ne, lhs, rhs}), nil
		default:
			return relational(iNode)
		}
	default:
		return relational(iNode)
	}
}

func relational(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Binary:
		iBinaryField := iNode.GetField().(*interlang.BinaryField)
		iType := iBinaryField.GetTType()
		pType, err := ConvertTypeFromInterLang(iType)
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
			return NewNode(Binary, &BinaryField{pType, Lt, lhs, rhs}), nil
		case interlang.Le:
			return NewNode(Binary, &BinaryField{pType, Le, lhs, rhs}), nil
		case interlang.Gt:
			return NewNode(Binary, &BinaryField{pType, Gt, lhs, rhs}), nil
		case interlang.Ge:
			return NewNode(Binary, &BinaryField{pType, Ge, lhs, rhs}), nil
		default:
			return add(iNode)
		}
	default:
		return add(iNode)
	}
}

func add(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Binary:
		iBinaryField := iNode.GetField().(*interlang.BinaryField)
		iType := iBinaryField.GetTType()
		pType, err := ConvertTypeFromInterLang(iType)
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
			return NewNode(Binary, &BinaryField{pType, Add, lhs, rhs}), nil
		case interlang.Sub:
			return NewNode(Binary, &BinaryField{pType, Sub, lhs, rhs}), nil
		default:
			return mul(iNode)
		}
	default:
		return mul(iNode)
	}
}

func mul(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Binary:
		iBinaryField := iNode.GetField().(*interlang.BinaryField)
		iType := iBinaryField.GetTType()
		pType, err := ConvertTypeFromInterLang(iType)
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
			return NewNode(Binary, &BinaryField{pType, Mul, lhs, rhs}), nil
		case interlang.Div:
			return NewNode(Binary, &BinaryField{pType, Div, lhs, rhs}), nil
		case interlang.Mod:
			return NewNode(Binary, &BinaryField{pType, Mod, lhs, rhs}), nil
		default:
			return unary(iNode)
		}
	default:
		return unary(iNode)
	}
}

func unary(iNode *interlang.Node) (*Node, error) {
	return primary(iNode)
}

func primary(iNode *interlang.Node) (*Node, error) {
	switch iNode.GetKind() {
	case interlang.Ident:
		iIdentField := iNode.GetField().(*interlang.IdentField)
		return NewNode(Ident, &IdentField{S: iIdentField.S}), nil
	case interlang.Literal:
		return literal(iNode)
	case interlang.Call:
		log.Panicf("call unimplemented: %v", iNode)
	default:
		log.Panicf("unexpected primary node: %v", iNode)
	}
	return nil, nil
}

func literal(iNode *interlang.Node) (*Node, error) {
	iLitField := iNode.GetField().(*interlang.LiteralField)
	switch iLitField.GetTType() {
	case interlang.String:
		return NewNode(Literal, &LiteralField{TType: String, S: iLitField.S}), nil
	case interlang.Integer:
		return NewNode(Literal, &LiteralField{TType: Integer, I: iLitField.I}), nil
	case interlang.Bool:
		//return NewNode(Literal, &LiteralField{TType: B, S: iLitField.S}), nil
		log.Panicf("unsupported literal: %v", iNode)
	default:
		log.Panicf("unsupported literal: %v", iNode)
	}
	return nil, nil
}
