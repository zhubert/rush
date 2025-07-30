package aot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAOTCompilerBasic(t *testing.T) {
	// Create a temporary source file
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "test.rush")
	
	source := `5 + 10`
	
	err := os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}
	
	// Create AOT compiler with default options
	options := DefaultCompileOptions()
	options.OutputPath = filepath.Join(tempDir, "test_executable")
	compiler := NewAOTCompiler(options)
	
	// Test compilation
	err = compiler.CompileProgram(sourceFile)
	if err != nil {
		t.Fatalf("AOT compilation failed: %v", err)
	}
	
	// Verify executable was created
	if _, err := os.Stat(options.OutputPath); os.IsNotExist(err) {
		t.Fatalf("Executable not created at %s", options.OutputPath)
	}
	
	// Verify statistics
	stats := compiler.GetStats()
	if stats.CompileTime == 0 {
		t.Error("Expected non-zero compile time")
	}
	if stats.ExecutableSize == 0 {
		t.Error("Expected non-zero executable size")
	}
}

func TestCompileOptionsValidation(t *testing.T) {
	tests := []struct {
		name           string
		options        *CompileOptions
		expectError    bool
		errorContains  string
	}{
		{
			name:    "Valid default options",
			options: DefaultCompileOptions(),
			expectError: false,
		},
		{
			name: "Valid custom options",
			options: &CompileOptions{
				OptimizationLevel: 2,
				TargetOS:         "linux",
				TargetArch:       "arm64",
				MinimalRuntime:   true,
				StaticLinking:    true,
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewAOTCompiler(tt.options)
			if compiler == nil && !tt.expectError {
				t.Error("Expected compiler to be created successfully")
			}
		})
	}
}

func TestOptimizationLevels(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "test.rush")
	
	source := `1 + 2 * 3`
	
	err := os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}
	
	optimizationLevels := []int{0, 1, 2}
	var stats []*AOTStats
	
	for _, level := range optimizationLevels {
		options := DefaultCompileOptions()
		options.OptimizationLevel = level
		options.OutputPath = filepath.Join(tempDir, "test_O"+string(rune('0'+level)))
		
		compiler := NewAOTCompiler(options)
		err := compiler.CompileProgram(sourceFile)
		if err != nil {
			t.Fatalf("AOT compilation failed at optimization level %d: %v", level, err)
		}
		
		stats = append(stats, compiler.GetStats())
	}
	
	// Verify that higher optimization levels generally take more time
	// (This is a heuristic test, actual results may vary)
	for i := 1; i < len(stats); i++ {
		if stats[i].AnalysisTime < stats[i-1].AnalysisTime {
			t.Logf("Note: Higher optimization level %d took less analysis time than level %d", 
				   optimizationLevels[i], optimizationLevels[i-1])
		}
	}
}

func TestCrossCompilation(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "test.rush")
	
	source := `42`
	
	err := os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}
	
	targets := []struct {
		os   string
		arch string
	}{
		{"darwin", "arm64"},
		{"linux", "arm64"},
		// Windows PE not implemented yet
		// {"windows", "arm64"},
	}
	
	for _, target := range targets {
		t.Run(target.os+"_"+target.arch, func(t *testing.T) {
			options := DefaultCompileOptions()
			options.TargetOS = target.os
			options.TargetArch = target.arch
			options.OutputPath = filepath.Join(tempDir, "test_"+target.os+"_"+target.arch)
			
			compiler := NewAOTCompiler(options)
			err := compiler.CompileProgram(sourceFile)
			
			if target.os == "windows" {
				// Windows PE not implemented yet, expect error
				if err == nil {
					t.Error("Expected error for Windows PE compilation, but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Cross compilation failed for %s/%s: %v", target.os, target.arch, err)
			}
			
			// Verify executable exists
			if _, err := os.Stat(options.OutputPath); os.IsNotExist(err) {
				t.Errorf("Cross-compiled executable not created for %s/%s", target.os, target.arch)
			}
		})
	}
}

func TestMinimalRuntimeOption(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "test.rush")
	
	source := `42`
	
	err := os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}
	
	// Test with full runtime
	fullOptions := DefaultCompileOptions()
	fullOptions.MinimalRuntime = false
	fullOptions.OutputPath = filepath.Join(tempDir, "test_full")
	fullCompiler := NewAOTCompiler(fullOptions)
	
	err = fullCompiler.CompileProgram(sourceFile)
	if err != nil {
		t.Fatalf("Compilation with full runtime failed: %v", err)
	}
	
	// Test with minimal runtime
	minOptions := DefaultCompileOptions()
	minOptions.MinimalRuntime = true
	minOptions.OutputPath = filepath.Join(tempDir, "test_minimal")
	minCompiler := NewAOTCompiler(minOptions)
	
	err = minCompiler.CompileProgram(sourceFile)
	if err != nil {
		t.Fatalf("Compilation with minimal runtime failed: %v", err)
	}
	
	// Compare executable sizes (minimal should be smaller or equal for now)
	fullStats := fullCompiler.GetStats()
	minStats := minCompiler.GetStats()
	
	// NOTE: Current simple LLVM implementation produces same size for both runtimes
	// TODO: Implement actual minimal runtime that produces smaller executables
	if minStats.ExecutableSize > fullStats.ExecutableSize {
		t.Errorf("Minimal runtime executable should not be larger than full runtime: min=%d, full=%d", 
				 minStats.ExecutableSize, fullStats.ExecutableSize)
	} else {
		t.Logf("Runtime sizes - minimal: %d bytes, full: %d bytes", minStats.ExecutableSize, fullStats.ExecutableSize)
	}
}

