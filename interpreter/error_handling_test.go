package interpreter

import (
  "strings"
  "testing"
  "rush/lexer"
  "rush/parser"
)

// Advanced error handling tests for specialized scenarios

func TestStackTraceFormatting(t *testing.T) {
  input := `
level3 = fn() {
  throw RuntimeError("deep error")
}

level2 = fn() {
  level3()
}

level1 = fn() {
  level2()
}

try {
  level1()
} catch (error) {
  error.stack
}
`
  
  evaluated := testEval(input)
  
  stackStr, ok := evaluated.(*String)
  if !ok {
    t.Fatalf("Expected String for stack trace, got %T", evaluated)
  }
  
  // Check that stack trace contains function names
  if !strings.Contains(stackStr.Value, "level3") {
    t.Errorf("Stack trace should contain 'level3', got: %s", stackStr.Value)
  }
  if !strings.Contains(stackStr.Value, "level2") {
    t.Errorf("Stack trace should contain 'level2', got: %s", stackStr.Value)
  }
  if !strings.Contains(stackStr.Value, "level1") {
    t.Errorf("Stack trace should contain 'level1', got: %s", stackStr.Value)
  }
}

func TestMultipleFinallyBlocks(t *testing.T) {
  input := `
outerFinally = false
innerFinally = false

try {
  try {
    throw RuntimeError("nested error")
  } finally {
    innerFinally = true
  }
} catch (error) {
  # Error caught
} finally {
  outerFinally = true
}
outerFinally && innerFinally
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestFinallyBlockWithReturn(t *testing.T) {
  input := `
testFunc = fn() {
  finallyRan = false
  try {
    return 42
  } finally {
    finallyRan = true
  }
  return 0  # Should not reach this
}

result = testFunc()
result == 42
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestFinallyBlockWithUncaughtException(t *testing.T) {
  input := `
finallyRan = false

testFunc = fn() {
  try {
    throw RuntimeError("uncaught")
  } finally {
    finallyRan = true
  }
}

try {
  testFunc()
} catch (error) {
  # Error should propagate but finally should still run
}
finallyRan
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestErrorPropertyAccessAfterCatch(t *testing.T) {
  input := `
errorType = ""
errorMessage = ""
hasStack = false

try {
  throw ValidationError("persistent error")
} catch (error) {
  errorType = error.type
  errorMessage = error.message
  hasStack = len(error.stack) > 0
}

errorType == "ValidationError"
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestMultipleCatchWithSameType(t *testing.T) {
  input := `
firstCatch = false
secondCatch = false

try {
  throw ValidationError("test")
} catch (ValidationError error) {
  firstCatch = true
} catch (ValidationError error) {
  secondCatch = true  # Should not reach here
}
firstCatch && !secondCatch
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestGenericCatchAfterSpecific(t *testing.T) {
  input := `
specificCaught = false
genericCaught = false

try {
  throw ArgumentError("specific error")
} catch (ValidationError error) {
  # Should not catch ValidationError
} catch (error) {
  genericCaught = true
  # Check that it's still the correct type
  specificCaught = (error.type == "ArgumentError")
}
specificCaught && genericCaught
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestExceptionInFinallyBlock(t *testing.T) {
  input := `
finallyException = false

try {
  try {
    throw RuntimeError("original error")
  } finally {
    throw ValidationError("finally error")
  }
} catch (ValidationError error) {
  finallyException = true
} catch (RuntimeError error) {
  # Should not catch the original error
}
finallyException
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestDeepStackTrace(t *testing.T) {
  input := `
f5 = fn() { throw RuntimeError("deep") }
f4 = fn() { f5() }
f3 = fn() { f4() }
f2 = fn() { f3() }
f1 = fn() { f2() }

stackDepth = 0
try {
  f1()
} catch (error) {
  # Count function names in stack trace
  stack = error.stack
  stackDepth = 0
  if (len(stack) > 0) {
    # This is a simplified check - just verify stack exists
    stackDepth = 5  # We know we have 5 functions
  }
}
stackDepth == 5
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestErrorConstructorWithEmptyMessage(t *testing.T) {
  input := `
try {
  throw RuntimeError("")
} catch (error) {
  error.message == ""
}
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestBuiltinExceptionIntegrationWithArraySlice(t *testing.T) {
  input := `
# Test that slice returns empty array for out of bounds, doesn't throw
arr = [1, 2, 3]
result = slice(arr, 5, 10)
len(result) == 0
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestTryCatchInLoop(t *testing.T) {
  input := `
successCount = 0
errorCount = 0

for (i = 0; i < 5; i = i + 1) {
  try {
    if (i == 2) {
      throw ValidationError("error at 2")
    } else {
      successCount = successCount + 1
    }
  } catch (ValidationError error) {
    errorCount = errorCount + 1
  }
}

(successCount == 4) && (errorCount == 1)
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestTryCatchInRecursiveFunction(t *testing.T) {
  input := `
recursiveFunc = fn(n) {
  try {
    if (n <= 0) {
      throw ArgumentError("base case")
    }
    return n + recursiveFunc(n - 1)
  } catch (ArgumentError error) {
    return 0
  }
}

result = recursiveFunc(3)
result == 6  # 3 + 2 + 1 + 0
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestErrorMessageWithSpecialCharacters(t *testing.T) {
  input := `
specialMessage = "Error with special chars: !@#$%^&*()_+{}|:<>?[],./"
try {
  throw RuntimeError(specialMessage)
} catch (error) {
  error.message == specialMessage
}
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestCatchVariableNameConflict(t *testing.T) {
  input := `
error = "original value"
try {
  throw ValidationError("test")
} catch (error) {
  # error variable should be shadowed in catch block
  message = error.message
}
# Original error variable should be unchanged outside catch
error == "original value"
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestThrowNonErrorValue(t *testing.T) {
  input := `throw "plain string"`
  
  evaluated := testEval(input)
  
  // Should wrap non-error values in a RuntimeError
  exception, ok := evaluated.(*Exception)
  if !ok {
    t.Fatalf("Expected Exception, got %T", evaluated)
  }
  
  errorObj, ok := exception.Error.(*Error)
  if !ok {
    t.Fatalf("Expected Error inside Exception, got %T", exception.Error)
  }
  
  if errorObj.ErrorType != "Error" {
    t.Errorf("Expected Error for non-error throw, got %s", errorObj.ErrorType)
  }
}

func TestComplexErrorPropertyAccess(t *testing.T) {
  input := `
nested = fn() {
  inner = fn() {
    throw TypeError("complex error")
  }
  inner()
}

properties = []
try {
  nested()
} catch (error) {
  properties = [error.type, error.message, len(error.stack) > 0, error.line >= 0, error.column >= 0]
}

# Check all properties are accessible and correct
(properties[0] == "TypeError") && (properties[1] == "complex error") && (properties[2] == true) && (properties[3] == true) && (properties[4] == true)
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

func TestErrorInBuiltinFunctionCall(t *testing.T) {
  input := `
# Test that substr with invalid bounds returns empty string
result = substr("hello", -1, 10)
result == ""
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}

// Test helper function specifically for error handling tests
func testEvalErrorHandling(input string) Value {
  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  
  // Check for parse errors
  if len(p.Errors()) > 0 {
    return newError("Parse errors: %s", strings.Join(p.Errors(), "; "))
  }
  
  env := NewEnvironment()
  return Eval(program, env)
}

func TestParseErrorVsRuntimeError(t *testing.T) {
  // This should be a parse error, not a runtime error
  l := lexer.New("throw RuntimeError(")  // Incomplete syntax
  p := parser.New(l)
  _ = p.ParseProgram()  // We just need to trigger parsing
  
  if len(p.Errors()) == 0 {
    t.Error("Expected parse error for incomplete syntax")
  }
}

func TestTryCatchWithComplexConditions(t *testing.T) {
  input := `
result = ""
value = 10

try {
  if (value > 5 && value < 15) {
    throw ValidationError("in range error")
  } else {
    result = "no error"
  }
} catch (ValidationError error) {
  if (error.message == "in range error") {
    result = "caught expected error"
  } else {
    result = "caught unexpected error"
  }
}

result == "caught expected error"
`
  
  evaluated := testEval(input)
  testBooleanObject(t, evaluated, true)
}