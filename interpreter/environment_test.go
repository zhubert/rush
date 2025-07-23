package interpreter

import (
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()
	
	// Test that environment is initialized properly
	if env.store == nil {
		t.Fatal("Environment store should not be nil")
	}
	
	if env.outer != nil {
		t.Fatal("New environment should not have outer environment")
	}
	
	if env.moduleResolver == nil {
		t.Fatal("Environment should have module resolver")
	}
	
	if env.currentDir != "." {
		t.Fatalf("Expected current directory '.', got %q", env.currentDir)
	}
	
	if env.exports == nil {
		t.Fatal("Environment exports should not be nil")
	}
	
	if env.callStack == nil {
		t.Fatal("Environment call stack should not be nil")
	}
	
	// Test that built-ins are loaded
	if _, ok := env.Get("len"); !ok {
		t.Fatal("Built-in function 'len' should be available")
	}
	
	if _, ok := env.Get("print"); !ok {
		t.Fatal("Built-in function 'print' should be available")
	}
}

func TestNewEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.SetCurrentDir("/test/dir")
	outer.PushCall("testFunc", 10, 5)
	
	inner := NewEnclosedEnvironment(outer)
	
	if inner.outer != outer {
		t.Fatal("Enclosed environment should reference outer environment")
	}
	
	if inner.GetCurrentDir() != "/test/dir" {
		t.Fatalf("Expected current directory '/test/dir', got %q", inner.GetCurrentDir())
	}
	
	// Should inherit call stack
	stack := inner.GetCallStack()
	if len(stack) != 1 {
		t.Fatalf("Expected call stack length 1, got %d", len(stack))
	}
	
	if stack[0].FunctionName != "testFunc" {
		t.Fatalf("Expected function name 'testFunc', got %q", stack[0].FunctionName)
	}
	
	// Should still have built-ins
	if _, ok := inner.Get("len"); !ok {
		t.Fatal("Built-in function 'len' should be available in enclosed environment")
	}
}

func TestBasicVariableOperations(t *testing.T) {
	env := NewEnvironment()
	
	// Test Set and Get
	testValue := &Integer{Value: 42}
	env.Set("x", testValue)
	
	retrieved, ok := env.Get("x")
	if !ok {
		t.Fatal("Variable 'x' should exist")
	}
	
	if retrieved != testValue {
		t.Fatal("Retrieved value should be the same as set value")
	}
	
	// Test non-existent variable
	_, ok = env.Get("nonexistent")
	if ok {
		t.Fatal("Non-existent variable should return false")
	}
}

func TestEnvironmentVariableScoping(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("outerVar", &Integer{Value: 10})
	outer.Set("sharedVar", &String{Value: "outer"})
	
	inner := NewEnclosedEnvironment(outer)
	inner.Set("innerVar", &Integer{Value: 20})
	
	// Test outer variable access from inner
	val, ok := inner.Get("outerVar")
	if !ok {
		t.Fatal("Inner environment should access outer variables")
	}
	
	intVal, ok := val.(*Integer)
	if !ok || intVal.Value != 10 {
		t.Fatal("Outer variable should have correct value")
	}
	
	// Test inner variable not accessible from outer
	_, ok = outer.Get("innerVar")
	if ok {
		t.Fatal("Outer environment should not access inner variables")
	}
	
	// Test variable update propagation (not shadowing)
	// When a variable exists in outer scope, Set updates it there
	inner.Set("sharedVar", &String{Value: "updated"})
	
	// Both inner and outer should see the updated value
	val, ok = inner.Get("sharedVar")
	if !ok {
		t.Fatal("Updated variable should be accessible from inner")
	}
	
	strVal, ok := val.(*String)
	if !ok || strVal.Value != "updated" {
		t.Fatal("Variable should have updated value")
	}
	
	val, ok = outer.Get("sharedVar")
	if !ok {
		t.Fatal("Variable should still exist in outer scope")
	}
	
	strVal, ok = val.(*String)
	if !ok || strVal.Value != "updated" {
		t.Fatal("Outer scope should see updated value")
	}
}

func TestVariableUpdatePropagation(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("count", &Integer{Value: 0})
	
	inner := NewEnclosedEnvironment(outer)
	
	// Update variable from inner environment
	inner.Set("count", &Integer{Value: 5})
	
	// Check that outer environment sees the update
	val, ok := outer.Get("count")
	if !ok {
		t.Fatal("Variable should exist in outer environment")
	}
	
	intVal, ok := val.(*Integer)
	if !ok || intVal.Value != 5 {
		t.Fatalf("Expected value 5, got %d", intVal.Value)
	}
	
	// Test new variable creation in inner
	inner.Set("newVar", &String{Value: "new"})
	
	// Should exist in inner
	val, ok = inner.Get("newVar")
	if !ok {
		t.Fatal("New variable should exist in inner environment")
	}
	
	// Should not exist in outer
	_, ok = outer.Get("newVar")
	if ok {
		t.Fatal("New variable should not be created in outer environment")
	}
}

