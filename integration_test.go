package main

import (
  "bytes"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "testing"
)

func TestCompletePrograms(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "Variable Assignment and Arithmetic",
      program: `
a = 10
b = 20
c = a + b
print(c)
`,
      expected: "30",
    },
    {
      name: "Function Definition and Call",
      program: `
add = fn(x, y) { return x + y }
result = add(15, 25)
print(result)
`,
      expected: "40",
    },
    {
      name: "Recursive Function - Factorial",
      program: `
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
result = factorial(5)
print(result)
`,
      expected: "120",
    },
    {
      name: "If-Else Conditional",
      program: `
score = 85
if (score >= 90) {
  grade = "A"
} else {
  if (score >= 80) {
    grade = "B"
  } else {
    grade = "C"
  }
}
print(grade)
`,
      expected: "B",
    },
    {
      name: "Boolean Logic",
      program: `
age = 25
hasLicense = true
canDrive = (age >= 16) && hasLicense
print(canDrive)
`,
      expected: "true",
    },
    {
      name: "Array Operations",
      program: `
numbers = [1, 2, 3, 4, 5]
first = numbers[0]
last = numbers[4]
sum = first + last
print(sum)
`,
      expected: "6",
    },
    {
      name: "String Operations",
      program: `
greeting = "Hello, World!"
hello = substr(greeting, 0, 5)
print(hello)
`,
      expected: "Hello",
    },
    {
      name: "Array Built-ins",
      program: `
original = [1, 2, 3]
extended = push(original, 4)
length = len(extended)
print(length)
`,
      expected: "4",
    },
    {
      name: "While Loop",
      program: `
i = 0
sum = 0
while (i < 5) {
  sum = sum + i
  i = i + 1
}
print(sum)
`,
      expected: "10",
    },
    {
      name: "For Loop",
      program: `
total = 0
for (j = 1; j <= 5; j = j + 1) {
  total = total + j
}
print(total)
`,
      expected: "15",
    },
    {
      name: "Nested Loops and Arrays",
      program: `
matrix = [[1, 2], [3, 4]]
sum = 0
for (row = 0; row < 2; row = row + 1) {
  for (col = 0; col < 2; col = col + 1) {
    element = matrix[row][col]
    sum = sum + element
  }
}
print(sum)
`,
      expected: "10",
    },
    {
      name: "Higher-Order Functions",
      program: `
apply = fn(f, x) { return f(x) }
square = fn(x) { return x * x }
result = apply(square, 8)
print(result)
`,
      expected: "64",
    },
    {
      name: "String Split and Array Access",
      program: `
message = "Hello,World,Rush"
parts = split(message, ",")
middle = parts[1]
print(middle)
`,
      expected: "World",
    },
    {
      name: "Mixed Types and Type Checking",
      program: `
value = 42
valueType = type(value)
print(valueType)
`,
      expected: "INTEGER",
    },
    {
      name: "Complex Expression Evaluation",
      program: `
result = (10 + 5) * 2 - 8 / 4
print(result)
`,
      expected: "28",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Create a temporary file with the program
      tmpfile, err := os.CreateTemp("", "rush_test_*.rush")
      if err != nil {
        t.Fatal(err)
      }
      defer os.Remove(tmpfile.Name())

      if _, err := tmpfile.Write([]byte(tt.program)); err != nil {
        t.Fatal(err)
      }
      if err := tmpfile.Close(); err != nil {
        t.Fatal(err)
      }

      // Run the Rush interpreter
      cmd := exec.Command("go", "run", "cmd/rush/main.go", tmpfile.Name())
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      if err := cmd.Run(); err != nil {
        t.Fatalf("Program execution failed: %v\nOutput: %s", err, out.String())
      }

      fullOutput := strings.TrimSpace(out.String())
      
      // Extract just the printed content (lines that don't start with "Rush interpreter", "Result:", or "Execution complete!")
      lines := strings.Split(fullOutput, "\n")
      var printedLines []string
      for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" && 
           !strings.HasPrefix(line, "Rush interpreter") && 
           !strings.HasPrefix(line, "Rush tree-walking interpreter") && 
           line != "Result: null" && 
           !strings.HasPrefix(line, "Execution complete!") {
          printedLines = append(printedLines, line)
        }
      }
      
      output := strings.Join(printedLines, "\n")
      if output != tt.expected {
        t.Errorf("Expected output %q, got %q (full output: %q)", tt.expected, output, fullOutput)
      }
    })
  }
}

