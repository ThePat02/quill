package main

import (
	"fmt"
	"os"
	"quill/pkg/scanner"
	"quill/pkg/token"
)

var hadError bool = false

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
		hadError = false
	}
}

func run(source string) {
	scanner := scanner.New(source, Error)
	tokens := scanner.ScanTokens()

	printTokens(tokens)

	if hadError {
		fmt.Println("Errors were encountered during scanning.")
		return
	}
}

func printTokens(tokens []token.Token) {
	for _, token := range tokens {
		fmt.Printf("%s\n", token.ToString())
	}
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("Error at line %d, %s: %s\n", line, where, message)
	hadError = true
}
