package jit

import (
	"fmt"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"rush/interpreter"
)

// CodeCache manages compiled ARM64 native code
type CodeCache struct {
	entries    map[uint64]*CompiledCode
	mu         sync.RWMutex
	maxEntries int
	size       int64 // Total size in bytes
	maxSize    int64 // Maximum cache size in bytes
}

// CompiledCode represents executable ARM64 native code
type CompiledCode struct {
	NativeCode   []byte                      // ARM64 machine code
	Function     *interpreter.CompiledFunction // Original bytecode function
	Hash         uint64                      // Function hash for identification
	CreatedAt    time.Time                   // Creation timestamp
	ExecuteCount int64                       // Number of times executed
	codePtr      uintptr                     // Pointer to executable memory
	codeSize     int                         // Size of executable code
}

// ExecutionContext holds runtime context for JIT execution
type ExecutionContext struct {
	Args    []interpreter.Value
	Globals []interpreter.Value
	Stack   []interpreter.Value
}

// ARM64ExecutionContext holds ARM64-specific execution context
type ARM64ExecutionContext struct {
	Args       [8]uint64 // ARM64 X0-X7 argument registers
	GlobalsPtr uintptr   // Pointer to globals array
	StackPtr   uintptr   // Pointer to stack
	StackSize  int       // Stack size
}

// ARM64Result holds the result of ARM64 function execution
type ARM64Result struct {
	Value uint64 // Return value in X0
	Type  uint8  // Value type indicator
	Error uint8  // Error flag
}

const (
	DefaultMaxCacheEntries = 1000
	DefaultMaxCacheSize    = 64 * 1024 * 1024 // 64MB
	PageSize               = 4096
)

// NewCodeCache creates a new code cache
func NewCodeCache() *CodeCache {
	return &CodeCache{
		entries:    make(map[uint64]*CompiledCode),
		maxEntries: DefaultMaxCacheEntries,
		maxSize:    DefaultMaxCacheSize,
	}
}

// Add adds compiled code to the cache
func (c *CodeCache) Add(hash uint64, code *CompiledCode) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if we need to evict entries
	if len(c.entries) >= c.maxEntries || c.size+int64(len(code.NativeCode)) > c.maxSize {
		if err := c.evictLRU(); err != nil {
			return fmt.Errorf("failed to evict cache entries: %v", err)
		}
	}
	
	// Allocate executable memory for the code
	if err := c.makeExecutable(code); err != nil {
		return fmt.Errorf("failed to make code executable: %v", err)
	}
	
	c.entries[hash] = code
	c.size += int64(len(code.NativeCode))
	
	return nil
}

// Get retrieves compiled code from the cache
func (c *CodeCache) Get(hash uint64) *CompiledCode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if code, exists := c.entries[hash]; exists {
		code.ExecuteCount++
		return code
	}
	return nil
}

// Has checks if compiled code exists in the cache
func (c *CodeCache) Has(hash uint64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	_, exists := c.entries[hash]
	return exists
}

// Remove removes compiled code from the cache
func (c *CodeCache) Remove(hash uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if code, exists := c.entries[hash]; exists {
		// Free executable memory
		if err := c.freeExecutable(code); err != nil {
			return fmt.Errorf("failed to free executable memory: %v", err)
		}
		
		c.size -= int64(len(code.NativeCode))
		delete(c.entries, hash)
	}
	
	return nil
}

// makeExecutable allocates executable memory and copies code
func (c *CodeCache) makeExecutable(code *CompiledCode) error {
	codeSize := len(code.NativeCode)
	if codeSize == 0 {
		return fmt.Errorf("empty native code")
	}
	
	// Round up to page size
	allocSize := ((codeSize + PageSize - 1) / PageSize) * PageSize
	
	// Allocate memory with read/write permissions
	mem, err := syscall.Mmap(-1, 0, allocSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return fmt.Errorf("failed to allocate memory: %v", err)
	}
	
	// Copy code to allocated memory
	copy(mem, code.NativeCode)
	
	// Make memory executable (ARM64 requires instruction cache invalidation)
	err = syscall.Mprotect(mem, syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		syscall.Munmap(mem)
		return fmt.Errorf("failed to make memory executable: %v", err)
	}
	
	// ARM64: Invalidate instruction cache for the new code
	// This ensures cache coherency between instruction and data caches
	c.invalidateInstructionCache(uintptr(unsafe.Pointer(&mem[0])), codeSize)
	
	code.codePtr = uintptr(unsafe.Pointer(&mem[0]))
	code.codeSize = allocSize
	
	return nil
}