func TestFileExecution(t *testing.T) {
  // Test that existing example files can be executed without errors
  exampleFiles := []string{
    "examples/test.rush",
    "examples/test_control_flow.rush",
    "examples/test_boolean.rush",
    "examples/function_test.rush",
    "examples/test_while.rush",
    "examples/test_for.rush",
    "examples/test_indexing.rush",
  }

  for _, filename := range exampleFiles {
    t.Run(filepath.Base(filename), func(t *testing.T) {
      if _, err := os.Stat(filename); os.IsNotExist(err) {
        t.Skipf("Example file %s does not exist", filename)
        return
      }

      cmd := exec.Command("go", "run", "cmd/rush/main.go", filename)
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      if err := cmd.Run(); err != nil {
        t.Errorf("Failed to execute %s: %v\nOutput: %s", filename, err, out.String())
      }
    })
  }
}

func TestErrorHandling(t *testing.T) {
  errorTests := []struct {
    name    string
    program string
    hasError bool
  }{
    {
      name: "Syntax Error",
      program: `
a = 5 +
`,
      hasError: true,
    },
    {
      name: "Runtime Error - Undefined Variable",
      program: `
print(undefinedVariable)
`,
      hasError: true,
    },
    {
      name: "Runtime Error - Type Mismatch",
      program: `
result = true + false
`,
      hasError: true,
    },
    {
      name: "Runtime Error - Wrong Function Arguments",
      program: `
add = fn(x, y) { return x + y }
result = add(5)
`,
      hasError: true,
    },
    {
      name: "Valid Program",
      program: `
print("Hello, World!")
`,
      hasError: false,
    },
  }

  for _, tt := range errorTests {
    t.Run(tt.name, func(t *testing.T) {
      // Create a temporary file with the program
      tmpfile, err := os.CreateTemp("", "rush_error_test_*.rush")
      if err != nil {
        t.Fatal(err)
      }
      defer os.Remove(tmpfile.Name())

      if _, err := tmpfile.Write([]byte(tt.program)); err != nil {
        t.Fatal(err)
      }
      if err := tmpfile.Close(); err != nil {
        t.Fatal(err)
      }

      // Run the Rush interpreter
      cmd := exec.Command("go", "run", "cmd/rush/main.go", tmpfile.Name())
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      err = cmd.Run()
      
      if tt.hasError && err == nil {
        t.Errorf("Expected error but program executed successfully. Output: %s", out.String())
      } else if !tt.hasError && err != nil {
        t.Errorf("Expected success but got error: %v\nOutput: %s", err, out.String())
      }
    })
  }
}

func TestREPLBasics(t *testing.T) {
  // Test basic REPL functionality by sending commands to stdin
  cmd := exec.Command("go", "run", "cmd/rush/main.go")
  
  stdin, err := cmd.StdinPipe()
  if err != nil {
    t.Fatal(err)
  }

  var out bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &out

  // Start the command
  if err := cmd.Start(); err != nil {
    t.Fatal(err)
  }

  // Send commands to REPL
  commands := []string{
    "5 + 3\n",
    "name = \"Rush\"\n",
    "print(name)\n",
    ":quit\n",
  }

  for _, command := range commands {
    if _, err := stdin.Write([]byte(command)); err != nil {
      t.Fatal(err)
    }
  }

  stdin.Close()

  // Wait for command to finish
  if err := cmd.Wait(); err != nil {
    // REPL might exit with non-zero status, which is normal
    // We just check that it produces some output
  }

  output := out.String()
  if !strings.Contains(output, "⛤") {
    t.Errorf("REPL output should contain prompt symbol ⛤, got: %s", output)
  }
}

