package interpreter

import (
	"fmt"
	"strings"
)

// ValueType represents the type of a value
type ValueType string

const (
	INTEGER_VALUE ValueType = "INTEGER"
	FLOAT_VALUE   ValueType = "FLOAT"
	STRING_VALUE  ValueType = "STRING"
	BOOLEAN_VALUE ValueType = "BOOLEAN"
	ARRAY_VALUE   ValueType = "ARRAY"
	NULL_VALUE    ValueType = "NULL"
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

// Null represents null/nil values
type Null struct{}

func (n *Null) Type() ValueType { return NULL_VALUE }
func (n *Null) Inspect() string { return "null" }

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