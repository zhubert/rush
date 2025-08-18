# Rush Built-in Functions API Reference

This document provides comprehensive documentation for all built-in functions available in the Rush programming language.

## Table of Contents

1. [Execution Modes](#execution-modes)
2. [Core Built-in Functions](#core-built-in-functions)
3. [String Built-in Functions](#string-built-in-functions)
4. [Array Built-in Functions](#array-built-in-functions)
5. [Hash Built-in Functions](#hash-built-in-functions)
6. [Math Built-in Functions](#math-built-in-functions)
7. [JSON Built-in Functions](#json-built-in-functions)
8. [Regular Expression Built-in Functions](#regular-expression-built-in-functions)
9. [File System Built-in Functions](#file-system-built-in-functions)
10. [Standard Library Modules](#standard-library-modules)
11. [Error Types](#error-types)

## Execution Modes

Rush provides three execution modes with different performance characteristics:

### Tree-Walking Interpreter (Default)

```bash
rush program.rush
```

Direct AST evaluation. Best for development and debugging.

### Bytecode Virtual Machine

```bash
rush -bytecode program.rush
rush -bytecode -log-level=info program.rush  # With statistics
```

Compiles to bytecode and executes on a virtual machine. Better performance than interpreter.

**Log Levels:**
- `none` - No logging (default)
- `error` - Only errors
- `warn` - Warnings and errors
- `info` - Execution statistics and performance metrics
- `debug` - Detailed execution flow
- `trace` - Instruction-by-instruction tracing

### JIT Compilation (ARM64)

```bash
rush -jit program.rush
rush -jit -log-level=info program.rush  # With statistics
```

Just-In-Time compilation to native ARM64 machine code. Highest performance for compute-intensive programs.

**JIT Features:**
- Hot path detection (functions become "hot" after 100+ calls)
- Native ARM64 code generation
- Adaptive optimization
- Automatic fallback to bytecode VM
- Performance statistics

**JIT Statistics Output:**
```
JIT Statistics:
  Compilations attempted: 45
  Compilations succeeded: 12
  JIT hits: 156
  JIT misses: 23
  Deoptimizations: 2
```

## Core Built-in Functions

These functions are available globally without any imports.

### `print(...)`

Outputs values to standard output, separated by spaces.

**Syntax:**
```rush
print(value1, value2, ...)
```

**Parameters:**
- `...` (variadic): Any number of values of any type

**Returns:**
- `null`

**Description:**
Prints all provided arguments to standard output, separated by spaces, followed by a newline. Values are automatically converted to their string representation.

**Examples:**
```rush
print("Hello, World!")
# Output: Hello, World!

print("The answer is", 42)
# Output: The answer is 42

print(true, 3.14, [1, 2, 3])
# Output: true 3.14 [1, 2, 3]

print()  # Empty call - prints just a newline
# Output: (empty line)
```

**Type Conversions:**
- `INTEGER`: Printed as decimal number
- `FLOAT`: Printed as decimal number with decimal point
- `STRING`: Printed as-is (without quotes)
- `BOOLEAN`: Printed as "true" or "false"
- `ARRAY`: Printed in format `[element1, element2, ...]`
- `FUNCTION`: Printed as `<function>`
- `NULL`: Printed as `null`

## Utility Functions

### `len(collection)`

Returns the length of a collection (array or string).

**Syntax:**
```rush
len(collection)
```

**Parameters:**
- `collection` (`ARRAY` | `STRING`): The collection to measure

**Returns:**
- `INTEGER`: The number of elements in the array or characters in the string

**Errors:**
- Returns error if argument is not an array or string
- Returns error if wrong number of arguments provided

**Examples:**
```rush
len([1, 2, 3, 4, 5])    # Returns: 5
len("hello")            # Returns: 5
len("")                 # Returns: 0
len([])                 # Returns: 0

# Nested arrays
len([[1, 2], [3, 4]])   # Returns: 2 (two sub-arrays)

# Error cases
len(42)                 # Error: argument to `len` not supported, got INTEGER
len("hello", "world")   # Error: wrong number of arguments. got=2, want=1
```

### `type(value)`

Returns the type of a value as a string.

**Syntax:**
```rush
type(value)
```

**Parameters:**
- `value` (any type): The value to get the type of

**Returns:**
- `STRING`: The type name

**Type Names:**
- `"INTEGER"` - for integer values
- `"FLOAT"` - for floating-point values
- `"STRING"` - for string values
- `"BOOLEAN"` - for boolean values
- `"ARRAY"` - for array values
- `"FUNCTION"` - for function values
- `"NULL"` - for null values

**Examples:**
```rush
type(42)                # Returns: "INTEGER"
type(3.14)              # Returns: "FLOAT"
type("hello")           # Returns: "STRING"
type(true)              # Returns: "BOOLEAN"
type([1, 2, 3])         # Returns: "ARRAY"
type(fn() {})           # Returns: "FUNCTION"

# Use in conditional logic
value = 42
if (type(value) == "INTEGER") {
  print("It's an integer!")
}
```

### `to_string(value)`

Converts any value to its string representation.

**Syntax:**
```rush
to_string(value)
```

**Parameters:**
- `value` (any): Value to convert to string

**Returns:**
- `STRING`: String representation of the value

**Examples:**
```rush
to_string(42)          # Returns: "42"
to_string(3.14)        # Returns: "3.14"
to_string(true)        # Returns: "true"
to_string([1, 2, 3])   # Returns: "[1, 2, 3]"
to_string(null)        # Returns: "null"
to_string(fn() {})     # Returns: "<function>"
```

## String Built-in Functions

### `substr(string, start, length)`

Extracts a substring from a string.

**Syntax:**
```rush
substr(string, start, length)
```

**Parameters:**
- `string` (`STRING`): The source string
- `start` (`INTEGER`): The starting index (0-based)
- `length` (`INTEGER`): The number of characters to extract

**Returns:**
- `STRING`: The extracted substring

**Behavior:**
- If `start` is negative or beyond string length, returns empty string
- If `start + length` exceeds string length, returns substring to end of string
- If `length` is 0 or negative, returns empty string

**Errors:**
- Returns error if first argument is not a string
- Returns error if wrong number of arguments provided

**Examples:**
```rush
substr("hello", 0, 5)   # Returns: "hello"
substr("hello", 1, 3)   # Returns: "ell"
substr("hello", 0, 2)   # Returns: "he"
substr("hello", 2, 3)   # Returns: "llo"

# Edge cases
substr("hello", 10, 1)  # Returns: "" (start beyond string)
substr("hello", -1, 1)  # Returns: "" (negative start)
substr("hello", 0, 10)  # Returns: "hello" (length exceeds string)
substr("hello", 2, 0)   # Returns: "" (zero length)

# Error cases
substr(42, 0, 1)        # Error: argument to `substr` must be STRING, got INTEGER
substr("hello")         # Error: wrong number of arguments. got=1, want=3
```

### `split(string, separator)`

Splits a string into an array of substrings using a separator.

**Syntax:**
```rush
split(string, separator)
```

**Parameters:**
- `string` (`STRING`): The string to split
- `separator` (`STRING`): The separator to split on

**Returns:**
- `ARRAY`: Array of string parts

**Behavior:**
- If separator is empty string (`""`), splits into individual characters
- If separator is not found, returns array with original string as single element
- If string is empty, returns array with single empty string
- Consecutive separators create empty string elements

**Errors:**
- Returns error if first argument is not a string
- Returns error if second argument is not a string
- Returns error if wrong number of arguments provided

**Examples:**
```rush
split("hello,world", ",")     # Returns: ["hello", "world"]
split("a b c", " ")           # Returns: ["a", "b", "c"]
split("hello", "")            # Returns: ["h", "e", "l", "l", "o"]
split("", ",")                # Returns: [""]
split("hello", "xyz")         # Returns: ["hello"]

# Multiple character separator
split("hello--world", "--")   # Returns: ["hello", "world"]

# Consecutive separators
split("a,,b", ",")            # Returns: ["a", "", "b"]

# Error cases
split(42, ",")                # Error: first argument to `split` must be STRING, got INTEGER
split("hello", 42)            # Error: second argument to `split` must be STRING, got INTEGER
split("hello")                # Error: wrong number of arguments. got=1, want=2
```

## Standard Library String Functions

The following functions are available through the `std/string` module and provide advanced string manipulation capabilities. Note that boolean predicate functions use the `?` suffix.

### `contains?(string, substring)`

Checks if a string contains a substring.

**Syntax:**
```rush
string.contains?(string, substring)
```

**Parameters:**
- `string` (`STRING`): The source string
- `substring` (`STRING`): The substring to search for

**Returns:**
- `BOOLEAN`: `true` if substring is found, `false` otherwise

**Examples:**
```rush
text = "hello world"
has_world = string.contains?(text, "world")
# Returns: true

has_foo = string.contains?(text, "foo")
# Returns: false
```

### `starts_with?(string, prefix)`

Checks if a string starts with a specific prefix.

**Syntax:**
```rush
string.starts_with?(string, prefix)
```

**Parameters:**
- `string` (`STRING`): The source string
- `prefix` (`STRING`): The prefix to check for

**Returns:**
- `BOOLEAN`: `true` if string starts with prefix, `false` otherwise

**Examples:**
```rush
text = "hello world"
starts_hello = string.starts_with?(text, "hello")
# Returns: true

starts_world = string.starts_with?(text, "world")
# Returns: false
```

### `ends_with?(string, suffix)`

Checks if a string ends with a specific suffix.

**Syntax:**
```rush
string.ends_with?(string, suffix)
```

**Parameters:**
- `string` (`STRING`): The source string
- `suffix` (`STRING`): The suffix to check for

**Returns:**
- `BOOLEAN`: `true` if string ends with suffix, `false` otherwise

**Examples:**
```rush
text = "hello world"
ends_world = string.ends_with?(text, "world")
# Returns: true

ends_hello = string.ends_with?(text, "hello")
# Returns: false
```

### `is_whitespace_char?(character)`

Checks if a single character is whitespace (space, tab, newline, or carriage return).

**Syntax:**
```rush
string.is_whitespace_char?(character)
```

**Parameters:**
- `character` (`STRING`): A single character string

**Returns:**
- `BOOLEAN`: `true` if character is whitespace, `false` otherwise

**Examples:**
```rush
is_space = string.is_whitespace_char?(" ")
# Returns: true

is_letter = string.is_whitespace_char?("a")
# Returns: false

is_tab = string.is_whitespace_char?("\t")
# Returns: true
```

### `ord(character)` and `chr(code)`

Convert between characters and ASCII codes.

**Syntax:**
```rush
ord(character)
chr(ascii_code)
```

**Parameters:**
- `character` (STRING): Single character string
- `ascii_code` (INTEGER): ASCII code value (0-127)

**Returns:**
- `ord`: INTEGER - ASCII code of the character
- `chr`: STRING - Character corresponding to ASCII code

**Examples:**
```rush
ord("A")               # Returns: 65
ord("a")               # Returns: 97
chr(65)                # Returns: "A"
chr(97)                # Returns: "a"
ord("0")               # Returns: 48
chr(48)                # Returns: "0"
```

## Array Built-in Functions

### `push(array, element)`

Returns a new array with an element added to the end.

**Syntax:**
```rush
push(array, element)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `element` (any type): The element to add

**Returns:**
- `ARRAY`: New array with the element appended

**Behavior:**
- Original array is not modified (immutable operation)
- Element can be of any type
- Returns new array with all original elements plus the new element at the end

**Errors:**
- Returns error if first argument is not an array
- Returns error if wrong number of arguments provided

**Examples:**
```rush
push([1, 2, 3], 4)           # Returns: [1, 2, 3, 4]
push([], 1)                  # Returns: [1]
push(["a", "b"], "c")        # Returns: ["a", "b", "c"]

# Mixed types
push([1, "hello"], true)     # Returns: [1, "hello", true]

# Nested arrays
push([[1, 2]], [3, 4])       # Returns: [[1, 2], [3, 4]]

# Error cases
push(42, 1)                  # Error: first argument to `push` must be ARRAY, got INTEGER
push([1, 2, 3])              # Error: wrong number of arguments. got=1, want=2
```

### `pop(array)`

Returns the last element of an array.

**Syntax:**
```rush
pop(array)
```

**Parameters:**
- `array` (`ARRAY`): The source array

**Returns:**
- (any type): The last element of the array, or `null` if array is empty

**Behavior:**
- Original array is not modified
- If array is empty, returns `null`
- Only returns the element, does not remove it from the array

**Errors:**
- Returns error if argument is not an array
- Returns error if wrong number of arguments provided

**Examples:**
```rush
pop([1, 2, 3])               # Returns: 3
pop([42])                    # Returns: 42
pop(["hello", "world"])      # Returns: "world"
pop([])                      # Returns: null

# Mixed types
pop([1, "hello", true])      # Returns: true

# Nested arrays
pop([[1, 2], [3, 4]])        # Returns: [3, 4]

# Error cases
pop(42)                      # Error: argument to `pop` must be ARRAY, got INTEGER
pop([1, 2, 3], 1)           # Error: wrong number of arguments. got=2, want=1
```

### `slice(array, start, end)`

Returns a portion of an array from start index to end index (exclusive).

**Syntax:**
```rush
slice(array, start, end)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `start` (`INTEGER`): The starting index (inclusive, 0-based)
- `end` (`INTEGER`): The ending index (exclusive)

**Returns:**
- `ARRAY`: New array containing elements from start to end-1

**Behavior:**
- Original array is not modified
- `start` is inclusive, `end` is exclusive
- If `start` is negative or >= array length, returns empty array
- If `end` > array length, uses array length as end
- If `end` <= `start`, returns empty array

**Errors:**
- Returns error if first argument is not an array
- Returns error if wrong number of arguments provided

**Examples:**
```rush
slice([1, 2, 3, 4, 5], 1, 4)    # Returns: [2, 3, 4]
slice([1, 2, 3], 0, 2)          # Returns: [1, 2]
slice([1, 2, 3], 1, 3)          # Returns: [2, 3]

# Edge cases
slice([1, 2, 3], 0, 10)         # Returns: [1, 2, 3] (end > length)
slice([1, 2, 3], 5, 10)         # Returns: [] (start >= length)
slice([1, 2, 3], -1, 2)         # Returns: [] (negative start)
slice([1, 2, 3], 2, 1)          # Returns: [] (end <= start)

# Full array copy
slice([1, 2, 3], 0, 3)          # Returns: [1, 2, 3]

# Error cases
slice(42, 0, 1)                 # Error: first argument to `slice` must be ARRAY, got INTEGER
slice([1, 2, 3])                # Error: wrong number of arguments. got=1, want=3
```

## Standard Library Array Functions

The following functions are available through the `std/array` module and provide advanced array manipulation capabilities.

### `map(array, function)`

Transforms each element of an array using the provided function.

**Syntax:**
```rush
array.map(array, function)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `function` (`FUNCTION`): A function that takes one argument and returns a value

**Returns:**
- `ARRAY`: New array with transformed elements

**Examples:**
```rush
numbers = [1, 2, 3, 4]
doubled = array.map(numbers, fn(x) { return x * 2 })
# Returns: [2, 4, 6, 8]

words = ["hello", "world"]
uppercased = array.map(words, fn(word) { return upper(word) })
# Returns: ["HELLO", "WORLD"]
```

### `filter(array, function)`

Returns a new array containing only elements that satisfy the predicate function.

**Syntax:**
```rush
array.filter(array, function)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `function` (`FUNCTION`): A predicate function that returns true/false

**Returns:**
- `ARRAY`: New array with filtered elements

**Examples:**
```rush
numbers = [1, 2, 3, 4, 5, 6]
evens = array.filter(numbers, fn(x) { return x % 2 == 0 })
# Returns: [2, 4, 6]

words = ["cat", "elephant", "dog"]
long_words = array.filter(words, fn(word) { return len(word) > 3 })
# Returns: ["elephant"]
```

### `reduce(array, function, initial)`

Reduces an array to a single value using the provided accumulator function.

**Syntax:**
```rush
array.reduce(array, function, initial)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `function` (`FUNCTION`): Function that takes (accumulator, current_value) and returns new accumulator
- `initial` (any type): Initial value for the accumulator

**Returns:**
- (any type): The final accumulated value

**Examples:**
```rush
numbers = [1, 2, 3, 4]
sum = array.reduce(numbers, fn(acc, x) { return acc + x }, 0)
# Returns: 10

words = ["hello", "beautiful", "world"]
longest = array.reduce(words, fn(acc, word) { 
  if (len(word) > len(acc)) { return word }
  return acc 
}, "")
# Returns: "beautiful"
```

### `find(array, function)`

Returns the first element that satisfies the predicate function.

**Syntax:**
```rush
array.find(array, function)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `function` (`FUNCTION`): A predicate function that returns true/false

**Returns:**
- (any type): The first matching element, or `null` if not found

**Examples:**
```rush
numbers = [1, 3, 5, 2, 4]
first_even = array.find(numbers, fn(x) { return x % 2 == 0 })
# Returns: 2

fruits = ["apple", "banana", "cherry"]
long_fruit = array.find(fruits, fn(fruit) { return len(fruit) > 5 })
# Returns: "banana"

not_found = array.find([1, 3, 5], fn(x) { return x % 2 == 0 })
# Returns: null
```

### `index_of(array, element)`

Returns the index of the first occurrence of the specified element.

**Syntax:**
```rush
array.index_of(array, element)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `element` (any type): The element to search for

**Returns:**
- `INTEGER`: The index of the element (0-based), or -1 if not found

**Examples:**
```rush
fruits = ["apple", "banana", "cherry", "banana"]
index = array.index_of(fruits, "banana")
# Returns: 1 (first occurrence)

numbers = [10, 20, 30]
missing = array.index_of(numbers, 25)
# Returns: -1 (not found)
```

### `includes?(array, element)`

Checks if an array contains the specified element.

**Syntax:**
```rush
array.includes?(array, element)
```

**Parameters:**
- `array` (`ARRAY`): The source array
- `element` (any type): The element to search for

**Returns:**
- `BOOLEAN`: `true` if element is found, `false` otherwise

**Examples:**
```rush
numbers = [10, 20, 30]
has_twenty = array.includes?(numbers, 20)
# Returns: true

has_forty = array.includes?(numbers, 40)
# Returns: false
```

### `reverse(array)`

Returns a new array with elements in reverse order.

**Syntax:**
```rush
array.reverse(array)
```

**Parameters:**
- `array` (`ARRAY`): The source array

**Returns:**
- `ARRAY`: New array with elements reversed

**Examples:**
```rush
numbers = [1, 2, 3, 4, 5]
reversed = array.reverse(numbers)
# Returns: [5, 4, 3, 2, 1]

words = ["first", "second", "third"]
backwards = array.reverse(words)
# Returns: ["third", "second", "first"]
```

### `sort(array)`

Returns a new array with elements sorted in ascending order.

**Syntax:**
```rush
array.sort(array)
```

**Parameters:**
- `array` (`ARRAY`): The source array

**Returns:**
- `ARRAY`: New array with elements sorted

**Behavior:**
- Sorts numbers in ascending numerical order
- Sorts strings in lexicographical order
- Original array is not modified

**Examples:**
```rush
numbers = [3, 1, 4, 1, 5, 9, 2, 6]
sorted_nums = array.sort(numbers)
# Returns: [1, 1, 2, 3, 4, 5, 6, 9]

words = ["zebra", "apple", "banana"]
sorted_words = array.sort(words)
# Returns: ["apple", "banana", "zebra"]
```

### `length(array)`

Returns the number of elements in an array (alias for built-in `len` function).

**Syntax:**
```rush
array.length(array)
```

**Parameters:**
- `array` (`ARRAY`): The source array

**Returns:**
- `INTEGER`: The number of elements in the array

**Examples:**
```rush
numbers = [1, 2, 3, 4, 5]
count = array.length(numbers)
# Returns: 5

empty = array.length([])
# Returns: 0
```

## Hash Built-in Functions

These functions provide comprehensive operations for working with hash/dictionary data structures.

### `builtin_hash_keys(hash)`

Returns an array containing all keys from the hash in insertion order.

**Syntax:**
```rush
builtin_hash_keys(hash)
```

**Parameters:**
- `hash` (`HASH`): The hash to extract keys from

**Returns:**
- `ARRAY`: Array of keys in insertion order

**Behavior:**
- Keys are returned in the order they were inserted into the hash
- Works with all hashable key types (string, integer, boolean, float)

**Errors:**
- Returns error if argument is not a hash
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice", "age": 30, "active": true}
keys = builtin_hash_keys(person)
# Returns: ["name", "age", "active"]

empty_keys = builtin_hash_keys({})
# Returns: []
```

### `builtin_hash_values(hash)`

Returns an array containing all values from the hash in insertion order.

**Syntax:**
```rush
builtin_hash_values(hash)
```

**Parameters:**
- `hash` (`HASH`): The hash to extract values from

**Returns:**
- `ARRAY`: Array of values in key insertion order

**Behavior:**
- Values are returned in the same order as their corresponding keys
- Values can be of any type

**Errors:**
- Returns error if argument is not a hash
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice", "age": 30}
values = builtin_hash_values(person)
# Returns: ["Alice", 30]
```

### `builtin_hash_has_key(hash, key)`

Checks if a hash contains a specific key.

**Syntax:**
```rush
builtin_hash_has_key(hash, key)
```

**Parameters:**
- `hash` (`HASH`): The hash to check
- `key` (hashable): The key to look for

**Returns:**
- `BOOLEAN`: True if key exists, false otherwise

**Behavior:**
- Works with all hashable key types
- Case-sensitive for string keys

**Errors:**
- Returns error if first argument is not a hash
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice", "age": 30}
has_name = builtin_hash_has_key(person, "name")    # Returns: true
has_city = builtin_hash_has_key(person, "city")    # Returns: false
has_answer = builtin_hash_has_key({42: "answer"}, 42)  # Returns: true
```

### `builtin_hash_get(hash, key, default?)`

Gets a value from the hash with optional default.

**Syntax:**
```rush
builtin_hash_get(hash, key)
builtin_hash_get(hash, key, default)
```

**Parameters:**
- `hash` (`HASH`): The hash to get value from
- `key` (hashable): The key to look up
- `default` (any type, optional): Value to return if key doesn't exist

**Returns:**
- Any type: The value associated with the key, or default/null if not found

**Behavior:**
- Returns the value if key exists
- Returns default value if provided and key doesn't exist
- Returns null if no default provided and key doesn't exist

**Errors:**
- Returns error if first argument is not a hash
- Returns error if wrong number of arguments provided (must be 2 or 3)

**Examples:**
```rush
person = {"name": "Alice", "age": 30}
name = builtin_hash_get(person, "name")           # Returns: "Alice"
city = builtin_hash_get(person, "city")           # Returns: null
city = builtin_hash_get(person, "city", "NYC")    # Returns: "NYC"
```

### `builtin_hash_set(hash, key, value)`

Returns a new hash with the key-value pair added or updated.

**Syntax:**
```rush
builtin_hash_set(hash, key, value)
```

**Parameters:**
- `hash` (`HASH`): The source hash
- `key` (hashable): The key to set
- `value` (any type): The value to associate with the key

**Returns:**
- `HASH`: New hash with the key-value pair added/updated

**Behavior:**
- Original hash is not modified (immutable operation)
- If key exists, its value is updated in the new hash
- If key doesn't exist, it's added to the new hash
- Maintains insertion order for new keys

**Errors:**
- Returns error if first argument is not a hash
- Returns error if key is not hashable
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice"}
updated = builtin_hash_set(person, "age", 30)
# Returns: {"name": "Alice", "age": 30}

modified = builtin_hash_set(person, "name", "Bob")
# Returns: {"name": "Bob"}
```

### `builtin_hash_delete(hash, key)`

Returns a new hash with the specified key removed.

**Syntax:**
```rush
builtin_hash_delete(hash, key)
```

**Parameters:**
- `hash` (`HASH`): The source hash
- `key` (hashable): The key to remove

**Returns:**
- `HASH`: New hash without the specified key

**Behavior:**
- Original hash is not modified (immutable operation)
- If key exists, it's removed from the new hash
- If key doesn't exist, returns a copy of the original hash
- Maintains order of remaining keys

**Errors:**
- Returns error if first argument is not a hash
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice", "age": 30, "city": "NYC"}
without_age = builtin_hash_delete(person, "age")
# Returns: {"name": "Alice", "city": "NYC"}

same = builtin_hash_delete(person, "nonexistent")
# Returns: {"name": "Alice", "age": 30, "city": "NYC"}
```

### `builtin_hash_merge(hash1, hash2)`

Returns a new hash combining two hashes, with hash2 values overwriting hash1.

**Syntax:**
```rush
builtin_hash_merge(hash1, hash2)
```

**Parameters:**
- `hash1` (`HASH`): The base hash
- `hash2` (`HASH`): The hash to merge in (takes precedence)

**Returns:**
- `HASH`: New hash combining both inputs

**Behavior:**
- Original hashes are not modified (immutable operation)
- Keys from hash2 overwrite keys from hash1 if they exist in both
- Maintains insertion order: hash1 keys first, then new keys from hash2
- Updated keys maintain their original position

**Errors:**
- Returns error if either argument is not a hash
- Returns error if wrong number of arguments provided

**Examples:**
```rush
person = {"name": "Alice", "age": 25}
updates = {"age": 30, "city": "NYC"}
merged = builtin_hash_merge(person, updates)
# Returns: {"name": "Alice", "age": 30, "city": "NYC"}

empty_merge = builtin_hash_merge({}, {"key": "value"})
# Returns: {"key": "value"}
```

## Math Built-in Functions

These math functions are implemented as builtins and used by the standard library.

### Mathematical Operations

```rush
builtin_abs(-42)           # Returns: 42
builtin_min(5, 3, 8, 1)    # Returns: 1
builtin_max(5, 3, 8, 1)    # Returns: 8
builtin_sqrt(16)           # Returns: 4.0
builtin_pow(2, 8)          # Returns: 256.0
```

### Rounding Functions

```rush
builtin_floor(3.7)         # Returns: 3.0
builtin_ceil(3.2)          # Returns: 4.0
builtin_round(3.6)         # Returns: 4.0
```

### Random Functions

```rush
builtin_random()           # Returns: random float [0,1)
builtin_random_int(1, 10)  # Returns: random integer [1,10]
```

### Array Math Functions

```rush
builtin_sum([1, 2, 3, 4])     # Returns: 10.0
builtin_average([1, 2, 3, 4]) # Returns: 2.5
```

### Type Predicates

```rush
builtin_is_number?(42)     # Returns: true
builtin_is_integer?(3.14)  # Returns: false
```

## JSON Built-in Functions

Rush provides comprehensive JSON processing capabilities with object-oriented dot notation for parsing, manipulation, and serialization.

### JSON Parsing and Serialization

Rush provides a clean static method API for JSON processing.

#### `JSON.parse(json_string)`

Static method that parses a JSON string and returns a JSON object with dot notation methods.

**Syntax:**
```rush
JSON.parse(json_string)
```

**Parameters:**
- `json_string` (`STRING`): A valid JSON string to parse

**Returns:**
- `JSON`: A JSON object containing the parsed data with dot notation methods

**Errors:**
- Returns error if argument is not a string
- Returns error if JSON string is invalid or malformed

**Examples:**
```rush
# Parse basic JSON types
JSON.parse("42")                    # JSON object containing integer 42
JSON.parse("\"hello\"")             # JSON object containing string "hello"
JSON.parse("true")                  # JSON object containing boolean true
JSON.parse("null")                  # JSON object containing null

# Parse complex structures
user = JSON.parse("{\"name\": \"John\", \"age\": 30}")
data = JSON.parse("[1, 2, 3]")
nested = JSON.parse("{\"user\": {\"profile\": {\"city\": \"SF\"}}}")

# Method chaining with parsed JSON
name = JSON.parse("{\"user\": {\"name\": \"Alice\"}}").get("user").get("name")

# Error cases
JSON.parse("{invalid json}")        # Error: invalid JSON
JSON.parse(42)                      # Error: argument must be STRING
```

#### `JSON.stringify(value)`

Static method that converts a Rush value to a JSON string representation.

**Syntax:**
```rush
JSON.stringify(value)
```

**Parameters:**
- `value` (any): A Rush value to convert to JSON

**Returns:**
- `STRING`: JSON string representation of the value

**Examples:**
```rush
# Stringify various types
JSON.stringify(42)                          # "42"
JSON.stringify("hello")                     # "\"hello\""
JSON.stringify(true)                        # "true"
JSON.stringify([1, 2, 3])                   # "[1,2,3]"
JSON.stringify({"name": "John", "age": 30}) # "{\"age\":30,\"name\":\"John\"}"

# Round-trip parsing and stringifying
original = "{\"key\": \"value\"}"
parsed = JSON.parse(original)
back_to_string = JSON.stringify(parsed)
```

**Returns:**
- `STRING`: JSON string representation of the value

**Errors:**
- Returns error if value cannot be serialized to JSON
- Returns error for hash objects with non-string keys

**Examples:**
```rush
# Basic value serialization
json_stringify("hello")             # Returns: "\"hello\""
json_stringify(42)                  # Returns: "42"
json_stringify(true)                # Returns: "true"
json_stringify(null)                # Returns: "null"

# Array serialization
json_stringify([1, 2, 3])           # Returns: "[1,2,3]"
json_stringify(["a", "b"])          # Returns: "[\"a\",\"b\"]"

# Hash serialization (object)
json_stringify({"name": "John"})    # Returns: "{\"name\":\"John\"}"

# Nested structures
nested = {"user": {"name": "Alice", "data": [1, 2]}}
json_stringify(nested)              # Returns: "{\"user\":{\"data\":[1,2],\"name\":\"Alice\"}}"

# Error case - non-string keys not supported in JSON
invalid = {42: "value"}
json_stringify(invalid)             # Error: JSON object keys must be strings
```

### JSON Object Methods

When you parse JSON using `json_parse()`, you get a JSON object that supports comprehensive dot notation methods for manipulation and querying.

#### Properties

##### `.data`
Returns the underlying parsed data structure.

```rush
obj = json_parse("{\"name\": \"John\"}")
obj.data                            # Returns: Hash with "name" -> "John"
```

##### `.type`
Returns the type of the underlying data as a string.

```rush
json_parse("42").type               # Returns: "INTEGER"
json_parse("\"hello\"").type        # Returns: "STRING"
json_parse("{}").type               # Returns: "HASH"
json_parse("[]").type               # Returns: "ARRAY"
```

##### `.valid`
Always returns `true` for successfully parsed JSON objects.

```rush
json_parse("{}").valid              # Returns: true
```

#### Data Access Methods

##### `.get(key)`
Retrieves a value by key (for objects) or index (for arrays).

**Parameters:**
- `key` (`STRING` for objects, `INTEGER` for arrays): The key or index to retrieve

**Returns:**
- The value at the specified key/index, or `null` if not found

```rush
# Object access
user = json_parse("{\"name\": \"John\", \"age\": 30}")
user.get("name")                    # Returns: "John"
user.get("age")                     # Returns: 30
user.get("missing")                 # Returns: null

# Array access
arr = json_parse("[\"a\", \"b\", \"c\"]")
arr.get(0)                          # Returns: "a"
arr.get(1)                          # Returns: "b"
arr.get(10)                         # Returns: null
```

##### `.set(key, value)`
Sets a value at the specified key/index. Returns a new JSON object.

**Parameters:**
- `key` (`STRING` for objects, `INTEGER` for arrays): The key or index
- `value` (any): The value to set

**Returns:**
- `JSON`: A new JSON object with the updated value

```rush
# Object modification
user = json_parse("{\"name\": \"John\"}")
updated = user.set("age", 30)
updated.get("age")                  # Returns: 30

# Array modification
arr = json_parse("[1, 2, 3]")
updated = arr.set(1, 99)
updated.get(1)                      # Returns: 99

# Method chaining
result = json_parse("{}")
          .set("name", "Alice")
          .set("age", 25)
          .set("active", true)
```

##### `.has?(key)`
Checks if a key exists in an object or if an index is valid in an array.

**Parameters:**
- `key` (`STRING` for objects, `INTEGER` for arrays): The key or index to check

**Returns:**
- `BOOLEAN`: `true` if key/index exists, `false` otherwise

```rush
# Object key checking
user = json_parse("{\"name\": \"John\", \"age\": 30}")
user.has?("name")                   # Returns: true
user.has?("email")                  # Returns: false

# Array index checking
arr = json_parse("[1, 2, 3]")
arr.has?(1)                         # Returns: true
arr.has?(5)                         # Returns: false
```

#### Collection Methods

##### `.keys()`
Returns all keys (for objects) or indices (for arrays) as an array.

**Returns:**
- `ARRAY`: Array of keys or indices

```rush
# Object keys
obj = json_parse("{\"b\": 2, \"a\": 1}")
obj.keys()                          # Returns: ["b", "a"] (insertion order)

# Array indices
arr = json_parse("[\"x\", \"y\", \"z\"]")
arr.keys()                          # Returns: [0, 1, 2]
```

##### `.values()`
Returns all values as an array.

**Returns:**
- `ARRAY`: Array of all values

```rush
# Object values
obj = json_parse("{\"name\": \"John\", \"age\": 30}")
obj.values()                        # Returns: ["John", 30] (insertion order)

# Array values (returns the array itself)
arr = json_parse("[1, 2, 3]")
arr.values()                        # Returns: [1, 2, 3]
```

##### `.length()` / `.size()`
Returns the number of elements in the JSON data.

**Returns:**
- `INTEGER`: Number of elements

```rush
# Object length
json_parse("{\"a\": 1, \"b\": 2}").length()    # Returns: 2

# Array length
json_parse("[1, 2, 3, 4]").size()              # Returns: 4

# String length
json_parse("\"hello\"").length()                # Returns: 5
```

#### Formatting Methods

##### `.pretty([indent])`
Returns a pretty-formatted JSON string with indentation.

**Parameters:**
- `indent` (`STRING`, optional): Custom indentation string (default: "  ")

**Returns:**
- `STRING`: Pretty-formatted JSON string

```rush
obj = json_parse("{\"name\": \"John\", \"data\": [1, 2]}")

# Default indentation
obj.pretty()
# Returns:
# {
#   "data": [
#     1,
#     2  
#   ],
#   "name": "John"
# }

# Custom indentation
obj.pretty("    ")                  # 4-space indentation
obj.pretty("\t")                    # Tab indentation
```

##### `.compact()`
Returns a compact JSON string without extra whitespace.

**Returns:**
- `STRING`: Compact JSON string

```rush
obj = json_parse("{ \"name\" : \"John\" , \"age\" : 30 }")
obj.compact()                       # Returns: "{\"age\":30,\"name\":\"John\"}"
```

#### Advanced Methods

##### `.validate()`
Validates the JSON object. Always returns `true` for parsed JSON.

**Returns:**
- `BOOLEAN`: Always `true` for valid JSON objects

```rush
json_parse("{}").validate()         # Returns: true
json_parse("\"test\"").validate()   # Returns: true
```

##### `.merge(other)`
Merges two JSON objects. Only works with object types (hashes).

**Parameters:**
- `other` (`JSON`): Another JSON object to merge

**Returns:**
- `JSON`: New JSON object with merged data

```rush
obj1 = json_parse("{\"a\": 1, \"b\": 2}")
obj2 = json_parse("{\"b\": 3, \"c\": 4}")

merged = obj1.merge(obj2)
merged.get("a")                     # Returns: 1
merged.get("b")                     # Returns: 3 (obj2 overwrites)
merged.get("c")                     # Returns: 4

# Error case
json_parse("[1, 2]").merge(json_parse("[3, 4]"))  # Error: can only merge objects
```

##### `.path(path_string)`
Retrieves a value using a dot-separated path string.

**Parameters:**
- `path_string` (`STRING`): Dot-separated path (e.g., "user.profile.name")

**Returns:**
- The value at the path, or `null` if path doesn't exist

```rush
data = json_parse("{
  \"user\": {
    \"profile\": {
      \"name\": \"Alice\",
      \"contacts\": [\"email@example.com\", \"555-1234\"]
    }
  }
}")

# Navigate nested objects
data.path("user.profile.name")      # Returns: "Alice"
data.path("user.profile.contacts.0") # Returns: "email@example.com"

# Non-existent paths
data.path("user.missing.field")     # Returns: null
data.path("nonexistent")            # Returns: null
```

### Method Chaining

All JSON methods that return JSON objects support method chaining for fluent data manipulation:

```rush
# Complex chaining example
result = json_parse("{}")
          .set("user", json_parse("{\"name\": \"John\"}"))
          .set("timestamp", "2024-01-15")
          .merge(json_parse("{\"active\": true}"))

# Access chained result
result.get("user").get("name")      # Returns: "John"
result.get("active")                # Returns: true

# Data transformation chain
api_data = json_parse("{\"users\": [{\"id\": 1, \"name\": \"Alice\"}]}")
first_user = api_data.path("users.0")
                    .set("last_seen", "2024-01-15")
                    .set("verified", true)
```

### JSON Usage Patterns

#### Configuration Processing
```rush
config = json_parse(file_content)
database_url = config.path("database.connection.url")
debug_mode = config.path("app.debug")
```

#### API Response Handling
```rush
response = json_parse(api_response)
if response.has("error") {
    print("API Error:", response.get("error"))
} else {
    data = response.get("data")
    print("Success:", data.pretty())
}
```

#### Data Transformation
```rush
# Transform and format data
source = json_parse(input_json)
transformed = source.set("processed_at", current_time())
                   .set("version", "2.0")
                   .merge(metadata)

output = transformed.pretty()
```

## Regular Expression Built-in Functions

Rush provides comprehensive regular expression support with pattern matching, text processing, and string manipulation capabilities.

### Regular Expression Constructor

#### `Regexp(pattern)`

Creates a regular expression object from a pattern string.

**Syntax:**
```rush
Regexp(pattern)
```

**Parameters:**
- `pattern` (string): The regular expression pattern

**Returns:**
- `Regexp`: A regular expression object

**Description:**
Creates a compiled regular expression object that can be used for pattern matching, finding, and replacing text. The pattern follows Go's regular expression syntax.

**Examples:**
```rush
# Create basic patterns
email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
number_pattern = Regexp("\\d+")
word_pattern = Regexp("\\b\\w+\\b")

# Use case-insensitive flag
hello_pattern = Regexp("(?i)hello")
```

### Regular Expression Object Properties

#### `regexp.pattern`

Gets the original pattern string used to create the regular expression.

**Type:** `string` (read-only property)

**Examples:**
```rush
pattern = Regexp("\\d+")
print(pattern.pattern)  # Output: \d+
```

### Regular Expression Object Methods

#### `regexp.matches?(text)`

Tests whether the regular expression matches any part of the input text.

**Syntax:**
```rush
regexp.matches?(text)
```

**Parameters:**
- `text` (string): The text to test against the pattern

**Returns:**
- `boolean`: `true` if the pattern matches, `false` otherwise

**Examples:**
```rush
email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
print(email_pattern.matches?("Contact john@example.com"))  # true
print(email_pattern.matches?("No email here"))            # false
```

#### `regexp.find_first(text)`

Finds the first match of the regular expression in the text.

**Syntax:**
```rush
regexp.find_first(text)
```

**Parameters:**
- `text` (string): The text to search

**Returns:**
- `string`: The first match, or `null` if no match found

**Examples:**
```rush
number_pattern = Regexp("\\d+")
result = number_pattern.find_first("The answer is 42 and 24")
print(result)  # Output: 42
```

#### `regexp.find_all(text)`

Finds all matches of the regular expression in the text.

**Syntax:**
```rush
regexp.find_all(text)
```

**Parameters:**
- `text` (string): The text to search

**Returns:**
- `array`: Array of all matches, or empty array if no matches found

**Examples:**
```rush
number_pattern = Regexp("\\d+")
results = number_pattern.find_all("The answer is 42 and 24")
print(results)  # Output: ["42", "24"]

email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
emails = email_pattern.find_all("Contact john@example.com or support@rush-lang.org")
print(emails)  # Output: ["john@example.com", "support@rush-lang.org"]
```

#### `regexp.replace(text, replacement)`

Replaces all matches of the regular expression with the replacement string.

**Syntax:**
```rush
regexp.replace(text, replacement)
```

**Parameters:**
- `text` (string): The text to perform replacements on
- `replacement` (string): The replacement string

**Returns:**
- `string`: New string with all matches replaced

**Examples:**
```rush
number_pattern = Regexp("\\d+")
result = number_pattern.replace("I have 5 apples and 3 oranges", "X")
print(result)  # Output: I have X apples and X oranges

email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
safe_text = email_pattern.replace("Contact john@example.com for help", "[EMAIL]")
print(safe_text)  # Output: Contact [EMAIL] for help
```

### String Methods with Regular Expression Support

Several string methods support regular expression objects as arguments:

#### `string.match(regexp)`

Finds all matches of the regular expression in the string.

**Examples:**
```rush
text = "The price is $123.45 and tax is $12.34"
price_pattern = Regexp("\\$\\d+\\.\\d+")
prices = text.match(price_pattern)
print(prices)  # Output: ["$123.45", "$12.34"]
```

#### `string.matches?(regexp)`

Tests if the string matches the regular expression pattern.

**Examples:**
```rush
email = "user@example.com"
email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")
is_valid = email.matches?(email_pattern)
print(is_valid)  # Output: true
```

#### `string.replace(regexp, replacement)`

Replaces all matches of the regular expression with the replacement string.

**Examples:**
```rush
text = "Phone: 123-456-7890 or 987-654-3210"
phone_pattern = Regexp("\\d{3}-\\d{3}-\\d{4}")
private_text = text.replace(phone_pattern, "[PHONE]")
print(private_text)  # Output: Phone: [PHONE] or [PHONE]
```

#### `string.split(regexp)`

Splits the string using the regular expression as a delimiter.

**Examples:**
```rush
text = "apple,banana;cherry|date"
delimiter_pattern = Regexp("[,;|]")
fruits = text.split(delimiter_pattern)
print(fruits)  # Output: ["apple", "banana", "cherry", "date"]
```

### Common Regular Expression Patterns

Here are some useful patterns for common use cases:

```rush
# Email validation
email_pattern = Regexp("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")

# Phone numbers (US format)
phone_pattern = Regexp("\\(?\\d{3}\\)?[-\\s]?\\d{3}[-\\s]?\\d{4}")

# URLs
url_pattern = Regexp("https?://[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}(?:/\\S*)?")

# IP addresses
ip_pattern = Regexp("\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b")

# Numbers (integers and floats)
number_pattern = Regexp("-?\\d+(?:\\.\\d+)?")

# Words only (alphabetic characters)
word_pattern = Regexp("\\b[a-zA-Z]+\\b")

# Whitespace
whitespace_pattern = Regexp("\\s+")

# Case-insensitive matching
case_insensitive = Regexp("(?i)pattern")
```

## File System Built-in Functions

Rush provides comprehensive file system operations through three core built-in constructors and their associated methods using dot notation.

### File System Constructors

#### `file(path)`

Creates a File object for file operations.

**Syntax:**
```rush
file(path_string)
```

**Parameters:**
- `path_string` (string): The file path

**Returns:**
- `File`: A File object for the specified path

**Security:**
- Path traversal (`..`) is blocked for security

**Example:**
```rush
f = file("data.txt")
print(f.path)     # Outputs: data.txt
print(f.is_open)  # Outputs: false
```

#### `directory(path)`

Creates a Directory object for directory operations.

**Syntax:**
```rush
directory(path_string)
```

**Parameters:**
- `path_string` (string): The directory path

**Returns:**
- `Directory`: A Directory object for the specified path

**Example:**
```rush
dir = directory("/tmp/mydir")
print(dir.path)   # Outputs: /tmp/mydir
```

#### `path(value)`

Creates a Path object for path manipulation.

**Syntax:**
```rush
path(path_string)
```

**Parameters:**
- `path_string` (string): The path string

**Returns:**
- `Path`: A Path object with path manipulation methods

**Example:**
```rush
p = path("/tmp/test")
print(p.value)    # Outputs: /tmp/test
```

### File Object Methods and Properties

File objects support comprehensive file operations with proper error handling.

#### Properties

- `path`: Returns the file path as a string
- `is_open`: Returns true if the file is currently open

#### Methods

##### `open(mode)`

Opens the file in the specified mode.

**Parameters:**
- `mode` (string, optional): File mode - "r" (read), "w" (write), "a" (append), "r+" (read/write), "w+" (write/read), "a+" (append/read). Default: "r"

**Returns:**
- `File`: The File object (for method chaining)

**Example:**
```rush
f = file("data.txt").open("w")
# or
f = file("data.txt").open()  # defaults to read mode
```

##### `read()`

Reads the entire file content.

**Returns:**
- `String`: The file content

**Example:**
```rush
content = file("data.txt").open("r").read()
print(content)
```

##### `write(content)`

Writes content to the file.

**Parameters:**
- `content` (string): Content to write

**Returns:**
- `Integer`: Number of bytes written

**Example:**
```rush
bytes_written = file("output.txt").open("w").write("Hello, Rush!")
print("Wrote " + bytes_written + " bytes")
```

##### `close()`

Closes the file.

**Returns:**
- `File`: The File object

**Example:**
```rush
f = file("data.txt").open("r")
content = f.read()
f.close()
```

##### `exists?()`

Checks if the file exists.

**Returns:**
- `Boolean`: true if file exists, false otherwise

**Example:**
```rush
if file("config.txt").exists?() {
    print("Config file found")
}
```

##### `size()`

Gets the file size in bytes.

**Returns:**
- `Integer`: File size in bytes

**Example:**
```rush
file_size = file("data.txt").size()
print("File size: " + file_size + " bytes")
```

##### `delete()`

Deletes the file.

**Returns:**
- `Boolean`: true on successful deletion

**Example:**
```rush
if file("temp.txt").delete() {
    print("File deleted successfully")
}
```

### Directory Object Methods and Properties

Directory objects provide directory manipulation functionality.

#### Properties

- `path`: Returns the directory path as a string

#### Methods

##### `create()`

Creates the directory (including parent directories if needed).

**Returns:**
- `Directory`: The Directory object

**Example:**
```rush
directory("/tmp/new/nested/dir").create()
```

##### `list()`

Lists the contents of the directory.

**Returns:**
- `Array`: Array of strings containing directory entry names

**Example:**
```rush
contents = directory("/tmp").list()
for entry in contents {
    print(entry)
}
```

##### `exists?()`

Checks if the directory exists.

**Returns:**
- `Boolean`: true if directory exists, false otherwise

**Example:**
```rush
if directory("/home/user").exists?() {
    print("Directory exists")
}
```

##### `delete()`

Deletes the directory and all its contents.

**Returns:**
- `Boolean`: true on successful deletion

**Example:**
```rush
directory("/tmp/old_data").delete()
```

### Path Object Methods and Properties

Path objects provide cross-platform path manipulation.

#### Properties

- `value`: Returns the path string

#### Methods

##### `join(other)`

Joins this path with another path component.

**Parameters:**
- `other` (string): Path component to join

**Returns:**
- `Path`: New Path object with joined path

**Example:**
```rush
full_path = path("/tmp").join("data").join("file.txt")
print(full_path.value)  # /tmp/data/file.txt (Unix) or \tmp\data\file.txt (Windows)
```

##### `basename()`

Gets the final component of the path.

**Returns:**
- `String`: The basename

**Example:**
```rush
name = path("/tmp/data/file.txt").basename()
print(name)  # Outputs: file.txt
```

##### `dirname()`

Gets the directory component of the path.

**Returns:**
- `String`: The directory path

**Example:**
```rush
dir = path("/tmp/data/file.txt").dirname()
print(dir)  # Outputs: /tmp/data
```

##### `absolute()`

Converts the path to an absolute path.

**Returns:**
- `Path`: New Path object with absolute path

**Example:**
```rush
abs_path = path("file.txt").absolute()
print(abs_path.value)  # Outputs full absolute path
```

##### `clean()`

Cleans the path by removing redundant components.

**Returns:**
- `Path`: New Path object with cleaned path

**Example:**
```rush
clean_path = path("/tmp//data/../file.txt").clean()
print(clean_path.value)  # Outputs: /tmp/file.txt
```

### File System Usage Examples

#### Basic File Operations
```rush
# Write to a file
file("greeting.txt").open("w").write("Hello, Rush!").close()

# Read from a file
if file("greeting.txt").exists?() {
    content = file("greeting.txt").open("r").read()
    print("File content: " + content)
    file("greeting.txt").close()
}

# Get file information
if file("greeting.txt").exists?() {
    size = file("greeting.txt").size()
    print("File size: " + size + " bytes")
}
```

#### Directory Operations
```rush
# Create a directory structure
data_dir = directory("project/data")
data_dir.create()

# List directory contents
if data_dir.exists?() {
    files = data_dir.list()
    print("Directory contains " + len(files) + " items")
}
```

#### Path Manipulation
```rush
# Build paths cross-platform
config_path = path("config").join("app.conf").absolute()
config_file = file(config_path.value)

if config_file.exists?() {
    print("Config found at: " + config_path.value)
}

# Path components
full_path = path("/home/user/documents/file.txt")
print("Directory: " + full_path.dirname())
print("Filename: " + full_path.basename())
```

#### Method Chaining
```rush
# Path method chaining
final_path = path("/tmp")
    .join("project")
    .join("data")
    .join("output.txt")
    .clean()

print("Final path: " + final_path.value)
```

## Standard Library Modules

Rush includes comprehensive standard library modules that can be imported.

### Math Module (`std/math`)

**Import:**
```rush
import { function_name } from "std/math"
```

**Constants:**
- `PI` - π (3.141592653589793)
- `E` - Euler's number (2.718281828459045)

**Functions:**
- `abs(number)` - Absolute value
- `min(...numbers)` - Minimum value  
- `max(...numbers)` - Maximum value
- `floor(number)` - Floor function
- `ceil(number)` - Ceiling function
- `round(number)` - Round to nearest integer
- `sqrt(number)` - Square root
- `pow(base, exponent)` - Power function
- `random()` - Random float [0,1)
- `random_int(min, max)` - Random integer in range
- `sum(array)` - Sum of numeric array
- `average(array)` - Average of numeric array
- `is_number?(value)` - Check if value is a number
- `is_integer?(value)` - Check if value is an integer

**Example:**
```rush
import { sqrt, PI, sum } from "std/math"
radius = 5
area = PI * sqrt(radius)
numbers = [1, 2, 3, 4, 5]
total = sum(numbers)  # 15.0
```

### String Module (`std/string`)

**Functions:**
- `trim(string)` - Remove leading/trailing whitespace
- `ltrim(string)` - Remove leading whitespace
- `rtrim(string)` - Remove trailing whitespace
- `upper(string)` - Convert to uppercase
- `lower(string)` - Convert to lowercase
- `length(string)` - Get string length (alias for `len`)
- `contains?(string, search)` - Check if string contains substring
- `replace(string, old, new)` - Replace all occurrences
- `starts_with?(string, prefix)` - Check string prefix
- `ends_with?(string, suffix)` - Check string suffix
- `join(array, separator)` - Join array elements into string
- `is_whitespace_char?(char)` - Check if character is whitespace

**Example:**
```rush
import { trim, upper, replace } from "std/string"
text = "  hello world  "
clean = trim(text)           # "hello world"
loud = upper(clean)          # "HELLO WORLD"
greeting = replace(loud, "WORLD", "RUSH")  # "HELLO RUSH"
```

### Array Module (`std/array`)

**Functions:**
- `length(array)` - Get array length (alias for `len`)
- `map(array, function)` - Transform each element
- `filter(array, predicate)` - Filter elements by predicate
- `reduce(array, function, initial)` - Reduce array to single value
- `find(array, predicate)` - Find first matching element
- `index_of(array, element)` - Find index of element (-1 if not found)
- `includes?(array, element)` - Check if array contains element
- `reverse(array)` - Reverse array order
- `sort(array)` - Sort array (bubble sort implementation)

**Example:**
```rush
import { map, filter, reduce } from "std/array"
numbers = [1, 2, 3, 4, 5]
doubled = map(numbers, fn(x) { x * 2 })      # [2, 4, 6, 8, 10]
evens = filter(numbers, fn(x) { x % 2 == 0 }) # [2, 4]
sum = reduce(numbers, fn(acc, x) { acc + x }, 0) # 15
```

## Error Types

All built-in functions follow consistent error handling patterns:

### Common Error Types

1. **Type Mismatch Errors**
   ```rush
   len(42)  # Error: argument to `len` not supported, got INTEGER
   ```

2. **Argument Count Errors**
   ```rush
   len()           # Error: wrong number of arguments. got=0, want=1
   len("a", "b")   # Error: wrong number of arguments. got=2, want=1
   ```

3. **Parameter Type Errors**
   ```rush
   substr(42, 0, 1)        # Error: argument to `substr` must be STRING, got INTEGER
   split("hello", 42)      # Error: second argument to `split` must be STRING, got INTEGER
   ```

### Error Message Format

Error messages follow these patterns:
- `"argument to `<function>` not supported, got <TYPE>"` - for unsupported types
- `"wrong number of arguments. got=<N>, want=<M>"` - for incorrect argument count
- `"<ordinal> argument to `<function>` must be <TYPE>, got <TYPE>"` - for parameter type mismatch

### Handling Errors in Code

Since Rush doesn't have exception handling, errors are returned as values:

```rush
# Check if result is an error
result = len("hello")
if (type(result) == "ERROR") {
  print("Error occurred:", result)
} else {
  print("Length:", result)
}

# Defensive programming
safe_substr = fn(str, start, length) {
  if (type(str) != "STRING") {
    return ""  # Return safe default
  }
  return substr(str, start, length)
}
```

### Function Compatibility

All built-in functions are designed to work together:

```rush
# Chaining operations
text = "hello,world,rush"
words = split(text, ",")
first_word = words[0]
first_char = substr(first_word, 0, 1)
print("First character of first word:", first_char)  # "h"

# Array processing pipeline
numbers = [1, 2, 3, 4, 5]
last_number = pop(numbers)
extended = push(numbers, last_number * 2)
subset = slice(extended, 1, 4)
print("Result length:", len(subset))
```

This completes the API reference for all built-in functions in Rush. Each function is designed to be intuitive and composable, following consistent patterns for error handling and type checking.