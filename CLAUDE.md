# Rush Programming Language - Developer Guide

## Development Best Practices

- Always write Go style unit and integration tests for new language features.
- Always update the README.md and CLAUDE.md on phase completion.
- Test new operators, built-in functions, and syntax additions thoroughly.
- Run `go test ./...` before committing changes.
- Whenever we begin work on a new phase of development, always create a new git branch.
- Always reference the related Github Issue in a Github Pull Request
- Boolean predicate functions in Rush need to be camel cased, for instance `is_whitespace?`

## Project Planning & Tracking

For detailed project plans, implementation roadmaps, and feature tracking, see:
https://github.com/zhubert/rush/issues

**Current Phase**: Maintenance & Bug Fixes - **ONGOING**
**Previous Phase**: Phase 13.4 - JIT Compilation - **COMPLETE**

## Recent Updates (January 2025)

### Bug Fixes & Cleanup
- **VM Stack Underflow Fix**: Resolved critical panic in nested if-else statements during bytecode execution
- **LLVM Removal**: Cleaned up unused LLVM infrastructure and dependencies from AOT compilation system
- **CI Simplification**: Streamlined GitHub Actions workflow to use standard Go toolchain only

### Current Execution Modes
Rush supports three high-performance execution modes:
1. **Tree-walking Interpreter** (default) - Direct AST evaluation
2. **Bytecode VM** (`-bytecode` flag) - Stack-based virtual machine
3. **JIT Compilation** (`-jit` flag) - Native ARM64 machine code generation

## Codebase Architecture

### Key Directories

```
rush/
├── cmd/rush/           # CLI entry point and main application
├── lexer/             # Tokenization - converts source code to tokens
├── parser/            # AST generation - converts tokens to Abstract Syntax Tree
├── ast/               # AST node definitions for all language constructs
├── interpreter/       # Runtime evaluation and built-in function implementations
├── vm/                # Bytecode virtual machine and stack-based execution
├── bytecode/          # Bytecode instruction definitions and serialization
├── compiler/          # Bytecode compiler (AST → bytecode)
├── jit/               # Just-In-Time compilation system (ARM64 target)
├── examples/          # Example Rush programs for testing and demonstration
├── std/              # Standard library modules (math.rush)
├── docs/              # User-facing documentation
└── tests/             # Test suite (mirrors main directory structure)
```

### Navigation Shortcuts

- **New operators**: Add to `lexer/lexer.go` → `parser/parser.go` → `interpreter/interpreter.go` + `compiler/compiler.go` + `vm/vm.go`
- **Built-in functions**: Add to `interpreter/builtins.go`
- **Standard library**: Add modules to `std/` directory
- **Language constructs**: Define AST nodes in `ast/ast.go`
- **Bytecode operations**: Add to `bytecode/instruction.go` and implement in `vm/vm.go`
- **Compilation**: Extend `compiler/compiler.go` for new AST nodes

## Development Workflow

### Common Commands

