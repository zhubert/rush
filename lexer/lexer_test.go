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

// Edge case and error handling tests

func TestEmptyInput(t *testing.T) {
  input := ""
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != EOF {
    t.Fatalf("Expected EOF, got %q", tok.Type)
  }
  
  if tok.Literal != "" {
    t.Fatalf("EOF literal should be empty, got %q", tok.Literal)
  }
}

func TestWhitespaceOnly(t *testing.T) {
  input := "   \t  \n  \r\n  "
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != SEMICOLON {
    t.Fatalf("Expected SEMICOLON for first newline, got %q", tok.Type)
  }
  
  tok = l.NextToken()
  if tok.Type != SEMICOLON {
    t.Fatalf("Expected SEMICOLON for second newline, got %q", tok.Type)
  }
  
  tok = l.NextToken()
  if tok.Type != EOF {
    t.Fatalf("Expected EOF at end, got %q", tok.Type)
  }
}

func TestUnterminatedString(t *testing.T) {
  input := `"unterminated string`
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != STRING {
    t.Fatalf("Expected STRING for unterminated string, got %q", tok.Type)
  }
  
  // The lexer reads until EOF for unterminated strings
  if tok.Literal != "unterminated string" {
    t.Fatalf("String literal wrong. expected=%q, got=%q", "unterminated string", tok.Literal)
  }
}

func TestStringWithEscapes(t *testing.T) {
  input := `"hello\nworld\t\"quoted\""`
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != STRING {
    t.Fatalf("Expected STRING, got %q", tok.Type)
  }
  
  // The lexer now processes escape sequences properly
  expected := "hello\nworld\t\"quoted\""
  if tok.Literal != expected {
    t.Fatalf("String literal wrong. expected=%q, got=%q", expected, tok.Literal)
  }
}

func TestSingleQuoteStrings(t *testing.T) {
  input := `'hello world'`
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != STRING {
    t.Fatalf("Expected STRING, got %q", tok.Type)
  }
  
  expected := "hello world"
  if tok.Literal != expected {
    t.Fatalf("String literal wrong. expected=%q, got=%q", expected, tok.Literal)
  }
}

func TestSingleQuoteStringWithEscapes(t *testing.T) {
  input := `'hello\nworld\'quoted\''`
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != STRING {
    t.Fatalf("Expected STRING, got %q", tok.Type)
  }
  
  expected := "hello\nworld'quoted'"
  if tok.Literal != expected {
    t.Fatalf("String literal wrong. expected=%q, got=%q", expected, tok.Literal)
  }
}

func TestFloatNumbers(t *testing.T) {
  tests := []struct {
    input    string
    expected string
    expectedType TokenType
  }{
    {"3.14", "3.14", FLOAT},
    {"0.5", "0.5", FLOAT},
    {"123.456", "123.456", FLOAT},
    {"5.", "5", INT}, // The lexer only sees '5', '.' would be separate token
  }

  for _, tt := range tests {
    l := New(tt.input)
    tok := l.NextToken()
    
    if tok.Type != tt.expectedType {
      t.Fatalf("Expected %q for %q, got %q", tt.expectedType, tt.input, tok.Type)
    }
    
    if tok.Literal != tt.expected {
      t.Fatalf("Literal wrong for %q. expected=%q, got=%q", 
        tt.input, tt.expected, tok.Literal)
    }
  }
}

func TestInvalidFloats(t *testing.T) {
  tests := []string{
    "3.14.15",  // Multiple dots
    "3..14",    // Double dots
  }

  for _, input := range tests {
    l := New(input)
    tok := l.NextToken()
    
    // Should tokenize as separate tokens, not ILLEGAL
    if tok.Type == ILLEGAL {
      t.Fatalf("Should not be ILLEGAL for %q, got %q", input, tok.Type)
    }
  }
}

func TestLongNumbers(t *testing.T) {
  input := "123456789012345678901234567890"
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != INT {
    t.Fatalf("Expected INT for long number, got %q", tok.Type)
  }
  
  if tok.Literal != input {
    t.Fatalf("Long number literal wrong. expected=%q, got=%q", input, tok.Literal)
  }
}

