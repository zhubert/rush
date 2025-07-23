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
  
  // Use tolerance-based comparison for floating-point numbers
  tolerance := 1e-9
  if abs(result.Value - expected) > tolerance {
    t.Errorf("object has wrong value. got=%f, want=%f",
      result.Value, expected)
    return false
  }

  return true
}

func abs(x float64) float64 {
  if x < 0 {
    return -x
  }
  return x
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

// Core language feature tests

func TestBasicVariableAssignments(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Integer assignments
    {"a = 42; a", 42},
    {"b = -15; b", -15},
    
    // Float assignments
    {"pi = 3.14; pi", 3.14},
    {"negative = -2.5; negative", -2.5},
    
    // String assignments
    {"greeting = \"hello world\"; greeting", "hello world"},
    {"message = \"Rush programming language\"; message", "Rush programming language"},
    
    // Boolean assignments
    {"isTrue = true; isTrue", true},
    {"isFalse = false; isFalse", false},
    
    // Variable chaining
    {"x = 5; y = x; y", 5},
    {"a = 10; b = a + 5; b", 15},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    }
  }
}

func TestArithmeticExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Basic arithmetic
    {"a = 5; b = a + 3; b", 8},
    {"result = 10 * 2; result", 20},
    {"difference = 100 - 42; difference", 58},
    {"quotient = 20 / 4; quotient", 5.0}, // Division results in float
    
    // Mixed arithmetic  
    {"sum = 42 + 10; sum", 52},
    {"product = 3.14 * 2.5; product", 7.85},
    
    // Complex expressions
    {"result = (10 + 5) * 2 - 8; result", 22},
    {"complex = 2 * (5 + 10) / 3; complex", 10.0},
    
    // Mixed integer and float
    {"pi = 3.14; area = pi * 10; area", 31.4},
    {"mixed = 5 + 2.5; mixed", 7.5},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    }
  }
}

func TestStringOperations(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    // String concatenation
    {"greeting = \"Hello, \"; name = \"Rush\"; message = greeting + name; message", "Hello, Rush"},
    {"part1 = \"Rush \"; part2 = \"programming \"; part3 = \"language\"; full = part1 + part2 + part3; full", "Rush programming language"},
    
    // String with expressions
    {"base = \"Result: \"; num = 42; combined = base + \"Value is \" + \"42\"; combined", "Result: Value is 42"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testStringObject(t, evaluated, tt.expected)
  }
}

func TestBooleanOperationsBasic(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    // Basic boolean values
    {"a = true; a", true},
    {"b = false; b", false},
    
    // Boolean expressions
    {"isGreater = 42 > 30; isGreater", true},
    {"isEqual = 3.14 == 3.14; isEqual", true},
    {"isLess = 5 < 10; isLess", true},
    
    // Negation
    {"isTrue = true; isFalse = !isTrue; isFalse", false},
    {"original = false; negated = !original; negated", true},
    
    // Comparisons
    {"comparison = 5 > 3; comparison", true},
    {"equality = 10 == 10; equality", true},
    {"inequality = 7 != 8; inequality", true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestArrayLiteralsAndBasicOperations(t *testing.T) {
  tests := []struct {
    input       string
    expectedLen int
    elements    []interface{}
  }{
    // Basic array literals
    {"numbers = [1, 2, 3, 4, 5]; numbers", 5, []interface{}{1, 2, 3, 4, 5}},
    {"mixed = [42, 3.14, \"hello\", true]; mixed", 4, []interface{}{42, 3.14, "hello", true}},
    
    // Arrays with variables
    {"a = 10; b = 20; arr = [a, b, a + b]; arr", 3, []interface{}{10, 20, 30}},
    
    // Empty array
    {"empty = []; empty", 0, []interface{}{}},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    array, ok := evaluated.(*Array)
    if !ok {
      t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
      continue
    }
    
    if len(array.Elements) != tt.expectedLen {
      t.Errorf("array has wrong length. expected=%d, got=%d", tt.expectedLen, len(array.Elements))
      continue
    }
    
    for i, expectedElement := range tt.elements {
      switch expected := expectedElement.(type) {
      case int:
        testIntegerObject(t, array.Elements[i], int64(expected))
      case float64:
        testFloatObject(t, array.Elements[i], expected)
      case string:
        testStringObject(t, array.Elements[i], expected)
      case bool:
        testBooleanObject(t, array.Elements[i], expected)
      }
    }
  }
}

func TestMixedTypeOperations(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Mixed number operations
    {"intVal = 5; floatVal = 2.5; result = intVal + floatVal; result", 7.5},
    {"mixed = 10 * 1.5; mixed", 15.0},
    
    // Complex mixed operations
    {"a = 5; pi = 3.14; area = pi * a; area", 15.7},
    {"base = 10; height = 2.5; area = base * height; area", 25.0},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case float64:
      testFloatObject(t, evaluated, expected)
    }
  }
}