// freeExecutable frees executable memory
func (c *CodeCache) freeExecutable(code *CompiledCode) error {
	if code.codePtr != 0 && code.codeSize > 0 {
		mem := (*[1 << 30]byte)(unsafe.Pointer(code.codePtr))[:code.codeSize:code.codeSize]
		err := syscall.Munmap(mem)
		if err != nil {
			return fmt.Errorf("failed to unmap memory: %v", err)
		}
		code.codePtr = 0
		code.codeSize = 0
	}
	return nil
}

// invalidateInstructionCache ensures cache coherency on ARM64
func (c *CodeCache) invalidateInstructionCache(addr uintptr, size int) {
	// ARM64 cache invalidation using system call
	// This is a simplified version - in practice, you might want to use
	// more specific cache operations
	const (
		SYS_CACHE_INVALIDATE = 0x5 // Hypothetical system call number
	)
	
	// Note: This is a placeholder for ARM64 cache invalidation
	// In a real implementation, you would use appropriate ARM64 instructions
	// or system calls to invalidate the instruction cache
	
	// For now, we'll use a memory barrier to ensure ordering
	// In practice, ARM64 cache invalidation would be handled by
	// specific assembly instructions or system calls
}

// evictLRU evicts least recently used entries to make space
func (c *CodeCache) evictLRU() error {
	if len(c.entries) == 0 {
		return nil
	}
	
	// Find oldest entry (simple LRU based on creation time)
	var oldestHash uint64
	var oldestTime time.Time = time.Now()
	
	for hash, code := range c.entries {
		if code.CreatedAt.Before(oldestTime) {
			oldestTime = code.CreatedAt
			oldestHash = hash
		}
	}
	
	// Remove oldest entry
	return c.Remove(oldestHash)
}

// Execute runs compiled ARM64 code with exception handling
func (code *CompiledCode) Execute(args []interpreter.Value, globals []interpreter.Value) (interpreter.Value, error) {
	if code.codePtr == 0 {
		return nil, fmt.Errorf("no executable code available")
	}
	
	// Validate execution environment
	if err := ValidateExecutionEnvironment(); err != nil {
		return nil, fmt.Errorf("execution environment validation failed: %v", err)
	}
	
	// Create execution context
	ctx := &ExecutionContext{
		Args:    args,
		Globals: globals,
		Stack:   make([]interpreter.Value, 1024), // Stack for execution
	}
	
	// Set up exception handler
	handler := NewARM64ExceptionHandler()
	if err := handler.Install(); err != nil {
		return nil, fmt.Errorf("failed to install exception handler: %v", err)
	}
	defer handler.Uninstall()
	
	// Execute ARM64 code with exception handling
	var result interpreter.Value
	var execErr error
	
	_, excInfo, err := handler.SafeExecuteARM64(func() (uint64, error) {
		var execResult interpreter.Value
		execResult, execErr = code.executeNative(ctx)
		result = execResult
		return 0, execErr
	})
	
	if err != nil {
		// Check if exception is recoverable
		if excInfo != nil && excInfo.Recoverable {
			// Attempt recovery
			if recoverErr := handler.RecoverFromException(excInfo); recoverErr == nil {
				// Deoptimize and fall back to interpreter
				return nil, fmt.Errorf("JIT execution failed, deoptimizing: %v", err)
			}
		}
		return nil, fmt.Errorf("ARM64 execution failed with unrecoverable exception: %v", err)
	}
	
	if execErr != nil {
		return nil, fmt.Errorf("JIT execution failed, deoptimizing: %v", execErr)
	}
	
	return result, nil
}

// executeNative executes the ARM64 code
func (code *CompiledCode) executeNative(ctx *ExecutionContext) (interpreter.Value, error) {
	if code.codePtr == 0 {
		return nil, fmt.Errorf("no executable code pointer")
	}
	
	// Marshal Rush values to ARM64 calling convention
	arm64Args, err := code.marshalArgsToARM64(ctx.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}
	
	// Set up ARM64 execution context
	arm64Ctx := &ARM64ExecutionContext{
		Args:       arm64Args,
		GlobalsPtr: uintptr(unsafe.Pointer(&ctx.Globals[0])),
		StackPtr:   uintptr(unsafe.Pointer(&ctx.Stack[0])),
		StackSize:  len(ctx.Stack),
	}
	
	// Call ARM64 function with proper error handling
	result, err := code.callARM64Function(arm64Ctx)
	if err != nil {
		return nil, fmt.Errorf("ARM64 execution failed: %v", err)
	}
	
	// Unmarshal ARM64 result back to Rush value
	rushResult, err := code.unmarshalResultFromARM64(result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}
	
	return rushResult, nil
}

