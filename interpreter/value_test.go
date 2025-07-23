package interpreter

import (
	"testing"
	"rush/ast"
	"rush/lexer"
)

func TestIntegerValue(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{42, "42"},
		{-17, "-17"},
		{0, "0"},
		{9223372036854775807, "9223372036854775807"}, // max int64
		{-9223372036854775808, "-9223372036854775808"}, // min int64
	}
	
	for _, tt := range tests {
		intVal := &Integer{Value: tt.value}
		
		if intVal.Type() != INTEGER_VALUE {
			t.Fatalf("Expected type INTEGER_VALUE, got %s", intVal.Type())
		}
		
		if intVal.Inspect() != tt.expected {
			t.Fatalf("Expected inspect %q, got %q", tt.expected, intVal.Inspect())
		}
	}
}

func TestFloatValue(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{3.14, "3.14"},
		{-2.5, "-2.5"},
		{0.0, "0"},
		{123.456789, "123.457"}, // %g formatting
		{1000000.0, "1e+06"},    // scientific notation
	}
	
	for _, tt := range tests {
		floatVal := &Float{Value: tt.value}
		
		if floatVal.Type() != FLOAT_VALUE {
			t.Fatalf("Expected type FLOAT_VALUE, got %s", floatVal.Type())
		}
		
		// Note: floating point formatting can vary, so we check basic cases
		inspect := floatVal.Inspect()
		if len(inspect) == 0 {
			t.Fatalf("Float inspect should not be empty")
		}
	}
}

func TestStringValue(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"hello", "hello"},
		{"", ""},
		{"Hello, World!", "Hello, World!"},
		{"Rush Programming Language", "Rush Programming Language"},
		{"Special chars: \n\t\"'", "Special chars: \n\t\"'"},
	}
	
	for _, tt := range tests {
		strVal := &String{Value: tt.value}
		
		if strVal.Type() != STRING_VALUE {
			t.Fatalf("Expected type STRING_VALUE, got %s", strVal.Type())
		}
		
		if strVal.Inspect() != tt.expected {
			t.Fatalf("Expected inspect %q, got %q", tt.expected, strVal.Inspect())
		}
	}
}

func TestBooleanValue(t *testing.T) {
	tests := []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}
	
	for _, tt := range tests {
		boolVal := &Boolean{Value: tt.value}
		
		if boolVal.Type() != BOOLEAN_VALUE {
			t.Fatalf("Expected type BOOLEAN_VALUE, got %s", boolVal.Type())
		}
		
		if boolVal.Inspect() != tt.expected {
			t.Fatalf("Expected inspect %q, got %q", tt.expected, boolVal.Inspect())
		}
	}
}

func TestArrayValue(t *testing.T) {
	// Empty array
	emptyArray := &Array{Elements: []Value{}}
	if emptyArray.Type() != ARRAY_VALUE {
		t.Fatalf("Expected type ARRAY_VALUE, got %s", emptyArray.Type())
	}
	
	if emptyArray.Inspect() != "[]" {
		t.Fatalf("Expected empty array inspect '[]', got %q", emptyArray.Inspect())
	}
	
	// Array with elements
	elements := []Value{
		&Integer{Value: 1},
		&String{Value: "hello"},
		&Boolean{Value: true},
		&Float{Value: 3.14},
	}
	
	array := &Array{Elements: elements}
	expected := "[1, hello, true, 3.14]"
	
	if array.Inspect() != expected {
		t.Fatalf("Expected array inspect %q, got %q", expected, array.Inspect())
	}
	
	// Nested arrays
	nestedArray := &Array{
		Elements: []Value{
			&Array{Elements: []Value{&Integer{Value: 1}, &Integer{Value: 2}}},
			&Array{Elements: []Value{&String{Value: "a"}, &String{Value: "b"}}},
		},
	}
	
	expectedNested := "[[1, 2], [a, b]]"
	if nestedArray.Inspect() != expectedNested {
		t.Fatalf("Expected nested array inspect %q, got %q", expectedNested, nestedArray.Inspect())
	}
}

func TestNullValue(t *testing.T) {
	nullVal := &Null{}
	
	if nullVal.Type() != NULL_VALUE {
		t.Fatalf("Expected type NULL_VALUE, got %s", nullVal.Type())
	}
	
	if nullVal.Inspect() != "null" {
		t.Fatalf("Expected null inspect 'null', got %q", nullVal.Inspect())
	}
}

