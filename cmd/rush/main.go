package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"rush/ast"
	"rush/bytecode"
	"rush/compiler"
	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
	"rush/vm"
)

func main() {
	// Define command line flags
	bytecodeMode := flag.Bool("bytecode", false, "Use bytecode compilation and VM execution")
	useCache := flag.Bool("cache", false, "Enable bytecode caching")
	clearCache := flag.Bool("clear-cache", false, "Clear bytecode cache and exit")
	cacheStats := flag.Bool("cache-stats", false, "Show cache statistics and exit")
	flag.Parse()

	// Handle cache management commands
	if *clearCache {
		err := bytecode.ClearCache()
		if err != nil {
			fmt.Printf("Error clearing cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cache cleared successfully")
		return
	}

	if *cacheStats {
		fileCount, totalSize, err := bytecode.GetCacheStats()
		if err != nil {
			fmt.Printf("Error getting cache stats: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cache statistics:\n")
		fmt.Printf("  Files: %d\n", fileCount)
		fmt.Printf("  Total size: %d bytes (%.2f KB)\n", totalSize, float64(totalSize)/1024)
		return
	}

	// Get remaining arguments after flag parsing
	args := flag.Args()
	if len(args) < 1 {
		// Start REPL mode
		startREPL(*bytecodeMode)
		return
	}

	filename := args[0]
	
	// Read the source file
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Execute the file using the selected mode
	if *bytecodeMode {
		fmt.Printf("Rush bytecode compiler - executing file: %s\n", filename)
		err := executeFileBytecode(filename, string(input), *useCache)
		if err != nil {
			fmt.Printf("Execution error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Rush tree-walking interpreter - executing file: %s\n", filename)
		err := executeFileTreeWalking(filename, string(input))
		if err != nil {
			fmt.Printf("Execution error: %v\n", err)
			os.Exit(1)
		}
	}
	
	fmt.Println("\nExecution complete!")
}

func startREPL(bytecodeMode bool) {
	if bytecodeMode {
		fmt.Println("Rush Interactive REPL (Bytecode Mode)")
	} else {
		fmt.Println("Rush Interactive REPL (Tree-Walking Mode)")
	}
	fmt.Println("Type ':help' for help, ':quit' to exit")
	fmt.Print("⛤ ")

	scanner := bufio.NewScanner(os.Stdin)
	env := interpreter.NewEnvironment()
	globals := make([]interpreter.Value, vm.GlobalsSize)

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
		if bytecodeMode {
			globals = evaluateInputBytecode(line, globals)
		} else {
			evaluateInputTreeWalking(line, env)
		}
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

// executeFileTreeWalking executes a file using the tree-walking interpreter
func executeFileTreeWalking(filename, source string) error {
	// Create lexer
	l := lexer.New(source)
	
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
		return fmt.Errorf("parse errors occurred")
	}
	
	// Phase 2: Interpret and execute
	env := interpreter.NewEnvironment()
	result := interpreter.Eval(program, env)
	
	if result != nil {
		if result.Type() == "ERROR" || result.Type() == "EXCEPTION" {
			return fmt.Errorf("runtime error: %s", result.Inspect())
		}
		if result.Type() != "NULL" {
			fmt.Printf("Result: %s\n", result.Inspect())
		}
	}
	
	return nil
}

// executeFileBytecode executes a file using bytecode compilation and VM
func executeFileBytecode(filename, source string, useCache bool) error {
	sourceHash := bytecode.HashSource(source)
	
	// Try to load from cache first
	var instructions bytecode.Instructions
	var constants []interpreter.Value
	var err error
	
	if useCache {
		instructions, constants, err = bytecode.LoadFromCache(filename, sourceHash)
		if err == nil {
			fmt.Println("Using cached bytecode")
		}
	}
	
	// If cache miss or cache disabled, compile from source
	if instructions == nil {
		fmt.Println("Compiling to bytecode...")
		
		// Parse the source
		l := lexer.New(source)
		p := parser.New(l)
		program := p.ParseProgram()
		
		errors := p.Errors()
		if len(errors) > 0 {
			fmt.Println("Parse errors:")
			for _, err := range errors {
				fmt.Printf("  %s\n", err)
			}
			return fmt.Errorf("parse errors occurred")
		}
		
		// Compile to bytecode
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			return fmt.Errorf("compilation error: %w", err)
		}
		
		compiledBytecode := comp.Bytecode()
		instructions = compiledBytecode.Instructions
		constants = compiledBytecode.Constants
		
		// Save to cache if enabled
		if useCache {
			err = bytecode.SaveToCache(filename, instructions, constants, sourceHash)
			if err != nil {
				fmt.Printf("Warning: failed to save to cache: %v\n", err)
			}
		}
	}
	
	// Execute with VM
	machine := vm.New(&compiler.Bytecode{
		Instructions: instructions,
		Constants:    constants,
	})
	
	err = machine.Run()
	if err != nil {
		return fmt.Errorf("VM error: %w", err)
	}
	
	// Get result
	stackTop := machine.StackTop()
	if stackTop != nil && stackTop.Type() != "NULL" {
		fmt.Printf("Result: %s\n", stackTop.Inspect())
	}
	
	return nil
}

func evaluateInputTreeWalking(input string, env *interpreter.Environment) {
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
		if result.Type() == "ERROR" || result.Type() == "EXCEPTION" {
			fmt.Printf("Error: %s\n", result.Inspect())
		} else if result.Type() != "NULL" {
			fmt.Printf("%s\n", result.Inspect())
		}
	}
}

func evaluateInputBytecode(input string, globals []interpreter.Value) []interpreter.Value {
	// Parse the input
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		return globals
	}
	
	// Compile to bytecode with REPL mode (don't pop last expression)
	comp := compiler.New()
	
	// Modify the program to avoid popping the last expression for REPL display
	if len(program.Statements) > 0 {
		if lastStmt, ok := program.Statements[len(program.Statements)-1].(*ast.ExpressionStatement); ok {
			// Compile all statements except the last
			for _, stmt := range program.Statements[:len(program.Statements)-1] {
				err := comp.Compile(stmt)
				if err != nil {
					fmt.Printf("Compilation error: %v\n", err)
					return globals
				}
			}
			// Compile the last expression without popping
			err := comp.Compile(lastStmt.Expression)
			if err != nil {
				fmt.Printf("Compilation error: %v\n", err)
				return globals
			}
		} else {
			// Normal compilation for non-expression statements
			err := comp.Compile(program)
			if err != nil {
				fmt.Printf("Compilation error: %v\n", err)
				return globals
			}
		}
	} else {
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			return globals
		}
	}
	
	// Execute with VM
	machine := vm.NewWithGlobalsStore(comp.Bytecode(), globals)
	err := machine.Run()
	if err != nil {
		fmt.Printf("VM error: %v\n", err)
		return globals
	}
	
	// Get result and print if not null
	stackTop := machine.StackTop()
	if stackTop != nil && stackTop.Type() != "NULL" {
		fmt.Printf("%s\n", stackTop.Inspect())
	}
	
	return globals
}