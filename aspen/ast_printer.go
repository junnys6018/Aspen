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

func (p *AstPrinter) parenthesize(name string, exprOrStmts ...interface{}) {
	p.builder.WriteRune('(')

	p.builder.WriteString(name)

	for _, exprOrStmt := range exprOrStmts {
		if exprOrStmt == nil {
			continue
		}
		p.builder.WriteRune(' ')
		switch exprOrStmt := exprOrStmt.(type) {
		case Expression:
			p.VisitExpressionNode(exprOrStmt)
		case Statement:
			p.VisitStatementNode(exprOrStmt)
		}
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

func (p *AstPrinter) VisitAssignment(expr *AssignmentExpression) interface{} {
	p.parenthesize(fmt.Sprintf("= (identifier %s)", expr.name), expr.value)
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

func (p *AstPrinter) VisitBlock(stmt *BlockStatement) interface{} {
	// have to convert stmt.statements into a []interface{} to pass into parenthesize...
	in := make([]interface{}, 0, len(stmt.statements))
	for _, stmt := range stmt.statements {
		in = append(in, stmt)
	}

	p.parenthesize("block", in...)
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
