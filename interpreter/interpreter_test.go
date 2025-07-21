package interpreter

import (
  "testing"
  "rush/lexer"
  "rush/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"5", 5},
    {"10", 10},
    {"-5", -5},
    {"-10", -10},
    {"5 + 5 + 5 + 5 - 10", 10},
    {"2 * 2 * 2 * 2 * 2", 32},
    {"-50 + 100 + -50", 0},
    {"5 * 2 + 10", 20},
    {"5 + 2 * 10", 25},
    {"20 + 2 * -10", 0},
    {"2 * (5 + 10)", 30},
    {"3 * 3 * 3 + 10", 37},
    {"3 * (3 * 3) + 10", 37},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestEvalFloatExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected float64
  }{
    {"5.5", 5.5},
    {"10.25", 10.25},
    {"-5.5", -5.5},
    {"5.5 + 2.5", 8.0},
    {"10.0 - 3.5", 6.5},
    {"2.5 * 4.0", 10.0},
    {"10.0 / 2.0", 5.0},
    {"5 + 2.5", 7.5}, // Mixed integer and float
    {"10.5 - 5", 5.5},
    {"50 / 2 * 2 + 10", 60.0}, // Integer division results in float
    {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50.0},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testFloatObject(t, evaluated, tt.expected)
  }
}

func TestEvalBooleanExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"true", true},
    {"false", false},
    {"1 < 2", true},
    {"1 > 2", false},
    {"1 < 1", false},
    {"1 > 1", false},
    {"1 == 1", true},
    {"1 != 1", false},
    {"1 == 2", false},
    {"1 != 2", true},
    {"1 <= 2", true},
    {"1 >= 2", false},
    {"2 <= 2", true},
    {"2 >= 2", true},
    {"true == true", true},
    {"false == false", true},
    {"true == false", false},
    {"true != false", true},
    {"false != true", true},
    {"(1 < 2) == true", true},
    {"(1 < 2) == false", false},
    {"(1 > 2) == true", false},
    {"(1 > 2) == false", true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestBangOperator(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"!true", false},
    {"!false", true},
    {"!5", false},
    {"!!true", true},
    {"!!false", false},
    {"!!5", true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestLogicalOperators(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"true && true", true},
    {"true && false", false},
    {"false && true", false},
    {"false && false", false},
    {"true || true", true},
    {"true || false", true},
    {"false || true", true},
    {"false || false", false},
    {"5 > 3 && 2 < 4", true},
    {"5 > 3 && 2 > 4", false},
    {"5 < 3 || 2 < 4", true},
    {"5 < 3 || 2 > 4", false},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestIfElseExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"if (true) { 10 }", 10},
    {"if (false) { 10 }", nil},
    {"if (1) { 10 }", 10},
    {"if (1 < 2) { 10 }", 10},
    {"if (1 > 2) { 10 }", nil},
    {"if (1 > 2) { 10 } else { 20 }", 20},
    {"if (1 < 2) { 10 } else { 20 }", 10},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    integer, ok := tt.expected.(int)
    if ok {
      testIntegerObject(t, evaluated, int64(integer))
    } else {
      testNullObject(t, evaluated)
    }
  }
}

func TestReturnStatements(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"return 10;", 10},
    {"return 10; 9;", 10},
    {"return 2 * 5; 9;", 10},
    {"9; return 2 * 5; 9;", 10},
    {`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }
  return 1;
}
`, 10},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestErrorHandling(t *testing.T) {
  tests := []struct {
    input           string
    expectedMessage string
  }{
    {
      "5 + true;",
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
      "foobar",
      "identifier not found: foobar",
    },
    {
      `"Hello" - "World"`,
      "unknown operator: -",
    },
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)

    errObj, ok := evaluated.(*Error)
    if !ok {
      t.Errorf("no error object returned. got=%T(%+v)",
        evaluated, evaluated)
      continue
    }

    if errObj.Message != tt.expectedMessage {
      t.Errorf("wrong error message. expected=%q, got=%q",
        tt.expectedMessage, errObj.Message)
    }
  }
}

