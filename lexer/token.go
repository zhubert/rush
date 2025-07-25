package lexer

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Identifiers and literals
	IDENT  // variable names, function names
	INT    // 42
	FLOAT  // 3.14
	STRING // "foo"
	TRUE   // true
	FALSE  // false

	// Operators
	ASSIGN // =
	PLUS   // +
	MINUS  // -
	MULT   // *
	DIV    // /
	MOD    // %
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=
	AND    // &&
	OR     // ||
	NOT    // !

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	DOT       // .

	// Keywords
	FN     // fn
	IF     // if
	ELSE   // else
	FOR    // for
	WHILE  // while
	RETURN // return
	IMPORT  // import
	EXPORT  // export
	FROM    // from
	TRY     // try
	CATCH   // catch
	FINALLY // finally
	THROW   // throw
	CLASS   // class
	INITIALIZE // initialize
	SUPER   // super
	INSTANCE_VAR // @
	BREAK   // break
	CONTINUE // continue
	SWITCH   // switch
	CASE     // case
	DEFAULT  // default
	AS       // as
)

// Token represents a single token
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Column   int
}

// tokenTypeNames maps token types to their string representations
var tokenTypeNames = map[TokenType]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	COMMENT:   "COMMENT",
	IDENT:     "IDENT",
	INT:       "INT",
	FLOAT:     "FLOAT",
	STRING:    "STRING",
	TRUE:      "TRUE",
	FALSE:     "FALSE",
	ASSIGN:    "=",
	PLUS:      "+",
	MINUS:     "-",
	MULT:      "*",
	DIV:       "/",
	MOD:       "%",
	EQ:        "==",
	NOT_EQ:    "!=",
	LT:        "<",
	GT:        ">",
	LTE:       "<=",
	GTE:       ">=",
	AND:       "&&",
	OR:        "||",
	NOT:       "!",
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	DOT:       ".",
	FN:        "fn",
	IF:        "if",
	ELSE:      "else",
	FOR:       "for",
	WHILE:     "while",
	RETURN:    "return",
	IMPORT:    "import",
	EXPORT:    "export",
	FROM:      "from",
	TRY:       "try",
	CATCH:     "catch",
	FINALLY:   "finally",
	THROW:     "throw",
	CLASS:     "class",
	INITIALIZE: "initialize",
	SUPER:     "super",
	INSTANCE_VAR: "@",
	BREAK:     "break",
	CONTINUE:  "continue",
	SWITCH:    "switch",
	CASE:      "case",
	DEFAULT:   "default",
	AS:        "as",
}

// String returns the string representation of a token type
func (t TokenType) String() string {
	if name, ok := tokenTypeNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// keywords maps string literals to their token types
var keywords = map[string]TokenType{
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"return": RETURN,
	"import":  IMPORT,
	"export":  EXPORT,
	"from":    FROM,
	"try":     TRY,
	"catch":   CATCH,
	"finally": FINALLY,
	"throw":   THROW,
	"class":   CLASS,
	"initialize": INITIALIZE,
	"super":   SUPER,
	"break":   BREAK,
	"continue": CONTINUE,
	"switch":  SWITCH,
	"case":    CASE,
	"default": DEFAULT,
	"as":      AS,
	"true":    TRUE,
	"false":   FALSE,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}