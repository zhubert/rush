# Rush Programming Language

[![Tests](https://img.shields.io/badge/tests-passing-green)](./tests)
[![Phase](https://img.shields.io/badge/phase-9%20complete%2C%2010--13%20planned-blue)](#phases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Rush is a modern, dynamically-typed programming language designed for simplicity and expressiveness. Built with Go, Rush features clean syntax inspired by modern languages with C-style control flow.

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/zhubert/rush.git
cd rush

# Install Rush system-wide (requires sudo)
make install

# Or just build without installing
make build
```

After installation, you can use Rush from anywhere:

```bash
# Execute a Rush file
rush examples/comprehensive_demo.rush

# Start interactive REPL
rush
```

### Development Mode

If you want to run Rush from source without installing:

```bash
# Execute a Rush file from source
make dev FILE=examples/comprehensive_demo.rush

# Start REPL from source  
make repl
```

### Hello World

Create `hello.rush`:
```rush
print("Hello, Rush!")
```

Run it:
```bash
# If installed
rush hello.rush

# From source
make dev FILE=hello.rush
```

## ✨ Features

- **Dynamic Typing**: Variables can hold any type of value
- **First-Class Functions**: Functions are values that can be passed around
- **Module System**: Import/export functionality for code organization
- **Error Handling**: Try/catch/finally/throw mechanisms with type-specific catching
- **Arrays**: Built-in support for dynamic arrays
- **String Manipulation**: Comprehensive string operations
- **Control Flow**: If/else, while, and for loops
- **Interactive REPL**: Explore Rush interactively
- **Clean Syntax**: Easy to read and write

## 📚 Language Overview

### Variables and Types
```rush
# Variables are dynamically typed
name = "Rush"           # String
version = 1.0           # Float
count = 42              # Integer
ready = true            # Boolean
items = [1, 2, 3]       # Array
```

### Functions
```rush
# Function definition
add = fn(x, y) {
  return x + y
}

# Function call
result = add(10, 20)    # 30

# Recursive functions
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
```

### Control Flow
```rush
# If-else statements
if (score >= 90) {
  grade = "A"
} else {
  grade = "B"
}

# While loops
i = 0
while (i < 10) {
  print(i)
  i = i + 1
}

# For loops
for (j = 0; j < 5; j = j + 1) {
  print("Count: " + type(j))
}
```

### Arrays and Strings
```rush
# Array operations
numbers = [1, 2, 3, 4, 5]
first = numbers[0]              # Array indexing
length = len(numbers)           # Get length
extended = push(numbers, 6)     # Add element
subset = slice(numbers, 1, 4)   # Get slice

# String operations
text = "Hello, World!"
chars = split(text, ", ")       # Split string
hello = substr(text, 0, 5)      # Substring
```

### Module System
```rush
# math.rush - Define a module
export add = fn(x, y) { return x + y }
export PI = 3.14159
export square = fn(x) { return x * x }

# main.rush - Use the module
import { add, PI, square } from "./math"
result = add(10, 5)             # 15
area = PI * square(3)           # ~28.27
print("Result:", result, "Area:", area)
```

### Standard Library
Rush includes a standard library with useful modules (available after installation):

```rush
# Import math functions
import { PI, sin, cos, sqrt } from "std/math"
print("π =", PI)
print("sin(π/2) =", sin(PI / 2))

# Import array utilities
import { map, filter, reduce } from "std/array" 
numbers = [1, 2, 3, 4, 5]
doubled = map(numbers, fn(x) { return x * 2 })

# Import string utilities  
import { trim, upper, lower } from "std/string"
text = trim("  hello world  ")
```

### Error Handling
```rush
# Try-catch-finally blocks
try {
  result = divide(10, 0)
  print("Result:", result)
} catch (error) {
  print("Error occurred:", error.message)
} finally {
  print("Cleanup complete")
}

# Type-specific error catching
try {
  validateInput(data)
} catch (ValidationError error) {
  print("Validation failed:", error.message)
} catch (TypeError error) {
  print("Type error:", error.message)
}

# Throwing custom errors
throw RuntimeError("Something went wrong")
```

## 🛠️ Built-in Functions

### I/O Functions
- `print(...)` - Output values to console

### Utility Functions
- `len(collection)` - Get length of array or string
- `type(value)` - Get type of value as string

### String Functions
- `substr(string, start, length)` - Extract substring
- `split(string, separator)` - Split string into array

### Array Functions
- `push(array, element)` - Add element to array
- `pop(array)` - Get last element of array
- `slice(array, start, end)` - Extract array slice

## 📖 Documentation

- **[Language Specification](docs/LANGUAGE_SPECIFICATION.md)** - Complete syntax and grammar reference
- **[User Guide](docs/USER_GUIDE.md)** - Comprehensive tutorial and examples
- **[API Reference](docs/API_REFERENCE.md)** - Built-in functions documentation

## 🎮 Examples

### Basic Calculator
```rush
calculator = fn(op, a, b) {
  if (op == "+") {
    return a + b
  } else {
    if (op == "-") {
      return a - b
    } else {
      if (op == "*") {
        return a * b
      } else {
        if (op == "/") {
          return a / b
        } else {
          return 0
        }
      }
    }
  }
}

print("5 + 3 =", calculator("+", 5, 3))
print("10 * 2 =", calculator("*", 10, 2))
```

### Array Processing
```rush
# Filter even numbers
filter_evens = fn(arr) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    num = arr[i]
    if ((num / 2) * 2 == num) {  # Even check
      result = push(result, num)
    }
  }
  return result
}

numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
evens = filter_evens(numbers)
print("Even numbers:", evens)
```

### More Examples
- **[Comprehensive Demo](examples/comprehensive_demo.rush)** - All language features
- **[Algorithms Demo](examples/algorithms_demo.rush)** - Sorting, searching, and algorithms
- **[Games Demo](examples/game_demo.rush)** - Interactive games and simulations

## 🧪 Interactive REPL

Start the REPL for interactive exploration:

```bash
go run cmd/rush/main.go
```

```
Rush Interactive REPL
Type ':help' for help, ':quit' to exit
⛤ greet = fn(name) { return "Hello, " + name + "!" }
⛤ greet("Rush")
Hello, Rush!
⛤ numbers = [1, 2, 3, 4, 5]
⛤ len(numbers)
5
⛤ :quit
```

## 🏗️ Architecture

Rush is implemented in Go with a clean, modular architecture:

```
rush/
├── cmd/rush/           # CLI entry point
├── lexer/             # Tokenization
├── parser/            # AST generation
├── ast/               # AST node definitions
├── interpreter/       # Runtime evaluation
├── examples/          # Example programs
├── docs/              # Documentation
└── tests/             # Test suite
```

## ✅ Testing

Rush includes comprehensive tests covering all components:

```bash
# Run all tests
go test ./...

# Run specific test suites
go test ./lexer      # Lexer tests
go test ./parser     # Parser tests
go test ./interpreter # Interpreter tests
go test .            # Integration tests
```

### Test Coverage
- **Lexer Tests**: Tokenization of all language constructs
- **Parser Tests**: AST generation for all syntax
- **Interpreter Tests**: Runtime evaluation and built-ins
- **Integration Tests**: End-to-end program execution
- **Example Tests**: All example programs execute successfully

## 📈 Development Phases

Rush was developed through a systematic phased approach:

### Core Language (Completed)
- **✅ Phase 1**: Foundation & Basic Parsing
- **✅ Phase 2**: Basic Interpreter  
- **✅ Phase 3**: Control Flow
- **✅ Phase 4**: Functions
- **✅ Phase 5**: Loops & Arrays
- **✅ Phase 6**: Core Polish & Tools
- **✅ Phase 7**: Testing & Documentation

### Advanced Features (In Progress)
- **✅ Phase 8**: Module System - Import/export functionality
- **✅ Phase 9**: Error Handling - Try/catch/finally/throw mechanisms
- **🔄 Phase 10**: Object System - Classes and inheritance
- **🔄 Phase 11**: Enhanced Standard Library - File I/O, JSON, regex
- **🔄 Phase 12**: Performance & Compilation - Bytecode optimization
- **🔄 Phase 13**: Developer Tooling - Debugger, language server, IDE integration

**📋 Development Tracking**: Phases under development are tracked in the [Initial Development GitHub Project](https://github.com/users/zhubert/projects/1).

## 🤝 Contributing

Rush is an experimental language built for learning and exploration. Feel free to:

1. Explore the codebase
2. Run the examples
3. Experiment in the REPL
4. Suggest improvements

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎯 Design Goals

Rush was designed with these principles:

1. **Simplicity**: Easy to learn and understand syntax
2. **Expressiveness**: Powerful features without complexity
3. **Consistency**: Uniform patterns throughout the language
4. **Interactivity**: REPL-first development experience
5. **Clarity**: Clear error messages and predictable behavior

---