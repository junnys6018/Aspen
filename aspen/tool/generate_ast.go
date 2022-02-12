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

func (n *ASTNode) write(w io.Writer) {
	fmt.Fprintf(w, "type %sExpression struct {\n", n.name)
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

func (m *ASTMethod) write(w io.Writer, nodeName string) {
	fmt.Fprintf(w, "func (expr *%sExpression) %s(%s) %s {\n", nodeName, m.name, m.arguments.asArguments(), m.returnType)

	m.generator(w, nodeName)

	io.WriteString(w, "}\n")
}

type ASTNodes struct {
	nodes   []ASTNode
	methods []ASTMethod
}

func (a *ASTNodes) defineNode(name string, fields Fields) {
	a.nodes = append(a.nodes, ASTNode{name, fields})
}

func (a *ASTNodes) defineMethod(name string, arguments Fields, returnType string, generator ASTMethodGenerator) {
	a.methods = append(a.methods, ASTMethod{name, arguments, returnType, generator})
}

func (a *ASTNodes) writeVisitorInterface(w io.Writer) {
	io.WriteString(w, "type ExpressionVisitor interface {\n")

	for _, node := range a.nodes {
		fmt.Fprintf(w, "Visit%[1]s(expr *%[1]sExpression) interface{}\n", node.name)
	}

	io.WriteString(w, "}\n")
}

func (a *ASTNodes) writeExpressionInterface(w io.Writer) {
	io.WriteString(w, "type Expression interface {\n")

	for _, method := range a.methods {
		fmt.Fprintf(w, "%s(%s) %s\n", method.name, method.arguments.asArguments(), method.returnType)
	}

	io.WriteString(w, "}\n")
}

func (a *ASTNodes) writeNodes(w io.Writer) {
	for _, node := range a.nodes {
		node.write(w)
		for _, method := range a.methods {
			method.write(w, node.name)
		}
	}
}

func (a *ASTNodes) write(w io.Writer) {
	io.WriteString(w, "package main\n")

	a.writeVisitorInterface(w)
	a.writeExpressionInterface(w)
	a.writeNodes(w)
}

func GenerateASTCode() {
	a := &ASTNodes{}

	a.defineNode("Binary", Fields{
		{"left", "Expression"},
		{"right", "Expression"},
		{"operator", "Token"},
	})

	a.defineNode("Unary", Fields{
		{"operand", "Expression"},
		{"operator", "Token"},
	})

	a.defineNode("Literal", Fields{
		{"value", "Token"},
	})

	a.defineNode("Grouping", Fields{
		{"expr", "Expression"},
	})

	a.defineMethod("Accept", Fields{
		{"visitor", "ExpressionVisitor"},
	}, "interface{}", func(w io.Writer, nodeName string) {
		fmt.Fprintf(w, "return visitor.Visit%s(expr)\n", nodeName)
	})

	a.defineMethod("String", Fields{}, "string", func(w io.Writer, nodeName string) {
		io.WriteString(w, "printer := AstPrinter{}\n")
		io.WriteString(w, "printer.Visit(expr)\n")
		io.WriteString(w, "return printer.builder.String()\n")
	})

	const path = "../ast.go"

	file, err := os.Create(path)

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	a.write(file)
	file.Close()

	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