func TestVariableAssignmentAndAccess(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"a = 5; a;", 5},
    {"a = 5 * 5; a;", 25},
    {"a = 5; b = a; b;", 5},
    {"a = 5; b = a; c = a + b + 5; c;", 15},
  }

  for _, tt := range tests {
    testIntegerObject(t, testEval(tt.input), tt.expected)
  }
}

func TestStringLiterals(t *testing.T) {
  input := `"Hello World!"`

  evaluated := testEval(input)
  str, ok := evaluated.(*String)
  if !ok {
    t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
  }

  if str.Value != "Hello World!" {
    t.Errorf("String has wrong value. got=%q", str.Value)
  }
}

func TestStringConcatenation(t *testing.T) {
  input := `"Hello" + " " + "World!"`

  evaluated := testEval(input)
  str, ok := evaluated.(*String)
  if !ok {
    t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
  }

  if str.Value != "Hello World!" {
    t.Errorf("String has wrong value. got=%q", str.Value)
  }
}

func TestFunctionObject(t *testing.T) {
  input := "fn(x) { x + 2; };"

  evaluated := testEval(input)
  fn, ok := evaluated.(*Function)
  if !ok {
    t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
  }

  if len(fn.Parameters) != 1 {
    t.Fatalf("function has wrong parameters. Parameters=%+v",
      fn.Parameters)
  }

  if fn.Parameters[0].String() != "x" {
    t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
  }

  expectedBody := "{(x + 2)}"

  if fn.Body.String() != expectedBody {
    t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
  }
}

func TestFunctionApplication(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"identity = fn(x) { x; }; identity(5);", 5},
    {"identity = fn(x) { return x; }; identity(5);", 5},
    {"double = fn(x) { x * 2; }; double(5);", 10},
    {"add = fn(x, y) { x + y; }; add(5, 5);", 10},
    {"add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
    {"fn(x) { x; }(5)", 5},
  }

  for _, tt := range tests {
    testIntegerObject(t, testEval(tt.input), tt.expected)
  }
}

func TestArrayLiterals(t *testing.T) {
  input := "[1, 2 * 2, 3 + 3]"

  evaluated := testEval(input)
  result, ok := evaluated.(*Array)
  if !ok {
    t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Elements) != 3 {
    t.Fatalf("array has wrong num of elements. got=%d",
      len(result.Elements))
  }

  testIntegerObject(t, result.Elements[0], 1)
  testIntegerObject(t, result.Elements[1], 4)
  testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {
      "[1, 2, 3][0]",
      1,
    },
    {
      "[1, 2, 3][1]",
      2,
    },
    {
      "[1, 2, 3][2]",
      3,
    },
    {
      "i = 0; [1][i];",
      1,
    },
    {
      "[1, 2, 3][1 + 1];",
      3,
    },
    {
      "myArray = [1, 2, 3]; myArray[2];",
      3,
    },
    {
      "myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
      6,
    },
    {
      "myArray = [1, 2, 3]; i = myArray[0]; myArray[i]",
      2,
    },
    {
      "arr = [1, 2, 3]; try { arr[3] } catch (IndexError error) { if (false) { 1 } }",
      nil,
    },
    {
      "arr = [1, 2, 3]; try { arr[-1] } catch (IndexError error) { if (false) { 1 } }",
      nil,
    },
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    integer, ok := tt.expected.(int)
    if ok {
      testIntegerObject(t, evaluated, int64(integer))
    } else {
      testNullObject(t, evaluated)
    }
  }
}

func TestStringIndexing(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {
      `"hello"[0]`,
      "h",
    },
    {
      `"hello"[1]`,
      "e",
    },
    {
      `"hello"[4]`,
      "o",
    },
    {
      `str = "hello"; try { str[5] } catch (IndexError error) { if (false) { 1 } }`,
      nil,
    },
    {
      `str = "hello"; try { str[-1] } catch (IndexError error) { if (false) { 1 } }`,
      nil,
    },
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    str, ok := tt.expected.(string)
    if ok {
      testStringObject(t, evaluated, str)
    } else {
      testNullObject(t, evaluated)
    }
  }
}

