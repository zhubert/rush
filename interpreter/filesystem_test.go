package interpreter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemConstructors(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`file("test.txt")`, "#<File:test.txt (closed)>"},
		{`directory("testdir")`, "#<Directory:testdir>"},
		{`path("/tmp/test")`, "#<Path:/tmp/test>"},
		{`file(123)`, "argument to `file` must be STRING, got INTEGER"},
		{`directory(true)`, "argument to `directory` must be STRING, got BOOLEAN"},
		{`path([])`, "argument to `path` must be STRING, got ARRAY"},
		{`file("../etc/passwd")`, "invalid file path: path traversal not allowed"},
		{`directory("../")`, "invalid directory path: path traversal not allowed"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case string:
			if evaluated.Inspect() != expected {
				// Check if it's an error
				if errObj, ok := evaluated.(*Error); ok {
					if errObj.Message != expected {
						t.Errorf("wrong error message. expected=%q, got=%q",
							expected, errObj.Message)
					}
				} else {
					t.Errorf("wrong result. expected=%q, got=%q",
						expected, evaluated.Inspect())
				}
			}
		}
	}
}

func TestFileOperations(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "rush_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, Rush!"

	// Test 1: Check file doesn't exist initially
	t.Run("file exists check - nonexistent", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").exists?()`)
		testBooleanObject(t, evaluated, false)
	})

	// Test 2: Create and write to file
	t.Run("create and write file", func(t *testing.T) {
		// Write content to file using actual file operations
		if err := ioutil.WriteFile(testFile, []byte(testContent), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	})

	// Test 3: Check file exists after creation
	t.Run("file exists check - existing", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").exists?()`)
		testBooleanObject(t, evaluated, true)
	})

	// Test 4: Get file size
	t.Run("get file size", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").size()`)
		testIntegerObject(t, evaluated, int64(len(testContent)))
	})

	// Test 5: Read file content
	t.Run("read file content", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").open("r").read()`)
		if str, ok := evaluated.(*String); ok {
			if str.Value != testContent {
				t.Errorf("wrong file content. expected=%q, got=%q", testContent, str.Value)
			}
		} else {
			t.Errorf("failed to read file content, got %T", evaluated)
		}
	})

	// Test 6: File path property
	t.Run("file path property", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").path`)
		if str, ok := evaluated.(*String); ok {
			if str.Value != testFile {
				t.Errorf("wrong file path. expected=%q, got=%q", testFile, str.Value)
			}
		} else {
			t.Errorf("path property should return string, got %T", evaluated)
		}
	})

	// Test 7: File is_open property when closed
	t.Run("file is_open property when closed", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").is_open`)
		testBooleanObject(t, evaluated, false)
	})

	// Test 8: Delete file
	t.Run("delete file", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").delete()`)
		testBooleanObject(t, evaluated, true)
	})

	// Test 9: Check file doesn't exist after deletion
	t.Run("file exists check after deletion", func(t *testing.T) {
		evaluated := testEval(`file("` + testFile + `").exists?()`)
		testBooleanObject(t, evaluated, false)
	})
}

func TestFileErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "read closed file",
			input:    `file("nonexistent.txt").read()`,
			expected: "file is not open: nonexistent.txt",
		},
		{
			name:     "write to closed file",
			input:    `file("test.txt").write("content")`,
			expected: "file is not open: test.txt",
		},
		{
			name:     "close already closed file",
			input:    `file("test.txt").close()`,
			expected: "file is not open: test.txt",
		},
		{
			name:     "size of nonexistent file",
			input:    `file("nonexistent.txt").size()`,
			expected: "file does not exist: nonexistent.txt",
		},
		{
			name:     "invalid file mode",
			input:    `file("test.txt").open("x")`,
			expected: "invalid file mode: x",
		},
		{
			name:     "write wrong argument type",
			input:    `file("test.txt").open("w").write(123)`,
			expected: "file content argument must be STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("expected error object, got %T", evaluated)
				return
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expected, errObj.Message)
			}
		})
	}
}

