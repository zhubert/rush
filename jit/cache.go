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

// Execute runs compiled ARM64 code
func (code *CompiledCode) Execute(args []interpreter.Value, globals []interpreter.Value) (interpreter.Value, error) {
	if code.codePtr == 0 {
		return nil, fmt.Errorf("no executable code available")
	}
	
	// Create execution context
	ctx := &ExecutionContext{
		Args:    args,
		Globals: globals,
		Stack:   make([]interpreter.Value, 1024), // Stack for execution
	}
	
	// Execute ARM64 code (this would be implemented in assembly)
	result, err := code.executeNative(ctx)
	if err != nil {
		// Deoptimization - fallback to interpreter
		return nil, fmt.Errorf("JIT execution failed, deoptimizing: %v", err)
	}
	
	return result, nil
}

// executeNative executes the ARM64 code
// This is a placeholder - actual implementation would involve calling
// the ARM64 code through function pointers and proper calling conventions
func (code *CompiledCode) executeNative(ctx *ExecutionContext) (interpreter.Value, error) {
	// This is a placeholder for actual ARM64 code execution
	// In a real implementation, this would:
	// 1. Set up ARM64 calling convention (X0-X7 for args)
	// 2. Call the compiled code at code.codePtr
	// 3. Handle the return value according to ARM64 ABI
	// 4. Convert ARM64 return values back to Rush values
	
	// For now, return an error to indicate deoptimization needed
	return nil, fmt.Errorf("ARM64 execution not yet implemented")
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