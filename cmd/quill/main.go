package main

import (
	"fmt"
	"os"
	"quill/pkg/interpreter"
	"quill/pkg/parser"
	"quill/pkg/scanner"
)

type Args struct {
	File string
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Println("Usage: quill [file]")
		return
	}

	if args.File != "" {
		runFile(args.File)
	} else {
		runPrompt()
	}
}

func parseArgs(args []string) (Args, error) {
	if len(args) == 0 {
		return Args{}, nil
	} else if len(args) > 1 {
		return Args{}, fmt.Errorf("too many arguments")
	}
	return Args{File: args[0]}, nil
}

func runFile(file string) {
	fileConent, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
		return
	}
	run(string(fileConent))
}

func runPrompt() {
	fmt.Println("Quill is running in prompt mode. Type your input and press Enter.")
	for {
		fmt.Print("> ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}
		if input == "quit" {
			return
		}
		run(input)
	}
}

func run(source string) {
	scanner := scanner.New(source)
	tokens, errors := scanner.ScanTokens()

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "ScannerError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

	parser := parser.New(tokens)
	program := parser.Parse()

	interpreter := interpreter.New(program)
	err := interpreter.Interpret()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return
	}
}
