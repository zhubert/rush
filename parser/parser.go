package parser

import (
	"fmt"
	"strconv"

	"rush/ast"
	"rush/lexer"
)

// Precedence levels for operator precedence parsing
const (
	_ int = iota
	LOWEST
	LOGICAL     // && and ||
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

// precedences maps token types to their precedence
var precedences = map[lexer.TokenType]int{
	lexer.EQ:      EQUALS,
	lexer.NOT_EQ:  EQUALS,
	lexer.LT:      LESSGREATER,
	lexer.GT:      LESSGREATER,
	lexer.LTE:     LESSGREATER,
	lexer.GTE:     LESSGREATER,
	lexer.PLUS:    SUM,
	lexer.MINUS:   SUM,
	lexer.DIV:     PRODUCT,
	lexer.MULT:    PRODUCT,
	lexer.MOD:     PRODUCT,
	lexer.AND:     LOGICAL,
	lexer.OR:      LOGICAL,
	lexer.LPAREN:  CALL,
	lexer.LBRACKET: INDEX,
	lexer.DOT:     INDEX, // module.member has same precedence as array[index]
}

// Parser parses tokens into an AST
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseHashLiteral)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.FN, p.parseFunctionLiteral)
	p.registerPrefix(lexer.INSTANCE_VAR, p.parseInstanceVariable)
	p.registerPrefix(lexer.SUPER, p.parseSuperExpression)

	// Initialize infix parse functions
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.DIV, p.parseInfixExpression)
	p.registerInfix(lexer.MULT, p.parseInfixExpression)
	p.registerInfix(lexer.MOD, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parsePropertyAccess)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken advances the parser tokens
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	
	// Skip comments in current token
	for p.curToken.Type == lexer.COMMENT {
		p.curToken = p.peekToken
		p.peekToken = p.l.NextToken()
	}
	
	// Skip comments in peek token
	for p.peekToken.Type == lexer.COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF {
		// Handle method chaining continuation - check if next line starts with DOT
		if p.curToken.Type == lexer.SEMICOLON {
			// Look ahead to see if next non-whitespace token is DOT
			tempLexer := *p.l
			nextToken := tempLexer.NextToken()
			
			// Skip the semicolon and continue parsing if next token is DOT
			if nextToken.Type == lexer.DOT {
				p.nextToken() // Skip semicolon
				continue      // Continue parsing the expression
			} else {
				// Normal semicolon - skip it
				p.nextToken()
				continue
			}
		}
		
		// Skip comments
		if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
			continue
		}

		// Parse normal statement
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses statements
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.IMPORT:
		return p.parseImportStatement()
	case lexer.EXPORT:
		return p.parseExportStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	case lexer.SWITCH:
		return p.parseSwitchStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.TRY:
		return p.parseTryStatement()
	case lexer.THROW:
		return p.parseThrowStatement()
	case lexer.CLASS:
		return p.parseClassDeclaration()
	case lexer.INSTANCE_VAR:
		return p.parseInstanceVariableStatement()
	default:
		// Check if this is an assignment statement (identifier = value)
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			return p.parseAssignmentStatement()
		}
		// Check if this is an array element assignment (identifier[index] = value)
		if p.isIndexAssignment() {
			return p.parseIndexAssignmentStatement()
		}
		// Otherwise, parse as expression statement
		return p.parseExpressionStatement()
	}
}

// parseAssignmentStatement parses assignment statements like "a = 5"
func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}

	if p.curToken.Type != lexer.IDENT {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

// isIndexAssignment checks if the current position represents an array index assignment
// Pattern: IDENT [ ... ] = 
func (p *Parser) isIndexAssignment() bool {
	// Simple check: if current token is IDENT and peek is LBRACKET, it might be array assignment
	// We'll do the full validation in parseIndexAssignmentStatement
	return p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.LBRACKET
}

// parseIndexAssignmentStatement parses array element assignments like "arr[0] = 5"
func (p *Parser) parseIndexAssignmentStatement() ast.Statement {
	// First, parse the left side as an index expression
	leftExpr := p.parseExpression(LOWEST)
	indexExpr, ok := leftExpr.(*ast.IndexExpression)
	if !ok {
		// This wasn't an index expression, fall back to expression statement
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: leftExpr,
		}
	}
	
	// Check if the next token is an assignment
	if p.peekToken.Type != lexer.ASSIGN {
		// This is just an index expression, not an assignment
		return &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: leftExpr,
		}
	}
	
	// This is indeed an index assignment
	p.nextToken() // move to =
	
	stmt := &ast.IndexAssignmentStatement{
		Token: p.curToken, // the = token
		Left:  indexExpr,
	}
	
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	
	// Optional semicolon
	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}
	
	return stmt
}

