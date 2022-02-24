package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Program []Statement

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

func ParseSource(source []rune) (Program, error) {
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

func TypeCheckSource(source []rune) (Program, error) {
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

func Initialize() {
	start := time.Now()
	DefineNativeFunction(SimpleFunction(TYPE_I64), "clock", func(args []interface{}) interface{} {
		return time.Since(start).Microseconds()
	})

	DefineNativeFunction(SimpleFunction(TYPE_STRING, TYPE_STRING), "__TESTFN__", func(args []interface{}) interface{} {
		arg0 := args[0].([]rune)
		return []rune(fmt.Sprintf("__TESTFN__(%s)", string(arg0)))
	})

	// string related functions

	DefineNativeFunction(SimpleFunction(TYPE_STRING, TYPE_I64), "itoa", func(args []interface{}) interface{} {
		arg0 := args[0].(int64)
		return []rune(fmt.Sprintf("%d", arg0))
	})

	DefineNativeFunction(SimpleFunction(TYPE_STRING, TYPE_DOUBLE), "ftoa", func(args []interface{}) interface{} {
		arg0 := args[0].(float64)
		return []rune(fmt.Sprintf("%f", arg0))
	})

	DefineNativeFunction(SimpleFunction(TYPE_I64, TYPE_STRING), "atoi", func(args []interface{}) interface{} {
		arg0 := string(args[0].([]rune))
		i, _ := strconv.Atoi(arg0)
		return i
	})

	DefineNativeFunction(SimpleFunction(TYPE_DOUBLE, TYPE_STRING), "atof", func(args []interface{}) interface{} {
		arg0 := string(args[0].([]rune))
		f, _ := strconv.ParseFloat(arg0, 64)
		return f
	})

	// type casting

	AddConversion(SimpleType(TYPE_I64), SimpleType(TYPE_U64), func(from interface{}) interface{} {
		v := from.(int64)
		return uint64(v)
	})

	AddConversion(SimpleType(TYPE_I64), SimpleType(TYPE_DOUBLE), func(from interface{}) interface{} {
		v := from.(int64)
		return float64(v)
	})

	AddConversion(SimpleType(TYPE_U64), SimpleType(TYPE_I64), func(from interface{}) interface{} {
		v := from.(uint64)
		return int64(v)
	})

	AddConversion(SimpleType(TYPE_U64), SimpleType(TYPE_DOUBLE), func(from interface{}) interface{} {
		v := from.(uint64)
		return float64(v)
	})

	AddConversion(SimpleType(TYPE_DOUBLE), SimpleType(TYPE_I64), func(from interface{}) interface{} {
		v := from.(float64)
		return int64(v)
	})

	AddConversion(SimpleType(TYPE_DOUBLE), SimpleType(TYPE_U64), func(from interface{}) interface{} {
		v := from.(float64)
		return uint64(v)
	})
}

func main() {
	Initialize()

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