func TestSetLocal(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Integer{Value: 10})
	
	inner := NewEnclosedEnvironment(outer)
	
	// Use SetLocal to shadow variable without updating outer
	inner.SetLocal("x", &Integer{Value: 20})
	
	// Inner should see the local value
	val, ok := inner.Get("x")
	if !ok {
		t.Fatal("Variable should exist in inner environment")
	}
	
	intVal, ok := val.(*Integer)
	if !ok || intVal.Value != 20 {
		t.Fatalf("Expected local value 20, got %d", intVal.Value)
	}
	
	// Outer should still have original value
	val, ok = outer.Get("x")
	if !ok {
		t.Fatal("Variable should still exist in outer environment")
	}
	
	intVal, ok = val.(*Integer)
	if !ok || intVal.Value != 10 {
		t.Fatalf("Expected outer value 10, got %d", intVal.Value)
	}
}

func TestDeepNesting(t *testing.T) {
	// Create nested environments
	level1 := NewEnvironment()
	level1.Set("var1", &Integer{Value: 1})
	
	level2 := NewEnclosedEnvironment(level1)
	level2.Set("var2", &Integer{Value: 2})
	
	level3 := NewEnclosedEnvironment(level2)
	level3.Set("var3", &Integer{Value: 3})
	
	level4 := NewEnclosedEnvironment(level3)
	level4.Set("var4", &Integer{Value: 4})
	
	// Test access from deepest level
	testCases := []struct {
		varName string
		expected int64
	}{
		{"var1", 1},
		{"var2", 2},
		{"var3", 3},
		{"var4", 4},
	}
	
	for _, tc := range testCases {
		val, ok := level4.Get(tc.varName)
		if !ok {
			t.Fatalf("Variable %q should be accessible from deep level", tc.varName)
		}
		
		intVal, ok := val.(*Integer)
		if !ok || intVal.Value != tc.expected {
			t.Fatalf("Expected %q = %d, got %d", tc.varName, tc.expected, intVal.Value)
		}
	}
	
	// Test that shallow levels cannot access deep variables
	_, ok1 := level1.Get("var4")
	if ok1 {
		t.Fatal("Level 1 should not access level 4 variables")
	}
	
	_, ok2 := level2.Get("var4")
	if ok2 {
		t.Fatal("Level 2 should not access level 4 variables")
	}
}

func TestModuleEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("globalVar", &Integer{Value: 100})
	
	moduleEnv := NewModuleEnvironment(outer)
	
	// Should have access to outer variables
	val, ok := moduleEnv.Get("globalVar")
	if !ok {
		t.Fatal("Module environment should access global variables")
	}
	
	intVal, ok := val.(*Integer)
	if !ok || intVal.Value != 100 {
		t.Fatal("Global variable should have correct value")
	}
	
	// Should have fresh exports map
	if len(moduleEnv.GetExports()) != 0 {
		t.Fatal("Module environment should start with empty exports")
	}
	
	// Test exports functionality
	moduleEnv.AddExport("testExport", &String{Value: "exported"})
	
	exports := moduleEnv.GetExports()
	if len(exports) != 1 {
		t.Fatalf("Expected 1 export, got %d", len(exports))
	}
	
	exportVal, ok := exports["testExport"]
	if !ok {
		t.Fatal("Export 'testExport' should exist")
	}
	
	strVal, ok := exportVal.(*String)
	if !ok || strVal.Value != "exported" {
		t.Fatal("Export should have correct value")
	}
}

func TestCallStack(t *testing.T) {
	env := NewEnvironment()
	
	// Initially empty
	stack := env.GetCallStack()
	if len(stack) != 0 {
		t.Fatal("Call stack should initially be empty")
	}
	
	// Push calls
	env.PushCall("main", 1, 1)
	env.PushCall("foo", 10, 5)
	env.PushCall("bar", 20, 10)
	
	stack = env.GetCallStack()
	if len(stack) != 3 {
		t.Fatalf("Expected 3 call frames, got %d", len(stack))
	}
	
	// Check frames
	expected := []struct {
		name   string
		line   int
		column int
	}{
		{"main", 1, 1},
		{"foo", 10, 5},
		{"bar", 20, 10},
	}
	
	for i, frame := range expected {
		if stack[i].FunctionName != frame.name {
			t.Fatalf("Frame %d: expected function %q, got %q", i, frame.name, stack[i].FunctionName)
		}
		if stack[i].Line != frame.line {
			t.Fatalf("Frame %d: expected line %d, got %d", i, frame.line, stack[i].Line)
		}
		if stack[i].Column != frame.column {
			t.Fatalf("Frame %d: expected column %d, got %d", i, frame.column, stack[i].Column)
		}
	}
	
	// Pop calls
	env.PopCall()
	stack = env.GetCallStack()
	if len(stack) != 2 {
		t.Fatalf("Expected 2 call frames after pop, got %d", len(stack))
	}
	
	if stack[len(stack)-1].FunctionName != "foo" {
		t.Fatalf("Expected top frame to be 'foo', got %q", stack[len(stack)-1].FunctionName)
	}
	
	// Pop all
	env.PopCall()
	env.PopCall()
	stack = env.GetCallStack()
	if len(stack) != 0 {
		t.Fatal("Call stack should be empty after popping all")
	}
	
	// Pop empty stack should not crash
	env.PopCall()
	stack = env.GetCallStack()
	if len(stack) != 0 {
		t.Fatal("Popping empty stack should not add frames")
	}
}

