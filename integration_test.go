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