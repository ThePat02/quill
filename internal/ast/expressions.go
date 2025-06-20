package ast

import "quill/internal/token"

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	if i == nil {
		return "<nil Identifier>"
	}
	return i.Value
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	if sl == nil {
		return "<nil StringLiteral>"
	}
	return "\"" + sl.Value + "\""
}

type TagList struct {
	Token token.Token // The '[' token
	Tags  []*Identifier
}

func (tl *TagList) expressionNode() {}
func (tl *TagList) String() string {
	if tl == nil {
		return "<nil TagList>"
	}
	if len(tl.Tags) == 0 {
		return "[]"
	}

	result := "["
	for i, tag := range tl.Tags {
		if i > 0 {
			result += ", "
		}
		if tag != nil {
			result += tag.String()
		}
	}
	result += "]"
	return result
}

// Boolean Literal
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	if bl == nil {
		return "<nil BooleanLiteral>"
	}
	return bl.Token.Lexeme
}

// Integer Literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string {
	if il == nil {
		return "<nil IntegerLiteral>"
	}
	return il.Token.Lexeme
}

// Infix Expression (for operators like +, -, ==, >=, etc.)
type InfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	if ie == nil {
		return "<nil InfixExpression>"
	}
	result := "("
	if ie.Left != nil {
		result += ie.Left.String()
	}
	result += " " + ie.Operator + " "
	if ie.Right != nil {
		result += ie.Right.String()
	}
	result += ")"
	return result
}

// Prefix Expression (for operators like !)
type PrefixExpression struct {
	Token    token.Token // the prefix token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	if pe == nil {
		return "<nil PrefixExpression>"
	}
	result := "(" + pe.Operator
	if pe.Right != nil {
		result += pe.Right.String()
	}
	result += ")"
	return result
}

// Variable Interpolation Expression (for {variable} in strings)
type InterpolatedString struct {
	Token token.Token
	Parts []Expression // mix of StringLiteral and Identifier
}

func (is *InterpolatedString) expressionNode() {}
func (is *InterpolatedString) String() string {
	if is == nil {
		return "<nil InterpolatedString>"
	}
	result := "\""
	for _, part := range is.Parts {
		if part == nil {
			continue
		}
		if ident, ok := part.(*Identifier); ok {
			result += "{" + ident.String() + "}"
		} else {
			result += part.String()
		}
	}
	result += "\""
	return result
}
