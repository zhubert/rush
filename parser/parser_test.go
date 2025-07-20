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