package main

import "fmt"

type TypeChecker struct{}

func (tc *TypeChecker) EmitError(token Token, message string) {
	panic(ErrorData{token.line, token.col, message})
}

func (tc *TypeChecker) Visit(expr Expression) interface{} {
	return expr.Accept(tc)
}

func (tc *TypeChecker) VisitBinary(expr *BinaryExpression) interface{} {
	leftType := tc.Visit(expr.left).(Type)
	rightType := tc.Visit(expr.right).(Type)

	check := func(condition bool) {
		if !condition {
			tc.EmitError(expr.operator, fmt.Sprintf("invalid operation: operator %v is not defined for %v and %v.", expr.operator, leftType, rightType))
		}
	}

	sameCategory := func() bool {
		return leftType.kind == rightType.kind // todo: does not account for composite types
	}

	bothNumeric := func() bool {
		return leftType.kind.IsNumeric() && sameCategory()
	}

	bothIntegral := func() bool {
		return leftType.kind.IsIntegral() && sameCategory()
	}

	switch expr.operator.tokenType {
	case TOKEN_AMP_AMP, TOKEN_PIPE_PIPE:
		check(leftType.kind == TYPE_BOOL && rightType.kind == TYPE_BOOL)
		return SimpleType(TYPE_BOOL)
	case TOKEN_EQUAL_EQUAL, TOKEN_BANG_EQUAL:
		check(sameCategory())
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
		check(sameCategory() && (leftType.kind.IsNumeric() || leftType.kind == TYPE_STRING))
		return leftType
	}

	Unreachable("TypeChecker::VisitUnary")
	return nil
}

func (tc *TypeChecker) VisitUnary(expr *UnaryExpression) interface{} {
	operandType := tc.Visit(expr.operand).(Type)

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
	case TOKEN_NIL:
		return SimpleType(TYPE_NIL)
	case TOKEN_INT:
		return SimpleType(TYPE_I64)
	case TOKEN_FLOAT:
		return SimpleType(TYPE_DOUBLE)
	case TOKEN_STRING:
		return SimpleType(TYPE_STRING)
	}

	Unreachable("TypeChecker::VisitLiteral")
	return nil
}

func (tc *TypeChecker) VisitGrouping(expr *GroupingExpression) interface{} {
	return tc.Visit(expr.expr)
}

func TypeCheck(ast Expression, errorReporter ErrorReporter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case ErrorData:
				// recover from any calls to panic with an argument of type `ErrorData` and push the error to the reporter
				errorReporter.Push(v.line, v.col, v.message)

				// override the returned error
				err = errorReporter
			default:
				// else re-panic
				panic(v)
			}
		}
	}()

	typeChecker := TypeChecker{}
	typeChecker.Visit(ast)
	return nil
}
