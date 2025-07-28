# Rush Language Specification

## Overview

Rush is a dynamically-typed, interpreted programming language with clean syntax inspired by modern languages. It features C-style control flow, first-class functions, and built-in support for arrays, hashes/dictionaries, and strings.

## Table of Contents

1. [Execution Architecture](#execution-architecture)
2. [Lexical Structure](#lexical-structure)
3. [Data Types](#data-types)
4. [Variables](#variables)
5. [Expressions](#expressions)
6. [Statements](#statements)
7. [Functions](#functions)
8. [Control Flow](#control-flow)
9. [Error Handling](#error-handling)
10. [Module System](#module-system)
11. [Built-in Functions](#built-in-functions)
12. [Grammar](#grammar)

## Execution Architecture

Rush implements a tiered execution system with three distinct modes:

### 1. Tree-Walking Interpreter (Level 0)

The default execution mode performs direct evaluation of the Abstract Syntax Tree (AST):

```
Source Code → Lexer → Parser → AST → Interpreter → Result
```

**Characteristics:**
- Direct AST traversal and evaluation
- No compilation step
- Immediate execution
- Best for development and debugging
- Slower performance for compute-intensive tasks

### 2. Bytecode Virtual Machine (Level 1)

The bytecode mode compiles Rush source to intermediate bytecode instructions:

```
Source Code → Lexer → Parser → AST → Compiler → Bytecode → VM → Result
```

**Characteristics:**
- Compilation to stack-based bytecode
- Virtual machine execution
- Better performance than interpreter
- Comprehensive logging and debugging support
- Good balance of speed and portability

**Bytecode Instructions Include:**
- `OpConstant` - Load constant value
- `OpAdd`, `OpSub`, `OpMul`, `OpDiv` - Arithmetic operations
- `OpEqual`, `OpNotEqual`, `OpGreaterThan` - Comparison operations
- `OpJump`, `OpJumpNotTruthy` - Control flow
- `OpCall`, `OpReturn` - Function operations
- `OpGetGlobal`, `OpSetGlobal` - Variable access

### 3. JIT Compilation (Level 2)

The JIT mode provides adaptive optimization with native code generation:

```
Source Code → Lexer → Parser → AST → Compiler → Bytecode → VM → JIT → ARM64 → Result
                                                            ↑
                                                      (Hot Functions)
```

**Characteristics:**
- ARM64 architecture support
- Hot path detection (threshold: 100+ function calls)
- Native machine code generation
- Automatic deoptimization and fallback
- Up to 15x performance improvement for hot functions
- Code cache with LRU eviction

**JIT Process:**
1. Functions start in bytecode VM
2. Execution count tracked per function
3. Hot functions (100+ calls) trigger JIT compilation
4. ARM64 code generated and cached
5. Subsequent calls execute native code
6. Fallback to bytecode VM on compilation failure

### Command Line Interface

```bash
# Tree-walking interpreter (default)
rush program.rush

# Bytecode virtual machine
rush -bytecode program.rush

# JIT compilation (ARM64 only)
rush -jit program.rush

# Performance monitoring
rush -bytecode -log-level=info program.rush
rush -jit -log-level=info program.rush
```

## Lexical Structure

### Comments

Single-line comments start with `#` and continue to the end of the line:

```rush
# This is a comment
x = 5  # This is also a comment
```

### Identifiers

Identifiers start with a letter or underscore, followed by letters, digits, or underscores:

```rush
variable_name
_private
userName
count2
```

### Keywords

Reserved words in Rush:

- `fn` - function definition
- `if` - conditional statement
- `else` - alternative branch
- `while` - while loop
- `for` - for loop
- `return` - return statement
- `import` - import statement
- `export` - export statement  
- `from` - from clause in imports
- `as` - alias clause in imports
- `try` - try block for error handling
- `catch` - catch block for error handling
- `finally` - finally block for cleanup
- `throw` - throw statement for raising errors
- `switch` - switch statement
- `case` - case clause in switch
- `default` - default clause in switch
- `break` - break statement
- `continue` - continue statement
- `class` - class definition
- `new` - object instantiation
- `super` - parent class reference
- `true` - boolean literal
- `false` - boolean literal
- `null` - null literal

### Operators

#### Arithmetic Operators
- `+` - addition
- `-` - subtraction
- `*` - multiplication
- `/` - division (always produces float result)
- `%` - modulo (remainder)

#### Comparison Operators
- `==` - equality
- `!=` - inequality
- `<` - less than
- `>` - greater than
- `<=` - less than or equal
- `>=` - greater than or equal

#### Logical Operators
- `&&` - logical AND (short-circuit)
- `||` - logical OR (short-circuit)
- `!` - logical NOT

#### Assignment Operator
- `=` - assignment

#### Access Operators
- `.` - module member access / object method call
- `[index]` - array/string indexing
- `[index] =` - array element assignment

### Delimiters
- `()` - parentheses for grouping and function calls
- `{}` - braces for code blocks and hash literals
- `[]` - brackets for array literals and indexing
- `,` - comma separator
- `:` - colon for hash key-value pairs
- `;` - statement terminator (optional, newlines also terminate statements)

## Data Types

### Integer

Whole numbers:

```rush
x = 42
y = -17
z = 0
```

### Float

Decimal numbers:

```rush
pi = 3.14159
temperature = -0.5
rate = 2.0
```

**Note**: Division operations always produce float results, even with integer operands.

### String

Text enclosed in double quotes:

```rush
name = "Rush"
message = "Hello, World!"
empty = ""
```

String indexing is supported:

```rush
first_char = "Hello"[0]  # Returns "H"
```

### Boolean

Logical values:

```rush
is_ready = true
is_complete = false
```

### Array

Ordered collections of values:

```rush
numbers = [1, 2, 3, 4, 5]
mixed = [42, "hello", true, [1, 2]]
empty_array = []
```

Array indexing starts at 0:

```rush
first = numbers[0]   # Returns 1
second = numbers[1]  # Returns 2
```

Out-of-bounds access returns `null`.

Array element assignment is supported:

```rush
numbers = [1, 2, 3]
numbers[0] = 10  # Array becomes [10, 2, 3]
```

### Hash/Dictionary

Key-value mappings with `{key: value}` syntax:

```rush
person = {"name": "Alice", "age": 30, "active": true}
config = {42: "answer", "pi": 3.14, true: "enabled"}
empty_hash = {}
```

Hash keys can be strings, integers, booleans, or floats:

```rush
mixed_keys = {
  "string": "value1",
  42: "value2", 
  true: "value3",
  3.14: "value4"
}
```

Hash indexing uses bracket notation:

```rush
name = person["name"]      # Returns "Alice"
age = person["age"]        # Returns 30
```

Hash assignment creates or updates key-value pairs:

```rush
person["city"] = "NYC"     # Add new key
person["age"] = 31         # Update existing key
```

Out-of-bounds access returns `null`:

```rush
unknown = person["unknown"]  # Returns null
```

### Function

First-class functions with closure support:

```rush
add = fn(x, y) { return x + y }
greet = fn(name) { return "Hello, " + name }

# Closures capture their environment
createCounter = fn() {
  count = 0
  return fn() {
    count = count + 1
    return count
  }
}

counter = createCounter()
print(counter())  # 1
print(counter())  # 2
```

### Error

Error objects represent exceptional conditions:

```rush
error = ValidationError("Invalid input")
```

Error objects have properties:
- `type` - the error type name (string)
- `message` - descriptive error message (string)
- `stack` - stack trace information (string)
- `line` - line number where error occurred (integer)
- `column` - column number where error occurred (integer)

### Null

Represents absence of value. Returned by:
- Out-of-bounds array access
- Functions without explicit return
- Failed operations

```rush
result = null
if (result == null) {
  print("No value")
}
```

### Class/Object

Object-oriented programming support:

```rush
class Animal {
  fn initialize(name) {
    @name = name  # Instance variable
  }
  
  fn speak() {
    return @name + " makes a sound"
  }
}

class Dog < Animal {  # Inheritance
  fn speak() {
    return @name + " barks"
  }
}

dog = Dog.new("Buddy")
print(dog.speak())  # "Buddy barks"
```

### File System Types

Rush provides built-in file system types for file and directory operations with dot notation support.

#### File

Represents a file in the file system:

```rush
# Create a file object
text_file = file("data.txt")
print(text_file.path)     # "data.txt"
print(text_file.is_open)  # false

# File operations
text_file.open("w").write("Hello, Rush!").close()
if text_file.exists?() {
    content = text_file.open("r").read()
    print(content)  # "Hello, Rush!"
}
```

#### Directory

Represents a directory in the file system:

```rush
# Create a directory object
project_dir = directory("my_project")
print(project_dir.path)  # "my_project"

# Directory operations
project_dir.create()
if project_dir.exists?() {
    files = project_dir.list()
    print("Directory contains " + len(files) + " items")
}
```

#### Path

Represents a file system path with cross-platform manipulation:

```rush
# Create a path object
config_path = path("/home/user")
print(config_path.value)  # "/home/user"

# Path manipulation
full_path = config_path.join("config").join("app.conf")
print(full_path.basename())  # "app.conf"
print(full_path.dirname())   # "/home/user/config"

# Clean and make absolute
clean_path = path("./data/../config.txt").clean().absolute()
```

**Features:**
- **Dot notation support**: All file system types support method chaining
- **Cross-platform**: Path operations work consistently across operating systems
- **Security**: Path traversal attacks (`..`) are blocked in constructors
- **Error handling**: Operations return appropriate error messages for invalid operations

### JSON

JSON objects provide comprehensive JSON processing capabilities with object-oriented dot notation:

```rush
# Parse JSON strings into JSON objects
user_data = json_parse("{\"name\": \"Alice\", \"age\": 30}")
array_data = json_parse("[1, 2, 3, 4, 5]")

# Convert Rush values to JSON strings
json_string = json_stringify({"name": "Bob", "active": true})
```

JSON objects support comprehensive dot notation methods:

```rush
# Data access
user_data.get("name")              # "Alice"
user_data.has?("email")            # false
user_data.keys()                   # ["name", "age"]
user_data.length()                 # 2

# Data modification (returns new JSON objects)
updated = user_data.set("email", "alice@example.com")
                  .set("verified", true)

# Path-based navigation
nested = json_parse("{\"user\": {\"profile\": {\"city\": \"NYC\"}}}")
city = nested.path("user.profile.city")  # "NYC"

# Formatting
user_data.compact()                # Compact JSON string
user_data.pretty()                 # Pretty-formatted with indentation
user_data.pretty("    ")           # Custom indentation

# Advanced operations
config1 = json_parse("{\"debug\": true}")
config2 = json_parse("{\"port\": 8080}")
merged = config1.merge(config2)    # Combine JSON objects
```

**Features:**
- **Full JSON standard compliance**: Supports all JSON data types
- **Object-oriented API**: Comprehensive dot notation methods
- **Method chaining**: Fluent operations for complex transformations
- **Path navigation**: Dot-separated path access for nested data
- **Immutable operations**: Methods return new objects, preserving originals
- **Rich formatting**: Compact and pretty-print options with custom indentation
- **Error handling**: Robust parsing and serialization error reporting

## Variables

Variables are dynamically typed and declared by assignment:

```rush
x = 5          # Integer
x = "hello"    # Now a string
x = [1, 2, 3]  # Now an array
```

Variable names must start with a letter or underscore.

## Expressions

### Literals
```rush
42           # Integer literal
3.14         # Float literal
"hello"      # String literal
true         # Boolean literal
[1, 2, 3]    # Array literal
```

### Arithmetic Expressions
```rush
x + y        # Addition
x - y        # Subtraction
x * y        # Multiplication
x / y        # Division (always float)
x % y        # Modulo (remainder)
-x           # Unary negation
```

### Comparison Expressions
```rush
x == y       # Equality
x != y       # Inequality
x < y        # Less than
x > y        # Greater than
x <= y       # Less than or equal
x >= y       # Greater than or equal
```

### Logical Expressions
```rush
x && y       # Logical AND
x || y       # Logical OR
!x           # Logical NOT
```

### Function Call Expressions
```rush
add(5, 3)
greet("World")
len(array)
```

### Index Expressions
```rush
array[0]     # Array indexing
string[1]    # String indexing
matrix[i][j] # Nested indexing
array[i] = value  # Array element assignment
```

### Grouping
```rush
(x + y) * z  # Parentheses for precedence
```

## Statements

### Assignment Statement
```rush
variable = expression
```

### Expression Statement
```rush
function_call()
5 + 3  # Valid but result is discarded
```

### Block Statement
```rush
{
  statement1
  statement2
}
```

### Return Statement
```rush
return expression
return  # Returns null
```

### Throw Statement
```rush
throw ErrorType("Error message")
throw expression  # Any value can be thrown
```

### Try-Catch-Finally Statement
```rush
try {
  # Code that might throw an error
} catch (error) {
  # Handle error
}

try {
  # Code that might throw an error
} catch (ValidationError error) {
  # Handle specific error type
} catch (error) {
  # Handle any other error
} finally {
  # Always executed cleanup code
}
```

## Functions

### Function Definition
```rush
function_name = fn(param1, param2, ...) {
  # function body
  return value  # optional
}
```

### Function Calls
```rush
result = function_name(arg1, arg2, ...)
```

### Anonymous Functions
```rush
# Functions are values and can be used directly
apply = fn(f, x) { return f(x) }
result = apply(fn(n) { return n * n }, 5)
```

### Recursion
```rush
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
```

## Control Flow

### If-Else Statements
```rush
if (condition) {
  # if block
} else {
  # else block
}

# If expressions (return values)
result = if (x > 0) { "positive" } else { "non-positive" }
```

### While Loops
```rush
while (condition) {
  # loop body
}
```

### For Loops
```rush
for (initialization; condition; update) {
  # loop body
}

# Example
for (i = 0; i < 10; i = i + 1) {
  print(i)
}
```

### Switch Statements
```rush
switch (value) {
  case 1, 2, 3: {
    print("Small number")
  }
  case 4, 5: {
    print("Medium number")
  }
  default: {
    print("Other value")
  }
}
```

### Break and Continue
```rush
for (i = 0; i < 10; i = i + 1) {
  if (i == 3) {
    continue  # Skip iteration
  }
  if (i == 7) {
    break     # Exit loop
  }
  print(i)
}
```

## Error Handling

Rush provides comprehensive error handling through try/catch/finally blocks and throw statements.

### Basic Try-Catch

```rush
try {
  result = divide(10, 0)
  print("Result:", result)
} catch (error) {
  print("Error occurred:", error.message)
}
```

### Try-Catch-Finally

The `finally` block always executes, regardless of whether an error occurs:

```rush
file = null
try {
  file = open("data.txt")
  content = read(file)
  process(content)
} catch (error) {
  print("Failed to process file:", error.message)
} finally {
  if (file != null) {
    close(file)
  }
}
```

### Multiple Catch Blocks

Different error types can be handled separately:

```rush
try {
  data = validateAndProcess(input)
} catch (ValidationError error) {
  print("Validation failed:", error.message)
  showValidationHelp()
} catch (NetworkError error) {
  print("Network issue:", error.message)
  retryConnection()
} catch (error) {
  print("Unexpected error:", error.message)
  logError(error)
}
```

### Throwing Errors

Use `throw` to raise errors:

```rush
validateAge = fn(age) {
  if (age < 0) {
    throw ValidationError("Age cannot be negative")
  }
  if (age > 150) {
    throw ValidationError("Age seems unrealistic")
  }
  return age
}

processUser = fn(userData) {
  try {
    age = validateAge(userData.age)
    return createUser(userData)
  } catch (ValidationError error) {
    return { error: error.message }
  }
}
```

### Built-in Error Types

- `Error` - Base error type
- `ValidationError` - Input validation failures
- `TypeError` - Type-related errors
- `IndexError` - Array/string index out of bounds
- `ArgumentError` - Function argument errors
- `RuntimeError` - General runtime errors

### Custom Error Types

Create custom error types by calling them as constructors:

```rush
# Throwing custom errors
throw CustomError("Something went wrong")
throw DatabaseError("Connection failed")

# Catching custom errors
try {
  connectToDatabase()
} catch (DatabaseError error) {
  print("Database error:", error.message)
} catch (error) {
  print("Other error:", error.message)
}
```

### Error Object Properties

All error objects have these properties:

```rush
try {
  riskyOperation()
} catch (error) {
  print("Error type:", error.type)        # e.g., "ValidationError"
  print("Message:", error.message)        # e.g., "Invalid input"
  print("Stack trace:", error.stack)      # Function call stack
  print("Line:", error.line)              # Line number where error occurred
  print("Column:", error.column)          # Column number where error occurred
}
```

### Error Propagation

Errors automatically propagate up the call stack until caught:

```rush
deepFunction = fn() {
  throw RuntimeError("Deep error")
}

middleFunction = fn() {
  deepFunction()  # Error propagates through here
}

topFunction = fn() {
  try {
    middleFunction()
  } catch (error) {
    print("Caught error from deep function:", error.message)
  }
}
```

### Best Practices

1. **Use specific error types** for different failure modes
2. **Always handle errors** that might occur in your code
3. **Use finally blocks** for cleanup that must always happen
4. **Include descriptive error messages** for debugging
5. **Don't catch and ignore errors** without good reason

## Module System

The module system allows code organization through import and export statements.

### Export Statements

Export statements make values available to other modules:

```rush
# Export with assignment
export add = fn(x, y) { return x + y }
export PI = 3.14159

# Export existing variable
multiply = fn(x, y) { return x * y }
export multiply
```

### Import Statements

Import statements bring values from other modules into the current scope:

```rush
# Import specific values
import { add, PI } from "./math"

# Import multiple values
import { func1, func2, var1 } from "./utilities"

# Import with aliasing
import { add as plus, subtract as minus } from "./math"
import { very_long_function_name as short } from "./utils"
```

### Module Paths

Module paths can be:
- **Relative**: `./module` or `../parent/module`
- **Absolute**: `/path/to/module`
- **Standard Library**: `std/math`, `std/string`, `std/array`

The `.rush` extension is added automatically if not specified.

### Standard Library

Rush includes a comprehensive standard library:

```rush
# Math operations
import { sqrt, PI, abs, min, max } from "std/math"
radius = 5
area = PI * sqrt(radius)

# String manipulation
import { trim, upper, split, join } from "std/string"
text = upper(trim("  hello world  "))
words = split(text, " ")

# Array utilities
import { map, filter, reduce } from "std/array"
doubled = map([1, 2, 3], fn(x) { x * 2 })
evens = filter([1, 2, 3, 4], fn(x) { x % 2 == 0 })
sum = reduce([1, 2, 3], fn(acc, x) { acc + x }, 0)
```

### Module Example

**math.rush:**
```rush
export add = fn(x, y) { return x + y }
export subtract = fn(x, y) { return x - y }
export PI = 3.14159
```

**main.rush:**
```rush
import { add, subtract, PI } from "./math"

result = add(10, 5)        # 15
diff = subtract(10, 5)     # 5
area = PI * 4 * 4          # ~50.27
```

### Module Isolation

Each module executes in its own scope. Variables not exported remain private to the module.

## Built-in Functions

### `print(...)`
Outputs values to standard output:
```rush
print("Hello")
print(42)
print(array)
```

### `len(collection)`
Returns the length of arrays or strings:
```rush
len([1, 2, 3])    # Returns 3
len("hello")      # Returns 5
```

### `type(value)`
Returns the type of a value as a string:
```rush
type(42)          # Returns "INTEGER"
type(3.14)        # Returns "FLOAT"
type("hello")     # Returns "STRING"
type(true)        # Returns "BOOLEAN"
type([1, 2])      # Returns "ARRAY"
type(fn() {})     # Returns "FUNCTION"
type(null)        # Returns "NULL"
```

### `to_string(value)`
Converts any value to its string representation:
```rush
to_string(42)         # Returns "42"
to_string(3.14)       # Returns "3.14"
to_string(true)       # Returns "true"
to_string([1, 2, 3])  # Returns "[1, 2, 3]"
to_string(null)       # Returns "null"
```

### String Functions

#### `substr(string, start, length)`
Extracts a substring:
```rush
substr("hello", 1, 3)  # Returns "ell"
```

#### `split(string, separator)`
Splits a string into an array:
```rush
split("a,b,c", ",")    # Returns ["a", "b", "c"]
split("hello", "")     # Returns ["h", "e", "l", "l", "o"]
```

#### `ord(character)` and `chr(code)`
Convert between characters and ASCII codes:
```rush
ord("A")              # Returns 65
chr(65)               # Returns "A"
```

### Array Functions

#### `push(array, element)`
Returns a new array with the element added:
```rush
push([1, 2], 3)        # Returns [1, 2, 3]
```

#### `pop(array)`
Returns the last element of an array:
```rush
pop([1, 2, 3])         # Returns 3
pop([])                # Returns null
```

#### `slice(array, start, end)`
Returns a portion of an array:
```rush
slice([1, 2, 3, 4, 5], 1, 4)  # Returns [2, 3, 4]
```

### Hash Functions

#### `builtin_hash_keys(hash)`
Returns an array of all keys in the hash:
```rush
person = {"name": "Alice", "age": 30}
keys = builtin_hash_keys(person)  # Returns ["name", "age"]
```

#### `builtin_hash_values(hash)`
Returns an array of all values in the hash:
```rush
person = {"name": "Alice", "age": 30}
values = builtin_hash_values(person)  # Returns ["Alice", 30]
```

#### `builtin_hash_has_key(hash, key)`
Returns true if the hash contains the specified key:
```rush
person = {"name": "Alice", "age": 30}
has_name = builtin_hash_has_key(person, "name")  # Returns true
has_city = builtin_hash_has_key(person, "city")  # Returns false
```

#### `builtin_hash_get(hash, key, default?)`
Gets a value from the hash with optional default:
```rush
person = {"name": "Alice", "age": 30}
name = builtin_hash_get(person, "name")        # Returns "Alice"
city = builtin_hash_get(person, "city")        # Returns null
city = builtin_hash_get(person, "city", "NYC") # Returns "NYC"
```

#### `builtin_hash_set(hash, key, value)`
Returns a new hash with the key-value pair added or updated:
```rush
person = {"name": "Alice"}
updated = builtin_hash_set(person, "age", 30)  # Returns {"name": "Alice", "age": 30}
```

#### `builtin_hash_delete(hash, key)`
Returns a new hash with the specified key removed:
```rush
person = {"name": "Alice", "age": 30}
updated = builtin_hash_delete(person, "age")  # Returns {"name": "Alice"}
```

#### `builtin_hash_merge(hash1, hash2)`
Returns a new hash combining both hashes (hash2 values overwrite hash1):
```rush
person = {"name": "Alice", "age": 25}
updates = {"age": 30, "city": "NYC"}
merged = builtin_hash_merge(person, updates)  # Returns {"name": "Alice", "age": 30, "city": "NYC"}
```

## Grammar

### Precedence (highest to lowest)

1. Primary expressions: literals, identifiers, parentheses
2. Postfix: function calls, array indexing, method calls
3. Unary: `-`, `!`
4. Multiplicative: `*`, `/`, `%`
5. Additive: `+`, `-`
6. Relational: `<`, `>`, `<=`, `>=`
7. Equality: `==`, `!=`
8. Logical AND: `&&`
9. Logical OR: `||`
10. Assignment: `=`

### Associativity

- Most operators are left-associative
- Assignment is right-associative

### EBNF Grammar

```ebnf
program = { statement } ;

statement = assignmentStatement
          | expressionStatement
          | blockStatement
          | ifStatement
          | whileStatement
          | forStatement
          | returnStatement
          | throwStatement
          | tryStatement ;

assignmentStatement = identifier "=" expression ;

expressionStatement = expression ;

blockStatement = "{" { statement } "}" ;

ifStatement = "if" "(" expression ")" blockStatement [ "else" blockStatement ] ;

whileStatement = "while" "(" expression ")" blockStatement ;

forStatement = "for" "(" assignmentStatement ";" expression ";" assignmentStatement ")" blockStatement ;

returnStatement = "return" [ expression ] ;

throwStatement = "throw" expression ;

tryStatement = "try" blockStatement { catchClause } [ finallyClause ] ;

catchClause = "catch" "(" [ identifier ] identifier ")" blockStatement ;

finallyClause = "finally" blockStatement ;

expression = logicalOrExpression ;

logicalOrExpression = logicalAndExpression { "||" logicalAndExpression } ;

logicalAndExpression = equalityExpression { "&&" equalityExpression } ;

equalityExpression = relationalExpression { ( "==" | "!=" ) relationalExpression } ;

relationalExpression = additiveExpression { ( "<" | ">" | "<=" | ">=" ) additiveExpression } ;

additiveExpression = multiplicativeExpression { ( "+" | "-" ) multiplicativeExpression } ;

multiplicativeExpression = unaryExpression { ( "*" | "/" ) unaryExpression } ;

unaryExpression = ( "-" | "!" ) unaryExpression | postfixExpression ;

postfixExpression = primaryExpression { ( "(" [ argumentList ] ")" | "[" expression "]" ) } ;

primaryExpression = identifier
                  | integerLiteral
                  | floatLiteral
                  | stringLiteral
                  | booleanLiteral
                  | arrayLiteral
                  | functionLiteral
                  | ifExpression
                  | "(" expression ")" ;

arrayLiteral = "[" [ expressionList ] "]" ;

functionLiteral = "fn" "(" [ parameterList ] ")" blockStatement ;

ifExpression = "if" "(" expression ")" blockStatement "else" blockStatement ;

argumentList = expression { "," expression } ;

expressionList = expression { "," expression } ;

parameterList = identifier { "," identifier } ;

identifier = letter { letter | digit | "_" } ;

integerLiteral = digit { digit } ;

floatLiteral = digit { digit } "." digit { digit } ;

stringLiteral = '"' { character } '"' ;

booleanLiteral = "true" | "false" ;

letter = "a" | "b" | ... | "z" | "A" | "B" | ... | "Z" | "_" ;

digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;

character = any Unicode character except '"' ;
```

## Examples

### Hello World
```rush
print("Hello, World!")
```

### Variables and Arithmetic
```rush
x = 10
y = 20
sum = x + y
print("Sum: " + type(sum))
```

### Functions
```rush
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}

result = factorial(5)
print(result)  # Outputs: 120
```

### Arrays
```rush
numbers = [1, 2, 3, 4, 5]
sum = 0
for (i = 0; i < len(numbers); i = i + 1) {
  sum = sum + numbers[i]
}
print("Sum: " + type(sum))
```

### String Processing
```rush
text = "Hello, World!"
words = split(text, ", ")
for (i = 0; i < len(words); i = i + 1) {
  print("Word " + type(i) + ": " + words[i])
}
```