package ast

import (
	"testing"
	"rush/lexer"
)

func TestProgramString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignmentStatement{
				Token: lexer.Token{Type: lexer.IDENT, Literal: "x"},
				Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
				Value: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "5"}, Value: 5},
			},
		},
	}

	expected := "x = 5"
	if program.String() != expected {
		t.Errorf("program.String() wrong. got=%q, want=%q", program.String(), expected)
	}
}

func TestAssignmentStatementString(t *testing.T) {
	stmt := &AssignmentStatement{
		Token: lexer.Token{Type: lexer.IDENT, Literal: "variable"},
		Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "variable"}, Value: "variable"},
		Value: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "hello"}, Value: "hello"},
	}

	expected := "variable = \"hello\""
	if stmt.String() != expected {
		t.Errorf("stmt.String() wrong. got=%q, want=%q", stmt.String(), expected)
	}

	if stmt.TokenLiteral() != "variable" {
		t.Errorf("stmt.TokenLiteral() wrong. got=%q, want=%q", stmt.TokenLiteral(), "variable")
	}
}

func TestIndexAssignmentStatementString(t *testing.T) {
	stmt := &IndexAssignmentStatement{
		Token: lexer.Token{Type: lexer.ASSIGN, Literal: "="},
		Left: &IndexExpression{
			Token: lexer.Token{Type: lexer.LBRACKET, Literal: "["},
			Left:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "arr"}, Value: "arr"},
			Index: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "0"}, Value: 0},
		},
		Value: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "42"}, Value: 42},
	}

	expected := "(arr[0]) = 42"
	if stmt.String() != expected {
		t.Errorf("stmt.String() wrong. got=%q, want=%q", stmt.String(), expected)
	}

	if stmt.TokenLiteral() != "=" {
		t.Errorf("stmt.TokenLiteral() wrong. got=%q, want=%q", stmt.TokenLiteral(), "=")
	}
}

func TestIdentifierString(t *testing.T) {
	ident := &Identifier{
		Token: lexer.Token{Type: lexer.IDENT, Literal: "foobar"},
		Value: "foobar",
	}

	if ident.String() != "foobar" {
		t.Errorf("ident.String() wrong. got=%q, want=%q", ident.String(), "foobar")
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() wrong. got=%q, want=%q", ident.TokenLiteral(), "foobar")
	}
}

func TestIntegerLiteralString(t *testing.T) {
	intLit := &IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: "42"},
		Value: 42,
	}

	if intLit.String() != "42" {
		t.Errorf("intLit.String() wrong. got=%q, want=%q", intLit.String(), "42")
	}

	if intLit.TokenLiteral() != "42" {
		t.Errorf("intLit.TokenLiteral() wrong. got=%q, want=%q", intLit.TokenLiteral(), "42")
	}
}

func TestFloatLiteralString(t *testing.T) {
	floatLit := &FloatLiteral{
		Token: lexer.Token{Type: lexer.FLOAT, Literal: "3.14"},
		Value: 3.14,
	}

	if floatLit.String() != "3.14" {
		t.Errorf("floatLit.String() wrong. got=%q, want=%q", floatLit.String(), "3.14")
	}

	if floatLit.TokenLiteral() != "3.14" {
		t.Errorf("floatLit.TokenLiteral() wrong. got=%q, want=%q", floatLit.TokenLiteral(), "3.14")
	}
}

func TestStringLiteralString(t *testing.T) {
	strLit := &StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: "hello world"},
		Value: "hello world",
	}

	expected := "\"hello world\""
	if strLit.String() != expected {
		t.Errorf("strLit.String() wrong. got=%q, want=%q", strLit.String(), expected)
	}

	if strLit.TokenLiteral() != "hello world" {
		t.Errorf("strLit.TokenLiteral() wrong. got=%q, want=%q", strLit.TokenLiteral(), "hello world")
	}
}

