package main

type ExpressionVisitor interface {
	VisitBinary(expr *BinaryExpression) interface{}
	VisitUnary(expr *UnaryExpression) interface{}
	VisitLiteral(expr *LiteralExpression) interface{}
	VisitGrouping(expr *GroupingExpression) interface{}
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
	printer.visit(expr)
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
	printer.visit(expr)
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
	printer.visit(expr)
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
	printer.visit(expr)
	return printer.builder.String()
}
