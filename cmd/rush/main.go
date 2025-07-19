package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rush <file.rush>")
		os.Exit(1)
	}

	filename := os.Args[1]
	
	// Read the source file
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Phase 1: Tokenize and parse
	fmt.Printf("Rush interpreter - executing file: %s\n", filename)
	
	// Create lexer
	l := lexer.New(string(input))
	
	// Create parser
	p := parser.New(l)
	
	// Parse the program
	program := p.ParseProgram()
	
	// Check for parse errors
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}
	
	// Phase 2: Interpret and execute
	env := interpreter.NewEnvironment()
	result := interpreter.Eval(program, env)
	
	if result != nil {
		if result.Type() == "ERROR" {
			fmt.Printf("Runtime error: %s\n", result.Inspect())
			os.Exit(1)
		}
		fmt.Printf("Result: %s\n", result.Inspect())
	}
	
	fmt.Println("\nExecution complete!")
}