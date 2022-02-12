package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		bytes, err := os.ReadFile(os.Args[1])

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot open file %v\n", os.Args[1])
			os.Exit(1)
		}

		source := []rune(string(bytes))

		errorReporter := NewErrorReporter(source)
		tokens, err := ScanTokens(source, errorReporter)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		errorReporter = NewErrorReporter(source)
		ast, err := Parse(tokens, errorReporter)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		errorReporter = NewErrorReporter(source)
		err = TypeCheck(ast, errorReporter)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		Interpret(ast)

	} else {
		fmt.Printf("usage: %v source.aspen\n", os.Args[0])
	}
}
