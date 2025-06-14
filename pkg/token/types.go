package token

type TokenType string

// Special Tokens
const (
	ILLEGAL TokenType = "ILLEGAL" // Illegal token, used for unrecognized characters
	NEWLINE TokenType = "NEWLINE" // Newline token, used to indicate line breaks in the source code
	EOF     TokenType = "EOF"     // End of file token, indicates no more input
)

// Generic Tokens
const (
	IDENT   TokenType = "IDENT"   // Identifiers, used for variable names, function names, etc.
	INT     TokenType = "INT"     // Integer literals
	STRING  TokenType = "STRING"  // String literals
	COMMENT TokenType = "COMMENT" // Comments in the source code

	// Structural
	COLON    TokenType = ":"
	COMMA    TokenType = ";"
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	STAR     TokenType = "*"
	ARROW    TokenType = "->"
	QUESTION TokenType = "?"
	EXLAM    TokenType = "!"
	GT       TokenType = ">"
	LT       TokenType = "<"
)

// Keywords
const (
	SCENE TokenType = "SCENE" // Scene keyword, used to define a scene in the script
	GOTO  TokenType = "GOTO"  // Goto keyword, used for jumping to a label in the script
	LABEL TokenType = "LABEL" // Label keyword, used to define a label for goto statements
)

var Keywords = map[string]TokenType{
	"scene": SCENE,
	"goto":  GOTO,
	"label": LABEL,
}
