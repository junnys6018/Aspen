package main

type Interpreter struct {
	environment Environment
}

func (i *Interpreter) VisitExpressionNode(expr Expression) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitStatementNode(stmt Statement) interface{} {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBinary(expr *BinaryExpression) interface{} {
	lhs := i.VisitExpressionNode(expr.left)
	rhs := i.VisitExpressionNode(expr.right)

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
	operand := i.VisitExpressionNode(expr.operand)
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
	case TOKEN_INT_LITERAL:
		return expr.value.value.(int64)
	case TOKEN_FLOAT_LITERAL:
		return expr.value.value.(float64)
	case TOKEN_STRING_LITERAL:
		return expr.value.value.([]rune)
	}

	Unreachable("Interpreter::VisitLiteral")
	return nil
}

func (i *Interpreter) VisitIdentifier(expr *IdentifierExpression) interface{} {
	return i.environment.Get(expr.name.String())
}

func (i *Interpreter) VisitGrouping(expr *GroupingExpression) interface{} {
	return i.VisitExpressionNode(expr.expr)
}

func (i *Interpreter) VisitAssignment(expr *AssignmentExpression) interface{} {
	value := i.VisitExpressionNode(expr.value)
	i.environment.Assign(expr.name.String(), value)
	return value
}

func (i *Interpreter) VisitExpression(stmt *ExpressionStatement) interface{} {
	i.VisitExpressionNode(stmt.expr)
	return nil
}

func (i *Interpreter) VisitPrint(stmt *PrintStatement) interface{} {
	value := i.VisitExpressionNode(stmt.expr)
	PrintValue(value)
	return nil
}

func (i *Interpreter) VisitLet(stmt *LetStatement) interface{} {
	value := i.VisitExpressionNode(stmt.initializer)

	i.environment.Define(stmt.name.String(), value)
	return nil
}

func (i *Interpreter) VisitBlock(stmt *BlockStatement) interface{} {
	enclosing := i.environment
	i.environment = NewEnvironment(&enclosing)

	for _, stmt := range stmt.statements {
		i.VisitStatementNode(stmt)
	}

	i.environment = enclosing
	return nil
}

func Interpret(ast Program) (err error) {
	interpreter := Interpreter{environment: NewEnvironment(nil)}

	for _, stmt := range ast {
		interpreter.VisitStatementNode(stmt)
	}
	return nil
}