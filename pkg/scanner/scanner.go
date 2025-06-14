package scanner

import "skribent/pkg/token"

type ErrorReporter func(line int, message string)

type Scanner struct {
	source        string
	tokens        []token.Token
	start         int
	current       int
	line          int
	errorReporter ErrorReporter
}

func New(source string, errorReporter ErrorReporter) *Scanner {
	return &Scanner{
		source:        source,
		tokens:        make([]token.Token, 0),
		start:         0,
		current:       0,
		line:          1,
		errorReporter: errorReporter,
	}
}

func (scanner *Scanner) ScanTokens() []token.Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, token.NewToken(token.EOF, "", nil, scanner.line))

	return scanner.tokens
}

func (scanner *Scanner) scanToken() {
	char := scanner.advance()
	switch char {
	// Ignored Characters
	case ' ', '\r', '\t':
		break

	// Special Cases
	case '\n':
		scanner.line++
		scanner.addToken(token.NEWLINE)

	// Single Character
	case ':':
		scanner.addToken(token.COLON)
	case ';':
		scanner.addToken(token.COMMA)
	case '(':
		scanner.addToken(token.LPAREN)
	case ')':
		scanner.addToken(token.RPAREN)
	case '{':
		scanner.addToken(token.LBRACE)
	case '}':
		scanner.addToken(token.RBRACE)
	case '[':
		scanner.addToken(token.LBRACKET)
	case ']':
		scanner.addToken(token.RBRACKET)
	case '=':
		scanner.addToken(token.ASSIGN)
	case '+':
		scanner.addToken(token.PLUS)
	case '-':
		scanner.addToken(token.MINUS)
	case '*':
		scanner.addToken(token.STAR)
	case '>':
		scanner.addToken(token.GT)
	case '<':
		scanner.addToken(token.LT)
	case '?':
		scanner.addToken(token.QUESTION)
	case '!':
		scanner.addToken(token.EXLAM)

	// Multi-Character
	case '#':
		for !scanner.isAtEnd() && scanner.peek() != '\n' {
			scanner.advance()
		}
		scanner.addToken(token.COMMENT)

	// Default
	default:
		scanner.errorReporter(scanner.line, "Unexpected character: "+string(char))
	}
}

func (scanner *Scanner) advance() byte {
	scanner.current++
	return scanner.source[scanner.current-1]
}

func (scanner *Scanner) addToken(tokenType token.TokenType) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, token.NewToken(tokenType, string(text), literal, scanner.line))
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() || scanner.source[scanner.current] != expected {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current]
}
