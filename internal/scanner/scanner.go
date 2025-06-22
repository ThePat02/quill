package scanner

import "quill/internal/token"

type ErrorReporter func(line int, message string)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

type ScannerError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (scanner *Scanner) ScanTokens() ([]token.Token, []ScannerError) {
	var errors []ScannerError = make([]ScannerError, 0)
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		err := scanner.scanToken()
		if err != nil {
			errors = append(errors, *err)
		}
	}

	scanner.tokens = append(scanner.tokens, token.NewToken(token.EOF, "", nil, scanner.line))

	return scanner.tokens, errors
}

func (scanner *Scanner) scanToken() *ScannerError {
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
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.EQ)
		} else {
			scanner.addToken(token.ASSIGN)
		}
	case '+':
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.PLUS_ASSIGN)
		} else {
			scanner.addToken(token.PLUS)
		}
	case '-':
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.MINUS_ASSIGN)
		} else {
			scanner.addToken(token.MINUS)
		}
	case '>':
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.GE)
		} else {
			scanner.addToken(token.GT)
		}
	case '<':
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.LE)
		} else {
			scanner.addToken(token.LT)
		}
	case '!':
		if scanner.peek() == '=' {
			scanner.advance()
			scanner.addToken(token.NE)
		} else {
			scanner.addToken(token.NOT)
		}
	case '&':
		if scanner.peek() == '&' {
			scanner.advance()
			scanner.addToken(token.AND)
		}
	case '|':
		if scanner.peek() == '|' {
			scanner.advance()
			scanner.addToken(token.OR)
		}

	// Literals
	case '"':
		err := scanner.scanString()
		if err != nil {
			return err
		}

	// Multi-Character
	case '#':
		for !scanner.isAtEnd() && scanner.peek() != '\n' {
			scanner.advance()
		}
		scanner.addToken(token.COMMENT)

	// Default
	default:
		// Digits
		if scanner.isDigit(char) {
			scanner.scanNumber()
			return nil
		}

		// Identifiers and Keywords
		if scanner.isAlpha(char) {
			scanner.scanIdentifier()
			return nil
		}

		// Handle unexpected characters
		scanner.addToken(token.ILLEGAL)
		return &ScannerError{
			Line:    scanner.line,
			Message: "Unexpected character: " + string(char),
		}
	}

	return nil
}
