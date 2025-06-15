package parser

import (
	"quill/pkg/ast"
	"quill/pkg/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isAtEnd() {
		if p.check(token.NEWLINE) || p.check(token.COMMENT) {
			p.advance() // Skip newlines and comments
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch {
	case p.check(token.LABEL):
		return p.parseLabelStatement()
	case p.check(token.GOTO):
		return p.parseGotoStatement()
	case p.check(token.CHOICE):
		return p.parseChoiceStatement()
	case p.check(token.END):
		return p.parseEndStatement()
	case p.check(token.IDENT):
		if p.checkNext(token.COLON) {
			return p.parseDialogStatement()
		}
	}

	p.advance() // Skip unrecognized token
	return nil
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	labelToken := p.peek()
	p.advance() // consume LABEL

	if !p.check(token.IDENT) {
		return nil
	}

	name := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume identifier

	return &ast.LabelStatement{
		Token: labelToken,
		Name:  name,
	}
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	gotoToken := p.peek()
	p.advance() // consume GOTO

	if !p.check(token.IDENT) {
		return nil
	}

	label := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume identifier

	return &ast.GotoStatement{
		Token: gotoToken,
		Label: label,
	}
}

func (p *Parser) parseEndStatement() *ast.EndStatement {
	endToken := p.peek()
	p.advance() // consume END

	return &ast.EndStatement{
		Token: endToken,
	}
}

func (p *Parser) parseDialogStatement() *ast.DialogStatement {
	character := &ast.Identifier{
		Token: p.peek(),
		Value: p.peek().Lexeme,
	}
	p.advance() // consume character name

	colonToken := p.peek()
	p.advance() // consume ':'

	if !p.check(token.STRING) {
		return nil
	}

	text := &ast.StringLiteral{
		Token: p.peek(),
		Value: p.peek().Literal.(string),
	}
	p.advance() // consume string

	return &ast.DialogStatement{
		Character: character,
		Colon:     colonToken,
		Text:      text,
	}
}

func (p *Parser) parseChoiceStatement() *ast.ChoiceStatement {
	choiceToken := p.peek()
	p.advance() // consume CHOICE

	if !p.check(token.LBRACE) {
		return nil
	}

	p.advance() // consume '{'

	var options []*ast.ChoiceOption

	for !p.check(token.RBRACE) && !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.advance()
			continue
		}

		option := p.parseChoiceOption()
		if option != nil {
			options = append(options, option)
		}

		// Consume optional comma
		if p.check(token.COMMA) {
			p.advance()
		}
	}

	if !p.check(token.RBRACE) {
		return nil
	}

	p.advance() // consume '}'

	return &ast.ChoiceStatement{
		Token:   choiceToken,
		Options: options,
	}
}

func (p *Parser) parseChoiceOption() *ast.ChoiceOption {
	if !p.check(token.STRING) {
		return nil
	}

	text := &ast.StringLiteral{
		Token: p.peek(),
		Value: p.peek().Literal.(string),
	}
	p.advance()

	if !p.check(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.ChoiceOption{
		Text: text,
		Body: body,
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	lbraceToken := p.peek()
	p.advance() // consume '{'

	var statements []ast.Statement

	for !p.check(token.RBRACE) && !p.isAtEnd() {
		// Skip newlines and comments
		if p.check(token.NEWLINE) || p.check(token.COMMENT) {
			p.advance()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	if !p.check(token.RBRACE) {
		return nil
	}

	p.advance() // consume '}'
	return &ast.BlockStatement{
		Token:      lbraceToken,
		Statements: statements,
	}
}
