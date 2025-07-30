package aot

import (
	"fmt"
)

// RuntimeBuilder creates minimal runtime for compiled Rush programs
type RuntimeBuilder struct {
	components []RuntimeComponent
}

// RuntimeComponent represents a part of the minimal runtime
type RuntimeComponent interface {
	Name() string
	GenerateCode() ([]byte, error)
	GetSymbols() map[string]Symbol
	IsRequired(options *CompileOptions) bool
}

// MinimalRuntime contains the essential runtime code
type MinimalRuntime struct {
	Code    []byte
	Symbols map[string]Symbol
	Size    int
}

// NewRuntimeBuilder creates a new runtime builder
func NewRuntimeBuilder() *RuntimeBuilder {
	return &RuntimeBuilder{
		components: []RuntimeComponent{
			NewStartupComponent(),
			NewMemoryManagerComponent(),
			NewErrorHandlerComponent(),
			NewBuiltinFunctionsComponent(),
			NewGarbageCollectorComponent(),
		},
	}
}

// BuildRuntime constructs the minimal runtime based on compile options
func (rb *RuntimeBuilder) BuildRuntime(options *CompileOptions) (*MinimalRuntime, error) {
	var allCode []byte
	allSymbols := make(map[string]Symbol)
	
	// Include only required components
	for _, component := range rb.components {
		if component.IsRequired(options) {
			code, err := component.GenerateCode()
			if err != nil {
				return nil, fmt.Errorf("failed to generate %s: %w", component.Name(), err)
			}
			
			// Adjust symbol offsets for concatenated code
			symbols := component.GetSymbols()
			for name, symbol := range symbols {
				symbol.Offset += len(allCode)
				allSymbols[name] = symbol
			}
			
			allCode = append(allCode, code...)
		}
	}
	
	return &MinimalRuntime{
		Code:    allCode,
		Symbols: allSymbols,
		Size:    len(allCode),
	}, nil
}

// Startup Component - Program initialization
type StartupComponent struct{}

func NewStartupComponent() *StartupComponent {
	return &StartupComponent{}
}

func (c *StartupComponent) Name() string {
	return "Startup"
}

func (c *StartupComponent) IsRequired(options *CompileOptions) bool {
	return true // Always required
}

func (c *StartupComponent) GenerateCode() ([]byte, error) {
	// ARM64 program startup code
	// This would initialize the runtime and call main
	code := []byte{
		// Simplified startup - in reality this would be much more complex
		0x00, 0x00, 0x80, 0xD2, // mov x0, #0    ; argc = 0
		0x01, 0x00, 0x80, 0xD2, // mov x1, #0    ; argv = NULL
		0x00, 0x00, 0x00, 0x94, // bl _main      ; Call main function
		0x00, 0x00, 0x80, 0xD2, // mov x0, #0    ; Exit code 0
		0x01, 0x00, 0x00, 0xD4, // svc #0        ; System call (exit)
	}
	return code, nil
}

func (c *StartupComponent) GetSymbols() map[string]Symbol {
	return map[string]Symbol{
		"_start": {
			Name:   "_start",
			Offset: 0,
			Size:   20, // 5 instructions * 4 bytes
			Type:   SymbolFunction,
		},
	}
}

// Memory Manager Component - Basic memory allocation
type MemoryManagerComponent struct{}

func NewMemoryManagerComponent() *MemoryManagerComponent {
	return &MemoryManagerComponent{}
}

func (c *MemoryManagerComponent) Name() string {
	return "MemoryManager"
}

func (c *MemoryManagerComponent) IsRequired(options *CompileOptions) bool {
	return true // Required for dynamic allocation
}

func (c *MemoryManagerComponent) GenerateCode() ([]byte, error) {
	// Simplified memory allocator
	// In a real implementation, this would be a proper malloc/free
	code := []byte{
		// rush_malloc function
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... actual malloc implementation would go here
		0xFD, 0x7B, 0xC1, 0xA8, // ldp x29, x30, [sp], #16
		0xC0, 0x03, 0x5F, 0xD6, // ret
		
		// rush_free function  
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... actual free implementation would go here
		0xFD, 0x7B, 0xC1, 0xA8, // ldp x29, x30, [sp], #16
		0xC0, 0x03, 0x5F, 0xD6, // ret
	}
	return code, nil
}

