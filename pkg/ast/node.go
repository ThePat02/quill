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
	var out string
	for _, stmt := range p.Statements {
		out += stmt.String() + "\n"
	}
	return out
}
