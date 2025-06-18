package ast

import "quill/internal/token"

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return "\"" + sl.Value + "\""
}

type TagList struct {
	Token token.Token // The '[' token
	Tags  []*Identifier
}

func (tl *TagList) expressionNode() {}
func (tl *TagList) String() string {
	if len(tl.Tags) == 0 {
		return "[]"
	}

	result := "["
	for i, tag := range tl.Tags {
		if i > 0 {
			result += ", "
		}
		result += tag.String()
	}
	result += "]"
	return result
}
