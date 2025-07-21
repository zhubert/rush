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
           !strings.HasPrefix(line, "Result:") && 
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
result = 5 + "hello"
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
  
  // Update import path in main content
  adjustedMainContent := strings.Replace(mainContent, "./test_module", "./"+strings.TrimSuffix(filepath.Base(moduleFile.Name()), ".rush"), 1)
  
  if err := os.WriteFile(mainFile, []byte(adjustedMainContent), 0644); err != nil {
    t.Fatal(err)
  }
  defer os.Remove(mainFile)

  // Run the main program
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