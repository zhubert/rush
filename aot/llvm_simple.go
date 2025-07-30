//go:build llvm14 || llvm15 || llvm16 || llvm17 || llvm18 || llvm19 || llvm20

package aot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rush/compiler"
	
	"tinygo.org/x/go-llvm"
)

// SimpleLLVMCodeGenerator is a minimal LLVM-based code generator for testing
type SimpleLLVMCodeGenerator struct{}

// SimpleNativeCode represents simple LLVM output
type SimpleNativeCode struct {
	ObjectFile string
	EntryPoint string
}

// NewSimpleLLVMCodeGenerator creates a minimal LLVM code generator
func NewSimpleLLVMCodeGenerator() *SimpleLLVMCodeGenerator {
	return &SimpleLLVMCodeGenerator{}
}

// GenerateExecutable creates a simple "Hello World" program using LLVM
func (g *SimpleLLVMCodeGenerator) GenerateExecutable(bytecode *compiler.Bytecode) (*SimpleNativeCode, error) {
	// Initialize LLVM - must be called before any LLVM operations
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
	
	// Create context and module
	context := llvm.NewContext()
	defer context.Dispose()
	
	module := context.NewModule("rush_program")
	defer module.Dispose()
	
	builder := context.NewBuilder()
	defer builder.Dispose()
	
	// Create main function type: int main()
	mainType := llvm.FunctionType(context.Int32Type(), nil, false)
	mainFunc := llvm.AddFunction(module, "main", mainType)
	
	// Create entry basic block
	entry := context.AddBasicBlock(mainFunc, "entry")
	builder.SetInsertPointAtEnd(entry)
	
	// For now, just return 0
	builder.CreateRet(llvm.ConstInt(context.Int32Type(), 0, false))
	
	// Verify module
	if err := llvm.VerifyModule(module, llvm.ReturnStatusAction); err != nil {
		return nil, fmt.Errorf("module verification failed: %w", err)
	}
	
	// Get target
	targetTriple := llvm.DefaultTargetTriple()
	target, err := llvm.GetTargetFromTriple(targetTriple)
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}
	
	// Create target machine
	machine := target.CreateTargetMachine(
		targetTriple,
		"generic",
		"",
		llvm.CodeGenLevelDefault,
		llvm.RelocDefault,
		llvm.CodeModelDefault,
	)
	defer machine.Dispose()
	
	// Set data layout and target
	targetData := machine.CreateTargetData()
	defer targetData.Dispose()
	module.SetDataLayout(targetData.String())
	module.SetTarget(targetTriple)
	
	// Create temporary object file
	tempDir := os.TempDir()
	objectFile := filepath.Join(tempDir, "rush_simple.o")
	
	// Compile to memory buffer first
	memBuf, err := machine.EmitToMemoryBuffer(module, llvm.ObjectFile)
	if err != nil {
		return nil, fmt.Errorf("failed to emit to memory buffer: %w", err)
	}
	defer memBuf.Dispose()
	
	// Write memory buffer to file
	if err := os.WriteFile(objectFile, memBuf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("failed to write object file: %w", err)
	}
	
	return &SimpleNativeCode{
		ObjectFile: objectFile,
		EntryPoint: "main",
	}, nil
}

// LinkExecutable links the object file to create an executable
func (g *SimpleLLVMCodeGenerator) LinkExecutable(objectFile, outputPath string) error {
	// Use clang to link
	cmd := exec.Command("clang", "-o", outputPath, objectFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("linking failed: %s\nOutput: %s", err, string(output))
	}
	return nil
}