func (c *MemoryManagerComponent) GetSymbols() map[string]Symbol {
	return map[string]Symbol{
		"rush_malloc": {
			Name:   "rush_malloc",
			Offset: 0,
			Size:   16,
			Type:   SymbolFunction,
		},
		"rush_free": {
			Name:   "rush_free",
			Offset: 16,
			Size:   16,
			Type:   SymbolFunction,
		},
	}
}

// Error Handler Component - Runtime error handling
type ErrorHandlerComponent struct{}

func NewErrorHandlerComponent() *ErrorHandlerComponent {
	return &ErrorHandlerComponent{}
}

func (c *ErrorHandlerComponent) Name() string {
	return "ErrorHandler"
}

func (c *ErrorHandlerComponent) IsRequired(options *CompileOptions) bool {
	return true // Always needed for error handling
}

func (c *ErrorHandlerComponent) GenerateCode() ([]byte, error) {
	// Basic error handler that prints error and exits
	code := []byte{
		// rush_panic function
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... error printing code would go here
		0x01, 0x00, 0x80, 0xD2, // mov x1, #0    ; Exit code 1
		0x01, 0x00, 0x00, 0xD4, // svc #0        ; Exit syscall
	}
	return code, nil
}

func (c *ErrorHandlerComponent) GetSymbols() map[string]Symbol {
	return map[string]Symbol{
		"rush_panic": {
			Name:   "rush_panic",
			Offset: 0,
			Size:   20,
			Type:   SymbolFunction,
		},
	}
}

// Builtin Functions Component - Essential built-in functions
type BuiltinFunctionsComponent struct{}

func NewBuiltinFunctionsComponent() *BuiltinFunctionsComponent {
	return &BuiltinFunctionsComponent{}
}

func (c *BuiltinFunctionsComponent) Name() string {
	return "BuiltinFunctions"
}

func (c *BuiltinFunctionsComponent) IsRequired(options *CompileOptions) bool {
	return !options.MinimalRuntime // Skip if minimal runtime requested
}

func (c *BuiltinFunctionsComponent) GenerateCode() ([]byte, error) {
	// Compiled versions of essential built-in functions
	code := []byte{
		// rush_print function
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... print implementation using system calls
		0xFD, 0x7B, 0xC1, 0xA8, // ldp x29, x30, [sp], #16
		0xC0, 0x03, 0x5F, 0xD6, // ret
		
		// rush_len function
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... length calculation implementation
		0xFD, 0x7B, 0xC1, 0xA8, // ldp x29, x30, [sp], #16
		0xC0, 0x03, 0x5F, 0xD6, // ret
	}
	return code, nil
}

func (c *BuiltinFunctionsComponent) GetSymbols() map[string]Symbol {
	return map[string]Symbol{
		"rush_print": {
			Name:   "rush_print",
			Offset: 0,
			Size:   20,
			Type:   SymbolFunction,
		},
		"rush_len": {
			Name:   "rush_len",
			Offset: 20,
			Size:   20,
			Type:   SymbolFunction,
		},
	}
}

// Garbage Collector Component - Optional garbage collection
type GarbageCollectorComponent struct{}

func NewGarbageCollectorComponent() *GarbageCollectorComponent {
	return &GarbageCollectorComponent{}
}

func (c *GarbageCollectorComponent) Name() string {
	return "GarbageCollector"
}

func (c *GarbageCollectorComponent) IsRequired(options *CompileOptions) bool {
	return !options.MinimalRuntime && options.OptimizationLevel < 2
}

func (c *GarbageCollectorComponent) GenerateCode() ([]byte, error) {
	// Simple mark-and-sweep garbage collector
	code := []byte{
		// rush_gc_collect function
		0xFD, 0x7B, 0xBF, 0xA9, // stp x29, x30, [sp, #-16]!
		0xFD, 0x03, 0x00, 0x91, // mov x29, sp
		// ... GC implementation would go here
		0xFD, 0x7B, 0xC1, 0xA8, // ldp x29, x30, [sp], #16
		0xC0, 0x03, 0x5F, 0xD6, // ret
	}
	return code, nil
}

func (c *GarbageCollectorComponent) GetSymbols() map[string]Symbol {
	return map[string]Symbol{
		"rush_gc_collect": {
			Name:   "rush_gc_collect",
			Offset: 0,
			Size:   16,
			Type:   SymbolFunction,
		},
	}
}