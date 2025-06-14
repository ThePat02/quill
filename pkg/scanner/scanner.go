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

	// Literals
	case '"':
		scanner.scanString()

	// Multi-Character
	case '#':
		for !scanner.isAtEnd() && scanner.peek() != '\n' {
			scanner.advance()
		}
		scanner.addToken(token.COMMENT)
	case '-':
		if scanner.peek() == '>' {
			scanner.advance()
			scanner.addToken(token.ARROW)
			return
		}
		scanner.addToken(token.MINUS)

	// Default
	default:
		// Digits
		if scanner.isDigit(char) {
			scanner.scanNumber()
			return
		}

		// Identifiers and Keywords
		if scanner.isAlpha(char) {
			scanner.scanIdentifier()
			return
		}

		// Handle unexpected characters
		scanner.addToken(token.ILLEGAL)
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

func (scanner *Scanner) scanString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		scanner.errorReporter(scanner.line, "Unterminated string.")
		return
	}

	scanner.advance()

	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addTokenWithLiteral(token.STRING, string(value))
}

func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}
	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) scanNumber() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	value := scanner.source[scanner.start:scanner.current]
	scanner.addTokenWithLiteral(token.INT, string(value))
}

func (scanner *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (scanner *Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}

func (scanner *Scanner) scanIdentifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.source[scanner.start:scanner.current]
	tokenType, exists := token.Keywords[string(text)]

	if exists {
		scanner.addToken(tokenType)
	} else {
		scanner.addToken(token.IDENT)
	}
}