func TestBooleanLiteralString(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		boolLit := &BooleanLiteral{
			Token: lexer.Token{Type: lexer.TRUE, Literal: tt.expected},
			Value: tt.input,
		}

		if boolLit.String() != tt.expected {
			t.Errorf("boolLit.String() wrong. got=%q, want=%q", boolLit.String(), tt.expected)
		}

		if boolLit.TokenLiteral() != tt.expected {
			t.Errorf("boolLit.TokenLiteral() wrong. got=%q, want=%q", boolLit.TokenLiteral(), tt.expected)
		}
	}
}

func TestInfixExpressionString(t *testing.T) {
	infixExp := &InfixExpression{
		Token:    lexer.Token{Type: lexer.PLUS, Literal: "+"},
		Left:     &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "5"}, Value: 5},
		Operator: "+",
		Right:    &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "10"}, Value: 10},
	}

	expected := "(5 + 10)"
	if infixExp.String() != expected {
		t.Errorf("infixExp.String() wrong. got=%q, want=%q", infixExp.String(), expected)
	}

	if infixExp.TokenLiteral() != "+" {
		t.Errorf("infixExp.TokenLiteral() wrong. got=%q, want=%q", infixExp.TokenLiteral(), "+")
	}
}

func TestPrefixExpressionString(t *testing.T) {
	prefixExp := &PrefixExpression{
		Token:    lexer.Token{Type: lexer.MINUS, Literal: "-"},
		Operator: "-",
		Right:    &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "5"}, Value: 5},
	}

	expected := "(-5)"
	if prefixExp.String() != expected {
		t.Errorf("prefixExp.String() wrong. got=%q, want=%q", prefixExp.String(), expected)
	}

	if prefixExp.TokenLiteral() != "-" {
		t.Errorf("prefixExp.TokenLiteral() wrong. got=%q, want=%q", prefixExp.TokenLiteral(), "-")
	}
}

func TestArrayLiteralString(t *testing.T) {
	arrayLit := &ArrayLiteral{
		Token: lexer.Token{Type: lexer.LBRACKET, Literal: "["},
		Elements: []Expression{
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "2"}, Value: 2},
			&StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "three"}, Value: "three"},
		},
	}

	expected := "[1, 2, \"three\"]"
	if arrayLit.String() != expected {
		t.Errorf("arrayLit.String() wrong. got=%q, want=%q", arrayLit.String(), expected)
	}

	if arrayLit.TokenLiteral() != "[" {
		t.Errorf("arrayLit.TokenLiteral() wrong. got=%q, want=%q", arrayLit.TokenLiteral(), "[")
	}
}

func TestHashLiteralString(t *testing.T) {
	hashLit := &HashLiteral{
		Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
		Pairs: []HashPair{
			{
				Key:   &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "name"}, Value: "name"},
				Value: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "Alice"}, Value: "Alice"},
			},
			{
				Key:   &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "42"}, Value: 42},
				Value: &BooleanLiteral{Token: lexer.Token{Type: lexer.TRUE, Literal: "true"}, Value: true},
			},
		},
	}

	expected := "{\"name\": \"Alice\", 42: true}"
	if hashLit.String() != expected {
		t.Errorf("hashLit.String() wrong. got=%q, want=%q", hashLit.String(), expected)
	}

	if hashLit.TokenLiteral() != "{" {
		t.Errorf("hashLit.TokenLiteral() wrong. got=%q, want=%q", hashLit.TokenLiteral(), "{")
	}
}

func TestExpressionStatementString(t *testing.T) {
	exprStmt := &ExpressionStatement{
		Token:      lexer.Token{Type: lexer.INT, Literal: "5"},
		Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "5"}, Value: 5},
	}

	expected := "5"
	if exprStmt.String() != expected {
		t.Errorf("exprStmt.String() wrong. got=%q, want=%q", exprStmt.String(), expected)
	}

	if exprStmt.TokenLiteral() != "5" {
		t.Errorf("exprStmt.TokenLiteral() wrong. got=%q, want=%q", exprStmt.TokenLiteral(), "5")
	}
}

