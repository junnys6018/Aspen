package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type ParserTestCase struct {
	fileName string
	source   string
	expect   string
}

func (tc *ParserTestCase) run(t *testing.T) {
	if tc == nil {
		return
	}

	tokens, err := ScanTokens([]rune(tc.source))

	if err != nil {
		t.Errorf("%s: failed to scan tokens\n %v", tc.fileName, err)
		return
	}

	ast, err := Parse(tokens)

	if err != nil {
		t.Errorf("%s: failed to parse source code\n %v", tc.fileName, err)
		return
	}

	astString := ast.String()

	if astString != tc.expect {
		t.Errorf("%s: expected ast to be %s, got %s", tc.fileName, tc.expect, astString)
		return
	}
}

func newParserTestCase(file string) *ParserTestCase {
	data, err := os.ReadFile(file)
	content := string(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s\n", file)
		return nil
	}

	var source, expect string

	for i, r := range content {
		if r == '\n' {
			expect = content[:i]
			source = content[i+1:]
			break
		}
	}

	return &ParserTestCase{file, source, expect}
}

func TestParser(t *testing.T) {
	matches, err := filepath.Glob("test_cases/parser/*.txt")

	if err != nil {
		t.Error("could not glob files")
		return
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := newParserTestCase(match)
		tc.run(t)
	}
}
