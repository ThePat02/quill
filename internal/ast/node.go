package ast

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
	if p == nil {
		return "<nil Program>"
	}
	var out string
	for _, stmt := range p.Statements {
		if stmt != nil {
			out += stmt.String() + "\n"
		}
	}
	return out
}