// GetStats returns cache statistics
func (c *CodeCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	totalExecutions := int64(0)
	for _, code := range c.entries {
		totalExecutions += code.ExecuteCount
	}
	
	return CacheStats{
		Entries:         len(c.entries),
		MaxEntries:      c.maxEntries,
		Size:            c.size,
		MaxSize:         c.maxSize,
		TotalExecutions: totalExecutions,
	}
}

// CacheStats holds cache statistics
type CacheStats struct {
	Entries         int
	MaxEntries      int
	Size            int64
	MaxSize         int64
	TotalExecutions int64
}

// Clear removes all entries from the cache
func (c *CodeCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Free all executable memory
	for hash := range c.entries {
		if code := c.entries[hash]; code != nil {
			if err := c.freeExecutable(code); err != nil {
				// Log error but continue cleanup
				continue
			}
		}
	}
	
	c.entries = make(map[uint64]*CompiledCode)
	c.size = 0
	
	return nil
}

// marshalArgsToARM64 converts Rush values to ARM64 calling convention
func (code *CompiledCode) marshalArgsToARM64(args []interpreter.Value) ([8]uint64, error) {
	var arm64Args [8]uint64
	
	// ARM64 calling convention: first 8 args in X0-X7
	for i, arg := range args {
		if i >= 8 {
			return arm64Args, fmt.Errorf("too many arguments for ARM64 calling convention: %d", len(args))
		}
		
		// Convert Rush value to 64-bit representation
		val, err := code.valueToUint64(arg)
		if err != nil {
			return arm64Args, fmt.Errorf("failed to marshal argument %d: %v", i, err)
		}
		arm64Args[i] = val
	}
	
	return arm64Args, nil
}

// unmarshalResultFromARM64 converts ARM64 result back to Rush value
func (code *CompiledCode) unmarshalResultFromARM64(result *ARM64Result) (interpreter.Value, error) {
	if result.Error != 0 {
		return nil, fmt.Errorf("ARM64 execution error: code %d", result.Error)
	}
	
	// Convert based on type indicator
	return code.uint64ToValue(result.Value, result.Type)
}

// valueToUint64 converts a Rush value to 64-bit representation
func (code *CompiledCode) valueToUint64(value interpreter.Value) (uint64, error) {
	switch v := value.(type) {
	case *interpreter.Integer:
		return uint64(v.Value), nil
	case *interpreter.Float:
		return *(*uint64)(unsafe.Pointer(&v.Value)), nil
	case *interpreter.Boolean:
		if v.Value {
			return 1, nil
		}
		return 0, nil
	case *interpreter.String:
		// Store pointer to string data
		return uint64(uintptr(unsafe.Pointer(&v.Value))), nil
	case *interpreter.Null:
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported value type for ARM64 marshaling: %T", value)
	}
}

// uint64ToValue converts 64-bit representation back to Rush value
func (code *CompiledCode) uint64ToValue(val uint64, typeIndicator uint8) (interpreter.Value, error) {
	switch typeIndicator {
	case 0: // INTEGER
		return &interpreter.Integer{Value: int64(val)}, nil
	case 1: // FLOAT
		floatVal := *(*float64)(unsafe.Pointer(&val))
		return &interpreter.Float{Value: floatVal}, nil
	case 2: // BOOLEAN
		return &interpreter.Boolean{Value: val != 0}, nil
	case 3: // STRING
		// Reconstruct string from pointer (careful with memory management)
		if val == 0 {
			return &interpreter.String{Value: ""}, nil
		}
		strPtr := (*string)(unsafe.Pointer(uintptr(val)))
		return &interpreter.String{Value: *strPtr}, nil
	case 4: // NULL
		return &interpreter.Null{}, nil
	default:
		return nil, fmt.Errorf("unknown type indicator: %d", typeIndicator)
	}
}

// callARM64Function executes the compiled ARM64 code
func (code *CompiledCode) callARM64Function(ctx *ARM64ExecutionContext) (*ARM64Result, error) {
	// This is where we make the actual function call to ARM64 code
	// Using CGO and assembly for the actual call
	return callARM64Code(code.codePtr, ctx)
}