// parseInstanceVariableStatement parses instance variable statements like "@name" or "@name = value"
func (p *Parser) parseInstanceVariableStatement() ast.Statement {
	// Parse the @name part first
	instVarExpr := p.parseInstanceVariable()
	if instVarExpr == nil {
		return nil
	}

	// Cast to InstanceVariable type
	instVar, ok := instVarExpr.(*ast.InstanceVariable)
	if !ok {
		return nil
	}

	// Check if this is an assignment (@name = value)
	if p.peekToken.Type == lexer.ASSIGN {
		// Convert to assignment statement
		p.nextToken() // consume =
		p.nextToken() // move to value
		
		stmt := &ast.AssignmentStatement{
			Token: instVar.Token,
			Name: &ast.Identifier{Token: instVar.Token, Value: "@" + instVar.Name.Value},
			Value: p.parseExpression(LOWEST),
		}
		
		// Optional semicolon or newline
		if p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
		
		return stmt
	} else {
		// Just an expression statement (@name)
		return &ast.ExpressionStatement{
			Token: instVar.Token,
			Expression: instVar,
		}
	}
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

// parseExpression parses expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Skip any leading semicolons (newlines)
	for p.curToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.shouldStopParsing(precedence) {
		// Skip semicolons before looking for infix operators
		for p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Parse functions for different expression types
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == lexer.TRUE}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	// Skip optional semicolons/newlines after opening paren
	for p.curToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	exp := p.parseExpression(LOWEST)

	// Skip optional semicolons/newlines before closing paren
	for p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = []ast.HashPair{}

	// Skip optional semicolons/newlines after opening brace
	for p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	// Handle empty hash {}
	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		return hash
	}

	p.nextToken()

	// Parse first key-value pair
	key := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	value := p.parseExpression(LOWEST)
	hash.Pairs = append(hash.Pairs, ast.HashPair{Key: key, Value: value})

	// Parse remaining key-value pairs
	for p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.SEMICOLON {
		// Skip comma or semicolon/newline
		p.nextToken()
		
		// Skip any additional semicolons/newlines
		for p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
		
		// Check if we've reached the end
		if p.peekToken.Type == lexer.RBRACE {
			break
		}
		
		p.nextToken()

		key := p.parseExpression(LOWEST)
		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs = append(hash.Pairs, ast.HashPair{Key: key, Value: value})
	}

	// Skip optional semicolons/newlines before closing brace
	for p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	// Skip optional semicolons/newlines after opening bracket
	for p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA || p.peekToken.Type == lexer.SEMICOLON {
		// Skip comma or semicolon/newline
		p.nextToken()
		
		// Skip any additional semicolons/newlines
		for p.peekToken.Type == lexer.SEMICOLON {
			p.nextToken()
		}
		
		// Check if we've reached the end
		if p.peekToken.Type == end {
			break
		}
		
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// Helper functions
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}


