# Rush Programming Language

[![Tests](https://img.shields.io/badge/tests-passing-green)](./tests)
[![Phase](https://img.shields.io/badge/phase-11%20complete-green)](#phases)
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

### Core Language
- **Dynamic Typing**: Variables can hold any type of value
- **First-Class Functions**: Functions with closures and higher-order support
- **Object-Oriented Programming**: Classes, inheritance, and method calls
- **Module System**: Import/export with aliasing for code organization
- **Error Handling**: Try/catch/finally/throw with typed error catching
- **Control Flow**: If/else, while, for loops, switch/case, break/continue
- **Interactive REPL**: Explore Rush interactively

### Data Types & Operations
- **Arrays**: Dynamic arrays with element assignment (`arr[i] = value`)
- **Hashes/Dictionaries**: Key-value mappings with `{key: value}` syntax
- **Strings**: String indexing and comprehensive manipulation
- **Numbers**: Integers and floats with modulo operator (`%`)
- **Booleans**: Logical operations with short-circuit evaluation
- **Null**: Explicit null handling

### Standard Library
- **Math Module** (`std/math`): Mathematical functions, constants, and utilities
- **String Module** (`std/string`): String manipulation and text processing
- **Array Module** (`std/array`): Functional programming utilities (map, filter, reduce)
- **Collections Module**: Built-in hash/dictionary operations and utilities
- **Import Aliasing**: Clean imports with `import { func as alias } from "module"`

### Development Experience
- **System Installation**: Install globally with `make install`
- **Development Tools**: REPL, file execution, comprehensive error messages
- **Clean Syntax**: Modern, readable syntax inspired by contemporary languages

## 📚 Language Overview

### Variables and Types
```rush
# Variables are dynamically typed
name = "Rush"           # String
version = 1.0           # Float
count = 42              # Integer
ready = true            # Boolean
items = [1, 2, 3]       # Array
person = {"name": "Alice", "age": 30}  # Hash/Dictionary
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

# Switch statements (Go-style automatic break)
grade = "B"
switch (grade) {
  case "A":
    result = "Excellent!"
  case "B", "C":
    result = "Good work"
  case "D":
    result = "Needs improvement"
  default:
    result = "Invalid grade"
}

# Modulo operator
remainder = 10 % 3        # 1
even = (number % 2) == 0  # Check if even
```

### Arrays, Strings, and Hashes
```rush
# Array operations
numbers = [1, 2, 3, 4, 5]
first = numbers[0]              # Array indexing
numbers[0] = 99                 # Array assignment
length = len(numbers)           # Get length
extended = push(numbers, 6)     # Add element
subset = slice(numbers, 1, 4)   # Get slice

# String operations
text = "Hello, World!"
chars = split(text, ", ")       # Split string
hello = substr(text, 0, 5)      # Substring

# Hash/Dictionary operations
import { keys, values, has_key?, get, set, delete, merge } from "std/collections"

person = {"name": "Alice", "age": 30, "active": true}
name = person["name"]           # Hash indexing
person["city"] = "NYC"          # Hash assignment
person_keys = keys(person)      # Get all keys
person_values = values(person)  # Get all values
has_age = has_key?(person, "age") # Check key existence
empty_hash = {}                 # Empty hash literal

# Multiline hash literals supported
config = {
  "database": {"host": "localhost", "port": 5432},
  "cache": {"enabled": true, "ttl": 3600}
}
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
Rush includes a comprehensive standard library with three main modules:

```rush
# Math operations
import { PI, E, sqrt, abs, min, max, sum } from "std/math"
numbers = [1, 2, 3, 4, 5]
print("π =", PI)                    # 3.141592653589793
print("√16 =", sqrt(16))            # 4.0
print("Sum:", sum(numbers))         # 15.0

# String manipulation
import { trim, upper, split, join, contains? } from "std/string"
text = "  hello, rush world  "
words = split(trim(text), ", ")     # ["hello", "rush world"]
sentence = join(words, " and ")     # "hello and rush world"
print(upper(sentence))              # "HELLO AND RUSH WORLD"

# Array utilities with functional programming
import { map, filter, reduce } from "std/array"
numbers = [1, 2, 3, 4, 5]
doubled = map(numbers, fn(x) { x * 2 })          # [2, 4, 6, 8, 10]
evens = filter(numbers, fn(x) { x % 2 == 0 })    # [2, 4]
total = reduce(numbers, fn(acc, x) { acc + x }, 0) # 15
```

### Object-Oriented Programming
```rush
# Define a class
class Animal {
  fn initialize(name) {
    @name = name
  }
  
  fn speak() {
    return @name + " makes a sound"
  }
}

# Inheritance
class Dog < Animal {
  fn speak() {
    return @name + " barks"
  }
}

# Create instances
dog = Dog.new("Buddy")
print(dog.speak())  # "Buddy barks"
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
- `ord(char)` - Get ASCII code of character
- `chr(code)` - Get character from ASCII code

### String Functions
- `substr(string, start, length)` - Extract substring
- `split(string, separator)` - Split string into array

### Array Functions
- `push(array, element)` - Add element to array
- `pop(array)` - Get last element of array
- `slice(array, start, end)` - Extract array slice

### Hash/Dictionary Functions (Collections Module)
```rush
import { keys, values, has_key?, get, set, delete, merge, filter, each } from "std/collections"
```
- `keys(hash)` - Get array of all keys
- `values(hash)` - Get array of all values
- `has_key?(hash, key)` - Check if key exists
- `get(hash, key, default?)` - Get value with optional default
- `set(hash, key, value)` - Create new hash with key-value pair
- `delete(hash, key)` - Create new hash without key
- `merge(hash1, hash2)` - Merge two hashes (hash2 overwrites)
- `filter(hash, predicate)` - Filter hash by key-value predicate
- `each(hash, callback)` - Iterate over hash with callback

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

# Array assignment example
evens[0] = 999          # Modify first element
print("Modified:", evens)
```

### Hash/Dictionary Processing
```rush
# User database example
users = {
  "alice": {"name": "Alice Smith", "age": 30, "role": "admin"},
  "bob": {"name": "Bob Jones", "age": 25, "role": "user"}
}

# Access nested data
alice_name = users["alice"]["name"]
print("User:", alice_name)

# Add new user
users["charlie"] = {"name": "Charlie Brown", "age": 35, "role": "user"}

# Check user permissions
check_admin = fn(username) {
  if (has_key?(users, username)) {
    user = users[username]
    return user["role"] == "admin"
  }
  return false
}

print("Alice is admin:", check_admin("alice"))  # true
print("Bob is admin:", check_admin("bob"))      # false
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
- **✅ Phase 10**: Object System - Classes, inheritance, and method enhancements
- **🔄 Phase 11**: Enhanced Standard Library - String/array/math/collections modules, JSON, file I/O
- **🔄 Phase 12**: Performance & Compilation - Bytecode optimization, JIT compilation
- **🔄 Phase 13**: Developer Tooling - Language server, debugger, IDE integration

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