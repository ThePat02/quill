package parser

import "quill/pkg/token"

func (p *Parser) peek() token.Token {
	if p.current >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() token.Token {
	if p.current < len(p.tokens) {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().Type == token.EOF
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) checkNext(tokenType token.TokenType) bool {
	if p.current+1 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.current+1].Type == tokenType
}