```bash
# Build and test
make build          # Build the Rush binary
make install        # Install Rush system-wide
go test ./...       # Run all tests
make dev FILE=example.rush  # Run file from source (development)
make repl          # Start REPL from source

# JIT compilation mode
rush -jit program.rush              # Enable JIT compilation
rush -jit -log-level=info program.rush  # JIT with statistics

# Testing specific components
go test ./lexer      # Test tokenization
go test ./parser     # Test AST generation
go test ./interpreter # Test runtime evaluation
go test ./compiler   # Test bytecode compilation
go test ./vm         # Test bytecode VM
go test ./jit        # Test JIT compilation
go test .           # Integration tests

# VM Execution with Logging (Bytecode Mode)
rush -bytecode -log-level=info program.rush     # Basic execution info
rush -bytecode -log-level=debug program.rush    # Detailed debugging
rush -bytecode -log-level=trace program.rush    # Full instruction tracing
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

## Current Development Context (Phase 12 Complete)

### Dot Notation Implementation

- **String methods**: Moved from `std/string.rush` to core interpreter with dot notation
- **Array methods**: Moved from `std/array.rush` to core interpreter with dot notation
- **Number methods**: Added mathematical operations with dot notation (abs, sqrt, pow, etc.)
- **Hash methods**: Built-in hash/dictionary operations with dot notation
- **JSON methods**: Comprehensive JSON processing with dot notation (parse, stringify, get, set, path, merge, format)
- **File system methods**: File, Directory, and Path objects with dot notation operations
- **Regular expressions**: Comprehensive Regexp type with dot notation methods and string integration
- **Method chaining**: All dot notation methods support seamless chaining
- **Standard library**: Now focuses on constants (PI, E) and multi-value operations

### Implementation Notes

- **Dot notation methods**: Implemented as Go built-ins in `interpreter/interpreter.go`
- **Method types**: StringMethod, ArrayMethod, NumberMethod, HashMethod, JSONMethod, FileMethod, DirectoryMethod, PathMethod, RegexpMethod value types in `interpreter/value.go`
- **Property access**: Extended `evalPropertyAccess` function handles all dot notation
- **Method application**: Separate apply functions (applyStringMethod, applyArrayMethod, applyJSONMethod, applyRegexpMethod, etc.)
- **JSON processing**: Core functions `json_parse()` and `json_stringify()` in `interpreter/builtins.go`
- **JSON object**: New JSON value type with comprehensive dot notation methods (get, set, has, keys, values, length, pretty, compact, validate, merge, path)
- **JSON features**: Full JSON standard compliance, method chaining, path navigation, immutable operations, rich formatting
- **Regexp processing**: `Regexp()` constructor function in `interpreter/builtins.go` creates regexp objects
- **Regexp object**: New Regexp value type with comprehensive dot notation methods (matches?, find_first, find_all, replace) and pattern property
- **String regexp integration**: String methods (match, replace, split, matches?) support both string and regexp arguments
- **Standard library**: Reduced to constants and multi-value operations in `std/` directory
- **Collections/hash operations**: Go built-ins with `builtin_hash_` prefix (legacy)
- **Module resolution**: Still happens in `interpreter/module.go`
- **Hash literals**: Use `{key: value}` syntax with support for string, integer, boolean, and float keys

## VM Debugging and Logging

The Rush VM includes comprehensive logging for debugging and performance analysis.

### VM Log Levels

- **none**: No logging (default for production)
- **error**: Only errors and critical issues
- **warn**: Warnings and errors
- **info**: Initialization, execution stats, performance metrics
- **debug**: Detailed execution flow, operations, stack state
- **trace**: Extremely verbose instruction-by-instruction execution

### When to Use Each Log Level

```bash
# Debug VM issues or understand execution flow
rush -bytecode -log-level=debug problematic_program.rush

# Performance analysis and optimization
rush -bytecode -log-level=info large_program.rush

# Deep debugging of instruction execution
rush -bytecode -log-level=trace simple_program.rush

# Production error tracking
rush -bytecode -log-level=error production_program.rush
```

### VM Debugging Workflow

1. **Start with info level** to see execution stats and identify performance issues
2. **Use debug level** to understand operation flow and identify logical errors
3. **Use trace level** for instruction-by-instruction debugging (very verbose)
4. **Check error level** for production deployments to catch critical issues

### Performance Metrics Available

- Execution time
- Instructions executed per second
- Stack operations count
- Function calls count
- Memory allocations
- Error count

## AI Assistant Guidelines

### Common Pitfalls

- Don't confuse built-in functions (Go) with standard library functions (Rush)
- Always test across lexer → parser → interpreter pipeline
- Remember to update precedence tables for new operators
- Standard library changes require testing both module loading and function execution
- VM operations need both interpreter and bytecode implementations
- Use VM logging to debug bytecode execution issues

### Quick Navigation

- Find built-ins: `rg "func.*args.*object.Object" interpreter/`
- Find dot notation methods: `rg "Method.*struct" interpreter/value.go`
- Find property access: `rg "evalPropertyAccess" interpreter/interpreter.go`
- Find method application: `rg "apply.*Method" interpreter/interpreter.go`
- Find JSON functions: `rg "json_parse\|json_stringify" interpreter/builtins.go`
- Find JSON methods: `rg "applyJSONMethod" interpreter/interpreter.go`
- Find JSON helpers: `rg "getJSONValue\|setJSONValue\|parseJSON" interpreter/interpreter.go`
- Find Regexp functions: `rg "Regexp" interpreter/builtins.go`
- Find Regexp methods: `rg "applyRegexpMethod" interpreter/interpreter.go`
- Find AST nodes: `rg "type.*struct" ast/ast.go`
- Find token types: `rg "= Token" token/token.go`
- Find test patterns: Look at existing `*_test.go` files in each directory
- Find VM operations: `rg "case bytecode\." vm/vm.go`
- Find bytecode definitions: `rg "Op.*=" bytecode/instruction.go`

### Dependencies

- Lexer → Parser → AST → Interpreter (linear dependency)
- Standard library depends on interpreter and module resolution
- Examples depend on all language features being implemented
```