# Rush Programming Language - Developer Guide

## Development Best Practices

- Always write Go style unit and integration tests for new language features.
- Always update the README.md and CLAUDE.md on phase completion.
- Test new operators, built-in functions, and syntax additions thoroughly.
- Run `go test ./...` before committing changes.

## Project Planning & Tracking

For detailed project plans, implementation roadmaps, and feature tracking, see:
https://github.com/zhubert/rush/issues

**Current Phase**: Phase 11 - Enhanced Standard Library (string operations, array utilities, math functions, file I/O, JSON handling)

## Codebase Architecture

### Key Directories

```
rush/
├── cmd/rush/           # CLI entry point and main application
├── lexer/             # Tokenization - converts source code to tokens
├── parser/            # AST generation - converts tokens to Abstract Syntax Tree
├── ast/               # AST node definitions for all language constructs
├── interpreter/       # Runtime evaluation and built-in function implementations
├── examples/          # Example Rush programs for testing and demonstration
├── std/              # Standard library modules (math.rush, string.rush, array.rush)
├── docs/              # User-facing documentation
└── tests/             # Test suite (mirrors main directory structure)
```

### Navigation Shortcuts

- **New operators**: Add to `lexer/lexer.go` → `parser/parser.go` → `interpreter/interpreter.go`
- **Built-in functions**: Add to `interpreter/builtins.go`
- **Standard library**: Add modules to `std/` directory
- **Language constructs**: Define AST nodes in `ast/ast.go`
- **Error types**: Define in `interpreter/errors.go`

## Development Workflow

### Common Commands

```bash
# Build and test
make build          # Build the Rush binary
make install        # Install Rush system-wide
go test ./...       # Run all tests
make dev FILE=example.rush  # Run file from source (development)
make repl          # Start REPL from source

# Testing specific components
go test ./lexer      # Test tokenization
go test ./parser     # Test AST generation
go test ./interpreter # Test runtime evaluation
go test .           # Integration tests
```

### Running Examples

```bash
# After installation
rush examples/comprehensive_demo.rush

# From source during development
make dev FILE=examples/comprehensive_demo.rush
```

## Implementation Patterns

### Adding New Language Features

1. **Lexer** (`lexer/lexer.go`):

   - Add token type to `token/token.go`
   - Add tokenization logic in `NextToken()`

2. **Parser** (`parser/parser.go`):

   - Add AST node type to `ast/ast.go`
   - Add parsing function (e.g., `parseInfixExpression`)
   - Register precedence if needed

3. **Interpreter** (`interpreter/interpreter.go`):

   - Add evaluation case in `Eval()` function
   - Handle the new AST node type

4. **Testing**:
   - Add lexer tests for tokenization
   - Add parser tests for AST generation
   - Add interpreter tests for evaluation
   - Add integration tests for end-to-end behavior

### Built-in Function Pattern

```go
// In interpreter/builtins.go
"function_name": &object.Builtin{
    Fn: func(args ...object.Object) object.Object {
        // Argument validation
        if len(args) != expected {
            return newError("wrong number of arguments")
        }

        // Type checking
        arg, ok := args[0].(*object.String)
        if !ok {
            return newError("argument must be STRING")
        }

        // Implementation
        result := doSomething(arg.Value)
        return &object.String{Value: result}
    },
},
```

### Error Handling Pattern

```go
// Standard error creation
newError("descriptive error message: got %T", actualValue)

// Type validation
if len(args) != 1 {
    return newError("wrong number of arguments. got=%d, want=1", len(args))
}
```

## Current Development Context (Phase 11)

### Standard Library Development

- Focus on `std/` directory modules
- String operations in `std/string.rush`
- Array utilities in `std/array.rush`
- Math functions in `std/math.rush`
- Upcoming: file I/O, JSON handling, collections

### Implementation Notes

- Standard library functions are Rush functions, not Go built-ins
- Module resolution happens in `interpreter/module.go`
- Test standard library functions with Rush code examples

## AI Assistant Guidelines

### Common Pitfalls

- Don't confuse built-in functions (Go) with standard library functions (Rush)
- Always test across lexer → parser → interpreter pipeline
- Remember to update precedence tables for new operators
- Standard library changes require testing both module loading and function execution

### Quick Navigation

- Find built-ins: `rg "func.*args.*object.Object" interpreter/`
- Find AST nodes: `rg "type.*struct" ast/ast.go`
- Find token types: `rg "= Token" token/token.go`
- Find test patterns: Look at existing `*_test.go` files in each directory

### Dependencies

- Lexer → Parser → AST → Interpreter (linear dependency)
- Standard library depends on interpreter and module resolution
- Examples depend on all language features being implemented

