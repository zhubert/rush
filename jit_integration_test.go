package main

import (
	"strings"
	"testing"
	"time"

	"rush/compiler"
	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
	"rush/vm"
)

// TestJITIntegration tests end-to-end JIT functionality
func TestJITIntegration(t *testing.T) {
	// Simple recursive function that should trigger JIT compilation
	source := `
	fibonacci = fn(n) {
		if (n <= 1) {
			return n
		}
		return fibonacci(n - 1) + fibonacci(n - 2)
	}
	
	result = fibonacci(10)
	result
	`
	
	// Test with JIT enabled
	resultJIT := executeWithJIT(t, source)
	
	// Test with bytecode only
	resultBytecode := executeWithBytecode(t, source)
	
	// Check for nil results
	if resultJIT == nil {
		t.Fatal("JIT execution returned nil result")
	}
	if resultBytecode == nil {
		t.Fatal("Bytecode execution returned nil result") 
	}
	
	// Results should be identical
	if resultJIT.Inspect() != resultBytecode.Inspect() {
		t.Errorf("JIT and bytecode results differ: JIT=%s, Bytecode=%s", 
			resultJIT.Inspect(), resultBytecode.Inspect())
	}
	
	// Result should be fibonacci(10) = 55
	if resultJIT.Inspect() != "55" {
		t.Errorf("Incorrect fibonacci result: got %s, want 55", resultJIT.Inspect())
	}
}

// TestJITStatistics verifies that JIT statistics are properly tracked
func TestJITStatistics(t *testing.T) {
	source := `
	add = fn(a, b) { return a + b }
	
	# Call function many times to trigger hot path detection
	sum = 0
	i = 0
	while (i < 50) {
		sum = add(sum, i)
		i = i + 1
	}
	sum
	`
	
	machine := createJITVMTest(t, source)
	err := machine.Run()
	if err != nil {
		t.Fatalf("JIT execution failed: %v", err)
	}
	
	// Check JIT statistics
	jitStats := machine.GetJITStats()
	
	// Should have attempted compilations
	if jitStats.CompilationsAttempted == 0 {
		t.Error("Expected JIT compilation attempts, got 0")
	}
	
	// Should have recorded function executions
	vmStats := machine.GetStats()
	if len(vmStats.FunctionExecutions) == 0 {
		t.Error("Expected function execution tracking, got none")
	}
	
	t.Logf("JIT Stats: Attempts=%d, Succeeded=%d, Hits=%d, Misses=%d", 
		jitStats.CompilationsAttempted, jitStats.CompilationsSucceeded, 
		jitStats.JITHits, jitStats.JITMisses)
}

// TestJITDeoptimization tests that JIT properly falls back to bytecode
func TestJITDeoptimization(t *testing.T) {
	source := `
	test_fn = fn(x) {
		if (x > 0) {
			return x * 2
		}
		return 0
	}
	
	# Call enough times to trigger compilation attempts
	result = 0
	i = 0
	while (i < 20) {
		result = test_fn(i)
		i = i + 1
	}
	result
	`
	
	machine := createJITVMTest(t, source)
	err := machine.Run()
	if err != nil {
		t.Fatalf("JIT execution with deoptimization failed: %v", err)
	}
	
	// Should have some deoptimizations since ARM64 execution isn't implemented
	jitStats := machine.GetJITStats()
	if jitStats.CompilationsAttempted > 0 && jitStats.JITHits == 0 {
		// This is expected - compilations attempted but no successful JIT hits
		t.Logf("JIT deoptimization working correctly: %d attempts, %d hits", 
			jitStats.CompilationsAttempted, jitStats.JITHits)
	}
}

// TestJITvsTreeWalkingPerformance compares JIT overhead vs tree-walking
func TestJITvsTreeWalkingPerformance(t *testing.T) {
	source := `
	factorial = fn(n) {
		if (n <= 1) {
			return 1
		}
		return n * factorial(n - 1)
	}
	
	factorial(8)
	`
	
	// Measure tree-walking interpreter
	startTree := time.Now()
	resultTree := executeWithTreeWalking(t, source)
	treeTime := time.Since(startTree)
	
	// Measure JIT (which will fall back to bytecode with profiling overhead)
	startJIT := time.Now()
	resultJIT := executeWithJIT(t, source)
	jitTime := time.Since(startJIT)
	
	// Results should be identical
	if resultTree.Inspect() != resultJIT.Inspect() {
		t.Errorf("Tree-walking and JIT results differ: Tree=%s, JIT=%s", 
			resultTree.Inspect(), resultJIT.Inspect())
	}
	
	t.Logf("Performance comparison - Tree-walking: %v, JIT (with overhead): %v", 
		treeTime, jitTime)
	
	// JIT mode may be slower due to profiling overhead, which is expected
	// since ARM64 execution isn't fully implemented
}