func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Error handling
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("line %d:%d: expected next token to be %s, got %s instead",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("line %d:%d: no prefix parse function for %s found", 
		p.curToken.Line, p.curToken.Column, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == lexer.ELSE {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		// Skip comments and newlines
		if p.curToken.Type == lexer.COMMENT || p.curToken.Type == lexer.SEMICOLON {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	return &ast.BreakStatement{Token: p.curToken}
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	return &ast.ContinueStatement{Token: p.curToken}
}

func (p *Parser) parseSwitchStatement() *ast.SwitchStatement {
	stmt := &ast.SwitchStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse case and default clauses
	for p.peekToken.Type != lexer.RBRACE && p.peekToken.Type != lexer.EOF {
		p.nextToken()
		
		// Skip semicolons and comments
		if p.curToken.Type == lexer.SEMICOLON || p.curToken.Type == lexer.COMMENT {
			continue
		}
		
		if p.curToken.Type == lexer.CASE {
			caseClause := p.parseCaseClause()
			if caseClause != nil {
				stmt.Cases = append(stmt.Cases, caseClause)
			}
		} else if p.curToken.Type == lexer.DEFAULT {
			if stmt.Default != nil {
				msg := fmt.Sprintf("line %d:%d: switch statement can only have one default clause", p.curToken.Line, p.curToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			}
			stmt.Default = p.parseDefaultClause()
		} else {
			msg := fmt.Sprintf("line %d:%d: expected 'case' or 'default', got %s", p.curToken.Line, p.curToken.Column, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return stmt
}

func (p *Parser) parseCaseClause() *ast.CaseClause {
	clause := &ast.CaseClause{Token: p.curToken}

	p.nextToken()
	clause.Values = append(clause.Values, p.parseExpression(LOWEST))

	// Handle multiple values separated by commas
	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // consume comma
		p.nextToken() // move to next value
		clause.Values = append(clause.Values, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Parse the body - collect statements until we hit case, default, or }
	clause.Body = p.parseCaseBody()

	return clause
}

func (p *Parser) parseDefaultClause() *ast.DefaultClause {
	clause := &ast.DefaultClause{Token: p.curToken}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Parse the body - collect statements until we hit case, default, or }
	clause.Body = p.parseCaseBody()

	return clause
}

func (p *Parser) parseCaseBody() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for p.peekToken.Type != lexer.CASE && 
		p.peekToken.Type != lexer.DEFAULT && 
		p.peekToken.Type != lexer.RBRACE && 
		p.peekToken.Type != lexer.EOF {
		
		p.nextToken()
		if p.curToken.Type == lexer.SEMICOLON || p.curToken.Type == lexer.COMMENT {
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse init statement (can be assignment or expression)
	p.nextToken()
	if p.curToken.Type != lexer.SEMICOLON {
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			stmt.Init = p.parseForAssignmentStatement()
		} else {
			stmt.Init = p.parseForExpressionStatement()
		}
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Parse condition
	p.nextToken()
	if p.curToken.Type != lexer.SEMICOLON {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Parse update statement - can be assignment or expression
	p.nextToken()
	if p.curToken.Type != lexer.RPAREN {
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			stmt.Update = p.parseForAssignmentStatement()
		} else {
			stmt.Update = p.parseForExpressionStatement()
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}

	if p.curToken.Type != lexer.IDENT {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseForExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// parseImportStatement parses import statements like "import { func, var } from "module""
// and "import { func as alias } from "module""
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Items = []*ast.ImportItem{}

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		return stmt
	}

	p.nextToken()

	// Parse first import item
	item := p.parseImportItem()
	if item == nil {
		return nil
	}
	stmt.Items = append(stmt.Items, item)

	// Parse remaining import items
	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		item := p.parseImportItem()
		if item == nil {
			return nil
		}
		stmt.Items = append(stmt.Items, item)
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	if !p.expectPeek(lexer.FROM) {
		return nil
	}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	stmt.Module = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}

// parseImportItem parses a single import item like "func" or "func as alias"
func (p *Parser) parseImportItem() *ast.ImportItem {
	if p.curToken.Type != lexer.IDENT {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	
	item := &ast.ImportItem{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	
	// Check for "as alias" syntax
	if p.peekToken.Type == lexer.AS {
		p.nextToken() // consume 'as'
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		item.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	
	return item
}

// parseExportStatement parses export statements like "export var = 42"
func (p *Parser) parseExportStatement() *ast.ExportStatement {
	stmt := &ast.ExportStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check if there's an assignment
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // consume =
		p.nextToken() // move to value
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

// parseModuleAccess parses module member access like "module.function"
func (p *Parser) parseModuleAccess(left ast.Expression) ast.Expression {
	moduleAccess := &ast.ModuleAccess{Token: p.curToken}

	// Left side should be an identifier (module name)
	if ident, ok := left.(*ast.Identifier); ok {
		moduleAccess.Module = ident
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier before '.', got %T", left))
		return nil
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	moduleAccess.Member = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return moduleAccess
}

// parsePropertyAccess parses property access like "object.property" or class instantiation like "ClassName.new()"
func (p *Parser) parsePropertyAccess(left ast.Expression) ast.Expression {
	dotToken := p.curToken // store the '.' token

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	propertyName := p.curToken.Literal

	// Check if this is a "new" method call for class instantiation
	if propertyName == "new" && p.peekToken.Type == lexer.LPAREN {
		// This is class instantiation like "ClassName.new()"
		if className, ok := left.(*ast.Identifier); ok {
			newExpr := &ast.NewExpression{
				Token:     dotToken,
				ClassName: className,
			}

			// Parse the arguments
			p.nextToken() // move to '('
			newExpr.Arguments = p.parseExpressionList(lexer.RPAREN)
			return newExpr
		}
	}

	// Regular property access
	propertyAccess := &ast.PropertyAccess{
		Token:  dotToken,
		Object: left,
		Property: &ast.Identifier{Token: p.curToken, Value: propertyName},
	}

	return propertyAccess
}

// parseThrowStatement parses throw statements like "throw ErrorType("message")"
func (p *Parser) parseThrowStatement() *ast.ThrowStatement {
	stmt := &ast.ThrowStatement{Token: p.curToken}

	p.nextToken()

	stmt.Expression = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseTryStatement parses try-catch-finally blocks
func (p *Parser) parseTryStatement() *ast.TryStatement {
	stmt := &ast.TryStatement{Token: p.curToken}

	// Parse try block
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.TryBlock = p.parseBlockStatement()

	// Parse catch clauses (there can be multiple)
	stmt.CatchClauses = []*ast.CatchClause{}
	for p.peekToken.Type == lexer.CATCH {
		p.nextToken() // consume catch token
		catchClause := p.parseCatchClause()
		if catchClause != nil {
			stmt.CatchClauses = append(stmt.CatchClauses, catchClause)
		}
	}

	// Parse optional finally block
	if p.peekToken.Type == lexer.FINALLY {
		p.nextToken() // consume finally token
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
		stmt.FinallyBlock = p.parseBlockStatement()
	}

	// Must have at least one catch clause or a finally block
	if len(stmt.CatchClauses) == 0 && stmt.FinallyBlock == nil {
		p.errors = append(p.errors, "try statement must have at least one catch clause or a finally block")
		return nil
	}

	return stmt
}

// parseCatchClause parses catch clauses like "catch (ErrorType error) { ... }"
func (p *Parser) parseCatchClause() *ast.CatchClause {
	clause := &ast.CatchClause{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Check if we have an error type followed by variable name
	// catch (ErrorType error) or catch (error)
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	firstIdent := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check if there's a second identifier (error variable)
	if p.peekToken.Type == lexer.IDENT {
		// First identifier is the error type, second is the variable
		clause.ErrorType = firstIdent
		p.nextToken()
		clause.ErrorVar = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else {
		// Only one identifier, it's the error variable
		clause.ErrorVar = firstIdent
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	clause.Body = p.parseBlockStatement()

	return clause
}

// parseClassDeclaration parses class declarations like "class ClassName < SuperClass { methods... }"
func (p *Parser) parseClassDeclaration() ast.Statement {
  stmt := &ast.ClassDeclaration{Token: p.curToken}

  if !p.expectPeek(lexer.IDENT) {
    return nil
  }

  stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

  // Check for inheritance (optional)
  if p.peekToken.Type == lexer.LT {
    p.nextToken() // consume '<'
    if !p.expectPeek(lexer.IDENT) {
      return nil
    }
    stmt.SuperClass = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
  }

  if !p.expectPeek(lexer.LBRACE) {
    return nil
  }

  stmt.Body = p.parseClassBody()
  return stmt
}

// parseClassBody parses class body with special handling for method declarations
func (p *Parser) parseClassBody() *ast.BlockStatement {
  block := &ast.BlockStatement{Token: p.curToken}
  block.Statements = []ast.Statement{}

  p.nextToken()

  for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
    // Skip comments and newlines
    if p.curToken.Type == lexer.COMMENT || p.curToken.Type == lexer.SEMICOLON {
      p.nextToken()
      continue
    }

    var stmt ast.Statement
    
    // Handle fn as method declarations within class body
    if p.curToken.Type == lexer.FN {
      stmt = p.parseMethodDeclaration()
    } else {
      // Parse other statements normally
      stmt = p.parseStatement()
    }
    
    if stmt != nil {
      block.Statements = append(block.Statements, stmt)
    }
    p.nextToken()
  }

  return block
}

// parseInstanceVariable parses instance variables like "@variable"
func (p *Parser) parseInstanceVariable() ast.Expression {
  inst := &ast.InstanceVariable{Token: p.curToken}

  if !p.expectPeek(lexer.IDENT) {
    return nil
  }

  inst.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
  return inst
}

// parseSuperExpression parses super() calls to parent methods
func (p *Parser) parseSuperExpression() ast.Expression {
  expr := &ast.SuperExpression{Token: p.curToken}

  if !p.expectPeek(lexer.LPAREN) {
    return nil
  }

  expr.Arguments = p.parseExpressionList(lexer.RPAREN)
  return expr
}

// parseMethodDeclaration parses method declarations like "fn methodName() { body }"
func (p *Parser) parseMethodDeclaration() ast.Statement {
  method := &ast.MethodDeclaration{Token: p.curToken}

  // Method name can be IDENT or INITIALIZE keyword
  if p.curToken.Type != lexer.FN {
    return nil
  }
  
  p.nextToken() // Move past 'fn'
  
  if p.curToken.Type == lexer.IDENT || p.curToken.Type == lexer.INITIALIZE {
    method.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
  } else {
    return nil
  }

  if !p.expectPeek(lexer.LPAREN) {
    return nil
  }

  method.Parameters = p.parseFunctionParameters()

  if !p.expectPeek(lexer.LBRACE) {
    return nil
  }

  method.Body = p.parseBlockStatement()

  return method
}


// parseNewExpression is handled by parsePropertyAccess when it encounters "ClassName.new"
// This is called from the property access parsing when we see ".new("

// shouldStopParsing determines if we should stop parsing infix expressions
// It returns true if we hit a semicolon and the next token isn't a DOT (method chaining)
func (p *Parser) shouldStopParsing(precedence int) bool {
	// If we're not at a semicolon, use normal precedence rules
	if p.peekToken.Type != lexer.SEMICOLON {
		return precedence >= p.peekPrecedence()
	}
	
	// We're at a semicolon - check if this is method chaining or a continuation operator
	// Create a temporary copy of the lexer state to look ahead
	tempLexer := *p.l
	nextToken := tempLexer.NextToken() // Get the token after semicolon
	
	// If the next significant token is a DOT, it's method chaining
	if nextToken.Type == lexer.DOT {
		// Skip the semicolon in the actual parser
		p.nextToken()
		return false // Don't stop parsing
	}
	
	// If the next token is a binary operator (&&, ||, +, -, etc.), continue parsing
	if nextToken.Type == lexer.AND || nextToken.Type == lexer.OR ||
		nextToken.Type == lexer.PLUS || nextToken.Type == lexer.MINUS ||
		nextToken.Type == lexer.MULT || nextToken.Type == lexer.DIV ||
		nextToken.Type == lexer.EQ || nextToken.Type == lexer.NOT_EQ ||
		nextToken.Type == lexer.LT || nextToken.Type == lexer.GT ||
		nextToken.Type == lexer.LTE || nextToken.Type == lexer.GTE {
		// Skip the semicolon in the actual parser
		p.nextToken()
		return false // Don't stop parsing
	}
	
	// Otherwise, stop at semicolon
	return true
}