func TestSingleCharacterTokens(t *testing.T) {
  input := "(){}[],.;:@"
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {LPAREN, "("},
    {RPAREN, ")"},
    {LBRACE, "{"},
    {RBRACE, "}"},
    {LBRACKET, "["},
    {RBRACKET, "]"},
    {COMMA, ","},
    {DOT, "."},
    {SEMICOLON, ";"},
    {COLON, ":"},
    {INSTANCE_VAR, "@"},
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

func TestTwoCharacterOperators(t *testing.T) {
  input := "== != <= >= && ||"
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {EQ, "=="},
    {NOT_EQ, "!="},
    {LTE, "<="},
    {GTE, ">="},
    {AND, "&&"},
    {OR, "||"},
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

func TestMixedOperators(t *testing.T) {
  input := "= == ! != < <= > >="
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {ASSIGN, "="},
    {EQ, "=="},
    {NOT, "!"},
    {NOT_EQ, "!="},
    {LT, "<"},
    {LTE, "<="},
    {GT, ">"},
    {GTE, ">="},
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

func TestAllKeywords(t *testing.T) {
  input := `fn if else for while return import export from try catch finally throw 
            class initialize super break continue switch case default as true false`

  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {FN, "fn"},
    {IF, "if"},
    {ELSE, "else"},
    {FOR, "for"},
    {WHILE, "while"},
    {RETURN, "return"},
    {IMPORT, "import"},
    {EXPORT, "export"},
    {FROM, "from"},
    {TRY, "try"},
    {CATCH, "catch"},
    {FINALLY, "finally"},
    {THROW, "throw"},
    {SEMICOLON, "\n"},
    {CLASS, "class"},
    {INITIALIZE, "initialize"},
    {SUPER, "super"},
    {BREAK, "break"},
    {CONTINUE, "continue"},
    {SWITCH, "switch"},
    {CASE, "case"},
    {DEFAULT, "default"},
    {AS, "as"},
    {TRUE, "true"},
    {FALSE, "false"},
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

func TestIdentifierEdgeCases(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"_underscore", "_underscore"},
    {"var123", "var123"},
    {"CamelCase", "CamelCase"},
    {"snake_case", "snake_case"},
    {"a", "a"},
    {"A", "A"},
    {"myVeryLongIdentifierName", "myVeryLongIdentifierName"},
  }

  for _, tt := range tests {
    l := New(tt.input)
    tok := l.NextToken()
    
    if tok.Type != IDENT {
      t.Fatalf("Expected IDENT for %q, got %q", tt.input, tok.Type)
    }
    
    if tok.Literal != tt.expected {
      t.Fatalf("Identifier literal wrong for %q. expected=%q, got=%q", 
        tt.input, tt.expected, tok.Literal)
    }
  }
}

func TestLineAndColumnTracking(t *testing.T) {
  input := `x = 5
y = 10`

  l := New(input)
  
  // x token
  tok := l.NextToken()
  if tok.Line != 1 || tok.Column != 1 {
    t.Fatalf("Expected line=1, column=1 for 'x', got line=%d, column=%d", 
      tok.Line, tok.Column)
  }
  
  // = token  
  tok = l.NextToken()
  if tok.Line != 1 || tok.Column != 3 {
    t.Fatalf("Expected line=1, column=3 for '=', got line=%d, column=%d", 
      tok.Line, tok.Column)
  }
  
  // 5 token
  tok = l.NextToken()
  if tok.Line != 1 || tok.Column != 5 {
    t.Fatalf("Expected line=1, column=5 for '5', got line=%d, column=%d", 
      tok.Line, tok.Column)
  }
  
  // newline token - the lexer column tracking updates after reading the newline
  tok = l.NextToken()
  if tok.Line != 2 || tok.Column != 0 {
    t.Fatalf("Expected line=2, column=0 for newline, got line=%d, column=%d", 
      tok.Line, tok.Column)
  }
  
  // y token (on new line)
  tok = l.NextToken()
  if tok.Line != 2 || tok.Column != 1 {
    t.Fatalf("Expected line=2, column=1 for 'y', got line=%d, column=%d", 
      tok.Line, tok.Column)
  }
}

func TestMultilineComments(t *testing.T) {
  input := `# First comment
x = 5  # End of line comment
# Another comment`

  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {COMMENT, "# First comment"},
    {SEMICOLON, "\n"},
    {IDENT, "x"},
    {ASSIGN, "="},
    {INT, "5"},
    {COMMENT, "# End of line comment"},
    {SEMICOLON, "\n"},
    {COMMENT, "# Another comment"},
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

func TestIllegalCharacters(t *testing.T) {
  tests := []string{
    "$", "^", "&", "|", "~", "`", "\\", 
  }

  for _, input := range tests {
    l := New(input)
    tok := l.NextToken()
    
    if tok.Type != ILLEGAL {
      t.Fatalf("Expected ILLEGAL for %q, got %q", input, tok.Type)
    }
    
    if tok.Literal != input {
      t.Fatalf("ILLEGAL literal wrong for %q. expected=%q, got=%q", 
        input, input, tok.Literal)
    }
  }
}

func TestModuloOperator(t *testing.T) {
  input := "5 % 2"
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {INT, "5"},
    {MOD, "%"},
    {INT, "2"},
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

func TestComplexExpression(t *testing.T) {
  input := `arr[index].prop && (x >= 10 || y != "test")`
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {IDENT, "arr"},
    {LBRACKET, "["},
    {IDENT, "index"},
    {RBRACKET, "]"},
    {DOT, "."},
    {IDENT, "prop"},
    {AND, "&&"},
    {LPAREN, "("},
    {IDENT, "x"},
    {GTE, ">="},
    {INT, "10"},
    {OR, "||"},
    {IDENT, "y"},
    {NOT_EQ, "!="},
    {STRING, "test"},
    {RPAREN, ")"},
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

func TestZeroValue(t *testing.T) {
  input := "0"
  l := New(input)
  
  tok := l.NextToken()
  if tok.Type != INT {
    t.Fatalf("Expected INT for zero, got %q", tok.Type)
  }
  
  if tok.Literal != "0" {
    t.Fatalf("Zero literal wrong. expected=%q, got=%q", "0", tok.Literal)
  }
}

func TestNegativeNumbers(t *testing.T) {
  // Note: negative numbers are parsed as MINUS followed by number
  input := "+42 -17"
  
  tests := []struct {
    expectedType TokenType
    expectedLiteral string
  }{
    {PLUS, "+"},
    {INT, "42"},
    {MINUS, "-"},
    {INT, "17"},
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