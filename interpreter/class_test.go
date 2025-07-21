package interpreter

import (
	"testing"
	"rush/lexer"
	"rush/parser"
	"rush/ast"
)

func TestBasicClassDefinition(t *testing.T) {
	input := `
class Person {
}
`
	
	evaluated := testEvalClass(input)
	class, ok := evaluated.(*Class)
	if !ok {
		t.Fatalf("object is not Class. got=%T (%+v)", evaluated, evaluated)
	}
	
	if class.Name != "Person" {
		t.Errorf("class name wrong. expected=%q, got=%q", "Person", class.Name)
	}
	
	if len(class.Methods) != 0 {
		t.Errorf("class should have no methods. got=%d", len(class.Methods))
	}
}

func TestClassWithMethods(t *testing.T) {
	input := `
class Calculator {
  fn add(x, y) {
    return x + y
  }
  
  fn multiply(x, y) {
    return x * y
  }
}
`
	
	evaluated := testEvalClass(input)
	class, ok := evaluated.(*Class)
	if !ok {
		t.Fatalf("object is not Class. got=%T (%+v)", evaluated, evaluated)
	}
	
	if class.Name != "Calculator" {
		t.Errorf("class name wrong. expected=%q, got=%q", "Calculator", class.Name)
	}
	
	expectedMethods := []string{"add", "multiply"}
	if len(class.Methods) != len(expectedMethods) {
		t.Errorf("class should have %d methods. got=%d", len(expectedMethods), len(class.Methods))
	}
	
	for _, methodName := range expectedMethods {
		if method, exists := class.Methods[methodName]; !exists {
			t.Errorf("method %q not found in class", methodName)
		} else if method == nil {
			t.Errorf("method %q is nil", methodName)
		}
	}
}

func TestObjectInstantiation(t *testing.T) {
	input := `
class Person {
}

person = Person.new()
person
`
	
	evaluated := testEvalClass(input)
	obj, ok := evaluated.(*Object)
	if !ok {
		t.Fatalf("object is not Object. got=%T (%+v)", evaluated, evaluated)
	}
	
	if obj.Class.Name != "Person" {
		t.Errorf("object class name wrong. expected=%q, got=%q", "Person", obj.Class.Name)
	}
	
	if obj.InstanceVars == nil {
		t.Errorf("object instance variables should be initialized")
	}
}

func TestInitializeMethod(t *testing.T) {
	input := `
class Person {
  fn initialize(name) {
    @name = name
  }
}

person = Person.new("Alice")
person
`
	
	evaluated := testEvalClass(input)
	obj, ok := evaluated.(*Object)
	if !ok {
		t.Fatalf("object is not Object. got=%T (%+v)", evaluated, evaluated)
	}
	
	// Check that instance variable was set
	if nameVal, exists := obj.InstanceVars["name"]; !exists {
		t.Errorf("instance variable 'name' should be set")
	} else if str, ok := nameVal.(*String); !ok {
		t.Errorf("instance variable 'name' should be String. got=%T", nameVal)
	} else if str.Value != "Alice" {
		t.Errorf("instance variable 'name' wrong. expected=%q, got=%q", "Alice", str.Value)
	}
}

func TestMethodCalls(t *testing.T) {
	input := `
class Person {
  fn initialize(name) {
    @name = name
  }
  
  fn getName() {
    return @name
  }
  
  fn setName(newName) {
    @name = newName
  }
}

person = Person.new("Bob")
name1 = person.getName()
person.setName("Robert")
name2 = person.getName()
[name1, name2]
`
	
	evaluated := testEvalClass(input)
	result, ok := evaluated.(*Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	
	if len(result.Elements) != 2 {
		t.Fatalf("array should have 2 elements. got=%d", len(result.Elements))
	}
	
	// Check first name
	if str1, ok := result.Elements[0].(*String); !ok {
		t.Errorf("first element should be String. got=%T", result.Elements[0])
	} else if str1.Value != "Bob" {
		t.Errorf("first name wrong. expected=%q, got=%q", "Bob", str1.Value)
	}
	
	// Check second name  
	if str2, ok := result.Elements[1].(*String); !ok {
		t.Errorf("second element should be String. got=%T", result.Elements[1])
	} else if str2.Value != "Robert" {
		t.Errorf("second name wrong. expected=%q, got=%q", "Robert", str2.Value)
	}
}

func TestInstanceVariableAccess(t *testing.T) {
	input := `
class Counter {
  fn initialize(start) {
    @count = start
  }
  
  fn increment() {
    @count = @count + 1
  }
  
  fn getCount() {
    return @count
  }
}

counter = Counter.new(5)
counter.increment()
counter.increment()
counter.getCount()
`
	
	evaluated := testEvalClass(input)
	result, ok := evaluated.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
	}
	
	if result.Value != 7 {
		t.Errorf("counter value wrong. expected=7, got=%d", result.Value)
	}
}