func TestErrorHandlingIntegration(t *testing.T) {
  errorHandlingTests := []struct {
    name     string
    program  string
    expected string
    hasError bool
  }{
    {
      name: "Basic Try-Catch with Print",
      program: `
try {
  throw RuntimeError("test error")
} catch (error) {
  print("Caught:", error.message)
}
`,
      expected: "Caught: test error",
      hasError: false,
    },
    {
      name: "Type-Specific Exception Handling",
      program: `
try {
  throw ValidationError("validation failed")
} catch (ValidationError error) {
  print("Validation error:", error.message)
} catch (error) {
  print("Other error:", error.message)
}
`,
      expected: "Validation error: validation failed",
      hasError: false,
    },
    {
      name: "Finally Block Execution",
      program: `
print("Starting")
try {
  throw RuntimeError("error")
} catch (error) {
  print("Caught error")
} finally {
  print("Finally executed")
}
print("Done")
`,
      expected: "Starting\nCaught error\nFinally executed\nDone",
      hasError: false,
    },
    {
      name: "Array Index Exception Handling",
      program: `
arr = [1, 2, 3]
try {
  value = arr[10]
  print("No error")
} catch (IndexError error) {
  print("Index error caught")
}
`,
      expected: "Index error caught",
      hasError: false,
    },
    {
      name: "Stack Trace in Error Output",
      program: `
level2 = fn() {
  throw RuntimeError("deep error")
}

level1 = fn() {
  level2()
}

try {
  level1()
} catch (error) {
  print("Error type:", error.type)
  print("Has stack:", len(error.stack) > 0)
}
`,
      expected: "Error type: RuntimeError\nHas stack: true",
      hasError: false,
    },
    {
      name: "Nested Try-Catch Blocks",
      program: `
try {
  try {
    throw ValidationError("inner error")
  } catch (TypeError error) {
    print("Should not catch TypeError")
  }
  print("This should not execute")
} catch (ValidationError error) {
  print("Outer catch:", error.message)
}
`,
      expected: "Outer catch: inner error",
      hasError: false,
    },
    {
      name: "Exception Propagation Through Functions",
      program: `
deepFunc = fn() {
  throw ArgumentError("argument error")
}

middleFunc = fn() {
  deepFunc()
}

try {
  middleFunc()
} catch (ArgumentError error) {
  print("Propagated error:", error.message)
}
`,
      expected: "Propagated error: argument error",
      hasError: false,
    },
    {
      name: "All Error Types Work",
      program: `
errorTypes = ["Error", "ValidationError", "TypeError", "IndexError", "ArgumentError", "RuntimeError"]

for (i = 0; i < len(errorTypes); i = i + 1) {
  errorType = errorTypes[i]
  try {
    if (errorType == "Error") { throw Error("test") }
    if (errorType == "ValidationError") { throw ValidationError("test") }
    if (errorType == "TypeError") { throw TypeError("test") }
    if (errorType == "IndexError") { throw IndexError("test") }
    if (errorType == "ArgumentError") { throw ArgumentError("test") }
    if (errorType == "RuntimeError") { throw RuntimeError("test") }
  } catch (error) {
    print("Caught:", error.type)
  }
}
`,
      expected: "Caught: Error\nCaught: ValidationError\nCaught: TypeError\nCaught: IndexError\nCaught: ArgumentError\nCaught: RuntimeError",
      hasError: false,
    },
    {
      name: "Uncaught Exception Program Termination",
      program: `
print("Before exception")
throw RuntimeError("uncaught error")
print("After exception - should not print")
`,
      expected: "",
      hasError: true,
    },
    {
      name: "Exception in Built-in Function",
      program: `
emptyArr = []
try {
  value = pop(emptyArr)
} catch (IndexError error) {
  print("Built-in exception:", error.message)
}
`,
      expected: "Built-in exception: pop from empty array",
      hasError: false,
    },
    {
      name: "Error Properties Access",
      program: `
try {
  throw ValidationError("detailed error")
} catch (error) {
  print("Type:", error.type)
  print("Message:", error.message)
  print("Line:", error.line >= 0)
  print("Column:", error.column >= 0)
  print("Stack exists:", len(error.stack) >= 0)
}
`,
      expected: "Type: ValidationError\nMessage: detailed error\nLine: true\nColumn: true\nStack exists: true",
      hasError: false,
    },
    {
      name: "Try-Catch in Loop",
      program: `
for (i = 0; i < 3; i = i + 1) {
  try {
    if (i == 1) {
      throw ValidationError("error at 1")
    } else {
      print("Success at", i)
    }
  } catch (ValidationError error) {
    print("Error at", i)
  }
}
`,
      expected: "Success at 0\nError at 1\nSuccess at 2",
      hasError: false,
    },
  }

  for _, tt := range errorHandlingTests {
    t.Run(tt.name, func(t *testing.T) {
      // Create a temporary file with the program
      tmpfile, err := os.CreateTemp("", "rush_error_integration_*.rush")
      if err != nil {
        t.Fatal(err)
      }
      defer os.Remove(tmpfile.Name())

      if _, err := tmpfile.Write([]byte(tt.program)); err != nil {
        t.Fatal(err)
      }
      if err := tmpfile.Close(); err != nil {
        t.Fatal(err)
      }

      // Run the Rush interpreter
      cmd := exec.Command("go", "run", "cmd/rush/main.go", tmpfile.Name())
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      err = cmd.Run()
      
      if tt.hasError && err == nil {
        t.Errorf("Expected error but program executed successfully. Output: %s", out.String())
        return
      } else if !tt.hasError && err != nil {
        t.Errorf("Expected success but got error: %v\nOutput: %s", err, out.String())
        return
      }

      if !tt.hasError {
        fullOutput := strings.TrimSpace(out.String())
        
        // Extract just the printed content (lines that don't start with "Rush interpreter", "Result:", or "Execution complete!")
        lines := strings.Split(fullOutput, "\n")
        var printedLines []string
        for _, line := range lines {
          line = strings.TrimSpace(line)
          if line != "" && 
             !strings.HasPrefix(line, "Rush interpreter") && 
             !strings.HasPrefix(line, "Result:") && 
             !strings.HasPrefix(line, "Execution complete!") &&
             !strings.HasPrefix(line, "Error:") {
            printedLines = append(printedLines, line)
          }
        }
        
        output := strings.Join(printedLines, "\n")
        if output != tt.expected {
          t.Errorf("Expected output %q, got %q (full output: %q)", tt.expected, output, fullOutput)
        }
      }
    })
  }
}

