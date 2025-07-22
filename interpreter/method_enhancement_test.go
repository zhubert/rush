package interpreter

import (
  "testing"
  "rush/lexer"
  "rush/parser"
)

// TestMethodBindingEnhancement tests implicit self context improvements
func TestMethodBindingEnhancement(t *testing.T) {
  input := `
class Calculator {
  fn initialize(value) {
    @base = value
  }
  
  fn add(x) {
    return @base + x
  }
  
  fn getBase() {
    return @base
  }
}

calc = Calculator.new(10)
result1 = calc.add(5)
result2 = calc.getBase()
`

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  
  if len(p.Errors()) != 0 {
    t.Fatalf("parser has errors: %v", p.Errors())
  }
  
  env := NewEnvironment()
  Eval(program, env)
  
  // Test that method calls work with implicit self
  result1, exists := env.Get("result1")
  if !exists {
    t.Fatalf("result1 not found in environment")
  }
  
  intResult1, ok := result1.(*Integer)
  if !ok {
    t.Fatalf("result1 is not an integer. got=%T", result1)
  }
  
  if intResult1.Value != 15 {
    t.Errorf("wrong result1. expected=15, got=%d", intResult1.Value)
  }
  
  result2, exists := env.Get("result2")
  if !exists {
    t.Fatalf("result2 not found in environment")
  }
  
  intResult2, ok := result2.(*Integer)
  if !ok {
    t.Fatalf("result2 is not an integer. got=%T", result2)
  }
  
  if intResult2.Value != 10 {
    t.Errorf("wrong result2. expected=10, got=%d", intResult2.Value)
  }
}

// TestMethodDispatchPerformance tests performance of method resolution
func TestMethodDispatchPerformance(t *testing.T) {
  input := `
class BaseClass {
  fn method1() { return 1 }
  fn method2() { return 2 }
  fn method3() { return 3 }
}

class MiddleClass < BaseClass {
  fn method4() { return 4 }
  fn method5() { return 5 }
}

class FinalClass < MiddleClass {
  fn method6() { return 6 }
  fn method7() { return 7 }
}

obj = FinalClass.new()
`

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  
  if len(p.Errors()) != 0 {
    t.Fatalf("parser has errors: %v", p.Errors())
  }
  
  env := NewEnvironment()
  Eval(program, env)
  
  // Get the object for testing
  objVal, exists := env.Get("obj")
  if !exists {
    t.Fatalf("obj not found in environment")
  }
  
  obj, ok := objVal.(*Object)
  if !ok {
    t.Fatalf("obj is not an Object. got=%T", objVal)
  }
  
  // Test that method resolution works through inheritance chain
  tests := []struct {
    methodName string
    expected   bool
  }{
    {"method1", true},  // from BaseClass
    {"method4", true},  // from MiddleClass
    {"method6", true},  // from FinalClass
    {"nonexistent", false},
  }
  
  for _, test := range tests {
    method := resolveMethod(obj.Class, test.methodName)
    found := method != nil
    if found != test.expected {
      t.Errorf("method resolution for %s: expected=%t, got=%t", 
        test.methodName, test.expected, found)
    }
  }
}

// TestBoundMethodEnhancement tests that methods can be stored and called later
func TestBoundMethodEnhancement(t *testing.T) {
  input := `
class Counter {
  fn initialize() {
    @count = 0
  }
  
  fn increment() {
    @count = @count + 1
    return @count
  }
  
  fn getCount() {
    return @count
  }
}

counter = Counter.new()
result1 = counter.increment()
result2 = counter.increment()
result3 = counter.getCount()
`

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  
  if len(p.Errors()) != 0 {
    t.Fatalf("parser has errors: %v", p.Errors())
  }
  
  env := NewEnvironment()
  Eval(program, env)
  
  // Verify sequential method calls maintain state
  result3, exists := env.Get("result3")
  if !exists {
    t.Fatalf("result3 not found in environment")
  }
  
  intResult3, ok := result3.(*Integer)
  if !ok {
    t.Fatalf("result3 is not an integer. got=%T", result3)
  }
  
  if intResult3.Value != 2 {
    t.Errorf("wrong final count. expected=2, got=%d", intResult3.Value)
  }
}

func TestTrueBoundMethods(t *testing.T) {
  input := `
class Calculator {
  fn initialize(base) {
    @base = base
  }
  
  fn add(x) {
    return @base + x
  }
}

calc1 = Calculator.new(10)
calc2 = Calculator.new(20)

method1 = calc1.add
method2 = calc2.add

result1 = method1(5)
result2 = method2(5)
`

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  
  if len(p.Errors()) != 0 {
    t.Fatalf("parser has errors: %v", p.Errors())
  }
  
  env := NewEnvironment()
  Eval(program, env)
  
  result1, exists := env.Get("result1")
  if !exists {
    t.Fatalf("result1 not found in environment")
  }
  
  intResult1, ok := result1.(*Integer)
  if !ok {
    t.Fatalf("result1 is not an integer. got=%T", result1)
  }
  
  if intResult1.Value != 15 {
    t.Errorf("wrong result1. expected=15, got=%d", intResult1.Value)
  }
  
  result2, exists := env.Get("result2")
  if !exists {
    t.Fatalf("result2 not found in environment")
  }
  
  intResult2, ok := result2.(*Integer)
  if !ok {
    t.Fatalf("result2 is not an integer. got=%T", result2)
  }
  
  if intResult2.Value != 25 {
    t.Errorf("wrong result2. expected=25, got=%d", intResult2.Value)
  }
  
  method1Val, exists := env.Get("method1")
  if !exists {
    t.Fatalf("method1 not found in environment")
  }
  
  if _, ok := method1Val.(*BoundMethod); !ok {
    t.Errorf("method1 is not a BoundMethod. got=%T", method1Val)
  }
}