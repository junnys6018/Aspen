package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

const (
	// single character tokens
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_LEFT_SQUARE
	TOKEN_RIGHT_SQUARE
	TOKEN_COMMA
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR
	TOKEN_CARET
	TOKEN_PERCENT

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
	TOKEN_STRING_LITERAL
	TOKEN_FLOAT_LITERAL
	TOKEN_INT_LITERAL
	TOKEN_COMMENT

	// keywords
	TOKEN_ELSE
	TOKEN_FOR
	TOKEN_FN
	TOKEN_IF
	TOKEN_VOID
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_LET
	TOKEN_WHILE

	// types
	TOKEN_I64
	TOKEN_U64
	TOKEN_BOOL
	TOKEN_STRING
	TOKEN_DOUBLE

	TOKEN_EOF
)

type Token struct {
	tokenType TokenType
	line      int
	col       int
	value     interface{}
}

func (token Token) String() string {
	switch token.tokenType {
	case TOKEN_LEFT_PAREN:
		return "("
	case TOKEN_RIGHT_PAREN:
		return ")"
	case TOKEN_LEFT_BRACE:
		return "{"
	case TOKEN_RIGHT_BRACE:
		return "}"
	case TOKEN_LEFT_SQUARE:
		return "["
	case TOKEN_RIGHT_SQUARE:
		return "]"
	case TOKEN_COMMA:
		return ","
	case TOKEN_MINUS:
		return "-"
	case TOKEN_PLUS:
		return "+"
	case TOKEN_SEMICOLON:
		return ";"
	case TOKEN_SLASH:
		return "/"
	case TOKEN_STAR:
		return "*"
	case TOKEN_CARET:
		return "^"
	case TOKEN_PERCENT:
		return "%"
	case TOKEN_BANG:
		return "!"
	case TOKEN_BANG_EQUAL:
		return "!="
	case TOKEN_EQUAL:
		return "="
	case TOKEN_EQUAL_EQUAL:
		return "=="
	case TOKEN_GREATER:
		return ">"
	case TOKEN_GREATER_EQUAL:
		return ">="
	case TOKEN_LESS:
		return "<"
	case TOKEN_LESS_EQUAL:
		return "<="
	case TOKEN_AMP:
		return "&"
	case TOKEN_AMP_AMP:
		return "&&"
	case TOKEN_PIPE:
		return "|"
	case TOKEN_PIPE_PIPE:
		return "||"
	case TOKEN_IDENTIFIER:
		return fmt.Sprintf("%v", token.value.(string))
	case TOKEN_STRING_LITERAL:
		return fmt.Sprintf("\"%v\"", string(token.value.([]rune)))
	case TOKEN_FLOAT_LITERAL:
		return fmt.Sprintf("%.2f", token.value.(float64))
	case TOKEN_INT_LITERAL:
		return fmt.Sprintf("%d", token.value.(int64))
	case TOKEN_COMMENT:
		return "<comment>" // todo
	case TOKEN_ELSE:
		return "else"
	case TOKEN_FOR:
		return "for"
	case TOKEN_FN:
		return "fn"
	case TOKEN_IF:
		return "if"
	case TOKEN_VOID:
		return "void"
	case TOKEN_PRINT:
		return "print"
	case TOKEN_RETURN:
		return "return"
	case TOKEN_TRUE:
		return "true"
	case TOKEN_FALSE:
		return "false"
	case TOKEN_LET:
		return "let"
	case TOKEN_WHILE:
		return "while"
	case TOKEN_I64:
		return "i64"
	case TOKEN_U64:
		return "u64"
	case TOKEN_BOOL:
		return "bool"
	case TOKEN_STRING:
		return "string"
	case TOKEN_DOUBLE:
		return "double"
	case TOKEN_EOF:
		return "<eof>"
	}

	Unreachable("Token::String")
	return ""
}

type TokenStream []Token

