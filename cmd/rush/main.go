package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
)

func main() {
	if len(os.Args) < 2 {
		// Start REPL mode
		startREPL()
		return
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

func startREPL() {
	fmt.Println("Rush Interactive REPL")
	fmt.Println("Type ':help' for help, ':quit' to exit")
	fmt.Print("⛤ ")

	scanner := bufio.NewScanner(os.Stdin)
	env := interpreter.NewEnvironment()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Handle REPL commands
		if strings.HasPrefix(line, ":") {
			handleREPLCommand(line)
			fmt.Print("⛤ ")
			continue
		}
		
		// Skip empty lines
		if line == "" {
			fmt.Print("⛤ ")
			continue
		}
		
		// Evaluate the input
		evaluateInput(line, env)
		fmt.Print("⛤ ")
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

func handleREPLCommand(command string) {
	switch command {
	case ":help":
		fmt.Println("Available commands:")
		fmt.Println("  :help  - Show this help message")
		fmt.Println("  :quit  - Exit the REPL")
		fmt.Println("")
		fmt.Println("Enter Rush expressions to evaluate them interactively")
	case ":quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Type ':help' for available commands")
	}
}

func evaluateInput(input string, env *interpreter.Environment) {
	// Create lexer
	l := lexer.New(input)
	
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
		return
	}
	
	// Evaluate
	result := interpreter.Eval(program, env)
	
	if result != nil {
		if result.Type() == "ERROR" {
			fmt.Printf("Error: %s\n", result.Inspect())
		} else if result.Type() != "NULL" {
			fmt.Printf("%s\n", result.Inspect())
		}
	}
}