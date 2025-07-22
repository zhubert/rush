package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestArrayOperationImplementations(t *testing.T) {
	tests := []struct {
		name     string
		program  string
		expected string
	}{
		{
			name: "Map function implementation",
			program: `
numbers = [1, 2, 3]
map_fn = fn(arr, func) {
  result = []
  i = 0
  while (i < len(arr)) {
    result = push(result, func(arr[i]))
    i = i + 1
  }
  return result
}
doubled = map_fn(numbers, fn(x) { return x * 2 })
print(doubled)
`,
			expected: "[2, 4, 6]",
		},
		{
			name: "Filter function implementation",
			program: `
numbers = [1, 2, 3, 4, 5]
filter_fn = fn(arr, func) {
  result = []
  i = 0
  while (i < len(arr)) {
    if (func(arr[i])) {
      result = push(result, arr[i])
    }
    i = i + 1
  }
  return result
}
evens = filter_fn(numbers, fn(x) { return x % 2 == 0 })
print(evens)
`,
			expected: "[2, 4]",
		},
		{
			name: "Reduce function implementation",
			program: `
numbers = [1, 2, 3, 4]
reduce_fn = fn(arr, func, initial) {
  result = initial
  i = 0
  while (i < len(arr)) {
    result = func(result, arr[i])
    i = i + 1
  }
  return result
}
sum = reduce_fn(numbers, fn(acc, x) { return acc + x }, 0)
print(sum)
`,
			expected: "10",
		},
		{
			name: "Index of function implementation",
			program: `
numbers = [10, 20, 30]
index_of_fn = fn(arr, element) {
  i = 0
  while (i < len(arr)) {
    if (arr[i] == element) {
      return i
    }
    i = i + 1
  }
  return -1
}
index = index_of_fn(numbers, 20)
print(index)
`,
			expected: "1",
		},
		{
			name: "Includes function implementation",
			program: `
numbers = [10, 20, 30]
index_of_fn = fn(arr, element) {
  i = 0
  while (i < len(arr)) {
    if (arr[i] == element) {
      return i
    }
    i = i + 1
  }
  return -1
}
includes_fn = fn(arr, element) {
  return index_of_fn(arr, element) != -1
}
has_twenty = includes_fn(numbers, 20)
print(has_twenty)
`,
			expected: "true",
		},
		{
			name: "Reverse function implementation",
			program: `
numbers = [1, 2, 3]
reverse_fn = fn(arr) {
  result = []
  i = len(arr) - 1
  while (i >= 0) {
    result = push(result, arr[i])
    i = i - 1
  }
  return result
}
reversed = reverse_fn(numbers)
print(reversed)
`,
			expected: "[3, 2, 1]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a temporary file with the program
			tmpFile, err := os.CreateTemp("", "test*.rush")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(test.program); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Run the program
			cmd := exec.Command("go", "run", "cmd/rush/main.go", tmpFile.Name())
			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = &output

			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to run program: %v\nOutput: %s", err, output.String())
			}

			// Extract just the printed result (skip interpreter messages)
			lines := strings.Split(strings.TrimSpace(output.String()), "\n")
			var actual string
			for _, line := range lines {
				if !strings.Contains(line, "Rush interpreter") && !strings.Contains(line, "Execution complete") && !strings.Contains(line, "Result:") {
					if strings.TrimSpace(line) != "" {
						actual = strings.TrimSpace(line)
						break
					}
				}
			}

			if actual != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, actual)
			}
		})
	}
}