func TestIfExpressionString(t *testing.T) {
	ifExpr := &IfExpression{
		Token:     lexer.Token{Type: lexer.IF, Literal: "if"},
		Condition: &BooleanLiteral{Token: lexer.Token{Type: lexer.TRUE, Literal: "true"}, Value: true},
		Consequence: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "10"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "10"}, Value: 10},
				},
			},
		},
		Alternative: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "20"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "20"}, Value: 20},
				},
			},
		},
	}

	expected := "iftrue {10}else {20}"
	if ifExpr.String() != expected {
		t.Errorf("ifExpr.String() wrong. got=%q, want=%q", ifExpr.String(), expected)
	}

	if ifExpr.TokenLiteral() != "if" {
		t.Errorf("ifExpr.TokenLiteral() wrong. got=%q, want=%q", ifExpr.TokenLiteral(), "if")
	}
}

func TestBlockStatementString(t *testing.T) {
	blockStmt := &BlockStatement{
		Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
		Statements: []Statement{
			&ExpressionStatement{
				Token:      lexer.Token{Type: lexer.INT, Literal: "5"},
				Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "5"}, Value: 5},
			},
			&ExpressionStatement{
				Token:      lexer.Token{Type: lexer.INT, Literal: "10"},
				Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "10"}, Value: 10},
			},
		},
	}

	expected := "{510}"
	if blockStmt.String() != expected {
		t.Errorf("blockStmt.String() wrong. got=%q, want=%q", blockStmt.String(), expected)
	}

	if blockStmt.TokenLiteral() != "{" {
		t.Errorf("blockStmt.TokenLiteral() wrong. got=%q, want=%q", blockStmt.TokenLiteral(), "{")
	}
}

func TestFunctionLiteralString(t *testing.T) {
	fnLit := &FunctionLiteral{
		Token: lexer.Token{Type: lexer.FN, Literal: "fn"},
		Parameters: []*Identifier{
			{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
			{Token: lexer.Token{Type: lexer.IDENT, Literal: "y"}, Value: "y"},
		},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token: lexer.Token{Type: lexer.IDENT, Literal: "x"},
					Expression: &InfixExpression{
						Token:    lexer.Token{Type: lexer.PLUS, Literal: "+"},
						Left:     &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
						Operator: "+",
						Right:    &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "y"}, Value: "y"},
					},
				},
			},
		},
	}

	expected := "fn(x, y) {(x + y)}"
	if fnLit.String() != expected {
		t.Errorf("fnLit.String() wrong. got=%q, want=%q", fnLit.String(), expected)
	}

	if fnLit.TokenLiteral() != "fn" {
		t.Errorf("fnLit.TokenLiteral() wrong. got=%q, want=%q", fnLit.TokenLiteral(), "fn")
	}
}

func TestCallExpressionString(t *testing.T) {
	callExpr := &CallExpression{
		Token:    lexer.Token{Type: lexer.LPAREN, Literal: "("},
		Function: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "add"}, Value: "add"},
		Arguments: []Expression{
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "2"}, Value: 2},
		},
	}

	expected := "add(1, 2)"
	if callExpr.String() != expected {
		t.Errorf("callExpr.String() wrong. got=%q, want=%q", callExpr.String(), expected)
	}

	if callExpr.TokenLiteral() != "(" {
		t.Errorf("callExpr.TokenLiteral() wrong. got=%q, want=%q", callExpr.TokenLiteral(), "(")
	}
}

func TestReturnStatementString(t *testing.T) {
	returnStmt := &ReturnStatement{
		Token:       lexer.Token{Type: lexer.RETURN, Literal: "return"},
		ReturnValue: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "42"}, Value: 42},
	}

	expected := "return 42;"
	if returnStmt.String() != expected {
		t.Errorf("returnStmt.String() wrong. got=%q, want=%q", returnStmt.String(), expected)
	}

	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral() wrong. got=%q, want=%q", returnStmt.TokenLiteral(), "return")
	}
}

