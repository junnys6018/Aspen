package main

import (
	"fmt"
	"strings"
)

type ReferenceNode struct {
	// the set of undefined functions this node directly or indirectly references
	unresolvedReferences NodeSet

	// the list of nodes this node directly references
	references []*ReferenceNode

	// the set of nodes that reference this node (either directly or indirectly)
	isReferencedBy NodeSet

	// the tokens in which the references was made
	referenceLocations []*Token
}

type NodeSet map[*ReferenceNode]struct{}

type ReferenceGraph struct {
	// maps a function declaration to its corresponding node in the reference graph
	nodes map[*FunctionStatement]*ReferenceNode

	// the set of undefined functions
	undefinedFunctions NodeSet
}

func (g *ReferenceGraph) GetFunction(target *ReferenceNode) *FunctionStatement {
	for fn, node := range g.nodes {
		if node == target {
			return fn
		}
	}
	Unreachable("ReferenceGraph::GetFunction")
	return nil
}

func (g *ReferenceGraph) AddNode(fn *FunctionStatement) {
	g.nodes[fn] = &ReferenceNode{
		unresolvedReferences: make(NodeSet),
		isReferencedBy:       make(NodeSet),
	}

	// we mark a function as always referencing itself (even if it might not be the case)
	g.nodes[fn].isReferencedBy[g.nodes[fn]] = struct{}{}
}

func (g *ReferenceGraph) AddUndefinedNode(fn *FunctionStatement) {
	g.AddNode(fn)

	// mark the newly added node as undefined
	g.undefinedFunctions[g.nodes[fn]] = struct{}{}

	// since a node always references itself, we add the node to its own list
	// of unresolved references
	g.nodes[fn].unresolvedReferences[g.nodes[fn]] = struct{}{}
}

func (g *ReferenceGraph) AddEdge(from, to *FunctionStatement, loc *Token) {
	fromNode := g.nodes[from]
	toNode := g.nodes[to]

	fromNode.references = append(fromNode.references, toNode)
	fromNode.referenceLocations = append(fromNode.referenceLocations, loc)

	// any node that referenced `from` indirectly now references `to` indirectly
	for k := range fromNode.isReferencedBy {
		toNode.isReferencedBy[k] = struct{}{}
	}

	// any unresolved reference made by the `to` function will be inherited by any function that references `from`
	// this is why we mark a function as referencing itself, so that `from` also inherits the unresolved references
	for node := range toNode.unresolvedReferences {
		for from := range fromNode.isReferencedBy {
			from.unresolvedReferences[node] = struct{}{}
			node.isReferencedBy[from] = struct{}{}
		}
	}
}

func (g *ReferenceGraph) MarkNodeAsDefined(fn *FunctionStatement) {
	definedNode := g.nodes[fn]

	// for each function that references the undefined node
	for node := range definedNode.isReferencedBy {
		// ...remove that reference. The node is now defined
		delete(node.unresolvedReferences, definedNode)
	}

	// the node is now defined, remove its entry here
	delete(g.undefinedFunctions, definedNode)
}

func (g *ReferenceGraph) ReferencesUndefinedNode(fn *FunctionStatement) (bool, []*Token) {
	// does `fn` reference any undefined function?
	ok := len(g.nodes[fn].unresolvedReferences) == 0
	if ok {
		// if no return early
		return false, nil
	}

	// traverse through the graph to find the path to the undefined function
	parent := make(map[*ReferenceNode]*ReferenceNode)
	visitedNodes := make(NodeSet)
	stack := make([]*ReferenceNode, 0)
	stack = append(stack, g.nodes[fn])

	// the last element of the path
	var end *ReferenceNode

	// dfs our way to an undefined function
	for len(stack) != 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// mark node as visited
		visitedNodes[node] = struct{}{}

		if _, ok := g.undefinedFunctions[node]; ok {
			// the node we popped is undefined, we are done
			end = node
			break
		}

		for _, child := range node.references {
			_, visited := visitedNodes[child]
			if !visited {
				parent[child] = node
				stack = append(stack, child)
			}
		}
	}

	chain := make([]*Token, 0)
	chain = append(chain, &g.GetFunction(end).name)

	// traverse the parent chain starting from `end` until we reach the root node
	for {
		p, ok := parent[end]
		if !ok {
			// end has no parent, ie we have reached the root node
			break
		}

		// find the corresponding token and append it
		for i := range p.references {
			if p.references[i] == end {
				chain = append(chain, p.referenceLocations[i])
				break
			}
		}
		end = p
	}

	return true, chain
}

func NewReferenceGraph() *ReferenceGraph {
	return &ReferenceGraph{
		nodes:              make(map[*FunctionStatement]*ReferenceNode),
		undefinedFunctions: make(NodeSet),
	}
}

