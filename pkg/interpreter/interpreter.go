package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"quill/pkg/ast"
	"strconv"
	"strings"
)

type Interpreter struct {
	program *ast.Program
	labels  map[string]*ast.LabelStatement // Fix: should be pointer type
	scanner *bufio.Scanner
}

func New(program *ast.Program) *Interpreter {
	interpreter := &Interpreter{
		program: program,
		labels:  make(map[string]*ast.LabelStatement), // Fix: pointer type
		scanner: bufio.NewScanner(os.Stdin),
	}

	interpreter.collectLabels()

	return interpreter
}

func (i *Interpreter) Interpret() error {
	fmt.Println("--- Starting Dialog ---")
	return i.executeStatements(i.program.Statements)
}

func (i *Interpreter) executeStatements(statements []ast.Statement) error {
	for _, stmt := range statements {
		result, err := i.executeStatement(stmt)
		if err != nil {
			return err
		}

		// Handle control flow
		switch result {
		case "END":
			fmt.Println("--- Dialog Ended ---")
			return nil
		case "":
			// Continue to next statement
			continue
		default:
			// Handle GOTO
			if strings.HasPrefix(result, "GOTO:") {
				labelName := strings.TrimPrefix(result, "GOTO:")
				return i.gotoLabel(labelName)
			}
		}
	}

	fmt.Println("--- Dialog Completed ---")
	return nil
}

func (i *Interpreter) executeStatement(stmt ast.Statement) (string, error) {
	switch node := stmt.(type) {
	case *ast.LabelStatement:
		// Labels are just markers, no execution needed
		return "", nil

	case *ast.DialogStatement:
		return i.executeDialog(node)

	case *ast.ChoiceStatement:
		return i.executeChoice(node)

	case *ast.GotoStatement:
		return "GOTO:" + node.Label.Value, nil

	case *ast.EndStatement:
		return "END", nil

	case *ast.BlockStatement:
		return i.executeBlock(node)

	default:
		return "", fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeDialog(dialog *ast.DialogStatement) (string, error) {
	character := dialog.Character.Value
	text := dialog.Text.Value

	fmt.Printf("%s: %s\n", character, text)
	return "", nil
}

func (i *Interpreter) executeChoice(choice *ast.ChoiceStatement) (string, error) {
	fmt.Println("\nChoices:")
	for idx, option := range choice.Options {
		text := ""
		if stringLit, ok := option.Text.(*ast.StringLiteral); ok {
			text = stringLit.Value
		}
		fmt.Printf("%d. %s\n", idx+1, text)
	}

	fmt.Print("Enter your choice (number): ")
	if !i.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(i.scanner.Text())
	choiceNum, err := strconv.Atoi(input)
	if err != nil || choiceNum < 1 || choiceNum > len(choice.Options) {
		fmt.Println("Invalid choice. Please try again.")
		return i.executeChoice(choice) // Retry
	}

	selectedOption := choice.Options[choiceNum-1]
	fmt.Printf("You chose: %s\n", selectedOption.Text.(*ast.StringLiteral).Value)

	// Execute the body of the selected choice
	return i.executeBlock(selectedOption.Body)
}

func (i *Interpreter) executeBlock(block *ast.BlockStatement) (string, error) {
	for _, stmt := range block.Statements {
		result, err := i.executeStatement(stmt)
		if err != nil {
			return "", err
		}

		// If we hit a control flow statement, return it
		if result != "" {
			return result, nil
		}
	}

	return "", nil
}

func (i *Interpreter) gotoLabel(labelName string) error {
	label, exists := i.labels[labelName]
	if !exists {
		return fmt.Errorf("label '%s' not found", labelName)
	}

	// Find the index of this label in the program
	labelIndex := -1
	for idx, stmt := range i.program.Statements { // Fix: rename loop variable to avoid conflict
		if stmt == label {
			labelIndex = idx
			break
		}
	}

	if labelIndex == -1 {
		return fmt.Errorf("label '%s' not found in program", labelName)
	}

	// Execute from the statement after the label
	remainingStatements := i.program.Statements[labelIndex+1:]
	return i.executeStatements(remainingStatements)
}

func (i *Interpreter) collectLabels() {
	for _, stmt := range i.program.Statements {
		i.collectLabelsFromStatement(stmt)
	}
}

func (i *Interpreter) collectLabelsFromStatement(stmt ast.Statement) {
	switch node := stmt.(type) {
	case *ast.LabelStatement:
		i.labels[node.Name.Value] = node // Fix: store pointer, not value
	case *ast.ChoiceStatement:
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
