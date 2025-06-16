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

func run(source string) {
	scanner := scanner.New(source)
	tokens, scannerErrors := scanner.ScanTokens()

	if len(scannerErrors) > 0 {
		for _, err := range scannerErrors {
			fmt.Fprintf(os.Stderr, "ScannerError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

	parser := parser.New(tokens)
	program, parserErrors := parser.Parse()

	if len(parserErrors) > 0 {
		for _, err := range parserErrors {
			fmt.Fprintf(os.Stderr, "ParseError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

	interpreter := interpreter.New(program)
	interpreterErrors := interpreter.Interpret()

	if len(interpreterErrors) > 0 {
		for _, err := range interpreterErrors {
			fmt.Fprintf(os.Stderr, "InterpreterError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}
}