func (tokens TokenStream) String() string {
	builder := strings.Builder{}

	lastLine := -1
	for i, token := range tokens {
		if token.line != lastLine {
			fmt.Fprintf(&builder, "%4d ", token.line)
			lastLine = token.line
		} else {
			builder.WriteString("   | ")
		}

		fmt.Fprintf(&builder, "%2d %v", token.col, token)

		if i != len(tokens)-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

func IsLetter(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

var KEYWORDS = map[string]TokenType{
	"else":   TOKEN_ELSE,
	"for":    TOKEN_FOR,
	"fn":     TOKEN_FN,
	"if":     TOKEN_IF,
	"void":   TOKEN_VOID,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
	"let":    TOKEN_LET,
	"while":  TOKEN_WHILE,
	"i64":    TOKEN_I64,
	"u64":    TOKEN_U64,
	"bool":   TOKEN_BOOL,
	"string": TOKEN_STRING,
	"double": TOKEN_DOUBLE,
}

// note: this function can be optimised, see: https://craftinginterpreters.com/scanning-on-demand.html#tries-and-state-machines
func matchKeyword(s string) (keyword TokenType, isKeyword bool) {
	keyword, ok := KEYWORDS[s]
	if ok {
		return keyword, true
	}
	return TOKEN_EOF, false
}

func ScanTokens(source []rune, errorReporter ErrorReporter) (TokenStream, error) {
	line := 1
	col := 1
	tokens := make(TokenStream, 0)
	i := 0

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

	simpleToken := func(tokenType TokenType) {
		tokens = append(tokens, Token{tokenType, line, col, nil})
	}

	commentToken := func(start, end, line, col int) {
		tokens = append(tokens, Token{TOKEN_COMMENT, line, col, string(source[start:end])})
	}

	singleLineComment := func() {
		oldLine := line
		oldCol := col
		start := i

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

		end := i
		if source[end-1] == '\n' {
			end-- // remove trailing newline
		}

		commentToken(start, end, oldLine, oldCol)
	}

	multiLineComment := func() {
		oldLine := line
		oldCol := col
		start := i

		col += 2
		terminated := false

		for !isAtEnd() {
			next := advance()
			col++
			if next == '*' {
				if !isAtEnd() {
					next = advance()
					col++
					if next == '/' {
						terminated = true
						break
					}
				}
			} else if next == '\n' {
				line++
				col = 1
			}
		}

		if !terminated {
			errorReporter.Push(line, col, "comment not terminated.")
		} else {
			end := i - 2
			commentToken(start, end, oldLine, oldCol)
		}
	}

	conditionalToken := func(ifNoMatch TokenType, ifMatch TokenType, matcher rune) {
		if match(matcher) {
			simpleToken(ifMatch)
			col += 2
		} else {
			simpleToken(ifNoMatch)
			col++
		}
	}

	stringToken := func() {
		oldCol := col
		col++

		start := i

		for !isAtEnd() && peek() != '\n' && peek() != '"' {
			advance()
			col++
		}

		if isAtEnd() || peek() == '\n' {
			errorReporter.Push(line, col, "string literal not terminated.")
			return
		}

		end := i

		if !match('"') {
			Unreachable("lexer.go: stringToken()")
		}
		col++

		tokens = append(tokens, Token{TOKEN_STRING_LITERAL, line, oldCol, source[start:end]})
	}

	numberToken := func() {
		oldCol := col
		col++

		start := i - 1
		isInteger := true

		for !isAtEnd() && IsDigit(peek()) {
			advance()
			col++
		}

		if match('.') {
			col++
			isInteger = false
			for !isAtEnd() && IsDigit(peek()) {
				advance()
				col++
			}
		}

		end := i
		if isInteger {
			value, err := strconv.ParseInt(string(source[start:end]), 10, 64)
			if err != nil {
				Unreachable("lexer.go: numberToken()")
			}
			tokens = append(tokens, Token{TOKEN_INT_LITERAL, line, oldCol, value})
		} else {
			value, err := strconv.ParseFloat(string(source[start:end]), 64)
			if err != nil {
				Unreachable("lexer.go: numberToken()")
			}
			tokens = append(tokens, Token{TOKEN_FLOAT_LITERAL, line, oldCol, value})
		}
	}

	identifierToken := func() {
		oldCol := col
		col++

		start := i - 1

		for !isAtEnd() && (IsLetter(peek()) || IsDigit(peek())) {
			advance()
			col++
		}

		end := i

		identifier := string(source[start:end])

		if keyword, isKeyword := matchKeyword(identifier); isKeyword {
			tokens = append(tokens, Token{keyword, line, oldCol, nil})
		} else {
			tokens = append(tokens, Token{TOKEN_IDENTIFIER, line, oldCol, identifier})
		}
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
			simpleToken(TOKEN_LEFT_PAREN)
			col++
		case ')':
			simpleToken(TOKEN_RIGHT_PAREN)
			col++
		case '{':
			simpleToken(TOKEN_LEFT_BRACE)
			col++
		case '}':
			simpleToken(TOKEN_RIGHT_BRACE)
			col++
		case '[':
			simpleToken(TOKEN_LEFT_SQUARE)
			col++
		case ']':
			simpleToken(TOKEN_RIGHT_SQUARE)
			col++
		case ',':
			simpleToken(TOKEN_COMMA)
			col++
		case '-':
			simpleToken(TOKEN_MINUS)
			col++
		case '+':
			simpleToken(TOKEN_PLUS)
			col++
		case ';':
			simpleToken(TOKEN_SEMICOLON)
			col++
		case '*':
			simpleToken(TOKEN_STAR)
			col++
		case '^':
			simpleToken(TOKEN_CARET)
			col++
		case '%':
			simpleToken(TOKEN_PERCENT)
			col++
		case '/':
			if match('/') {
				singleLineComment()
			} else if match('*') {
				multiLineComment()
			} else {
				simpleToken(TOKEN_SLASH)
				col++
			}
		case '!':
			conditionalToken(TOKEN_BANG, TOKEN_BANG_EQUAL, '=')
		case '=':
			conditionalToken(TOKEN_EQUAL, TOKEN_EQUAL_EQUAL, '=')
		case '>':
			conditionalToken(TOKEN_GREATER, TOKEN_GREATER_EQUAL, '=')
		case '<':
			conditionalToken(TOKEN_LESS, TOKEN_LESS_EQUAL, '=')
		case '&':
			conditionalToken(TOKEN_AMP, TOKEN_AMP_AMP, '&')
		case '|':
			conditionalToken(TOKEN_PIPE, TOKEN_PIPE_PIPE, '|')
		case '"':
			stringToken()
		default:
			if IsDigit(r) {
				numberToken()
			} else if IsLetter(r) {
				identifierToken()
			} else {
				errorReporter.Push(line, col, fmt.Sprintf("unexpected token \"%c\".", r))
				col++
			}
		}
	}

	simpleToken(TOKEN_EOF)

	if errorReporter.HadError() {
		return tokens, errorReporter
	} else {
		return tokens, nil
	}
}
