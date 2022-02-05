package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	// single character tokens
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// one or two character tokens
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	TOKEN_AMP
	TOKEN_AMP_AMP
	TOKEN_PIPE
	TOKEN_PIPE_PIPE

	// literals
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_FLOAT
	TOKEN_INT

	// keywords
	TOKEN_ELSE
	TOKEN_FOR
	TOKEN_FN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_LET
	TOKEN_WHILE

	TOKEN_EOF
)

type Token struct {
	tokenType TokenType
	line      int
	col       int
	value     interface{}
}

type TokenStream []Token

func (tokens TokenStream) String() string {
	builder := strings.Builder{}
	for _, token := range tokens {
		fmt.Fprintf(&builder, "%+v\n", token)
	}
	return builder.String()
}

type ScanError struct {
	source []rune
	errors []int
}

func (e ScanError) Error() string {
	builder := strings.Builder{}

	for i, idx := range e.errors {
		builder.WriteString(ErrorString(e.source, fmt.Sprintf("unexpected token \"%c\".", e.source[idx]), idx))
		if i != len(e.errors)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}

func ScanTokens(source []rune) (TokenStream, error) {
	line := 1
	col := 1
	tokens := make(TokenStream, 0)
	i := 0

	err := ScanError{source, make([]int, 0)}

	simpleToken := func(tokens TokenStream, tokenType TokenType) TokenStream {
		tokens = append(tokens, Token{tokenType, line, col, nil})
		return tokens
	}

	advance := func() rune {
		i++
		return source[i-1]
	}

	isAtEnd := func() bool {
		return i == len(source)
	}

	peek := func() rune {
		return source[i]
	}

	match := func(r rune) bool {
		if !isAtEnd() && r == peek() {
			advance()
			return true
		}
		return false
	}

	conditionalToken := func(tokens TokenStream, ifNoMatch TokenType, ifMatch TokenType, matcher rune) TokenStream {
		if match(matcher) {
			tokens = simpleToken(tokens, ifMatch)
			col += 2
		} else {
			tokens = simpleToken(tokens, ifNoMatch)
			col++
		}
		return tokens
	}

	for !isAtEnd() {
		r := advance()

		// skip white space
		if r == ' ' || r == '\t' || r == '\r' {
			col++
			continue
		}

		if r == '\n' {
			line++
			col = 1
			continue
		}

		switch r {
		case '(':
			tokens = simpleToken(tokens, TOKEN_LEFT_PAREN)
			col++
		case ')':
			tokens = simpleToken(tokens, TOKEN_RIGHT_PAREN)
			col++
		case '{':
			tokens = simpleToken(tokens, TOKEN_LEFT_BRACE)
			col++
		case '}':
			tokens = simpleToken(tokens, TOKEN_RIGHT_BRACE)
			col++
		case ',':
			tokens = simpleToken(tokens, TOKEN_COMMA)
			col++
		case '-':
			tokens = simpleToken(tokens, TOKEN_MINUS)
			col++
		case '+':
			tokens = simpleToken(tokens, TOKEN_PLUS)
			col++
		case ';':
			tokens = simpleToken(tokens, TOKEN_SEMICOLON)
			col++
		case '*':
			tokens = simpleToken(tokens, TOKEN_STAR)
			col++
		case '/':
			if match('/') /* comment */ {
				col += 2
				for !isAtEnd() {
					next := advance()
					col++
					if next == '\n' {
						line++
						col = 1
						break
					}
				}
			} else /* token */ {
				tokens = simpleToken(tokens, TOKEN_SLASH)
				col++
			}
		case '!':
			tokens = conditionalToken(tokens, TOKEN_BANG, TOKEN_BANG_EQUAL, '=')
		case '=':
			tokens = conditionalToken(tokens, TOKEN_EQUAL, TOKEN_EQUAL_EQUAL, '=')
		case '>':
			tokens = conditionalToken(tokens, TOKEN_GREATER, TOKEN_GREATER_EQUAL, '=')
		case '<':
			tokens = conditionalToken(tokens, TOKEN_LESS, TOKEN_LESS_EQUAL, '=')
		case '&':
			tokens = conditionalToken(tokens, TOKEN_AMP, TOKEN_AMP_AMP, '&')
		case '|':
			tokens = conditionalToken(tokens, TOKEN_PIPE, TOKEN_PIPE_PIPE, '|')
		default:
			if unicode.IsDigit(r) {
				oldCol := col
				col++

				start := i - 1
				isInteger := true

				for !isAtEnd() && unicode.IsDigit(peek()) {
					advance()
					col++
				}

				if match('.') {
					col++
					isInteger = false
					for !isAtEnd() && unicode.IsDigit(peek()) {
						advance()
						col++
					}
				}

				end := i
				if isInteger {
					value, err := strconv.ParseInt(string(source[start:end]), 10, 64)
					if err != nil {
						panic("bug: should never error here")
					}
					tokens = append(tokens, Token{TOKEN_INT, line, oldCol, value})
				} else {
					value, err := strconv.ParseFloat(string(source[start:end]), 64)
					if err != nil {
						panic("bug: should never error here")
					}
					tokens = append(tokens, Token{TOKEN_FLOAT, line, oldCol, value})
				}

			} else {
				err.errors = append(err.errors, i-1)
			}
		}
	}

	tokens = simpleToken(tokens, TOKEN_EOF)

	if len(err.errors) == 0 {
		return tokens, nil
	} else {
		return nil, err
	}
}
