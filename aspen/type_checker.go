package main

import "fmt"

type TypeChecker struct {
	environment     Environment
	errorReporter   ErrorReporter
	currentFunction *FunctionType
}

func (tc *TypeChecker) FatalError(token Token, message string) {
	panic(ErrorData{token.line, token.col, message})
}

func (tc *TypeChecker) Error(token Token, message string) {
	tc.errorReporter.Push(token.line, token.col, message)
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
			tc.FatalError(expr.operator, fmt.Sprintf("invalid operation: operator %v is not defined for %v and %v.", expr.operator, leftType, rightType))
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
	case TOKEN_PIPE, TOKEN_CARET, TOKEN_AMP, TOKEN_PERCENT:
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
			tc.FatalError(expr.operator, fmt.Sprintf("invalid operation: operator %v is not defined for %v.", expr.operator, operandType))
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
		tc.FatalError(expr.name, fmt.Sprintf("undeclared identifier '%s'.", name))
	}

	expr.depth = tc.environment.GetDepth(name)

	return tc.environment.GetAt(name, expr.depth)
}

func (tc *TypeChecker) VisitGrouping(expr *GroupingExpression) interface{} {
	return tc.VisitExpressionNode(expr.expr)
}

func (tc *TypeChecker) VisitAssignment(expr *AssignmentExpression) interface{} {
	name := expr.name.String()

	if !tc.environment.IsDefined(name) {
		tc.FatalError(expr.name, fmt.Sprintf("undeclared identifier '%s'.", name))
	}

	expr.depth = tc.environment.GetDepth(name)

	identifierType := tc.environment.GetAt(name, expr.depth).(*Type)
	valueType := tc.VisitExpressionNode(expr.value).(*Type)

	if !TypesEqual(identifierType, valueType) {
		tc.FatalError(expr.name, fmt.Sprintf("cannot assign expression of type %v to '%s', which has type %v.", valueType, name, identifierType))
	}

	return identifierType
}

func (tc *TypeChecker) VisitCall(expr *CallExpression) interface{} {
	callee := tc.VisitExpressionNode(expr.callee).(*Type)

	if callee.kind != TYPE_FUNCTION {
		tc.FatalError(expr.loc, "callee is not a function.")
	}

	other := callee.other.(FunctionType)

	// check arity
	if len(expr.arguments) != other.Arity() {
		if len(expr.arguments) < other.Arity() {
			tc.FatalError(expr.loc, "not enough arguments in call to function.")
		} else {
			tc.FatalError(expr.loc, "too many arguments in call to function.")
		}
	}

	for i := range expr.arguments {
		arg := tc.VisitExpressionNode(expr.arguments[i]).(*Type)
		if !TypesEqual(arg, other.parameters[i]) {
			tc.Error(expr.loc,
				fmt.Sprintf("cannot use argument of type %v as the %s parameter to function call (expected %v).",
					arg,
					OrdinalSuffixOf(i+1),
					other.parameters[i]))
		}
	}

	return other.returnType
}

func (tc *TypeChecker) VisitExpression(stmt *ExpressionStatement) interface{} {
	tc.VisitExpressionNode(stmt.expr)
	return nil
}

func (tc *TypeChecker) VisitPrint(stmt *PrintStatement) interface{} {
	value := tc.VisitExpressionNode(stmt.expr).(*Type)
	if value.kind == TYPE_VOID {
		tc.Error(stmt.loc, "cannot print an expression of type void.")
	}
	return nil
}

func (tc *TypeChecker) VisitLet(stmt *LetStatement) interface{} {
	name := stmt.name.String()
	if tc.environment.IsDefinedLocally(name) {
		tc.FatalError(stmt.name, fmt.Sprintf("cannot redefine '%s'.", name))
	}

	// Slice and function types must be initialized
	if stmt.initializer == nil && (stmt.atype.kind == TYPE_SLICE || stmt.atype.kind == TYPE_FUNCTION) {
		tc.FatalError(stmt.name, fmt.Sprintf("'%s' must be initialized.", stmt.name.value))
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
			tc.FatalError(stmt.name, fmt.Sprintf("cannot assign expression of type %v to '%s', which has type %v.", atype, stmt.name.value, stmt.atype))
		}
	}

	tc.environment.Define(name, stmt.atype)

	return nil
}

