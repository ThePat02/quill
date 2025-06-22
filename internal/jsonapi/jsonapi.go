package jsonapi

import (
	"encoding/json"
	"quill/internal/interpreter"
	"quill/internal/parser"
	"quill/internal/scanner"
)

// JSONResult is the main wrapper for all API responses
type JSONResult struct {
	Success bool   `json:"success"`
	Type    string `json:"type,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// QuillInterpreter wraps the Go interpreter with JSON API
type QuillInterpreter struct {
	interpreter *interpreter.Interpreter
}

// NewQuillInterpreter creates a new JSON API interpreter from source code
func NewQuillInterpreter(source string) (*QuillInterpreter, string) {
	// Scan tokens
	scanner := scanner.New(source)
	tokens, scannerErrors := scanner.ScanTokens()

	if len(scannerErrors) > 0 {
		result := JSONResult{
			Success: false,
			Type:    "scanner_errors",
			Data:    scannerErrors,
			Error:   "Scanner errors occurred",
		}

		jsonBytes, _ := json.Marshal(result)
		return nil, string(jsonBytes)
	}

	// Parse program
	parser := parser.New(tokens)
	program, parserErrors := parser.Parse()

	if len(parserErrors) > 0 {
		result := JSONResult{
			Success: false,
			Type:    "parser_errors",
			Data:    parserErrors,
			Error:   "Parser errors occurred",
		}

		jsonBytes, _ := json.Marshal(result)
		return nil, string(jsonBytes)
	}

	// Create interpreter
	interp := interpreter.New(program)

	result := JSONResult{
		Success: true,
		Type:    "ready",
		Data:    nil,
	}

	jsonBytes, _ := json.Marshal(result)
	return &QuillInterpreter{interpreter: interp}, string(jsonBytes)
}

// Step executes the next step in the interpreter and returns JSON
func (qi *QuillInterpreter) Step() string {
	if qi.interpreter == nil {
		result := JSONResult{
			Success: false,
			Error:   "Interpreter not initialized",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	interpResult := qi.interpreter.Step()
	return qi.convertResultToJSON(interpResult)
}

// HandleChoiceInput handles choice input and returns JSON
func (qi *QuillInterpreter) HandleChoiceInput(choiceIndex int) string {
	if qi.interpreter == nil {
		result := JSONResult{
			Success: false,
			Error:   "Interpreter not initialized",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	interpResult := qi.interpreter.HandleChoiceInput(choiceIndex)
	return qi.convertResultToJSON(interpResult)
}

// GetState returns the current interpreter state as JSON
func (qi *QuillInterpreter) GetState() string {
	if qi.interpreter == nil {
		result := JSONResult{
			Success: false,
			Error:   "Interpreter not initialized",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	state := qi.interpreter.GetState()
	var stateStr string

	switch state {
	case interpreter.StateReady:
		stateStr = "ready"
	case interpreter.StateWaitingForChoice:
		stateStr = "waiting_for_choice"
	case interpreter.StateEnded:
		stateStr = "ended"
	case interpreter.StateError:
		stateStr = "error"
	default:
		stateStr = "unknown"
	}

	result := JSONResult{
		Success: true,
		Type:    "state",
		Data:    stateStr,
	}

	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}

// IsEnded returns whether the interpreter has ended as JSON
func (qi *QuillInterpreter) IsEnded() string {
	if qi.interpreter == nil {
		result := JSONResult{
			Success: false,
			Error:   "Interpreter not initialized",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	result := JSONResult{
		Success: true,
		Type:    "ended_status",
		Data:    qi.interpreter.IsEnded(),
	}

	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}

// IsWaitingForChoice returns whether the interpreter is waiting for choice input as JSON
func (qi *QuillInterpreter) IsWaitingForChoice() string {
	if qi.interpreter == nil {
		result := JSONResult{
			Success: false,
			Error:   "Interpreter not initialized",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	result := JSONResult{
		Success: true,
		Type:    "waiting_for_choice_status",
		Data:    qi.interpreter.IsWaitingForChoice(),
	}

	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}

// convertResultToJSON converts interpreter results to JSON format
func (qi *QuillInterpreter) convertResultToJSON(interpResult *interpreter.InterpreterResult) string {
	if interpResult == nil {
		result := JSONResult{
			Success: false,
			Error:   "Received nil result from interpreter",
		}
		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	var resultType string
	var success bool = true

	switch interpResult.Type {
	case interpreter.DialogResult:
		resultType = "dialog"
	case interpreter.ChoiceResult:
		resultType = "choice"
	case interpreter.EndResult:
		resultType = "end"
	case interpreter.ErrorResult:
		resultType = "runtime_error"
		success = false
	default:
		resultType = "unknown"
		success = false
	}

	result := JSONResult{
		Success: success,
		Type:    resultType,
		Data:    interpResult.Data,
	}

	if !success && interpResult.Data != nil {
		if errorData, ok := interpResult.Data.(interpreter.ErrorData); ok {
			result.Error = errorData.Message
		}
	}

	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}

// ParseOnly parses source code without creating an interpreter, returns JSON
func ParseOnly(source string) string {
	// Scan tokens
	scanner := scanner.New(source)
	tokens, scannerErrors := scanner.ScanTokens()

	if len(scannerErrors) > 0 {
		result := JSONResult{
			Success: false,
			Type:    "scanner_errors",
			Data:    scannerErrors,
			Error:   "Scanner errors occurred",
		}

		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	// Parse program
	parser := parser.New(tokens)
	program, parserErrors := parser.Parse()

	if len(parserErrors) > 0 {
		result := JSONResult{
			Success: false,
			Type:    "parser_errors",
			Data:    parserErrors,
			Error:   "Parser errors occurred",
		}

		jsonBytes, _ := json.Marshal(result)
		return string(jsonBytes)
	}

	// Return success with program info
	programInfo := map[string]interface{}{
		"statement_count": len(program.Statements),
		"parsed":          true,
	}

	result := JSONResult{
		Success: true,
		Type:    "parse_success",
		Data:    programInfo,
	}

	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}