func TestIndexExpressionString(t *testing.T) {
	indexExpr := &IndexExpression{
		Token: lexer.Token{Type: lexer.LBRACKET, Literal: "["},
		Left:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "myArray"}, Value: "myArray"},
		Index: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
	}

	expected := "(myArray[1])"
	if indexExpr.String() != expected {
		t.Errorf("indexExpr.String() wrong. got=%q, want=%q", indexExpr.String(), expected)
	}

	if indexExpr.TokenLiteral() != "[" {
		t.Errorf("indexExpr.TokenLiteral() wrong. got=%q, want=%q", indexExpr.TokenLiteral(), "[")
	}
}

func TestWhileStatementString(t *testing.T) {
	whileStmt := &WhileStatement{
		Token:     lexer.Token{Type: lexer.WHILE, Literal: "while"},
		Condition: &BooleanLiteral{Token: lexer.Token{Type: lexer.TRUE, Literal: "true"}, Value: true},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "1"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
				},
			},
		},
	}

	expected := "whiletrue {1}"
	if whileStmt.String() != expected {
		t.Errorf("whileStmt.String() wrong. got=%q, want=%q", whileStmt.String(), expected)
	}

	if whileStmt.TokenLiteral() != "while" {
		t.Errorf("whileStmt.TokenLiteral() wrong. got=%q, want=%q", whileStmt.TokenLiteral(), "while")
	}
}

func TestForStatementString(t *testing.T) {
	forStmt := &ForStatement{
		Token: lexer.Token{Type: lexer.FOR, Literal: "for"},
		Init: &AssignmentStatement{
			Token: lexer.Token{Type: lexer.IDENT, Literal: "i"},
			Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "i"}, Value: "i"},
			Value: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "0"}, Value: 0},
		},
		Condition: &InfixExpression{
			Token:    lexer.Token{Type: lexer.LT, Literal: "<"},
			Left:     &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "i"}, Value: "i"},
			Operator: "<",
			Right:    &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "10"}, Value: 10},
		},
		Update: &AssignmentStatement{
			Token: lexer.Token{Type: lexer.IDENT, Literal: "i"},
			Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "i"}, Value: "i"},
			Value: &InfixExpression{
				Token:    lexer.Token{Type: lexer.PLUS, Literal: "+"},
				Left:     &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "i"}, Value: "i"},
				Operator: "+",
				Right:    &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
			},
		},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "1"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
				},
			},
		},
	}

	expected := "for(i = 0;(i < 10);i = (i + 1)) {1}"
	if forStmt.String() != expected {
		t.Errorf("forStmt.String() wrong. got=%q, want=%q", forStmt.String(), expected)
	}

	if forStmt.TokenLiteral() != "for" {
		t.Errorf("forStmt.TokenLiteral() wrong. got=%q, want=%q", forStmt.TokenLiteral(), "for")
	}
}

func TestImportStatementString(t *testing.T) {
	importStmt := &ImportStatement{
		Token: lexer.Token{Type: lexer.IMPORT, Literal: "import"},
		Items: []*ImportItem{
			{
				Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "foo"}, Value: "foo"},
				Alias: nil,
			},
			{
				Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "bar"}, Value: "bar"},
				Alias: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "baz"}, Value: "baz"},
			},
		},
		Module: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "module"}, Value: "module"},
	}

	expected := "import { foo, bar as baz } from \"module\""
	if importStmt.String() != expected {
		t.Errorf("importStmt.String() wrong. got=%q, want=%q", importStmt.String(), expected)
	}

	if importStmt.TokenLiteral() != "import" {
		t.Errorf("importStmt.TokenLiteral() wrong. got=%q, want=%q", importStmt.TokenLiteral(), "import")
	}
}

