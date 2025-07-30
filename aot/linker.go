package aot

import (
	"fmt"
)

// ExecutableLinker combines compiled code with runtime into a standalone executable
type ExecutableLinker struct {
	format ExecutableFormat
}

// ExecutableFormat defines the target executable format
type ExecutableFormat int

const (
	FormatMachO ExecutableFormat = iota // macOS Mach-O
	FormatELF                           // Linux ELF
	FormatPE                            // Windows PE
)

// LinkedExecutable represents the final executable
type LinkedExecutable struct {
	Data       []byte
	Format     ExecutableFormat
	EntryPoint int64
	Symbols    map[string]Symbol
}

// NewExecutableLinker creates a new executable linker
func NewExecutableLinker() *ExecutableLinker {
	return &ExecutableLinker{
		format: FormatMachO, // Default to macOS
	}
}

// LinkExecutable combines native code and runtime into executable
func (l *ExecutableLinker) LinkExecutable(code *NativeCode, runtime *MinimalRuntime, options *CompileOptions) (*LinkedExecutable, error) {
	// Determine target format based on options
	format := l.getTargetFormat(options)
	
	switch format {
	case FormatMachO:
		return l.linkMachO(code, runtime, options)
	case FormatELF:
		return l.linkELF(code, runtime, options)
	case FormatPE:
		return l.linkPE(code, runtime, options)
	default:
		return nil, fmt.Errorf("unsupported executable format: %d", format)
	}
}

// getTargetFormat determines executable format from compile options
func (l *ExecutableLinker) getTargetFormat(options *CompileOptions) ExecutableFormat {
	switch options.TargetOS {
	case "darwin":
		return FormatMachO
	case "linux":
		return FormatELF
	case "windows":
		return FormatPE
	default:
		return FormatMachO // Default
	}
}

// linkMachO creates a macOS Mach-O executable
func (l *ExecutableLinker) linkMachO(code *NativeCode, runtime *MinimalRuntime, options *CompileOptions) (*LinkedExecutable, error) {
	// Simplified Mach-O structure
	// In a real implementation, this would create proper Mach-O headers
	
	// Calculate sizes and offsets
	headerSize := 1024          // Simplified header size
	runtimeOffset := headerSize
	codeOffset := runtimeOffset + len(runtime.Code)
	totalSize := codeOffset + len(code.Instructions)
	
	// Create executable data
	execData := make([]byte, totalSize)
	
	// Write simplified Mach-O header
	if err := l.writeMachOHeader(execData[:headerSize], codeOffset, len(code.Instructions)); err != nil {
		return nil, fmt.Errorf("failed to write Mach-O header: %w", err)
	}
	
	// Copy runtime code
	copy(execData[runtimeOffset:], runtime.Code)
	
	// Copy program code
	copy(execData[codeOffset:], code.Instructions)
	
	// Combine symbol tables
	allSymbols := make(map[string]Symbol)
	for name, symbol := range runtime.Symbols {
		symbol.Offset += runtimeOffset
		allSymbols[name] = symbol
	}
	for name, symbol := range code.Symbols {
		symbol.Offset += codeOffset
		allSymbols[name] = symbol
	}
	
	return &LinkedExecutable{
		Data:       execData,
		Format:     FormatMachO,
		EntryPoint: int64(runtimeOffset), // Start from runtime startup
		Symbols:    allSymbols,
	}, nil
}

// writeMachOHeader writes a simplified Mach-O header
func (l *ExecutableLinker) writeMachOHeader(header []byte, codeOffset, codeSize int) error {
	// This is a very simplified Mach-O header
	// A real implementation would need proper load commands, sections, etc.
	
	// Mach-O magic number for ARM64
	header[0] = 0xFE
	header[1] = 0xED
	header[2] = 0xFA
	header[3] = 0xCF
	
	// CPU type (ARM64)
	header[4] = 0x0C
	header[5] = 0x01
	header[6] = 0x00
	header[7] = 0x00
	
	// Entry point (simplified)
	entryPoint := uint64(codeOffset)
	for i := 0; i < 8; i++ {
		header[24+i] = byte(entryPoint >> (i * 8))
	}
	
	return nil
}

