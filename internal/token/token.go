package token

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	switch t.Lexeme {
	case "\n":
		return "\\n (" + string(t.Type) + ")"
	default:
		return t.Lexeme + " (" + string(t.Type) + ")"
	}
}
