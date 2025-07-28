package jit

import (
	"fmt"
	"sync"
	"time"

	"rush/bytecode"
	"rush/interpreter"
)

// CompilationThreshold defines when to trigger JIT compilation
const (
	DefaultHotThreshold     = 100   // Function calls before JIT compilation
	DefaultCompileTimeout   = 5     // Max seconds for compilation
	DefaultMaxCompiledFuncs = 1000  // Max number of JIT compiled functions
)

// JITCompiler manages Just-In-Time compilation for hot functions
type JITCompiler struct {
	cache           *CodeCache
	profiler        *ExecutionProfiler
	hotThreshold    int
	compileTimeout  time.Duration
	maxCompiledFuncs int
	mu              sync.RWMutex
	stats           *JITStats
}

// JITStats tracks JIT compilation statistics
type JITStats struct {
	CompilationsAttempted int64
	CompilationsSucceeded int64
	CompilationsFailed    int64
	CompilationTime       time.Duration
	JITHits               int64
	JITMisses             int64
	Deoptimizations       int64
	CacheEvictions        int64
}

// NewJITCompiler creates a new JIT compiler instance
func NewJITCompiler() *JITCompiler {
	return &JITCompiler{
		cache:            NewCodeCache(),
		profiler:         NewExecutionProfiler(),
		hotThreshold:     DefaultHotThreshold,
		compileTimeout:   time.Duration(DefaultCompileTimeout) * time.Second,
		maxCompiledFuncs: DefaultMaxCompiledFuncs,
		stats:            &JITStats{},
	}
}

// ShouldCompile determines if a function should be JIT compiled
func (j *JITCompiler) ShouldCompile(fnHash uint64) bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	
	// Check if already compiled
	if j.cache.Has(fnHash) {
		return false
	}
	
	// Check execution count threshold
	count := j.profiler.GetExecutionCount(fnHash)
	return count >= j.hotThreshold
}

// Compile attempts to JIT compile a function
func (j *JITCompiler) Compile(fn *interpreter.CompiledFunction, fnHash uint64) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	
	// Check cache again (double-checked locking)
	if j.cache.Has(fnHash) {
		return nil
	}
	
	j.stats.CompilationsAttempted++
	startTime := time.Now()
	
	// Create compilation context with timeout
	ctx := &CompilationContext{
		Function:  fn,
		Hash:      fnHash,
		Timeout:   j.compileTimeout,
	}
	
	// Compile to ARM64 native code
	compiledCode, err := j.compileToARM64(ctx)
	if err != nil {
		j.stats.CompilationsFailed++
		return fmt.Errorf("JIT compilation failed: %v", err)
	}
	
	// Add to cache
	err = j.cache.Add(fnHash, compiledCode)
	if err != nil {
		j.stats.CompilationsFailed++
		return fmt.Errorf("failed to cache compiled code: %v", err)
	}
	
	j.stats.CompilationsSucceeded++
	j.stats.CompilationTime += time.Since(startTime)
	
	return nil
}

// Execute runs JIT compiled code for a function
func (j *JITCompiler) Execute(fnHash uint64, args []interpreter.Value, globals []interpreter.Value) (interpreter.Value, error) {
	j.mu.RLock()
	compiledCode := j.cache.Get(fnHash)
	j.mu.RUnlock()
	
	if compiledCode == nil {
		j.stats.JITMisses++
		return nil, fmt.Errorf("no compiled code found for function")
	}
	
	j.stats.JITHits++
	
	// Execute ARM64 compiled code
	result, err := compiledCode.Execute(args, globals)
	if err != nil {
		// Handle deoptimization
		j.stats.Deoptimizations++
		return nil, err
	}
	
	return result, nil
}

// compileToARM64 compiles bytecode to ARM64 native code
func (j *JITCompiler) compileToARM64(ctx *CompilationContext) (*CompiledCode, error) {
	// Create ARM64 code generator
	codegen := NewARM64CodeGen()
	
	// Analyze bytecode instructions
	instructions := bytecode.Instructions(ctx.Function.Instructions)
	
	// Generate ARM64 code
	nativeCode, err := codegen.Generate(instructions)
	if err != nil {
		return nil, fmt.Errorf("ARM64 code generation failed: %v", err)
	}
	
	// Create executable code object
	compiledCode := &CompiledCode{
		NativeCode:   nativeCode,
		Function:     ctx.Function,
		Hash:         ctx.Hash,
		CreatedAt:    time.Now(),
	}
	
	return compiledCode, nil
}

// GetStats returns a copy of JIT compilation statistics
func (j *JITCompiler) GetStats() JITStats {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return *j.stats
}

// RecordExecution tracks function execution for hot path detection
func (j *JITCompiler) RecordExecution(fnHash uint64, executionTime time.Duration) {
	j.profiler.RecordExecution(fnHash, executionTime)
}

// CompilationContext holds context for JIT compilation
type CompilationContext struct {
	Function *interpreter.CompiledFunction
	Hash     uint64
	Timeout  time.Duration
}