func TestStackTrace(t *testing.T) {
	env := NewEnvironment()
	
	// Empty stack trace
	trace := env.GetStackTrace()
	if trace != "" {
		t.Fatal("Empty call stack should produce empty trace")
	}
	
	// Add some calls
	env.PushCall("main", 1, 1)
	env.PushCall("calculateSum", 15, 8)
	env.PushCall("helper", 30, 12)
	
	trace = env.GetStackTrace()
	if trace == "" {
		t.Fatal("Non-empty call stack should produce trace")
	}
	
	// Should contain function names in reverse order (most recent first)
	expectedOrder := []string{"helper", "calculateSum", "main"}
	for _, funcName := range expectedOrder {
		if !contains(trace, funcName) {
			t.Fatalf("Stack trace should contain function %q", funcName)
		}
	}
	
	// Should contain line numbers
	if !contains(trace, "line 30:12") {
		t.Fatal("Stack trace should contain line numbers")
	}
	
	if !contains(trace, "line 15:8") {
		t.Fatal("Stack trace should contain line numbers")
	}
}

func TestCallStackInheritance(t *testing.T) {
	outer := NewEnvironment()
	outer.PushCall("outerFunc", 5, 3)
	outer.PushCall("middleFunc", 10, 7)
	
	inner := NewEnclosedEnvironment(outer)
	
	// Should inherit call stack
	stack := inner.GetCallStack()
	if len(stack) != 2 {
		t.Fatalf("Expected inherited call stack length 2, got %d", len(stack))
	}
	
	if stack[0].FunctionName != "outerFunc" {
		t.Fatalf("Expected first frame 'outerFunc', got %q", stack[0].FunctionName)
	}
	
	if stack[1].FunctionName != "middleFunc" {
		t.Fatalf("Expected second frame 'middleFunc', got %q", stack[1].FunctionName)
	}
	
	// Adding to inner should not affect outer
	inner.PushCall("innerFunc", 15, 2)
	
	outerStack := outer.GetCallStack()
	if len(outerStack) != 2 {
		t.Fatalf("Outer call stack should remain length 2, got %d", len(outerStack))
	}
	
	innerStack := inner.GetCallStack()
	if len(innerStack) != 3 {
		t.Fatalf("Inner call stack should be length 3, got %d", len(innerStack))
	}
	
	if innerStack[2].FunctionName != "innerFunc" {
		t.Fatalf("Expected last frame 'innerFunc', got %q", innerStack[2].FunctionName)
	}
}

func TestMultipleVariableTypes(t *testing.T) {
	env := NewEnvironment()
	
	// Test various value types
	values := map[string]Value{
		"intVar":    &Integer{Value: 42},
		"floatVar":  &Float{Value: 3.14},
		"stringVar": &String{Value: "hello"},
		"boolVar":   &Boolean{Value: true},
		"arrayVar":  &Array{Elements: []Value{&Integer{Value: 1}, &Integer{Value: 2}}},
	}
	
	// Set all values
	for name, val := range values {
		env.Set(name, val)
	}
	
	// Retrieve and verify all values
	for name, expectedVal := range values {
		retrievedVal, ok := env.Get(name)
		if !ok {
			t.Fatalf("Variable %q should exist", name)
		}
		
		if retrievedVal != expectedVal {
			t.Fatalf("Variable %q should have correct value", name)
		}
	}
}

func TestScopeIsolation(t *testing.T) {
	// Create two separate environment chains
	chain1_level1 := NewEnvironment()
	chain1_level2 := NewEnclosedEnvironment(chain1_level1)
	
	chain2_level1 := NewEnvironment()
	chain2_level2 := NewEnclosedEnvironment(chain2_level1)
	
	// Set variables in each chain
	chain1_level1.Set("shared", &String{Value: "chain1"})
	chain2_level1.Set("shared", &String{Value: "chain2"})
	
	// Each chain should have its own values
	val1, ok1 := chain1_level2.Get("shared")
	val2, ok2 := chain2_level2.Get("shared")
	
	if !ok1 || !ok2 {
		t.Fatal("Both chains should have the shared variable")
	}
	
	str1, ok1 := val1.(*String)
	if !ok1 || str1.Value != "chain1" {
		t.Fatal("Chain 1 should have its own value")
	}
	
	str2, ok2 := val2.(*String)
	if !ok2 || str2.Value != "chain2" {
		t.Fatal("Chain 2 should have its own value")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    s[:len(substr)] == substr || 
		    s[len(s)-len(substr):] == substr ||
		    containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}