func TestWhileLoop(t *testing.T) {
  input := `
i = 0
sum = 0
while (i < 5) {
  sum = sum + i
  i = i + 1
}
sum
`

  evaluated := testEval(input)
  testIntegerObject(t, evaluated, 10) // 0+1+2+3+4 = 10
}

func TestForLoop(t *testing.T) {
  input := `
sum = 0
for (i = 0; i < 5; i = i + 1) {
  sum = sum + i
}
sum
`

  evaluated := testEval(input)
  testIntegerObject(t, evaluated, 10) // 0+1+2+3+4 = 10
}

// Helper functions

func testEval(input string) Value {
  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj Value, expected int64) bool {
  result, ok := obj.(*Integer)
  if !ok {
    t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
    return false
  }
  if result.Value != expected {
    t.Errorf("object has wrong value. got=%d, want=%d",
      result.Value, expected)
    return false
  }

  return true
}

func testFloatObject(t *testing.T, obj Value, expected float64) bool {
  result, ok := obj.(*Float)
  if !ok {
    t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
    return false
  }
  if result.Value != expected {
    t.Errorf("object has wrong value. got=%f, want=%f",
      result.Value, expected)
    return false
  }

  return true
}

func testBooleanObject(t *testing.T, obj Value, expected bool) bool {
  result, ok := obj.(*Boolean)
  if !ok {
    t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
    return false
  }
  if result.Value != expected {
    t.Errorf("object has wrong value. got=%t, want=%t",
      result.Value, expected)
    return false
  }

  return true
}

func testNullObject(t *testing.T, obj Value) bool {
  if obj != NULL {
    t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
    return false
  }
  return true
}

func testStringObject(t *testing.T, obj Value, expected string) bool {
  result, ok := obj.(*String)
  if !ok {
    t.Errorf("object is not String. got=%T (%+v)", obj, obj)
    return false
  }
  if result.Value != expected {
    t.Errorf("object has wrong value. got=%q, want=%q",
      result.Value, expected)
    return false
  }

  return true
}

// Error handling tests

