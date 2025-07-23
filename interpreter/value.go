package interpreter

import (
	"fmt"
	"strings"

	"rush/ast"
)

// ValueType represents the type of a value
type ValueType string

const (
	INTEGER_VALUE  ValueType = "INTEGER"
	FLOAT_VALUE    ValueType = "FLOAT"
	STRING_VALUE   ValueType = "STRING"
	BOOLEAN_VALUE  ValueType = "BOOLEAN"
	ARRAY_VALUE    ValueType = "ARRAY"
	HASH_VALUE     ValueType = "HASH"
	NULL_VALUE     ValueType = "NULL"
	FUNCTION_VALUE  ValueType = "FUNCTION"
	BUILTIN_VALUE   ValueType = "BUILTIN"
	RETURN_VALUE    ValueType = "RETURN_VALUE"
	BREAK_VALUE     ValueType = "BREAK_VALUE"
	CONTINUE_VALUE  ValueType = "CONTINUE_VALUE"
	EXCEPTION_VALUE ValueType = "EXCEPTION"
	CLASS_VALUE     ValueType = "CLASS"
	INSTANCE_VALUE  ValueType = "INSTANCE"
	BOUND_METHOD_VALUE ValueType = "BOUND_METHOD"
)

// Value represents a value in the Rush language
type Value interface {
	Type() ValueType
	Inspect() string
}

// Integer represents integer values
type Integer struct {
	Value int64
}

func (i *Integer) Type() ValueType { return INTEGER_VALUE }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Float represents floating-point values
type Float struct {
	Value float64
}

func (f *Float) Type() ValueType { return FLOAT_VALUE }
func (f *Float) Inspect() string { return fmt.Sprintf("%g", f.Value) }

// String represents string values
type String struct {
	Value string
}

func (s *String) Type() ValueType { return STRING_VALUE }
func (s *String) Inspect() string { return s.Value }

// Boolean represents boolean values
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ValueType { return BOOLEAN_VALUE }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Array represents array values
type Array struct {
	Elements []Value
}

func (a *Array) Type() ValueType { return ARRAY_VALUE }
func (a *Array) Inspect() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

// HashKey represents a key in a hash for efficient storage
type HashKey struct {
	Type  ValueType
	Value interface{} // int64, string, bool, float64
}

// Hash represents hash/dictionary values
type Hash struct {
	Pairs map[HashKey]Value
	Keys  []Value // maintain insertion order
}

func (h *Hash) Type() ValueType { return HASH_VALUE }
func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, key := range h.Keys {
		value := h.Pairs[CreateHashKey(key)]
		pairs = append(pairs, fmt.Sprintf("%s: %s", key.Inspect(), value.Inspect()))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

// CreateHashKey creates a HashKey from a Value for internal storage
func CreateHashKey(value Value) HashKey {
	switch val := value.(type) {
	case *Integer:
		return HashKey{Type: INTEGER_VALUE, Value: val.Value}
	case *String:
		return HashKey{Type: STRING_VALUE, Value: val.Value}
	case *Boolean:
		return HashKey{Type: BOOLEAN_VALUE, Value: val.Value}
	case *Float:
		return HashKey{Type: FLOAT_VALUE, Value: val.Value}
	default:
		// This should not happen in practice due to type validation
		return HashKey{Type: NULL_VALUE, Value: nil}
	}
}

// Null represents null/nil values
type Null struct{}

func (n *Null) Type() ValueType { return NULL_VALUE }
func (n *Null) Inspect() string { return "null" }

// Function represents function values
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ValueType { return FUNCTION_VALUE }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
}

// ReturnValue wraps values to signal a return statement
type ReturnValue struct {
	Value Value
}

func (rv *ReturnValue) Type() ValueType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// BreakValue signals a break statement
type BreakValue struct{}

func (bv *BreakValue) Type() ValueType { return BREAK_VALUE }
func (bv *BreakValue) Inspect() string { return "break" }

// ContinueValue signals a continue statement
type ContinueValue struct{}

func (cv *ContinueValue) Type() ValueType { return CONTINUE_VALUE }
func (cv *ContinueValue) Inspect() string { return "continue" }

// BuiltinFunction represents built-in functions
type BuiltinFunction struct {
	Fn func(args ...Value) Value
}

func (bf *BuiltinFunction) Type() ValueType { return BUILTIN_VALUE }
func (bf *BuiltinFunction) Inspect() string { return "builtin function" }

// Exception wraps errors to signal exception propagation
type Exception struct {
	Error Value // Can be any error value
}

func (ex *Exception) Type() ValueType { return EXCEPTION_VALUE }
func (ex *Exception) Inspect() string { return ex.Error.Inspect() }

// Class represents a class definition
type Class struct {
  Name       string
  SuperClass *Class
  Methods    map[string]*Function
  Env        *Environment
}

func (c *Class) Type() ValueType { return CLASS_VALUE }
func (c *Class) Inspect() string { return fmt.Sprintf("class %s", c.Name) }

// Object represents an instance of a class
type Object struct {
  Class            *Class
  InstanceVars     map[string]Value
  Env              *Environment
}

func (o *Object) Type() ValueType { return INSTANCE_VALUE }
func (o *Object) Inspect() string { return fmt.Sprintf("#<%s:0x%p>", o.Class.Name, o) }

// BoundMethod represents a method bound to a specific object instance
type BoundMethod struct {
  Method *Function
  Instance *Object
}

func (bm *BoundMethod) Type() ValueType { return BOUND_METHOD_VALUE }
func (bm *BoundMethod) Inspect() string { 
  return fmt.Sprintf("#<BoundMethod:0x%p>", bm) 
}

// IsTruthy returns whether a value is considered truthy
func IsTruthy(val Value) bool {
	switch val {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// Global singleton values for commonly used values
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// NewException creates a new exception for error propagation
func NewException(error Value) *Exception {
	return &Exception{Error: error}
}