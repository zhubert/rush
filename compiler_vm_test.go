package main

import (
	"testing"

	"rush/ast"
	"rush/compiler"
	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
	"rush/vm"
)

type integrationTestCase struct {
	input    string
	expected interface{}
}

func TestCompilerVMIntegration(t *testing.T) {
	tests := []integrationTestCase{
		// Basic arithmetic
		{"1 + 2 * 3", 7},
		{"(1 + 2) * 3", 9},
		{"10 - 5 + 3", 8},
		{"20 / 4 * 2", 10},
		{"15 % 7", 1},

		// Boolean operations
		{"true && false", false},
		{"true || false", true},
		{"!true", false},
		{"!!true", true},
		{"1 == 1", true},
		{"1 != 2", true},
		{"3 > 2", true},
		{"2 < 3", true},
		{"3 >= 3", true},
		{"2 <= 3", true},

		// String operations
		{`"hello" + " " + "world"`, "hello world"},
		{`"test"`, "test"},

		// Array operations
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[]", []int{}},
		{"[1, 2, 3][1]", 2},
		{"[1 + 1, 2 * 2, 3 - 1][2]", 2},

		// Hash operations
		{"{1: 2, 3: 4}[1]", 2},
		{"{1: 2, 3: 4}[3]", 4},
		{"{1: 2}[2]", interpreter.NULL},

		// Variables
		{"a = 5; a", 5},
		{"a = 5; b = 10; a + b", 15},
		{"a = 1; a = a + 1; a", 2},

		// Conditionals
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", interpreter.NULL},

		// Functions
		{"fn(x) { x * 2 }(5)", 10},
		{"add = fn(a, b) { a + b }; add(5, 3)", 8},
		{"fn() { 42 }()", 42},

		// Recursion
		{`
		factorial = fn(n) {
			if (n <= 1) {
				1
			} else {
				n * factorial(n - 1)
			}
		};
		factorial(5);
		`, 120},

		// Closures
		{`
		newAdder = fn(x) {
			fn(y) { x + y }
		};
		addTwo = newAdder(2);
		addTwo(3);
		`, 5},

		// Complex example
		{`
		counter = fn() {
			count = 0;
			fn() {
				count = count + 1;
				count
			}
		};
		c = counter();
		c() + c() + c();
		`, 6},

		// Built-in functions
		{`len("hello")`, 5},
		{`len([1, 2, 3, 4])`, 4},
		{`push([1, 2], 3)`, []int{1, 2, 3}},
	}

	runIntegrationTests(t, tests)
}

func TestComplexPrograms(t *testing.T) {
	tests := []integrationTestCase{
		// Fibonacci with memoization
		{`
		fibonacci = fn(n) {
			if (n == 0) { 
				0 
			} else { 
				if (n == 1) { 
					1 
				} else { 
					fibonacci(n-1) + fibonacci(n-2) 
				} 
			}
		};
		fibonacci(10);
		`, 55},


		// Higher order functions
		{`
		twice = fn(f, x) {
			f(f(x))
		};
		addOne = fn(x) { x + 1 };
		twice(addOne, 5);
		`, 7},

		// Nested scopes and closures
		{`
		makeMultiplier = fn(x) {
			fn(y) {
				fn(z) {
					x * y * z
				}
			}
		};
		multiplyByTwo = makeMultiplier(2);
		multiplyByTwoAndThree = multiplyByTwo(3);
		multiplyByTwoAndThree(4);
		`, 24},
	}

	runIntegrationTests(t, tests)
}

