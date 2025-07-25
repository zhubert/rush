package lexer


// Lexer tokenizes input source code
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokenType := INT
	
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	
	return l.input[position:l.position], tokenType
}

// readString reads a string literal with the given quote character
func (l *Lexer) readString(quote byte) string {
	var result []byte
	l.readChar() // skip opening quote
	
	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() != 0 {
			// Handle escape sequences
			l.readChar() // consume backslash
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			case '\'':
				result = append(result, '\'')
			default:
				// For unknown escape sequences, include both backslash and character
				result = append(result, '\\')
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}
	
	return string(result)
}

// readComment reads a comment starting with # or //
func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// NextToken returns the next token in the input
func (l *Lexer) NextToken() Token {
	var tok Token
	
	l.skipWhitespace()
	
	// Store current position for token
	line := l.line
	column := l.column
	
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ASSIGN, l.ch, line, column)
		}
	case '+':
		tok = newToken(PLUS, l.ch, line, column)
	case '-':
		tok = newToken(MINUS, l.ch, line, column)
	case '*':
		tok = newToken(MULT, l.ch, line, column)
	case '/':
		if l.peekChar() == '/' {
			// Handle // comment
			tok.Type = COMMENT
			tok.Literal = l.readComment()
			tok.Line = line
			tok.Column = column
			return tok // Don't advance past newline
		} else {
			tok = newToken(DIV, l.ch, line, column)
		}
	case '%':
		tok = newToken(MOD, l.ch, line, column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(NOT, l.ch, line, column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(LT, l.ch, line, column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(GT, l.ch, line, column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	case ',':
		tok = newToken(COMMA, l.ch, line, column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, line, column)
	case ':':
		tok = newToken(COLON, l.ch, line, column)
	case '(':
		tok = newToken(LPAREN, l.ch, line, column)
	case ')':
		tok = newToken(RPAREN, l.ch, line, column)
	case '{':
		tok = newToken(LBRACE, l.ch, line, column)
	case '}':
		tok = newToken(RBRACE, l.ch, line, column)
	case '[':
		tok = newToken(LBRACKET, l.ch, line, column)
	case ']':
		tok = newToken(RBRACKET, l.ch, line, column)
	case '.':
		// Only treat as DOT if not followed by a digit (which would be a float)
		if !isDigit(l.peekChar()) {
			tok = newToken(DOT, l.ch, line, column)
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
		tok.Line = line
		tok.Column = column
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
		tok.Line = line
		tok.Column = column
	case '@':
		tok = newToken(INSTANCE_VAR, l.ch, line, column)
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
		tok.Line = line
		tok.Column = column
		return tok // Don't advance past newline
	case '\n':
		tok = newToken(SEMICOLON, l.ch, line, column) // Treat newlines as statement separators
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = line
		tok.Column = column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = column
			return tok // readIdentifier already advanced position
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = line
			tok.Column = column
			return tok // readNumber already advanced position
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	}
	
	l.readChar()
	return tok
}

// newToken creates a new token with the given type and character
func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
		Column:  column,
	}
}

// isLetter checks if a character is a letter, underscore, or question mark (for boolean predicates)
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}