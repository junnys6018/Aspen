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

func (p *AstPrinter) VisitCall(expr *CallExpression) interface{} {
	in := make([]interface{}, len(expr.arguments)+1)
	in[0] = expr.callee
	for i, arg := range expr.arguments {
		in[i+1] = arg
	}

	p.parenthesize("call", in...)
	return nil
}

func (p *AstPrinter) VisitTypeCast(expr *TypeCastExpression) interface{} {
	p.parenthesize(fmt.Sprintf("cast %v", expr.to), expr.value)
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
	p.parenthesize("block", ConvertStatementList(stmt.statements)...)
	return nil
}

func (p *AstPrinter) VisitIf(stmt *IfStatement) interface{} {
	p.parenthesize("if", stmt.condition, stmt.thenBranch, stmt.elseBranch)
	return nil
}

func (p *AstPrinter) VisitWhile(stmt *WhileStatement) interface{} {
	p.parenthesize("while", stmt.condition, stmt.body)
	return nil
}

func (p *AstPrinter) VisitFunction(stmt *FunctionStatement) interface{} {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "fn %s ", stmt.name)

	fmt.Fprintf(&builder, "(return %v)", stmt.atype.returnType)

	for i := range stmt.parameters {
		builder.WriteRune(' ')
		fmt.Fprintf(&builder, "(param %v %v)", stmt.parameters[i], stmt.atype.parameters[i])
	}

	p.parenthesize(builder.String(), ConvertStatementList(stmt.body.statements)...)
	return nil
}

func (p *AstPrinter) VisitReturn(stmt *ReturnStatement) interface{} {
	p.parenthesize("return", stmt.value)
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

func ConvertStatementList(stmts []Statement) []interface{} {
	ret := make([]interface{}, len(stmts))
	for i, stmt := range stmts {
		ret[i] = stmt
	}

	return ret
}