func TestMultipleInstances(t *testing.T) {
	input := `
class Counter {
  fn initialize(start) {
    @count = start
  }
  
  fn increment() {
    @count = @count + 1
  }
  
  fn getCount() {
    return @count
  }
}

counter1 = Counter.new(0)
counter2 = Counter.new(10)

counter1.increment()
counter1.increment()
counter2.increment()

[counter1.getCount(), counter2.getCount()]
`
	
	evaluated := testEvalClass(input)
	result, ok := evaluated.(*Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	
	if len(result.Elements) != 2 {
		t.Fatalf("array should have 2 elements. got=%d", len(result.Elements))
	}
	
	// Check counter1 value
	if int1, ok := result.Elements[0].(*Integer); !ok {
		t.Errorf("first element should be Integer. got=%T", result.Elements[0])
	} else if int1.Value != 2 {
		t.Errorf("counter1 value wrong. expected=2, got=%d", int1.Value)
	}
	
	// Check counter2 value
	if int2, ok := result.Elements[1].(*Integer); !ok {
		t.Errorf("second element should be Integer. got=%T", result.Elements[1])
	} else if int2.Value != 11 {
		t.Errorf("counter2 value wrong. expected=11, got=%d", int2.Value)
	}
}

func TestInstanceVariableOutsideObjectContext(t *testing.T) {
	input := `@name = "test"`
	
	evaluated := testEvalClass(input)
	errObj, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
	}
	
	expectedMsg := "instance variable @name used outside of object context"
	if errObj.Message != expectedMsg {
		t.Errorf("error message wrong. expected=%q, got=%q", expectedMsg, errObj.Message)
	}
}

func TestUndefinedMethod(t *testing.T) {
	input := `
class Person {
}

person = Person.new()
person.nonexistent()
`
	
	evaluated := testEvalClass(input)
	errObj, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
	}
	
	expectedMsg := "undefined method nonexistent for class Person"
	if errObj.Message != expectedMsg {
		t.Errorf("error message wrong. expected=%q, got=%q", expectedMsg, errObj.Message)
	}
}

func TestInitializeWithWrongArguments(t *testing.T) {
	input := `
class Person {
  fn initialize(name, age) {
    @name = name
    @age = age
  }
}

Person.new("Alice")
`
	
	evaluated := testEvalClass(input)
	errObj, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
	}
	
	expectedMsg := "wrong number of arguments for initialize: want=2, got=1"
	if errObj.Message != expectedMsg {
		t.Errorf("error message wrong. expected=%q, got=%q", expectedMsg, errObj.Message)
	}
}

func TestMethodWithWrongArguments(t *testing.T) {
	input := `
class Calculator {
  fn add(x, y) {
    return x + y
  }
}

calc = Calculator.new()
calc.add(5)
`
	
	evaluated := testEvalClass(input)
	errObj, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
	}
	
	expectedMsg := "wrong number of arguments: want=2, got=1"
	if errObj.Message != expectedMsg {
		t.Errorf("error message wrong. expected=%q, got=%q", expectedMsg, errObj.Message)
	}
}

func TestClassInheritancePlaceholder(t *testing.T) {
	// Test for inheritance syntax parsing (not functional yet in Phase 10.1)
	input := `
class Animal {
  fn initialize(name) {
    @name = name
  }
}

class Dog < Animal {
  fn bark() {
    return "woof"
  }
}
`
	
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	
	if len(program.Statements) != 2 {
		t.Fatalf("program should have 2 statements. got=%d", len(program.Statements))
	}
	
	// Check that Dog class has Animal as superclass in AST
	dogStmt := program.Statements[1]
	dogClass, ok := dogStmt.(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("second statement is not ClassDeclaration. got=%T", dogStmt)
	}
	
	if dogClass.SuperClass == nil {
		t.Errorf("Dog class should have superclass")
	} else if dogClass.SuperClass.Value != "Animal" {
		t.Errorf("Dog superclass wrong. expected=%q, got=%q", "Animal", dogClass.SuperClass.Value)
	}
}

