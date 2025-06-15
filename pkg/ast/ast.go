package ast

import "quill/pkg/token"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, stmt := range p.Statements {
		out += stmt.String() + "\n"
	}
	return out
}

type LabelStatement struct {
	Token token.Token
	Name  *Identifier
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return ls.Token.String() + " " + ls.Name.String()
}

type GotoStatement struct {
	Token token.Token
	Label *Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return gs.Token.String() + " " + gs.Label.String()
}

type EndStatement struct {
	Token token.Token
}

func (es *EndStatement) statementNode() {}
func (es *EndStatement) String() string {
	return es.Token.String()
}

type DialogStatement struct {
	Character *Identifier
	Colon     token.Token
	Text      *StringLiteral
}

func (ds *DialogStatement) statementNode() {}
func (ds *DialogStatement) String() string {
	return ds.Character.String() + ds.Colon.String() + " " + ds.Text.String()
}

type ChoiceStatement struct {
	Token   token.Token
	Options []*ChoiceOption
}

type ChoiceOption struct {
	Text Expression
	Body *BlockStatement
}

func (co *ChoiceOption) String() string {
	return "Choice: " + co.Text.String() + "\n" + co.Body.String()
}

func (cs *ChoiceStatement) statementNode() {}
func (cs *ChoiceStatement) String() string {
	var optionsStr string
	for _, option := range cs.Options {
		optionsStr += option.String() + "\n"
	}
	return cs.Token.String() + "\n" + optionsStr
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out string
	for _, stmt := range bs.Statements {
		out += stmt.String() + "\n"
	}
	return bs.Token.String() + "\n" + out
}

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