func TestExportStatementString(t *testing.T) {
	exportStmt := &ExportStatement{
		Token: lexer.Token{Type: lexer.EXPORT, Literal: "export"},
		Name:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "myVar"}, Value: "myVar"},
		Value: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "42"}, Value: 42},
	}

	expected := "export myVar = 42"
	if exportStmt.String() != expected {
		t.Errorf("exportStmt.String() wrong. got=%q, want=%q", exportStmt.String(), expected)
	}

	if exportStmt.TokenLiteral() != "export" {
		t.Errorf("exportStmt.TokenLiteral() wrong. got=%q, want=%q", exportStmt.TokenLiteral(), "export")
	}
}

func TestPropertyAccessString(t *testing.T) {
	propAccess := &PropertyAccess{
		Token:    lexer.Token{Type: lexer.DOT, Literal: "."},
		Object:   &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "obj"}, Value: "obj"},
		Property: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "prop"}, Value: "prop"},
	}

	expected := "(obj.prop)"
	if propAccess.String() != expected {
		t.Errorf("propAccess.String() wrong. got=%q, want=%q", propAccess.String(), expected)
	}

	if propAccess.TokenLiteral() != "." {
		t.Errorf("propAccess.TokenLiteral() wrong. got=%q, want=%q", propAccess.TokenLiteral(), ".")
	}
}

func TestThrowStatementString(t *testing.T) {
	throwStmt := &ThrowStatement{
		Token: lexer.Token{Type: lexer.THROW, Literal: "throw"},
		Expression: &CallExpression{
			Token:    lexer.Token{Type: lexer.LPAREN, Literal: "("},
			Function: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "Error"}, Value: "Error"},
			Arguments: []Expression{
				&StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "message"}, Value: "message"},
			},
		},
	}

	expected := "throw Error(\"message\")"
	if throwStmt.String() != expected {
		t.Errorf("throwStmt.String() wrong. got=%q, want=%q", throwStmt.String(), expected)
	}

	if throwStmt.TokenLiteral() != "throw" {
		t.Errorf("throwStmt.TokenLiteral() wrong. got=%q, want=%q", throwStmt.TokenLiteral(), "throw")
	}
}

func TestCatchClauseString(t *testing.T) {
	catchClause := &CatchClause{
		Token:     lexer.Token{Type: lexer.CATCH, Literal: "catch"},
		ErrorType: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "Error"}, Value: "Error"},
		ErrorVar:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "err"}, Value: "err"},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "1"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
				},
			},
		},
	}

	expected := "catch (Error err) {1}"
	if catchClause.String() != expected {
		t.Errorf("catchClause.String() wrong. got=%q, want=%q", catchClause.String(), expected)
	}

	if catchClause.TokenLiteral() != "catch" {
		t.Errorf("catchClause.TokenLiteral() wrong. got=%q, want=%q", catchClause.TokenLiteral(), "catch")
	}
}

func TestTryStatementString(t *testing.T) {
	tryStmt := &TryStatement{
		Token: lexer.Token{Type: lexer.TRY, Literal: "try"},
		TryBlock: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "1"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
				},
			},
		},
		CatchClauses: []*CatchClause{
			{
				Token:     lexer.Token{Type: lexer.CATCH, Literal: "catch"},
				ErrorType: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "Error"}, Value: "Error"},
				ErrorVar:  &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "err"}, Value: "err"},
				Body: &BlockStatement{
					Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token:      lexer.Token{Type: lexer.INT, Literal: "2"},
							Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "2"}, Value: 2},
						},
					},
				},
			},
		},
		FinallyBlock: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.INT, Literal: "3"},
					Expression: &IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "3"}, Value: 3},
				},
			},
		},
	}

	expected := "try {1} catch (Error err) {2} finally {3}"
	if tryStmt.String() != expected {
		t.Errorf("tryStmt.String() wrong. got=%q, want=%q", tryStmt.String(), expected)
	}

	if tryStmt.TokenLiteral() != "try" {
		t.Errorf("tryStmt.TokenLiteral() wrong. got=%q, want=%q", tryStmt.TokenLiteral(), "try")
	}
}