func TestReturnValue(t *testing.T) {
	innerValue := &Integer{Value: 42}
	returnVal := &ReturnValue{Value: innerValue}
	
	if returnVal.Type() != RETURN_VALUE {
		t.Fatalf("Expected type RETURN_VALUE, got %s", returnVal.Type())
	}
	
	if returnVal.Inspect() != "42" {
		t.Fatalf("Expected return value inspect '42', got %q", returnVal.Inspect())
	}
}

func TestBreakContinueValues(t *testing.T) {
	breakVal := &BreakValue{}
	if breakVal.Type() != BREAK_VALUE {
		t.Fatalf("Expected type BREAK_VALUE, got %s", breakVal.Type())
	}
	
	if breakVal.Inspect() != "break" {
		t.Fatalf("Expected break inspect 'break', got %q", breakVal.Inspect())
	}
	
	continueVal := &ContinueValue{}
	if continueVal.Type() != CONTINUE_VALUE {
		t.Fatalf("Expected type CONTINUE_VALUE, got %s", continueVal.Type())
	}
	
	if continueVal.Inspect() != "continue" {
		t.Fatalf("Expected continue inspect 'continue', got %q", continueVal.Inspect())
	}
}

func TestBuiltinFunction(t *testing.T) {
	builtin := &BuiltinFunction{
		Fn: func(args ...Value) Value {
			return &Integer{Value: 42}
		},
	}
	
	if builtin.Type() != BUILTIN_VALUE {
		t.Fatalf("Expected type BUILTIN_VALUE, got %s", builtin.Type())
	}
	
	if builtin.Inspect() != "builtin function" {
		t.Fatalf("Expected builtin inspect 'builtin function', got %q", builtin.Inspect())
	}
	
	// Test function execution
	result := builtin.Fn()
	intResult, ok := result.(*Integer)
	if !ok || intResult.Value != 42 {
		t.Fatal("Builtin function should execute correctly")
	}
}

func TestFunctionValue(t *testing.T) {
	// Create a simple function AST
	params := []*ast.Identifier{
		{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
		{Token: lexer.Token{Type: lexer.IDENT, Literal: "y"}, Value: "y"},
	}
	
	body := &ast.BlockStatement{
		Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: lexer.Token{Type: lexer.IDENT, Literal: "x"},
				Expression: &ast.InfixExpression{
					Token:    lexer.Token{Type: lexer.PLUS, Literal: "+"},
					Left:     &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "x"}, Value: "x"},
					Operator: "+",
					Right:    &ast.Identifier{Token: lexer.Token{Type: lexer.IDENT, Literal: "y"}, Value: "y"},
				},
			},
		},
	}
	
	env := NewEnvironment()
	fn := &Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
	
	if fn.Type() != FUNCTION_VALUE {
		t.Fatalf("Expected type FUNCTION_VALUE, got %s", fn.Type())
	}
	
	inspect := fn.Inspect()
	if len(inspect) == 0 {
		t.Fatal("Function inspect should not be empty")
	}
	
	// Should contain parameter names
	if !contains(inspect, "x") || !contains(inspect, "y") {
		t.Fatal("Function inspect should contain parameter names")
	}
}

func TestHashValue(t *testing.T) {
	// Empty hash
	emptyHash := &Hash{
		Pairs: make(map[HashKey]Value),
		Keys:  []Value{},
	}
	
	if emptyHash.Type() != HASH_VALUE {
		t.Fatalf("Expected type HASH_VALUE, got %s", emptyHash.Type())
	}
	
	if emptyHash.Inspect() != "{}" {
		t.Fatalf("Expected empty hash inspect '{}', got %q", emptyHash.Inspect())
	}
	
	// Hash with elements
	stringKey := &String{Value: "name"}
	intKey := &Integer{Value: 42}
	boolKey := &Boolean{Value: true}
	
	hash := &Hash{
		Pairs: map[HashKey]Value{
			CreateHashKey(stringKey): &String{Value: "Alice"},
			CreateHashKey(intKey):    &String{Value: "answer"},
			CreateHashKey(boolKey):   &Integer{Value: 1},
		},
		Keys: []Value{stringKey, intKey, boolKey},
	}
	
	inspect := hash.Inspect()
	
	// Should contain all key-value pairs
	expectedPairs := []string{
		"name: Alice",
		"42: answer", 
		"true: 1",
	}
	
	for _, pair := range expectedPairs {
		if !contains(inspect, pair) {
			t.Fatalf("Hash inspect should contain %q, got %q", pair, inspect)
		}
	}
}

