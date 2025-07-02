package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"quill/internal/interpreter"
	"quill/internal/parser"
	"quill/internal/scanner"
	"strconv"
	"strings"
)

type Args struct {
	File      string
	Verbose   bool
	ParseOnly bool
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	runFile(args.File, args)
}

func parseArgs() (Args, error) {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")

	var parseOnly bool
	flag.BoolVar(&parseOnly, "p", false, "Parse only, do not run the program")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: quill [options] [file]\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	var file string
	if len(args) == 1 {
		file = args[0]
	}

	return Args{
		File:      file,
		Verbose:   verbose,
		ParseOnly: parseOnly,
	}, nil
}

func runFile(file string, args Args) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
		return
	}
	run(string(fileContent), args)
}

func run(source string, args Args) {
	scanner := scanner.New(source)
	tokens, scannerErrors := scanner.ScanTokens()

	if len(scannerErrors) > 0 {
		for _, err := range scannerErrors {
			fmt.Fprintf(os.Stderr, "ScannerError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

	if args.Verbose {
		fmt.Println("Tokens:")
		for _, token := range tokens {
			t := token.String()
			if t == "\n" {
				t = "\\n"
			}
			fmt.Printf("%s (%s)\n", t, token.Type)
		}
	}

	parser := parser.New(tokens)
	program, parserErrors := parser.Parse()

	if len(parserErrors) > 0 {
		for _, err := range parserErrors {
			fmt.Fprintf(os.Stderr, "ParseError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

	fmt.Println("File parsed successfully.")

	if args.Verbose {
		fmt.Println("Program:")
		fmt.Println(program)
	}

	if args.ParseOnly {
		fmt.Println("Parse only mode, exiting after parsing.")
		return
	}

	// Run the interpreter with the new result-based model
	runInterpreter(interpreter.New(program))
}

func runInterpreter(interp *interpreter.Interpreter) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("-- Starting script execution ---")

	for !interp.IsEnded() {
		result := interp.Step()

		switch result.Type {
		case interpreter.DialogResult:
			data := result.Data.(interpreter.DialogData)
			fmt.Printf("%s: %s", data.Character, data.Text)
			if len(data.Tags) > 0 {
				fmt.Printf(" [%s]", strings.Join(data.Tags, ", "))
			}
			fmt.Println()

		case interpreter.ChoiceResult:
			data := result.Data.(interpreter.ChoiceData)
			fmt.Println("\nChoices:")
			for _, option := range data.Options {
				fmt.Printf("%d. %s", option.Index+1, option.Text)
				if len(option.Tags) > 0 {
					fmt.Printf(" [%s]", strings.Join(option.Tags, ", "))
				}
				fmt.Println()
			}

			// Get user input for choice
			for {
				fmt.Print("\nEnter your choice (1-" + strconv.Itoa(len(data.Options)) + "): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
					return
				}

				input = strings.TrimSpace(input)
				choice, err := strconv.Atoi(input)
				if err != nil || choice < 1 || choice > len(data.Options) {
					fmt.Println("Invalid choice. Please try again.")
					continue
				}

				// Handle the choice (convert from 1-based to 0-based index)
				choiceResult := interp.HandleChoiceInput(choice - 1)
				if choiceResult.Type == interpreter.ErrorResult {
					errorData := choiceResult.Data.(interpreter.ErrorData)
					fmt.Fprintf(os.Stderr, "Error: %s\n", errorData.Message)
					return
				}

				// Process the result from handling the choice
				if choiceResult.Type == interpreter.DialogResult {
					dialogData := choiceResult.Data.(interpreter.DialogData)
					fmt.Printf("%s: %s", dialogData.Character, dialogData.Text)
					if len(dialogData.Tags) > 0 {
						fmt.Printf(" [%s]", strings.Join(dialogData.Tags, ", "))
					}
					fmt.Println()
				}
				// Note: We don't handle other result types here as they'll be processed in the next loop iteration

				break
			}

		case interpreter.ToolCallResult:
			data := result.Data.(interpreter.ToolCallData)
			fmt.Printf("\n--- Tool Call: %s ---\n", data.Function)
			fmt.Printf("Arguments: %v\n", data.Arguments)

			// Mock the external system response
			mockResult := mockToolCall(data.Function, data.Arguments)
			fmt.Printf("Mock result: %v\n", mockResult)

			// Send the result back to the interpreter
			toolResult := interp.HandleToolCallResponse(mockResult)
			if toolResult != nil {
				if toolResult.Type == interpreter.ErrorResult {
					errorData := toolResult.Data.(interpreter.ErrorData)
					fmt.Fprintf(os.Stderr, "Error: %s\n", errorData.Message)
					return
				}

				// Handle other result types if needed
				if toolResult.Type == interpreter.DialogResult {
					dialogData := toolResult.Data.(interpreter.DialogData)
					fmt.Printf("%s: %s", dialogData.Character, dialogData.Text)
					if len(dialogData.Tags) > 0 {
						fmt.Printf(" [%s]", strings.Join(dialogData.Tags, ", "))
					}
					fmt.Println()
				}
			}
			// If toolResult is nil, the LET statement completed and execution will continue in next loop iteration

		case interpreter.EndResult:
			fmt.Println("\n--- End of script ---")
			return

		case interpreter.ErrorResult:
			errorData := result.Data.(interpreter.ErrorData)
			fmt.Fprintf(os.Stderr, "Runtime Error at line %d: %s\n", errorData.Line, errorData.Message)
			return

		default:
			fmt.Fprintf(os.Stderr, "Unknown result type: %d\n", result.Type)
			return
		}
	}
}

func mockToolCall(functionName string, args []interface{}) interface{} {
	// Mock implementation of tool calls for testing
	switch functionName {
	case "getPlayerName":
		return "Player"

	case "getPlayerAge":
		return int64(25)

	case "getData":
		if len(args) > 0 {
			key := fmt.Sprintf("%v", args[0])
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
		if len(args) >= 2 {
			itemType := fmt.Sprintf("%v", args[0])
			level := int64(1)
			if levelArg, ok := args[1].(int64); ok {
				level = levelArg
			}

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
		if len(args) > 0 {
			if age, ok := args[0].(int64); ok {
				return age + 5
			}
		}
		return int64(5)

	default:
		return "Unknown function: " + functionName
	}
}
