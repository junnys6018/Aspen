package main

type Parser struct {
	tokens        TokenStream
	current       int
	errorReporter ErrorReporter
}

func (p *Parser) Synchronize() {
	p.Advance()

	for !p.IsAtEnd() {
		if p.Previous().tokenType == TOKEN_SEMICOLON {
			return
		}

		switch p.Peek().tokenType {
		case TOKEN_FN, TOKEN_LET, TOKEN_FOR, TOKEN_IF, TOKEN_WHILE, TOKEN_PRINT, TOKEN_RETURN:
			return
		}

		p.Advance()
	}
}

// Grammar

// Statements

func (p *Parser) Declaration() Statement {
	defer func() {
		if r := recover(); r != nil {
			err := r.(ErrorData)
			p.errorReporter.Push(err.line, err.col, err.message)
			p.Synchronize()
		}
	}()

	if p.Match(TOKEN_LET) {
		return p.LetStatement()
	}

	return p.Statement()
}

func (p *Parser) Statement() Statement {
	if p.Match(TOKEN_PRINT) {
		return p.PrintStatement()
	}
	if p.Match(TOKEN_LEFT_BRACE) {
		return p.BlockStatement()
	}
	return p.ExpressionStatement()
}

func (p *Parser) PrintStatement() Statement {
	expr := p.Expression()
	p.Consume(TOKEN_SEMICOLON, "expected \";\" after expression.")
	return &PrintStatement{expr: expr}
}

func (p *Parser) BlockStatement() Statement {
	statements := make([]Statement, 0)

	for !p.Check(TOKEN_RIGHT_BRACE) && !p.IsAtEnd() {
		statements = append(statements, p.Declaration())
	}

	p.Consume(TOKEN_RIGHT_BRACE, "expected \"}\" after block.")
	return &BlockStatement{statements: statements}
}

func (p *Parser) ExpressionStatement() Statement {
	expr := p.Expression()
	p.Consume(TOKEN_SEMICOLON, "expected \";\" after expression.")
	return &ExpressionStatement{expr: expr}
}

func (p *Parser) LetStatement() Statement {
	name := p.Consume(TOKEN_IDENTIFIER, "expected variable name.")
	atype := p.Type()

	var initializer Expression

	if p.Match(TOKEN_EQUAL) {
		initializer = p.Expression()
	}

	p.Consume(TOKEN_SEMICOLON, "expected \";\" after variable declaration.")
	return &LetStatement{name: *name, initializer: initializer, atype: atype}
}

// Expressions

func (p *Parser) Expression() Expression {
	return p.LogicOr()
}

