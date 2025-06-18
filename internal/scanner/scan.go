package scanner

import "quill/internal/token"

func (scanner *Scanner) scanString() *ScannerError {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		return &ScannerError{
			Line:    scanner.line,
			Message: "Unterminated string.",
		}
	}

	scanner.advance()

	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addTokenWithLiteral(token.STRING, string(value))
	return nil
}

func (scanner *Scanner) scanNumber() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	value := scanner.source[scanner.start:scanner.current]
	scanner.addTokenWithLiteral(token.INT, string(value))
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