func TestModuleErrorHandling(t *testing.T) {
  // Test error handling across module imports
  
  // Create a module with error handling
  moduleContent := `
export errorFunc = fn() {
  throw ValidationError("module error")
}

export safeFunc = fn() {
  try {
    errorFunc()
  } catch (ValidationError error) {
    return "handled in module"
  }
}
`
  
  mainContent := `
import { errorFunc, safeFunc } from "./test_module"

# Test 1: Exception propagates from module
try {
  errorFunc()
} catch (ValidationError error) {
  print("Caught from module:", error.message)
}

# Test 2: Exception handled within module
result = safeFunc()
print("Module result:", result)
`

  // Create temporary module file
  moduleFile, err := os.CreateTemp("", "test_module*.rush")
  if err != nil {
    t.Fatal(err)
  }
  defer os.Remove(moduleFile.Name())

  if _, err := moduleFile.Write([]byte(moduleContent)); err != nil {
    t.Fatal(err)
  }
  if err := moduleFile.Close(); err != nil {
    t.Fatal(err)
  }

  // Create main file in same directory
  dir := filepath.Dir(moduleFile.Name())
  mainFile := filepath.Join(dir, "main_test.rush")
  
  // Update import path in main content to use absolute path
  adjustedMainContent := strings.Replace(mainContent, "./test_module", strings.TrimSuffix(moduleFile.Name(), ".rush"), 1)
  
  if err := os.WriteFile(mainFile, []byte(adjustedMainContent), 0644); err != nil {
    t.Fatal(err)
  }
  defer os.Remove(mainFile)

  // Run the main program using go run from project root
  cmd := exec.Command("go", "run", "cmd/rush/main.go", mainFile)
  var out bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &out

  if err := cmd.Run(); err != nil {
    t.Errorf("Module error handling test failed: %v\nOutput: %s", err, out.String())
    return
  }

  fullOutput := strings.TrimSpace(out.String())
  expectedOutput := "Caught from module: module error\nModule result: handled in module"
  
  // Extract printed content
  lines := strings.Split(fullOutput, "\n")
  var printedLines []string
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if line != "" && 
       !strings.HasPrefix(line, "Rush interpreter") && 
       !strings.HasPrefix(line, "Result:") && 
       !strings.HasPrefix(line, "Execution complete!") {
      printedLines = append(printedLines, line)
    }
  }
  
  output := strings.Join(printedLines, "\n")
  if output != expectedOutput {
    t.Errorf("Expected module error handling output %q, got %q", expectedOutput, output)
  }
}

// Integration tests for converted example features

func TestConvertedBasicLanguageFeatures(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "Basic Variable Assignments and Arithmetic",
      program: `
# Basic variable assignments
a = 42
b = 3.14
c = "hello world"
d = true
e = false

# Array literal
numbers = [1, 2, 3, 4, 5]

# Arithmetic expressions
sum = a + 10
product = b * 2.5
difference = 100 - a

# Boolean expressions
isGreater = a > 30
isEqual = b == 3.14

# String assignment
message = "Rush programming language"

print("Sum:", sum)
print("Product:", product)
print("IsGreater:", isGreater)
print("Array length:", len(numbers))
`,
      expected: "Sum: 52\nProduct: 7.8500000000000005\nIsGreater: true\nArray length: 5",
    },
    {
      name: "Mixed Type Operations and String Concatenation",
      program: `
# Mixed number types
pi = 3.14
area = pi * 10

# String operations
greeting = "Hello, "
name = "Rush"
message = greeting + name

# Boolean operations
isTrue = true
isFalse = !isTrue
comparison = 5 > 3

# Array with mixed types
numbers = [1, 2, 3, 10, 16]

# Final calculations
final = area + 100

print("Area:", area)
print("Message:", message)
print("Is False:", isFalse)
print("Comparison:", comparison)
print("Final:", final)
`,
      expected: "Area: 31.400000000000002\nMessage: Hello, Rush\nIs False: false\nComparison: true\nFinal: 131.4",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runIntegrationTest(t, tt.program, tt.expected)
    })
  }
}

