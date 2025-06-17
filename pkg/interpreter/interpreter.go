package interpreter

import (
	"math/rand"
	"quill/pkg/ast"
	"time"
)

type ResultType int

const (
	DialogResult ResultType = iota
	ChoiceResult
	EndResult
	ErrorResult
)

type InterpreterResult struct {
	Type ResultType
	Data interface{}
}

type DialogData struct {
	Character string
	Text      string
	Tags      []string
}

type ChoiceOption struct {
	Index int
	Text  string
	Tags  []string
}

type ChoiceData struct {
	Options []ChoiceOption
}

type ExecutionState int

const (
	StateReady ExecutionState = iota
	StateWaitingForChoice
	StateEnded
	StateError
)

type executionFrame struct {
	statements []ast.Statement
	index      int
}

type Interpreter struct {
	program           *ast.Program
	labels            map[string]*ast.LabelStatement
	state             ExecutionState
	currentStatements []ast.Statement
	statementIndex    int
	pendingChoice     *ast.ChoiceStatement
	executionStack    []executionFrame
}

type InterpreterError struct {
	Message string
	Line    int
}

type ErrorData struct {
	Message string
	Line    int
}

func New(program *ast.Program) *Interpreter {
	interpreter := &Interpreter{
		program:           program,
		labels:            make(map[string]*ast.LabelStatement),
		state:             StateReady,
		currentStatements: program.Statements,
		statementIndex:    0,
		executionStack:    make([]executionFrame, 0),
	}

	interpreter.collectLabels()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	return interpreter
}

func (i *Interpreter) executeStatement(stmt ast.Statement) *InterpreterResult {
	switch node := stmt.(type) {
	case *ast.LabelStatement:
		// Labels are just markers, don't do anything special
		// Let the normal flow continue to the next statement
		return nil

	case *ast.DialogStatement:
		return i.executeDialog(node)

	case *ast.ChoiceStatement:
		return i.executeChoice(node)

	case *ast.RandomStatement:
		return i.executeRandom(node)

	case *ast.GotoStatement:
		return i.executeGoto(node)

	case *ast.EndStatement:
		i.state = StateEnded
		return &InterpreterResult{
			Type: EndResult,
			Data: nil,
		}

	case *ast.BlockStatement:
		return i.executeBlock(node)

	default:
		i.state = StateError
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "unknown statement type",
				Line:    i.getStatementLine(stmt),
			},
		}
	}
}

func (i *Interpreter) executeDialog(dialog *ast.DialogStatement) *InterpreterResult {
	character := dialog.Character.Value
	text := dialog.Text.Value

	// Extract tags if present
	var tags []string
	if dialog.Tags != nil && len(dialog.Tags.Tags) > 0 {
		tags = make([]string, len(dialog.Tags.Tags))
		for idx, tag := range dialog.Tags.Tags {
			tags[idx] = tag.Value
		}
	}

	return &InterpreterResult{
		Type: DialogResult,
		Data: DialogData{
			Character: character,
			Text:      text,
			Tags:      tags,
		},
	}
}

func (i *Interpreter) executeChoice(choice *ast.ChoiceStatement) *InterpreterResult {
	options := make([]ChoiceOption, len(choice.Options))

	for idx, option := range choice.Options {
		text := ""
		if stringLit, ok := option.Text.(*ast.StringLiteral); ok {
			text = stringLit.Value
		}

		var tags []string
		if option.Tags != nil && len(option.Tags.Tags) > 0 {
			tags = make([]string, len(option.Tags.Tags))
			for i, tag := range option.Tags.Tags {
				tags[i] = tag.Value
			}
		}

		options[idx] = ChoiceOption{
			Index: idx,
			Text:  text,
			Tags:  tags,
		}
	}

	// Store choice and wait for input
	i.pendingChoice = choice
	i.state = StateWaitingForChoice

	return &InterpreterResult{
		Type: ChoiceResult,
		Data: ChoiceData{
			Options: options,
		},
	}
}

func (i *Interpreter) executeRandom(random *ast.RandomStatement) *InterpreterResult {
	if len(random.Options) == 0 {
		i.state = StateError
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "RANDOM block has no options",
				Line:    random.Token.Line,
			},
		}
	}

	// Pick a random option
	selectedIndex := rand.Intn(len(random.Options))
	selectedOption := random.Options[selectedIndex]

	// Execute the selected option's body
	return i.executeBlock(selectedOption.Body)
}

func (i *Interpreter) executeBlock(block *ast.BlockStatement) *InterpreterResult {
	if len(block.Statements) == 0 {
		return i.Step()
	}

	// Push current context to stack
	i.executionStack = append(i.executionStack, executionFrame{
		statements: i.currentStatements,
		index:      i.statementIndex,
	})

	// Set up new execution context
	i.currentStatements = block.Statements
	i.statementIndex = 0

	return i.Step()
}

