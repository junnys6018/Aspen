package main

import (
	"fmt"
	"os"
)

const helpString = `useage: aspen [<options>] <path>

Options
    <path>
    The path to the aspen source file to execute

    -i or --interpret
    Execute the program using the tree walk implementation

    -l or --lex
    Do lexical analysis on the source code and print out the tokens scanned

    -p or --parse
    Parse the source code and print out the ast as an S-expression

    -t or --type-check
    Run the type checker on the program but do not execute it`

func OpenFile(path string) ([]rune, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("error: cannot open file %s", path)
	}

	return []rune(string(bytes)), nil
}

func ScanSource(source []rune) (TokenStream, error) {
	errorReporter := NewErrorReporter(source)
	tokens, err := ScanTokens(source, errorReporter)

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func ParseSource(source []rune) (Expression, error) {
	tokens, err := ScanSource(source)

	if err != nil {
		return nil, err
	}

	errorReporter := NewErrorReporter(source)
	ast, err := Parse(tokens, errorReporter)

	if err != nil {
		return nil, err
	}

	return ast, nil
}

func TypeCheckSource(source []rune) (Expression, error) {
	ast, err := ParseSource(source)

	if err != nil {
		return nil, err
	}

	errorReporter := NewErrorReporter(source)
	err = TypeCheck(ast, errorReporter)

	if err != nil {
		return nil, err
	}

	return ast, nil
}

func ExecuteSource(source []rune) error {
	ast, err := TypeCheckSource(source)

	if err != nil {
		return err
	}

	err = Interpret(ast)
	return err
}

func ExecuteFile(path string) error {
	source, err := OpenFile(path)
	if err != nil {
		return err
	}

	err = ExecuteSource(source)
	if err != nil {
		return err
	}

	return nil
}

func Check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 2 {
		path := os.Args[1]

		err := ExecuteFile(path)
		Check(err)
	} else if len(os.Args) == 3 {
		flag := os.Args[1]

		path := os.Args[2]
		source, err := OpenFile(path)
		Check(err)

		switch flag {
		case "-i", "--interpret":
			err = ExecuteSource(source)
			Check(err)
		case "-l", "--lex":
			tokens, err := ScanSource(source)
			Check(err)

			fmt.Println(tokens)
		case "-p", "--parse":
			ast, err := ParseSource(source)
			Check(err)

			fmt.Println(ast)
		case "-t", "--type-check":
			_, err := TypeCheckSource(source)
			Check(err)

			os.Exit(0)
		default:
			fmt.Println(helpString)
		}
	} else {
		fmt.Println(helpString)
	}
}
