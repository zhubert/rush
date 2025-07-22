package interpreter

import (
  "testing"
  "rush/lexer"
  "rush/parser"
)

func TestBuiltinFunctions(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {`len("")`, 0},
    {`len("four")`, 4},
    {`len("hello world")`, 11},
    {`len([1, 2, 3])`, 3},
    {`len([])`, 0},
    {`len([1, 2, 3, 4, 5])`, 5},
    {"len(1)", "argument to `len` not supported, got INTEGER"},
    {"len(\"one\", \"two\")", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestTypeFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"type(42)", "INTEGER"},
    {"type(3.14)", "FLOAT"},
    {"type(true)", "BOOLEAN"},
    {"type(false)", "BOOLEAN"},
    {"type(\"hello\")", "STRING"},
    {"type([1, 2, 3])", "ARRAY"},
    {"type(fn(x) { x })", "FUNCTION"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)
    testStringObject(t, evaluated, tt.expected)
  }
}

func TestPrintFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"print(\"hello\")", nil},
    {"print(42)", nil},
    {"print(true)", nil},
    {"print([1, 2, 3])", nil},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case nil:
      testNullObject(t, evaluated)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestSubstrFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"substr(\"hello\", 0, 5)", "hello"},
    {"substr(\"hello\", 1, 3)", "ell"},
    {"substr(\"hello\", 0, 2)", "he"},
    {"substr(\"hello\", 2, 3)", "llo"},
    {"substr(\"hello\", 10, 1)", ""},
    {"substr(\"hello\", -1, 1)", ""},
    {"substr(\"hello\", 0, 10)", "hello"},
    {"substr(42, 0, 1)", "argument to `substr` must be STRING, got INTEGER"},
    {"substr(\"hello\")", "wrong number of arguments. got=1, want=3"},
    {"substr(\"hello\", 0)", "wrong number of arguments. got=2, want=3"},
    {"substr(\"hello\", 0, 1, 2)", "wrong number of arguments. got=4, want=3"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case string:
      if errObj, ok := evaluated.(*Error); ok {
        if errObj.Message != expected {
          t.Errorf("wrong error message. expected=%q, got=%q",
            expected, errObj.Message)
        }
      } else {
        testStringObject(t, evaluated, expected)
      }
    }
  }
}

func TestSplitFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"split(\"hello,world\", \",\")", []string{"hello", "world"}},
    {"split(\"a b c\", \" \")", []string{"a", "b", "c"}},
    {"split(\"hello\", \"\")", []string{"h", "e", "l", "l", "o"}},
    {"split(\"\", \",\")", []string{""}},
    {"split(\"hello\", \"xyz\")", []string{"hello"}},
    {"split(42, \",\")", "first argument to `split` must be STRING, got INTEGER"},
    {"split(\"hello\", 42)", "second argument to `split` must be STRING, got INTEGER"},
    {"split(\"hello\")", "wrong number of arguments. got=1, want=2"},
    {"split(\"hello\", \",\", \"extra\")", "wrong number of arguments. got=3, want=2"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case []string:
      array, ok := evaluated.(*Array)
      if !ok {
        t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if len(array.Elements) != len(expected) {
        t.Errorf("wrong array length. expected=%d, got=%d",
          len(expected), len(array.Elements))
        continue
      }
      for i, expectedStr := range expected {
        testStringObject(t, array.Elements[i], expectedStr)
      }
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestPushFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"push([1, 2, 3], 4)", []int{1, 2, 3, 4}},
    {"push([], 1)", []int{1}},
    {"push([\"a\", \"b\"], \"c\")", []string{"a", "b", "c"}},
    {"push(42, 1)", "first argument to `push` must be ARRAY, got INTEGER"},
    {"push([1, 2, 3])", "wrong number of arguments. got=1, want=2"},
    {"push([1, 2, 3], 4, 5)", "wrong number of arguments. got=3, want=2"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case []int:
      array, ok := evaluated.(*Array)
      if !ok {
        t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if len(array.Elements) != len(expected) {
        t.Errorf("wrong array length. expected=%d, got=%d",
          len(expected), len(array.Elements))
        continue
      }
      for i, expectedInt := range expected {
        testIntegerObject(t, array.Elements[i], int64(expectedInt))
      }
    case []string:
      array, ok := evaluated.(*Array)
      if !ok {
        t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if len(array.Elements) != len(expected) {
        t.Errorf("wrong array length. expected=%d, got=%d",
          len(expected), len(array.Elements))
        continue
      }
      for i, expectedStr := range expected {
        testStringObject(t, array.Elements[i], expectedStr)
      }
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestPopFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"pop([1, 2, 3])", 3},
    {"pop([42])", 42},
    {"pop([\"hello\", \"world\"])", "world"},
    {"try { pop([]) } catch (IndexError error) { if (false) { 1 } }", nil},
    {"pop(42)", "argument to `pop` must be ARRAY, got INTEGER"},
    {"pop([1, 2, 3], 1)", "wrong number of arguments. got=2, want=1"},
    {"pop()", "wrong number of arguments. got=0, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      if errObj, ok := evaluated.(*Error); ok {
        if errObj.Message != expected {
          t.Errorf("wrong error message. expected=%q, got=%q",
            expected, errObj.Message)
        }
      } else {
        testStringObject(t, evaluated, expected)
      }
    case nil:
      testNullObject(t, evaluated)
    }
  }
}

func TestSliceFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"slice([1, 2, 3, 4, 5], 1, 4)", []int{2, 3, 4}},
    {"slice([1, 2, 3], 0, 2)", []int{1, 2}},
    {"slice([1, 2, 3], 1, 3)", []int{2, 3}},
    {"slice([1, 2, 3], 0, 10)", []int{1, 2, 3}},
    {"slice([1, 2, 3], 5, 10)", []int{}},
    {"slice([1, 2, 3], -1, 2)", []int{}},
    {"slice(42, 0, 1)", "first argument to `slice` must be ARRAY, got INTEGER"},
    {"slice([1, 2, 3])", "wrong number of arguments. got=1, want=3"},
    {"slice([1, 2, 3], 0)", "wrong number of arguments. got=2, want=3"},
    {"slice([1, 2, 3], 0, 1, 2)", "wrong number of arguments. got=4, want=3"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case []int:
      array, ok := evaluated.(*Array)
      if !ok {
        t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if len(array.Elements) != len(expected) {
        t.Errorf("wrong array length. expected=%d, got=%d",
          len(expected), len(array.Elements))
        continue
      }
      for i, expectedInt := range expected {
        testIntegerObject(t, array.Elements[i], int64(expectedInt))
      }
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestOrdFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {`ord("a")`, 97},
    {`ord("A")`, 65},
    {`ord("z")`, 122},
    {`ord("Z")`, 90},
    {`ord("0")`, 48},
    {`ord("9")`, 57},
    {`ord(" ")`, 32},
    {`ord("!")`, 33},
    {`ord("")`, "argument to `ord` must be a single character, got length 0"},
    {`ord("ab")`, "argument to `ord` must be a single character, got length 2"},
    {`ord(42)`, "argument to `ord` must be STRING, got INTEGER"},
    {`ord("a", "b")`, "wrong number of arguments. got=2, want=1"},
    {`ord()`, "wrong number of arguments. got=0, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestChrFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {`chr(97)`, "a"},
    {`chr(65)`, "A"},
    {`chr(122)`, "z"},
    {`chr(90)`, "Z"},
    {`chr(48)`, "0"},
    {`chr(57)`, "9"},
    {`chr(32)`, " "},
    {`chr(33)`, "!"},
    {`chr(127)`, string(byte(127))},
    {`chr(0)`, string(byte(0))},
    {`chr(-1)`, "argument to `chr` must be between 0 and 127, got -1"},
    {`chr(128)`, "argument to `chr` must be between 0 and 127, got 128"},
    {`chr("a")`, "argument to `chr` must be INTEGER, got STRING"},
    {`chr(65, 66)`, "wrong number of arguments. got=2, want=1"},
    {`chr()`, "wrong number of arguments. got=0, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case string:
      if errObj, ok := evaluated.(*Error); ok {
        if errObj.Message != expected {
          t.Errorf("wrong error message. expected=%q, got=%q",
            expected, errObj.Message)
        }
      } else {
        testStringObject(t, evaluated, expected)
      }
    }
  }
}

// Helper function specifically for builtin tests
func testEvalBuiltin(input string) Value {
  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  return Eval(program, env)
}

func TestToStringFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {`to_string(42)`, "42"},
    {`to_string(3.14)`, "3.14"},
    {`to_string(true)`, "true"},
    {`to_string(false)`, "false"},
    {`to_string("hello")`, "hello"},
    {`to_string([1, 2, 3])`, "[1, 2, 3]"},
    {`to_string([])`, "[]"},
    {`to_string(0)`, "0"},
    {`to_string(-42)`, "-42"},
    {`to_string(3.0)`, "3"},
    {`to_string([true, "test", 42])`, `[true, test, 42]`},
  }

  for _, tt := range tests {
    t.Run(tt.input, func(t *testing.T) {
      result := testEvalBuiltin(tt.input)
      
      if result.Type() != STRING_VALUE {
        t.Fatalf("Expected STRING, got %T (%+v)", result, result)
      }
      
      str, ok := result.(*String)
      if !ok {
        t.Fatalf("Expected *String, got %T", result)
      }
      
      if str.Value != tt.expected {
        t.Errorf("Expected %q, got %q", tt.expected, str.Value)
      }
    })
  }
}

