//go:build !llvm14 && !llvm15 && !llvm16 && !llvm17 && !llvm18 && !llvm19 && !llvm20

package aot

import (
	"fmt"
	"rush/compiler"
)

// SimpleLLVMCodeGenerator is a fallback when LLVM is not available
type SimpleLLVMCodeGenerator struct{}

// SimpleNativeCode represents fallback output
type SimpleNativeCode struct {
	ObjectFile string
	EntryPoint string
}

// NewSimpleLLVMCodeGenerator creates a fallback code generator
func NewSimpleLLVMCodeGenerator() *SimpleLLVMCodeGenerator {
	return &SimpleLLVMCodeGenerator{}
}

// GenerateExecutable returns an error when LLVM is not available
func (g *SimpleLLVMCodeGenerator) GenerateExecutable(bytecode *compiler.Bytecode) (*SimpleNativeCode, error) {
	return nil, fmt.Errorf("LLVM AOT compilation not available: build with -tags llvm18 (or other supported LLVM version)")
}

// LinkExecutable returns an error when LLVM is not available  
func (g *SimpleLLVMCodeGenerator) LinkExecutable(objectFile, outputPath string) error {
	return fmt.Errorf("LLVM AOT compilation not available: build with -tags llvm18 (or other supported LLVM version)")
}