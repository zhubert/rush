# Rush Language Specification

## Overview

Rush is a dynamically-typed, interpreted programming language with clean syntax inspired by modern languages. It features C-style control flow, first-class functions, and built-in support for arrays and strings.

## Table of Contents

1. [Lexical Structure](#lexical-structure)
2. [Data Types](#data-types)
3. [Variables](#variables)
4. [Expressions](#expressions)
5. [Statements](#statements)
6. [Functions](#functions)
7. [Control Flow](#control-flow)
8. [Module System](#module-system)
9. [Built-in Functions](#built-in-functions)
10. [Grammar](#grammar)

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
- `true` - boolean literal
- `false` - boolean literal

### Operators

#### Arithmetic Operators
- `+` - addition
- `-` - subtraction
- `*` - multiplication
- `/` - division (always produces float result)

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

#### Access Operator
- `.` - module member access

### Delimiters
- `()` - parentheses for grouping and function calls
- `{}` - braces for code blocks
- `[]` - brackets for array literals and indexing
- `,` - comma separator
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

### Function

First-class functions:

```rush
add = fn(x, y) { return x + y }
greet = fn(name) { return "Hello, " + name }
```

### Null

Represents absence of value. Returned by:
- Out-of-bounds array access
- Functions without explicit return
- Failed operations

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
```

### Module Paths

Module paths can be:
- **Relative**: `./module` or `../parent/module`
- **Absolute**: `/path/to/module`

The `.rush` extension is added automatically if not specified.

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

## Grammar

### Precedence (highest to lowest)

1. Primary expressions: literals, identifiers, parentheses
2. Postfix: function calls, array indexing
3. Unary: `-`, `!`
4. Multiplicative: `*`, `/`
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
          | returnStatement ;

assignmentStatement = identifier "=" expression ;

expressionStatement = expression ;

blockStatement = "{" { statement } "}" ;

ifStatement = "if" "(" expression ")" blockStatement [ "else" blockStatement ] ;

whileStatement = "while" "(" expression ")" blockStatement ;

forStatement = "for" "(" assignmentStatement ";" expression ";" assignmentStatement ")" blockStatement ;

returnStatement = "return" [ expression ] ;

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