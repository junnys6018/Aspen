package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type ParserTestCase struct {
	fileName string
	source   string
	expect   string
	kind     string
}

func (tc *ParserTestCase) Run(t *testing.T) {
	if tc == nil {
		return
	}

	source := []rune(tc.source)
	errorReporter := NewErrorReporter(source)
	tokens, err := ScanTokens(source, errorReporter)

	if err != nil {
		t.Errorf("%s: failed to scan tokens\n %v", tc.fileName, err)
		return
	}

	var astString string

	switch tc.kind {
	case "TYPE PROGRAM":
		errorReporter = NewErrorReporter(source)
		ast, err := Parse(tokens, errorReporter)

		if err != nil {
			t.Errorf("%s: failed to parse source code\n %v", tc.fileName, err)
			return
		}

		astString = ast.String()
	case "TYPE EXPRESSION":
		errorReporter = NewErrorReporter(source)
		parser := Parser{tokens: tokens, current: 0, errorReporter: errorReporter}
		expr := parser.Expression()
		astString = expr.String()
	}

	if astString != tc.expect {
		t.Errorf("%s: expected ast to be %s, got %s", tc.fileName, tc.expect, astString)
		return
	}
}

func NewParserTestCase(file string) TestCase {
	data, err := os.ReadFile(file)
	content := string(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s\n", file)
		return nil
	}

	var kind, source, expect string

	i := strings.Index(content, "\n")
	j := i + 1 + strings.Index(content[i+1:], "\n")

	kind = content[:i]
	expect = content[i+1 : j]
	source = content[j+1:]

	return &ParserTestCase{file, source, expect, kind}
}

func TestParser(t *testing.T) {
	matches, err := filepath.Glob("test_cases/parser/*.txt")

	if err != nil {
		t.Error("could not glob files")
		return
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := NewParserTestCase(match)
		tc.Run(t)
	}
}
