# Rush Programming Language

[![Tests](https://img.shields.io/badge/tests-passing-green)](./tests)
[![Phase](https://img.shields.io/badge/phase-maintenance%20%26%20fixes-blue)](#phases)
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

# Build with JIT support (ARM64 only)
make build
```

After installation, you can use Rush from anywhere:

```bash
# Execute a Rush file
rush examples/comprehensive_demo.rush

# Execute with bytecode VM for better performance
rush -bytecode examples/comprehensive_demo.rush

# Execute with JIT compilation (ARM64 only)
rush -jit examples/comprehensive_demo.rush

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

# With JIT compilation
rush -jit hello.rush
```

## ✨ Features

### Core Language
- **Dynamic Typing**: Variables can hold any type of value
- **First-Class Functions**: Functions with closures and higher-order support
- **Object-Oriented Programming**: Classes, inheritance, and method calls
- **Module System**: Import/export with aliasing for code organization
- **Error Handling**: Try/catch/finally/throw with typed error catching
- **Control Flow**: If/else, while, for loops, switch/case, break/continue
- **Regular Expressions**: Built-in regexp support with `Regexp()` constructor
- **Interactive REPL**: Explore Rush interactively

### Data Types & Operations
- **Arrays**: Dynamic arrays with element assignment and dot notation methods (`arr.length`, `arr.map()`)
- **Hashes/Dictionaries**: Key-value mappings with `{key: value}` syntax and dot notation methods
- **Strings**: String indexing and dot notation methods (`str.length`, `str.upper()`)
- **Numbers**: Integers and floats with modulo operator and dot notation methods (`num.abs()`, `num.sqrt()`)
- **Booleans**: Logical operations with short-circuit evaluation
- **Null**: Explicit null handling

### Built-in Dot Notation & Standard Library
- **String Dot Notation**: Built-in string methods (`str.trim()`, `str.upper()`, `str.split()`) - no imports needed!
- **Array Dot Notation**: Built-in array methods (`arr.map()`, `arr.filter()`, `arr.reduce()`) - no imports needed!
- **Number Dot Notation**: Built-in math methods (`num.abs()`, `num.sqrt()`, `num.pow()`) - no imports needed!
- **Hash Dot Notation**: Built-in hash/dictionary operations and utilities - no imports needed!
- **Math Module** (`std/math`): Mathematical constants (PI, E) and multi-value operations
- **Import Aliasing**: Clean imports with `import { func as alias } from "module"`

### Development Experience
- **System Installation**: Install globally with `make install`
- **Development Tools**: REPL, file execution, comprehensive error messages
- **Clean Syntax**: Modern, readable syntax inspired by contemporary languages
- **JIT Compilation**: ARM64 JIT compiler for high-performance execution
- **Tiered Execution**: Tree-walking interpreter, bytecode VM, and JIT compilation modes

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
# Array operations with dot notation
numbers = [1, 2, 3, 4, 5]
first = numbers[0]              # Array indexing
numbers[0] = 99                 # Array assignment
length = numbers.length         # Get length
numbers.push(6)                 # Add element (mutates array)
subset = numbers.slice(1, 4)    # Get slice

# Functional programming with arrays
doubled = numbers.map(fn(x) { x * 2 })           # Transform elements
evens = numbers.filter(fn(x) { x % 2 == 0 })     # Filter elements
total = numbers.reduce(fn(acc, x) { acc + x }, 0) # Reduce to single value

# String operations with dot notation
text = "Hello, World!"
chars = text.split(", ")       # Split string
hello = text.substr(0, 5)       # Substring
upper_text = text.upper()       # Convert to uppercase
length = text.length            # Get string length

# Hash/Dictionary operations with dot notation
person = {"name": "Alice", "age": 30, "active": true}
name = person["name"]           # Hash indexing
person["city"] = "NYC"          # Hash assignment

# Properties (no imports needed!)
person_keys = person.keys       # Get all keys
person_values = person.values   # Get all values  
count = person.length           # Get count
is_empty = person.empty         # Check if empty

# Methods
has_age = person.has_key?("age")           # Check key existence
safe_get = person.get("phone", "N/A")     # Get with default
new_person = person.set("email", "alice@example.com")  # Immutable set
filtered = person.filter(fn(k, v) { k != "active" })  # Filter keys

# Multiline hash literals supported
config = {
  "database": {"host": "localhost", "port": 5432},
  "cache": {"enabled": true, "ttl": 3600}
}

# Method chaining
result = person.set("score", 95).delete("age").merge({"grade": "A"})
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

### Built-in Dot Notation Methods
Rush features comprehensive dot notation methods built directly into the language - no imports needed!

```rush
# String operations with dot notation
text = "  hello, rush world  "
words = text.trim().split(", ")     # ["hello", "rush world"]
sentence = "and".join(words)        # "hello and rush world"
print(sentence.upper())              # "HELLO AND RUSH WORLD"

# Array operations with dot notation
numbers = [1, 2, 3, 4, 5]
doubled = numbers.map(fn(x) { x * 2 })          # [2, 4, 6, 8, 10]
evens = numbers.filter(fn(x) { x % 2 == 0 })    # [2, 4]
total = numbers.reduce(fn(acc, x) { acc + x }, 0) # 15

# Number operations with dot notation
result = (-16).abs().sqrt().pow(2)   # |(−16)| = 16, √16 = 4, 4² = 16
print("Result:", result)             # 16

# Method chaining works seamlessly
processed = "  Hello World  ".trim().lower().replace(" ", "_")
print(processed)                     # "hello_world"
```

### Standard Library
The standard library now focuses on constants and multi-value operations:

```rush
# Math constants and utilities
import { PI, E, min, max, sum } from "std/math"
numbers = [1, 2, 3, 4, 5]
print("π =", PI)                    # 3.141592653589793
print("Sum:", sum(numbers))         # 15.0
print("Max:", max(numbers))         # 5
```

### Regular Expressions
```rush
# Create regexp objects
email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
number_pattern = Regexp("\\d+")

# String methods with regexp support
text = "Contact john@example.com for info"
emails = text.match(email_pattern)        # Find matches
clean_text = text.replace(email_pattern, "[EMAIL]")  # Replace with string
parts = text.split(number_pattern)        # Split by pattern

# Regexp object methods
has_email = email_pattern.matches?(text)  # Boolean test
first_email = email_pattern.find_first(text)  # First match
all_emails = email_pattern.find_all(text)     # All matches
no_emails = email_pattern.replace(text, "[HIDDEN]")  # Replace all
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

### Regular Expression Functions
- `Regexp(pattern)` - Create a regular expression object from pattern string

### String Methods (Dot Notation)
No imports needed - all methods are built into string objects!

**Properties:**
- `string.length` - Get string length
- `string.empty` - Boolean, true if string is empty

**Methods:**
- `string.trim()` - Remove leading/trailing whitespace
- `string.ltrim()` - Remove leading whitespace
- `string.rtrim()` - Remove trailing whitespace
- `string.upper()` - Convert to uppercase
- `string.lower()` - Convert to lowercase
- `string.contains?(substring)` - Check if string contains substring
- `string.starts_with?(prefix)` - Check if string starts with prefix
- `string.ends_with?(suffix)` - Check if string ends with suffix
- `string.replace(old, new)` - Replace all occurrences (supports regexp objects)
- `string.substr(start, length)` - Extract substring
- `string.split(separator)` - Split into array (supports regexp objects)
- `string.join(array)` - Join array elements with string as separator
- `string.match(regexp)` - Find all matches of regexp pattern
- `string.matches?(regexp)` - Test if string matches regexp pattern

### Array Methods (Dot Notation)
No imports needed - all methods are built into array objects!

**Properties:**
- `array.length` - Get array length
- `array.empty` - Boolean, true if array is empty

**Methods:**
- `array.map(transform_fn)` - Transform each element
- `array.filter(predicate_fn)` - Filter elements by predicate
- `array.reduce(reducer_fn, initial)` - Reduce to single value
- `array.find(predicate_fn)` - Find first matching element
- `array.index_of(element)` - Get index of element (-1 if not found)
- `array.includes?(element)` - Check if array contains element
- `array.reverse()` - Create new reversed array
- `array.sort()` - Create new sorted array
- `array.push(element)` - Add element to end (mutates array)
- `array.pop()` - Remove and return last element
- `array.slice(start, end)` - Extract array slice

### Number Methods (Dot Notation)
No imports needed - all methods are built into number objects!

**Methods:**
- `number.abs()` - Absolute value
- `number.floor()` - Round down to nearest integer
- `number.ceil()` - Round up to nearest integer
- `number.round()` - Round to nearest integer
- `number.sqrt()` - Square root
- `number.pow(exponent)` - Raise to power

### Hash/Dictionary Methods (Dot Notation)
No imports needed - all methods are built into hash objects!

**Properties:**
- `hash.keys` - Array of all keys
- `hash.values` - Array of all values  
- `hash.length` / `hash.size` - Number of key-value pairs
- `hash.empty` - Boolean, true if hash has no pairs

**Methods:**
- `hash.has_key?(key)` - Check if key exists
- `hash.has_value?(value)` - Check if value exists
- `hash.get(key, default?)` - Get value with optional default
- `hash.set(key, value)` - Create new hash with key-value pair (immutable)
- `hash.delete(key)` - Create new hash without key (immutable)
- `hash.merge(other_hash)` - Merge two hashes (other overwrites conflicts)
- `hash.filter(predicate_fn)` - Filter hash by key-value predicate
- `hash.map_values(transform_fn)` - Transform all values
- `hash.each(callback_fn)` - Iterate over hash with callback
- `hash.select_keys(key_array)` - Create hash with only specified keys
- `hash.reject_keys(key_array)` - Create hash without specified keys
- `hash.invert()` - Create hash with keys and values swapped
- `hash.to_array()` - Convert to array of `[key, value]` pairs

**Standalone Function:**
- `array_to_hash(pairs)` - Convert array of `[key, value]` pairs to hash

### Regular Expression Methods (Dot Notation)
No imports needed - all methods are built into regexp objects!

**Properties:**
- `regexp.pattern` - Get the original pattern string

**Methods:**
- `regexp.matches?(text)` - Test if text matches the pattern
- `regexp.find_first(text)` - Find first match in text
- `regexp.find_all(text)` - Find all matches in text
- `regexp.replace(text, replacement)` - Replace all matches with replacement string

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
  if (users.has_key?(username)) {
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
- **[Regular Expression Demo](examples/regexp_demo.rush)** - Regular expression patterns and methods
- **[Algorithms Demo](examples/algorithms_demo.rush)** - Sorting, searching, and algorithms
- **[Games Demo](examples/game_demo.rush)** - Interactive games and simulations

## ⚡ Performance Modes

Rush offers multiple execution modes for different performance needs:

### Tree-Walking Interpreter (Default)
```bash
# Default mode - good for development and debugging
rush program.rush
```

### Bytecode Virtual Machine
```bash
# Faster execution with bytecode compilation
rush -bytecode program.rush

# With performance monitoring
rush -bytecode -log-level=info program.rush
```

### JIT Compilation (ARM64)
```bash
# Just-In-Time compilation for maximum performance
rush -jit program.rush

# JIT with detailed statistics
rush -jit -log-level=info program.rush
```

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
├── vm/                # Bytecode virtual machine
├── bytecode/          # Bytecode instruction definitions
├── compiler/          # Bytecode compiler (AST → bytecode)
├── jit/               # Just-In-Time ARM64 compiler
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
go test ./compiler   # Compiler tests
go test ./vm         # VM tests
go test ./jit        # JIT tests
go test .            # Integration tests

# Run JIT benchmarks (ARM64 only)
go test -bench=. ./jit_integration_test.go
```

### Test Coverage
- **Lexer Tests**: Tokenization of all language constructs
- **Parser Tests**: AST generation for all syntax
- **Interpreter Tests**: Runtime evaluation and built-ins
- **Compiler Tests**: Bytecode compilation and optimization
- **VM Tests**: Bytecode execution and stack operations
- **JIT Tests**: ARM64 code generation and optimization
- **Integration Tests**: End-to-end program execution across all modes
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
- **✅ Phase 11**: Enhanced Standard Library - String/array/math/collections modules, JSON, file I/O
- **✅ Phase 12**: Enhanced Standard Library & Regular Expressions - Comprehensive regexp support
- **✅ Phase 13**: Performance & Compilation - Bytecode VM, ARM64 JIT compilation system

### Current Focus (January 2025)
- **🔧 Maintenance & Bug Fixes**: Stability improvements and issue resolution
  - Fixed critical VM stack underflow in nested if-else statements
  - Removed unused LLVM infrastructure and dependencies
  - Simplified CI/CD pipeline for better build reliability

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