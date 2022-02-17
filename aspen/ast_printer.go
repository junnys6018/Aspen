package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	builder strings.Builder
}

func (p *AstPrinter) VisitExpressionNode(expr Expression) {
	expr.Accept(p)
}

func (p *AstPrinter) VisitStatementNode(stmt Statement) {
	stmt.Accept(p)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expression) {
	p.builder.WriteRune('(')

	p.builder.WriteString(name)

	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		p.builder.WriteRune(' ')
		p.VisitExpressionNode(expr)
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

func (p *AstPrinter) VisitIdentifier(expr *IdentifierExpression) interface{} {
	p.parenthesize(fmt.Sprintf("identifier %s", expr.name))
	return nil
}

func (p *AstPrinter) VisitGrouping(expr *GroupingExpression) interface{} {
	p.parenthesize("group", expr.expr)
	return nil
}

func (p *AstPrinter) VisitExpression(stmt *ExpressionStatement) interface{} {
	p.parenthesize("expr", stmt.expr)
	return nil
}

func (p *AstPrinter) VisitPrint(stmt *PrintStatement) interface{} {
	p.parenthesize("print", stmt.expr)
	return nil
}

func (p *AstPrinter) VisitLet(stmt *LetStatement) interface{} {
	p.parenthesize(fmt.Sprintf("let %s %s", stmt.name.value, stmt.atype), stmt.initializer)
	return nil
}

func (program Program) String() string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	for i, stmt := range program {
		fmt.Fprint(&builder, stmt)
		if i != len(program)-1 {
			builder.WriteRune(' ')
		}
	}
	builder.WriteRune(')')
	return builder.String()
}