func UnresolvedErrorMessage(chain []*Token, start *Token) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "reference to unresolved function '%v'.", chain[0])

	if len(chain) > 1 {
		fmt.Fprintf(&builder, "\n\n    %d:%d %v refers to\n", start.line, start.col, start)
		for i := len(chain) - 1; i > 0; i-- {
			token := chain[i]
			if i == 1 {
				fmt.Fprintf(&builder, "    %d:%d %v", token.line, token.col, token)
			} else {
				fmt.Fprintf(&builder, "    %d:%d %v refers to\n", token.line, token.col, token)
			}
		}
	}

	return builder.String()
}

type Scopes []map[string]*FunctionStatement

func (s Scopes) GetAt(name string, depth int) *FunctionStatement {
	return s[len(s)-depth-1][name]
}

func (s Scopes) GetGlobal(name string) *FunctionStatement {
	return s[0][name]
}

func (s Scopes) Define(name string, fn *FunctionStatement) {
	s[len(s)-1][name] = fn
}

type TypeChecker struct {
	environment     Environment
	errorReporter   ErrorReporter
	currentFunction *FunctionStatement

	referenceGraph *ReferenceGraph
	scopes         Scopes
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
	atype := tc.environment.GetAt(name, expr.depth)

	if atype.(*Type).kind == TYPE_FUNCTION {
		if fn := tc.scopes.GetAt(name, expr.depth); fn != nil {
			if tc.currentFunction == nil {
				// make sure reference to function in top level code is not undefined
				if err, chain := tc.referenceGraph.ReferencesUndefinedNode(fn); err {
					var location Token
					if len(chain) > 1 {
						location = *chain[0]
					} else {
						location = expr.name
					}
					tc.Error(location, UnresolvedErrorMessage(chain, &expr.name))
				}
			} else {
				// tc.currentFunction references fn
				tc.referenceGraph.AddEdge(tc.currentFunction, fn, &expr.name)
			}
		}
	}

	return atype
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

func (tc *TypeChecker) VisitTypeCast(expr *TypeCastExpression) interface{} {
	from := tc.VisitExpressionNode(expr.value).(*Type)
	expr.from = from

	if !IsConversionLegal(from, expr.to) {
		tc.FatalError(expr.loc, fmt.Sprintf("cannot cast expression of type %v to %v.", from, expr.to))
	}

	return expr.to
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
			stmt.initializer = &TypeCastExpression{
				from:  SimpleType(TYPE_I64),
				to:    SimpleType(TYPE_U64),
				value: &LiteralExpression{value: Token{tokenType: TOKEN_INT_LITERAL, value: int64(0)}},
			}
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

	tc.scopes = append(tc.scopes, make(map[string]*FunctionStatement))

	for _, stmt := range stmt.statements {
		tc.VisitStatementNode(stmt)
	}

	tc.scopes = tc.scopes[:len(tc.scopes)-1]

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

func (tc *TypeChecker) DefineFunction(name string, atype FunctionType) bool {
	if tc.environment.IsDefinedLocally(name) {
		return false
	}

	tc.environment.Define(name, &Type{kind: TYPE_FUNCTION, other: atype})
	return true
}

func (tc *TypeChecker) VisitFunction(stmt *FunctionStatement) interface{} {
	name := stmt.name.String()

	if tc.environment.enclosing != nil {
		// skip defining global functions, they were defined in the first pass
		if !tc.DefineFunction(name, stmt.atype) {
			tc.Error(stmt.name, fmt.Sprintf("cannot redefine '%s'.", name))
		}
		tc.referenceGraph.AddNode(stmt)
	} else {
		// mark global function as defined
		tc.referenceGraph.MarkNodeAsDefined(stmt)
	}

	tc.scopes.Define(name, stmt)

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
	tc.currentFunction = stmt
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

		returnType := tc.currentFunction.atype.returnType

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
func NewTypeChecker(errorReporter ErrorReporter) *TypeChecker {
	typeChecker := TypeChecker{
		environment:    NewEnvironment(nil),
		errorReporter:  errorReporter,
		scopes:         make(Scopes, 1),
		referenceGraph: NewReferenceGraph(),
	}

	typeChecker.scopes[0] = make(map[string]*FunctionStatement)

	return &typeChecker
}

func TypeCheck(ast Program, errorReporter ErrorReporter) (err error) {
	typeChecker := NewTypeChecker(errorReporter)

	// define native functions
	for name, fn := range NativeFunctions {
		typeChecker.DefineFunction(name, fn.atype)
	}

	// define global functions
	for _, stmt := range ast {
		fn, ok := stmt.(*FunctionStatement)
		if ok {
			name := fn.name.String()
			if !typeChecker.DefineFunction(name, fn.atype) {
				typeChecker.Error(fn.name, fmt.Sprintf("cannot redefine '%s'.", name))
			}
			typeChecker.referenceGraph.AddUndefinedNode(fn)
			typeChecker.scopes.Define(name, fn)
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