func TestConvertedBooleanLogicFeatures(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "Simple Conditionals",
      program: `
# Test basic if/else
x = 15
if (x > 10) {
    result = "greater"
} else {
    result = "smaller"
}

# Test boolean operators
a = true
b = false
andTest = a && b
orTest = a || b

print("Result:", result)
print("And test:", andTest)
print("Or test:", orTest)
`,
      expected: "Result: greater\nAnd test: false\nOr test: true",
    },
    {
      name: "Complex Boolean Logic",
      program: `
# Advanced boolean expressions
age = 25
hasLicense = true
isWeekend = false

canDrive = (age >= 16) && hasLicense
shouldStayHome = isWeekend || (age < 18)
complexLogic = ((age > 20) && hasLicense) || (isWeekend && (age >= 16))

# Nested boolean conditions
score = 85
if (score >= 90) {
    grade = "A"
} else {
    if (score >= 80) {
        grade = "B"
    } else {
        grade = "C"
    }
}

print("Can drive:", canDrive)
print("Should stay home:", shouldStayHome)
print("Complex logic:", complexLogic)
print("Grade:", grade)
`,
      expected: "Can drive: true\nShould stay home: false\nComplex logic: true\nGrade: B",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runIntegrationTest(t, tt.program, tt.expected)
    })
  }
}

func TestConvertedLoopFeatures(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "For Loop Operations",
      program: `
# Basic for loop
sum = 0
for (i = 0; i < 5; i = i + 1) {
    sum = sum + i
}

# For loop with array access
numbers = [2, 4, 6, 8, 10]
product = 1
for (j = 0; j < 5; j = j + 1) {
    product = product * numbers[j]
}

# Nested for loops
nestedSum = 0
for (x = 1; x <= 3; x = x + 1) {
    for (y = 1; y <= 2; y = y + 1) {
        nestedSum = nestedSum + (x * y)
    }
}

print("Sum:", sum)
print("Product:", product)
print("Nested sum:", nestedSum)
`,
      expected: "Sum: 10\nProduct: 3840\nNested sum: 18",
    },
    {
      name: "While Loop Operations",
      program: `
# Simple while loop
counter = 0
sum = 0
while (counter < 5) {
    sum = sum + counter
    counter = counter + 1
}

# While loop with array indexing
numbers = [1, 2, 3, 4, 5]
index = 0
product = 1
while (index < 5) {
    product = product * numbers[index]
    index = index + 1
}

# Conditional while loop
value = 1
iterations = 0
while (value < 100) {
    value = value * 2
    iterations = iterations + 1
}

print("Sum:", sum)
print("Product:", product)
print("Iterations:", iterations)
`,
      expected: "Sum: 10\nProduct: 120\nIterations: 7",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runIntegrationTest(t, tt.program, tt.expected)
    })
  }
}

func TestConvertedFunctionFeatures(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "Comprehensive Function Tests",
      program: `
# Simple function
test1 = fn() { return 42 }
result1 = test1()

# Function with parameters
test2 = fn(x) { return x * 2 }
result2 = test2(21)

# Multiple parameter function
test3 = fn(a, b, c) { return a + b + c }
result3 = test3(1, 2, 3)

# Recursive function (factorial)
factorial = fn(n) {
    if (n <= 1) {
        return 1
    } else {
        return n * factorial(n - 1)
    }
}
result4 = factorial(5)

# Higher-order function
apply = fn(f, x) { return f(x) }
square = fn(x) { return x * x }
result5 = apply(square, 7)

print("Simple:", result1)
print("Double:", result2)
print("Add three:", result3)
print("Factorial:", result4)
print("Higher-order:", result5)
`,
      expected: "Simple: 42\nDouble: 42\nAdd three: 6\nFactorial: 120\nHigher-order: 49",
    },
    {
      name: "Advanced Function Patterns",
      program: `
# Simple arithmetic function
add = fn(a, b) { return a + b }
multiply = fn(a, b) { return a * b }
mathResult = add(multiply(2, 3), 4)

# Function composition
double = fn(x) { return x * 2 }
addTen = fn(x) { return x + 10 }
compose = fn(f, g, x) { return f(g(x)) }
composed = compose(double, addTen, 5)

# Array processing function
sumArray = fn(arr) {
    total = 0
    for (i = 0; i < len(arr); i = i + 1) {
        total = total + arr[i]
    }
    return total
}
arraySum = sumArray([1, 2, 3, 4, 5])

print("Math result:", mathResult)
print("Composed:", composed)
print("Array sum:", arraySum)
`,
      expected: "Math result: 10\nComposed: 30\nArray sum: 15",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runIntegrationTest(t, tt.program, tt.expected)
    })
  }
}

