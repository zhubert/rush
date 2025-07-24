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
	COMPILED_FUNCTION_VALUE ValueType = "COMPILED_FUNCTION"
	CLOSURE_VALUE   ValueType = "CLOSURE"
	RETURN_VALUE    ValueType = "RETURN_VALUE"
	BREAK_VALUE     ValueType = "BREAK_VALUE"
	CONTINUE_VALUE  ValueType = "CONTINUE_VALUE"
	EXCEPTION_VALUE ValueType = "EXCEPTION"
	CLASS_VALUE     ValueType = "CLASS"
	INSTANCE_VALUE  ValueType = "INSTANCE"
	BOUND_METHOD_VALUE ValueType = "BOUND_METHOD"
	HASH_METHOD_VALUE   ValueType = "HASH_METHOD"
	STRING_METHOD_VALUE ValueType = "STRING_METHOD"
	ARRAY_METHOD_VALUE  ValueType = "ARRAY_METHOD"
	NUMBER_METHOD_VALUE ValueType = "NUMBER_METHOD"
	FILE_VALUE          ValueType = "FILE"
	DIRECTORY_VALUE     ValueType = "DIRECTORY"
	PATH_VALUE          ValueType = "PATH"
	FILE_METHOD_VALUE   ValueType = "FILE_METHOD"
	DIRECTORY_METHOD_VALUE ValueType = "DIRECTORY_METHOD"
	PATH_METHOD_VALUE   ValueType = "PATH_METHOD"
	JSON_VALUE          ValueType = "JSON"
	JSON_METHOD_VALUE   ValueType = "JSON_METHOD"
	JSON_NAMESPACE_VALUE ValueType = "JSON_NAMESPACE"
	TIME_VALUE          ValueType = "TIME"
	TIME_METHOD_VALUE   ValueType = "TIME_METHOD"
	TIME_NAMESPACE_VALUE ValueType = "TIME_NAMESPACE"
	DURATION_VALUE      ValueType = "DURATION"
	DURATION_METHOD_VALUE ValueType = "DURATION_METHOD"
	DURATION_NAMESPACE_VALUE ValueType = "DURATION_NAMESPACE"
	TIMEZONE_VALUE      ValueType = "TIMEZONE"
	TIMEZONE_METHOD_VALUE ValueType = "TIMEZONE_METHOD" 
	TIMEZONE_NAMESPACE_VALUE ValueType = "TIMEZONE_NAMESPACE"
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

// HashMethod represents a method bound to a specific hash instance
type HashMethod struct {
  Hash   *Hash
  Method string
}

func (hm *HashMethod) Type() ValueType { return HASH_METHOD_VALUE }
func (hm *HashMethod) Inspect() string { 
  return fmt.Sprintf("#<HashMethod:%s on %s>", hm.Method, hm.Hash.Inspect()) 
}

// StringMethod represents a method bound to a specific string instance
type StringMethod struct {
  String *String
  Method string
}

func (sm *StringMethod) Type() ValueType { return STRING_METHOD_VALUE }
func (sm *StringMethod) Inspect() string { 
  return fmt.Sprintf("#<StringMethod:%s on %s>", sm.Method, sm.String.Inspect()) 
}

// ArrayMethod represents a method bound to a specific array instance
type ArrayMethod struct {
  Array  *Array
  Method string
}

func (am *ArrayMethod) Type() ValueType { return ARRAY_METHOD_VALUE }
func (am *ArrayMethod) Inspect() string { 
  return fmt.Sprintf("#<ArrayMethod:%s on %s>", am.Method, am.Array.Inspect()) 
}

// NumberMethod represents a method bound to a specific number instance (integer or float)
type NumberMethod struct {
  Number Value  // Can be *Integer or *Float
  Method string
}

func (nm *NumberMethod) Type() ValueType { return NUMBER_METHOD_VALUE }
func (nm *NumberMethod) Inspect() string { 
  return fmt.Sprintf("#<NumberMethod:%s on %s>", nm.Method, nm.Number.Inspect()) 
}

// File represents a file in the filesystem
type File struct {
  Path string
  Handle interface{} // Will hold actual file handle when opened
  IsOpen bool
}

func (f *File) Type() ValueType { return FILE_VALUE }
func (f *File) Inspect() string {
  if f.IsOpen {
    return fmt.Sprintf("#<File:%s (open)>", f.Path)
  }
  return fmt.Sprintf("#<File:%s (closed)>", f.Path)
}

// Directory represents a directory in the filesystem
type Directory struct {
  Path string
}

func (d *Directory) Type() ValueType { return DIRECTORY_VALUE }
func (d *Directory) Inspect() string {
  return fmt.Sprintf("#<Directory:%s>", d.Path)
}

// Path represents a filesystem path
type Path struct {
  Value string
}

func (p *Path) Type() ValueType { return PATH_VALUE }
func (p *Path) Inspect() string {
  return fmt.Sprintf("#<Path:%s>", p.Value)
}

// FileMethod represents a method bound to a specific file instance
type FileMethod struct {
  File   *File
  Method string
}

func (fm *FileMethod) Type() ValueType { return FILE_METHOD_VALUE }
func (fm *FileMethod) Inspect() string {
  return fmt.Sprintf("#<FileMethod:%s on %s>", fm.Method, fm.File.Inspect())
}

// DirectoryMethod represents a method bound to a specific directory instance
type DirectoryMethod struct {
  Directory *Directory
  Method    string
}

