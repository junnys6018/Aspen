package main

import "fmt"

type TypeChecker struct {
	environment   Environment
	errorReporter ErrorReporter
}

func (tc *TypeChecker) EmitError(token Token, message string) {
	panic(ErrorData{token.line, token.col, message})
}

func (tc *TypeChecker) VisitExpressionNode(expr Expression) interface{} {
	return expr.Accept(tc)
}

func (tc *TypeChecker) VisitStatementNode(stmt Statement) interface{} {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case ErrorData:
				// recover from any calls to panic with an argument of type `ErrorData` and push the error to the reporter
				tc.errorReporter.Push(v.line, v.col, v.message)
			default:
				// else re-panic
				panic(v)
			}
		}
	}()

	stmt.Accept(tc)
	return nil
}

func (tc *TypeChecker) VisitBinary(expr *BinaryExpression) interface{} {
	leftType := tc.VisitExpressionNode(expr.left).(*Type)
	rightType := tc.VisitExpressionNode(expr.right).(*Type)

	check := func(condition bool) {
		if !condition {
			tc.EmitError(expr.operator, fmt.Sprintf("invalid operation: operator %v is not defined for %v and %v.", expr.operator, leftType, rightType))
		}
	}

	bothNumeric := func() bool {
		return leftType.kind.IsNumeric() && leftType.kind == rightType.kind
	}

	bothIntegral := func() bool {
		return leftType.kind.IsIntegral() && leftType.kind == rightType.kind
	}

	switch expr.operator.tokenType {
	case TOKEN_AMP_AMP, TOKEN_PIPE_PIPE:
		check(leftType.kind == TYPE_BOOL && rightType.kind == TYPE_BOOL)
		return SimpleType(TYPE_BOOL)
	case TOKEN_EQUAL_EQUAL, TOKEN_BANG_EQUAL:
		check(TypesEqual(leftType, rightType))
		return SimpleType(TYPE_BOOL)
	case TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL:
		check(bothNumeric())
		return SimpleType(TYPE_BOOL)
	case TOKEN_PIPE, TOKEN_CARET, TOKEN_AMP:
		check(bothIntegral())
		return leftType
	case TOKEN_MINUS, TOKEN_SLASH, TOKEN_STAR:
		check(bothNumeric())
		return leftType
	case TOKEN_PLUS:
		check(leftType.kind == rightType.kind && (leftType.kind.IsNumeric() || leftType.kind == TYPE_STRING))
		return leftType
	}

	Unreachable("TypeChecker::VisitUnary")
	return nil
}

func (tc *TypeChecker) VisitUnary(expr *UnaryExpression) interface{} {
	operandType := tc.VisitExpressionNode(expr.operand).(*Type)

	check := func(condition bool) {
		if !condition {
			tc.EmitError(expr.operator, fmt.Sprintf("invalid operation: operator %v is not defined for %v.", expr.operator, operandType))
		}
	}

	switch expr.operator.tokenType {
	case TOKEN_BANG:
		check(operandType.kind == TYPE_BOOL)
		return SimpleType(TYPE_BOOL)
	case TOKEN_MINUS:
		check(operandType.kind.IsNumeric())
		return operandType
	}

	Unreachable("TypeChecker::VisitUnary")
	return nil
}

func (tc *TypeChecker) VisitLiteral(expr *LiteralExpression) interface{} {
	switch expr.value.tokenType {
	case TOKEN_FALSE, TOKEN_TRUE:
		return SimpleType(TYPE_BOOL)
	case TOKEN_INT_LITERAL:
		return SimpleType(TYPE_I64)
	case TOKEN_FLOAT_LITERAL:
		return SimpleType(TYPE_DOUBLE)
	case TOKEN_STRING_LITERAL:
		return SimpleType(TYPE_STRING)
	}

	Unreachable("TypeChecker::VisitLiteral")
	return nil
}

func (tc *TypeChecker) VisitIdentifier(expr *IdentifierExpression) interface{} {
	name := expr.name.String()

	if !tc.environment.IsDefined(name) {
		tc.EmitError(expr.name, fmt.Sprintf("undeclared identifier '%s'.", name))
	}

	return tc.environment.Get(name)
}