func TestToStringErrors(t *testing.T) {
  tests := []struct {
    input    string
    hasError bool
  }{
    {`to_string()`, true},           // No arguments
    {`to_string(1, 2)`, true},      // Too many arguments
  }

  for _, tt := range tests {
    t.Run(tt.input, func(t *testing.T) {
      result := testEvalBuiltin(tt.input)
      
      if tt.hasError {
        if result.Type() != ERROR_VALUE {
          t.Fatalf("Expected ERROR, got %T (%+v)", result, result)
        }
      } else {
        if result.Type() == ERROR_VALUE {
          t.Fatalf("Unexpected error: %s", result.Inspect())
        }
      }
    })
  }
}

// Math function tests

func TestBuiltinAbsFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_abs(5)", 5},
    {"builtin_abs(-5)", 5},
    {"builtin_abs(0)", 0},
    {"builtin_abs(3.14)", 3.14},
    {"builtin_abs(-3.14)", 3.14},
    {"builtin_abs(0.0)", 0.0},
    {"builtin_abs(\"hello\")", "argument to `builtin_abs` must be INTEGER or FLOAT, got STRING"},
    {"builtin_abs()", "wrong number of arguments. got=0, want=1"},
    {"builtin_abs(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinMinFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_min(5, 3)", 3},
    {"builtin_min(3, 5)", 3},
    {"builtin_min(-5, -3)", -5},
    {"builtin_min(0, 1)", 0},
    {"builtin_min(1, 2, 3)", 1},
    {"builtin_min(5, 3, 7, 1, 9)", 1},
    {"builtin_min(3.14, 2.71)", 2.71},
    {"builtin_min(5, 3.14)", 3.14},
    {"builtin_min(3.14, 5)", 3.14},
    {"builtin_min(42)", 42},
    {"builtin_min()", "wrong number of arguments. got=0, want at least 1"},
    {"builtin_min(\"hello\", 5)", "arguments to `builtin_min` must be INTEGER or FLOAT, got STRING"},
    {"builtin_min(5, \"hello\")", "arguments to `builtin_min` must be INTEGER or FLOAT, got STRING"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinMaxFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_max(5, 3)", 5},
    {"builtin_max(3, 5)", 5},
    {"builtin_max(-5, -3)", -3},
    {"builtin_max(0, 1)", 1},
    {"builtin_max(1, 2, 3)", 3},
    {"builtin_max(5, 3, 7, 1, 9)", 9},
    {"builtin_max(3.14, 2.71)", 3.14},
    {"builtin_max(5, 3.14)", 5.0},
    {"builtin_max(3.14, 5)", 5.0},
    {"builtin_max(42)", 42},
    {"builtin_max()", "wrong number of arguments. got=0, want at least 1"},
    {"builtin_max(\"hello\", 5)", "arguments to `builtin_max` must be INTEGER or FLOAT, got STRING"},
    {"builtin_max(5, \"hello\")", "arguments to `builtin_max` must be INTEGER or FLOAT, got STRING"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinFloorFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_floor(5)", 5},
    {"builtin_floor(5.9)", 5},
    {"builtin_floor(5.1)", 5},
    {"builtin_floor(-5.1)", -6},
    {"builtin_floor(-5.9)", -6},
    {"builtin_floor(0.0)", 0},
    {"builtin_floor(0.9)", 0},
    {"builtin_floor(-0.9)", -1},
    {"builtin_floor(\"hello\")", "argument to `builtin_floor` must be INTEGER or FLOAT, got STRING"},
    {"builtin_floor()", "wrong number of arguments. got=0, want=1"},
    {"builtin_floor(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinCeilFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_ceil(5)", 5},
    {"builtin_ceil(5.1)", 6},
    {"builtin_ceil(5.9)", 6},
    {"builtin_ceil(-5.1)", -5},
    {"builtin_ceil(-5.9)", -5},
    {"builtin_ceil(0.0)", 0},
    {"builtin_ceil(0.1)", 1},
    {"builtin_ceil(-0.1)", 0},
    {"builtin_ceil(\"hello\")", "argument to `builtin_ceil` must be INTEGER or FLOAT, got STRING"},
    {"builtin_ceil()", "wrong number of arguments. got=0, want=1"},
    {"builtin_ceil(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinRoundFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_round(5)", 5},
    {"builtin_round(5.4)", 5},
    {"builtin_round(5.5)", 6},
    {"builtin_round(5.6)", 6},
    {"builtin_round(-5.4)", -5},
    {"builtin_round(-5.5)", -6},
    {"builtin_round(-5.6)", -6},
    {"builtin_round(0.0)", 0},
    {"builtin_round(0.4)", 0},
    {"builtin_round(0.5)", 1},
    {"builtin_round(-0.5)", -1},
    {"builtin_round(\"hello\")", "argument to `builtin_round` must be INTEGER or FLOAT, got STRING"},
    {"builtin_round()", "wrong number of arguments. got=0, want=1"},
    {"builtin_round(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinSqrtFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_sqrt(4)", 2.0},
    {"builtin_sqrt(9)", 3.0},
    {"builtin_sqrt(16.0)", 4.0},
    {"builtin_sqrt(0)", 0.0},
    {"builtin_sqrt(0.0)", 0.0},
    {"builtin_sqrt(1)", 1.0},
    {"builtin_sqrt(2)", 1.4142135623730951}, // math.Sqrt(2)
    {"builtin_sqrt(-1)", "argument to `builtin_sqrt` cannot be negative, got -1.000000"},
    {"builtin_sqrt(-4.0)", "argument to `builtin_sqrt` cannot be negative, got -4.000000"},
    {"builtin_sqrt(\"hello\")", "argument to `builtin_sqrt` must be INTEGER or FLOAT, got STRING"},
    {"builtin_sqrt()", "wrong number of arguments. got=0, want=1"},
    {"builtin_sqrt(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinPowFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_pow(2, 3)", 8.0},
    {"builtin_pow(5, 0)", 1.0},
    {"builtin_pow(2, -1)", 0.5},
    {"builtin_pow(9, 0.5)", 3.0},
    {"builtin_pow(2.0, 3.0)", 8.0},
    {"builtin_pow(10, 2)", 100.0},
    {"builtin_pow(1, 100)", 1.0},
    {"builtin_pow(0, 5)", 0.0},
    {"builtin_pow(\"hello\", 2)", "first argument to `builtin_pow` must be INTEGER or FLOAT, got STRING"},
    {"builtin_pow(2, \"hello\")", "second argument to `builtin_pow` must be INTEGER or FLOAT, got STRING"},
    {"builtin_pow(2)", "wrong number of arguments. got=1, want=2"},
    {"builtin_pow()", "wrong number of arguments. got=0, want=2"},
    {"builtin_pow(1, 2, 3)", "wrong number of arguments. got=3, want=2"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinRandomFunction(t *testing.T) {
  tests := []struct {
    input    string
    hasError bool
    errorMsg string
  }{
    {"builtin_random()", false, ""},
    {"builtin_random(1)", true, "wrong number of arguments. got=1, want=0"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    if tt.hasError {
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != tt.errorMsg {
        t.Errorf("wrong error message. expected=%q, got=%q",
          tt.errorMsg, errObj.Message)
      }
    } else {
      // Check that it's a float between 0.0 and 1.0
      float, ok := evaluated.(*Float)
      if !ok {
        t.Errorf("object is not Float. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if float.Value < 0.0 || float.Value >= 1.0 {
        t.Errorf("random value out of range [0.0, 1.0). got=%f", float.Value)
      }
    }
  }
}

func TestBuiltinRandomIntFunction(t *testing.T) {
  tests := []struct {
    input    string
    hasError bool
    errorMsg string
    minVal   int64
    maxVal   int64
  }{
    {"builtin_random_int(1, 5)", false, "", 1, 5},
    {"builtin_random_int(0, 0)", false, "", 0, 0},
    {"builtin_random_int(-5, 5)", false, "", -5, 5},
    {"builtin_random_int(5, 1)", true, "first argument to `builtin_random_int` cannot be greater than second argument", 0, 0},
    {"builtin_random_int(\"hello\", 5)", true, "first argument to `builtin_random_int` must be INTEGER, got STRING", 0, 0},
    {"builtin_random_int(1, \"hello\")", true, "second argument to `builtin_random_int` must be INTEGER, got STRING", 0, 0},
    {"builtin_random_int(1)", true, "wrong number of arguments. got=1, want=2", 0, 0},
    {"builtin_random_int()", true, "wrong number of arguments. got=0, want=2", 0, 0},
    {"builtin_random_int(1, 2, 3)", true, "wrong number of arguments. got=3, want=2", 0, 0},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    if tt.hasError {
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != tt.errorMsg {
        t.Errorf("wrong error message. expected=%q, got=%q",
          tt.errorMsg, errObj.Message)
      }
    } else {
      // Check that it's an integer in the specified range
      integer, ok := evaluated.(*Integer)
      if !ok {
        t.Errorf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if integer.Value < tt.minVal || integer.Value > tt.maxVal {
        t.Errorf("random int out of range [%d, %d]. got=%d", tt.minVal, tt.maxVal, integer.Value)
      }
    }
  }
}

func TestBuiltinSumFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_sum([1, 2, 3])", 6},
    {"builtin_sum([1.5, 2.5, 3.0])", 7.0},
    {"builtin_sum([1, 2.5, 3])", 6.5},
    {"builtin_sum([])", 0},
    {"builtin_sum([42])", 42},
    {"builtin_sum([-1, -2, -3])", -6},
    {"builtin_sum([0, 0, 0])", 0},
    {"builtin_sum([1, \"hello\", 3])", "array elements for `builtin_sum` must be INTEGER or FLOAT, got STRING"},
    {"builtin_sum(42)", "argument to `builtin_sum` must be ARRAY, got INTEGER"},
    {"builtin_sum()", "wrong number of arguments. got=0, want=1"},
    {"builtin_sum([1, 2], [3, 4])", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinAverageFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"builtin_average([1, 2, 3])", 2.0},
    {"builtin_average([1.5, 2.5, 3.0])", 2.3333333333333335},
    {"builtin_average([1, 2.5, 3])", 2.1666666666666665},
    {"builtin_average([42])", 42.0},
    {"builtin_average([-1, -2, -3])", -2.0},
    {"builtin_average([0, 0, 0])", 0.0},
    {"builtin_average([10, 20])", 15.0},
    {"builtin_average([])", "cannot calculate average of empty array"},
    {"builtin_average([1, \"hello\", 3])", "array elements for `builtin_average` must be INTEGER or FLOAT, got STRING"},
    {"builtin_average(42)", "argument to `builtin_average` must be ARRAY, got INTEGER"},
    {"builtin_average()", "wrong number of arguments. got=0, want=1"},
    {"builtin_average([1, 2], [3, 4])", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)

    switch expected := tt.expected.(type) {
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      errObj, ok := evaluated.(*Error)
      if !ok {
        t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
        continue
      }
      if errObj.Message != expected {
        t.Errorf("wrong error message. expected=%q, got=%q",
          expected, errObj.Message)
      }
    }
  }
}

func TestBuiltinIsNumberFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"builtin_is_number(42)", true},
    {"builtin_is_number(3.14)", true},
    {"builtin_is_number(0)", true},
    {"builtin_is_number(0.0)", true},
    {"builtin_is_number(-5)", true},
    {"builtin_is_number(-3.14)", true},
    {"builtin_is_number(\"hello\")", false},
    {"builtin_is_number(true)", false},
    {"builtin_is_number(false)", false},
    {"builtin_is_number([])", false},
    {"builtin_is_number([1, 2, 3])", false},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestBuiltinIsIntegerFunction(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"builtin_is_integer(42)", true},
    {"builtin_is_integer(0)", true},
    {"builtin_is_integer(-5)", true},
    {"builtin_is_integer(3.14)", false},
    {"builtin_is_integer(0.0)", false},
    {"builtin_is_integer(-3.14)", false},
    {"builtin_is_integer(\"hello\")", false},
    {"builtin_is_integer(true)", false},
    {"builtin_is_integer(false)", false},
    {"builtin_is_integer([])", false},
    {"builtin_is_integer([1, 2, 3])", false},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestBuiltinIsNumberErrors(t *testing.T) {
  tests := []struct {
    input    string
    errorMsg string
  }{
    {"builtin_is_number()", "wrong number of arguments. got=0, want=1"},
    {"builtin_is_number(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)
    
    errObj, ok := evaluated.(*Error)
    if !ok {
      t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
      continue
    }
    if errObj.Message != tt.errorMsg {
      t.Errorf("wrong error message. expected=%q, got=%q",
        tt.errorMsg, errObj.Message)
    }
  }
}

func TestBuiltinIsIntegerErrors(t *testing.T) {
  tests := []struct {
    input    string
    errorMsg string
  }{
    {"builtin_is_integer()", "wrong number of arguments. got=0, want=1"},
    {"builtin_is_integer(1, 2)", "wrong number of arguments. got=2, want=1"},
  }

  for _, tt := range tests {
    evaluated := testEvalBuiltin(tt.input)
    
    errObj, ok := evaluated.(*Error)
    if !ok {
      t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
      continue
    }
    if errObj.Message != tt.errorMsg {
      t.Errorf("wrong error message. expected=%q, got=%q",
        tt.errorMsg, errObj.Message)
    }
  }
}