func TestStaticAnalyzer(t *testing.T) {
	analyzer := NewStaticAnalyzer()
	
	// This test requires a proper bytecode object
	// For now, we'll test the analyzer creation
	if analyzer == nil {
		t.Fatal("Failed to create static analyzer")
	}
	
	if len(analyzer.passes) == 0 {
		t.Error("Expected optimization passes to be registered")
	}
}

func TestOptimizationPasses(t *testing.T) {
	passes := []OptimizationPass{
		NewDeadCodeEliminationPass(),
		NewConstantPropagationPass(),
		NewFunctionInliningPass(),
		NewTailCallOptimizationPass(),
		NewLoopOptimizationPass(),
	}
	
	for _, pass := range passes {
		t.Run(pass.Name(), func(t *testing.T) {
			if pass.Name() == "" {
				t.Error("Optimization pass should have a name")
			}
			
			// Test enablement at different levels
			if !pass.IsEnabled(2) && pass.Name() != "FunctionInlining" && pass.Name() != "LoopOptimization" {
				t.Errorf("Pass %s should be enabled at level 2", pass.Name())
			}
		})
	}
}

func TestRuntimeComponents(t *testing.T) {
	builder := NewRuntimeBuilder()
	if builder == nil {
		t.Fatal("Failed to create runtime builder")
	}
	
	if len(builder.components) == 0 {
		t.Error("Expected runtime components to be registered")
	}
	
	// Test building runtime with different options
	options := DefaultCompileOptions()
	runtime, err := builder.BuildRuntime(options)
	if err != nil {
		t.Fatalf("Failed to build runtime: %v", err)
	}
	
	if len(runtime.Code) == 0 {
		t.Error("Expected runtime to contain code")
	}
	
	if len(runtime.Symbols) == 0 {
		t.Error("Expected runtime to contain symbols")
	}
}

func TestExecutableFormats(t *testing.T) {
	linker := NewExecutableLinker()
	if linker == nil {
		t.Fatal("Failed to create executable linker")
	}
	
	// Create dummy native code and runtime for testing
	nativeCode := &NativeCode{
		Instructions: []byte{0x01, 0x02, 0x03, 0x04},
		Symbols:     make(map[string]Symbol),
		EntryPoint:  0,
	}
	
	runtime := &MinimalRuntime{
		Code:    []byte{0x10, 0x20, 0x30, 0x40},
		Symbols: make(map[string]Symbol),
		Size:    4,
	}
	
	testFormats := []struct {
		name    string
		targetOS string
		expectError bool
	}{
		{"Mach-O", "darwin", false},
		{"ELF", "linux", false},
		{"PE", "windows", true}, // Not implemented yet
	}
	
	for _, tf := range testFormats {
		t.Run(tf.name, func(t *testing.T) {
			options := DefaultCompileOptions()
			options.TargetOS = tf.targetOS
			
			executable, err := linker.LinkExecutable(nativeCode, runtime, options)
			
			if tf.expectError {
				if err == nil {
					t.Errorf("Expected error for %s format, but got none", tf.name)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Failed to link %s executable: %v", tf.name, err)
			}
			
			if len(executable.Data) == 0 {
				t.Errorf("Expected %s executable to contain data", tf.name)
			}
			
			if executable.EntryPoint < 0 {
				t.Errorf("Expected valid entry point for %s executable", tf.name)
			}
		})
	}
}

func TestPerformanceRequirements(t *testing.T) {
	// Test that AOT compilation meets performance targets
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "perf_test.rush")
	
	// Create a simple program for testing
	source := `1 + 2 + 3 + 4 + 5`
	
	err := os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write test source file: %v", err)
	}
	
	options := DefaultCompileOptions()
	options.OptimizationLevel = 2
	options.OutputPath = filepath.Join(tempDir, "perf_test")
	
	compiler := NewAOTCompiler(options)
	
	start := time.Now()
	err = compiler.CompileProgram(sourceFile)
	compileTime := time.Since(start)
	
	if err != nil {
		t.Fatalf("Performance test compilation failed: %v", err)
	}
	
	// AOT compilation should be reasonably fast (under 10 seconds for this test)
	if compileTime > 10*time.Second {
		t.Errorf("AOT compilation took too long: %v (expected < 10s)", compileTime)
	}
	
	stats := compiler.GetStats()
	
	// Executable should be reasonably sized (under 10MB for this simple program)
	if stats.ExecutableSize > 10*1024*1024 {
		t.Errorf("Executable too large: %d bytes (expected < 10MB)", stats.ExecutableSize)
	}
	
	t.Logf("Performance test results:")
	t.Logf("  Compile time: %v", compileTime)
	t.Logf("  Executable size: %d bytes", stats.ExecutableSize)
	t.Logf("  Bytecode instructions: %d", stats.BytecodeInstructions)
	t.Logf("  Native instructions: %d", stats.NativeInstructions)
}

func TestErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test compilation of invalid source
	invalidSource := `1 +` // Syntax error
	sourceFile := filepath.Join(tempDir, "invalid.rush")
	
	err := os.WriteFile(sourceFile, []byte(invalidSource), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid source file: %v", err)
	}
	
	options := DefaultCompileOptions()
	options.OutputPath = filepath.Join(tempDir, "invalid_executable")
	compiler := NewAOTCompiler(options)
	
	err = compiler.CompileProgram(sourceFile)
	if err == nil {
		t.Error("Expected compilation error for invalid source, but got none")
	}
	
	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("Expected parse error, got: %v", err)
	}
	
	// Test compilation of non-existent file
	err = compiler.CompileProgram("nonexistent.rush")
	if err == nil {
		t.Error("Expected error for non-existent file, but got none")
	}
}