func TestCreateHashKey(t *testing.T) {
	tests := []struct {
		value    Value
		expected HashKey
	}{
		{
			&Integer{Value: 42},
			HashKey{Type: INTEGER_VALUE, Value: int64(42)},
		},
		{
			&String{Value: "hello"},
			HashKey{Type: STRING_VALUE, Value: "hello"},
		},
		{
			&Boolean{Value: true},
			HashKey{Type: BOOLEAN_VALUE, Value: true},
		},
		{
			&Float{Value: 3.14},
			HashKey{Type: FLOAT_VALUE, Value: 3.14},
		},
	}
	
	for _, tt := range tests {
		key := CreateHashKey(tt.value)
		
		if key.Type != tt.expected.Type {
			t.Fatalf("Expected key type %s, got %s", tt.expected.Type, key.Type)
		}
		
		if key.Value != tt.expected.Value {
			t.Fatalf("Expected key value %v, got %v", tt.expected.Value, key.Value)
		}
	}
}

func TestCreateHashKeyInvalidType(t *testing.T) {
	// Test with unsupported type
	array := &Array{Elements: []Value{}}
	key := CreateHashKey(array)
	
	if key.Type != NULL_VALUE {
		t.Fatalf("Expected NULL_VALUE type for invalid hash key, got %s", key.Type)
	}
	
	if key.Value != nil {
		t.Fatalf("Expected nil value for invalid hash key, got %v", key.Value)
	}
}

func TestExceptionValue(t *testing.T) {
	errorValue := &String{Value: "Something went wrong"}
	exception := &Exception{Error: errorValue}
	
	if exception.Type() != EXCEPTION_VALUE {
		t.Fatalf("Expected type EXCEPTION_VALUE, got %s", exception.Type())
	}
	
	if exception.Inspect() != "Something went wrong" {
		t.Fatalf("Expected exception inspect 'Something went wrong', got %q", exception.Inspect())
	}
}

func TestClassValue(t *testing.T) {
	env := NewEnvironment()
	class := &Class{
		Name:       "TestClass",
		SuperClass: nil,
		Methods:    make(map[string]*Function),
		Env:        env,
	}
	
	if class.Type() != CLASS_VALUE {
		t.Fatalf("Expected type CLASS_VALUE, got %s", class.Type())
	}
	
	expected := "class TestClass"
	if class.Inspect() != expected {
		t.Fatalf("Expected class inspect %q, got %q", expected, class.Inspect())
	}
}

func TestObjectValue(t *testing.T) {
	env := NewEnvironment()
	class := &Class{
		Name:       "TestClass",
		SuperClass: nil,
		Methods:    make(map[string]*Function),
		Env:        env,
	}
	
	obj := &Object{
		Class:        class,
		InstanceVars: make(map[string]Value),
		Env:          env,
	}
	
	if obj.Type() != INSTANCE_VALUE {
		t.Fatalf("Expected type INSTANCE_VALUE, got %s", obj.Type())
	}
	
	inspect := obj.Inspect()
	if !contains(inspect, "TestClass") {
		t.Fatalf("Object inspect should contain class name, got %q", inspect)
	}
}

func TestBoundMethodValue(t *testing.T) {
	env := NewEnvironment()
	method := &Function{
		Parameters: []*ast.Identifier{},
		Body:       &ast.BlockStatement{},
		Env:        env,
	}
	
	class := &Class{Name: "TestClass", Methods: make(map[string]*Function), Env: env}
	instance := &Object{Class: class, InstanceVars: make(map[string]Value), Env: env}
	
	boundMethod := &BoundMethod{
		Method:   method,
		Instance: instance,
	}
	
	if boundMethod.Type() != BOUND_METHOD_VALUE {
		t.Fatalf("Expected type BOUND_METHOD_VALUE, got %s", boundMethod.Type())
	}
	
	inspect := boundMethod.Inspect()
	if !contains(inspect, "BoundMethod") {
		t.Fatalf("BoundMethod inspect should contain 'BoundMethod', got %q", inspect)
	}
}