func TestConvertedIndexingFeatures(t *testing.T) {
  tests := []struct {
    name     string
    program  string
    expected string
  }{
    {
      name: "Array and String Indexing",
      program: `
# Array indexing
numbers = [10, 20, 30, 40, 50]
first = numbers[0]
second = numbers[1]
last = numbers[4]

# String indexing
message = "hello"
char1 = message[0]
char2 = message[1]
char5 = message[4]

# Expression-based indexing
index = 2
third = numbers[index]
fourth = numbers[index + 1]

result = first + second + last

print("Array result:", result)
print("String chars:", char1 + char2 + char5)
print("Third:", third)
print("Fourth:", fourth)
`,
      expected: "Array result: 80\nString chars: heo\nThird: 30\nFourth: 40",
    },
    {
      name: "Indexing with Exception Handling",
      program: `
# Valid array access
numbers = [1, 2, 3]
validAccess = numbers[1]

# Exception handling for out-of-bounds
caught1 = false
try {
    outOfBounds = numbers[10]
} catch (IndexError error) {
    caught1 = true
}

# String bounds checking
message = "hello"
caught2 = false
try {
    badChar = message[20]
} catch (IndexError error) {
    caught2 = true
}

# Safe indexing function
safeGet = fn(arr, index) {
    try {
        return arr[index]
    } catch (IndexError error) {
        return -1
    }
}

safeValue = safeGet(numbers, 1)
invalidValue = safeGet(numbers, 10)

print("Valid access:", validAccess)
print("Caught array:", caught1)
print("Caught string:", caught2)
print("Safe valid:", safeValue)
print("Safe invalid:", invalidValue)
`,
      expected: "Valid access: 2\nCaught array: true\nCaught string: true\nSafe valid: 2\nSafe invalid: -1",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      runIntegrationTest(t, tt.program, tt.expected)
    })
  }
}

// Helper function for running integration tests
func runIntegrationTest(t *testing.T, program, expected string) {
  // Create a temporary file with the program
  tmpfile, err := os.CreateTemp("", "rush_integration_test_*.rush")
  if err != nil {
    t.Fatal(err)
  }
  defer os.Remove(tmpfile.Name())

  if _, err := tmpfile.Write([]byte(program)); err != nil {
    t.Fatal(err)
  }
  if err := tmpfile.Close(); err != nil {
    t.Fatal(err)
  }

  // Run the Rush interpreter
  cmd := exec.Command("go", "run", "cmd/rush/main.go", tmpfile.Name())
  var out bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &out

  if err := cmd.Run(); err != nil {
    t.Fatalf("Program execution failed: %v\nOutput: %s", err, out.String())
  }

  fullOutput := strings.TrimSpace(out.String())
  
  // Extract just the printed content
  lines := strings.Split(fullOutput, "\n")
  var printedLines []string
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if line != "" && 
       !strings.HasPrefix(line, "Rush interpreter") && 
       line != "Result: null" && 
       !strings.HasPrefix(line, "Execution complete!") {
      printedLines = append(printedLines, line)
    }
  }
  
  output := strings.Join(printedLines, "\n")
  if output != expected {
    t.Errorf("Expected output %q, got %q (full output: %q)", expected, output, fullOutput)
  }
}

