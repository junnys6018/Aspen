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
		return leftType.category == rightType.category
	}

	bothNumeric := func() bool {
		return leftType.category.IsNumeric() && sameCategory()
	}

	switch expr.operator.tokenType {
	case TOKEN_AMP_AMP, TOKEN_PIPE_PIPE:
		check(leftType.category == TYPE_BOOL && rightType.category == TYPE_BOOL)
		return SimpleType(TYPE_BOOL)
	case TOKEN_EQUAL_EQUAL, TOKEN_BANG_EQUAL:
		check(sameCategory())
		return leftType
	case TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL:
		check(bothNumeric())
		return SimpleType(TYPE_BOOL)
	case TOKEN_PIPE, TOKEN_CARET, TOKEN_AMP, TOKEN_MINUS, TOKEN_SLASH, TOKEN_STAR:
		check(bothNumeric())
		return leftType
	case TOKEN_PLUS:
		check(sameCategory() && (leftType.category.IsNumeric() || leftType.category == TYPE_STRING))
		return leftType
	default:
		panic("SemanticAnalyzer::VisitUnary unknown type")
	}
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
		check(operandType.category == TYPE_BOOL)
		return SimpleType(TYPE_BOOL)
	case TOKEN_MINUS:
		check(operandType.category.IsNumeric())
		return operandType
	default:
		panic("SemanticAnalyzer::VisitUnary unknown type")
	}
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
	default:
		panic("SemanticAnalyzer::VisitLiteral unknown type")
	}
}

func (tc *TypeChecker) VisitGrouping(expr *GroupingExpression) interface{} {
	return tc.Visit(expr.expr)
}

func TypeCheck(ast Expression, errorReporter ErrorReporter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case ErrorData:
				errorReporter.Push(v.line, v.col, v.message)
				err = errorReporter
			case string:
				panic(v)
			}
		}
	}()

	semanticAnalyzer := TypeChecker{}
	semanticAnalyzer.Visit(ast)
	return nil
}
