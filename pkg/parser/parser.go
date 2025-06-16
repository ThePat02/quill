package parser

import (
	"quill/pkg/ast"
	"quill/pkg/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

type ParseError struct {
	Line    int
	Message string
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (*ast.Program, []ParseError) {
	var errors []ParseError = make([]ParseError, 0)

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isAtEnd() {
		if p.check(token.NEWLINE) || p.check(token.COMMENT) {
			p.advance() // Skip newlines and comments
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			errors = append(errors, *err)
			// Skip to next statement after error
			p.synchronize()
			continue
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, errors
}

func (p *Parser) parseStatement() (ast.Statement, *ParseError) {
	switch {
	case p.check(token.LABEL):
		return p.parseLabelStatement()
	case p.check(token.GOTO):
		return p.parseGotoStatement()
	case p.check(token.CHOICE):
		return p.parseChoiceStatement()
	case p.check(token.RANDOM):
		return p.parseRandomStatement()
	case p.check(token.END):
		return p.parseEndStatement()
	case p.check(token.IDENT):
		if p.checkNext(token.COLON) {
			return p.parseDialogStatement()
		}
	}

	// Unknown token error
	current := p.peek()
	p.advance() // Skip unrecognized token
	return nil, &ParseError{
		Line:    current.Line,
		Message: "Unexpected token: " + current.Lexeme,
	}
}

func (p *Parser) parseLabelStatement() (ast.Statement, *ParseError) {
	labelToken := p.peek()
	p.advance() // consume LABEL

	if !p.check(token.IDENT) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected identifier after LABEL",
		}
	}

	name := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume identifier

	return &ast.LabelStatement{
		Token: labelToken,
		Name:  name,
	}, nil
}

func (p *Parser) parseGotoStatement() (ast.Statement, *ParseError) {
	gotoToken := p.peek()
	p.advance() // consume GOTO

	if !p.check(token.IDENT) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected identifier after GOTO",
		}
	}

	label := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume identifier

	return &ast.GotoStatement{
		Token: gotoToken,
		Label: label,
	}, nil
}

func (p *Parser) parseEndStatement() (ast.Statement, *ParseError) {
	endToken := p.peek()
	p.advance() // consume END

	return &ast.EndStatement{
		Token: endToken,
	}, nil
}

func (p *Parser) parseDialogStatement() (ast.Statement, *ParseError) {
	character := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume character name

	colonToken := p.peek()
	p.advance() // consume ':'

	if !p.check(token.STRING) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected string literal after ':'",
		}
	}

	text := &ast.StringLiteral{
		Token: p.peek(),
		Value: p.peek().Literal.(string),
	}
	p.advance() // consume string

	// Parse optional tags
	var tags *ast.TagList
	if p.check(token.LBRACKET) {
		var err *ParseError
		tags, err = p.parseTagList()
		if err != nil {
			return nil, err
		}
	}

	return &ast.DialogStatement{
		Character: character,
		Colon:     colonToken,
		Text:      text,
		Tags:      tags,
	}, nil
}

func (p *Parser) parseTagList() (*ast.TagList, *ParseError) {
	lbracketToken := p.peek()
	p.advance() // consume '['

	var tags []*ast.Identifier

	for !p.check(token.RBRACKET) && !p.isAtEnd() {
		if !p.check(token.IDENT) {
			return nil, &ParseError{
				Line:    p.peek().Line,
				Message: "Expected identifier in tag list",
			}
		}

		tag := &ast.Identifier{
			Token: p.peek(),
			Value: p.peek().Lexeme,
		}
		p.advance()
		tags = append(tags, tag)

		// Handle comma separation
		if p.check(token.COMMA) {
			p.advance()
		} else if !p.check(token.RBRACKET) {
			return nil, &ParseError{
				Line:    p.peek().Line,
				Message: "Expected ',' or ']' in tag list",
			}
		}
	}

	if !p.check(token.RBRACKET) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected ']' to close tag list",
		}
	}

	p.advance() // consume ']'

	return &ast.TagList{
		Token: lbracketToken,
		Tags:  tags,
	}, nil
}

func (p *Parser) parseChoiceStatement() (ast.Statement, *ParseError) {
	choiceToken := p.peek()
	p.advance() // consume CHOICE

	if !p.check(token.LBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '{' after CHOICE",
		}
	}

	p.advance() // consume '{'

	var options []*ast.ChoiceOption

	for !p.check(token.RBRACE) && !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.advance()
			continue
		}

		option, err := p.parseChoiceOption()
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}

		// Consume optional comma
		if p.check(token.COMMA) {
			p.advance()
		}
	}

	if !p.check(token.RBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '}' to close CHOICE block",
		}
	}

	p.advance() // consume '}'

	return &ast.ChoiceStatement{
		Token:   choiceToken,
		Options: options,
	}, nil
}

func (p *Parser) parseChoiceOption() (*ast.ChoiceOption, *ParseError) {
	if !p.check(token.STRING) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected string literal for choice option",
		}
	}

	text := &ast.StringLiteral{
		Token: p.peek(),
		Value: p.peek().Literal.(string),
	}
	p.advance()

	if !p.check(token.LBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '{' after choice option text",
		}
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	// Parse optional tags after the body
	var tags *ast.TagList
	if p.check(token.LBRACKET) {
		tags, err = p.parseTagList()
		if err != nil {
			return nil, err
		}
	}

	return &ast.ChoiceOption{
		Text: text,
		Body: body,
		Tags: tags,
	}, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, *ParseError) {
	lbraceToken := p.peek()
	p.advance() // consume '{'

	var statements []ast.Statement

	for !p.check(token.RBRACE) && !p.isAtEnd() {
		// Skip newlines and comments
		if p.check(token.NEWLINE) || p.check(token.COMMENT) {
			p.advance()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	if !p.check(token.RBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '}' to close block",
		}
	}

	p.advance() // consume '}'
	return &ast.BlockStatement{
		Token:      lbraceToken,
		Statements: statements,
	}, nil
}

func (p *Parser) parseRandomStatement() (ast.Statement, *ParseError) {
	randomToken := p.peek()
	p.advance() // consume RANDOM

	if !p.check(token.LBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '{' after RANDOM",
		}
	}

	p.advance() // consume '{'

	var options []*ast.RandomOption

	for !p.check(token.RBRACE) && !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.advance()
			continue
		}

		option, err := p.parseRandomOption()
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}

		// Consume optional comma
		if p.check(token.COMMA) {
			p.advance()
		}
	}

	if !p.check(token.RBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '}' to close RANDOM block",
		}
	}

	p.advance() // consume '}'

	return &ast.RandomStatement{
		Token:   randomToken,
		Options: options,
	}, nil
}

func (p *Parser) parseRandomOption() (*ast.RandomOption, *ParseError) {
	if !p.check(token.LBRACE) {
		return nil, &ParseError{
			Line:    p.peek().Line,
			Message: "Expected '{' for random option",
		}
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	// Parse optional tags after the body
	var tags *ast.TagList
	if p.check(token.LBRACKET) {
		tags, err = p.parseTagList()
		if err != nil {
			return nil, err
		}
	}

	return &ast.RandomOption{
		Body: body,
		Tags: tags,
	}, nil
}

// Helper function for error recovery
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.tokens[p.current-1].Type == token.NEWLINE {
			return
		}

		switch p.peek().Type {
		case token.LABEL, token.GOTO, token.CHOICE, token.END:
			return
		}

		p.advance()
	}
}
