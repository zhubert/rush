package lexer

import "testing"

func TestBasicTokens(t *testing.T) {
  input := `let add = fn(x, y) {
  x + y
}`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {IDENT, "let"},
    {IDENT, "add"},
    {ASSIGN, "="},
    {FN, "fn"},
    {LPAREN, "("},
    {IDENT, "x"},
    {COMMA, ","},
    {IDENT, "y"},
    {RPAREN, ")"},
    {LBRACE, "{"},
    {SEMICOLON, "\n"},
    {IDENT, "x"},
    {PLUS, "+"},
    {IDENT, "y"},
    {SEMICOLON, "\n"},
    {RBRACE, "}"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal: %q)",
        i, tt.expectedType, tok.Type, tok.Literal)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}

func TestSwitchCaseTokens(t *testing.T) {
	input := `switch (grade) {
    case "A":
        result = "excellent"
    case "B", "C":
        result = "good"
    default:
        result = "other"
}`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{SWITCH, "switch"},
		{LPAREN, "("},
		{IDENT, "grade"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{SEMICOLON, "\n"},
		{CASE, "case"},
		{STRING, "A"},
		{COLON, ":"},
		{SEMICOLON, "\n"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{STRING, "excellent"},
		{SEMICOLON, "\n"},
		{CASE, "case"},
		{STRING, "B"},
		{COMMA, ","},
		{STRING, "C"},
		{COLON, ":"},
		{SEMICOLON, "\n"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{STRING, "good"},
		{SEMICOLON, "\n"},
		{DEFAULT, "default"},
		{COLON, ":"},
		{SEMICOLON, "\n"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{STRING, "other"},
		{SEMICOLON, "\n"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal: %q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestColonToken(t *testing.T) {
	input := `:`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{COLON, ":"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
  input := `+ - * / < > == != <= >= && || !`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {PLUS, "+"},
    {MINUS, "-"},
    {MULT, "*"},
    {DIV, "/"},
    {LT, "<"},
    {GT, ">"},
    {EQ, "=="},
    {NOT_EQ, "!="},
    {LTE, "<="},
    {GTE, ">="},
    {AND, "&&"},
    {OR, "||"},
    {NOT, "!"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, tok.Type)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}

func TestKeywords(t *testing.T) {
  input := `if else while for return true false import export from as`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {IF, "if"},
    {ELSE, "else"},
    {WHILE, "while"},
    {FOR, "for"},
    {RETURN, "return"},
    {TRUE, "true"},
    {FALSE, "false"},
    {IMPORT, "import"},
    {EXPORT, "export"},
    {FROM, "from"},
    {AS, "as"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, tok.Type)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}

func TestStringsAndNumbers(t *testing.T) {
  input := `"hello world" 42 3.14`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {STRING, "hello world"},
    {INT, "42"},
    {FLOAT, "3.14"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, tok.Type)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}

func TestArrayAndIndexing(t *testing.T) {
  input := `[1, 2, 3] arr[0]`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {LBRACKET, "["},
    {INT, "1"},
    {COMMA, ","},
    {INT, "2"},
    {COMMA, ","},
    {INT, "3"},
    {RBRACKET, "]"},
    {IDENT, "arr"},
    {LBRACKET, "["},
    {INT, "0"},
    {RBRACKET, "]"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, tok.Type)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}

func TestComments(t *testing.T) {
  input := `# This is a comment
x = 5`

  tests := []struct {
    expectedType    TokenType
    expectedLiteral string
  }{
    {COMMENT, "# This is a comment"},
    {SEMICOLON, "\n"},
    {IDENT, "x"},
    {ASSIGN, "="},
    {INT, "5"},
    {EOF, ""},
  }

  l := New(input)

  for i, tt := range tests {
    tok := l.NextToken()

    if tok.Type != tt.expectedType {
      t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, tok.Type)
    }

    if tok.Literal != tt.expectedLiteral {
      t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedLiteral, tok.Literal)
    }
  }
}