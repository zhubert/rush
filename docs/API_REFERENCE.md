# Rush Built-in Functions API Reference

This document provides comprehensive documentation for all built-in functions available in the Rush programming language.

## Table of Contents

1. [I/O Functions](#io-functions)
2. [Utility Functions](#utility-functions)
3. [String Functions](#string-functions)
4. [Array Functions](#array-functions)
5. [Error Handling](#error-handling)

## I/O Functions

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

## String Functions

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

## Array Functions

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

## Error Handling

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