func TestVariableScoping(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Basic scoping
    {"outer = 10; if (true) { inner = 20; outer = outer + inner }; outer", 30},
    
    // Variable shadowing in blocks
    {"x = 5; if (true) { x = 10 }; x", 10},
    
    // Multiple assignments
    {"a = 1; b = 2; c = 3; result = a + b + c; result", 6},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

// Boolean logic and conditional tests

func TestBooleanLogicOperators(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    // Basic AND operations
    {"a = true; b = false; andResult = a && b; andResult", false},
    {"a = true; b = true; andResult = a && b; andResult", true},
    {"a = false; b = false; andResult = a && b; andResult", false},
    {"a = false; b = true; andResult = a && b; andResult", false},
    
    // Basic OR operations  
    {"a = true; b = false; orResult = a || b; orResult", true},
    {"a = true; b = true; orResult = a || b; orResult", true},
    {"a = false; b = false; orResult = a || b; orResult", false},
    {"a = false; b = true; orResult = a || b; orResult", true},
    
    // Negation operations
    {"a = true; notA = !a; notA", false},
    {"b = false; notB = !b; notB", true},
    
    // Complex boolean expressions
    {"complexAnd = (5 > 3) && (10 < 20); complexAnd", true},
    {"complexOr = (5 < 3) || (10 > 20); complexOr", false},
    {"complexMixed = (5 > 3) || (10 > 20); complexMixed", true},
    
    // Nested boolean expressions  
    {"result = (true && false) || (true && true); result", true},
    {"result = (false || true) && (true || false); result", true},
    {"result = !(false && true) && (true || false); result", true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestSimpleConditionals(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Basic if statements
    {"x = 15; if (x > 10) { result = \"greater\" } else { result = \"smaller\" }; result", "greater"},
    {"x = 5; if (x > 10) { result = \"greater\" } else { result = \"smaller\" }; result", "smaller"},
    
    // Boolean operators in conditionals
    {"a = true; b = false; if (a && b) { result = \"both true\" } else { result = \"not both\" }; result", "not both"},
    {"a = true; b = false; if (a || b) { result = \"at least one\" } else { result = \"none\" }; result", "at least one"},
    
    // Numeric comparisons
    {"age = 25; if (age >= 18) { status = \"adult\" } else { status = \"minor\" }; status", "adult"},
    {"age = 16; if (age >= 18) { status = \"adult\" } else { status = \"minor\" }; status", "minor"},
    
    // Complex boolean conditions
    {"age = 25; hasLicense = true; if ((age >= 16) && hasLicense) { canDrive = true } else { canDrive = false }; canDrive", true},
    {"age = 15; hasLicense = true; if ((age >= 16) && hasLicense) { canDrive = true } else { canDrive = false }; canDrive", false},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    }
  }
}

func TestComplexBooleanExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    // Short-circuit evaluation
    {"false && (5 / 0)", false}, // Should not evaluate 5/0
    {"true || (5 / 0)", true},   // Should not evaluate 5/0
    
    // Multiple operators
    {"(5 > 3) && (10 < 20) && (7 == 7)", true},
    {"(5 > 3) || (10 > 20) || (7 != 7)", true},
    {"(5 < 3) && (10 < 20) && (7 == 7)", false},
    
    // Combination with arithmetic
    {"x = 10; y = 5; (x + y > 12) && (x - y < 8)", true},
    {"a = 8; b = 3; (a * b > 20) || (a / b < 3)", true},
    
    // Nested parentheses
    {"((5 > 3) && (10 < 20)) || ((2 > 5) && (1 < 0))", true},
    {"((5 < 3) || (10 > 20)) && ((2 < 5) || (1 > 0))", false},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    if !testBooleanObject(t, evaluated, tt.expected) {
      t.Errorf("Test failed for input: %s", tt.input)
    }
  }
}