func (tc *TypeChecker) CheckBlock(stmt *BlockStatement, environment Environment) {
	enclosing := tc.environment
	tc.environment = environment

	for _, stmt := range stmt.statements {
		tc.VisitStatementNode(stmt)
	}

	tc.environment = enclosing
}

func (tc *TypeChecker) VisitBlock(stmt *BlockStatement) interface{} {
	enclosing := tc.environment
	tc.CheckBlock(stmt, NewEnvironment(&enclosing))
	return nil
}

func (tc *TypeChecker) VisitIf(stmt *IfStatement) interface{} {
	// check that the condition is a bool
	condition := tc.VisitExpressionNode(stmt.condition).(*Type)
	if condition.kind != TYPE_BOOL {
		tc.Error(stmt.loc, "expected an expression of type bool.")
	}

	// visit the then and else block
	tc.VisitStatementNode(stmt.thenBranch)
	if stmt.elseBranch != nil {
		tc.VisitStatementNode(stmt.elseBranch)
	}

	return nil
}

func (tc *TypeChecker) VisitWhile(stmt *WhileStatement) interface{} {
	// check that the condition is a bool
	condition := tc.VisitExpressionNode(stmt.condition).(*Type)
	if condition.kind != TYPE_BOOL {
		tc.Error(stmt.loc, "expected an expression of type bool.")
	}

	tc.VisitStatementNode(stmt.body)
	return nil
}

func (tc *TypeChecker) DefineFunction(stmt *FunctionStatement) {
	name := stmt.name.String()
	if tc.environment.IsDefinedLocally(name) {
		tc.Error(stmt.name, fmt.Sprintf("cannot redefine '%s'.", name))
	}

	tc.environment.Define(name, &Type{kind: TYPE_FUNCTION, other: stmt.atype})
}

func (tc *TypeChecker) VisitFunction(stmt *FunctionStatement) interface{} {
	if tc.environment.enclosing != nil {
		// skip defining global functions, they were defined in the first pass
		tc.DefineFunction(stmt)
	}

	enclosing := tc.environment
	environment := NewEnvironment(&enclosing)

	for i := range stmt.parameters {
		environment.Define(stmt.parameters[i].String(), stmt.atype.parameters[i])
	}

	if !stmt.atype.returnType.IsVoid() {
		// check that the last statement in the body of the function is a return statement
		if len(stmt.body.statements) == 0 {
			tc.Error(stmt.name, "missing return.")
		} else {
			_, ok := stmt.body.statements[len(stmt.body.statements)-1].(*ReturnStatement)
			if !ok {
				tc.Error(stmt.name, "missing return.")
			}
		}
	}

	enclosingFn := tc.currentFunction
	tc.currentFunction = &stmt.atype
	tc.CheckBlock(stmt.body, environment)
	tc.currentFunction = enclosingFn
	return nil
}

func (tc *TypeChecker) VisitReturn(stmt *ReturnStatement) interface{} {
	if tc.currentFunction == nil {
		tc.FatalError(stmt.loc, "cannot return from top level code.")
	} else {
		var value *Type
		if stmt.value == nil {
			value = SimpleType(TYPE_VOID)
		} else {
			value = tc.VisitExpressionNode(stmt.value).(*Type)
		}

		returnType := tc.currentFunction.returnType

		if returnType.IsVoid() && !value.IsVoid() {
			tc.Error(stmt.loc, "no return values expected.")
		} else if !returnType.IsVoid() && value.IsVoid() {
			tc.Error(stmt.loc, fmt.Sprintf("function must return an expression of type %v.", returnType))
		} else if !TypesEqual(value, returnType) {
			tc.Error(stmt.loc, fmt.Sprintf("cannot return an expression of type %v (%v expected).", value, returnType))
		}
	}
	return nil
}

func TypeCheck(ast Program, errorReporter ErrorReporter) (err error) {
	// initialize global environment
	environment := NewEnvironment(nil)

	// copy native functions into global environment
	for k, v := range NativeFunctions {
		environment.values[k] = &Type{kind: TYPE_FUNCTION, other: v.atype}
	}

	typeChecker := TypeChecker{environment: environment, errorReporter: errorReporter}

	// first pass: define global functions
	for _, stmt := range ast {
		fnDecl, ok := stmt.(*FunctionStatement)
		if ok {
			typeChecker.DefineFunction(fnDecl)
		}
	}

	for _, stmt := range ast {
		typeChecker.VisitStatementNode(stmt)
	}

	if typeChecker.errorReporter.HadError() {
		return errorReporter
	}

	return nil
}
