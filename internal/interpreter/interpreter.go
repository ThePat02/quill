package interpreter

import (
	"fmt"
	"math/rand"
	"quill/internal/ast"
	"quill/internal/token"
)

type ResultType int

const (
	DialogResult ResultType = iota
	ChoiceResult
	ToolCallResult
	EndResult
	ErrorResult
)

type InterpreterResult struct {
	Type ResultType
	Data interface{}
}

type DialogData struct {
	Character string   `json:"character"`
	Text      string   `json:"text"`
	Tags      []string `json:"tags"`
}

type ChoiceOption struct {
	Index int      `json:"index"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type ChoiceData struct {
	Options []ChoiceOption `json:"options"`
}

type ToolCallData struct {
	Function  string        `json:"function"`
	Arguments []interface{} `json:"arguments"`
}

type ExecutionState int

const (
	StateReady ExecutionState = iota
	StateWaitingForChoice
	StateWaitingForToolCall
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
	variables         map[string]interface{}
	state             ExecutionState
	currentStatements []ast.Statement
	statementIndex    int
	pendingChoice     *ast.ChoiceStatement
	pendingToolCall   *ast.ToolCall
	pendingAssignment *ast.Identifier // Variable to assign tool call result to
	executionStack    []executionFrame
}

type InterpreterError struct {
	Message string
	Line    int
}

type ErrorData struct {
	Message string `json:"message"`
	Line    int    `json:"line"`
}

func New(program *ast.Program) *Interpreter {
	interpreter := &Interpreter{
		program:           program,
		labels:            make(map[string]*ast.LabelStatement),
		variables:         make(map[string]interface{}),
		state:             StateReady,
		currentStatements: program.Statements,
		statementIndex:    0,
		executionStack:    make([]executionFrame, 0),
	}

	interpreter.collectLabels()

	// Note: As of Go 1.20, rand.Seed is deprecated and no longer needed
	// The default source is automatically seeded with a random value

	return interpreter
}

func (i *Interpreter) executeStatement(stmt ast.Statement) *InterpreterResult {
	switch node := stmt.(type) {
	case *ast.LetStatement:
		return i.executeLetStatement(node)
	case *ast.AssignStatement:
		return i.executeAssignStatement(node)
	case *ast.IfStatement:
		return i.executeIfStatement(node)
	case *ast.LabelStatement:
		// Labels are just markers, don't do anything special
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

func (i *Interpreter) executeLetStatement(letStmt *ast.LetStatement) *InterpreterResult {
	// Check if the value is a tool call
	if toolCall, ok := letStmt.Value.(*ast.ToolCall); ok {
		// Evaluate arguments
		var args []interface{}
		for _, arg := range toolCall.Arguments {
			value, err := i.evaluateExpression(arg)
			if err != nil {
				return err
			}
			args = append(args, value)
		}

		// Store context for when we get the response
		i.pendingToolCall = toolCall
		i.pendingAssignment = letStmt.Name
		i.state = StateWaitingForToolCall
		// Don't decrement statement index - we'll complete this statement when we get the response

		return &InterpreterResult{
			Type: ToolCallResult,
			Data: ToolCallData{
				Function:  toolCall.Function,
				Arguments: args,
			},
		}
	}

	// Regular expression evaluation
	value, err := i.evaluateExpression(letStmt.Value)
	if err != nil {
		return err
	}

	i.variables[letStmt.Name.Value] = value
	return nil // Continue to next statement
}

func (i *Interpreter) executeAssignStatement(assignStmt *ast.AssignStatement) *InterpreterResult {
	currentValue, exists := i.variables[assignStmt.Name.Value]
	if !exists {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Variable '" + assignStmt.Name.Value + "' not defined",
				Line:    assignStmt.Name.Token.Line,
			},
		}
	}

	newValue, err := i.evaluateExpression(assignStmt.Value)
	if err != nil {
		return err
	}

	switch assignStmt.Operator.Type {
	case token.ASSIGN:
		i.variables[assignStmt.Name.Value] = newValue
	case token.PLUS_ASSIGN:
		if currentInt, ok := currentValue.(int64); ok {
			if newInt, ok := newValue.(int64); ok {
				i.variables[assignStmt.Name.Value] = currentInt + newInt
			} else {
				return &InterpreterResult{
					Type: ErrorResult,
					Data: ErrorData{
						Message: "Cannot add non-integer to integer",
						Line:    assignStmt.Operator.Line,
					},
				}
			}
		}
	case token.MINUS_ASSIGN:
		if currentInt, ok := currentValue.(int64); ok {
			if newInt, ok := newValue.(int64); ok {
				i.variables[assignStmt.Name.Value] = currentInt - newInt
			} else {
				return &InterpreterResult{
					Type: ErrorResult,
					Data: ErrorData{
						Message: "Cannot subtract non-integer from integer",
						Line:    assignStmt.Operator.Line,
					},
				}
			}
		}
	}

	return nil // Continue to next statement
}

func (i *Interpreter) executeIfStatement(ifStmt *ast.IfStatement) *InterpreterResult {
	condition, err := i.evaluateExpression(ifStmt.Condition)
	if err != nil {
		return err
	}

	conditionBool, ok := condition.(bool)
	if !ok {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "IF condition must be a boolean",
				Line:    ifStmt.Token.Line,
			},
		}
	}

	if conditionBool {
		return i.executeBlock(ifStmt.Consequence)
	} else if ifStmt.Alternative != nil {
		return i.executeBlock(ifStmt.Alternative)
	}

	return nil // Continue to next statement
}

func (i *Interpreter) executeDialog(dialog *ast.DialogStatement) *InterpreterResult {
	character := dialog.Character.Value

	// Evaluate the text expression (handles both StringLiteral and InterpolatedString)
	textResult, err := i.evaluateExpression(dialog.Text)
	if err != nil {
		return err
	}

	text := textResult.(string)

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

		// Handle both regular strings and interpolated strings
		if stringLit, ok := option.Text.(*ast.StringLiteral); ok {
			text = i.interpolateString(stringLit.Value)
		} else if interpolated, ok := option.Text.(*ast.InterpolatedString); ok {
			result, err := i.evaluateExpression(interpolated)
			if err != nil {
				return err
			}
			text = result.(string)
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

	if i.state == StateWaitingForToolCall {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Cannot step while waiting for tool call response",
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

func (i *Interpreter) HandleToolCallResponse(result interface{}) *InterpreterResult {
	if i.state != StateWaitingForToolCall || i.pendingToolCall == nil {
		return &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Not waiting for tool call response",
				Line:    0,
			},
		}
	}

	// Store the result in the pending assignment variable
	if i.pendingAssignment != nil {
		i.variables[i.pendingAssignment.Value] = result
	}

	// Clear pending state
	i.pendingToolCall = nil
	i.pendingAssignment = nil
	i.state = StateReady

	// Return nil to indicate the LET statement is now complete and execution should continue
	return nil
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

func (i *Interpreter) IsWaitingForToolCall() bool {
	return i.state == StateWaitingForToolCall
}

func (i *Interpreter) evaluateExpression(expr ast.Expression) (interface{}, *InterpreterResult) {
	switch node := expr.(type) {
	case *ast.Identifier:
		value, exists := i.variables[node.Value]
		if !exists {
			return nil, &InterpreterResult{
				Type: ErrorResult,
				Data: ErrorData{
					Message: "Variable '" + node.Value + "' not defined",
					Line:    node.Token.Line,
				},
			}
		}
		return value, nil

	case *ast.IntegerLiteral:
		return node.Value, nil

	case *ast.BooleanLiteral:
		return node.Value, nil

	case *ast.StringLiteral:
		return node.Value, nil

	case *ast.InterpolatedString:
		result := ""
		for _, part := range node.Parts {
			if ident, ok := part.(*ast.Identifier); ok {
				// Variable interpolation
				value, exists := i.variables[ident.Value]
				if !exists {
					return nil, &InterpreterResult{
						Type: ErrorResult,
						Data: ErrorData{
							Message: "Variable '" + ident.Value + "' not defined",
							Line:    ident.Token.Line,
						},
					}
				}
				result += i.valueToString(value)
			} else if str, ok := part.(*ast.StringLiteral); ok {
				// String literal part
				result += str.Value
			} else if toolCall, ok := part.(*ast.ToolCall); ok {
				// Tool call interpolation
				toolResult, err := i.executeToolCallImmediate(toolCall)
				if err != nil {
					return nil, err
				}
				result += i.valueToString(toolResult)
			}
		}
		return result, nil

	case *ast.InfixExpression:
		return i.evaluateInfixExpression(node)

	case *ast.PrefixExpression:
		return i.evaluatePrefixExpression(node)

	case *ast.ToolCall:
		// For tool calls in expressions (like interpolated strings), we need immediate execution
		// Only LET statements pause for tool calls
		return i.executeToolCallImmediate(node)

	default:
		return nil, &InterpreterResult{
			Type: ErrorResult,
			Data: ErrorData{
				Message: "Unknown expression type",
				Line:    0,
			},
		}
	}
}

func (i *Interpreter) evaluateInfixExpression(expr *ast.InfixExpression) (interface{}, *InterpreterResult) {
	// Special handling for null coalescing operator
	if expr.Operator == "??" {
		left, err := i.evaluateExpression(expr.Left)
		if err != nil {
			// If left side fails, evaluate right side
			right, err2 := i.evaluateExpression(expr.Right)
			if err2 != nil {
				return nil, err // Return the original left error
			}
			return right, nil
		}

		// Check if left value is "falsy" (nil, false, 0, empty string)
		if i.isFalsy(left) {
			right, err := i.evaluateExpression(expr.Right)
			if err != nil {
				return nil, err
			}
			return right, nil
		}

		return left, nil
	}

	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case "+":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt + rightInt, nil
			}
		}
	case "-":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt - rightInt, nil
			}
		}
	case ">":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt > rightInt, nil
			}
		}
	case "<":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt < rightInt, nil
			}
		}
	case ">=":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt >= rightInt, nil
			}
		}
	case "<=":
		if leftInt, ok := left.(int64); ok {
			if rightInt, ok := right.(int64); ok {
				return leftInt <= rightInt, nil
			}
		}
	case "&&":
		if leftBool, ok := left.(bool); ok {
			if rightBool, ok := right.(bool); ok {
				return leftBool && rightBool, nil
			}
		}
	case "||":
		if leftBool, ok := left.(bool); ok {
			if rightBool, ok := right.(bool); ok {
				return leftBool || rightBool, nil
			}
		}
	}

	return nil, &InterpreterResult{
		Type: ErrorResult,
		Data: ErrorData{
			Message: "Invalid operation: " + expr.Operator,
			Line:    expr.Token.Line,
		},
	}
}

func (i *Interpreter) evaluatePrefixExpression(expr *ast.PrefixExpression) (interface{}, *InterpreterResult) {
	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "!":
		if rightBool, ok := right.(bool); ok {
			return !rightBool, nil
		}
	}

	return nil, &InterpreterResult{
		Type: ErrorResult,
		Data: ErrorData{
			Message: "Invalid prefix operation: " + expr.Operator,
			Line:    expr.Token.Line,
		},
	}
}

func (i *Interpreter) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (interp *Interpreter) interpolateString(text string) string {
	result := ""
	for i := 0; i < len(text); i++ {
		if text[i] == '{' {
			// Find the closing brace
			j := i + 1
			for j < len(text) && text[j] != '}' {
				j++
			}

			if j < len(text) {
				// Extract variable name and substitute
				varName := text[i+1 : j]
				if value, exists := interp.variables[varName]; exists {
					result += interp.valueToString(value)
				} else {
					// Variable not found, keep the original text
					result += "{" + varName + "}"
				}
				i = j // Skip past the closing brace
			} else {
				// No closing brace found, add the character as-is
				result += string(text[i])
			}
		} else {
			result += string(text[i])
		}
	}
	return result
}

func (i *Interpreter) executeToolCallImmediate(toolCall *ast.ToolCall) (interface{}, *InterpreterResult) {
	// Evaluate all arguments first
	var args []interface{}
	for _, arg := range toolCall.Arguments {
		value, err := i.evaluateExpression(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}

	// Execute immediately with mock data
	result := i.callTool(toolCall.Function, args)
	return result, nil
}

func (i *Interpreter) callTool(functionName string, args []interface{}) interface{} {
	// Mock implementation of tool calls
	// In a real system, this would interface with external APIs, databases, etc.

	switch functionName {
	case "getPlayerName":
		// Return a default player name or from some external system
		return "Player"

	case "getPlayerAge":
		// Return a default age or from some external system
		return int64(25)

	case "getData":
		// Simulate getting data based on the first argument
		if len(args) > 0 {
			key := i.valueToString(args[0])
			switch key {
			case "gold":
				return int64(100)
			case "health":
				return int64(80)
			default:
				return "Unknown"
			}
		}
		return "No data"

	case "getItemPrice":
		// Simulate getting item price based on item type and level
		if len(args) >= 2 {
			itemType := i.valueToString(args[0])
			level, ok := args[1].(int64)
			if !ok {
				level = 1
			}

			// Simple price calculation: base price * level
			basePrice := int64(10)
			if itemType == "potion" {
				basePrice = 5
			} else if itemType == "weapon" {
				basePrice = 50
			} else if itemType == "armor" {
				basePrice = 30
			}

			return basePrice * level
		}
		return int64(0)

	case "agePlusFive":
		// Add 5 to the age argument
		if len(args) > 0 {
			if age, ok := args[0].(int64); ok {
				return age + 5
			}
		}
		return int64(5)

	default:
		// Unknown function, return a placeholder
		return "Unknown function: " + functionName
	}
}

func (i *Interpreter) isFalsy(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return true
	case bool:
		return !v
	case int64:
		return v == 0
	case string:
		return v == ""
	default:
		return false
	}
}
