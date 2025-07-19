# Language Design Requirements & Syntax Examples

## User Requirements Summary
Creating a general purpose programming language built with Go.

## Syntax Examples Provided

### Variable Declaration
```
a = "foo"
```
- Simple assignment without explicit type declarations or keywords
- Dynamic typing (inferred from the value)
- No semicolons required

### Conditionals
```
if (i > 12) {
  "foo"
} else {
  "bar"
}
```
- C-style if/else syntax with parentheses around conditions
- Braces for blocks
- Expression-based conditionals (returns a value)
- String literals using double quotes
- Comparison operators (`>`)

### Functions
```
fn add(a, b) {
  return a + b
}
```
- `fn` keyword for function declarations
- Parameter lists in parentheses without type annotations
- Explicit `return` statement
- Block syntax with braces

### Function Calls
```
add(5, 3)
```
- Standard function call syntax

### Loops
- **For loops**: C-style `for (i = 0; i < 10; i++) { ... }`
- **While loops**: `while (condition) { ... }`

### Arrays/Lists
```
[1, 2, 3]
```
- Standard array literal syntax with square brackets

### Comments
```
# This is a comment
```
- Comments begin the line with `#`

### Literals
- **Boolean literals**: `true`, `false`
- **Number literals**: `42`, `3.14`
- **String literals**: `"foo"`

### Operators
- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `==`, `!=`, `>`, `<`, `>=`, `<=`
- Logical: `&&`, `||`

## Complete Language Example
```
# Variable declarations
a = "foo"
b = 42
c = 3.14
d = true
e = [1, 2, 3]

# Function definition
fn add(a, b) {
  return a + b
}

# Function call
result = add(5, 3)

# Expression-based conditional
value = if (i > 12) {
  "foo"
} else {
  "bar"
}

# C-style for loop
for (i = 0; i < 10; i++) {
  # loop body
}

# While loop
while (condition) {
  # loop body
}
```

## Design Philosophy
- Clean, simple syntax
- Expression-based where possible
- Familiar C-style control flow
- Dynamic typing
- No semicolons required
- Functional programming features (expressions return values)