package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func writeBinaryOperatorFunction(w io.Writer, operator string, functionName string, includeFloat bool) {
	fmt.Fprintf(w, "func %s(lhs, rhs interface{}) interface{} {\n", functionName)
	fmt.Fprintln(w, "switch lhsV := lhs.(type) {")
	fmt.Fprintln(w, "case int64:")
	fmt.Fprintln(w, "rhsV := rhs.(int64)")
	fmt.Fprintf(w, "return lhsV %s rhsV\n", operator)
	fmt.Fprintln(w, "case uint64:")
	fmt.Fprintln(w, "rhsV := rhs.(uint64)")
	fmt.Fprintf(w, "return lhsV %s rhsV\n", operator)
	if includeFloat {
		fmt.Fprintln(w, "case float64:")
		fmt.Fprintln(w, "rhsV := rhs.(float64)")
		fmt.Fprintf(w, "return lhsV %s rhsV\n", operator)
	}
	fmt.Fprintln(w, "}")

	fmt.Fprintf(w, "Unreachable(\"value.go: %s\")\n", functionName)
	fmt.Fprintln(w, "return nil")
	fmt.Fprintln(w, "}")
}

func GenerateValueCode() {
	path := "../value.go"

	file, err := os.Create(path)

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	defer file.Close()

	fmt.Fprintln(file, "package main")

	writeBinaryOperatorFunction(file, ">", "OperatorGreater", true)
	writeBinaryOperatorFunction(file, ">=", "OperatorGreaterEqual", true)
	writeBinaryOperatorFunction(file, "<", "OperatorLess", true)
	writeBinaryOperatorFunction(file, "<=", "OperatorLessEqual", true)
	writeBinaryOperatorFunction(file, "|", "OperatorPipe", false)
	writeBinaryOperatorFunction(file, "^", "OperatorCaret", false)
	writeBinaryOperatorFunction(file, "&", "OperatorAmp", false)
	writeBinaryOperatorFunction(file, "-", "OperatorMinus", true)
	writeBinaryOperatorFunction(file, "+", "OperatorPlus", true)
	writeBinaryOperatorFunction(file, "/", "OperatorSlash", true)
	writeBinaryOperatorFunction(file, "*", "OperatorStar", true)

	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
