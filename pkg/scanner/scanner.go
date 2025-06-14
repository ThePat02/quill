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
		scanner.addToken(token.SEMICOLON)
	case ',':
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
