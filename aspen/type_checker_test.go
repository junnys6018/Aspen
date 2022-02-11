package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

type TypeCheckerTestCase struct {
	fileName string
	source   []rune
	ast      Expression
	errors   []ErrorData
}

func (tc *TypeCheckerTestCase) run(t *testing.T) {
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

	scanLine := func(line string, lineNumber *int, col *int, message *string) {
		var end int

		// skip leading space
		for unicode.IsSpace(rune(line[end])) {
			end++
		}
		// skip line:col part
		for !unicode.IsSpace(rune(line[end])) {
			end++
		}

		fmt.Sscanf(line, "%d:%d", lineNumber, col)

		*message = line[end+1:]
	}

	scanner := bufio.NewScanner(strings.NewReader(comment))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var (
			lineNumber, col int
			message         string
		)

		scanLine(line, &lineNumber, &col, &message)

		errors = append(errors, ErrorData{lineNumber, col, message})
	}

	return &TypeCheckerTestCase{file, source, ast, errors}
}

func TestTypeChecker(t *testing.T) {
	matches, err := filepath.Glob("test_cases/type_checker/*.txt")

	if err != nil {
		t.Error("could not glob files")
		return
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := NewTypeCheckerTestCase(match, t)
		tc.run(t)
	}
}
