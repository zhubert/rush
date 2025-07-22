package parser

import (
  "testing"
  "rush/lexer"
  "rush/ast"
)

func TestVariableAssignments(t *testing.T) {
  input := `
x = 5
name = "test"
flag = true
`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements. got=%d",
      len(program.Statements))
  }

  tests := []struct {
    expectedName string
  }{
    {"x"},
    {"name"},
    {"flag"},
  }

  for i, tt := range tests {
    stmt := program.Statements[i]
    if !testAssignStatement(t, stmt, tt.expectedName) {
      return
    }
  }
}

func TestArithmeticExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"5", "5"},
    {"10", "10"},
    {"-5", "(-5)"},
    {"-10", "(-10)"},
    {"5 + 5", "(5 + 5)"},
    {"5 - 5", "(5 - 5)"},
    {"5 * 5", "(5 * 5)"},
    {"5 / 5", "(5 / 5)"},
    {"2 * 3 + 4", "((2 * 3) + 4)"},
    {"2 + 3 * 4", "(2 + (3 * 4))"},
    {"(2 + 3) * 4", "((2 + 3) * 4)"},
  }

  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()

    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain 1 statement. got=%d",
        len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
        program.Statements[0])
    }

    if stmt.Expression.String() != tt.expected {
      t.Errorf("Expected %q, got %q", tt.expected, stmt.Expression.String())
    }
  }
}

func TestBooleanExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"true", "true"},
    {"false", "false"},
    {"3 > 5", "(3 > 5)"},
    {"3 < 5", "(3 < 5)"},
    {"3 == 3", "(3 == 3)"},
    {"3 != 5", "(3 != 5)"},
    {"3 <= 5", "(3 <= 5)"},
    {"3 >= 5", "(3 >= 5)"},
    {"true && false", "(true && false)"},
    {"true || false", "(true || false)"},
    {"!true", "(!true)"},
  }

  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()

    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain 1 statement. got=%d",
        len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
        program.Statements[0])
    }

    if stmt.Expression.String() != tt.expected {
      t.Errorf("Expected %q, got %q", tt.expected, stmt.Expression.String())
    }
  }
}

func TestIfExpressions(t *testing.T) {
  input := `if (x < y) { x } else { y }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
  }

  exp, ok := stmt.Expression.(*ast.IfExpression)
  if !ok {
    t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
      stmt.Expression)
  }

  if exp.Condition.String() != "(x < y)" {
    t.Errorf("exp.Condition is not %q. got=%q", "(x < y)", exp.Condition.String())
  }

  if len(exp.Consequence.Statements) != 1 {
    t.Errorf("consequence is not 1 statement. got=%d",
      len(exp.Consequence.Statements))
  }

  if len(exp.Alternative.Statements) != 1 {
    t.Errorf("alternative is not 1 statement. got=%d",
      len(exp.Alternative.Statements))
  }
}

func TestFunctionLiterals(t *testing.T) {
  input := `fn(x, y) { x + y }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
  }

  function, ok := stmt.Expression.(*ast.FunctionLiteral)
  if !ok {
    t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
      stmt.Expression)
  }

  if len(function.Parameters) != 2 {
    t.Fatalf("function literal parameters wrong. want 2, got=%d",
      len(function.Parameters))
  }

  if function.Parameters[0].Value != "x" {
    t.Fatalf("parameter[0] is not 'x'. got=%q", function.Parameters[0].Value)
  }

  if function.Parameters[1].Value != "y" {
    t.Fatalf("parameter[1] is not 'y'. got=%q", function.Parameters[1].Value)
  }

  if len(function.Body.Statements) != 1 {
    t.Fatalf("function.Body.Statements has not 1 statement. got=%d",
      len(function.Body.Statements))
  }
}

func TestCallExpressions(t *testing.T) {
  input := `add(1, 2 * 3, 4 + 5)`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
  }

  exp, ok := stmt.Expression.(*ast.CallExpression)
  if !ok {
    t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
      stmt.Expression)
  }

  if len(exp.Arguments) != 3 {
    t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
  }

  if exp.Arguments[0].String() != "1" {
    t.Errorf("argument 0 is not %q. got=%q", "1", exp.Arguments[0].String())
  }

  if exp.Arguments[1].String() != "(2 * 3)" {
    t.Errorf("argument 1 is not %q. got=%q", "(2 * 3)", exp.Arguments[1].String())
  }

  if exp.Arguments[2].String() != "(4 + 5)" {
    t.Errorf("argument 2 is not %q. got=%q", "(4 + 5)", exp.Arguments[2].String())
  }
}

