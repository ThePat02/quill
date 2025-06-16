package main

import (
	"flag"
	"fmt"
	"os"
	"quill/pkg/interpreter"
	"quill/pkg/parser"
	"quill/pkg/scanner"
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
	fileConent, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
		return
	}
	run(string(fileConent), args)
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

	interpreter := interpreter.New(program)
	interpreterErrors := interpreter.Interpret()

	if len(interpreterErrors) > 0 {
		for _, err := range interpreterErrors {
			fmt.Fprintf(os.Stderr, "InterpreterError at line %d: %s\n", err.Line, err.Message)
		}
		return
	}

}
