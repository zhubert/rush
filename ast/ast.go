package ast

import (
	"bytes"
	"strings"

	"rush/lexer"
)

// Node represents any node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents statements (don't produce values)
type Statement interface {
	Node
	statementNode()
}

// Expression represents expressions (produce values)
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// AssignmentStatement represents variable assignments like "a = 5"
type AssignmentStatement struct {
	Token lexer.Token // the identifier token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	return out.String()
}

// IndexAssignmentStatement represents array element assignments like "arr[0] = 5"
type IndexAssignmentStatement struct {
	Token lexer.Token      // the '=' token
	Left  *IndexExpression // the array[index] being assigned to
	Value Expression       // the value being assigned
}

func (ias *IndexAssignmentStatement) statementNode()       {}
func (ias *IndexAssignmentStatement) TokenLiteral() string { return ias.Token.Literal }
func (ias *IndexAssignmentStatement) String() string {
	var out bytes.Buffer
	if ias.Left != nil {
		out.WriteString(ias.Left.String())
	}
	out.WriteString(" = ")
	if ias.Value != nil {
		out.WriteString(ias.Value.String())
	}
	return out.String()
}

// Identifier represents identifiers like variable names
type Identifier struct {
	Token lexer.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents integer literals like 5, 10, 42
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents float literals like 3.14, 2.5
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents string literals like "hello"
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents boolean literals like true, false
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// InfixExpression represents infix expressions like "a + b", "x > y"
type InfixExpression struct {
	Token    lexer.Token // the operator token, e.g. +, -, *, /
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// PrefixExpression represents prefix expressions like "!true", "-5"
type PrefixExpression struct {
	Token    lexer.Token // the prefix token, e.g. !, -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral represents array literals like [1, 2, 3]
type ArrayLiteral struct {
	Token    lexer.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// ExpressionStatement represents expressions used as statements
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IfExpression represents if expressions like "if (condition) { consequence } else { alternative }"
type IfExpression struct {
	Token       lexer.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// BlockStatement represents a block of statements like "{ statement1; statement2; }"
type BlockStatement struct {
	Token      lexer.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

// FunctionLiteral represents function declarations like "fn(x, y) { return x + y; }"
type FunctionLiteral struct {
	Token      lexer.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents function calls like "add(1, 2)"
type CallExpression struct {
	Token     lexer.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ReturnStatement represents return statements like "return 5;"
type ReturnStatement struct {
	Token       lexer.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// IndexExpression represents array/string indexing like "arr[1]" or "str[0]"
type IndexExpression struct {
	Token lexer.Token // the '[' token
	Left  Expression  // the array or string being indexed
	Index Expression  // the index expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// WhileStatement represents while loop statements like "while (condition) { body }"
type WhileStatement struct {
	Token     lexer.Token // the 'while' token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement represents for loop statements like "for (init; condition; update) { body }"
type ForStatement struct {
	Token     lexer.Token // the 'for' token
	Init      Statement   // initialization statement (can be assignment or expression)
	Condition Expression  // loop condition
	Update    Statement   // update statement (can be assignment or expression)
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for(")
	if fs.Init != nil {
		out.WriteString(fs.Init.String())
	}
	out.WriteString(";")
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString(";")
	if fs.Update != nil {
		out.WriteString(fs.Update.String())
	}
	out.WriteString(") ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// ImportItem represents a single import with optional alias
type ImportItem struct {
	Name  *Identifier // original name
	Alias *Identifier // alias (can be nil)
}

// ImportStatement represents import statements like "import { func, var } from "module""
type ImportStatement struct {
	Token     lexer.Token      // the 'import' token
	Items     []*ImportItem    // imported items with optional aliases
	Module    *StringLiteral   // module path
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("import { ")
	items := []string{}
	for _, item := range is.Items {
		if item.Alias != nil {
			items = append(items, item.Name.String()+" as "+item.Alias.String())
		} else {
			items = append(items, item.Name.String())
		}
	}
	out.WriteString(strings.Join(items, ", "))
	out.WriteString(" } from ")
	out.WriteString(is.Module.String())
	return out.String()
}

// ExportStatement represents export statements like "export func" or "export var = 42"
type ExportStatement struct {
	Token     lexer.Token // the 'export' token
	Name      *Identifier // exported name
	Value     Expression  // exported value (can be nil for function declarations)
}

func (es *ExportStatement) statementNode()       {}
func (es *ExportStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("export ")
	out.WriteString(es.Name.String())
	if es.Value != nil {
		out.WriteString(" = ")
		out.WriteString(es.Value.String())
	}
	return out.String()
}

// PropertyAccess represents property access like "object.property" or "module.function"
type PropertyAccess struct {
	Token  lexer.Token // the dot token
	Object Expression  // the object being accessed (can be any expression)
	Property *Identifier // property name
}

func (pa *PropertyAccess) expressionNode()      {}
func (pa *PropertyAccess) TokenLiteral() string { return pa.Token.Literal }
func (pa *PropertyAccess) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pa.Object.String())
	out.WriteString(".")
	out.WriteString(pa.Property.String())
	out.WriteString(")")
	return out.String()
}

// ModuleAccess represents module member access like "module.function" (deprecated, use PropertyAccess)
type ModuleAccess struct {
	Token  lexer.Token // the identifier token (module name)
	Module *Identifier // module name
	Member *Identifier // member name
}

func (ma *ModuleAccess) expressionNode()      {}
func (ma *ModuleAccess) TokenLiteral() string { return ma.Token.Literal }
func (ma *ModuleAccess) String() string {
	var out bytes.Buffer
	out.WriteString(ma.Module.String())
	out.WriteString(".")
	out.WriteString(ma.Member.String())
	return out.String()
}

// ThrowStatement represents throw statements like "throw ErrorType("message")"
type ThrowStatement struct {
	Token      lexer.Token // the 'throw' token
	Expression Expression  // the expression to throw
}

func (ts *ThrowStatement) statementNode()       {}
func (ts *ThrowStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *ThrowStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ts.TokenLiteral() + " ")
	if ts.Expression != nil {
		out.WriteString(ts.Expression.String())
	}
	return out.String()
}

// CatchClause represents catch clauses like "catch (ErrorType error) { ... }"
type CatchClause struct {
	Token     lexer.Token     // the 'catch' token
	ErrorType *Identifier     // optional error type filter (can be nil)
	ErrorVar  *Identifier     // error variable name
	Body      *BlockStatement // catch block body
}

func (cc *CatchClause) TokenLiteral() string { return cc.Token.Literal }
func (cc *CatchClause) String() string {
	var out bytes.Buffer
	out.WriteString("catch (")
	if cc.ErrorType != nil {
		out.WriteString(cc.ErrorType.String())
		out.WriteString(" ")
	}
	out.WriteString(cc.ErrorVar.String())
	out.WriteString(") ")
	out.WriteString(cc.Body.String())
	return out.String()
}

// TryStatement represents try-catch-finally blocks
type TryStatement struct {
	Token        lexer.Token      // the 'try' token
	TryBlock     *BlockStatement  // try block
	CatchClauses []*CatchClause   // catch clauses (can be multiple)
	FinallyBlock *BlockStatement  // finally block (optional)
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	var out bytes.Buffer
	out.WriteString("try ")
	out.WriteString(ts.TryBlock.String())
	
	for _, catchClause := range ts.CatchClauses {
		out.WriteString(" ")
		out.WriteString(catchClause.String())
	}
	
	if ts.FinallyBlock != nil {
		out.WriteString(" finally ")
		out.WriteString(ts.FinallyBlock.String())
	}
	
	return out.String()
}

// ClassDeclaration represents class definitions like "class ClassName { methods... }"
type ClassDeclaration struct {
  Token      lexer.Token // the 'class' token
  Name       *Identifier
  SuperClass *Identifier     // optional superclass (can be nil)
  Methods    []*MethodDeclaration
  Body       *BlockStatement // class body containing methods and instance variables
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ClassDeclaration) String() string {
  var out bytes.Buffer
  out.WriteString("class ")
  out.WriteString(cd.Name.String())
  if cd.SuperClass != nil {
    out.WriteString(" < ")
    out.WriteString(cd.SuperClass.String())
  }
  out.WriteString(" ")
  out.WriteString(cd.Body.String())
  return out.String()
}

// MethodDeclaration represents method definitions like "fn methodName() { body }"
type MethodDeclaration struct {
  Token      lexer.Token // the 'fn' token
  Name       *Identifier
  Parameters []*Identifier
  Body       *BlockStatement
}

func (md *MethodDeclaration) statementNode()       {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }
func (md *MethodDeclaration) String() string {
  var out bytes.Buffer
  params := []string{}
  for _, p := range md.Parameters {
    params = append(params, p.String())
  }
  out.WriteString("fn ")
  out.WriteString(md.Name.String())
  out.WriteString("(")
  out.WriteString(strings.Join(params, ", "))
  out.WriteString(") ")
  out.WriteString(md.Body.String())
  return out.String()
}

// InstanceVariable represents instance variables like "@variable"
type InstanceVariable struct {
  Token lexer.Token // the '@' token
  Name  *Identifier
}

func (iv *InstanceVariable) expressionNode()      {}
func (iv *InstanceVariable) TokenLiteral() string { return iv.Token.Literal }
func (iv *InstanceVariable) String() string {
  var out bytes.Buffer
  out.WriteString("@")
  out.WriteString(iv.Name.String())
  return out.String()
}

// NewExpression represents object instantiation like "ClassName.new()"
type NewExpression struct {
  Token     lexer.Token  // the identifier token (class name)
  ClassName *Identifier
  Arguments []Expression
}

func (ne *NewExpression) expressionNode()      {}
func (ne *NewExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NewExpression) String() string {
  var out bytes.Buffer
  args := []string{}
  for _, a := range ne.Arguments {
    args = append(args, a.String())
  }
  out.WriteString(ne.ClassName.String())
  out.WriteString(".new(")
  out.WriteString(strings.Join(args, ", "))
  out.WriteString(")")
  return out.String()
}

// SuperExpression represents super() calls to parent methods
type SuperExpression struct {
  Token     lexer.Token  // the 'super' token
  Arguments []Expression
}

func (se *SuperExpression) expressionNode()      {}
func (se *SuperExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SuperExpression) String() string {
  var out bytes.Buffer
  args := []string{}
  for _, a := range se.Arguments {
    args = append(args, a.String())
  }
  out.WriteString("super(")
  out.WriteString(strings.Join(args, ", "))
  out.WriteString(")")
  return out.String()
}

// BreakStatement represents break statements
type BreakStatement struct {
	Token lexer.Token // the 'break' token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.TokenLiteral() }

// ContinueStatement represents continue statements
type ContinueStatement struct {
	Token lexer.Token // the 'continue' token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return cs.TokenLiteral() }

// SwitchStatement represents switch statements
type SwitchStatement struct {
	Token   lexer.Token    // the 'switch' token
	Value   Expression     // the expression to match against
	Cases   []*CaseClause  // the case clauses
	Default *DefaultClause // optional default clause
}

func (ss *SwitchStatement) statementNode()       {}
func (ss *SwitchStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SwitchStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.TokenLiteral())
	out.WriteString(" (")
	if ss.Value != nil {
		out.WriteString(ss.Value.String())
	}
	out.WriteString(") {")
	
	for _, c := range ss.Cases {
		out.WriteString(c.String())
	}
	
	if ss.Default != nil {
		out.WriteString(ss.Default.String())
	}
	
	out.WriteString("}")
	return out.String()
}

// CaseClause represents a case clause in a switch statement
type CaseClause struct {
	Token  lexer.Token   // the 'case' token
	Values []Expression  // the values to match (supports multiple values like case 1, 2, 3:)
	Body   *BlockStatement // the statements to execute
}

func (cc *CaseClause) statementNode()       {}
func (cc *CaseClause) TokenLiteral() string { return cc.Token.Literal }
func (cc *CaseClause) String() string {
	var out bytes.Buffer
	out.WriteString(cc.TokenLiteral())
	out.WriteString(" ")
	
	for i, value := range cc.Values {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(value.String())
	}
	
	out.WriteString(":")
	if cc.Body != nil {
		out.WriteString(cc.Body.String())
	}
	return out.String()
}

// DefaultClause represents a default clause in a switch statement
type DefaultClause struct {
	Token lexer.Token     // the 'default' token
	Body  *BlockStatement // the statements to execute
}

func (dc *DefaultClause) statementNode()       {}
func (dc *DefaultClause) TokenLiteral() string { return dc.Token.Literal }
func (dc *DefaultClause) String() string {
	var out bytes.Buffer
	out.WriteString(dc.TokenLiteral())
	out.WriteString(":")
	if dc.Body != nil {
		out.WriteString(dc.Body.String())
	}
	return out.String()
}