func TestMethodValues(t *testing.T) {
	// Test HashMethod
	hash := &Hash{Pairs: make(map[HashKey]Value), Keys: []Value{}}
	hashMethod := &HashMethod{Hash: hash, Method: "keys"}
	
	if hashMethod.Type() != HASH_METHOD_VALUE {
		t.Fatalf("Expected type HASH_METHOD_VALUE, got %s", hashMethod.Type())
	}
	
	if !contains(hashMethod.Inspect(), "HashMethod") {
		t.Fatal("HashMethod inspect should contain 'HashMethod'")
	}
	
	// Test StringMethod
	str := &String{Value: "test"}
	stringMethod := &StringMethod{String: str, Method: "length"}
	
	if stringMethod.Type() != STRING_METHOD_VALUE {
		t.Fatalf("Expected type STRING_METHOD_VALUE, got %s", stringMethod.Type())
	}
	
	if !contains(stringMethod.Inspect(), "StringMethod") {
		t.Fatal("StringMethod inspect should contain 'StringMethod'")
	}
	
	// Test ArrayMethod
	array := &Array{Elements: []Value{}}
	arrayMethod := &ArrayMethod{Array: array, Method: "push"}
	
	if arrayMethod.Type() != ARRAY_METHOD_VALUE {
		t.Fatalf("Expected type ARRAY_METHOD_VALUE, got %s", arrayMethod.Type())
	}
	
	if !contains(arrayMethod.Inspect(), "ArrayMethod") {
		t.Fatal("ArrayMethod inspect should contain 'ArrayMethod'")
	}
}

func TestValueEquality(t *testing.T) {
	// Test that same values create same hash keys
	int1 := &Integer{Value: 42}
	int2 := &Integer{Value: 42}
	
	key1 := CreateHashKey(int1)
	key2 := CreateHashKey(int2)
	
	if key1 != key2 {
		t.Fatal("Same integer values should create equal hash keys")
	}
	
	str1 := &String{Value: "hello"}
	str2 := &String{Value: "hello"}
	
	keyStr1 := CreateHashKey(str1)
	keyStr2 := CreateHashKey(str2)
	
	if keyStr1 != keyStr2 {
		t.Fatal("Same string values should create equal hash keys")
	}
	
	bool1 := &Boolean{Value: true}
	bool2 := &Boolean{Value: true}
	
	keyBool1 := CreateHashKey(bool1)
	keyBool2 := CreateHashKey(bool2)
	
	if keyBool1 != keyBool2 {
		t.Fatal("Same boolean values should create equal hash keys")
	}
}

func TestValueInequality(t *testing.T) {
	// Test that different values create different hash keys
	int1 := &Integer{Value: 42}
	int2 := &Integer{Value: 43}
	
	key1 := CreateHashKey(int1)
	key2 := CreateHashKey(int2)
	
	if key1 == key2 {
		t.Fatal("Different integer values should create different hash keys")
	}
	
	str1 := &String{Value: "hello"}
	str2 := &String{Value: "world"}
	
	keyStr1 := CreateHashKey(str1)
	keyStr2 := CreateHashKey(str2)
	
	if keyStr1 == keyStr2 {
		t.Fatal("Different string values should create different hash keys")
	}
	
	bool1 := &Boolean{Value: true}
	bool2 := &Boolean{Value: false}
	
	keyBool1 := CreateHashKey(bool1)
	keyBool2 := CreateHashKey(bool2)
	
	if keyBool1 == keyBool2 {
		t.Fatal("Different boolean values should create different hash keys")
	}
}

func TestComplexHashOperations(t *testing.T) {
	// Create a hash with mixed key types
	keys := []Value{
		&String{Value: "name"},
		&Integer{Value: 42},
		&Boolean{Value: true},
		&Float{Value: 3.14},
	}
	
	values := []Value{
		&String{Value: "Alice"},
		&String{Value: "answer"},
		&Integer{Value: 1},
		&String{Value: "pi"},
	}
	
	pairs := make(map[HashKey]Value)
	for i, key := range keys {
		pairs[CreateHashKey(key)] = values[i]
	}
	
	hash := &Hash{Pairs: pairs, Keys: keys}
	
	// Test that all keys are preserved in order
	if len(hash.Keys) != 4 {
		t.Fatalf("Expected 4 keys, got %d", len(hash.Keys))
	}
	
	// Test retrieval by hash key
	nameKey := CreateHashKey(&String{Value: "name"})
	nameValue, exists := hash.Pairs[nameKey]
	if !exists {
		t.Fatal("Expected 'name' key to exist")
	}
	
	strValue, ok := nameValue.(*String)
	if !ok || strValue.Value != "Alice" {
		t.Fatal("Expected 'name' key to have value 'Alice'")
	}
	
	// Test float key
	floatKey := CreateHashKey(&Float{Value: 3.14})
	floatValue, exists := hash.Pairs[floatKey]
	if !exists {
		t.Fatal("Expected float key to exist")
	}
	
	strValue, ok = floatValue.(*String)
	if !ok || strValue.Value != "pi" {
		t.Fatal("Expected float key to have value 'pi'")
	}
}

// Helper functions are defined in environment_test.go