func TestModuleImportExport(t *testing.T) {
  tests := []struct {
    name         string
    moduleContent string
    mainContent   string
    expected      string
  }{
    {
      name: "Basic Function Export and Import",
      moduleContent: `export add = fn(x, y) { return x + y }
export PI = 3.14159`,
      mainContent: `import { add, PI } from "./math_module"
result = add(10, 5)
print("Result:", result)
print("PI:", PI)`,
      expected: "Result: 15\nPI: 3.14159",
    },
    {
      name: "Multiple Function Exports",
      moduleContent: `export multiply = fn(x, y) { return x * y }
export subtract = fn(x, y) { return x - y }
export MAX_VALUE = 100`,
      mainContent: `import { multiply, subtract, MAX_VALUE } from "./calc_module"
prod = multiply(6, 7)
diff = subtract(50, 8)
print("Product:", prod)
print("Difference:", diff)
print("Max:", MAX_VALUE)`,
      expected: "Product: 42\nDifference: 42\nMax: 100",
    },
    {
      name: "Higher-Order Function Export",
      moduleContent: `export apply = fn(f, x) { return f(x) }
export square = fn(x) { return x * x }`,
      mainContent: `import { apply, square } from "./functional_module"
result = apply(square, 8)
print("Applied square:", result)`,
      expected: "Applied square: 64",
    },
    {
      name: "Module with Complex Logic",
      moduleContent: `export factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
export fibonacci = fn(n) {
  if (n <= 1) {
    return n
  } else {
    return fibonacci(n - 1) + fibonacci(n - 2)
  }
}`,
      mainContent: `import { factorial, fibonacci } from "./recursive_module"
fact5 = factorial(5)
fib7 = fibonacci(7)
print("Factorial 5:", fact5)
print("Fibonacci 7:", fib7)`,
      expected: "Factorial 5: 120\nFibonacci 7: 13",
    },
    {
      name: "Array and String Operations Module",
      moduleContent: `export getFirst = fn(arr) { return arr[0] }
export getLast = fn(arr) { return arr[len(arr) - 1] }
export concat = fn(str1, str2) { return str1 + str2 }`,
      mainContent: `import { getFirst, getLast, concat } from "./utils_module"
numbers = [10, 20, 30, 40, 50]
first = getFirst(numbers)
last = getLast(numbers)
message = concat("Hello", " World")
print("First:", first)
print("Last:", last)
print("Message:", message)`,
      expected: "First: 10\nLast: 50\nMessage: Hello World",
    },
    {
      name: "Import Aliasing - Basic",
      moduleContent: `export add = fn(x, y) { return x + y }
export PI = 3.14159`,
      mainContent: `import { add as sum, PI as pi } from "./math_alias_module"
result = sum(10, 5)
print("Result:", result)
print("Pi:", pi)`,
      expected: "Result: 15\nPi: 3.14159",
    },
    {
      name: "Import Aliasing - Mixed with Regular Imports",
      moduleContent: `export multiply = fn(x, y) { return x * y }
export divide = fn(x, y) { return x / y }
export MAX_VALUE = 100`,
      mainContent: `import { multiply as mult, divide, MAX_VALUE as max } from "./mixed_alias_module"
product = mult(6, 7)
quotient = divide(50, 2)
print("Product:", product)
print("Quotient:", quotient)
print("Max:", max)`,
      expected: "Product: 42\nQuotient: 25\nMax: 100",
    },
    {
      name: "Import Aliasing - Avoiding Name Conflicts",
      moduleContent: `export print = fn(msg) { return "Module: " + msg }
export format = fn(str) { return "[" + str + "]" }`,
      mainContent: `import { print as modulePrint, format } from "./conflict_module"
result = modulePrint("Hello")
formatted = format("Test")
print("Module result:", result)
print("Formatted:", formatted)`,
      expected: "Module result: Module: Hello\nFormatted: [Test]",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Create temporary directory for the test
      tmpDir, err := os.MkdirTemp("", "module_integration_test")
      if err != nil {
        t.Fatal(err)
      }
      defer os.RemoveAll(tmpDir)

      // Determine module filename from test name
      var moduleFileName string
      switch tt.name {
      case "Basic Function Export and Import":
        moduleFileName = "math_module.rush"
      case "Multiple Function Exports":
        moduleFileName = "calc_module.rush"
      case "Higher-Order Function Export":
        moduleFileName = "functional_module.rush"
      case "Module with Complex Logic":
        moduleFileName = "recursive_module.rush"
      case "Array and String Operations Module":
        moduleFileName = "utils_module.rush"
      case "Import Aliasing - Basic":
        moduleFileName = "math_alias_module.rush"
      case "Import Aliasing - Mixed with Regular Imports":
        moduleFileName = "mixed_alias_module.rush"
      case "Import Aliasing - Avoiding Name Conflicts":
        moduleFileName = "conflict_module.rush"
      }

      // Create module file
      modulePath := filepath.Join(tmpDir, moduleFileName)
      if err := os.WriteFile(modulePath, []byte(tt.moduleContent), 0644); err != nil {
        t.Fatal(err)
      }

      // Create main file with adjusted import path using absolute module path
      moduleNameWithoutExt := strings.TrimSuffix(moduleFileName, ".rush")
      adjustedMainContent := strings.Replace(tt.mainContent, "./"+moduleNameWithoutExt, strings.TrimSuffix(modulePath, ".rush"), 1)
      
      mainPath := filepath.Join(tmpDir, "main.rush")
      if err := os.WriteFile(mainPath, []byte(adjustedMainContent), 0644); err != nil {
        t.Fatal(err)
      }

      // Run the main program using go run from project root
      cmd := exec.Command("go", "run", "cmd/rush/main.go", mainPath)
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      if err := cmd.Run(); err != nil {
        t.Fatalf("Program execution failed: %v\nOutput: %s", err, out.String())
      }

      fullOutput := strings.TrimSpace(out.String())
      
      // Extract just the printed content
      lines := strings.Split(fullOutput, "\n")
      var printedLines []string
      for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" && 
           !strings.HasPrefix(line, "Rush interpreter") && 
           !strings.HasPrefix(line, "Rush tree-walking interpreter") && 
           line != "Result: null" && 
           !strings.HasPrefix(line, "Execution complete!") {
          printedLines = append(printedLines, line)
        }
      }
      
      output := strings.Join(printedLines, "\n")
      if output != tt.expected {
        t.Errorf("Expected output %q, got %q (full output: %q)", tt.expected, output, fullOutput)
      }
    })
  }
}

