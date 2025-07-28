# Rush Language User Guide

Welcome to Rush, a modern, dynamically-typed programming language designed for simplicity and expressiveness.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Performance Modes](#performance-modes)
3. [Basic Concepts](#basic-concepts)
4. [Working with Data](#working-with-data)
5. [Functions](#functions)
6. [Control Flow](#control-flow)
7. [Object-Oriented Programming](#object-oriented-programming)
8. [Error Handling](#error-handling)
9. [Module System](#module-system)
10. [JSON Processing](#json-processing)
11. [File System Operations](#file-system-operations)
12. [Standard Library](#standard-library)
13. [Built-in Functions](#built-in-functions)
14. [Examples](#examples)
15. [REPL Usage](#repl-usage)
16. [Common Patterns](#common-patterns)
17. [Tips and Best Practices](#tips-and-best-practices)

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

# Execute with bytecode VM for better performance
rush -bytecode your_program.rush

# Execute with JIT compilation for maximum performance (ARM64 only)
rush -jit your_program.rush

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

## Performance Modes

Rush offers three execution modes, each optimized for different use cases:

### Tree-Walking Interpreter (Default)

The default execution mode uses a tree-walking interpreter that directly evaluates the Abstract Syntax Tree (AST).

```bash
# Execute with interpreter (default)
rush program.rush
```

**Best for:**
- Development and debugging
- Short-running scripts
- Interactive REPL sessions

**Characteristics:**
- Immediate execution
- Full debugging support
- Lowest memory usage
- Slower execution for compute-intensive tasks

### Bytecode Virtual Machine

The bytecode VM compiles Rush code to bytecode instructions and executes them on a virtual machine.

```bash
# Execute with bytecode VM
rush -bytecode program.rush

# With performance monitoring
rush -bytecode -log-level=info program.rush
```

**Best for:**
- Production applications
- Long-running programs
- Better performance than interpreter

**Characteristics:**
- Compilation to bytecode
- Faster execution than interpreter
- Performance monitoring available
- Good balance of speed and compatibility

### JIT Compilation (ARM64)

The Just-In-Time compiler provides the highest performance by compiling hot functions to native ARM64 machine code.

```bash
# Execute with JIT compilation
rush -jit program.rush

# With JIT statistics
rush -jit -log-level=info program.rush
```

**Best for:**
- Compute-intensive applications
- Functions called repeatedly
- Maximum performance requirements

**Characteristics:**
- Native ARM64 code generation
- Adaptive optimization (functions become "hot" after 100+ calls)
- Up to 15x performance improvement for hot functions
- Automatic fallback to bytecode VM
- ARM64 architecture only

### Performance Monitoring

All execution modes support performance monitoring through log levels:

```bash
# Basic execution information
rush -bytecode -log-level=info program.rush

# Detailed debugging information
rush -bytecode -log-level=debug program.rush

# Full execution tracing
rush -bytecode -log-level=trace program.rush
```

Example JIT statistics output:
```
JIT Statistics:
  Compilations attempted: 45
  Compilations succeeded: 12
  JIT hits: 156
  JIT misses: 23
  Deoptimizations: 2
  Hot threshold: 100 calls
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

### Hashes/Dictionaries

Store key-value pairs:

```rush
# Hash creation
person = {"name": "Alice", "age": 30, "active": true}
config = {42: "answer", "pi": 3.14, true: "enabled"}
empty = {}

# Hash access
name = person["name"]      # "Alice"
age = person["age"]        # 30

# Hash assignment
person["city"] = "NYC"     # Add new key
person["age"] = 31         # Update existing key

# Check for keys
has_name = builtin_hash_has_key(person, "name")  # true
has_phone = builtin_hash_has_key(person, "phone")  # false

# Get all keys and values
keys = builtin_hash_keys(person)    # ["name", "age", "active", "city"]
values = builtin_hash_values(person) # ["Alice", 31, true, "NYC"]

# Hash length
size = len(person)         # 4

# Nested hashes
users = {
  "alice": {"role": "admin", "score": 95},
  "bob": {"role": "user", "score": 87}
}
alice_role = users["alice"]["role"]  # "admin"
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

## JSON Processing

Rush provides comprehensive JSON processing capabilities with an elegant object-oriented API. The JSON module offers parsing, manipulation, querying, and serialization with full support for method chaining and dot notation.

### Overview

JSON processing in Rush is designed around two core functions:
- `json_parse()` - Parse JSON strings into JSON objects
- `json_stringify()` - Convert Rush values to JSON strings

JSON objects support comprehensive dot notation methods for data manipulation, making complex JSON operations simple and readable.

### Basic JSON Operations

#### Parsing JSON Data

Use `json_parse()` to convert JSON strings into Rush JSON objects:

```rush
# Parse basic values
number_data = json_parse("42")
string_data = json_parse("\"Hello, Rush!\"")
boolean_data = json_parse("true")
null_data = json_parse("null")

# Parse complex structures
user_json = json_parse("{
  \"id\": 1001,
  \"name\": \"Alice Johnson\",
  \"email\": \"alice@example.com\",
  \"active\": true
}")

array_json = json_parse("[1, 2, 3, 4, 5]")

# Nested structures
complex_json = json_parse("{
  \"users\": [
    {\"name\": \"Alice\", \"role\": \"admin\"},
    {\"name\": \"Bob\", \"role\": \"user\"}
  ],
  \"metadata\": {
    \"version\": \"2.0\",
    \"created\": \"2024-01-15\"
  }
}")
```

#### Converting to JSON Strings

Use `json_stringify()` to convert Rush values to JSON:

```rush
# Convert basic values
json_stringify("hello")          # "\"hello\""
json_stringify(42)               # "42"
json_stringify(true)             # "true"
json_stringify(null)             # "null"

# Convert arrays
json_stringify([1, 2, 3])        # "[1,2,3]"

# Convert hash objects
user = {"name": "John", "age": 30}
json_stringify(user)             # "{\"age\":30,\"name\":\"John\"}"

# Complex nested structures
data = {"users": [{"name": "Alice"}], "count": 1}
json_stringify(data)             # "{\"count\":1,\"users\":[{\"name\":\"Alice\"}]}"
```

### JSON Object Properties and Methods

When you parse JSON, you get a JSON object with rich dot notation support:

#### Core Properties

```rush
user = json_parse("{\"name\": \"John\", \"age\": 30}")

# Access underlying data
user.data                        # Returns the parsed hash/array/value
user.type                        # Returns the data type as string
user.valid                       # Always true for parsed JSON
```

#### Data Access Methods

```rush
# Get values by key or index
user.get("name")                 # "John"
user.get("age")                  # 30
user.get("missing")              # null

# Array access
arr = json_parse("[\"a\", \"b\", \"c\"]")
arr.get(0)                       # "a"
arr.get(1)                       # "b"

# Check if keys/indices exist
user.has?("name")                # true
user.has?("email")               # false
arr.has?(2)                      # true
arr.has?(5)                      # false
```

#### Collection Operations

```rush
# Get all keys or values
user.keys()                      # ["name", "age"]
user.values()                    # ["John", 30]

# Get collection size
user.length()                    # 2
user.size()                      # Same as length()
```

#### Data Modification

JSON objects are immutable - modification methods return new objects:

```rush
# Set new values
updated_user = user.set("email", "john@example.com")
                  .set("active", true)

# Verify changes
updated_user.get("email")        # "john@example.com"
updated_user.get("active")       # true

# Original remains unchanged
user.has?("email")               # false
```

### Advanced JSON Operations

#### Path-Based Navigation

Use dot-separated paths to navigate nested structures:

```rush
data = json_parse("{
  \"company\": {
    \"employees\": [
      {\"name\": \"Alice\", \"department\": \"Engineering\"},
      {\"name\": \"Bob\", \"department\": \"Marketing\"}
    ],
    \"info\": {
      \"founded\": 2020,
      \"location\": \"San Francisco\"
    }
  }
}")

# Navigate with paths
data.path("company.info.founded")           # 2020
data.path("company.employees.0.name")       # "Alice"
data.path("company.employees.1.department") # "Marketing"

# Non-existent paths return null
data.path("company.missing.field")          # null
```

#### JSON Object Merging

Combine JSON objects while maintaining immutability:

```rush
base_config = json_parse("{
  \"database\": {\"host\": \"localhost\", \"port\": 5432},
  \"debug\": false
}")

overrides = json_parse("{
  \"database\": {\"host\": \"production.db\"},
  \"logging\": {\"level\": \"info\"}
}")

# Merge objects (overrides takes precedence)
final_config = base_config.merge(overrides)

# Access merged data
final_config.path("database.host")          # "production.db"
final_config.path("database.port")          # 5432 (preserved)
final_config.path("logging.level")          # "info" (added)
```

#### JSON Formatting

Control JSON output formatting:

```rush
data = json_parse("{\"name\": \"John\", \"items\": [1, 2, 3]}")

# Compact format (no extra whitespace)
data.compact()
# Returns: "{\"items\":[1,2,3],\"name\":\"John\"}"

# Pretty format with default indentation
data.pretty()
# Returns:
# {
#   "items": [
#     1,
#     2,
#     3
#   ],
#   "name": "John"
# }

# Custom indentation
data.pretty("    ")              # 4-space indents
data.pretty("\t")                # Tab indents
```

### Method Chaining

JSON methods support fluent chaining for complex operations:

```rush
# Build and transform data with chaining
result = json_parse("{}")
          .set("timestamp", "2024-01-15")
          .set("user", json_parse("{\"name\": \"Alice\"}"))
          .set("status", "active")
          .merge(json_parse("{\"version\": \"2.0\"}"))

# Access chained results
result.get("user").get("name")   # "Alice"
result.get("version")            # "2.0"

# Complex data transformation
api_response = json_parse("{\"data\": {\"users\": [{\"id\": 1, \"name\": \"Alice\"}]}}")
processed = api_response.path("data.users.0")
                       .set("processed_at", "2024-01-15")
                       .set("verified", true)
```

### JSON Processing Patterns

#### Configuration File Processing

```rush
# Load and parse configuration
config_content = file("config.json").read()
config = json_parse(config_content)

# Extract configuration values
database_url = config.path("database.connection.url")
api_key = config.path("services.api.key")
debug_mode = config.path("app.debug")

# Validate required settings
if !config.has?("database") {
    print("Error: Database configuration missing")
}
```

#### API Response Handling

```rush
# Parse API response
api_data = json_parse(http_response_body)

# Check for errors
if api_data.has?("error") {
    error_msg = api_data.get("error")
    print("API Error:", error_msg)
    return
}

# Process successful response
users = api_data.path("data.users")
user_count = users.length()

print("Loaded", user_count, "users")
for i = 0; i < user_count; i = i + 1 {
    user = users.get(i)
    print("- " + user.get("name"))
}
```

#### Data Transformation and Export

```rush
# Transform data structure
source_data = json_parse(input_json)

# Add metadata and transform
enhanced_data = source_data.set("processed_at", current_timestamp())
                          .set("version", "2.0")
                          .set("processor", "Rush JSON Module")

# Export as formatted JSON
formatted_output = enhanced_data.pretty("  ")
file("output.json").write(formatted_output)
```

#### Dynamic JSON Construction

```rush
# Build JSON programmatically
report = json_parse("{}")

# Add sections dynamically
report = report.set("header", json_parse("{
  \"title\": \"System Report\",
  \"generated\": \"2024-01-15\"
}"))

# Add array data
metrics = json_parse("[]")
for i = 0; i < data_points.length(); i = i + 1 {
    point = json_parse("{}")
           .set("timestamp", data_points.get(i).timestamp)
           .set("value", data_points.get(i).value)
    metrics = metrics.set(i, point)
}

report = report.set("metrics", metrics)

# Export final report
final_json = report.pretty()
```

### Error Handling

JSON operations can fail in several ways:

```rush
# Parsing errors
try {
    invalid_json = json_parse("{invalid: json}")
} catch error {
    print("Parse error:", error.message)
}

# Serialization errors
try {
    # Hash with non-string keys can't be serialized
    invalid_hash = {42: "numeric key"}
    json_stringify(invalid_hash)
} catch error {
    print("Stringify error:", error.message)
}

# Safe parsing with validation
parse_safely = fn(json_string) {
    try {
        parsed = json_parse(json_string)
        return parsed
    } catch error {
        print("Invalid JSON:", error.message)
        return null
    }
}
```

### Performance Tips

1. **Reuse parsed objects**: Parse once, access many times
2. **Use path navigation**: More efficient than nested `.get()` calls
3. **Batch modifications**: Chain `.set()` calls rather than creating intermediate objects
4. **Choose appropriate formatting**: Use `.compact()` for storage, `.pretty()` for debugging

### JSON Best Practices

1. **Validate input**: Always handle parse errors gracefully
2. **Use meaningful variable names**: JSON objects can contain complex nested data
3. **Leverage method chaining**: Creates more readable data transformation pipelines
4. **Use path navigation**: Simplifies access to deeply nested values
5. **Handle missing data**: Use `.has?()` to check existence before accessing
6. **Prefer immutable operations**: JSON objects return new instances, don't modify originals

## File System Operations

Rush provides comprehensive file system operations through built-in types with dot notation support. This makes working with files, directories, and paths intuitive and safe.

### Working with Files

Files in Rush are represented by File objects that provide methods for reading, writing, and manipulating files.

#### Creating File Objects

```rush
# Create a file object
config_file = file("config.txt")
log_file = file("logs/application.log")
```

#### Basic File Operations

**Writing to Files:**

```rush
# Simple file writing
message_file = file("greeting.txt")
message_file.open("w").write("Hello, Rush programming!").close()

# Check if write was successful
if message_file.exists?() {
    print("File created successfully")
    print("File size:", message_file.size(), "bytes")
}
```

**Reading from Files:**

```rush
# Read entire file content
if file("greeting.txt").exists?() {
    content = file("greeting.txt").open("r").read()
    print("File content:", content)
}

# More controlled reading
data_file = file("data.txt")
if data_file.exists?() {
    f = data_file.open("r")
    content = f.read()
    f.close()
    print("Data:", content)
}
```

**File Information:**

```rush
log_file = file("app.log")

# Check existence
if log_file.exists?() {
    print("Log file found")
    print("Path:", log_file.path)
    print("Size:", log_file.size(), "bytes")
    print("Currently open:", log_file.is_open)
} else {
    print("Log file not found")
}
```

#### File Modes

Rush supports several file modes for different operations:

```rush
# Read mode (default)
file("data.txt").open("r")    # or just .open()

# Write mode (creates new file or truncates existing)
file("output.txt").open("w")

# Append mode (creates new file or appends to existing)
file("log.txt").open("a")

# Read/write modes
file("data.txt").open("r+")   # Read and write (file must exist)
file("data.txt").open("w+")   # Read and write (truncates or creates)
file("log.txt").open("a+")    # Read and append (creates if needed)
```

#### Practical File Example

```rush
# Log file manager
fn log_message(message) {
    log_file = file("application.log")
    
    # Append message with timestamp
    timestamp = "2024-01-01 12:00:00"  # In real code, use date/time functions
    log_entry = "[" + timestamp + "] " + message + "\n"
    
    log_file.open("a").write(log_entry).close()
    print("Logged:", message)
}

# Usage
log_message("Application started")
log_message("User logged in")
log_message("Processing data")

# Read log file
if file("application.log").exists?() {
    logs = file("application.log").open("r").read()
    print("Complete log:")
    print(logs)
}
```

### Working with Directories

Directory objects provide functionality for creating, listing, and managing directories.

#### Basic Directory Operations

```rush
# Create directory objects
project_dir = directory("my_project")
temp_dir = directory("/tmp/rush_temp")

# Create directories
project_dir.create()
print("Project directory created at:", project_dir.path)

# Check if directory exists
if project_dir.exists?() {
    print("Directory exists and is ready to use")
}
```

#### Listing Directory Contents

```rush
# List current directory
current_dir = directory(".")
if current_dir.exists?() {
    files = current_dir.list()
    print("Current directory contains", len(files), "items:")
    for item in files {
        print("  -", item)
    }
}

# Create project structure
project = directory("website_project")
project.create()

# Create subdirectories
directory("website_project/css").create()
directory("website_project/js").create()
directory("website_project/images").create()

# List project structure
project_files = project.list()
print("Project structure:")
for item in project_files {
    print("  ", item)
}
```

#### Directory Management Example

```rush
# Project setup script
fn setup_project(project_name) {
    # Create main project directory
    main_dir = directory(project_name)
    if main_dir.exists?() {
        print("Project", project_name, "already exists!")
        return false
    }
    
    # Create directory structure
    main_dir.create()
    directory(project_name + "/src").create()
    directory(project_name + "/docs").create()
    directory(project_name + "/tests").create()
    
    # Create initial files
    readme = file(project_name + "/README.md")
    readme.open("w").write("# " + project_name + "\n\nA Rush project.\n").close()
    
    main_file = file(project_name + "/src/main.rush")
    main_file.open("w").write('print("Hello from ' + project_name + '!")\n').close()
    
    print("Project", project_name, "created successfully!")
    return true
}

# Create a new project
setup_project("my_app")

# Verify project structure
app_dir = directory("my_app")
if app_dir.exists?() {
    contents = app_dir.list()
    print("Created project with:", contents)
}
```

### Working with Paths

Path objects provide cross-platform path manipulation and are essential for building portable file system operations.

#### Path Construction and Manipulation

```rush
# Create path objects
home_path = path("/home/user")
config_path = path("config")

# Join paths (cross-platform)
full_config = home_path.join("config").join("app.conf")
print("Config path:", full_config.value)

# Path information
full_path = path("/home/user/documents/report.pdf")
print("Full path:", full_path.value)
print("Directory:", full_path.dirname())
print("Filename:", full_path.basename())
```

#### Path Cleaning and Normalization

```rush
# Clean messy paths
messy_path = path("./project/../config//settings.txt")
clean_path = messy_path.clean()
print("Original:", messy_path.value)
print("Cleaned:", clean_path.value)

# Make paths absolute
relative_path = path("data/input.txt")
absolute_path = relative_path.absolute()
print("Relative:", relative_path.value)
print("Absolute:", absolute_path.value)
```

#### Cross-Platform Path Building

```rush
# Build paths that work on any operating system
fn build_project_path(project, module, filename) {
    return path("projects")
        .join(project)
        .join("src")
        .join(module)
        .join(filename)
        .clean()
}

# Usage
script_path = build_project_path("webapp", "auth", "login.rush")
print("Script path:", script_path.value)

# Create the file using the path
script_file = file(script_path.value)
# First ensure parent directories exist
parent_dir = directory(script_path.dirname())
parent_dir.create()

# Now create the file
script_file.open("w").write("# Login module\nprint(\"User authentication\")").close()
```

### Practical File System Examples

#### Configuration File Manager

```rush
fn load_config(config_name) {
    config_path = path("config").join(config_name + ".conf").absolute()
    config_file = file(config_path.value)
    
    if config_file.exists?() {
        return config_file.open("r").read()
    } else {
        print("Config file not found:", config_path.value)
        return null
    }
}

fn save_config(config_name, data) {
    # Ensure config directory exists
    config_dir = directory("config")
    if !config_dir.exists?() {
        config_dir.create()
    }
    
    config_path = path("config").join(config_name + ".conf")
    config_file = file(config_path.value)
    
    config_file.open("w").write(data).close()
    print("Config saved:", config_path.value)
}

# Usage
save_config("database", "host=localhost\nport=5432\nname=myapp")
db_config = load_config("database")
if db_config != null {
    print("Database config:")
    print(db_config)
}
```

#### File Backup System

```rush
fn backup_file(source_path) {
    source = file(source_path)
    
    if !source.exists?() {
        print("Source file not found:", source_path)
        return false
    }
    
    # Create backup filename with timestamp
    source_info = path(source_path)
    backup_name = source_info.basename() + ".backup"
    backup_path = source_info.dirname() + "/" + backup_name
    
    # Read source and write backup
    content = source.open("r").read()
    backup_file = file(backup_path)
    backup_file.open("w").write(content).close()
    
    print("Backup created:", backup_path)
    return true
}

fn restore_backup(original_path) {
    source_info = path(original_path)
    backup_path = source_info.dirname() + "/" + source_info.basename() + ".backup"
    backup = file(backup_path)
    
    if !backup.exists?() {
        print("Backup not found:", backup_path)
        return false
    }
    
    content = backup.open("r").read()
    original = file(original_path)
    original.open("w").write(content).close()
    
    print("File restored from backup")
    return true
}

# Usage
file("important.txt").open("w").write("Critical data").close()
backup_file("important.txt")

# Simulate data loss
file("important.txt").open("w").write("Corrupted!").close()

# Restore from backup
restore_backup("important.txt")
```

#### Directory Synchronization

```rush
fn sync_directories(source_dir, target_dir) {
    source = directory(source_dir)
    target = directory(target_dir)
    
    if !source.exists?() {
        print("Source directory not found:", source_dir)
        return false
    }
    
    # Create target directory if it doesn't exist
    if !target.exists?() {
        target.create()
        print("Created target directory:", target_dir)
    }
    
    # Copy all files from source to target
    files = source.list()
    copied_count = 0
    
    for filename in files {
        source_file_path = path(source_dir).join(filename).value
        target_file_path = path(target_dir).join(filename).value
        
        source_file = file(source_file_path)
        if source_file.exists?() {
            content = source_file.open("r").read()
            target_file = file(target_file_path)
            target_file.open("w").write(content).close()
            copied_count = copied_count + 1
        }
    }
    
    print("Synchronized", copied_count, "files from", source_dir, "to", target_dir)
    return true
}

# Setup example directories
directory("documents").create()
file("documents/letter.txt").open("w").write("Dear friend...").close()
file("documents/report.txt").open("w").write("Annual report...").close()

# Sync to backup
sync_directories("documents", "documents_backup")

# Verify sync
backup_files = directory("documents_backup").list()
print("Backup contains:", backup_files)
```

### File System Best Practices

1. **Always check file existence** before performing operations
2. **Close files** after operations to free resources
3. **Use absolute paths** when possible for reliability
4. **Create parent directories** before creating files in nested structures
5. **Handle errors gracefully** with existence checks
6. **Use path objects** for cross-platform compatibility

### Security Considerations

- Path traversal (`..`) is automatically blocked in constructors
- File operations respect system permissions
- Always validate user input when constructing file paths
- Be careful with file modes - `w` mode will truncate existing files

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