func TestBooleanComparisonOperators(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    // Equality comparisons
    {"5 == 5", true},
    {"5 == 7", false},
    {"3.14 == 3.14", true},
    {"true == true", true},
    {"false == false", true},
    {"true == false", false},
    
    // Inequality comparisons
    {"5 != 7", true},
    {"5 != 5", false},
    {"true != false", true},
    {"true != true", false},
    
    // Relational comparisons
    {"10 > 5", true},
    {"5 > 10", false},
    {"10 < 15", true},
    {"15 < 10", false},
    {"10 >= 10", true},
    {"10 >= 5", true},
    {"5 >= 10", false},
    {"10 <= 10", true},
    {"5 <= 10", true},
    {"10 <= 5", false},
    
    // Mixed type comparisons
    {"5.0 == 5", true},
    {"3.5 > 3", true},
    {"2.0 < 3", true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

// Control flow tests

func TestAdvancedControlFlow(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Basic if/else with variable assignment
    {`
x = 10
if (x > 5) {
  result = "x is greater than 5"
} else {
  result = "x is not greater than 5"
}
result`, "x is greater than 5"},
    
    // Nested if statements
    {`
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
grade`, "B"},
    
    // Complex nested conditions with different scores
    {`
score = 95
if (score >= 90) {
  grade = "A"
} else {
  if (score >= 80) {
    grade = "B"
  } else {
    if (score >= 70) {
      grade = "C"
    } else {
      grade = "F"
    }
  }
}
grade`, "A"},
    
    {`
score = 75
if (score >= 90) {
  grade = "A"
} else {
  if (score >= 80) {
    grade = "B"
  } else {
    if (score >= 70) {
      grade = "C"
    } else {
      grade = "F"
    }
  }
}
grade`, "C"},
    
    {`
score = 65
if (score >= 90) {
  grade = "A"
} else {
  if (score >= 80) {
    grade = "B"
  } else {
    if (score >= 70) {
      grade = "C"
    } else {
      grade = "F"
    }
  }
}
grade`, "F"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

func TestComplexConditionalLogic(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Multiple condition combinations
    {`
age = 25
hasLicense = true
canDrive = (age >= 16) && hasLicense
canDrive`, true},
    
    {`
age = 15
hasLicense = true
canDrive = (age >= 16) && hasLicense
canDrive`, false},
    
    {`
age = 25
hasLicense = false
canDrive = (age >= 16) && hasLicense
canDrive`, false},
    
    // Complex decision making
    {`
temperature = 75
humidity = 60
if ((temperature > 70) && (humidity < 80)) {
  comfort = "comfortable"
} else {
  if (temperature > 85) {
    comfort = "too hot"
  } else {
    if (humidity > 90) {
      comfort = "too humid"
    } else {
      comfort = "not ideal"
    }
  }
}
comfort`, "comfortable"},
    
    // Multiple boolean variables
    {`
isWeekend = true
isRaining = false
hasWork = false
if (isWeekend && !isRaining && !hasWork) {
  activity = "go hiking"
} else {
  if (isRaining) {
    activity = "stay inside"
  } else {
    activity = "normal day"
  }
}
activity`, "go hiking"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

func TestControlFlowWithArithmetic(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Arithmetic in conditions
    {`
a = 10
b = 5
if ((a + b) > 12) {
  result = "sum is large"
} else {
  result = "sum is small"
}
result`, "sum is large"},
    
    // Multiple arithmetic comparisons
    {`
x = 8
y = 3
if ((x * y) > 20) {
  category = "high"
} else {
  if ((x + y) > 10) {
    category = "medium"
  } else {
    category = "low"
  }
}
category`, "high"},
    
    // Complex arithmetic in nested conditions
    {`
price = 50
discount = 10
tax = 5
finalPrice = price - discount + tax
if (finalPrice < 40) {
  affordability = "cheap"
} else {
  if (finalPrice < 60) {
    affordability = "reasonable"
  } else {
    affordability = "expensive"
  }
}
affordability`, "reasonable"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case string:
      testStringObject(t, evaluated, expected)
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

func TestControlFlowScoping(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Variable modifications in conditional blocks
    {`
counter = 0
if (true) {
  counter = counter + 10
  inner = 5
}
counter`, 10},
    
    // Nested scoping
    {`
outer = 1
if (true) {
  outer = outer + 1
  if (true) {
    outer = outer + 1
    if (true) {
      outer = outer + 1
    }
  }
}
outer`, 4},
    
    // Multiple variable assignments in conditions
    {`
status = "unknown"
level = 0
value = 15
if (value > 10) {
  status = "high"
  level = 3
} else {
  status = "low"
  level = 1
}
status + ":" + type(level)`, "high:INTEGER"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case string:
      testStringObject(t, evaluated, expected)
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

// Loop tests

func TestForLoops(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Basic for loop - sum
    {`
sum = 0
for (i = 0; i < 5; i = i + 1) {
  sum = sum + i
}
sum`, 10}, // 0+1+2+3+4 = 10
    
    // For loop with array access
    {`
numbers = [2, 4, 6, 8, 10]
product = 1
for (j = 0; j < 5; j = j + 1) {
  product = product * numbers[j]
}
product`, 3840}, // 2*4*6*8*10 = 3840
    
    // For loop that modifies outer scope variable
    {`
counter = 100
for (k = 0; k < 3; k = k + 1) {
  counter = counter + 10
}
counter`, 130}, // 100 + 3*10 = 130
    
    // Combined result matching original test
    {`
sum = 0
for (i = 0; i < 5; i = i + 1) {
  sum = sum + i
}
numbers = [2, 4, 6, 8, 10]
product = 1
for (j = 0; j < 5; j = j + 1) {
  product = product * numbers[j]
}
counter = 100
for (k = 0; k < 3; k = k + 1) {
  counter = counter + 10
}
result = sum + product + counter
result`, 3980}, // 10 + 3840 + 130 = 3980
    
    // For loop with different step
    {`
total = 0
for (x = 0; x < 10; x = x + 2) {
  total = total + x
}
total`, 20}, // 0+2+4+6+8 = 20
    
    // Nested for loops
    {`
total = 0
for (i = 1; i <= 3; i = i + 1) {
  for (j = 1; j <= 2; j = j + 1) {
    total = total + (i * j)
  }
}
total`, 18}, // (1*1 + 1*2) + (2*1 + 2*2) + (3*1 + 3*2) = 3 + 6 + 9 = 18
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestWhileLoops(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Simple counter while loop
    {`
counter = 0
sum = 0
while (counter < 5) {
  sum = sum + counter
  counter = counter + 1
}
sum`, 10}, // 0+1+2+3+4 = 10
    
    // While loop with array indexing
    {`
numbers = [1, 2, 3, 4, 5]
index = 0
product = 1
while (index < 5) {
  product = product * numbers[index]
  index = index + 1
}
product`, 120}, // 1*2*3*4*5 = 120
    
    // Combined result matching original test
    {`
counter = 0
sum = 0
while (counter < 5) {
  sum = sum + counter
  counter = counter + 1
}
numbers = [1, 2, 3, 4, 5]
index = 0
product = 1
while (index < 5) {
  product = product * numbers[index]
  index = index + 1
}
result = sum + product
result`, 130}, // 10 + 120 = 130
    
    // While loop with condition change
    {`
value = 1
count = 0
while (value < 100) {
  value = value * 2
  count = count + 1
}
count`, 7}, // 1->2->4->8->16->32->64->128 (7 iterations)
    
    // While loop with break-like condition
    {`
total = 0
num = 10
while (num > 0) {
  total = total + num
  num = num - 2
}
total`, 30}, // 10+8+6+4+2 = 30
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestNestedLoops(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Nested for loops with 2D calculation
    {`
total = 0
for (row = 0; row < 3; row = row + 1) {
  for (col = 0; col < 3; col = col + 1) {
    total = total + (row * col)
  }
}
total`, 9}, // (0*0+0*1+0*2) + (1*0+1*1+1*2) + (2*0+2*1+2*2) = 0 + 3 + 6 = 9
    
    // Nested loops with arrays
    {`
matrix = [[1, 2], [3, 4]]
sum = 0
for (i = 0; i < 2; i = i + 1) {
  for (j = 0; j < 2; j = j + 1) {
    sum = sum + matrix[i][j]
  }
}
sum`, 10}, // 1+2+3+4 = 10
    
    // Mixed for and while loops
    {`
result = 0
for (outer = 1; outer <= 2; outer = outer + 1) {
  inner = 1
  while (inner <= 3) {
    result = result + (outer * inner)
    inner = inner + 1
  }
}
result`, 18}, // (1*1+1*2+1*3) + (2*1+2*2+2*3) = 6 + 12 = 18
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestLoopScoping(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Loop variable scoping
    {`
outerVar = 5
for (i = 0; i < 3; i = i + 1) {
  outerVar = outerVar + i
  loopVar = i * 2
}
outerVar`, 8}, // 5 + 0 + 1 + 2 = 8
    
    // Multiple loop variables
    {`
x = 0
y = 0
for (a = 1; a <= 2; a = a + 1) {
  for (b = 1; b <= 2; b = b + 1) {
    x = x + a
    y = y + b
  }
}
x + y`, 12}, // x: (1+1) + (2+2) = 6, y: (1+2) + (1+2) = 6, total = 12
    
    // Loop modifying multiple outer variables
    {`
sum1 = 0
sum2 = 0
counter = 1
while (counter <= 4) {
  sum1 = sum1 + counter
  sum2 = sum2 + (counter * counter)
  counter = counter + 1
}
sum1 + sum2`, 40}, // sum1: 1+2+3+4=10, sum2: 1+4+9+16=30, total=40
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

// Advanced function tests

func TestComprehensiveFunctionFeatures(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Simple function call
    {`
test1 = fn() { return 42 }
result1 = test1()
result1`, 42},
    
    // Function with parameters
    {`
test2 = fn(x) { return x * 2 }
result2 = test2(21)
result2`, 42}, // 21 * 2 = 42
    
    // Function with multiple parameters
    {`
test3 = fn(a, b, c) { return a + b + c }
result3 = test3(1, 2, 3)
result3`, 6}, // 1 + 2 + 3 = 6
    
    // Function with local scope
    {`
test4 = fn(x) {
  local = x + 10
  return local
}
result4 = test4(5)
result4`, 15}, // 5 + 10 = 15
    
    // Recursive function (factorial)
    {`
factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
result5 = factorial(5)
result5`, 120}, // 5! = 120
    
    // Higher-order function
    {`
apply = fn(f, x) {
  return f(x)
}
square = fn(x) { return x * x }
result6 = apply(square, 7)
result6`, 49}, // 7 * 7 = 49
    
    // Combined result matching original test
    {`
test1 = fn() { return 42 }
result1 = test1()

test2 = fn(x) { return x * 2 }
result2 = test2(21)

test3 = fn(a, b, c) { return a + b + c }
result3 = test3(1, 2, 3)

test4 = fn(x) {
  local = x + 10
  return local
}
result4 = test4(5)

factorial = fn(n) {
  if (n <= 1) {
    return 1
  } else {
    return n * factorial(n - 1)
  }
}
result5 = factorial(5)

apply = fn(f, x) {
  return f(x)
}
square = fn(x) { return x * x }
result6 = apply(square, 7)

final = result1 + result2 + result3 + result4 + result5 + result6
final`, 274}, // 42 + 42 + 6 + 15 + 120 + 49 = 274
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestAdvancedFunctionScoping(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Simple function with parameter
    {`
simpleFunc = fn(x) {
  return x * 2
}
result = simpleFunc(5)
result`, 10}, // 5 * 2 = 10
    
    // Function accessing outer scope
    {`
outerVar = 10
testFunc = fn(x) {
  return x + outerVar
}
result = testFunc(5)
result`, 15}, // 5 + 10 = 15
    
    // Function with arithmetic operations
    {`
calculate = fn(a, b) {
  return (a + b) * 2
}
result = calculate(3, 7)
result`, 20}, // (3 + 7) * 2 = 20
    
    // Function with conditional
    {`
maxValue = fn(a, b) {
  if (a > b) {
    return a
  } else {
    return b
  }
}
result = maxValue(8, 5)
result`, 8}, // max(8, 5) = 8
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestFunctionTypesAndPatterns(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Function returning different types
    {`
getString = fn() { return "hello" }
result = getString()
result`, "hello"},
    
    {`
getBool = fn(x) { return x > 5 }
result = getBool(10)
result`, true},
    
    // Function with array operations
    {`
sumArray = fn(arr) {
  total = 0
  for (i = 0; i < len(arr); i = i + 1) {
    total = total + arr[i]
  }
  return total
}
result = sumArray([1, 2, 3, 4, 5])
result`, 15},
    
    // Function composition
    {`
double = fn(x) { return x * 2 }
addTen = fn(x) { return x + 10 }
compose = fn(f, g, x) {
  return f(g(x))
}
result = compose(double, addTen, 5)
result`, 30}, // double(addTen(5)) = double(15) = 30
    
    // Recursive with different base cases
    {`
fibonacci = fn(n) {
  if (n <= 1) {
    return n
  } else {
    return fibonacci(n - 1) + fibonacci(n - 2)
  }
}
result = fibonacci(7)
result`, 13}, // fib(7) = 13
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    }
  }
}

func TestHigherOrderFunctions(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Map-like function
    {`
mapFunc = fn(arr, func) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    result = push(result, func(arr[i]))
  }
  return result
}
double = fn(x) { return x * 2 }
mapped = mapFunc([1, 2, 3], double)
sum = mapped[0] + mapped[1] + mapped[2]
sum`, 12}, // [2, 4, 6] sum = 12
    
    // Filter-like function
    {`
filterFunc = fn(arr, predicate) {
  result = []
  for (i = 0; i < len(arr); i = i + 1) {
    if (predicate(arr[i])) {
      result = push(result, arr[i])
    }
  }
  return result
}
isGreaterThan3 = fn(x) { return x > 3 }
filtered = filterFunc([1, 2, 3, 4, 5, 6], isGreaterThan3)
sum = 0
for (i = 0; i < len(filtered); i = i + 1) {
  sum = sum + filtered[i]
}
sum`, 15}, // [4, 5, 6] sum = 15
    
    // Reduce-like function
    {`
reduceFunc = fn(arr, func, initial) {
  accumulator = initial
  for (i = 0; i < len(arr); i = i + 1) {
    accumulator = func(accumulator, arr[i])
  }
  return accumulator
}
add = fn(a, b) { return a + b }
result = reduceFunc([1, 2, 3, 4, 5], add, 0)
result`, 15}, // 0+1+2+3+4+5 = 15
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

// Array and string indexing tests (Phase 9 compatible)

func TestArrayIndexingWithExceptionHandling(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Valid array indexing
    {`
numbers = [10, 20, 30, 40, 50]
first = numbers[0]
first`, 10},
    
    {`
numbers = [10, 20, 30, 40, 50]
second = numbers[1]
second`, 20},
    
    {`
numbers = [10, 20, 30, 40, 50]
last = numbers[4]
last`, 50},
    
    // Array indexing with expressions
    {`
numbers = [10, 20, 30, 40, 50]
index = 2
third = numbers[index]
third`, 30},
    
    {`
numbers = [10, 20, 30, 40, 50]
index = 2
fourth = numbers[index + 1]
fourth`, 40},
    
    // Combined valid indexing result matching original test
    {`
numbers = [10, 20, 30, 40, 50]
first = numbers[0]
second = numbers[1]
last = numbers[4]
index = 2
third = numbers[index]
fourth = numbers[index + 1]
result = first + second + last
result`, 80}, // 10 + 20 + 50 = 80
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    }
  }
}

func TestStringIndexingWithExceptionHandling(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    // Valid string indexing
    {`
message = "hello"
char1 = message[0]
char1`, "h"},
    
    {`
message = "hello"
char2 = message[1]
char2`, "e"},
    
    {`
message = "hello"
char5 = message[4]
char5`, "o"},
    
    // String indexing with variables
    {`
message = "world"
index = 1
char = message[index]
char`, "o"},
    
    // Complex string operations
    {`
word1 = "Rush"
word2 = "Lang"
combined = word1[0] + word2[0]
combined`, "RL"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testStringObject(t, evaluated, tt.expected)
  }
}

func TestIndexingExceptionHandling(t *testing.T) {
  tests := []struct {
    input         string
    expectedCatch bool
  }{
    // Array out of bounds - positive index
    {`
caught = false
numbers = [10, 20, 30]
try {
  outOfBounds = numbers[10]
} catch (IndexError error) {
  caught = true
}
caught`, true},
    
    // Array out of bounds - negative index  
    {`
caught = false
numbers = [10, 20, 30]
try {
  outOfBounds = numbers[-1]
} catch (IndexError error) {
  caught = true
}
caught`, true},
    
    // String out of bounds - positive index
    {`
caught = false
message = "hello"
try {
  outOfBounds = message[20]
} catch (IndexError error) {
  caught = true
}
caught`, true},
    
    // String out of bounds - negative index
    {`
caught = false
message = "hello"
try {
  outOfBounds = message[-5]
} catch (IndexError error) {
  caught = true
}
caught`, true},
    
    // Array bounds checking with expressions
    {`
caught = false
arr = [1, 2, 3]
maxIndex = len(arr)
try {
  value = arr[maxIndex]
} catch (IndexError error) {
  caught = true
}
caught`, true},
    
    // Nested array indexing with exception
    {`
caught = false
matrix = [[1, 2], [3, 4]]
try {
  value = matrix[2][0]  # First access should fail
} catch (IndexError error) {
  caught = true
}
caught`, true},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expectedCatch)
  }
}

func TestIndexingWithDifferentTypes(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    // Array of different types
    {`
mixed = [42, "hello", true, 3.14]
first = mixed[0]
first`, 42},
    
    {`
mixed = [42, "hello", true, 3.14]
second = mixed[1]
second`, "hello"},
    
    {`
mixed = [42, "hello", true, 3.14]
third = mixed[2]
third`, true},
    
    // Nested arrays
    {`
nested = [[1, 2], [3, 4], [5, 6]]
element = nested[1][0]
element`, 3},
    
    {`
nested = [["a", "b"], ["c", "d"]]
element = nested[0][1]
element`, "b"},
    
    // Array of strings with string indexing
    {`
words = ["hello", "world", "rush"]
firstWord = words[0]
firstChar = firstWord[0]
firstChar`, "h"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    
    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    }
  }
}

func TestIndexingInLoopsAndFunctions(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    // Indexing in for loops
    {`
numbers = [1, 2, 3, 4, 5]
sum = 0
for (i = 0; i < len(numbers); i = i + 1) {
  sum = sum + numbers[i]
}
sum`, 15}, // 1+2+3+4+5 = 15
    
    // Function with array indexing
    {`
getElement = fn(arr, index) {
  return arr[index]
}
data = [10, 20, 30]
result = getElement(data, 1)
result`, 20},
    
    // Function with safe indexing
    {`
safeGet = fn(arr, index) {
  try {
    return arr[index]
  } catch (IndexError error) {
    return -1
  }
}
data = [10, 20, 30]
valid = safeGet(data, 1)
invalid = safeGet(data, 10)
valid + invalid`, 19}, // 20 + (-1) = 19
    
    // Complex indexing with nested structures
    {`
matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
total = 0
for (row = 0; row < 3; row = row + 1) {
  for (col = 0; col < 3; col = col + 1) {
    total = total + matrix[row][col]
  }
}
total`, 45}, // 1+2+3+4+5+6+7+8+9 = 45
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestSwitchStatements(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {
      `grade = "A"
      result = ""
      switch (grade) {
        case "A":
          result = "excellent"
        case "B":
          result = "good"
        default:
          result = "other"
      }
      result`,
      "excellent",
    },
    {
      `grade = "B"
      result = ""
      switch (grade) {
        case "A":
          result = "excellent"
        case "B", "C":
          result = "good"
        default:
          result = "other"
      }
      result`,
      "good",
    },
    {
      `grade = "C"
      result = ""
      switch (grade) {
        case "A":
          result = "excellent"
        case "B", "C":
          result = "good"
        default:
          result = "other"
      }
      result`,
      "good",
    },
    {
      `grade = "F"
      result = ""
      switch (grade) {
        case "A":
          result = "excellent"
        case "B", "C":
          result = "good"
        default:
          result = "other"
      }
      result`,
      "other",
    },
    {
      `num = 42
      result = ""
      switch (num) {
        case 1:
          result = "one"
        case 2, 3:
          result = "two or three"
        case 42:
          result = "answer"
        default:
          result = "unknown"
      }
      result`,
      "answer",
    },
    {
      `num = 2
      result = ""
      switch (num) {
        case 1:
          result = "one"
        case 2, 3:
          result = "two or three"
        case 42:
          result = "answer"
        default:
          result = "unknown"
      }
      result`,
      "two or three",
    },
    {
      `num = 99
      result = ""
      switch (num) {
        case 1:
          result = "one"
        case 2, 3:
          result = "two or three"
        case 42:
          result = "answer"
        default:
          result = "unknown"
      }
      result`,
      "unknown",
    },
    {
      `flag = true
      result = ""
      switch (flag) {
        case true:
          result = "yes"
        case false:
          result = "no"
      }
      result`,
      "yes",
    },
    {
      `flag = false
      result = ""
      switch (flag) {
        case true:
          result = "yes"
        case false:
          result = "no"
      }
      result`,
      "no",
    },
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)

    switch expected := tt.expected.(type) {
    case int:
      testIntegerObject(t, evaluated, int64(expected))
    case string:
      testStringObject(t, evaluated, expected)
    case bool:
      testBooleanObject(t, evaluated, expected)
    }
  }
}

func TestSwitchWithoutDefault(t *testing.T) {
  input := `
  x = 5
  result = "none"
  switch (x) {
    case 1:
      result = "one"
    case 2:
      result = "two"
  }
  result`

  evaluated := testEval(input)
  testStringObject(t, evaluated, "none")
}

func TestSwitchAutomaticBreak(t *testing.T) {
  // Test that switch has automatic break behavior (Go-style)
  input := `
  count = 0
  x = 1
  switch (x) {
    case 1:
      count = count + 1
    case 2:
      count = count + 10
  }
  count`

  evaluated := testEval(input)
  testIntegerObject(t, evaluated, 1) // Should be 1, not 11 (no fall-through)
}

func TestSwitchWithReturnStatement(t *testing.T) {
  input := `
  test = fn(x) {
    switch (x) {
      case 1:
        return "one"
      case 2:
        return "two"
      default:
        return "other"
    }
    return "unreachable"
  }
  test(1)`

  evaluated := testEval(input)
  testStringObject(t, evaluated, "one")
}

func TestSwitchInLoop(t *testing.T) {
  input := `
  result = []
  for (i = 1; i <= 3; i = i + 1) {
    value = ""
    switch (i) {
      case 1:
        value = "first"
      case 2:
        value = "second"
      case 3:
        value = "third"
    }
    result = push(result, value)
  }
  len(result)`

  evaluated := testEval(input)
  testIntegerObject(t, evaluated, 3)
}

func TestSwitchWithFloats(t *testing.T) {
  input := `
  x = 3.14
  result = ""
  switch (x) {
    case 1.0:
      result = "one"
    case 3.14:
      result = "pi"
    default:
      result = "other"
  }
  result`

  evaluated := testEval(input)
  testStringObject(t, evaluated, "pi")
}

func TestHashLiterals(t *testing.T) {
  input := `hash = {"name": "Alice", "age": 30, 42: "answer", true: "yes"}
  hash`

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
  }

  expected := map[HashKey]int64{
    HashKey{Type: STRING_VALUE, Value: "name"}:  0,
    HashKey{Type: STRING_VALUE, Value: "age"}:   0,
    HashKey{Type: INTEGER_VALUE, Value: int64(42)}: 0,
    HashKey{Type: BOOLEAN_VALUE, Value: true}:   0,
  }

  if len(result.Pairs) != len(expected) {
    t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
  }

  for expectedKey := range expected {
    _, ok := result.Pairs[expectedKey]
    if !ok {
      t.Errorf("no pair for given key in Pairs")
    }
  }
}

func TestEmptyHashLiteral(t *testing.T) {
  input := "{}"

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Pairs) != 0 {
    t.Fatalf("Hash.Pairs has wrong length. got=%d", len(result.Pairs))
  }
}

func TestHashIndexExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {
      `{"foo": 5}["foo"]`,
      5,
    },
    {
      `{"foo": 5}["bar"]`,
      nil,
    },
    {
      `hash = {"one": 1, "two": 2, "three": 3}
       hash["two"]`,
      2,
    },
    {
      `{4: 4}[4]`,
      4,
    },
    {
      `{true: 5}[true]`,
      5,
    },
    {
      `{false: 5}[false]`,
      5,
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

func TestHashIndexAssignment(t *testing.T) {
  input := `
  hash = {"existing": "value"}
  hash["new"] = "added"
  hash["existing"] = "updated"
  hash`

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
  }

  // Check that we have 2 key-value pairs
  if len(result.Pairs) != 2 {
    t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
  }

  // Check "new" key
  newKey := HashKey{Type: STRING_VALUE, Value: "new"}
  newValue, exists := result.Pairs[newKey]
  if !exists {
    t.Errorf("no pair for key 'new' in Pairs")
  } else {
    testStringObject(t, newValue, "added")
  }

  // Check "existing" key was updated
  existingKey := HashKey{Type: STRING_VALUE, Value: "existing"}
  existingValue, exists := result.Pairs[existingKey]
  if !exists {
    t.Errorf("no pair for key 'existing' in Pairs")
  } else {
    testStringObject(t, existingValue, "updated")
  }
}

func TestBuiltinHashKeys(t *testing.T) {
  input := `
  hash = {"name": "Alice", "age": 30, 42: "answer"}
  builtin_hash_keys(hash)`

  evaluated := testEval(input)
  result, ok := evaluated.(*Array)
  if !ok {
    t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Elements) != 3 {
    t.Fatalf("wrong num of elements. got=%d", len(result.Elements))
  }

  // Keys should be in insertion order
  expectedKeys := []interface{}{"name", "age", int64(42)}
  for i, expected := range expectedKeys {
    switch exp := expected.(type) {
    case string:
      testStringObject(t, result.Elements[i], exp)
    case int64:
      testIntegerObject(t, result.Elements[i], exp)
    }
  }
}

func TestBuiltinHashValues(t *testing.T) {
  input := `
  hash = {"name": "Alice", "age": 30}
  builtin_hash_values(hash)`

  evaluated := testEval(input)
  result, ok := evaluated.(*Array)
  if !ok {
    t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Elements) != 2 {
    t.Fatalf("wrong num of elements. got=%d", len(result.Elements))
  }

  // Values should be in same order as keys
  testStringObject(t, result.Elements[0], "Alice")
  testIntegerObject(t, result.Elements[1], 30)
}

func TestBuiltinHashHasKey(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {`builtin_hash_has_key({"name": "Alice"}, "name")`, true},
    {`builtin_hash_has_key({"name": "Alice"}, "age")`, false},
    {`builtin_hash_has_key({42: "answer"}, 42)`, true},
    {`builtin_hash_has_key({true: "yes"}, true)`, true},
    {`builtin_hash_has_key({true: "yes"}, false)`, false},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testBooleanObject(t, evaluated, tt.expected)
  }
}

func TestBuiltinHashGet(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {`builtin_hash_get({"name": "Alice"}, "name")`, "Alice"},
    {`builtin_hash_get({"name": "Alice"}, "age")`, nil},
    {`builtin_hash_get({"name": "Alice"}, "age", "unknown")`, "unknown"},
    {`builtin_hash_get({42: "answer"}, 42)`, "answer"},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    if tt.expected == nil {
      testNullObject(t, evaluated)
    } else if str, ok := tt.expected.(string); ok {
      testStringObject(t, evaluated, str)
    }
  }
}

func TestBuiltinHashSet(t *testing.T) {
  input := `
  original = {"name": "Alice"}
  updated = builtin_hash_set(original, "age", 30)
  updated`

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("object is not Hash. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Pairs) != 2 {
    t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
  }

  // Check both keys exist
  nameKey := HashKey{Type: STRING_VALUE, Value: "name"}
  ageKey := HashKey{Type: STRING_VALUE, Value: "age"}

  if _, exists := result.Pairs[nameKey]; !exists {
    t.Errorf("name key not found in result")
  }
  if _, exists := result.Pairs[ageKey]; !exists {
    t.Errorf("age key not found in result")
  }
}

func TestBuiltinHashDelete(t *testing.T) {
  input := `
  original = {"name": "Alice", "age": 30}
  result = builtin_hash_delete(original, "age")
  result`

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("object is not Hash. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Pairs) != 1 {
    t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
  }

  nameKey := HashKey{Type: STRING_VALUE, Value: "name"}
  ageKey := HashKey{Type: STRING_VALUE, Value: "age"}

  if _, exists := result.Pairs[nameKey]; !exists {
    t.Errorf("name key should still exist")
  }
  if _, exists := result.Pairs[ageKey]; exists {
    t.Errorf("age key should be deleted")
  }
}

func TestBuiltinHashMerge(t *testing.T) {
  input := `
  hash1 = {"name": "Alice", "age": 25}
  hash2 = {"age": 30, "city": "NYC"}
  result = builtin_hash_merge(hash1, hash2)
  result`

  evaluated := testEval(input)
  result, ok := evaluated.(*Hash)
  if !ok {
    t.Fatalf("object is not Hash. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Pairs) != 3 {
    t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
  }

  // Check that age was overridden by hash2
  ageKey := HashKey{Type: STRING_VALUE, Value: "age"}
  ageValue, exists := result.Pairs[ageKey]
  if !exists {
    t.Errorf("age key not found")
  } else {
    testIntegerObject(t, ageValue, 30) // Should be 30 from hash2, not 25 from hash1
  }
}

func TestHashLen(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {`len({})`, 0},
    {`len({"one": 1})`, 1},
    {`len({"one": 1, "two": 2, "three": 3})`, 3},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)
    testIntegerObject(t, evaluated, tt.expected)
  }
}

func TestHashWithDifferentKeyTypes(t *testing.T) {
  input := `hash = {"string": "value1", 42: "value2", true: "value3", 3.14: "value4"}
  [hash["string"], hash[42], hash[true], hash[3.14]]`

  evaluated := testEval(input)
  result, ok := evaluated.(*Array)
  if !ok {
    t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
  }

  if len(result.Elements) != 4 {
    t.Fatalf("wrong num of elements. got=%d", len(result.Elements))
  }

  testStringObject(t, result.Elements[0], "value1")
  testStringObject(t, result.Elements[1], "value2")
  testStringObject(t, result.Elements[2], "value3")
  testStringObject(t, result.Elements[3], "value4")
}

func testHashObject(t *testing.T, obj Value, expected map[HashKey]int64) bool {
  result, ok := obj.(*Hash)
  if !ok {
    t.Errorf("object is not Hash. got=%T (%+v)", obj, obj)
    return false
  }

  if len(result.Pairs) != len(expected) {
    t.Errorf("hash has wrong num of pairs. got=%d", len(result.Pairs))
    return false
  }

  for expectedKey, expectedValue := range expected {
    pair, ok := result.Pairs[expectedKey]
    if !ok {
      t.Errorf("no pair for given key in Pairs")
      return false
    }

    err := testIntegerObject(t, pair, expectedValue)
    if !err {
      return false
    }
  }

  return true
}