package aot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"rush/ast"
	"rush/compiler"
	"rush/parser"
	"rush/lexer"
)

// AOTCompiler manages Ahead-of-Time compilation for Rush programs
type AOTCompiler struct {
	analyzer     *StaticAnalyzer
	codeGen      *SimpleLLVMCodeGenerator
	options      *CompileOptions
	stats        *AOTStats
}

// CompileOptions configures AOT compilation behavior
type CompileOptions struct {
	OutputPath       string
	OptimizationLevel int  // 0=none, 1=basic, 2=aggressive
	DebugInfo        bool
	CrossCompile     bool
	TargetOS         string // "darwin", "linux", "windows"
	TargetArch       string // "arm64", "amd64"
	MinimalRuntime   bool   // Strip unused runtime components
	StaticLinking    bool   // Create fully static executable
}

// AOTStats tracks compilation statistics
type AOTStats struct {
	CompileTime         time.Duration
	AnalysisTime        time.Duration
	CodeGenTime         time.Duration
	LinkTime            time.Duration
	BytecodeInstructions int64
	NativeInstructions   int64
	ExecutableSize       int64
	OptimizationsApplied int64
}

// NewAOTCompiler creates a new AOT compiler instance
func NewAOTCompiler(options *CompileOptions) *AOTCompiler {
	return &AOTCompiler{
		analyzer: NewStaticAnalyzer(),
		codeGen:  NewSimpleLLVMCodeGenerator(),
		options:  options,
		stats:    &AOTStats{},
	}
}

// DefaultCompileOptions returns sensible default compilation options
func DefaultCompileOptions() *CompileOptions {
	return &CompileOptions{
		OptimizationLevel: 1,
		DebugInfo:        false,
		CrossCompile:     false,
		TargetOS:         "darwin",
		TargetArch:       "arm64",
		MinimalRuntime:   true,
		StaticLinking:    true,
	}
}

// CompileProgram performs complete AOT compilation of a Rush program
func (c *AOTCompiler) CompileProgram(sourcePath string) error {
	startTime := time.Now()
	defer func() {
		c.stats.CompileTime = time.Since(startTime)
	}()

	// Step 1: Parse source code to AST
	program, err := c.parseSource(sourcePath)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	// Step 2: Compile AST to bytecode
	bytecode, err := c.compileToBytecode(program)
	if err != nil {
		return fmt.Errorf("bytecode compilation failed: %w", err)
	}
	c.stats.BytecodeInstructions = int64(len(bytecode.Instructions))

	// Step 3: Static analysis and optimization
	analysisStart := time.Now()
	optimizedBytecode, err := c.analyzer.AnalyzeAndOptimize(bytecode)
	if err != nil {
		return fmt.Errorf("static analysis failed: %w", err)
	}
	c.stats.AnalysisTime = time.Since(analysisStart)

	// Step 4: Generate LLVM IR and compile to object file
	codeGenStart := time.Now()
	llvmCode, err := c.codeGen.GenerateExecutable(optimizedBytecode)
	if err != nil {
		return fmt.Errorf("LLVM code generation failed: %w", err)
	}
	c.stats.CodeGenTime = time.Since(codeGenStart)

	// Step 5: Link executable using LLVM/Clang
	linkStart := time.Now()
	outputPath := c.getOutputPath(sourcePath)
	if err := c.linkWithRuntime(llvmCode.ObjectFile, outputPath); err != nil {
		return fmt.Errorf("linking failed: %w", err)
	}
	c.stats.LinkTime = time.Since(linkStart)

	// Get executable size
	if stat, err := os.Stat(outputPath); err == nil {
		c.stats.ExecutableSize = stat.Size()
	}
	return nil
}

// parseSource parses Rush source code into AST
func (c *AOTCompiler) parseSource(sourcePath string) (*ast.Program, error) {
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("read source file: %w", err)
	}

	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.Errors())
	}

	return program, nil
}

// compileToBytecode compiles AST to bytecode
func (c *AOTCompiler) compileToBytecode(program *ast.Program) (*compiler.Bytecode, error) {
	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		return nil, err
	}
	return comp.Bytecode(), nil
}

// getOutputPath determines the output executable path
func (c *AOTCompiler) getOutputPath(sourcePath string) string {
	if c.options.OutputPath != "" {
		return c.options.OutputPath
	}

	// Default: same directory as source, with executable extension
	dir := filepath.Dir(sourcePath)
	base := filepath.Base(sourcePath)
	name := base[:len(base)-len(filepath.Ext(base))]
	
	if c.options.TargetOS == "windows" {
		name += ".exe"
	}
	
	return filepath.Join(dir, name)
}

// linkWithRuntime links the object file with the Rush runtime and system libraries
func (c *AOTCompiler) linkWithRuntime(objectFile, outputPath string) error {
	// First, compile the runtime.c file
	runtimeObj, err := c.compileRuntime()
	if err != nil {
		return fmt.Errorf("failed to compile runtime: %w", err)
	}
	defer os.Remove(runtimeObj) // Clean up temporary file
	
	// Use clang to link everything together
	return c.codeGen.LinkExecutable(objectFile, outputPath)
}

// compileRuntime compiles the C runtime to an object file
func (c *AOTCompiler) compileRuntime() (string, error) {
	// Get path to runtime.c
	runtimePath := "runtime/runtime.c"
	
	// Create temporary object file
	tempDir := os.TempDir()
	runtimeObj := filepath.Join(tempDir, "rush_runtime.o")
	
	// Compile runtime.c to object file
	cmd := exec.Command("clang", "-c", "-o", runtimeObj, runtimePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime compilation failed: %s\nOutput: %s", err, string(output))
	}
	
	return runtimeObj, nil
}

// GetStats returns compilation statistics
func (c *AOTCompiler) GetStats() *AOTStats {
	return c.stats
}

// PrintStats prints compilation statistics to stdout
func (c *AOTCompiler) PrintStats() {
	fmt.Printf("AOT Compilation Statistics:\n")
	fmt.Printf("  Total Time: %v\n", c.stats.CompileTime)
	fmt.Printf("  Analysis Time: %v\n", c.stats.AnalysisTime)
	fmt.Printf("  Code Generation Time: %v\n", c.stats.CodeGenTime)
	fmt.Printf("  Link Time: %v\n", c.stats.LinkTime)
	fmt.Printf("  Bytecode Instructions: %d\n", c.stats.BytecodeInstructions)
	fmt.Printf("  Native Instructions: %d\n", c.stats.NativeInstructions)
	fmt.Printf("  Executable Size: %d bytes\n", c.stats.ExecutableSize)
	fmt.Printf("  Optimizations Applied: %d\n", c.stats.OptimizationsApplied)
}