func TestArrayLiterals(t *testing.T) {
  input := `[1, 2 * 2, 3 + 3]`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
  }

  array, ok := stmt.Expression.(*ast.ArrayLiteral)
  if !ok {
    t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
  }

  if len(array.Elements) != 3 {
    t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
  }

  if array.Elements[0].String() != "1" {
    t.Errorf("array.Elements[0] is not %q. got=%q", "1", array.Elements[0].String())
  }

  if array.Elements[1].String() != "(2 * 2)" {
    t.Errorf("array.Elements[1] is not %q. got=%q", "(2 * 2)", array.Elements[1].String())
  }

  if array.Elements[2].String() != "(3 + 3)" {
    t.Errorf("array.Elements[2] is not %q. got=%q", "(3 + 3)", array.Elements[2].String())
  }
}

func TestIndexExpressions(t *testing.T) {
  input := `myArray[1 + 1]`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
  }

  indexExp, ok := stmt.Expression.(*ast.IndexExpression)
  if !ok {
    t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
  }

  if indexExp.Left.String() != "myArray" {
    t.Errorf("indexExp.Left is not %q. got=%q", "myArray", indexExp.Left.String())
  }

  if indexExp.Index.String() != "(1 + 1)" {
    t.Errorf("indexExp.Index is not %q. got=%q", "(1 + 1)", indexExp.Index.String())
  }
}

func TestWhileStatements(t *testing.T) {
  input := `while (x < 10) { x = x + 1 }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.WhileStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
      program.Statements[0])
  }

  if stmt.Condition.String() != "(x < 10)" {
    t.Errorf("stmt.Condition is not %q. got=%q", "(x < 10)", stmt.Condition.String())
  }

  if len(stmt.Body.Statements) != 1 {
    t.Errorf("stmt.Body does not contain 1 statement. got=%d",
      len(stmt.Body.Statements))
  }
}

func TestForStatements(t *testing.T) {
  input := `for (i = 0; i < 10; i = i + 1) { sum = sum + i }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ForStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T",
      program.Statements[0])
  }

  if stmt.Condition.String() != "(i < 10)" {
    t.Errorf("stmt.Condition is not %q. got=%q", "(i < 10)", stmt.Condition.String())
  }

  if len(stmt.Body.Statements) != 1 {
    t.Errorf("stmt.Body does not contain 1 statement. got=%d",
      len(stmt.Body.Statements))
  }
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
  assignStmt, ok := s.(*ast.AssignmentStatement)
  if !ok {
    t.Errorf("s not *ast.AssignmentStatement. got=%T", s)
    return false
  }

  if assignStmt.Name.Value != name {
    t.Errorf("assignStmt.Name.Value not '%s'. got=%s", name, assignStmt.Name.Value)
    return false
  }

  if assignStmt.Name.TokenLiteral() != name {
    t.Errorf("assignStmt.Name.TokenLiteral() not '%s'. got=%s",
      name, assignStmt.Name.TokenLiteral())
    return false
  }

  return true
}

func TestClassDeclaration(t *testing.T) {
  input := `class Animal {
    fn speak() {
      return "sound"
    }
  }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ClassDeclaration)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ClassDeclaration. got=%T",
      program.Statements[0])
  }

  if stmt.Name.Value != "Animal" {
    t.Errorf("stmt.Name.Value not 'Animal'. got=%s", stmt.Name.Value)
  }

  if stmt.SuperClass != nil {
    t.Errorf("stmt.SuperClass should be nil for base class. got=%v", stmt.SuperClass)
  }

  if len(stmt.Body.Statements) != 1 {
    t.Errorf("stmt.Body does not contain 1 statement. got=%d",
      len(stmt.Body.Statements))
  }
}

func TestClassInheritance(t *testing.T) {
  input := `class Dog < Animal {
    fn bark() {
      return "woof"
    }
  }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ClassDeclaration)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ClassDeclaration. got=%T",
      program.Statements[0])
  }

  if stmt.Name.Value != "Dog" {
    t.Errorf("stmt.Name.Value not 'Dog'. got=%s", stmt.Name.Value)
  }

  if stmt.SuperClass == nil {
    t.Fatalf("stmt.SuperClass should not be nil for inherited class")
  }

  if stmt.SuperClass.Value != "Animal" {
    t.Errorf("stmt.SuperClass.Value not 'Animal'. got=%s", stmt.SuperClass.Value)
  }

  if len(stmt.Body.Statements) != 1 {
    t.Errorf("stmt.Body does not contain 1 statement. got=%d",
      len(stmt.Body.Statements))
  }
}

