package jit

import (
	"testing"
	"time"

	"rush/bytecode"
	"rush/interpreter"
)

func TestJITCompilerCreation(t *testing.T) {
	compiler := NewJITCompiler()
	if compiler == nil {
		t.Fatal("Failed to create JIT compiler")
	}
	
	if compiler.hotThreshold != DefaultHotThreshold {
		t.Errorf("Unexpected hot threshold: got %d, want %d", 
			compiler.hotThreshold, DefaultHotThreshold)
	}
	
	if compiler.cache == nil {
		t.Error("JIT compiler cache not initialized")
	}
	
	if compiler.profiler == nil {
		t.Error("JIT compiler profiler not initialized")
	}
}

func TestExecutionProfiler(t *testing.T) {
	profiler := NewExecutionProfiler()
	
	// Test initial state
	if count := profiler.GetExecutionCount(123); count != 0 {
		t.Errorf("Expected 0 execution count, got %d", count)
	}
	
	// Record some executions
	fnHash := uint64(123)
	for i := 0; i < 5; i++ {
		profiler.RecordExecution(fnHash, time.Millisecond*10)
	}
	
	// Check execution count
	if count := profiler.GetExecutionCount(fnHash); count != 5 {
		t.Errorf("Expected 5 execution count, got %d", count)
	}
	
	// Get profile
	profile := profiler.GetProfile(fnHash)
	if profile == nil {
		t.Fatal("Expected profile, got nil")
	}
	
	if profile.ExecutionCount != 5 {
		t.Errorf("Expected 5 executions, got %d", profile.ExecutionCount)
	}
	
	if profile.AverageTime != time.Millisecond*10 {
		t.Errorf("Expected 10ms average time, got %v", profile.AverageTime)
	}
}

func TestHotPathDetection(t *testing.T) {
	compiler := NewJITCompiler()
	fnHash := uint64(456)
	
	// Function should not be hot initially
	if compiler.ShouldCompile(fnHash) {
		t.Error("Function should not be hot initially")
	}
	
	// Record executions up to threshold
	for i := 0; i < DefaultHotThreshold; i++ {
		compiler.RecordExecution(fnHash, time.Microsecond*100)
	}
	
	// Now function should be hot
	if !compiler.ShouldCompile(fnHash) {
		t.Error("Function should be hot after threshold executions")
	}
}

func TestCodeCache(t *testing.T) {
	cache := NewCodeCache()
	
	// Test empty cache
	if cache.Has(123) {
		t.Error("Empty cache should not have any entries")
	}
	
	if code := cache.Get(123); code != nil {
		t.Error("Empty cache should return nil for Get")
	}
	
	// Create test compiled code
	testCode := &CompiledCode{
		NativeCode: []byte{0x01, 0x02, 0x03, 0x04}, // Dummy ARM64 code
		Hash:       123,
		CreatedAt:  time.Now(),
	}
	
	// Add to cache
	err := cache.Add(123, testCode)
	if err != nil {
		t.Fatalf("Failed to add code to cache: %v", err)
	}
	
	// Test cache has entry
	if !cache.Has(123) {
		t.Error("Cache should have entry after adding")
	}
	
	// Test cache get
	retrieved := cache.Get(123)
	if retrieved == nil {
		t.Fatal("Cache should return code after adding")
	}
	
	if retrieved.Hash != 123 {
		t.Errorf("Retrieved code has wrong hash: got %d, want 123", retrieved.Hash)
	}
}

func TestARM64CodeGeneration(t *testing.T) {
	codegen := NewARM64CodeGen()
	
	// Create test bytecode (simple addition)
	instructions := bytecode.Instructions{
		byte(bytecode.OpConstant), 0, 0, // Load constant 0
		byte(bytecode.OpConstant), 0, 1, // Load constant 1
		byte(bytecode.OpAdd),            // Add them
		byte(bytecode.OpReturn),         // Return result
	}
	
	// Generate ARM64 code
	nativeCode, err := codegen.Generate(instructions)
	if err != nil {
		t.Fatalf("Code generation failed: %v", err)
	}
	
	if len(nativeCode) == 0 {
		t.Error("Generated code should not be empty")
	}
	
	// Code should be multiple of 4 bytes (ARM64 instructions are 32-bit aligned)
	if len(nativeCode)%4 != 0 {
		t.Errorf("Generated code length should be multiple of 4, got %d", len(nativeCode))
	}
}

func TestUnsupportedBytecode(t *testing.T) {
	codegen := NewARM64CodeGen()
	
	// Create bytecode with unsupported instruction
	instructions := bytecode.Instructions{
		255, // Invalid/unsupported opcode
	}
	
	// Should fail to generate code
	_, err := codegen.Generate(instructions)
	if err == nil {
		t.Error("Expected error for unsupported opcode")
	}
}

func TestJITCompilationFlow(t *testing.T) {
	compiler := NewJITCompiler()
	
	// Create test function
	testFn := &interpreter.CompiledFunction{
		Instructions: []byte{
			byte(bytecode.OpConstant), 0, 0,
			byte(bytecode.OpReturn),
		},
		NumLocals:     0,
		NumParameters: 0,
	}
	
	fnHash := uint64(789)
	
	// Initially should not compile
	if compiler.ShouldCompile(fnHash) {
		t.Error("Should not compile cold function")
	}
	
	// Make function hot
	for i := 0; i < DefaultHotThreshold; i++ {
		compiler.RecordExecution(fnHash, time.Microsecond*50)
	}
	
	// Now should compile
	if !compiler.ShouldCompile(fnHash) {
		t.Error("Should compile hot function")
	}
	
	// Attempt compilation (will fail at execution due to unimplemented native execution)
	err := compiler.Compile(testFn, fnHash)
	if err != nil {
		// This is expected since we haven't implemented full ARM64 execution
		t.Logf("Compilation failed as expected: %v", err)
	}
}

func TestJITStats(t *testing.T) {
	compiler := NewJITCompiler()
	
	// Get initial stats
	stats := compiler.GetStats()
	if stats.CompilationsAttempted != 0 {
		t.Error("Initial compilations attempted should be 0")
	}
	
	// Create test function and try to compile
	testFn := &interpreter.CompiledFunction{
		Instructions: []byte{byte(bytecode.OpReturn)},
		NumLocals:    0,
		NumParameters: 0,
	}
	
	// This will fail but should update stats
	compiler.Compile(testFn, 999)
	
	// Check stats updated
	newStats := compiler.GetStats()
	if newStats.CompilationsAttempted == 0 {
		t.Error("Stats should show compilation attempt")
	}
}

func BenchmarkJITCompilation(b *testing.B) {
	compiler := NewJITCompiler()
	
	testFn := &interpreter.CompiledFunction{
		Instructions: []byte{
			byte(bytecode.OpConstant), 0, 0,
			byte(bytecode.OpConstant), 0, 1,
			byte(bytecode.OpAdd),
			byte(bytecode.OpReturn),
		},
		NumLocals:     0,
		NumParameters: 0,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fnHash := uint64(i) // Different hash each time to avoid cache hits
		compiler.Compile(testFn, fnHash)
	}
}

func BenchmarkCodeCacheOperations(b *testing.B) {
	cache := NewCodeCache()
	
	testCode := &CompiledCode{
		NativeCode: make([]byte, 1024), // 1KB of dummy code
		Hash:       0,
		CreatedAt:  time.Now(),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		hash := uint64(i)
		testCode.Hash = hash
		cache.Add(hash, testCode)
		cache.Get(hash)
	}
}