func TestModuleErrorCases(t *testing.T) {
  errorTests := []struct {
    name         string
    moduleContent string
    mainContent   string
    hasError     bool
  }{
    {
      name: "Import Non-existent Module",
      moduleContent: "", // No module file created
      mainContent: `import { add } from "./nonexistent_module"
result = add(1, 2)`,
      hasError: true,
    },
    {
      name: "Import Non-existent Export",
      moduleContent: `export multiply = fn(x, y) { return x * y }`,
      mainContent: `import { nonexistentFunc } from "./test_module"
result = nonexistentFunc(1, 2)`,
      hasError: true,
    },
    {
      name: "Module with Syntax Error",
      moduleContent: `export broken = fn(x, y { return x + y`, // Missing closing parenthesis
      mainContent: `import { broken } from "./broken_module"
result = broken(1, 2)`,
      hasError: true,
    },
    {
      name: "Valid Module Import",
      moduleContent: `export working = fn(x) { return x * 2 }`,
      mainContent: `import { working } from "./valid_module"
result = working(21)
print("Result:", result)`,
      hasError: false,
    },
  }

  for _, tt := range errorTests {
    t.Run(tt.name, func(t *testing.T) {
      // Create temporary directory for the test
      tmpDir, err := os.MkdirTemp("", "module_error_test")
      if err != nil {
        t.Fatal(err)
      }
      defer os.RemoveAll(tmpDir)

      // Create module file only if content is provided
      if tt.moduleContent != "" {
        var moduleFileName string
        switch tt.name {
        case "Import Non-existent Export":
          moduleFileName = "test_module.rush"
        case "Module with Syntax Error":
          moduleFileName = "broken_module.rush"
        case "Valid Module Import":
          moduleFileName = "valid_module.rush"
        }

        if moduleFileName != "" {
          modulePath := filepath.Join(tmpDir, moduleFileName)
          if err := os.WriteFile(modulePath, []byte(tt.moduleContent), 0644); err != nil {
            t.Fatal(err)
          }
        }
      }

      // Create main file with adjusted import paths
      adjustedMainContent := tt.mainContent
      // Replace relative import paths with absolute paths for existing modules
      if tt.moduleContent != "" {
        var moduleFileName string
        switch tt.name {
        case "Import Non-existent Export":
          moduleFileName = "test_module.rush"
          modulePath := filepath.Join(tmpDir, moduleFileName)
          adjustedMainContent = strings.Replace(adjustedMainContent, "./test_module", strings.TrimSuffix(modulePath, ".rush"), 1)
        case "Module with Syntax Error":
          moduleFileName = "broken_module.rush"
          modulePath := filepath.Join(tmpDir, moduleFileName)
          adjustedMainContent = strings.Replace(adjustedMainContent, "./broken_module", strings.TrimSuffix(modulePath, ".rush"), 1)
        case "Valid Module Import":
          moduleFileName = "valid_module.rush"
          modulePath := filepath.Join(tmpDir, moduleFileName)
          adjustedMainContent = strings.Replace(adjustedMainContent, "./valid_module", strings.TrimSuffix(modulePath, ".rush"), 1)
        }
      }
      
      mainPath := filepath.Join(tmpDir, "main.rush")
      if err := os.WriteFile(mainPath, []byte(adjustedMainContent), 0644); err != nil {
        t.Fatal(err)
      }

      // Run the main program using go run from project root
      cmd := exec.Command("go", "run", "cmd/rush/main.go", mainPath)
      var out bytes.Buffer
      cmd.Stdout = &out
      cmd.Stderr = &out

      err = cmd.Run()
      
      if tt.hasError && err == nil {
        t.Errorf("Expected error but program executed successfully. Output: %s", out.String())
      } else if !tt.hasError && err != nil {
        t.Errorf("Expected success but got error: %v\nOutput: %s", err, out.String())
      }
    })
  }
}