func TestMultiLevelInheritance(t *testing.T) {
  input := `
  class Animal {
    fn move() { return "moving" }
  }
  
  class Dog < Animal {
    fn bark() { return "woof" }
  }
  
  class GermanShepherd < Dog {
    fn guard() { return "guarding" }
  }
  `

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements. got=%d",
      len(program.Statements))
  }

  tests := []struct {
    className   string
    superClass  *string
  }{
    {"Animal", nil},
    {"Dog", stringPtr("Animal")},
    {"GermanShepherd", stringPtr("Dog")},
  }

  for i, tt := range tests {
    stmt, ok := program.Statements[i].(*ast.ClassDeclaration)
    if !ok {
      t.Fatalf("program.Statements[%d] is not ast.ClassDeclaration. got=%T",
        i, program.Statements[i])
    }

    if stmt.Name.Value != tt.className {
      t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.className, stmt.Name.Value)
    }

    if tt.superClass == nil {
      if stmt.SuperClass != nil {
        t.Errorf("stmt.SuperClass should be nil for %s. got=%v", tt.className, stmt.SuperClass)
      }
    } else {
      if stmt.SuperClass == nil {
        t.Fatalf("stmt.SuperClass should not be nil for %s", tt.className)
      }
      if stmt.SuperClass.Value != *tt.superClass {
        t.Errorf("stmt.SuperClass.Value not '%s' for %s. got=%s", 
          *tt.superClass, tt.className, stmt.SuperClass.Value)
      }
    }
  }
}

func TestSwitchStatements(t *testing.T) {
  input := `switch (grade) {
    case "A":
      result = "excellent"
    case "B", "C":
      result = "good"
    default:
      result = "other"
  }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(p.Errors()) > 0 {
    for _, err := range p.Errors() {
      t.Errorf("Parser error: %s", err)
    }
    t.FailNow()
  }

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.SwitchStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.SwitchStatement. got=%T",
      program.Statements[0])
  }

  if stmt.Value.String() != "grade" {
    t.Errorf("stmt.Value is not 'grade'. got=%q", stmt.Value.String())
  }

  if len(stmt.Cases) != 2 {
    t.Fatalf("stmt.Cases does not contain 2 cases. got=%d", len(stmt.Cases))
  }

  // Test first case: case "A"
  firstCase := stmt.Cases[0]
  if len(firstCase.Values) != 1 {
    t.Errorf("firstCase.Values should have 1 value. got=%d", len(firstCase.Values))
  }
  if firstCase.Values[0].String() != "\"A\"" {
    t.Errorf("firstCase.Values[0] should be '\"A\"'. got=%s", firstCase.Values[0].String())
  }

  // Test second case: case "B", "C"
  secondCase := stmt.Cases[1]
  if len(secondCase.Values) != 2 {
    t.Errorf("secondCase.Values should have 2 values. got=%d", len(secondCase.Values))
  }
  if secondCase.Values[0].String() != "\"B\"" {
    t.Errorf("secondCase.Values[0] should be '\"B\"'. got=%s", secondCase.Values[0].String())
  }
  if secondCase.Values[1].String() != "\"C\"" {
    t.Errorf("secondCase.Values[1] should be '\"C\"'. got=%s", secondCase.Values[1].String())
  }

  // Test default clause
  if stmt.Default == nil {
    t.Fatal("stmt.Default should not be nil")
  }
  if len(stmt.Default.Body.Statements) != 1 {
    t.Errorf("default clause should have 1 statement. got=%d", len(stmt.Default.Body.Statements))
  }
}

func TestSwitchWithoutDefault(t *testing.T) {
  input := `switch (x) {
    case 1:
      y = 10
    case 2:
      y = 20
  }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(p.Errors()) > 0 {
    for _, err := range p.Errors() {
      t.Errorf("Parser error: %s", err)
    }
    t.FailNow()
  }

  if len(program.Statements) != 1 {
    t.Fatalf("program.Statements does not contain 1 statement. got=%d",
      len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.SwitchStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.SwitchStatement. got=%T",
      program.Statements[0])
  }

  if len(stmt.Cases) != 2 {
    t.Fatalf("stmt.Cases does not contain 2 cases. got=%d", len(stmt.Cases))
  }

  if stmt.Default != nil {
    t.Error("stmt.Default should be nil when no default clause is present")
  }
}