func (p *Parser) LogicOr() Expression {
	expr := p.LogicAnd()

	for p.Match(TOKEN_PIPE_PIPE) {
		operator := p.Previous()
		right := p.LogicAnd()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) LogicAnd() Expression {
	expr := p.Equality()

	for p.Match(TOKEN_AMP_AMP) {
		operator := p.Previous()
		right := p.Equality()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) Equality() Expression {
	expr := p.Comparison()

	for p.Match(TOKEN_EQUAL_EQUAL, TOKEN_BANG_EQUAL) {
		operator := p.Previous()
		right := p.Comparison()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) Comparison() Expression {
	expr := p.BitOr()

	for p.Match(TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL) {
		operator := p.Previous()
		right := p.BitOr()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) BitOr() Expression {
	expr := p.BitXor()

	for p.Match(TOKEN_PIPE) {
		operator := p.Previous()
		right := p.BitXor()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) BitXor() Expression {
	expr := p.BitAnd()

	for p.Match(TOKEN_CARET) {
		operator := p.Previous()
		right := p.BitAnd()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) BitAnd() Expression {
	expr := p.Term()

	for p.Match(TOKEN_AMP) {
		operator := p.Previous()
		right := p.Term()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) Term() Expression {
	expr := p.Factor()

	for p.Match(TOKEN_MINUS, TOKEN_PLUS) {
		operator := p.Previous()
		right := p.Factor()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) Factor() Expression {
	expr := p.Unary()

	for p.Match(TOKEN_SLASH, TOKEN_STAR) {
		operator := p.Previous()
		right := p.Unary()
		expr = &BinaryExpression{left: expr, right: right, operator: *operator}
	}

	return expr
}

func (p *Parser) Unary() Expression {
	if p.Match(TOKEN_BANG, TOKEN_MINUS) {
		operator := p.Previous()
		right := p.Unary()
		return &UnaryExpression{operand: right, operator: *operator}
	}

	return p.Primary()
}

func (p *Parser) Primary() Expression {
	if p.Match(TOKEN_FALSE, TOKEN_TRUE, TOKEN_INT_LITERAL, TOKEN_FLOAT_LITERAL, TOKEN_STRING_LITERAL) {
		return &LiteralExpression{value: *p.Previous()}
	}

	if p.Match(TOKEN_LEFT_PAREN) {
		expr := p.Expression()
		p.Consume(TOKEN_RIGHT_PAREN, "expected \")\" after expression.")
		return &GroupingExpression{expr: expr}
	}

	if p.Match(TOKEN_IDENTIFIER) {
		return &IdentifierExpression{name: *p.Previous()}
	}

	token := p.Peek()

	panic(ErrorData{token.line, token.col, "expected expression."})
}

// Types

func (p *Parser) Type() *Type {
	return p.Function()
}

func (p *Parser) Function() *Type {
	if p.Match(TOKEN_FN) {
		parameters := make([]*Type, 0)
		p.Consume(TOKEN_LEFT_PAREN, "expected \"(\".")

		if !p.Check(TOKEN_RIGHT_PAREN) {
			parameters = append(parameters, p.Type())
			for p.Match(TOKEN_COMMA) {
				parameters = append(parameters, p.Type())
			}
		}

		p.Consume(TOKEN_RIGHT_PAREN, "expected \")\".")

		var returnType *Type
		if p.Match(TOKEN_VOID) {
			returnType = SimpleType(TYPE_VOID)
		} else {
			returnType = p.Type()
		}
		return &Type{kind: TYPE_FUNCTION, other: FunctionType{parameters: parameters, returnType: returnType}}
	}

	return p.Slice()
}

func (p *Parser) Slice() *Type {
	atype := p.Primitive()

	for p.Match(TOKEN_LEFT_SQUARE) {
		p.Consume(TOKEN_RIGHT_SQUARE, "expected \"]\" after type definition.")
		atype = &Type{kind: TYPE_SLICE, other: SliceType{of: atype}}
	}
	return atype
}

func (p *Parser) Primitive() *Type {
	convert := func(tokenType TokenType) TypeEnum {
		switch tokenType {
		case TOKEN_I64:
			return TYPE_I64
		case TOKEN_U64:
			return TYPE_U64
		case TOKEN_BOOL:
			return TYPE_BOOL
		case TOKEN_STRING:
			return TYPE_STRING
		case TOKEN_DOUBLE:
			return TYPE_DOUBLE
		}
		Unreachable("Parser::Primitive convert()")
		return 0
	}

	if p.Match(TOKEN_I64, TOKEN_U64, TOKEN_BOOL, TOKEN_STRING, TOKEN_DOUBLE) {
		return SimpleType(convert(p.Previous().tokenType))
	}

	if p.Match(TOKEN_LEFT_PAREN) {
		atype := p.Type()
		p.Consume(TOKEN_RIGHT_PAREN, "expected \")\" after type definition.")
		return atype
	}

	token := p.Peek()
	panic(ErrorData{token.line, token.col, "expected type definition."})
}

// Helpers

func (p *Parser) Consume(tokenType TokenType, message string) *Token {
	token := p.Peek()
	if token.tokenType == tokenType {
		return p.Advance()
	}

	panic(ErrorData{token.line, token.col, message})
}

func (p *Parser) Match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.Check(tokenType) {
			p.Advance()
			return true
		}
	}
	return false
}

func (p *Parser) Check(tokenType TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().tokenType == tokenType
}

func (p *Parser) Advance() *Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.Previous()
}

func (p *Parser) IsAtEnd() bool {
	return p.Peek().tokenType == TOKEN_EOF
}

func (p *Parser) Peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) Previous() *Token {
	return &p.tokens[p.current-1]
}

func Parse(tokens TokenStream, errorReporter ErrorReporter) (Program, error) {
	// remove comment tokens
	filteredTokens := make(TokenStream, 0, len(tokens))
	for _, token := range tokens {
		if token.tokenType != TOKEN_COMMENT {
			filteredTokens = append(filteredTokens, token)
		}
	}

	parser := Parser{tokens: filteredTokens, current: 0, errorReporter: errorReporter}

	statements := make(Program, 0)

	for !parser.IsAtEnd() {
		statements = append(statements, parser.Declaration())
	}

	if errorReporter.HadError() {
		return nil, errorReporter
	} else {
		return statements, nil
	}
}
