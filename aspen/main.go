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

		source := string(bytes)

		tokens, err := ScanTokens([]rune(source))

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		ast, err := Parse(tokens)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		fmt.Printf("%v\n", ast)
	} else {
		fmt.Printf("usage: %v source.aspen\n", os.Args[0])
	}
}