func (tc *TypeChecker) VisitGrouping(expr *GroupingExpression) interface{} {
	return tc.VisitExpressionNode(expr.expr)
}

func (tc *TypeChecker) VisitAssignment(expr *AssignmentExpression) interface{} {
	name := expr.name.String()

	if !tc.environment.IsDefined(name) {
		tc.EmitError(expr.name, fmt.Sprintf("undeclared identifier '%s'.", name))
	}

	identifierType := tc.environment.Get(name).(*Type)
	valueType := tc.VisitExpressionNode(expr.value).(*Type)

	if !TypesEqual(identifierType, valueType) {
		tc.EmitError(expr.name, fmt.Sprintf("cannot assign expression of type %v to '%s', which has type %v.", valueType, name, identifierType))
	}

	return identifierType
}

func (tc *TypeChecker) VisitExpression(stmt *ExpressionStatement) interface{} {
	tc.VisitExpressionNode(stmt.expr)
	return nil
}

func (tc *TypeChecker) VisitPrint(stmt *PrintStatement) interface{} {
	tc.VisitExpressionNode(stmt.expr)
	return nil
}

func (tc *TypeChecker) VisitLet(stmt *LetStatement) interface{} {
	// Slice and function types must be initialized
	if stmt.initializer == nil && (stmt.atype.kind == TYPE_SLICE || stmt.atype.kind == TYPE_FUNCTION) {
		tc.EmitError(stmt.name, fmt.Sprintf("'%s' must be initialized.", stmt.name.value))
	}

	if stmt.initializer == nil {
		// Insert a default value for the initializer
		switch stmt.atype.kind {
		case TYPE_I64:
			stmt.initializer = &LiteralExpression{value: Token{tokenType: TOKEN_INT_LITERAL, value: int64(0)}}
		case TYPE_U64:
			// todo: set initializer to the expression `u64(0)`
		case TYPE_BOOL:
			stmt.initializer = &LiteralExpression{value: Token{tokenType: TOKEN_FALSE}}
		case TYPE_STRING:
			stmt.initializer = &LiteralExpression{value: Token{tokenType: TOKEN_STRING_LITERAL, value: []rune("")}}
		case TYPE_DOUBLE:
			stmt.initializer = &LiteralExpression{value: Token{tokenType: TOKEN_FLOAT_LITERAL, value: float64(0)}}
		default:
			Unreachable("TypeChecker::VisitLet")
		}
	} else {
		// Type check the initializer
		atype := tc.VisitExpressionNode(stmt.initializer).(*Type)
		if !TypesEqual(stmt.atype, atype) {
			tc.EmitError(stmt.name, fmt.Sprintf("cannot assign expression of type %v to '%s', which has type %v.", atype, stmt.name.value, stmt.atype))
		}
	}

	tc.environment.Define(stmt.name.String(), stmt.atype)

	return nil
}

func (tc *TypeChecker) VisitBlock(stmt *BlockStatement) interface{} {
	enclosing := tc.environment
	tc.environment = NewEnvironment(&enclosing)

	for _, stmt := range stmt.statements {
		tc.VisitStatementNode(stmt)
	}

	tc.environment = enclosing
	return nil
}

func (tc *TypeChecker) VisitIf(stmt *IfStatement) interface{} {
	// visit the then and else block first
	tc.VisitStatementNode(stmt.thenBranch)
	if stmt.elseBranch != nil {
		tc.VisitStatementNode(stmt.elseBranch)
	}

	// then type check the condition
	condition := tc.VisitExpressionNode(stmt.condition).(*Type)
	if condition.kind != TYPE_BOOL {
		tc.EmitError(stmt.loc, "expected an expression of type bool.")
	}

	return nil
}

func TypeCheck(ast Program, errorReporter ErrorReporter) (err error) {
	typeChecker := TypeChecker{environment: NewEnvironment(nil), errorReporter: errorReporter}

	for _, stmt := range ast {
		typeChecker.VisitStatementNode(stmt)
	}

	if typeChecker.errorReporter.HadError() {
		return errorReporter
	}

	return nil
}