func TestDirectoryOperations(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "rush_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "testsubdir")

	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "directory exists check - nonexistent",
			input:    `directory("` + testDir + `").exists?()`,
			expected: false,
		},
		{
			name:     "create directory",
			input:    `directory("` + testDir + `").create()`,
			expected: "#<Directory:" + testDir + ">",
		},
		{
			name:     "directory exists check - existing",
			input:    `directory("` + testDir + `").exists?()`,
			expected: true,
		},
		{
			name:     "directory path property",
			input:    `directory("` + testDir + `").path`,
			expected: testDir,
		},
		{
			name:     "list empty directory",
			input:    `directory("` + testDir + `").list()`,
			expected: "[]",
		},
		{
			name:     "delete directory",
			input:    `directory("` + testDir + `").delete()`,
			expected: true,
		},
		{
			name:     "directory exists check after deletion",
			input:    `directory("` + testDir + `").exists?()`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case bool:
				testBooleanObject(t, evaluated, expected)
			case string:
				if evaluated.Inspect() != expected {
					t.Errorf("wrong result. expected=%q, got=%q",
						expected, evaluated.Inspect())
				}
			}
		})
	}
}

func TestPathOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "path value property",
			input:    `path("/tmp/test").value`,
			expected: "/tmp/test",
		},
		{
			name:     "path join",
			input:    `path("/tmp").join("test.txt")`,
			expected: "#<Path:/tmp/test.txt>",
		},
		{
			name:     "path basename",
			input:    `path("/tmp/test.txt").basename()`,
			expected: "test.txt",
		},
		{
			name:     "path dirname",
			input:    `path("/tmp/test.txt").dirname()`,
			expected: "/tmp",
		},
		{
			name:     "path clean",
			input:    `path("/tmp//test/../file.txt").clean()`,
			expected: "#<Path:/tmp/file.txt>",
		},
		{
			name:     "path absolute from relative",
			input:    `path("test.txt").absolute()`,
			expected: "starts with /", // Will be absolute path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case string:
				if expected == "starts with /" {
					// Special case for absolute path test
					if path, ok := evaluated.(*Path); ok {
						if !filepath.IsAbs(path.Value) {
							t.Errorf("expected absolute path, got %q", path.Value)
						}
					} else {
						t.Errorf("expected Path object, got %T", evaluated)
					}
				} else if str, ok := evaluated.(*String); ok {
					if str.Value != expected {
						t.Errorf("wrong string value. expected=%q, got=%q",
							expected, str.Value)
					}
				} else {
					if evaluated.Inspect() != expected {
						t.Errorf("wrong result. expected=%q, got=%q",
							expected, evaluated.Inspect())
					}
				}
			}
		})
	}
}

func TestPathErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path join wrong argument type",
			input:    `path("/tmp").join(123)`,
			expected: "path join argument must be STRING",
		},
		{
			name:     "path join wrong number of arguments",
			input:    `path("/tmp").join("a", "b")`,
			expected: "wrong number of arguments for path.join: want=1, got=2",
		},
		{
			name:     "path basename with arguments",
			input:    `path("/tmp").basename("extra")`,
			expected: "wrong number of arguments for path.basename: want=0, got=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("expected error object, got %T", evaluated)
				return
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expected, errObj.Message)
			}
		})
	}
}

func TestFileSystemMethodChaining(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "rush_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path method chaining",
			input:    `path("/tmp").join("test").join("file.txt").clean()`,
			expected: "#<Path:/tmp/test/file.txt>",
		},
		{
			name:     "path to string conversion",
			input:    `path("/tmp").join("test.txt").basename()`,
			expected: "test.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			if str, ok := evaluated.(*String); ok {
				if str.Value != tt.expected {
					t.Errorf("wrong string value. expected=%q, got=%q",
						tt.expected, str.Value)
				}
			} else {
				if evaluated.Inspect() != tt.expected {
					t.Errorf("wrong result. expected=%q, got=%q",
						tt.expected, evaluated.Inspect())
				}
			}
		})
	}
}

