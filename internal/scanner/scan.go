package scanner

import "quill/internal/token"

func (scanner *Scanner) scanString() *ScannerError {
	toolCallDepth := 0

	for scanner.peek() != '"' || toolCallDepth > 0 {
		if scanner.isAtEnd() {
			return &ScannerError{
				Line:    scanner.line,
				Message: "Unterminated string.",
			}
		}

		char := scanner.peek()

		// Handle tool call nesting
		if char == '<' {
			toolCallDepth++
		} else if char == '>' && toolCallDepth > 0 {
			toolCallDepth--
		} else if char == '"' && toolCallDepth > 0 {
			// This is a quote inside a tool call, just advance past it
			scanner.advance()
			continue
		}

		if char == '\n' {
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

	scanner.advance() // consume closing quote

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

func (scanner *Scanner) scanToolCall() *ScannerError {
	// We're already past the '<' character
	start := scanner.current - 1 // Include the '<' in the token

	// Scan until we find the closing '>'
	for !scanner.isAtEnd() && scanner.peek() != '>' {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		return &ScannerError{
			Line:    scanner.line,
			Message: "Unterminated tool call.",
		}
	}

	scanner.advance() // consume the '>'

	// Extract the full tool call content including < and >
	value := scanner.source[start:scanner.current]
	scanner.addTokenWithLiteral(token.TOOL_CALL, string(value))
	return nil
}
