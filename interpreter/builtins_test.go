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
    {"pop([])", nil},
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

// Helper function specifically for builtin tests
func testEvalBuiltin(input string) Value {
  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  return Eval(program, env)
}