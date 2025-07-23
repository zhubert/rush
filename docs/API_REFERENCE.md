# Rush Built-in Functions API Reference

This document provides comprehensive documentation for all built-in functions available in the Rush programming language.

## Table of Contents

1. [Core Built-in Functions](#core-built-in-functions)
2. [String Built-in Functions](#string-built-in-functions)
3. [Array Built-in Functions](#array-built-in-functions)
4. [Math Built-in Functions](#math-built-in-functions)
5. [Standard Library Modules](#standard-library-modules)
6. [Error Types](#error-types)

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