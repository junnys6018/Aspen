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
	case TOKEN_PERCENT:
		return OperatorModulus(lhs, rhs)
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
	return i.environment.GetAt(expr.name.String(), expr.depth)
}

func (i *Interpreter) VisitGrouping(expr *GroupingExpression) interface{} {
	return i.VisitExpressionNode(expr.expr)
}

func (i *Interpreter) VisitAssignment(expr *AssignmentExpression) interface{} {
	value := i.VisitExpressionNode(expr.value)
	i.environment.AssignAt(expr.name.String(), expr.depth, value)
	return value
}

func (i *Interpreter) VisitCall(expr *CallExpression) interface{} {
	callee := i.VisitExpressionNode(expr.callee).(AspenFunction)

	arguments := make([]interface{}, callee.Arity())
	for j := range arguments {
		arguments[j] = i.VisitExpressionNode(expr.arguments[j])
	}

	return callee.Call(i, arguments)
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

func (i *Interpreter) ExecuteBlock(stmt *BlockStatement, environment Environment) {
	enclosing := i.environment
	i.environment = environment

	// we defer because VisitStatementNode can panic when we return
	defer func() {
		i.environment = enclosing
	}()

	for _, stmt := range stmt.statements {
		i.VisitStatementNode(stmt)
	}
}

func (i *Interpreter) VisitBlock(stmt *BlockStatement) interface{} {
	enclosing := i.environment
	i.ExecuteBlock(stmt, NewEnvironment(&enclosing))
	return nil
}

func (i *Interpreter) VisitIf(stmt *IfStatement) interface{} {
	condition := i.VisitExpressionNode(stmt.condition).(bool)

	if condition {
		i.VisitStatementNode(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.VisitStatementNode(stmt.elseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhile(stmt *WhileStatement) interface{} {
	for i.VisitExpressionNode(stmt.condition).(bool) {
		i.VisitStatementNode(stmt.body)
	}
	return nil
}

func (i *Interpreter) VisitFunction(stmt *FunctionStatement) interface{} {
	i.environment.Define(stmt.name.String(), &UserFunction{declaration: stmt, closure: i.environment})
	return nil
}

func (i *Interpreter) VisitReturn(stmt *ReturnStatement) interface{} {
	var value ReturnValue

	if stmt.value != nil {
		value.value = i.VisitExpressionNode(stmt.value)
	}

	panic(value)
}

func Interpret(ast Program) (err error) {
	interpreter := Interpreter{environment: NewGlobalEnvironment()}

	for _, stmt := range ast {
		interpreter.VisitStatementNode(stmt)
	}
	return nil
}