func (i *Interpreter) executeGoto(gotoStmt *ast.GotoStatement) *InterpreterResult {
	labelName := gotoStmt.Label.Value
	label, exists := i.labels[labelName]
	if !exists {
		i.state = StateError
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "label '" + labelName + "' not found",
				Line:    gotoStmt.Token.Line,
			},
		}
	}

	// Find the index of this label in the program
	labelIndex := -1
	for idx, stmt := range i.program.Statements {
		if stmt == label {
			labelIndex = idx
			break
		}
	}

	if labelIndex == -1 {
		i.state = StateError
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "label '" + labelName + "' not found in program",
				Line:    label.Token.Line,
			},
		}
	}

	// Clear execution stack and jump to label
	i.executionStack = make([]executionFrame, 0)
	i.currentStatements = i.program.Statements
	i.statementIndex = labelIndex // Don't add 1 here, Step() will increment it

	return i.Step()
}

func (i *Interpreter) collectLabels() {
	for _, stmt := range i.program.Statements {
		i.collectLabelsFromStatement(stmt)
	}
}

func (i *Interpreter) collectLabelsFromStatement(stmt ast.Statement) {
	switch node := stmt.(type) {
	case *ast.LabelStatement:
		i.labels[node.Name.Value] = node
	case *ast.ChoiceStatement:
		for _, option := range node.Options {
			i.collectLabelsFromBlock(option.Body)
		}
	case *ast.RandomStatement:
		for _, option := range node.Options {
			i.collectLabelsFromBlock(option.Body)
		}
	case *ast.BlockStatement:
		i.collectLabelsFromBlock(node)
	}
}

func (i *Interpreter) collectLabelsFromBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		i.collectLabelsFromStatement(stmt)
	}
}

// Helper function to get line number from statement
func (i *Interpreter) getStatementLine(stmt ast.Statement) int {
	switch node := stmt.(type) {
	case *ast.LabelStatement:
		return node.Token.Line
	case *ast.DialogStatement:
		return node.Character.Token.Line
	case *ast.ChoiceStatement:
		return node.Token.Line
	case *ast.RandomStatement:
		return node.Token.Line
	case *ast.GotoStatement:
		return node.Token.Line
	case *ast.EndStatement:
		return node.Token.Line
	case *ast.BlockStatement:
		return node.Token.Line
	default:
		return 0
	}
}

func (i *Interpreter) Step() *InterpreterResult {
	if i.state == StateEnded {
		return &InterpreterResult{
			Type: EndResult,
			Data: nil,
		}
	}

	if i.state == StateError {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Interpreter in error state",
				Line:    0,
			},
		}
	}

	if i.state == StateWaitingForChoice {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Cannot step while waiting for choice input",
				Line:    0,
			},
		}
	}

	// Execute next statement
	if i.statementIndex >= len(i.currentStatements) {
		// No more statements, check if we can pop from stack
		if len(i.executionStack) > 0 {
			frame := i.executionStack[len(i.executionStack)-1]
			i.executionStack = i.executionStack[:len(i.executionStack)-1]
			i.currentStatements = frame.statements
			i.statementIndex = frame.index
			return i.Step()
		} else {
			// Program completed
			i.state = StateEnded
			return &InterpreterResult{
				Type: EndResult,
				Data: nil,
			}
		}
	}

	stmt := i.currentStatements[i.statementIndex]
	i.statementIndex++

	result := i.executeStatement(stmt)

	// If executeStatement returns nil (like for labels), continue to next statement
	if result == nil {
		return i.Step()
	}

	return result
}

func (i *Interpreter) HandleChoiceInput(choiceIndex int) *InterpreterResult {
	if i.state != StateWaitingForChoice || i.pendingChoice == nil {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Not waiting for choice input",
				Line:    0,
			},
		}
	}

	if choiceIndex < 0 || choiceIndex >= len(i.pendingChoice.Options) {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Invalid choice index",
				Line:    i.pendingChoice.Token.Line,
			},
		}
	}

	// Execute the selected choice's body
	selectedOption := i.pendingChoice.Options[choiceIndex]
	i.pendingChoice = nil
	i.state = StateReady

	// Push current execution context and execute choice body
	return i.executeBlock(selectedOption.Body)
}

// Helper methods for external use
func (i *Interpreter) GetState() ExecutionState {
	return i.state
}

func (i *Interpreter) IsEnded() bool {
	return i.state == StateEnded
}

func (i *Interpreter) IsWaitingForChoice() bool {
	return i.state == StateWaitingForChoice
}