func (dm *DirectoryMethod) Type() ValueType { return DIRECTORY_METHOD_VALUE }
func (dm *DirectoryMethod) Inspect() string {
  return fmt.Sprintf("#<DirectoryMethod:%s on %s>", dm.Method, dm.Directory.Inspect())
}

// PathMethod represents a method bound to a specific path instance
type PathMethod struct {
  Path   *Path
  Method string
}

func (pm *PathMethod) Type() ValueType { return PATH_METHOD_VALUE }
func (pm *PathMethod) Inspect() string {
  return fmt.Sprintf("#<PathMethod:%s on %s>", pm.Method, pm.Path.Inspect())
}

// JSON represents a JSON object with parsed data
type JSON struct {
  Data Value // The underlying parsed data (Hash, Array, String, Integer, Float, Boolean, or Null)
}

func (j *JSON) Type() ValueType { return JSON_VALUE }
func (j *JSON) Inspect() string {
  return fmt.Sprintf("#<JSON:%s>", j.Data.Inspect())
}

// JSONMethod represents a method bound to a specific JSON instance
type JSONMethod struct {
  JSON   *JSON
  Method string
}

func (jm *JSONMethod) Type() ValueType { return JSON_METHOD_VALUE }
func (jm *JSONMethod) Inspect() string {
  return fmt.Sprintf("#<JSONMethod:%s on %s>", jm.Method, jm.JSON.Inspect())
}

// JSONNamespace represents the JSON namespace with static methods
type JSONNamespace struct{}

func (jn *JSONNamespace) Type() ValueType { return JSON_NAMESPACE_VALUE }
func (jn *JSONNamespace) Inspect() string {
  return "#<JSONNamespace>"
}

// Time represents a specific moment in time
type Time struct {
  Value    int64     // Unix timestamp in nanoseconds
  Location string    // Timezone location (e.g., "UTC", "Local", "America/New_York")
}

func (t *Time) Type() ValueType { return TIME_VALUE }
func (t *Time) Inspect() string {
  return fmt.Sprintf("#<Time:%d in %s>", t.Value, t.Location)
}

// Duration represents a span of time
type Duration struct {
  Value int64 // Duration in nanoseconds
}

func (d *Duration) Type() ValueType { return DURATION_VALUE }
func (d *Duration) Inspect() string {
  return fmt.Sprintf("#<Duration:%dns>", d.Value)
}

// TimeZone represents timezone information
type TimeZone struct {
  Name   string  // Timezone name (e.g., "UTC", "America/New_York")
  Offset int     // Offset from UTC in seconds
}

func (tz *TimeZone) Type() ValueType { return TIMEZONE_VALUE }
func (tz *TimeZone) Inspect() string {
  return fmt.Sprintf("#<TimeZone:%s offset:%d>", tz.Name, tz.Offset)
}

// TimeMethod represents a method bound to a specific time instance
type TimeMethod struct {
  Time   *Time
  Method string
}

func (tm *TimeMethod) Type() ValueType { return TIME_METHOD_VALUE }
func (tm *TimeMethod) Inspect() string {
  return fmt.Sprintf("#<TimeMethod:%s on %s>", tm.Method, tm.Time.Inspect())
}

// DurationMethod represents a method bound to a specific duration instance
type DurationMethod struct {
  Duration *Duration
  Method   string
}

func (dm *DurationMethod) Type() ValueType { return DURATION_METHOD_VALUE }
func (dm *DurationMethod) Inspect() string {
  return fmt.Sprintf("#<DurationMethod:%s on %s>", dm.Method, dm.Duration.Inspect())
}

// TimeZoneMethod represents a method bound to a specific timezone instance
type TimeZoneMethod struct {
  TimeZone *TimeZone
  Method   string
}

func (tzm *TimeZoneMethod) Type() ValueType { return TIMEZONE_METHOD_VALUE }
func (tzm *TimeZoneMethod) Inspect() string {
  return fmt.Sprintf("#<TimeZoneMethod:%s on %s>", tzm.Method, tzm.TimeZone.Inspect())
}

// TimeNamespace represents the Time namespace with static methods
type TimeNamespace struct{}

func (tn *TimeNamespace) Type() ValueType { return TIME_NAMESPACE_VALUE }
func (tn *TimeNamespace) Inspect() string {
  return "#<TimeNamespace>"
}

// DurationNamespace represents the Duration namespace with static methods
type DurationNamespace struct{}

func (dn *DurationNamespace) Type() ValueType { return DURATION_NAMESPACE_VALUE }
func (dn *DurationNamespace) Inspect() string {
  return "#<DurationNamespace>"
}

// TimeZoneNamespace represents the TimeZone namespace with static methods
type TimeZoneNamespace struct{}

func (tzn *TimeZoneNamespace) Type() ValueType { return TIMEZONE_NAMESPACE_VALUE }
func (tzn *TimeZoneNamespace) Inspect() string {
  return "#<TimeZoneNamespace>"
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

// CompiledFunction represents a compiled function
type CompiledFunction struct {
	Instructions  []byte // Bytecode instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() ValueType { return COMPILED_FUNCTION_VALUE }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

// Closure represents a closure (function with captured variables)
type Closure struct {
	Fn   *CompiledFunction
	Free []Value
}

func (c *Closure) Type() ValueType { return CLOSURE_VALUE }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}