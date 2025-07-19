package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	fmt.Printf("Rush interpreter - parsing file: %s\n", filename)
	
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
	
	// Print the parsed AST
	fmt.Println("\nParsed AST:")
	fmt.Println(program.String())
	
	fmt.Println("\nPhase 1 complete: Successfully parsed!")
}