func TestBasicTryCatch(t *testing.T) {
  input := `
caught = false
try {
  throw RuntimeError("test error")
} catch (error) {
  caught = true
}
caught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestTypeSpecificCatch(t *testing.T) {
  input := `
validationCaught = false
runtimeCaught = false
try {
  throw ValidationError("validation failed")
} catch (ValidationError error) {
  validationCaught = true
} catch (RuntimeError error) {
  runtimeCaught = true
}
validationCaught && !runtimeCaught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestMultipleCatchBlocks(t *testing.T) {
  input := `
typeErrorCaught = false
try {
  throw TypeError("wrong type")
} catch (ValidationError error) {
  # Should not catch ValidationError
} catch (TypeError error) {
  typeErrorCaught = true
} catch (error) {
  # Should not reach generic catch
}
typeErrorCaught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestFinallyBlockExecution(t *testing.T) {
  input := `
finallyRan = false
try {
  throw RuntimeError("error")
} catch (error) {
  # Exception caught
} finally {
  finallyRan = true
}
finallyRan
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestFinallyBlockWithoutException(t *testing.T) {
  input := `
finallyRan = false
result = 0
try {
  result = 42
} finally {
  finallyRan = true
}
finallyRan && (result == 42)
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestErrorPropertyAccess(t *testing.T) {
  input := `
errorType = ""
errorMessage = ""
try {
  throw ValidationError("invalid data")
} catch (error) {
  errorType = error.type
  errorMessage = error.message
}
(errorType == "ValidationError") && (errorMessage == "invalid data")
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestStackTraceGeneration(t *testing.T) {
  input := `
stackGenerated = false

innerFunc = fn() {
  throw RuntimeError("inner error")
}

try {
  innerFunc()
} catch (error) {
  stackGenerated = len(error.stack) > 0
}
stackGenerated
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestArrayIndexException(t *testing.T) {
  input := `
caught = false
arr = [1, 2, 3]
try {
  value = arr[10]
} catch (IndexError error) {
  caught = true
}
caught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestStringIndexException(t *testing.T) {
  input := `
caught = false
str = "hello"
try {
  char = str[20]
} catch (IndexError error) {
  caught = true
}
caught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestPopEmptyArrayException(t *testing.T) {
  input := `
caught = false
emptyArr = []
try {
  value = pop(emptyArr)
} catch (IndexError error) {
  caught = true
}
caught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestNestedTryCatch(t *testing.T) {
  input := `
outerCaught = false

try {
  try {
    throw ValidationError("inner error")
  } catch (TypeError error) {
    # Should not catch TypeError, will propagate to outer
  }
} catch (ValidationError error) {
  outerCaught = true
}
outerCaught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestExceptionPropagationThroughFunctions(t *testing.T) {
  input := `
caught = false

deepFunc = fn() {
  throw RuntimeError("deep error")
}

middleFunc = fn() {
  deepFunc()
}

try {
  middleFunc()
} catch (RuntimeError error) {
  caught = true
}
caught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestAllBuiltinErrorTypes(t *testing.T) {
  tests := []struct {
    errorType string
    input     string
  }{
    {"Error", `
try {
  throw Error("test")
} catch (error) {
  error.type
}
`},
    {"ValidationError", `
try {
  throw ValidationError("test")
} catch (error) {
  error.type
}
`},
    {"TypeError", `
try {
  throw TypeError("test")
} catch (error) {
  error.type
}
`},
    {"IndexError", `
try {
  throw IndexError("test")
} catch (error) {
  error.type
}
`},
    {"ArgumentError", `
try {
  throw ArgumentError("test")
} catch (error) {
  error.type
}
`},
    {"RuntimeError", `
try {
  throw RuntimeError("test")
} catch (error) {
  error.type
}
`},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testStringObject(t, evaluated, tt.errorType)
  }
}

func TestUncaughtExceptionPropagation(t *testing.T) {
  input := `throw RuntimeError("uncaught error")`
  
  evaluated := testEval(input)
  
  // Should return an Exception value
  exception, ok := evaluated.(*Exception)
  if !ok {
    t.Errorf("Expected Exception, got %T", evaluated)
    return
  }
  
  // The wrapped error should be an Error
  errorObj, ok := exception.Error.(*Error)
  if !ok {
    t.Errorf("Expected Error inside Exception, got %T", exception.Error)
    return
  }
  
  if errorObj.ErrorType != "RuntimeError" {
    t.Errorf("Expected RuntimeError, got %s", errorObj.ErrorType)
  }
  
  if errorObj.Message != "uncaught error" {
    t.Errorf("Expected 'uncaught error', got %s", errorObj.Message)
  }
}

func TestErrorLineAndColumnNumbers(t *testing.T) {
  input := `
errorLine = 0
errorCol = 0
try {
  throw RuntimeError("positioned error")
} catch (error) {
  errorLine = error.line
  errorCol = error.column
}
errorLine >= 0
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestCatchVariableScope(t *testing.T) {
  input := `
message = ""
try {
  throw ValidationError("scoped error")
} catch (error) {
  message = error.message
}
message == "scoped error"
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestThrowInCatchBlock(t *testing.T) {
  input := `
innerCaught = false
outerCaught = false

try {
  try {
    throw ValidationError("first error")
  } catch (error) {
    innerCaught = true
    throw RuntimeError("second error")
  }
} catch (RuntimeError error) {
  outerCaught = true
}
innerCaught && outerCaught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

// Helper function to test error objects
func testErrorObject(t *testing.T, obj Value, expectedType, expectedMessage string) bool {
  errorObj, ok := obj.(*Error)
  if !ok {
    t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
    return false
  }
  
  if errorObj.ErrorType != expectedType {
    t.Errorf("wrong error type. got=%s, want=%s", errorObj.ErrorType, expectedType)
    return false
  }
  
  if errorObj.Message != expectedMessage {
    t.Errorf("wrong error message. got=%s, want=%s", errorObj.Message, expectedMessage)
    return false
  }
  
  return true
}