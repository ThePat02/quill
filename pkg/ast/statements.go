package ast

import "quill/pkg/token"

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
	Tags      *TagList
}

func (ds *DialogStatement) statementNode() {}
func (ds *DialogStatement) String() string {
	return ds.Character.String() + ds.Colon.String() + " " + ds.Text.String() + " | " + ds.Tags.String()
}

type ChoiceStatement struct {
	Token   token.Token
	Options []*ChoiceOption
}

type ChoiceOption struct {
	Text Expression
	Body *BlockStatement
	Tags *TagList
}

func (co *ChoiceOption) String() string {
	result := "Choice: " + co.Text.String()
	if co.Tags != nil {
		result += " " + co.Tags.String()
	}
	result += "\n" + co.Body.String()
	return result
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

type RandomStatement struct {
	Token   token.Token
	Options []*RandomOption
}

type RandomOption struct {
	Body *BlockStatement
	Tags *TagList
}

func (ro *RandomOption) String() string {
	result := "Random Option:"
	if ro.Tags != nil {
		result += " " + ro.Tags.String()
	}
	result += "\n" + ro.Body.String()
	return result
}

func (rs *RandomStatement) statementNode() {}
func (rs *RandomStatement) String() string {
	var optionsStr string
	for _, option := range rs.Options {
		optionsStr += option.String() + "\n"
	}
	return rs.Token.String() + "\n" + optionsStr
}