func TestSwitchIntegerCases(t *testing.T) {
  input := `switch (num) {
    case 1, 2, 3:
      result = "low"
    case 42:
      result = "answer"
  }`

  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()

  if len(p.Errors()) > 0 {
    for _, err := range p.Errors() {
      t.Errorf("Parser error: %s", err)
    }
    t.FailNow()
  }

  stmt, ok := program.Statements[0].(*ast.SwitchStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.SwitchStatement. got=%T",
      program.Statements[0])
  }

  // Test first case with multiple integer values
  firstCase := stmt.Cases[0]
  if len(firstCase.Values) != 3 {
    t.Errorf("firstCase.Values should have 3 values. got=%d", len(firstCase.Values))
  }

  expectedValues := []string{"1", "2", "3"}
  for i, expectedValue := range expectedValues {
    if firstCase.Values[i].String() != expectedValue {
      t.Errorf("firstCase.Values[%d] should be '%s'. got=%s", i, expectedValue, firstCase.Values[i].String())
    }
  }
}

func TestSwitchParseErrors(t *testing.T) {
  tests := []struct {
    input    string
    errorMsg string
  }{
    {
      `switch (x) {
        case 1:
          y = 10
        default:
          y = 20
        default:
          y = 30
      }`,
      "switch statement can only have one default clause",
    },
  }

  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    p.ParseProgram()

    if len(p.Errors()) == 0 {
      t.Errorf("Expected parser errors for input: %s", tt.input)
      continue
    }

    found := false
    for _, err := range p.Errors() {
      if err == tt.errorMsg || len(err) > len(tt.errorMsg) {
        found = true
        break
      }
    }

    if !found {
      t.Errorf("Expected error containing '%s', but got: %v", tt.errorMsg, p.Errors())
    }
  }
}

func stringPtr(s string) *string {
  return &s
}

func TestImportAliasingParsing(t *testing.T) {
  tests := []struct {
    name   string
    input  string
    expected func(*ast.ImportStatement) bool
  }{
    {
      name:  "Basic import without alias",
      input: `import { foo } from "module"`,
      expected: func(stmt *ast.ImportStatement) bool {
        if len(stmt.Items) != 1 {
          t.Errorf("Expected 1 import item, got %d", len(stmt.Items))
          return false
        }
        item := stmt.Items[0]
        if item.Name.Value != "foo" {
          t.Errorf("Expected import name 'foo', got %s", item.Name.Value)
          return false
        }
        if item.Alias != nil {
          t.Errorf("Expected no alias, but got alias: %s", item.Alias.Value)
          return false
        }
        return true
      },
    },
    {
      name:  "Import with alias",
      input: `import { foo as bar } from "module"`,
      expected: func(stmt *ast.ImportStatement) bool {
        if len(stmt.Items) != 1 {
          t.Errorf("Expected 1 import item, got %d", len(stmt.Items))
          return false
        }
        item := stmt.Items[0]
        if item.Name.Value != "foo" {
          t.Errorf("Expected import name 'foo', got %s", item.Name.Value)
          return false
        }
        if item.Alias == nil {
          t.Errorf("Expected alias, but got nil")
          return false
        }
        if item.Alias.Value != "bar" {
          t.Errorf("Expected alias 'bar', got %s", item.Alias.Value)
          return false
        }
        return true
      },
    },
    {
      name:  "Multiple imports with mixed aliases",
      input: `import { foo, bar as baz, qux } from "module"`,
      expected: func(stmt *ast.ImportStatement) bool {
        if len(stmt.Items) != 3 {
          t.Errorf("Expected 3 import items, got %d", len(stmt.Items))
          return false
        }
        
        // Check first item (no alias)
        if stmt.Items[0].Name.Value != "foo" || stmt.Items[0].Alias != nil {
          t.Errorf("First item should be 'foo' with no alias")
          return false
        }
        
        // Check second item (with alias)
        if stmt.Items[1].Name.Value != "bar" || stmt.Items[1].Alias == nil || stmt.Items[1].Alias.Value != "baz" {
          t.Errorf("Second item should be 'bar as baz'")
          return false
        }
        
        // Check third item (no alias)
        if stmt.Items[2].Name.Value != "qux" || stmt.Items[2].Alias != nil {
          t.Errorf("Third item should be 'qux' with no alias")
          return false
        }
        
        return true
      },
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      l := lexer.New(tt.input)
      p := New(l)
      program := p.ParseProgram()

      if len(p.Errors()) > 0 {
        t.Fatalf("Parser errors: %v", p.Errors())
      }

      if len(program.Statements) != 1 {
        t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
      }

      stmt, ok := program.Statements[0].(*ast.ImportStatement)
      if !ok {
        t.Fatalf("Expected ImportStatement, got %T", program.Statements[0])
      }

      if !tt.expected(stmt) {
        t.Fatalf("Test validation failed")
      }
    })
  }
}