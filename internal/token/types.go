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
	COLON         TokenType = ":"
	COMMA         TokenType = ","
	SEMICOLON     TokenType = ";"
	LPAREN        TokenType = "("
	RPAREN        TokenType = ")"
	LBRACE        TokenType = "{"
	RBRACE        TokenType = "}"
	LBRACKET      TokenType = "["
	RBRACKET      TokenType = "]"
	ASSIGN        TokenType = "="
	PLUS          TokenType = "+"
	PLUS_ASSIGN   TokenType = "+="
	MINUS         TokenType = "-"
	MINUS_ASSIGN  TokenType = "-="
	STAR          TokenType = "*"
	ARROW         TokenType = "->"
	QUESTION      TokenType = "?"
	EQ            TokenType = "=="
	NE            TokenType = "!="
	GT            TokenType = ">"
	LT            TokenType = "<"
	GE            TokenType = ">="
	LE            TokenType = "<="
	AND           TokenType = "&&"
	OR            TokenType = "||"
	NOT           TokenType = "!"
	NULL_COALESCE TokenType = "??"        // Null coalescing operator
	TOOL_CALL     TokenType = "TOOL_CALL" // Tool call token for <function; args>
)

// Keywords
const (
	SCENE  TokenType = "SCENE"  // Unused keyword
	RANDOM TokenType = "RANDOM" // Random keyword, used to indicate a random choice in the script
	GOTO   TokenType = "GOTO"   // Goto keyword, used for jumping to a label in the script
	LABEL  TokenType = "LABEL"  // Label keyword, used to define a label for goto statements
	CHOICE TokenType = "CHOICE" // Choice keyword, used to define a choice in the script
	END    TokenType = "END"    // End keyword, used to indicate the end of a script or block

	// Variable and logic keywords
	LET   TokenType = "LET"   // Let keyword, used to define a variable
	IF    TokenType = "IF"    // If keyword, used for conditional statements
	ELSE  TokenType = "ELSE"  // Else keyword, used for alternative paths in conditional statements
	TRUE  TokenType = "TRUE"  // True keyword, used for boolean true values
	FALSE TokenType = "FALSE" // False keyword, used for boolean false values
)

var Keywords = map[string]TokenType{
	"SCENE":  SCENE,
	"RANDOM": RANDOM,
	"GOTO":   GOTO,
	"LABEL":  LABEL,
	"CHOICE": CHOICE,
	"END":    END,
	"LET":    LET,
	"IF":     IF,
	"ELSE":   ELSE,
	"TRUE":   TRUE,
	"FALSE":  FALSE,
}