// linkELF creates a Linux ELF executable
func (l *ExecutableLinker) linkELF(code *NativeCode, runtime *MinimalRuntime, options *CompileOptions) (*LinkedExecutable, error) {
	// Simplified ELF structure
	headerSize := 64 // ELF64 header size
	programHeaderSize := 56 // Program header size
	totalHeaderSize := headerSize + programHeaderSize
	
	runtimeOffset := totalHeaderSize
	codeOffset := runtimeOffset + len(runtime.Code)
	totalSize := codeOffset + len(code.Instructions)
	
	execData := make([]byte, totalSize)
	
	// Write ELF header
	if err := l.writeELFHeader(execData[:headerSize], codeOffset); err != nil {
		return nil, fmt.Errorf("failed to write ELF header: %w", err)
	}
	
	// Write program header
	if err := l.writeELFProgramHeader(execData[headerSize:totalHeaderSize], runtimeOffset, totalSize-totalHeaderSize); err != nil {
		return nil, fmt.Errorf("failed to write ELF program header: %w", err)
	}
	
	// Copy runtime and code
	copy(execData[runtimeOffset:], runtime.Code)
	copy(execData[codeOffset:], code.Instructions)
	
	// Combine symbol tables
	allSymbols := make(map[string]Symbol)
	for name, symbol := range runtime.Symbols {
		symbol.Offset += runtimeOffset
		allSymbols[name] = symbol
	}
	for name, symbol := range code.Symbols {
		symbol.Offset += codeOffset
		allSymbols[name] = symbol
	}
	
	return &LinkedExecutable{
		Data:       execData,
		Format:     FormatELF,
		EntryPoint: int64(runtimeOffset),
		Symbols:    allSymbols,
	}, nil
}

// writeELFHeader writes a simplified ELF header
func (l *ExecutableLinker) writeELFHeader(header []byte, entryPoint int) error {
	// ELF magic number
	header[0] = 0x7F
	header[1] = 'E'
	header[2] = 'L'
	header[3] = 'F'
	
	// 64-bit, little-endian
	header[4] = 2 // EI_CLASS: 64-bit
	header[5] = 1 // EI_DATA: little-endian
	header[6] = 1 // EI_VERSION: current
	
	// Object file type: executable
	header[16] = 2 // ET_EXEC
	header[17] = 0
	
	// Machine type: ARM64
	header[18] = 0xB7 // EM_AARCH64
	header[19] = 0x00
	
	// Entry point
	entryAddr := uint64(0x400000 + entryPoint) // Standard load address
	for i := 0; i < 8; i++ {
		header[24+i] = byte(entryAddr >> (i * 8))
	}
	
	return nil
}

// writeELFProgramHeader writes a simplified ELF program header
func (l *ExecutableLinker) writeELFProgramHeader(header []byte, offset, size int) error {
	// Program header type: LOAD
	header[0] = 1 // PT_LOAD
	
	// Flags: executable and readable
	header[4] = 5 // PF_X | PF_R
	
	// File offset
	fileOffset := uint64(offset)
	for i := 0; i < 8; i++ {
		header[8+i] = byte(fileOffset >> (i * 8))
	}
	
	// Virtual address (standard load address)
	virtualAddr := uint64(0x400000 + offset)
	for i := 0; i < 8; i++ {
		header[16+i] = byte(virtualAddr >> (i * 8))
		header[24+i] = byte(virtualAddr >> (i * 8)) // Physical address same as virtual
	}
	
	// File size and memory size
	segmentSize := uint64(size)
	for i := 0; i < 8; i++ {
		header[32+i] = byte(segmentSize >> (i * 8)) // File size
		header[40+i] = byte(segmentSize >> (i * 8)) // Memory size
	}
	
	return nil
}

// linkPE creates a Windows PE executable
func (l *ExecutableLinker) linkPE(code *NativeCode, runtime *MinimalRuntime, options *CompileOptions) (*LinkedExecutable, error) {
	return nil, fmt.Errorf("PE format not yet implemented")
}

// RelocateSymbols adjusts symbol addresses for the final executable
func (l *ExecutableLinker) RelocateSymbols(symbols map[string]Symbol, baseAddress int64) map[string]Symbol {
	relocated := make(map[string]Symbol)
	for name, symbol := range symbols {
		symbol.Offset += int(baseAddress)
		relocated[name] = symbol
	}
	return relocated
}