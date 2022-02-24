package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unicode"
)

type TypeCheckerTestCase struct {
	fileName string
	source   []rune
	ast      Program
	errors   []ErrorData
}

func (tc *TypeCheckerTestCase) Run(t *testing.T) {
	if tc == nil {
		return
	}

	errorReporter := NewErrorReporter(tc.source)
	err := TypeCheck(tc.ast, errorReporter)

	if err == nil {
		t.Errorf("%s: expected err to be non-nil", tc.fileName)
		return
	}

	errors := err.(*AspenError).data

	if len(tc.errors) != len(errors) {
		t.Errorf("%s: expected len(errors) to be %v got %v", tc.fileName, len(tc.errors), len(errors))
		return
	}

	for i, err := range errors {
		if err != tc.errors[i] {
			t.Errorf("%s: expected errors[%d] to be %v got %v", tc.fileName, i, tc.errors[i], err)
		}
	}
}

func NewTypeCheckerTestCase(file string, t *testing.T) *TypeCheckerTestCase {
	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s\n", file)
		return nil
	}

	source := []rune(string(data))
	errorReporter := NewErrorReporter(source)
	tokens, err := ScanTokens(source, errorReporter)

	if err != nil {
		t.Errorf("%s: failed to scan tokens\n %v", file, err)
		return nil
	}

	errorReporter = NewErrorReporter(source)
	ast, err := Parse(tokens, errorReporter)

	if err != nil {
		t.Errorf("%s: failed to parse source code\n %v", file, err)
		return nil
	}

	comment := tokens[0].value.(string)
	var errors []ErrorData

	pos := 0

	for pos < len(comment) {
		var (
			lineNumber, col int
			message         string
		)

		// skip leading space
		for unicode.IsSpace(rune(comment[pos])) {
			pos++
		}

		fmt.Sscanf(comment[pos:], "%d:%d", &lineNumber, &col)

		// skip line:col part
		for !unicode.IsSpace(rune(comment[pos])) {
			pos++
		}

		// skip leading space
		for unicode.IsSpace(rune(comment[pos])) {
			pos++
		}

		var terminator byte

		if comment[pos] == '`' {
			terminator = '`'
			pos++
		} else {
			terminator = '\n'
		}

		end := pos
		for comment[end] != terminator {
			end++
		}
		message = comment[pos:end]
		pos = end + 1

		// consume newline
		if terminator == '`' {
			pos++
		}

		errors = append(errors, ErrorData{lineNumber, col, message})
	}

	return &TypeCheckerTestCase{file, source, ast, errors}
}

func TestTypeChecker(t *testing.T) {
	Initialize()
	matches, err := filepath.Glob("test_cases/type_checker/*.txt")

	if err != nil {
		t.Error("could not glob files")
		return
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := NewTypeCheckerTestCase(match, t)
		tc.Run(t)
	}
}