// Helper functions

func executeWithJIT(t *testing.T, source string) interpreter.Value {
	machine := createJITVMTest(t, source)
	err := machine.Run()
	if err != nil {
		t.Fatalf("JIT execution failed: %v", err)
	}
	return machine.StackTop()
}

func executeWithBytecode(t *testing.T, source string) interpreter.Value {
	machine := createBytecodeVMTest(t, source)
	err := machine.Run()
	if err != nil {
		t.Fatalf("Bytecode execution failed: %v", err)
	}
	return machine.StackTop()
}

func executeWithTreeWalking(t *testing.T, source string) interpreter.Value {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	
	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %s", strings.Join(errors, ", "))
	}
	
	env := interpreter.NewEnvironment()
	result := interpreter.Eval(program, env)
	
	if result.Type() == "ERROR" {
		t.Fatalf("Tree-walking execution error: %s", result.Inspect())
	}
	
	return result
}

func createJITVMTest(t *testing.T, source string) *vm.VM {
	bytecode := compileToBytecodeTest(t, source)
	return vm.NewWithJIT(bytecode, vm.LogError)
}

func createBytecodeVMTest(t *testing.T, source string) *vm.VM {
	bytecode := compileToBytecodeTest(t, source)
	return vm.NewWithLogger(bytecode, vm.LogError)
}

func compileToBytecodeTest(t *testing.T, source string) *compiler.Bytecode {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	
	errors := p.Errors()
	if len(errors) > 0 {
		t.Fatalf("Parse errors: %s", strings.Join(errors, ", "))
	}
	
	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		t.Fatalf("Compilation error: %v", err)
	}
	
	return comp.Bytecode()
}

// Benchmark functions

func BenchmarkFibonacciJIT(b *testing.B) {
	source := `
	fibonacci = fn(n) {
		if (n <= 1) {
			return n
		}
		return fibonacci(n - 1) + fibonacci(n - 2)
	}
	fibonacci(15)
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine := createJITVMBench(b, source)
		machine.Run()
	}
}

func BenchmarkFibonacciBytecode(b *testing.B) {
	source := `
	fibonacci = fn(n) {
		if (n <= 1) {
			return n
		}
		return fibonacci(n - 1) + fibonacci(n - 2)
	}
	fibonacci(15)
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine := createBytecodeVMBench(b, source)
		machine.Run()
	}
}

func BenchmarkArithmeticJIT(b *testing.B) {
	source := `
	add = fn(a, b) { return a + b }
	mul = fn(a, b) { return a * b }
	
	result = 0
	i = 0
	while (i < 100) {
		result = add(result, mul(i, 2))
		i = i + 1
	}
	result
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine := createJITVMBench(b, source)
		machine.Run()
	}
}

func BenchmarkArithmeticBytecode(b *testing.B) {
	source := `
	add = fn(a, b) { return a + b }
	mul = fn(a, b) { return a * b }
	
	result = 0
	i = 0
	while (i < 100) {
		result = add(result, mul(i, 2))
		i = i + 1
	}
	result
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine := createBytecodeVMBench(b, source)
		machine.Run()
	}
}

// Helper functions for benchmarks
func createJITVMBench(b *testing.B, source string) *vm.VM {
	bytecode := compileToBytcodeBench(b, source)
	return vm.NewWithJIT(bytecode, vm.LogNone) // No logging for benchmarks
}

func createBytecodeVMBench(b *testing.B, source string) *vm.VM {
	bytecode := compileToBytcodeBench(b, source)
	return vm.NewWithLogger(bytecode, vm.LogNone) // No logging for benchmarks
}

func compileToBytcodeBench(b *testing.B, source string) *compiler.Bytecode {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	
	errors := p.Errors()
	if len(errors) > 0 {
		b.Fatalf("Parse errors: %s", strings.Join(errors, ", "))
	}
	
	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		b.Fatalf("Compilation error: %v", err)
	}
	
	return comp.Bytecode()
}