func TestCompilerVMErrorHandling(t *testing.T) {
	errorTests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true",
			"unknown operator: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"unknown operator: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"undefined variable foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Rush"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range errorTests {
		program := parseProgram(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			if err.Error() != tt.expected {
				t.Errorf("wrong compiler error. want=%q, got=%q",
					tt.expected, err.Error())
			}
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err == nil {
			t.Errorf("expected VM error but resulted in none.")
			continue
		}

		if err.Error() != tt.expected {
			t.Errorf("wrong VM error. want=%q, got=%q",
				tt.expected, err.Error())
		}
	}
}

func TestIndexAssignmentIntegration(t *testing.T) {
	tests := []integrationTestCase{
		{`
		arr = [1, 2, 3];
		arr[1] = 5;
		arr;
		`, []int{1, 5, 3}},

		{`
		hash = {1: "one", 2: "two"};
		hash[2] = "TWO";
		hash[1];
		`, "one"},

		{`
		arr = [1, 2, 3];
		arr[0] = arr[0] + 10;
		arr[0];
		`, 11},
	}

	runIntegrationTests(t, tests)
}

func TestWhileLoopIntegration(t *testing.T) {
	tests := []integrationTestCase{
		{`
		i = 0;
		sum = 0;
		while (i < 5) {
			sum = sum + i;
			i = i + 1;
		}
		sum;
		`, 10},

		{`
		count = 0;
		while (count < 3) {
			count = count + 1;
		}
		count;
		`, 3},
	}

	runIntegrationTests(t, tests)
}

func TestForLoopIntegration(t *testing.T) {
	tests := []integrationTestCase{
		{`
		sum = 0;
		for (i = 0; i < 5; i = i + 1) {
			sum = sum + i;
		}
		sum;
		`, 10},

		{`
		product = 1;
		for (i = 1; i <= 4; i = i + 1) {
			product = product * i;
		}
		product;
		`, 24},
	}

	runIntegrationTests(t, tests)
}

func TestPropertyAccessIntegration(t *testing.T) {
	tests := []integrationTestCase{
		{`"hello".length`, 5},
		{`[1, 2, 3, 4].length`, 4},
		{`[].length`, 0},
	}

	runIntegrationTests(t, tests)
}

func runIntegrationTests(t *testing.T, tests []integrationTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parseProgram(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := machine.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func parseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual interpreter.Value) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, int64(expected), actual)
	case bool:
		testBooleanObject(t, bool(expected), actual)
	case string:
		testStringObject(t, expected, actual)
	case []int:
		testIntegerArrayObject(t, expected, actual)
	case *interpreter.Null:
		if actual != interpreter.NULL {
			t.Errorf("object is not NULL: %T (%+v)", actual, actual)
		}
	default:
		// For NULL values passed as nil interface{}
		if expected == interpreter.NULL {
			if actual != interpreter.NULL {
				t.Errorf("object is not NULL: %T (%+v)", actual, actual)
			}
		}
	}
}

func testIntegerObject(t *testing.T, expected int64, actual interpreter.Value) {
	result, ok := actual.(*interpreter.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
}

func testBooleanObject(t *testing.T, expected bool, actual interpreter.Value) {
	result, ok := actual.(*interpreter.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
}

func testStringObject(t *testing.T, expected string, actual interpreter.Value) {
	result, ok := actual.(*interpreter.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", actual, actual)
		return
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}
}

func testIntegerArrayObject(t *testing.T, expected []int, actual interpreter.Value) {
	result, ok := actual.(*interpreter.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", actual, actual)
		return
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("wrong num of elements. want=%d, got=%d",
			len(expected), len(result.Elements))
		return
	}

	for i, expectedElem := range expected {
		testIntegerObject(t, int64(expectedElem), result.Elements[i])
	}
}

func BenchmarkCompilerVMArithmetic(b *testing.B) {
	input := "1 + 2 * 3 - 4 / 2"
	program := parseProgram(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			b.Fatal(err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompilerVMFibonacci(b *testing.B) {
	input := `
	fibonacci = fn(n) {
		if (n == 0) { 0 }
		else if (n == 1) { 1 }
		else { fibonacci(n-1) + fibonacci(n-2) }
	};
	fibonacci(20);
	`
	program := parseProgram(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			b.Fatal(err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompilerVMComplexProgram(b *testing.B) {
	input := `
	reduce = fn(arr, initial, f) {
		iter = fn(arr, result) {
			if (len(arr) == 0) {
				result
			} else {
				iter(rest(arr), f(result, first(arr)));
			}
		};
		iter(arr, initial);
	};
	
	sum = fn(arr) {
		add = fn(a, b) { a + b };
		reduce(arr, 0, add);
	};
	
	arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
	sum(arr);
	`
	program := parseProgram(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			b.Fatal(err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}