func TestSwitchStatementString(t *testing.T) {
	switchStmt := &SwitchStatement{
		Token: lexer.Token{Type: lexer.SWITCH, Literal: "switch"},
		Value: &Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
		Cases: []*CaseClause{
			{
				Token: lexer.Token{Type: lexer.CASE, Literal: "case"},
				Values: []Expression{
					&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
				},
				Body: &BlockStatement{
					Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
					Statements: []Statement{
						&ExpressionStatement{
							Token:      lexer.Token{Type: lexer.STRING, Literal: "one"},
							Expression: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "one"}, Value: "one"},
						},
					},
				},
			},
		},
		Default: &DefaultClause{
			Token: lexer.Token{Type: lexer.DEFAULT, Literal: "default"},
			Body: &BlockStatement{
				Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
				Statements: []Statement{
					&ExpressionStatement{
						Token:      lexer.Token{Type: lexer.STRING, Literal: "other"},
						Expression: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "other"}, Value: "other"},
					},
				},
			},
		},
	}

	expected := "switch (x) {case 1:{\"one\"}default:{\"other\"}}"
	if switchStmt.String() != expected {
		t.Errorf("switchStmt.String() wrong. got=%q, want=%q", switchStmt.String(), expected)
	}

	if switchStmt.TokenLiteral() != "switch" {
		t.Errorf("switchStmt.TokenLiteral() wrong. got=%q, want=%q", switchStmt.TokenLiteral(), "switch")
	}
}

func TestCaseClauseString(t *testing.T) {
	caseClause := &CaseClause{
		Token: lexer.Token{Type: lexer.CASE, Literal: "case"},
		Values: []Expression{
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
			&IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "2"}, Value: 2},
		},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.STRING, Literal: "low"},
					Expression: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "low"}, Value: "low"},
				},
			},
		},
	}

	expected := "case 1, 2:{\"low\"}"
	if caseClause.String() != expected {
		t.Errorf("caseClause.String() wrong. got=%q, want=%q", caseClause.String(), expected)
	}

	if caseClause.TokenLiteral() != "case" {
		t.Errorf("caseClause.TokenLiteral() wrong. got=%q, want=%q", caseClause.TokenLiteral(), "case")
	}
}

func TestDefaultClauseString(t *testing.T) {
	defaultClause := &DefaultClause{
		Token: lexer.Token{Type: lexer.DEFAULT, Literal: "default"},
		Body: &BlockStatement{
			Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
			Statements: []Statement{
				&ExpressionStatement{
					Token:      lexer.Token{Type: lexer.STRING, Literal: "other"},
					Expression: &StringLiteral{Token: lexer.Token{Type: lexer.STRING, Literal: "other"}, Value: "other"},
				},
			},
		},
	}

	expected := "default:{\"other\"}"
	if defaultClause.String() != expected {
		t.Errorf("defaultClause.String() wrong. got=%q, want=%q", defaultClause.String(), expected)
	}

	if defaultClause.TokenLiteral() != "default" {
		t.Errorf("defaultClause.TokenLiteral() wrong. got=%q, want=%q", defaultClause.TokenLiteral(), "default")
	}
}

func TestEmptyProgramString(t *testing.T) {
	program := &Program{Statements: []Statement{}}
	
	if program.String() != "" {
		t.Errorf("empty program.String() should be empty. got=%q", program.String())
	}
	
	if program.TokenLiteral() != "" {
		t.Errorf("empty program.TokenLiteral() should be empty. got=%q", program.TokenLiteral())
	}
}

func TestNilExpressionStatementString(t *testing.T) {
	exprStmt := &ExpressionStatement{
		Token:      lexer.Token{Type: lexer.SEMICOLON, Literal: ";"},
		Expression: nil,
	}

	if exprStmt.String() != "" {
		t.Errorf("nil expression statement.String() should be empty. got=%q", exprStmt.String())
	}
}