func TestMethodResolutionOrder(t *testing.T) {
  input := `
  class Animal {
    fn speak() {
      return "generic sound"
    }
    
    fn move() {
      return "moving"
    }
  }
  
  class Dog < Animal {
    fn speak() {
      return "woof"
    }
    
    fn wagTail() {
      return "wagging"
    }
  }
  
  d = Dog.new()
  speak_result = d.speak()
  move_result = d.move()
  wag_result = d.wagTail()
  `

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  result := Eval(program, env)
  if isError(result) {
    t.Fatalf("evaluation failed: %s", result.Inspect())
  }

  // Test overridden method
  speakResult, exists := env.Get("speak_result")
  if !exists {
    t.Fatal("speak_result variable not found")
  }
  if speakResult.Inspect() != `woof` {
    t.Errorf("expected speak_result to be 'woof', got %s", speakResult.Inspect())
  }

  // Test inherited method
  moveResult, exists := env.Get("move_result")
  if !exists {
    t.Fatal("move_result variable not found")
  }
  if moveResult.Inspect() != `moving` {
    t.Errorf("expected move_result to be 'moving', got %s", moveResult.Inspect())
  }

  // Test child-specific method
  wagResult, exists := env.Get("wag_result")
  if !exists {
    t.Fatal("wag_result variable not found")
  }
  if wagResult.Inspect() != `wagging` {
    t.Errorf("expected wag_result to be 'wagging', got %s", wagResult.Inspect())
  }
}

func TestMultiLevelInheritance(t *testing.T) {
  input := `
  class Animal {
    fn speak() {
      return "generic sound"
    }
    
    fn breathe() {
      return "breathing"
    }
  }
  
  class Mammal < Animal {
    fn speak() {
      return "mammal sound"
    }
    
    fn nurture() {
      return "nursing young"
    }
  }
  
  class Dog < Mammal {
    fn speak() {
      return "woof"
    }
    
    fn wagTail() {
      return "wagging"
    }
  }
  
  d = Dog.new()
  speak_result = d.speak()
  breathe_result = d.breathe()
  nurture_result = d.nurture()
  wag_result = d.wagTail()
  `

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  result := Eval(program, env)
  if isError(result) {
    t.Fatalf("evaluation failed: %s", result.Inspect())
  }

  tests := []struct {
    varName  string
    expected string
  }{
    {"speak_result", `woof`},      // overridden in Dog
    {"breathe_result", `breathing`}, // inherited from Animal
    {"nurture_result", `nursing young`}, // inherited from Mammal
    {"wag_result", `wagging`},     // Dog-specific method
  }

  for _, tt := range tests {
    result, exists := env.Get(tt.varName)
    if !exists {
      t.Fatalf("%s variable not found", tt.varName)
    }
    if result.Inspect() != tt.expected {
      t.Errorf("expected %s to be %s, got %s", tt.varName, tt.expected, result.Inspect())
    }
  }
}

func TestSuperMethodCalls(t *testing.T) {
  input := `
  class Animal {
    fn speak() {
      return "generic sound"
    }
    
    fn move() {
      return "moving"
    }
  }
  
  class Dog < Animal {
    fn speak() {
      return super() + " - but woof!"
    }
    
    fn move() {
      return super() + " quickly"
    }
  }
  
  d = Dog.new()
  speak_result = d.speak()
  move_result = d.move()
  `

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  result := Eval(program, env)
  if isError(result) {
    t.Fatalf("evaluation failed: %s", result.Inspect())
  }

  // Test super call in overridden method
  speakResult, exists := env.Get("speak_result")
  if !exists {
    t.Fatal("speak_result variable not found")
  }
  if speakResult.Inspect() != `generic sound - but woof!` {
    t.Errorf("expected speak_result to be 'generic sound - but woof!', got %s", speakResult.Inspect())
  }

  // Test super call in another method
  moveResult, exists := env.Get("move_result")
  if !exists {
    t.Fatal("move_result variable not found")
  }
  if moveResult.Inspect() != `moving quickly` {
    t.Errorf("expected move_result to be 'moving quickly', got %s", moveResult.Inspect())
  }
}

func TestSuperWithArguments(t *testing.T) {
  input := `
  class Animal {
    fn speak(sound) {
      return "Animal says: " + sound
    }
  }
  
  class Dog < Animal {
    fn speak(sound) {
      return super(sound) + " loudly!"
    }
  }
  
  d = Dog.new()
  result = d.speak("woof")
  `

  l := lexer.New(input)
  p := parser.New(l)
  program := p.ParseProgram()
  env := NewEnvironment()

  evalResult := Eval(program, env)
  if isError(evalResult) {
    t.Fatalf("evaluation failed: %s", evalResult.Inspect())
  }

  result, exists := env.Get("result")
  if !exists {
    t.Fatal("result variable not found")
  }
  if result.Inspect() != `Animal says: woof loudly!` {
    t.Errorf("expected result to be 'Animal says: woof loudly!', got %s", result.Inspect())
  }
}

// Helper function for testing classes
func testEvalClass(input string) Value {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	
	return Eval(program, env)
}