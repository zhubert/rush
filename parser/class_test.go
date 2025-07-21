package parser

import (
	"testing"
	"rush/ast"
	"rush/lexer"
)

func TestClassDefinitionParsing(t *testing.T) {
	input := `class Person {}`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDeclaration. got=%T", program.Statements[0])
	}
	
	if stmt.Name.Value != "Person" {
		t.Errorf("class name is not 'Person'. got=%q", stmt.Name.Value)
	}
	
	if stmt.SuperClass != nil {
		t.Errorf("class should not have superclass. got=%v", stmt.SuperClass)
	}
}

func TestClassInheritanceParsing(t *testing.T) {
	input := `class Dog < Animal {}`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDeclaration. got=%T", program.Statements[0])
	}
	
	if stmt.Name.Value != "Dog" {
		t.Errorf("class name is not 'Dog'. got=%q", stmt.Name.Value)
	}
	
	if stmt.SuperClass == nil {
		t.Fatalf("class should have superclass")
	}
	
	if stmt.SuperClass.Value != "Animal" {
		t.Errorf("superclass name is not 'Animal'. got=%q", stmt.SuperClass.Value)
	}
}

func TestMethodDefinitionParsing(t *testing.T) {
	input := `
class Calculator {
  fn add(x, y) {
    return x + y
  }
  
  fn initialize(value) {
    @value = value
  }
}
`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDeclaration. got=%T", program.Statements[0])
	}
	
	if len(stmt.Body.Statements) != 2 {
		t.Fatalf("class body should have 2 statements. got=%d", len(stmt.Body.Statements))
	}
	
	// Test first method (add)
	method1, ok := stmt.Body.Statements[0].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("first statement is not MethodDeclaration. got=%T", stmt.Body.Statements[0])
	}
	
	if method1.Name.Value != "add" {
		t.Errorf("first method name wrong. expected='add', got=%q", method1.Name.Value)
	}
	
	if len(method1.Parameters) != 2 {
		t.Errorf("first method should have 2 parameters. got=%d", len(method1.Parameters))
	}
	
	expectedParams := []string{"x", "y"}
	for i, expectedParam := range expectedParams {
		if method1.Parameters[i].Value != expectedParam {
			t.Errorf("parameter %d wrong. expected=%q, got=%q", i, expectedParam, method1.Parameters[i].Value)
		}
	}
	
	// Test second method (initialize)
	method2, ok := stmt.Body.Statements[1].(*ast.MethodDeclaration)
	if !ok {
		t.Fatalf("second statement is not MethodDeclaration. got=%T", stmt.Body.Statements[1])
	}
	
	if method2.Name.Value != "initialize" {
		t.Errorf("second method name wrong. expected='initialize', got=%q", method2.Name.Value)
	}
	
	if len(method2.Parameters) != 1 {
		t.Errorf("second method should have 1 parameter. got=%d", len(method2.Parameters))
	}
	
	if method2.Parameters[0].Value != "value" {
		t.Errorf("parameter name wrong. expected='value', got=%q", method2.Parameters[0].Value)
	}
}

func TestInstanceVariableParsing(t *testing.T) {
	input := `@name`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	instVar, ok := stmt.Expression.(*ast.InstanceVariable)
	if !ok {
		t.Fatalf("expression is not ast.InstanceVariable. got=%T", stmt.Expression)
	}
	
	if instVar.Name.Value != "name" {
		t.Errorf("instance variable name wrong. expected='name', got=%q", instVar.Name.Value)
	}
}

func TestInstanceVariableAssignmentParsing(t *testing.T) {
	input := `@name = "Alice"`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.AssignmentStatement. got=%T", program.Statements[0])
	}
	
	if stmt.Name.Value != "@name" {
		t.Errorf("assignment name wrong. expected='@name', got=%q", stmt.Name.Value)
	}
	
	strLit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("assignment value is not ast.StringLiteral. got=%T", stmt.Value)
	}
	
	if strLit.Value != "Alice" {
		t.Errorf("string value wrong. expected='Alice', got=%q", strLit.Value)
	}
}

func TestNewExpressionParsing(t *testing.T) {
	input := `Person.new()`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	newExpr, ok := stmt.Expression.(*ast.NewExpression)
	if !ok {
		t.Fatalf("expression is not ast.NewExpression. got=%T", stmt.Expression)
	}
	
	if newExpr.ClassName.Value != "Person" {
		t.Errorf("class name wrong. expected='Person', got=%q", newExpr.ClassName.Value)
	}
	
	if len(newExpr.Arguments) != 0 {
		t.Errorf("new expression should have 0 arguments. got=%d", len(newExpr.Arguments))
	}
}

func TestNewExpressionWithArgumentsParsing(t *testing.T) {
	input := `Person.new("Alice", 30)`
	
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	newExpr, ok := stmt.Expression.(*ast.NewExpression)
	if !ok {
		t.Fatalf("expression is not ast.NewExpression. got=%T", stmt.Expression)
	}
	
	if newExpr.ClassName.Value != "Person" {
		t.Errorf("class name wrong. expected='Person', got=%q", newExpr.ClassName.Value)
	}
	
	if len(newExpr.Arguments) != 2 {
		t.Fatalf("new expression should have 2 arguments. got=%d", len(newExpr.Arguments))
	}
	
	// Check first argument
	strArg, ok := newExpr.Arguments[0].(*ast.StringLiteral)
	if !ok {
		t.Errorf("first argument should be StringLiteral. got=%T", newExpr.Arguments[0])
	} else if strArg.Value != "Alice" {
		t.Errorf("first argument wrong. expected='Alice', got=%q", strArg.Value)
	}
	
	// Check second argument
	intArg, ok := newExpr.Arguments[1].(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("second argument should be IntegerLiteral. got=%T", newExpr.Arguments[1])
	} else if intArg.Value != 30 {
		t.Errorf("second argument wrong. expected=30, got=%d", intArg.Value)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}