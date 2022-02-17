package main

type ExpressionVisitor interface {
	VisitBinary(expr *BinaryExpression) interface{}
	VisitUnary(expr *UnaryExpression) interface{}
	VisitLiteral(expr *LiteralExpression) interface{}
	VisitGrouping(expr *GroupingExpression) interface{}
	VisitIdentifier(expr *IdentifierExpression) interface{}
}
type Expression interface {
	Accept(visitor ExpressionVisitor) interface{}
	String() string
}
type BinaryExpression struct {
	left     Expression
	right    Expression
	operator Token
}

func (expr *BinaryExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitBinary(expr)
}
func (expr *BinaryExpression) String() string {
	printer := AstPrinter{}
	printer.VisitExpressionNode(expr)
	return printer.builder.String()
}

type UnaryExpression struct {
	operand  Expression
	operator Token
}

func (expr *UnaryExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitUnary(expr)
}
func (expr *UnaryExpression) String() string {
	printer := AstPrinter{}
	printer.VisitExpressionNode(expr)
	return printer.builder.String()
}

type LiteralExpression struct {
	value Token
}

func (expr *LiteralExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteral(expr)
}
func (expr *LiteralExpression) String() string {
	printer := AstPrinter{}
	printer.VisitExpressionNode(expr)
	return printer.builder.String()
}

type GroupingExpression struct {
	expr Expression
}

func (expr *GroupingExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitGrouping(expr)
}
func (expr *GroupingExpression) String() string {
	printer := AstPrinter{}
	printer.VisitExpressionNode(expr)
	return printer.builder.String()
}

type IdentifierExpression struct {
	name Token
}

func (expr *IdentifierExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitIdentifier(expr)
}
func (expr *IdentifierExpression) String() string {
	printer := AstPrinter{}
	printer.VisitExpressionNode(expr)
	return printer.builder.String()
}

type StatementVisitor interface {
	VisitExpression(stmt *ExpressionStatement) interface{}
	VisitPrint(stmt *PrintStatement) interface{}
	VisitLet(stmt *LetStatement) interface{}
	VisitBlock(stmt *BlockStatement) interface{}
}
type Statement interface {
	Accept(visitor StatementVisitor) interface{}
	String() string
}
type ExpressionStatement struct {
	expr Expression
}

func (stmt *ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpression(stmt)
}
func (stmt *ExpressionStatement) String() string {
	printer := AstPrinter{}
	printer.VisitStatementNode(stmt)
	return printer.builder.String()
}

type PrintStatement struct {
	expr Expression
}

func (stmt *PrintStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitPrint(stmt)
}
func (stmt *PrintStatement) String() string {
	printer := AstPrinter{}
	printer.VisitStatementNode(stmt)
	return printer.builder.String()
}

type LetStatement struct {
	name        Token
	initializer Expression
	atype       *Type
}

func (stmt *LetStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitLet(stmt)
}
func (stmt *LetStatement) String() string {
	printer := AstPrinter{}
	printer.VisitStatementNode(stmt)
	return printer.builder.String()
}

type BlockStatement struct {
	statements []Statement
}

func (stmt *BlockStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitBlock(stmt)
}
func (stmt *BlockStatement) String() string {
	printer := AstPrinter{}
	printer.VisitStatementNode(stmt)
	return printer.builder.String()
}
