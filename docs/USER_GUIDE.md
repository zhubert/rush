# Rush Language User Guide

Welcome to Rush, a modern, dynamically-typed programming language designed for simplicity and expressiveness.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Concepts](#basic-concepts)
3. [Working with Data](#working-with-data)
4. [Functions](#functions)
5. [Control Flow](#control-flow)
6. [Object-Oriented Programming](#object-oriented-programming)
7. [Error Handling](#error-handling)
8. [Module System](#module-system)
9. [Standard Library](#standard-library)
10. [Built-in Functions](#built-in-functions)
11. [Examples](#examples)
12. [REPL Usage](#repl-usage)
13. [Common Patterns](#common-patterns)
14. [Tips and Best Practices](#tips-and-best-practices)

## Getting Started

### Installation

1. Clone the Rush repository
2. Ensure Go is installed on your system
3. Navigate to the Rush directory

### Running Rush Programs

#### After Installation:
```bash
# Execute a file
rush your_program.rush

# Start interactive REPL
rush
```

#### Development Mode (from source):
```bash
# Execute a file
make dev FILE=your_program.rush

# Start REPL
make repl

# Build project
make build

# Install system-wide
make install
```

### Your First Program

Create a file called `hello.rush`:

```rush
print("Hello, Rush!")
```

Run it:
```bash
rush hello.rush
```

## Basic Concepts

### Variables

Variables in Rush are dynamically typed - you don't need to declare their type:

```rush
name = "Alice"        # String
age = 25             # Integer
height = 5.8         # Float
is_student = true    # Boolean
```

Variables can change types:

```rush
x = 42              # x is an integer
x = "now a string"  # x is now a string
x = [1, 2, 3]      # x is now an array
```

### Comments

Use `#` for comments:

```rush
# This is a single-line comment
x = 5  # Comment at end of line

# Comments can explain complex logic
# or provide documentation
```

## Working with Data

### Numbers

Rush supports integers and floats:

```rush
# Integers
count = 42
negative = -17

# Floats
pi = 3.14159
temperature = -0.5

# Arithmetic operations
sum = 10 + 5        # 15
difference = 10 - 3 # 7
product = 4 * 6     # 24
quotient = 15 / 3   # 5.0 (always float)
```

**Important**: Division always produces a float, even with integer operands.

### Strings

Work with text using strings:

```rush
# String literals
greeting = "Hello, World!"
name = "Rush"

# String concatenation
message = "Welcome to " + name + "!"

# String indexing (0-based)
first_char = "Hello"[0]  # "H"
last_char = "Hello"[4]   # "o"

# String length
length = len("Hello")    # 5
```

### Booleans

For true/false values:

```rush
is_ready = true
is_complete = false

# Boolean operations
can_proceed = is_ready && !is_complete
should_wait = !is_ready || is_complete
```

### Arrays

Store collections of values:

```rush
# Array creation
numbers = [1, 2, 3, 4, 5]
names = ["Alice", "Bob", "Charlie"]
mixed = [42, "hello", true, [1, 2]]

# Array indexing
first = numbers[0]      # 1
second = numbers[1]     # 2

# Array length
size = len(numbers)     # 5

# Nested arrays
matrix = [[1, 2], [3, 4], [5, 6]]
element = matrix[1][0]  # 3
```

## Functions

Functions are first-class values in Rush:

### Basic Function Definition

```rush
# Simple function
add = fn(x, y) {
  return x + y
}

result = add(5, 3)  # 8
```

### Functions with Logic

```rush
# Function with conditional logic
max = fn(a, b) {
  if (a > b) {
    return a
  } else {
    return b
  }
}

largest = max(10, 15)  # 15
```

### Recursive Functions

```rush
# Factorial function
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}

result = factorial(5)  # 120
```

### Higher-Order Functions

Functions can take other functions as arguments:

```rush
# Function that applies another function twice
apply_twice = fn(f, x) {
  return f(f(x))
}

square = fn(n) { return n * n }
result = apply_twice(square, 3)  # 81 (3² = 9, 9² = 81)
```

### Anonymous Functions

```rush
# Functions can be created inline
numbers = [1, 2, 3, 4, 5]
transform = fn(arr, func) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    result = push(result, func(arr[i]))
  }
  return result
}

squared = transform(numbers, fn(x) { return x * x })
```

## Control Flow

### If-Else Statements

```rush
age = 18

if (age >= 18) {
  print("You can vote!")
} else {
  print("You cannot vote yet.")
}
```

### If Expressions

If statements can return values:

```rush
status = if (age >= 18) { "adult" } else { "minor" }
```

### While Loops

```rush
count = 1
while (count <= 5) {
  print("Count: " + type(count))
  count = count + 1
}
```

### For Loops

```rush
# Traditional for loop
for (i = 0; i < 10; i = i + 1) {
  print("Number: " + type(i))
}

# Iterating over arrays
numbers = [1, 2, 3, 4, 5]
sum = 0
for (i = 0; i < len(numbers); i = i + 1) {
  sum = sum + numbers[i]
}
```

### Nested Loops

```rush
# Working with 2D arrays
matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]

for (row = 0; row < len(matrix); row = row + 1) {
  for (col = 0; col < len(matrix[row]); col = col + 1) {
    element = matrix[row][col]
    print("Element at [" + type(row) + "," + type(col) + "]: " + type(element))
  }
}
```

## Object-Oriented Programming

Rush supports object-oriented programming with classes, inheritance, and method calls.

### Class Definition

```rush
# Define a class
class Animal {
  fn initialize(name) {
    @name = name  # Instance variable with @ prefix
  }
  
  fn speak() {
    return @name + " makes a sound"
  }
  
  fn get_name() {
    return @name
  }
}

# Create an instance
animal = Animal.new("Generic Animal")
print(animal.speak())  # "Generic Animal makes a sound"
```

### Inheritance

```rush
# Child class inherits from parent
class Dog < Animal {
  fn initialize(name, breed) {
    super(name)  # Call parent constructor
    @breed = breed
  }
  
  fn speak() {
    return @name + " barks"
  }
  
  fn get_breed() {
    return @breed
  }
}

# Create a dog instance
dog = Dog.new("Buddy", "Golden Retriever")
print(dog.speak())        # "Buddy barks"
print(dog.get_name())     # "Buddy"
print(dog.get_breed())    # "Golden Retriever"
```

### Method Calls

```rush
# Methods are called using dot notation
result = object.method_name(arguments)

# Instance variables use @ prefix
class Counter {
  fn initialize() {
    @count = 0
  }
  
  fn increment() {
    @count = @count + 1
    return @count
  }
  
  fn get_count() {
    return @count
  }
}

counter = Counter.new()
print(counter.increment())  # 1
print(counter.increment())  # 2
print(counter.get_count())  # 2
```

## Error Handling

Rush provides comprehensive error handling with try/catch/finally blocks.

### Basic Try-Catch

```rush
try {
  result = risky_operation()
  print("Success:", result)
} catch (error) {
  print("Something went wrong:", error.message)
}
```

### Try-Catch-Finally

The finally block always executes:

```rush
file = null
try {
  file = open_file("data.txt")
  content = read_file(file)
  process_content(content)
} catch (error) {
  print("Error processing file:", error.message)
} finally {
  if (file != null) {
    close_file(file)
  }
  print("Cleanup completed")
}
```

### Multiple Catch Blocks

Handle different error types separately:

```rush
try {
  data = validate_and_process(input)
} catch (ValidationError error) {
  print("Validation failed:", error.message)
  show_help()
} catch (NetworkError error) {
  print("Network issue:", error.message)
  retry_connection()
} catch (error) {
  print("Unexpected error:", error.message)
  log_error(error)
}
```

### Throwing Errors

```rush
validate_age = fn(age) {
  if (age < 0) {
    throw ValidationError("Age cannot be negative")
  }
  if (age > 150) {
    throw ValidationError("Age seems unrealistic")
  }
  return age
}

try {
  age = validate_age(-5)
} catch (ValidationError error) {
  print("Invalid age:", error.message)
}
```

### Error Types

Built-in error types include:
- `Error` - Base error type
- `ValidationError` - Input validation failures
- `TypeError` - Type-related errors
- `IndexError` - Array/string index out of bounds
- `ArgumentError` - Function argument errors
- `RuntimeError` - General runtime errors

## Standard Library

Rush includes a comprehensive standard library with modules for common operations.

### Math Module (`std/math`)

Mathematical functions and constants:

```rush
import { PI, E, sqrt, abs, min, max, pow } from "std/math"

# Constants
print("π =", PI)  # 3.141592653589793
print("e =", E)   # 2.718281828459045

# Basic operations
print(abs(-42))           # 42
print(min(5, 3, 8, 1))   # 1
print(max(5, 3, 8, 1))   # 8
print(sqrt(16))          # 4.0
print(pow(2, 8))         # 256.0

# Rounding functions
import { floor, ceil, round } from "std/math"
print(floor(3.7))        # 3.0
print(ceil(3.2))         # 4.0
print(round(3.6))        # 4.0

# Random functions
import { random, random_int } from "std/math"
print(random())          # Random float [0, 1)
print(random_int(1, 10)) # Random integer [1, 10]

# Array operations
import { sum, average } from "std/math"
numbers = [1, 2, 3, 4, 5]
print(sum(numbers))      # 15.0
print(average(numbers))  # 3.0

# Type predicates
import { is_number?, is_integer? } from "std/math"
print(is_number?(42))    # true
print(is_integer?(3.14)) # false
```

### String Module (`std/string`)

String manipulation functions:

```rush
import { trim, upper, lower, contains? } from "std/string"

text = "  Hello, World!  "
print(trim(text))        # "Hello, World!"
print(upper(text))       # "  HELLO, WORLD!  "
print(lower(text))       # "  hello, world!  "
print(contains?(text, "World"))  # true

# More string functions
import { replace, starts_with?, ends_with?, join } from "std/string"

sentence = "Hello, beautiful world!"
print(replace(sentence, "world", "universe"))  # "Hello, beautiful universe!"
print(starts_with?(sentence, "Hello"))         # true
print(ends_with?(sentence, "world!"))          # true

words = ["apple", "banana", "cherry"]
print(join(words, ", "))  # "apple, banana, cherry"

# Character operations
import { is_whitespace_char? } from "std/string"
print(is_whitespace_char?(" "))   # true
print(is_whitespace_char?("a"))   # false
```

### Array Module (`std/array`)

Functional programming utilities for arrays:

```rush
import { map, filter, reduce, find } from "std/array"

numbers = [1, 2, 3, 4, 5]

# Transform each element
doubled = map(numbers, fn(x) { x * 2 })
print(doubled)  # [2, 4, 6, 8, 10]

# Filter elements
evens = filter(numbers, fn(x) { x % 2 == 0 })
print(evens)    # [2, 4]

# Reduce to single value
sum = reduce(numbers, fn(acc, x) { acc + x }, 0)
print(sum)      # 15

# Find first matching element
first_even = find(numbers, fn(x) { x % 2 == 0 })
print(first_even)  # 2

# More array utilities
import { index_of, includes?, reverse, sort } from "std/array"

fruits = ["apple", "banana", "cherry"]
print(index_of(fruits, "banana"))     # 1
print(includes?(fruits, "grape"))     # false
print(reverse(fruits))                # ["cherry", "banana", "apple"]

numbers = [3, 1, 4, 1, 5]
print(sort(numbers))                  # [1, 1, 3, 4, 5]
```

### Import Aliasing

Use aliases for cleaner code:

```rush
import { 
  very_long_function_name as short,
  another_long_name as brief
} from "std/string"

result = short("some text")
```

## Module System

The module system helps you organize code into reusable files. You can export values from one file and import them in another.

### Creating a Module

Create a module by using `export` statements:

**math.rush:**
```rush
# Export functions
export add = fn(x, y) {
  return x + y
}

export subtract = fn(x, y) {
  return x - y
}

# Export constants
export PI = 3.14159

# Private helper function (not exported)
helper = fn(x) {
  return x * 2
}

export square = fn(x) {
  return helper(x) * helper(x) / 4  # Uses private helper
}
```

### Using a Module

Import values from a module using `import` statements:

**calculator.rush:**
```rush
# Import specific functions and constants
import { add, subtract, PI } from "./math"

# Use imported functions
result1 = add(10, 5)           # 15
result2 = subtract(10, 5)      # 5

# Use imported constants
circumference = 2 * PI * 5     # ~31.416

print("Addition:", result1)
print("Subtraction:", result2)
print("Circumference:", circumference)
```

### Module Paths

Rush supports different types of module paths:

#### Relative Paths
```rush
import { func } from "./same_directory"      # Same directory
import { func } from "../parent_directory"   # Parent directory
import { func } from "./sub/nested"         # Subdirectory
```

#### Module File Extensions
The `.rush` extension is added automatically:
```rush
import { func } from "./math"        # Loads math.rush
import { func } from "./math.rush"   # Also loads math.rush
```

### Best Practices

1. **Keep modules focused**: Each module should have a single, clear purpose
2. **Export only what's needed**: Keep internal helpers private
3. **Use descriptive names**: Make module and export names clear
4. **Organize by functionality**: Group related functions together
5. **Avoid circular dependencies**: Don't have modules import each other

## Built-in Functions

### I/O Functions

#### `print(...)`
Output values to the console:

```rush
print("Hello")           # String
print(42)               # Number
print(true)             # Boolean
print([1, 2, 3])        # Array
print("Value:", 42)     # Multiple arguments
```

### Utility Functions

#### `len(collection)`
Get the length of arrays or strings:

```rush
array_length = len([1, 2, 3, 4])     # 4
string_length = len("hello")          # 5
```

#### `type(value)`
Get the type of any value:

```rush
print(type(42))          # "INTEGER"
print(type(3.14))        # "FLOAT"
print(type("hello"))     # "STRING"
print(type(true))        # "BOOLEAN"
print(type([1, 2]))      # "ARRAY"
print(type(fn() {}))     # "FUNCTION"
```

### String Functions

#### `substr(string, start, length)`
Extract part of a string:

```rush
text = "Hello, World!"
hello = substr(text, 0, 5)     # "Hello"
world = substr(text, 7, 5)     # "World"
```

#### `split(string, separator)`
Split a string into an array:

```rush
csv = "apple,banana,cherry"
fruits = split(csv, ",")        # ["apple", "banana", "cherry"]

sentence = "Hello World"
words = split(sentence, " ")    # ["Hello", "World"]

# Split into characters
chars = split("hello", "")      # ["h", "e", "l", "l", "o"]
```

### Array Functions

#### `push(array, element)`
Add an element to an array (returns new array):

```rush
numbers = [1, 2, 3]
extended = push(numbers, 4)     # [1, 2, 3, 4]
```

#### `pop(array)`
Get the last element of an array:

```rush
numbers = [1, 2, 3, 4, 5]
last = pop(numbers)             # 5
```

#### `slice(array, start, end)`
Extract a portion of an array:

```rush
numbers = [1, 2, 3, 4, 5]
middle = slice(numbers, 1, 4)   # [2, 3, 4]
```

## Examples

### Working with Data

```rush
# Student grade calculator
students = ["Alice", "Bob", "Charlie"]
grades = [85, 92, 78]

total = 0
for (i = 0; i < len(grades); i = i + 1) {
  print(students[i] + ": " + type(grades[i]))
  total = total + grades[i]
}

average = total / len(grades)
print("Class average: " + type(average))
```

### Text Processing

```rush
# Word counter
text = "The quick brown fox jumps over the lazy dog"
words = split(text, " ")

print("Text: " + text)
print("Word count: " + type(len(words)))

# Find longest word
longest = ""
for (i = 0; i < len(words); i = i + 1) {
  word = words[i]
  if (len(word) > len(longest)) {
    longest = word
  }
}

print("Longest word: " + longest + " (" + type(len(longest)) + " characters)")
```

### Algorithm Implementation

```rush
# Bubble sort
sort_array = fn(arr) {
  n = len(arr)
  # Create a copy
  result = []
  for (i = 0; i < n; i = i + 1) {
    result = push(result, arr[i])
  }
  
  # Sort
  for (i = 0; i < n - 1; i = i + 1) {
    for (j = 0; j < n - i - 1; j = j + 1) {
      if (result[j] > result[j + 1]) {
        # Swap
        temp = result[j]
        result[j] = result[j + 1]
        result[j + 1] = temp
      }
    }
  }
  
  return result
}

unsorted = [64, 34, 25, 12, 22, 11, 90]
sorted_array = sort_array(unsorted)

print("Original: " + type(len(unsorted)) + " elements")
print("Sorted: " + type(len(sorted_array)) + " elements")
```

## REPL Usage

The Rush REPL (Read-Eval-Print Loop) provides an interactive environment:

### Starting the REPL

```bash
go run cmd/rush/main.go
```

### REPL Commands

- `:help` - Show help information
- `:quit` - Exit the REPL

### Interactive Examples

```
Rush Interactive REPL
Type ':help' for help, ':quit' to exit
⛤ x = 10
⛤ y = 20
⛤ x + y
30
⛤ greet = fn(name) { return "Hello, " + name }
⛤ greet("World")
Hello, World
⛤ :quit
```

### Multi-line Input

The REPL supports multi-line expressions:

```
⛤ factorial = fn(n) {
⛤   if (n <= 1) {
⛤     return 1
⛤   } else {
⛤     return n * factorial(n - 1)
⛤   }
⛤ }
⛤ factorial(5)
120
```

## Common Patterns

### Array Processing

```rush
# Map operation (transform each element)
map = fn(arr, func) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    result = push(result, func(arr[i]))
  }
  return result
}

numbers = [1, 2, 3, 4, 5]
doubled = map(numbers, fn(x) { return x * 2 })
```

### Filter Operation

```rush
# Filter operation (select elements that match condition)
filter = fn(arr, predicate) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    if (predicate(arr[i])) {
      result = push(result, arr[i])
    }
  }
  return result
}

numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
evens = filter(numbers, fn(x) { 
  doubled = x * 2
  return doubled / 2 == x / 2 * 2  # Even number check
})
```

### Reduce Operation

```rush
# Reduce operation (combine elements into single value)
reduce = fn(arr, func, initial) {
  accumulator = initial
  for (i = 0; i < len(arr); i = i + 1) {
    accumulator = func(accumulator, arr[i])
  }
  return accumulator
}

numbers = [1, 2, 3, 4, 5]
sum = reduce(numbers, fn(acc, x) { return acc + x }, 0)
product = reduce(numbers, fn(acc, x) { return acc * x }, 1)
```

### Object-like Patterns

```rush
# Using arrays to simulate objects
create_person = fn(name, age) {
  return [name, age]  # [0] = name, [1] = age
}

get_name = fn(person) { return person[0] }
get_age = fn(person) { return person[1] }

alice = create_person("Alice", 25)
print("Name: " + get_name(alice))
print("Age: " + type(get_age(alice)))
```

## Tips and Best Practices

### Code Organization

1. **Use descriptive names**:
   ```rush
   # Good
   calculate_average = fn(grades) { ... }
   
   # Less clear
   calc = fn(g) { ... }
   ```

2. **Keep functions small and focused**:
   ```rush
   # Good - single responsibility
   is_even = fn(n) { return (n / 2) * 2 == n }
   filter_evens = fn(arr) { return filter(arr, is_even) }
   ```

3. **Use meaningful variable names**:
   ```rush
   # Good
   user_count = len(users)
   max_attempts = 3
   
   # Less clear
   uc = len(u)
   ma = 3
   ```

### Performance Considerations

1. **Avoid unnecessary array copying**:
   ```rush
   # Less efficient - creates many intermediate arrays
   result = []
   for (i = 0; i < 1000; i = i + 1) {
     result = push(result, i)
   }
   
   # More efficient - build once
   build_range = fn(n) {
     result = []
     for (i = 0; i < n; i = i + 1) {
       result = push(result, i)
     }
     return result
   }
   ```

2. **Use early returns in functions**:
   ```rush
   # Good
   find_item = fn(arr, target) {
     for (i = 0; i < len(arr); i = i + 1) {
       if (arr[i] == target) {
         return i  # Early return
       }
     }
     return -1
   }
   ```

### Error Handling

Since Rush doesn't have exceptions, use return values to indicate errors:

```rush
safe_divide = fn(a, b) {
  if (b == 0) {
    return [false, 0]  # [success, result]
  }
  return [true, a / b]
}

result = safe_divide(10, 2)
if (result[0]) {
  print("Result: " + type(result[1]))
} else {
  print("Error: Division by zero")
}
```

### Type Checking

Use the `type()` function for runtime type checking:

```rush
process_value = fn(value) {
  value_type = type(value)
  
  if (value_type == "INTEGER" || value_type == "FLOAT") {
    return value * 2
  } else {
    if (value_type == "STRING") {
      return value + value
    } else {
      return value
    }
  }
}
```

## Next Steps

1. **Explore the examples**: Check out `examples/` directory for more complex programs
2. **Read the language specification**: See `docs/LANGUAGE_SPECIFICATION.md` for detailed syntax
3. **Practice in the REPL**: Experiment with Rush interactively
4. **Build projects**: Create your own Rush programs

Happy coding with Rush! 🚀