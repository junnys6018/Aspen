package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

type LexerTestCase struct {
	fileName    string
	test        string
	shouldError bool
	expect      []Token
	errors      []ErrorData
}

func LexerTestValuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	// check if a holds a []rune
	if v1, ok := a.([]rune); ok {
		// if so, check if b holds a []rune
		if v2, ok := b.([]rune); ok {
			// if so, do a deep comparison of the underlying []runes
			if len(v1) != len(v2) {
				return false
			}

			for i, r := range v1 {
				if r != v2[i] {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}

	// if not, compare a and b with normal interface comparison semantics
	return a == b
}

func (tc *LexerTestCase) Run(t *testing.T) {
	if tc == nil {
		return
	}

	source := []rune(tc.test)
	errorReporter := NewErrorReporter(source)
	tokens, err := ScanTokens(source, errorReporter)

	if tc.shouldError {
		if err == nil {
			t.Errorf("%s: expected err to be non-nil", tc.fileName)
			return
		}

		errors := err.(*AspenError).data
		if len(errors) != len(tc.errors) {
			t.Errorf("%s: expected len(errors) to be %d, got %d", tc.fileName, len(tc.errors), len(errors))
			return
		}

		for i, err := range errors {
			if err != tc.errors[i] {
				t.Errorf("%s: expected errors[%d] to be %v got %v", tc.fileName, i, tc.errors[i], err)
			}
		}
	} else {
		if err != nil {
			t.Errorf("%s: expected err to be nil", tc.fileName)
			return
		}

		if len(tokens) != len(tc.expect) {
			t.Errorf("%s: expected len(tokens) to be %d, got %d", tc.fileName, len(tc.expect), len(tokens))
			return
		}

		for i, expectToken := range tc.expect {
			if expectToken.tokenType != tokens[i].tokenType ||
				(expectToken.line != -1 && expectToken.line != tokens[i].line) ||
				(expectToken.col != -1 && expectToken.col != tokens[i].col) ||
				!LexerTestValuesEqual(expectToken.value, tokens[i].value) {
				t.Errorf("%s: expected tokens[%d] to be %+v, got %+v", tc.fileName, i, expectToken, tokens[i])
			}
		}
	}
}

func toTokenType(tokenType string) TokenType {
	switch tokenType {
	case "TOKEN_LEFT_PAREN":
		return TOKEN_LEFT_PAREN
	case "TOKEN_RIGHT_PAREN":
		return TOKEN_RIGHT_PAREN
	case "TOKEN_LEFT_BRACE":
		return TOKEN_LEFT_BRACE
	case "TOKEN_RIGHT_BRACE":
		return TOKEN_RIGHT_BRACE
	case "TOKEN_COMMA":
		return TOKEN_COMMA
	case "TOKEN_MINUS":
		return TOKEN_MINUS
	case "TOKEN_PLUS":
		return TOKEN_PLUS
	case "TOKEN_SEMICOLON":
		return TOKEN_SEMICOLON
	case "TOKEN_SLASH":
		return TOKEN_SLASH
	case "TOKEN_STAR":
		return TOKEN_STAR
	case "TOKEN_CARET":
		return TOKEN_CARET
	case "TOKEN_BANG":
		return TOKEN_BANG
	case "TOKEN_BANG_EQUAL":
		return TOKEN_BANG_EQUAL
	case "TOKEN_EQUAL":
		return TOKEN_EQUAL
	case "TOKEN_EQUAL_EQUAL":
		return TOKEN_EQUAL_EQUAL
	case "TOKEN_GREATER":
		return TOKEN_GREATER
	case "TOKEN_GREATER_EQUAL":
		return TOKEN_GREATER_EQUAL
	case "TOKEN_LESS":
		return TOKEN_LESS
	case "TOKEN_LESS_EQUAL":
		return TOKEN_LESS_EQUAL
	case "TOKEN_AMP":
		return TOKEN_AMP
	case "TOKEN_AMP_AMP":
		return TOKEN_AMP_AMP
	case "TOKEN_PIPE":
		return TOKEN_PIPE
	case "TOKEN_PIPE_PIPE":
		return TOKEN_PIPE_PIPE
	case "TOKEN_IDENTIFIER":
		return TOKEN_IDENTIFIER
	case "TOKEN_STRING":
		return TOKEN_STRING
	case "TOKEN_FLOAT":
		return TOKEN_FLOAT
	case "TOKEN_INT":
		return TOKEN_INT
	case "TOKEN_COMMENT":
		return TOKEN_COMMENT
	case "TOKEN_ELSE":
		return TOKEN_ELSE
	case "TOKEN_FOR":
		return TOKEN_FOR
	case "TOKEN_FN":
		return TOKEN_FN
	case "TOKEN_IF":
		return TOKEN_IF
	case "TOKEN_NIL":
		return TOKEN_NIL
	case "TOKEN_PRINT":
		return TOKEN_PRINT
	case "TOKEN_RETURN":
		return TOKEN_RETURN
	case "TOKEN_TRUE":
		return TOKEN_TRUE
	case "TOKEN_FALSE":
		return TOKEN_FALSE
	case "TOKEN_LET":
		return TOKEN_LET
	case "TOKEN_WHILE":
		return TOKEN_WHILE
	case "TOKEN_EOF":
		return TOKEN_EOF
	}

	Unreachable("lexer_test.go: toTokenType()")
	return 0
}

func LexerTestGetValue(line string, tokenType string) (interface{}, error) {
	end := 0
	// skip line:col part
	for !unicode.IsSpace(rune(line[end])) {
		end++
	}
	// skip whitespace
	for unicode.IsSpace(rune(line[end])) {
		end++
	}
	// skip token name part
	for !unicode.IsSpace(rune(line[end])) {
		end++
	}
	// skip whitespace
	for unicode.IsSpace(rune(line[end])) {
		end++
	}

	valueString := line[end:]

	switch tokenType {
	case "TOKEN_INT":
		var i int64
		fmt.Sscanf(valueString, "%d", &i)
		return i, nil
	case "TOKEN_STRING":
		unescaped, err := UnescapeString(valueString)
		if err != nil {
			return nil, err
		}
		return []rune(unescaped), nil
	case "TOKEN_FLOAT":
		var f float64
		fmt.Sscanf(valueString, "%f", &f)
		return f, nil
	case "TOKEN_COMMENT":
		unescaped, err := UnescapeString(valueString)
		if err != nil {
			return nil, err
		}
		return unescaped, nil
	case "TOKEN_IDENTIFIER":
		return valueString, nil
	default:
		return nil, errors.New("unknown token")
	}
}

func NewLexerTestCase(file string, t *testing.T) TestCase {
	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s\n", file)
		return nil
	}

	var tc LexerTestCase
	tc.fileName = file

	scanner := bufio.NewScanner(bytes.NewReader(data))

	scanner.Scan()
	line := scanner.Text()

	// scan status
	var status string
	fmt.Sscanf(line, "EXPECT %s", &status)
	if status != "FAILURE" && status != "SUCCESS" {
		t.Errorf("failed to parse %s", file)
	}

	tc.shouldError = status == "FAILURE"

	// scan expect block
	if tc.shouldError {
		for scanner.Scan() {
			line := scanner.Text()

			if line == "BEGIN TOKENS" {
				break
			}

			var (
				lineNumber int
				col        int
				message    string
			)

			ScanErrorMessage(line, &lineNumber, &col, &message)

			tc.errors = append(tc.errors, ErrorData{lineNumber, col, message})
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()

			if line == "BEGIN TOKENS" {
				break
			}

			var (
				lineNumber int
				col        int
				tokenType  string
			)
			fmt.Sscanf(line, "%d:%d %s", &lineNumber, &col, &tokenType)

			var value interface{}
			hasValue := tokenType == "TOKEN_INT" || tokenType == "TOKEN_STRING" || tokenType == "TOKEN_FLOAT" || tokenType == "TOKEN_COMMENT" || tokenType == "TOKEN_IDENTIFIER"

			if hasValue {
				value, err = LexerTestGetValue(line, tokenType)
				if err != nil {
					t.Errorf("failed to parse %s: %v", file, err)
					return nil
				}
			}

			tc.expect = append(tc.expect, Token{tokenType: toTokenType(tokenType), line: lineNumber, col: col, value: value})
		}
	}

	// scan tokens block
	source := string(data)

	idx := strings.Index(source, "BEGIN TOKENS")
	if idx >= 0 {
		tc.test = source[idx+len("BEGIN TOKENS")+1:]
	} else {
		t.Errorf("failed to parse %s", file)
		return nil
	}

	return &tc
}

func TestLexer(t *testing.T) {
	matches, err := filepath.Glob("test_cases/lexer/*.txt")

	if err != nil {
		t.Error("could not glob files")
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := NewLexerTestCase(match, t)
		tc.Run(t)
	}
}
