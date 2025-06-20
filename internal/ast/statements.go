package ast

import "quill/internal/token"

type LabelStatement struct {
	Token token.Token
	Name  *Identifier
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	if ls == nil {
		return "<nil LabelStatement>"
	}
	result := ls.Token.Lexeme
	if ls.Name != nil {
		result += " " + ls.Name.String()
	}
	return result
}

type GotoStatement struct {
	Token token.Token
	Label *Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	if gs == nil {
		return "<nil GotoStatement>"
	}
	result := gs.Token.Lexeme
	if gs.Label != nil {
		result += " " + gs.Label.String()
	}
	return result
}

type EndStatement struct {
	Token token.Token
}

func (es *EndStatement) statementNode() {}
func (es *EndStatement) String() string {
	if es == nil {
		return "<nil EndStatement>"
	}
	return es.Token.Lexeme
}

type DialogStatement struct {
	Character *Identifier
	Colon     token.Token
	Text      Expression
	Tags      *TagList
}

func (ds *DialogStatement) statementNode() {}
func (ds *DialogStatement) String() string {
	if ds == nil {
		return "<nil DialogStatement>"
	}
	result := ""
	if ds.Character != nil {
		result += ds.Character.String()
	}
	result += ": "
	if ds.Text != nil {
		result += ds.Text.String()
	}
	if ds.Tags != nil {
		result += " " + ds.Tags.String()
	}
	return result
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
	if co == nil {
		return "<nil ChoiceOption>"
	}
	result := ""
	if co.Text != nil {
		result += co.Text.String()
	}
	if co.Body != nil {
		result += " " + co.Body.String()
	}
	if co.Tags != nil {
		result += " " + co.Tags.String()
	}
	return result
}

func (cs *ChoiceStatement) statementNode() {}
func (cs *ChoiceStatement) String() string {
	if cs == nil {
		return "<nil ChoiceStatement>"
	}
	result := cs.Token.Lexeme + " {\n"
	for _, option := range cs.Options {
		if option != nil {
			result += "  " + option.String() + "\n"
		}
	}
	result += "}"
	return result
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	if bs == nil {
		return "<nil BlockStatement>"
	}
	result := "{\n"
	for _, stmt := range bs.Statements {
		if stmt != nil {
			result += "  " + stmt.String() + "\n"
		}
	}
	result += "}"
	return result
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
	if ro == nil {
		return "<nil RandomOption>"
	}
	result := ""
	if ro.Body != nil {
		result += ro.Body.String()
	}
	if ro.Tags != nil {
		result += " " + ro.Tags.String()
	}
	return result
}

func (rs *RandomStatement) statementNode() {}
func (rs *RandomStatement) String() string {
	if rs == nil {
		return "<nil RandomStatement>"
	}
	result := rs.Token.Lexeme + " {\n"
	for _, option := range rs.Options {
		if option != nil {
			result += "  " + option.String() + "\n"
		}
	}
	result += "}"
	return result
}

// Variable Declaration Statement
type LetStatement struct {
	Token token.Token // the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	if ls == nil {
		return "<nil LetStatement>"
	}
	result := ls.Token.Lexeme
	if ls.Name != nil {
		result += " " + ls.Name.String()
	}
	result += " = "
	if ls.Value != nil {
		result += ls.Value.String()
	}
	return result
}

// Assignment Statement
type AssignStatement struct {
	Name     *Identifier
	Operator token.Token // =, +=, -=
	Value    Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) String() string {
	if as == nil {
		return "<nil AssignStatement>"
	}
	result := ""
	if as.Name != nil {
		result += as.Name.String()
	}
	result += " " + as.Operator.Lexeme + " "
	if as.Value != nil {
		result += as.Value.String()
	}
	return result
}

// If Statement
type IfStatement struct {
	Token       token.Token // the IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // can be nil
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	if is == nil {
		return "<nil IfStatement>"
	}
	result := is.Token.Lexeme + " "
	if is.Condition != nil {
		result += is.Condition.String()
	}
	result += " "
	if is.Consequence != nil {
		result += is.Consequence.String()
	}
	if is.Alternative != nil {
		result += " ELSE " + is.Alternative.String()
	}
	return result
}
