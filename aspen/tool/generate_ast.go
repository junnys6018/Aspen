package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Field struct {
	identifier string
	typeName   string
}

type Fields []Field

func (f Fields) asArguments() string {
	b := &strings.Builder{}
	for i, argument := range f {
		fmt.Fprintf(b, "%s %s", argument.identifier, argument.typeName)
		if i != len(f)-1 {
			b.WriteRune(',')
		}
	}

	return b.String()
}

type ASTNode struct {
	name   string
	fields Fields
}

func (n *ASTNode) write(w io.Writer, kind string) {
	fmt.Fprintf(w, "type %s%s struct {\n", n.name, kind)
	for _, field := range n.fields {
		fmt.Fprintf(w, "%s %s\n", field.identifier, field.typeName)
	}
	io.WriteString(w, "}\n")
}

type ASTMethodGenerator func(w io.Writer, nodeName string)

type ASTMethod struct {
	name       string
	arguments  Fields
	returnType string
	generator  ASTMethodGenerator
}

func (m *ASTMethod) write(w io.Writer, nodeName string, kind string, shortHandKind string) {
	fmt.Fprintf(w, "func (%s *%s%s) %s(%s) %s {\n", shortHandKind, nodeName, kind, m.name, m.arguments.asArguments(), m.returnType)

	m.generator(w, nodeName)

	io.WriteString(w, "}\n")
}

type ASTNodes struct {
	kind          string
	shortHandKind string
	nodes         []ASTNode
	methods       []ASTMethod
}

func (a *ASTNodes) defineNode(name string, fields Fields) {
	a.nodes = append(a.nodes, ASTNode{name, fields})
}

func (a *ASTNodes) defineMethod(name string, arguments Fields, returnType string, generator ASTMethodGenerator) {
	a.methods = append(a.methods, ASTMethod{name, arguments, returnType, generator})
}

func (a *ASTNodes) writeVisitorInterface(w io.Writer) {
	fmt.Fprintf(w, "type %sVisitor interface {\n", a.kind)

	for _, node := range a.nodes {
		fmt.Fprintf(w, "Visit%[1]s(%[2]s *%[1]s%[3]s) interface{}\n", node.name, a.shortHandKind, a.kind)
	}

	io.WriteString(w, "}\n")
}

func (a *ASTNodes) writeNodeInterface(w io.Writer) {
	fmt.Fprintf(w, "type %s interface {\n", a.kind)

	for _, method := range a.methods {
		fmt.Fprintf(w, "%s(%s) %s\n", method.name, method.arguments.asArguments(), method.returnType)
	}

	io.WriteString(w, "}\n")
}

func (a *ASTNodes) writeNodes(w io.Writer) {
	for _, node := range a.nodes {
		node.write(w, a.kind)
		for _, method := range a.methods {
			method.write(w, node.name, a.kind, a.shortHandKind)
		}
	}
}

func (a *ASTNodes) write(w io.Writer) {
	a.writeVisitorInterface(w)
	a.writeNodeInterface(w)
	a.writeNodes(w)
}

func GenerateASTCode() {
	exprNodes := &ASTNodes{kind: "Expression", shortHandKind: "expr"}

	exprNodes.defineNode("Binary", Fields{
		{"left", "Expression"},
		{"right", "Expression"},
		{"operator", "Token"},
	})

	exprNodes.defineNode("Unary", Fields{
		{"operand", "Expression"},
		{"operator", "Token"},
	})

	exprNodes.defineNode("Literal", Fields{
		{"value", "Token"},
	})

	exprNodes.defineNode("Grouping", Fields{
		{"expr", "Expression"},
	})

	exprNodes.defineNode("Identifier", Fields{
		{"name", "Token"},
		{"depth", "int"},
	})

	exprNodes.defineNode("Assignment", Fields{
		{"name", "Token"},
		{"value", "Expression"},
		{"depth", "int"},
	})

	exprNodes.defineNode("Call", Fields{
		{"callee", "Expression"},
		{"arguments", "[]Expression"},
		{"loc", "Token"},
	})

	exprNodes.defineNode("TypeCast", Fields{
		{"from", "*Type"},
		{"to", "*Type"},
		{"value", "Expression"},
		{"loc", "Token"},
	})

	exprNodes.defineMethod("Accept", Fields{
		{"visitor", "ExpressionVisitor"},
	}, "interface{}", func(w io.Writer, nodeName string) {
		fmt.Fprintf(w, "return visitor.Visit%s(expr)\n", nodeName)
	})

	exprNodes.defineMethod("String", Fields{}, "string", func(w io.Writer, nodeName string) {
		io.WriteString(w, "printer := AstPrinter{}\n")
		io.WriteString(w, "printer.VisitExpressionNode(expr)\n")
		io.WriteString(w, "return printer.builder.String()\n")
	})

	stmtNodes := &ASTNodes{kind: "Statement", shortHandKind: "stmt"}

	stmtNodes.defineNode("Expression", Fields{
		{"expr", "Expression"},
	})

	stmtNodes.defineNode("Print", Fields{
		{"expr", "Expression"},
		{"loc", "Token"},
	})

	stmtNodes.defineNode("Let", Fields{
		{"name", "Token"},
		{"initializer", "Expression"},
		{"atype", "*Type"},
	})

	stmtNodes.defineNode("Block", Fields{
		{"statements", "[]Statement"},
	})

	stmtNodes.defineNode("If", Fields{
		{"condition", "Expression"},
		{"thenBranch", "Statement"},
		{"elseBranch", "Statement"},
		{"loc", "Token"},
	})

	stmtNodes.defineNode("While", Fields{
		{"condition", "Expression"},
		{"body", "Statement"},
		{"loc", "Token"},
	})

	stmtNodes.defineNode("Function", Fields{
		{"name", "Token"},
		{"parameters", "[]Token"},
		{"body", "*BlockStatement"},
		{"atype", "FunctionType"},
	})

	stmtNodes.defineNode("Return", Fields{
		{"value", "Expression"},
		{"loc", "Token"},
	})

	stmtNodes.defineMethod("Accept", Fields{
		{"visitor", "StatementVisitor"},
	}, "interface{}", func(w io.Writer, nodeName string) {
		fmt.Fprintf(w, "return visitor.Visit%s(stmt)\n", nodeName)
	})

	stmtNodes.defineMethod("String", Fields{}, "string", func(w io.Writer, nodeName string) {
		io.WriteString(w, "printer := AstPrinter{}\n")
		io.WriteString(w, "printer.VisitStatementNode(stmt)\n")
		io.WriteString(w, "return printer.builder.String()\n")
	})

	const path = "../ast.go"

	file, err := os.Create(path)

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString("package main\n")

	exprNodes.write(file)
	stmtNodes.write(file)

	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
