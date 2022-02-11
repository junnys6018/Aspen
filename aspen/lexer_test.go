package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type LexerTestCase struct {
	fileName    string
	test        string
	shouldError bool
	expect      []Token
	errors      ScanError
}

func valuesEqual(a, b interface{}) bool {
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

func (tc *LexerTestCase) run(t *testing.T) {
	if tc == nil {
		return
	}

	tokens, err := ScanTokens([]rune(tc.test))

	if tc.shouldError {
		if err == nil {
			t.Errorf("%s: expected err to be non-nil", tc.fileName)
			return
		}

		errors := err.(ScanError)
		if len(errors.data) != len(tc.errors.data) {
			t.Errorf("%s: expected len(errors) to be %d, got %d", tc.fileName, len(tc.errors.data), len(errors.data))
			return
		}

		for i, datum := range tc.errors.data {
			if datum != errors.data[i] {
				t.Errorf("%s: expected errors[%d] to be %v, got %v", tc.fileName, i, datum, errors.data[i])
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
				!valuesEqual(expectToken.value, tokens[i].value) {
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
	default:
		panic(fmt.Sprintf("unknown token %s", tokenType))
	}
}

func newLexerTestCase(file string) *LexerTestCase {
	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s\n", file)
		return nil
	}

	var decoded struct {
		Test   string
		Expect struct {
			Error  bool
			Result []interface{}
		}
	}

	err = json.Unmarshal(data, &decoded)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshall %s\n", file)
		return nil
	}

	var tc LexerTestCase

	tc.fileName = file
	tc.test = decoded.Test
	tc.shouldError = decoded.Expect.Error

	if tc.shouldError {
		for _, err := range decoded.Expect.Result {
			m := err.(map[string]interface{})

			tc.errors.data = append(tc.errors.data, ScanErrorData{int(m["line"].(float64)), int(m["col"].(float64)), m["message"].(string)})
		}
	} else {
		for _, token := range decoded.Expect.Result {
			m := token.(map[string]interface{})

			tokenTypeI, ok := m["type"]
			if !ok {
				fmt.Fprintf(os.Stderr, "%s: bad json\n", file)
				return nil
			}

			tokenType := tokenTypeI.(string)

			lineI, ok := m["line"]

			var line int
			if ok {
				line = int(lineI.(float64))
			} else {
				line = -1
			}

			colI, ok := m["col"]

			var col int
			if ok {
				col = int(colI.(float64))
			} else {
				col = -1
			}

			value, ok := m["value"]
			if ok {
				if tokenType == "TOKEN_INT" {
					value = int64(value.(float64))
				} else if tokenType == "TOKEN_STRING" {
					value = []rune(value.(string))
				}
			}

			tc.expect = append(tc.expect, Token{toTokenType(tokenType), line, col, value})
		}
	}
	return &tc
}

func TestLexer(t *testing.T) {
	matches, err := filepath.Glob("test_cases/lexer/*.json")

	if err != nil {
		t.Error("could not glob files")
	}

	for _, match := range matches {
		fmt.Printf("%s\n", match)
		tc := newLexerTestCase(match)
		tc.run(t)
	}
}
