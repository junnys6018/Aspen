package main

import "strings"

type AstPrinter struct {
	builder strings.Builder
}

func (p *AstPrinter) Visit(expr Expression) {
	expr.Accept(p)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expression) {
	p.builder.WriteRune('(')

	p.builder.WriteString(name)

	for _, expr := range exprs {
		p.builder.WriteRune(' ')
		p.Visit(expr)
	}

	p.builder.WriteRune(')')
}

func (p *AstPrinter) VisitBinary(expr *BinaryExpression) interface{} {
	p.parenthesize(expr.operator.String(), expr.left, expr.right)
	return nil
}

func (p *AstPrinter) VisitUnary(expr *UnaryExpression) interface{} {
	p.parenthesize(expr.operator.String(), expr.operand)
	return nil
}

func (p *AstPrinter) VisitLiteral(expr *LiteralExpression) interface{} {
	p.builder.WriteString(expr.value.String())
	return nil
}

func (p *AstPrinter) VisitGrouping(expr *GroupingExpression) interface{} {
	p.parenthesize("group", expr.expr)
	return nil
}
