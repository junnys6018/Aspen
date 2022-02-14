package main

type Interpreter struct{}

func (i *Interpreter) Visit(expr Expression) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(expr *BinaryExpression) interface{} {
	lhs := i.Visit(expr.left)
	rhs := i.Visit(expr.right)

	switch expr.operator.tokenType {
	case TOKEN_AMP_AMP:
		return lhs.(bool) && rhs.(bool)
	case TOKEN_PIPE_PIPE:
		return lhs.(bool) || rhs.(bool)
	case TOKEN_EQUAL_EQUAL:
		return ValuesEqual(lhs, rhs)
	case TOKEN_BANG_EQUAL:
		return !ValuesEqual(lhs, rhs)
	case TOKEN_GREATER:
		return OperatorGreater(lhs, rhs)
	case TOKEN_GREATER_EQUAL:
		return OperatorGreaterEqual(lhs, rhs)
	case TOKEN_LESS:
		return OperatorLess(lhs, rhs)
	case TOKEN_LESS_EQUAL:
		return OperatorLessEqual(lhs, rhs)
	case TOKEN_PIPE:
		return OperatorPipe(lhs, rhs)
	case TOKEN_CARET:
		return OperatorCaret(lhs, rhs)
	case TOKEN_AMP:
		return OperatorAmp(lhs, rhs)
	case TOKEN_MINUS:
		return OperatorMinus(lhs, rhs)
	case TOKEN_SLASH:
		return OperatorSlash(lhs, rhs)
	case TOKEN_STAR:
		return OperatorStar(lhs, rhs)
	case TOKEN_PLUS:
		switch lhsV := lhs.(type) {
		case []rune:
			rhsV := rhs.([]rune)
			return AddString(lhsV, rhsV)
		default:
			return OperatorPlus(lhs, rhs)
		}
	}

	Unreachable("Interpreter::VisitBinary")
	return nil
}

func (i *Interpreter) VisitUnary(expr *UnaryExpression) interface{} {
	operand := i.Visit(expr.operand)
	switch expr.operator.tokenType {
	case TOKEN_BANG:
		return !operand.(bool)
	case TOKEN_MINUS:
		switch v := operand.(type) {
		case int64:
			return -v
		case uint64:
			return -v
		case float64:
			return -v
		}
	}

	Unreachable("Interpreter::VisitUnary")
	return nil
}

func (i *Interpreter) VisitLiteral(expr *LiteralExpression) interface{} {
	switch expr.value.tokenType {
	case TOKEN_FALSE:
		return false
	case TOKEN_TRUE:
		return true
	case TOKEN_NIL:
		return nil
	case TOKEN_INT:
		return expr.value.value.(int64)
	case TOKEN_FLOAT:
		return expr.value.value.(float64)
	case TOKEN_STRING:
		return expr.value.value.([]rune)
	}

	Unreachable("Interpreter::VisitLiteral")
	return nil
}

func (i *Interpreter) VisitGrouping(expr *GroupingExpression) interface{} {
	return i.Visit(expr.expr)
}

func Interpret(ast Expression) (err error) {
	interpreter := Interpreter{}
	value := interpreter.Visit(ast)
	PrintValue(value)
	return nil
}
