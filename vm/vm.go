package vm

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strings"
	"time"

	"rush/bytecode"
	"rush/compiler"
	"rush/interpreter"
	"rush/jit"
)

const (
	StackSize   = 2048 // Stack size for VM execution
	GlobalsSize = 65536 // Size of global variables storage
	MaxFrames   = 1024 // Maximum call frames
)

// LogLevel defines the verbosity of VM logging
type LogLevel int

const (
	LogNone LogLevel = iota  // No logging
	LogError                 // Only errors
	LogWarn                  // Warnings and errors
	LogInfo                  // General information
	LogDebug                 // Detailed debugging info
	LogTrace                 // Extremely verbose execution tracing
)

// VMLogger handles VM logging with configurable levels
type VMLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewVMLogger creates a new VM logger
func NewVMLogger(level LogLevel) *VMLogger {
	return &VMLogger{
		level:  level,
		logger: log.New(os.Stderr, "[VM] ", log.LstdFlags|log.Lmicroseconds),
	}
}

// Log methods for different levels
func (l *VMLogger) Error(format string, args ...interface{}) {
	if l.level >= LogError {
		l.logger.Printf("ERROR: "+format, args...)
	}
}

func (l *VMLogger) Warn(format string, args ...interface{}) {
	if l.level >= LogWarn {
		l.logger.Printf("WARN: "+format, args...)
	}
}

func (l *VMLogger) Info(format string, args ...interface{}) {
	if l.level >= LogInfo {
		l.logger.Printf("INFO: "+format, args...)
	}
}

func (l *VMLogger) Debug(format string, args ...interface{}) {
	if l.level >= LogDebug {
		l.logger.Printf("DEBUG: "+format, args...)
	}
}

func (l *VMLogger) Trace(format string, args ...interface{}) {
	if l.level >= LogTrace {
		l.logger.Printf("TRACE: "+format, args...)
	}
}

// Frame represents a call frame for function execution
type Frame struct {
	cl          *interpreter.Closure // Closure being executed
	ip          int                  // Instruction pointer
	basePointer int                  // Base pointer for local variables
	self        *interpreter.Object  // Current object context for instance variables
}

// NewFrame creates a new call frame
func NewFrame(cl *interpreter.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
		self:        nil, // No object context by default
	}
}

// NewFrameWithSelf creates a new call frame with object context
func NewFrameWithSelf(cl *interpreter.Closure, basePointer int, self *interpreter.Object) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
		self:        self,
	}
}

// Instructions returns the instructions for this frame
func (f *Frame) Instructions() bytecode.Instructions {
	return bytecode.Instructions(f.cl.Fn.Instructions)
}

// VM represents the virtual machine
type VM struct {
	constants    []interpreter.Value // Constant pool
	stack        []interpreter.Value // Execution stack
	sp           int                 // Stack pointer (points to next free slot)
	globals      []interpreter.Value // Global variables
	frames       []*Frame            // Call frames stack
	framesIndex  int                 // Current frame index
	logger       *VMLogger           // Logger for debugging and monitoring
	stats        *VMStats            // Execution statistics
	
	// JIT-specific fields
	jitCompiler  *jit.JITCompiler    // JIT compiler instance
	jitEnabled   bool                // Whether JIT compilation is enabled
}

// VMStats tracks execution statistics
type VMStats struct {
	StartTime         time.Time
	InstructionCount  int64
	StackOperations   int64
	FunctionCalls     int64
	MemoryAllocations int64
	Errors            int64
	
	// Per-function execution counters
	FunctionExecutions map[uint64]int64    // Function hash -> execution count
	FunctionTimings    map[uint64]time.Duration // Function hash -> total execution time
	
	// JIT-specific statistics
	JITCompilations    int64
	JITHits            int64
	JITMisses          int64
	JITDeoptimizations int64
	JITCompilationTime time.Duration
}

// New creates a new virtual machine
func New(bytecode *compiler.Bytecode) *VM {
	return NewWithLogger(bytecode, LogNone)
}

// NewWithLogger creates a new virtual machine with specified log level
func NewWithLogger(bytecode *compiler.Bytecode, logLevel LogLevel) *VM {
	mainFn := &interpreter.CompiledFunction{
		Instructions:  []byte(bytecode.Instructions),
		NumLocals:     0, // Main execution has no local variables
		NumParameters: 0, // Main execution has no parameters
	}
	mainClosure := &interpreter.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	logger := NewVMLogger(logLevel)
	stats := &VMStats{
		StartTime:          time.Now(),
		FunctionExecutions: make(map[uint64]int64),
		FunctionTimings:    make(map[uint64]time.Duration),
	}

	vm := &VM{
		constants:   bytecode.Constants,
		stack:       make([]interpreter.Value, StackSize),
		sp:          0,
		globals:     make([]interpreter.Value, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
		logger:      logger,
		stats:       stats,
		jitCompiler: nil,
		jitEnabled:  false,
	}

	logger.Info("VM initialized with %d constants, %d stack size, %d globals size", 
		len(bytecode.Constants), StackSize, GlobalsSize)
	logger.Debug("Main function has %d instructions", len(mainFn.Instructions))

	return vm
}

// NewWithGlobalsStore creates a VM with existing global state
func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []interpreter.Value) *VM {
	vm := New(bytecode)
	vm.globals = s
	vm.logger.Debug("VM initialized with %d existing global variables", len(s))
	return vm
}

// NewWithGlobalsStoreAndLogger creates a VM with existing global state and logger
func NewWithGlobalsStoreAndLogger(bytecode *compiler.Bytecode, s []interpreter.Value, logLevel LogLevel) *VM {
	vm := NewWithLogger(bytecode, logLevel)
	vm.globals = s
	vm.logger.Debug("VM initialized with %d existing global variables", len(s))
	return vm
}

// NewWithJIT creates a new virtual machine with JIT compilation enabled
func NewWithJIT(bytecode *compiler.Bytecode, logLevel LogLevel) *VM {
	vm := NewWithLogger(bytecode, logLevel)
	vm.jitCompiler = jit.NewJITCompiler()
	vm.jitEnabled = true
	vm.logger.Info("JIT compilation enabled")
	return vm
}

// NewWithJITAndGlobalsStore creates a VM with JIT compilation and existing global state
func NewWithJITAndGlobalsStore(bytecode *compiler.Bytecode, s []interpreter.Value) *VM {
	vm := NewWithJIT(bytecode, LogNone)
	vm.globals = s
	vm.logger.Debug("VM initialized with JIT and %d existing global variables", len(s))
	return vm
}

// SetLogLevel changes the logging level of the VM
func (vm *VM) SetLogLevel(level LogLevel) {
	vm.logger.level = level
	vm.logger.Info("Log level changed to %v", level)
}

// GetStats returns a copy of the current execution statistics
func (vm *VM) GetStats() VMStats {
	stats := *vm.stats
	stats.InstructionCount = vm.stats.InstructionCount
	return stats
}

// GetJITStats returns JIT compilation statistics
func (vm *VM) GetJITStats() jit.JITStats {
	if vm.jitEnabled && vm.jitCompiler != nil {
		return vm.jitCompiler.GetStats()
	}
	return jit.JITStats{} // Return empty stats if JIT is disabled
}

// RecordFunctionExecution tracks function execution for profiling
func (vm *VM) RecordFunctionExecution(fnHash uint64, executionTime time.Duration) {
	vm.stats.FunctionExecutions[fnHash]++
	vm.stats.FunctionTimings[fnHash] += executionTime
	
	// Also record in JIT profiler if enabled
	if vm.jitEnabled && vm.jitCompiler != nil {
		vm.jitCompiler.RecordExecution(fnHash, executionTime)
	}
}

// GetFunctionStats returns execution statistics for a specific function
func (vm *VM) GetFunctionStats(fnHash uint64) (count int64, totalTime time.Duration) {
	count = vm.stats.FunctionExecutions[fnHash]
	totalTime = vm.stats.FunctionTimings[fnHash]
	return
}

// UpdateJITStats updates JIT-related statistics
func (vm *VM) UpdateJITStats(compilations, hits, misses, deoptimizations int64, compilationTime time.Duration) {
	vm.stats.JITCompilations += compilations
	vm.stats.JITHits += hits
	vm.stats.JITMisses += misses
	vm.stats.JITDeoptimizations += deoptimizations
	vm.stats.JITCompilationTime += compilationTime
}

// generateFunctionHash creates a unique hash for a compiled function
func (vm *VM) generateFunctionHash(fn *interpreter.CompiledFunction) uint64 {
	hasher := fnv.New64a()
	
	// Hash the instructions
	hasher.Write(fn.Instructions)
	
	// Hash function metadata
	binary.Write(hasher, binary.LittleEndian, int32(fn.NumLocals))
	binary.Write(hasher, binary.LittleEndian, int32(fn.NumParameters))
	
	return hasher.Sum64()
}

// PrintStats logs current execution statistics
func (vm *VM) PrintStats() {
	elapsed := time.Since(vm.stats.StartTime)
	vm.logger.Info("=== VM EXECUTION STATS ===")
	vm.logger.Info("Execution time: %v", elapsed)
	vm.logger.Info("Instructions executed: %d", vm.stats.InstructionCount)
	vm.logger.Info("Stack operations: %d", vm.stats.StackOperations)
	vm.logger.Info("Function calls: %d", vm.stats.FunctionCalls)
	vm.logger.Info("Memory allocations: %d", vm.stats.MemoryAllocations)
	vm.logger.Info("Errors encountered: %d", vm.stats.Errors)
	
	// Print per-function statistics
	if len(vm.stats.FunctionExecutions) > 0 {
		vm.logger.Info("Function executions: %d unique functions", len(vm.stats.FunctionExecutions))
	}
	
	// Print JIT statistics if enabled
	if vm.jitEnabled {
		vm.logger.Info("=== JIT STATS ===")
		vm.logger.Info("JIT compilations: %d", vm.stats.JITCompilations)
		vm.logger.Info("JIT hits: %d", vm.stats.JITHits)
		vm.logger.Info("JIT misses: %d", vm.stats.JITMisses)
		vm.logger.Info("JIT deoptimizations: %d", vm.stats.JITDeoptimizations)
		vm.logger.Info("JIT compilation time: %v", vm.stats.JITCompilationTime)
		
		if (vm.stats.JITHits + vm.stats.JITMisses) > 0 {
			hitRate := float64(vm.stats.JITHits) / float64(vm.stats.JITHits + vm.stats.JITMisses) * 100
			vm.logger.Info("JIT hit rate: %.2f%%", hitRate)
		}
	}
	
	if elapsed.Nanoseconds() > 0 {
		ips := float64(vm.stats.InstructionCount) / elapsed.Seconds()
		vm.logger.Info("Instructions per second: %.2f", ips)
	}
}

// StackTop returns the top element of the stack
func (vm *VM) StackTop() interpreter.Value {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// Run executes the bytecode instructions
func (vm *VM) Run() error {
	vm.logger.Info("Starting VM execution")
	defer func() {
		vm.logger.Info("VM execution completed")
		if vm.logger.level >= LogInfo {
			vm.PrintStats()
		}
	}()

	var ip int
	var ins bytecode.Instructions
	var op bytecode.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++
		vm.stats.InstructionCount++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = bytecode.Opcode(ins[ip])

		vm.logger.Trace("IP:%d OP:%s SP:%d Frame:%d", ip, vm.getOpcodeName(op), vm.sp, vm.framesIndex-1)

		switch op {
		case bytecode.OpConstant:
			constIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			vm.logger.Debug("Loading constant[%d]: %s", constIndex, vm.constants[constIndex].Inspect())
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				vm.logger.Error("Failed to push constant: %v", err)
				vm.stats.Errors++
				return err
			}

		case bytecode.OpPop:
			popped := vm.pop()
			vm.logger.Debug("Popped: %s", popped.Inspect())

		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod:
			vm.logger.Debug("Executing binary operation: %s", vm.getOpcodeName(op))
			err := vm.executeBinaryOperation(op)
			if err != nil {
				vm.logger.Error("Binary operation failed: %v", err)
				vm.stats.Errors++
				return err
			}

		case bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGreaterThan, bytecode.OpGreaterEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case bytecode.OpAnd, bytecode.OpOr:
			err := vm.executeLogicalOperation(op)
			if err != nil {
				return err
			}

		case bytecode.OpTrue:
			err := vm.push(interpreter.TRUE)
			if err != nil {
				return err
			}

		case bytecode.OpFalse:
			err := vm.push(interpreter.FALSE)
			if err != nil {
				return err
			}

		case bytecode.OpNull:
			err := vm.push(interpreter.NULL)
			if err != nil {
				return err
			}

		case bytecode.OpNot:
			err := vm.executeNotOperation()
			if err != nil {
				return err
			}

		case bytecode.OpMinus:
			err := vm.executeMinusOperation()
			if err != nil {
				return err
			}

		case bytecode.OpJump:
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.logger.Debug("Jumping to position %d", pos)
			vm.currentFrame().ip = pos - 1

		case bytecode.OpJumpNotTruthy:
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !interpreter.IsTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case bytecode.OpJumpTruthy:
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if interpreter.IsTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case bytecode.OpSetGlobal:
			globalIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()

		case bytecode.OpGetGlobal:
			globalIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case bytecode.OpSetLocal:
			localIndex := int(ins[ip+1])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()
			value := vm.pop()
			vm.stack[frame.basePointer+localIndex] = value

		case bytecode.OpGetLocal:
			localIndex := int(ins[ip+1])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()
			value := vm.stack[frame.basePointer+localIndex]
			err := vm.push(value)
			if err != nil {
				return err
			}

		case bytecode.OpArray:
			numElements := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.safeSetSP(vm.sp - numElements)

			err := vm.push(array)
			if err != nil {
				return err
			}

		case bytecode.OpHash:
			numPairs := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			numElements := numPairs * 2
			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.safeSetSP(vm.sp - numElements)

			err = vm.push(hash)
			if err != nil {
				return err
			}

		case bytecode.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case bytecode.OpSetIndex:
			value := vm.pop()
			index := vm.pop()
			left := vm.pop()

			err := vm.executeSetIndexExpression(left, index, value)
			if err != nil {
				return err
			}

		case bytecode.OpCall:
			numArgs := int(ins[ip+1])
			vm.currentFrame().ip += 1

			vm.logger.Debug("Calling function with %d arguments", numArgs)
			vm.stats.FunctionCalls++
			err := vm.executeCall(numArgs)
			if err != nil {
				vm.logger.Error("Function call failed: %v", err)
				vm.stats.Errors++
				return err
			}

		case bytecode.OpReturn:
			returnValue := vm.pop()
			vm.logger.Debug("Returning value: %s", returnValue.Inspect())

			frame := vm.popFrame()
			vm.logger.Debug("Popped frame, returning to frame %d", vm.framesIndex-1)
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				vm.logger.Error("Failed to push return value: %v", err)
				vm.stats.Errors++
				return err
			}

		case bytecode.OpReturnVoid:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(interpreter.NULL)
			if err != nil {
				return err
			}

		case bytecode.OpGetBuiltin:
			builtinIndex := int(ins[ip+1])
			vm.currentFrame().ip += 1

			definition := interpreter.Builtins[builtinIndex]
			err := vm.push(&interpreter.BuiltinFunction{
				Fn: func(args ...interpreter.Value) interpreter.Value {
					// Get the actual builtin function
					if builtin, ok := interpreter.GetBuiltin(definition); ok {
						return builtin.Fn(args...)
					}
					return interpreter.NULL
				},
			})
			if err != nil {
				return err
			}

		case bytecode.OpClosure:
			constIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			numFree := int(ins[ip+3])
			vm.currentFrame().ip += 3

			err := vm.pushClosure(constIndex, numFree)
			if err != nil {
				return err
			}

		case bytecode.OpGetFree:
			freeIndex := int(ins[ip+1])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure.Free[freeIndex])
			if err != nil {
				return err
			}

		case bytecode.OpSetFree:
			freeIndex := int(ins[ip+1])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().cl
			value := vm.pop()
			currentClosure.Free[freeIndex] = value

		case bytecode.OpCurrentClosure:
			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure)
			if err != nil {
				return err
			}

		case bytecode.OpThrow:
			// Pop the exception value from stack
			exception := vm.pop()
			// For now, just return an error - proper exception handling needs a try/catch stack
			return fmt.Errorf("exception thrown: %s", exception.Inspect())

		case bytecode.OpTryBegin:
			// For now, just skip the try/catch mechanics
			// In a full implementation, this would set up exception handlers
			catchOffset := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			_ = catchOffset // Skip for now

		case bytecode.OpTryEnd:
			// Mark end of try block
			// In a full implementation, this would clean up exception handlers

		case bytecode.OpCatch:
			// For now, just skip catch handling
			errorTypeIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			_ = errorTypeIndex // Skip for now

		case bytecode.OpFinally:
			// Finally block - just continue execution

		case bytecode.OpImport:
			moduleIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			moduleName := vm.constants[moduleIndex].(*interpreter.String).Value
			// For now, just push a placeholder module object
			// In a full implementation, this would load the actual module
			moduleObj := &interpreter.Hash{
				Pairs: make(map[interpreter.HashKey]interpreter.Value),
				Keys:  make([]interpreter.Value, 0),
			}
			err := vm.push(moduleObj)
			if err != nil {
				return err
			}
			_ = moduleName // Use in full implementation

		case bytecode.OpExport:
			exportIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			exportName := vm.constants[exportIndex].(*interpreter.String).Value
			exportValue := vm.pop()
			// For now, just continue - in full implementation would register export
			_ = exportName
			_ = exportValue

		case bytecode.OpGetProperty:
			propertyIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			object := vm.pop()
			propertyName := vm.constants[propertyIndex].(*interpreter.String).Value
			
			err := vm.executePropertyAccess(object, propertyName)
			if err != nil {
				return err
			}

		case bytecode.OpClass:
			nameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			methodCount := int(ins[ip+3])
			vm.currentFrame().ip += 3
			
			className := vm.constants[nameIndex].(*interpreter.String).Value
			
			// Create new class
			class := &interpreter.Class{
				Name:            className,
				Methods:         make(map[string]*interpreter.Function), // For interpreter compatibility
				CompiledMethods: make(map[string]*interpreter.CompiledFunction),
				SuperClass:      nil,
			}
			
			err := vm.push(class)
			if err != nil {
				return err
			}
			
			// The methodCount parameter is for future use
			_ = methodCount

		case bytecode.OpInherit:
			// Pop superclass and current class
			superClass := vm.pop()
			currentClass := vm.pop()
			
			// Set inheritance
			if class, ok := currentClass.(*interpreter.Class); ok {
				if super, ok := superClass.(*interpreter.Class); ok {
					class.SuperClass = super
				} else {
					return fmt.Errorf("superclass must be a class, got %T", superClass)
				}
			} else {
				return fmt.Errorf("inheritance target must be a class, got %T", currentClass)
			}
			
			// Push class back on stack
			err := vm.push(currentClass)
			if err != nil {
				return err
			}

		case bytecode.OpMethod:
			methodNameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			methodName := vm.constants[methodNameIndex].(*interpreter.String).Value
			
			// Get closure and class from stack
			closure := vm.pop()
			currentClass := vm.pop()
			
			if class, ok := currentClass.(*interpreter.Class); ok {
				if closureObj, ok := closure.(*interpreter.Closure); ok {
					// Store compiled method in class
					class.CompiledMethods[methodName] = closureObj.Fn
				} else {
					return fmt.Errorf("method must be a closure, got %T", closure)
				}
			} else {
				return fmt.Errorf("method definition target must be a class, got %T", currentClass)
			}
			
			// Push class back on stack
			err := vm.push(currentClass)
			if err != nil {
				return err
			}

		case bytecode.OpInvoke:
			methodNameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			numArgs := int(ins[ip+3])
			vm.currentFrame().ip += 3
			
			methodName := vm.constants[methodNameIndex].(*interpreter.String).Value
			
			// Get object and arguments from stack
			args := make([]interpreter.Value, numArgs)
			for i := numArgs - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			object := vm.pop()
			
			err := vm.executeMethodCall(object, methodName, args)
			if err != nil {
				return err
			}

		case bytecode.OpGetInstance:
			varNameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			varName := vm.constants[varNameIndex].(*interpreter.String).Value
			
			// Get current object context from frame
			currentFrame := vm.currentFrame()
			if currentFrame.self == nil {
				return fmt.Errorf("instance variable @%s used outside of object context", varName)
			}
			
			// Get instance variable from the object
			if value, exists := currentFrame.self.InstanceVars[varName]; exists {
				err := vm.push(value)
				if err != nil {
					return err
				}
			} else {
				// Instance variable not set yet, return NULL
				err := vm.push(interpreter.NULL)
				if err != nil {
					return err
				}
			}

		case bytecode.OpSetInstance:
			varNameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			varName := vm.constants[varNameIndex].(*interpreter.String).Value
			value := vm.pop()
			
			// Get current object context from frame
			currentFrame := vm.currentFrame()
			if currentFrame.self == nil {
				return fmt.Errorf("instance variable @%s assigned outside of object context", varName)
			}
			
			// Set instance variable on the object
			currentFrame.self.InstanceVars[varName] = value

		case bytecode.OpGetSuper:
			methodNameIndex := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			
			methodName := vm.constants[methodNameIndex].(*interpreter.String).Value
			
			// Get super method - needs proper class context tracking
			_ = methodName
			return fmt.Errorf("super method access requires class context")

		default:
			return fmt.Errorf("unknown opcode: %d", op)
		}
	}

	return nil
}

// Helper methods

func (vm *VM) push(o interpreter.Value) error {
	if vm.sp < 0 {
		vm.logger.Error("Stack pointer negative before push: sp=%d", vm.sp)
		panic(fmt.Sprintf("stack pointer negative before push: sp=%d", vm.sp))
	}
	if vm.sp >= StackSize {
		vm.logger.Error("Stack overflow: sp=%d, max=%d", vm.sp, StackSize)
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	vm.stats.StackOperations++

	vm.logger.Trace("Pushed: %s (SP now %d)", o.Inspect(), vm.sp)
	return nil
}

func (vm *VM) pop() interpreter.Value {
	if vm.sp <= 0 {
		vm.logger.Error("Stack underflow: sp=%d", vm.sp)
		panic(fmt.Sprintf("stack underflow: sp=%d", vm.sp))
	}
	o := vm.stack[vm.sp-1]
	vm.sp--
	vm.stats.StackOperations++

	vm.logger.Trace("Popped: %s (SP now %d)", o.Inspect(), vm.sp)
	return o
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

// safeSetSP sets the stack pointer, ensuring it doesn't go below local variable space
func (vm *VM) safeSetSP(newSP int) {
	frame := vm.currentFrame()
	minSP := frame.basePointer
	
	// Only apply local variable constraints if we're in a closure (function frame)
	// For the main execution frame (cl == nil), we can set SP more freely
	if frame.cl != nil {
		minSP += frame.cl.Fn.NumLocals
		if newSP < minSP {
			newSP = minSP
		}
	} else {
		// In main execution frame, just ensure we don't go below base pointer
		if newSP < minSP {
			newSP = minSP
		}
	}
	
	vm.sp = newSP
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
	vm.stats.MemoryAllocations++
	vm.logger.Debug("Pushed frame %d", vm.framesIndex-1)
}

func (vm *VM) popFrame() *Frame {
	frame := vm.frames[vm.framesIndex-1]
	vm.framesIndex--
	vm.logger.Debug("Popped frame, now at frame %d", vm.framesIndex-1)
	return frame
}

func (vm *VM) executeBinaryOperation(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == interpreter.INTEGER_VALUE && rightType == interpreter.INTEGER_VALUE:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == interpreter.FLOAT_VALUE && rightType == interpreter.FLOAT_VALUE:
		return vm.executeBinaryFloatOperation(op, left, right)
	case leftType == interpreter.INTEGER_VALUE && rightType == interpreter.FLOAT_VALUE:
		return vm.executeBinaryMixedNumberOperation(op, left, right)
	case leftType == interpreter.FLOAT_VALUE && rightType == interpreter.INTEGER_VALUE:
		return vm.executeBinaryMixedNumberOperation(op, left, right)
	case leftType == interpreter.STRING_VALUE && rightType == interpreter.STRING_VALUE:
		return vm.executeBinaryStringOperation(op, left, right)
	case leftType == interpreter.STRING_VALUE || rightType == interpreter.STRING_VALUE || leftType == interpreter.BUILTIN_VALUE || rightType == interpreter.BUILTIN_VALUE:
		if op == bytecode.OpAdd {
			return vm.executeBinaryStringCoercionOperation(op, left, right)
		}
		// For non-addition operations with builtin, fall through to error
		leftTypeName := vm.getTypeName(leftType)
		rightTypeName := vm.getTypeName(rightType) 
		opName := vm.getOperatorName(op)
		return fmt.Errorf("unknown operator: %s %s %s", leftTypeName, opName, rightTypeName)
	default:
		leftTypeName := vm.getTypeName(leftType)
		rightTypeName := vm.getTypeName(rightType) 
		opName := vm.getOperatorName(op)
		return fmt.Errorf("unknown operator: %s %s %s", leftTypeName, opName, rightTypeName)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op bytecode.Opcode, left, right interpreter.Value) error {
	leftVal := left.(*interpreter.Integer).Value
	rightVal := right.(*interpreter.Integer).Value

	var result int64
	switch op {
	case bytecode.OpAdd:
		result = leftVal + rightVal
	case bytecode.OpSub:
		result = leftVal - rightVal
	case bytecode.OpMul:
		result = leftVal * rightVal
	case bytecode.OpDiv:
		if rightVal == 0 {
			return fmt.Errorf("division by zero")
		}
		result = leftVal / rightVal
	case bytecode.OpMod:
		if rightVal == 0 {
			return fmt.Errorf("division by zero")
		}
		result = leftVal % rightVal
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&interpreter.Integer{Value: result})
}

func (vm *VM) executeBinaryFloatOperation(op bytecode.Opcode, left, right interpreter.Value) error {
	leftVal := left.(*interpreter.Float).Value
	rightVal := right.(*interpreter.Float).Value

	var result float64
	switch op {
	case bytecode.OpAdd:
		result = leftVal + rightVal
	case bytecode.OpSub:
		result = leftVal - rightVal
	case bytecode.OpMul:
		result = leftVal * rightVal
	case bytecode.OpDiv:
		if rightVal == 0 {
			return fmt.Errorf("division by zero")
		}
		result = leftVal / rightVal
	default:
		return fmt.Errorf("unknown float operator: %d", op)
	}

	return vm.push(&interpreter.Float{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op bytecode.Opcode, left, right interpreter.Value) error {
	leftVal := left.(*interpreter.String).Value
	rightVal := right.(*interpreter.String).Value

	var result string
	switch op {
	case bytecode.OpAdd:
		result = leftVal + rightVal
	default:
		opName := vm.getOperatorName(op)
		return fmt.Errorf("unknown operator: STRING %s STRING", opName)
	}

	return vm.push(&interpreter.String{Value: result})
}

func (vm *VM) executeComparison(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == interpreter.INTEGER_VALUE && right.Type() == interpreter.INTEGER_VALUE {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBoolToPushBool(right == left))
	case bytecode.OpNotEqual:
		return vm.push(nativeBoolToPushBool(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%T %T)", op, left, right)
	}
}

func (vm *VM) executeIntegerComparison(op bytecode.Opcode, left, right interpreter.Value) error {
	leftVal := left.(*interpreter.Integer).Value
	rightVal := right.(*interpreter.Integer).Value

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBoolToPushBool(rightVal == leftVal))
	case bytecode.OpNotEqual:
		return vm.push(nativeBoolToPushBool(rightVal != leftVal))
	case bytecode.OpGreaterThan:
		return vm.push(nativeBoolToPushBool(leftVal > rightVal))
	case bytecode.OpGreaterEqual:
		return vm.push(nativeBoolToPushBool(leftVal >= rightVal))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeLogicalOperation(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	switch op {
	case bytecode.OpAnd:
		if !interpreter.IsTruthy(left) {
			return vm.push(left)
		}
		return vm.push(right)
	case bytecode.OpOr:
		if interpreter.IsTruthy(left) {
			return vm.push(left)
		}
		return vm.push(right)
	default:
		return fmt.Errorf("unknown logical operator: %d", op)
	}
}

func (vm *VM) executeNotOperation() error {
	operand := vm.pop()

	switch operand {
	case interpreter.TRUE:
		return vm.push(interpreter.FALSE)
	case interpreter.FALSE:
		return vm.push(interpreter.TRUE)
	case interpreter.NULL:
		return vm.push(interpreter.TRUE)
	default:
		return vm.push(interpreter.FALSE)
	}
}

func (vm *VM) executeMinusOperation() error {
	operand := vm.pop()

	switch operand := operand.(type) {
	case *interpreter.Integer:
		return vm.push(&interpreter.Integer{Value: -operand.Value})
	case *interpreter.Float:
		return vm.push(&interpreter.Float{Value: -operand.Value})
	default:
		typeName := vm.getTypeName(operand.Type())
		return fmt.Errorf("unknown operator: -%s", typeName)
	}
}

func (vm *VM) buildArray(startIndex, endIndex int) interpreter.Value {
	elements := make([]interpreter.Value, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &interpreter.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (interpreter.Value, error) {
	hashedPairs := make(map[interpreter.HashKey]interpreter.Value)
	keys := make([]interpreter.Value, 0)

	for i := startIndex; i < endIndex; i = i + 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		// Check if key is hashable
		switch key.(type) {
		case *interpreter.Integer, *interpreter.String, *interpreter.Boolean, *interpreter.Float:
			// Valid hash key
		default:
			typeName := vm.getTypeName(key.Type())
		return nil, fmt.Errorf("unusable as hash key: %s", typeName)
		}

		hashed := interpreter.CreateHashKey(key)
		if _, exists := hashedPairs[hashed]; !exists {
			keys = append(keys, key)
		}
		hashedPairs[hashed] = value
	}

	return &interpreter.Hash{Pairs: hashedPairs, Keys: keys}, nil
}

func (vm *VM) executeIndexExpression(left, index interpreter.Value) error {
	switch {
	case left.Type() == interpreter.ARRAY_VALUE && index.Type() == interpreter.INTEGER_VALUE:
		return vm.executeArrayIndex(left, index)
	case left.Type() == interpreter.STRING_VALUE && index.Type() == interpreter.INTEGER_VALUE:
		return vm.executeStringIndex(left, index)
	case left.Type() == interpreter.HASH_VALUE:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %T", left)
	}
}

func (vm *VM) executeArrayIndex(array, index interpreter.Value) error {
	arrayObject := array.(*interpreter.Array)
	i := index.(*interpreter.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(interpreter.NULL)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeStringIndex(str, index interpreter.Value) error {
	stringObject := str.(*interpreter.String)
	i := index.(*interpreter.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if i < 0 || i > max {
		return fmt.Errorf("IndexError: string index %d out of range [0:%d]", i, max+1)
	}

	return vm.push(&interpreter.String{Value: string(stringObject.Value[i])})
}

func (vm *VM) executeHashIndex(hash, index interpreter.Value) error {
	hashObject := hash.(*interpreter.Hash)

	// Check if index is hashable
	switch index.(type) {
	case *interpreter.Integer, *interpreter.String, *interpreter.Boolean, *interpreter.Float:
		// Valid hash key
	default:
		typeName := vm.getTypeName(index.Type())
		return fmt.Errorf("unusable as hash key: %s", typeName)
	}

	hashKey := interpreter.CreateHashKey(index)
	value, ok := hashObject.Pairs[hashKey]
	if !ok {
		return vm.push(interpreter.NULL)
	}

	return vm.push(value)
}

func (vm *VM) executeSetIndexExpression(left, index, value interpreter.Value) error {
	switch {
	case left.Type() == interpreter.ARRAY_VALUE && index.Type() == interpreter.INTEGER_VALUE:
		return vm.executeArraySetIndex(left, index, value)
	case left.Type() == interpreter.HASH_VALUE:
		return vm.executeHashSetIndex(left, index, value)
	default:
		return fmt.Errorf("set index operator not supported: %T", left)
	}
}

func (vm *VM) executeArraySetIndex(array, index, value interpreter.Value) error {
	arrayObject := array.(*interpreter.Array)
	i := index.(*interpreter.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return fmt.Errorf("array index out of bounds: %d", i)
	}

	arrayObject.Elements[i] = value
	return vm.push(value) // Return the assigned value
}

func (vm *VM) executeHashSetIndex(hash, index, value interpreter.Value) error {
	hashObject := hash.(*interpreter.Hash)

	// Check if index is hashable
	switch index.(type) {
	case *interpreter.Integer, *interpreter.String, *interpreter.Boolean, *interpreter.Float:
		// Valid hash key
	default:
		typeName := vm.getTypeName(index.Type())
		return fmt.Errorf("unusable as hash key: %s", typeName)
	}

	hashKey := interpreter.CreateHashKey(index)
	
	// Add to keys if it's a new key
	if _, exists := hashObject.Pairs[hashKey]; !exists {
		hashObject.Keys = append(hashObject.Keys, index)
	}
	
	hashObject.Pairs[hashKey] = value
	return vm.push(value) // Return the assigned value
}

func (vm *VM) executePropertyAccess(object interpreter.Value, propertyName string) error {
	switch obj := object.(type) {
	case *interpreter.String:
		return vm.executeStringProperty(obj, propertyName)
	case *interpreter.Array:
		return vm.executeArrayProperty(obj, propertyName)
	case *interpreter.Hash:
		return vm.executeHashProperty(obj, propertyName)
	case *interpreter.Integer:
		return vm.executeNumberProperty(obj, propertyName)
	case *interpreter.Float:
		return vm.executeNumberProperty(obj, propertyName)
	case *interpreter.Object:
		return vm.executeObjectProperty(obj, propertyName)
	case *interpreter.BuiltinFunction:
		return vm.executeBuiltinFunctionProperty(obj, propertyName)
	case *interpreter.JSON:
		return vm.executeJSONProperty(obj, propertyName)
	case *interpreter.Regexp:
		return vm.executeRegexpProperty(obj, propertyName)
	case *interpreter.Error:
		// Errors don't have properties, just return the error itself
		return fmt.Errorf("cannot access property on error: %s", obj.Message)
	default:
		return fmt.Errorf("property access not supported for type: %T", object)
	}
}

func (vm *VM) executeStringProperty(str *interpreter.String, propertyName string) error {
	switch propertyName {
	case "length":
		return vm.push(&interpreter.Integer{Value: int64(len(str.Value))})
	case "upper":
		// Return a bound method for upper()
		return vm.push(&interpreter.StringMethod{String: str, Method: "upper"})
	case "lower":
		return vm.push(&interpreter.StringMethod{String: str, Method: "lower"})
	case "trim":
		return vm.push(&interpreter.StringMethod{String: str, Method: "trim"})
	case "split":
		return vm.push(&interpreter.StringMethod{String: str, Method: "split"})
	case "substring":
		return vm.push(&interpreter.StringMethod{String: str, Method: "substring"})
	case "contains":
		return vm.push(&interpreter.StringMethod{String: str, Method: "contains"})
	case "starts_with":
		return vm.push(&interpreter.StringMethod{String: str, Method: "starts_with"})
	case "ends_with":
		return vm.push(&interpreter.StringMethod{String: str, Method: "ends_with"})
	case "match":
		return vm.push(&interpreter.StringMethod{String: str, Method: "match"})
	case "matches?":
		return vm.push(&interpreter.StringMethod{String: str, Method: "matches?"})
	case "replace":
		return vm.push(&interpreter.StringMethod{String: str, Method: "replace"})
	default:
		return fmt.Errorf("unknown property '%s' for string", propertyName)
	}
}

func (vm *VM) executeArrayProperty(arr *interpreter.Array, propertyName string) error {
	switch propertyName {
	case "length":
		return vm.push(&interpreter.Integer{Value: int64(len(arr.Elements))})
	case "push":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "push"})
	case "pop":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "pop"})
	case "shift":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "shift"})
	case "unshift":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "unshift"})
	case "slice":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "slice"})
	case "join":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "join"})
	case "reverse":
		return vm.push(&interpreter.ArrayMethod{Array: arr, Method: "reverse"})
	case "first":
		if len(arr.Elements) > 0 {
			return vm.push(arr.Elements[0])
		}
		return vm.push(interpreter.NULL)
	case "last":
		if len(arr.Elements) > 0 {
			return vm.push(arr.Elements[len(arr.Elements)-1])
		}
		return vm.push(interpreter.NULL)
	default:
		return fmt.Errorf("unknown property '%s' for array", propertyName)
	}
}

func (vm *VM) executeHashProperty(hash *interpreter.Hash, propertyName string) error {
	switch propertyName {
	case "length", "size":
		return vm.push(&interpreter.Integer{Value: int64(len(hash.Keys))})
	case "keys":
		return vm.push(&interpreter.Array{Elements: hash.Keys})
	case "values":
		values := make([]interpreter.Value, 0, len(hash.Keys))
		for _, key := range hash.Keys {
			hashKey := interpreter.CreateHashKey(key)
			values = append(values, hash.Pairs[hashKey])
		}
		return vm.push(&interpreter.Array{Elements: values})
	case "empty":
		return vm.push(&interpreter.Boolean{Value: len(hash.Keys) == 0})
	case "has_key":
		return vm.push(&interpreter.HashMethod{Hash: hash, Method: "has_key"})
	case "get":
		return vm.push(&interpreter.HashMethod{Hash: hash, Method: "get"})
	case "set":
		return vm.push(&interpreter.HashMethod{Hash: hash, Method: "set"})
	case "delete":
		return vm.push(&interpreter.HashMethod{Hash: hash, Method: "delete"})
	default:
		return fmt.Errorf("unknown property '%s' for hash", propertyName)
	}
}

func (vm *VM) executeNumberProperty(num interpreter.Value, propertyName string) error {
	switch propertyName {
	case "abs":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "abs"})
	case "floor":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "floor"})
	case "ceil":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "ceil"})
	case "round":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "round"})
	case "sqrt":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "sqrt"})
	case "pow":
		return vm.push(&interpreter.NumberMethod{Number: num, Method: "pow"})
	default:
		return fmt.Errorf("unknown property '%s' for number", propertyName)
	}
}

func (vm *VM) executeObjectProperty(obj *interpreter.Object, propertyName string) error {
	// Check if the method exists in the class
	class := obj.Class
	if method, ok := class.CompiledMethods[propertyName]; ok {
		// Create a bound method-like structure for the object
		// For now, we'll push a closure that calls the method with the object context
		closure := &interpreter.Closure{Fn: method}
		
		// We need to create a bound method that carries the object context
		// For bytecode execution, we'll use a custom bound method type
		boundMethod := &ObjectBoundMethod{
			Object: obj,
			Method: closure,
		}
		return vm.push(boundMethod)
	}
	
	// Try superclass methods
	if class.SuperClass != nil {
		if method, ok := class.SuperClass.CompiledMethods[propertyName]; ok {
			closure := &interpreter.Closure{Fn: method}
			boundMethod := &ObjectBoundMethod{
				Object: obj,
				Method: closure,
			}
			return vm.push(boundMethod)
		}
	}
	
	return fmt.Errorf("undefined method '%s' for class %s", propertyName, class.Name)
}

func (vm *VM) executeBuiltinFunctionProperty(builtin *interpreter.BuiltinFunction, propertyName string) error {
	// Call the builtin function to get the namespace object
	namespaceObj := builtin.Fn()
	
	// Handle different namespace types
	switch namespace := namespaceObj.(type) {
	case *interpreter.JSONNamespace:
		return vm.executeJSONNamespaceProperty(namespace, propertyName)
	case *interpreter.TimeNamespace:
		return vm.executeTimeNamespaceProperty(namespace, propertyName)
	default:
		return fmt.Errorf("property access not supported for namespace type: %T", namespaceObj)
	}
}

func (vm *VM) executeJSONNamespaceProperty(namespace *interpreter.JSONNamespace, propertyName string) error {
	switch propertyName {
	case "parse":
		// Return a builtin function that will call applyJSONNamespaceMethod
		parseFunction := &interpreter.BuiltinFunction{
			Fn: func(args ...interpreter.Value) interpreter.Value {
				return interpreter.ApplyJSONNamespaceMethod(namespace, "parse", args...)
			},
		}
		return vm.push(parseFunction)
	case "stringify":
		// Return a builtin function that will call applyJSONNamespaceMethod
		stringifyFunction := &interpreter.BuiltinFunction{
			Fn: func(args ...interpreter.Value) interpreter.Value {
				return interpreter.ApplyJSONNamespaceMethod(namespace, "stringify", args...)
			},
		}
		return vm.push(stringifyFunction)
	default:
		return fmt.Errorf("undefined method %s for JSON namespace", propertyName)
	}
}

func (vm *VM) executeTimeNamespaceProperty(namespace *interpreter.TimeNamespace, propertyName string) error {
	// Placeholder for Time namespace methods - can be implemented later
	return fmt.Errorf("Time namespace methods not implemented in bytecode yet")
}

func (vm *VM) executeJSONProperty(jsonObj *interpreter.JSON, propertyName string) error {
	switch propertyName {
	// Simple properties (no parameters)
	case "data":
		return vm.push(jsonObj.Data)
	case "type":
		return vm.push(&interpreter.String{Value: string(jsonObj.Data.Type())})
	case "valid":
		return vm.push(&interpreter.Boolean{Value: true}) // If we have a JSON object, it's valid
	
	// Methods (with parameters) - return bound methods
	case "get", "set", "has", "keys", "values", "length", "size",
		 "pretty", "compact", "path", "validate", "merge":
		return vm.push(&interpreter.JSONMethod{JSON: jsonObj, Method: propertyName})
	
	default:
		return fmt.Errorf("unknown property %s for JSON", propertyName)
	}
}

func (vm *VM) executeRegexpProperty(regexpObj *interpreter.Regexp, propertyName string) error {
	switch propertyName {
	// Simple properties (no parameters)
	case "pattern":
		return vm.push(&interpreter.String{Value: regexpObj.Pattern})
	
	// Methods (with parameters) - return bound methods
	case "matches?", "find_all", "find_first", "replace":
		return vm.push(&interpreter.RegexpMethod{Regexp: regexpObj, Method: propertyName})
	
	default:
		return fmt.Errorf("unknown property %s for Regexp", propertyName)
	}
}

// ObjectBoundMethod represents a method bound to an object for bytecode execution
type ObjectBoundMethod struct {
	Object *interpreter.Object
	Method *interpreter.Closure
}

func (obm *ObjectBoundMethod) Type() interpreter.ValueType { return "OBJECT_BOUND_METHOD" }
func (obm *ObjectBoundMethod) Inspect() string { return "bound method" }

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]

	switch callee := callee.(type) {
	case *interpreter.Closure:
		return vm.callClosure(callee, numArgs)
	case *interpreter.BuiltinFunction:
		return vm.callBuiltin(callee, numArgs)
	case *interpreter.StringMethod:
		return vm.callStringMethod(callee, numArgs)
	case *interpreter.ArrayMethod:
		return vm.callArrayMethod(callee, numArgs)
	case *interpreter.HashMethod:
		return vm.callHashMethod(callee, numArgs)
	case *interpreter.NumberMethod:
		return vm.callNumberMethod(callee, numArgs)
	case *interpreter.JSONMethod:
		return vm.callJSONMethod(callee, numArgs)
	case *interpreter.RegexpMethod:
		return vm.callRegexpMethod(callee, numArgs)
	case *interpreter.Class:
		return vm.callClassConstructor(callee, numArgs)
	case *ObjectBoundMethod:
		return vm.callObjectBoundMethod(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function and non-builtin: %T", callee)
	}
}

func (vm *VM) callClosure(cl *interpreter.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	// JIT compilation and execution if enabled
	if vm.jitEnabled && vm.jitCompiler != nil {
		startTime := time.Now()
		
		// Generate function hash for profiling and JIT compilation
		fnHash := vm.generateFunctionHash(cl.Fn)
		
		// Record function execution for profiling
		vm.RecordFunctionExecution(fnHash, 0) // Will be updated with actual time later
		
		// Check if function should be JIT compiled
		if vm.jitCompiler.ShouldCompile(fnHash) {
			vm.logger.Debug("Attempting JIT compilation for function %d", fnHash)
			
			err := vm.jitCompiler.Compile(cl.Fn, fnHash)
			if err != nil {
				vm.logger.Debug("JIT compilation failed: %v, falling back to bytecode", err)
				vm.stats.JITMisses++
			} else {
				vm.logger.Debug("JIT compilation succeeded for function %d", fnHash)
				vm.stats.JITCompilations++
			}
		}
		
		// Attempt JIT execution if compiled code is available (only if we have compiled code)
		// We need to check the cache directly rather than ShouldCompile again
		if vm.jitCompiler.ShouldCompile(fnHash) && vm.jitCompiler.GetStats().CompilationsSucceeded > 0 {
			// Try to execute JIT compiled code
			args := make([]interpreter.Value, numArgs)
			for i := 0; i < numArgs; i++ {
				args[i] = vm.stack[vm.sp-numArgs+i]
			}
			
			result, err := vm.jitCompiler.Execute(fnHash, args, vm.globals)
			if err != nil {
				// JIT execution failed, deoptimize to bytecode
				vm.logger.Debug("JIT execution failed: %v, deoptimizing to bytecode", err)
				vm.stats.JITDeoptimizations++
				// Fall through to bytecode execution
			} else {
				// JIT execution succeeded
				executionTime := time.Since(startTime)
				vm.RecordFunctionExecution(fnHash, executionTime)
				vm.stats.JITHits++
				
				// Remove function and arguments from stack
				vm.sp -= numArgs + 1
				
				// Push result
				if result != nil {
					return vm.push(result)
				}
				return vm.push(interpreter.NULL)
			}
		}
		
		// Update execution time for bytecode fallback
		executionTime := time.Since(startTime)
		vm.RecordFunctionExecution(fnHash, executionTime)
	}

	// Bytecode execution (original implementation)
	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	// Initialize all local variable slots to NULL
	for i := vm.sp; i < frame.basePointer + cl.Fn.NumLocals; i++ {
		vm.stack[i] = interpreter.NULL
	}

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callClosureWithSelf(cl *interpreter.Closure, numArgs int, self *interpreter.Object) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	frame := NewFrameWithSelf(cl, vm.sp-numArgs, self)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *interpreter.BuiltinFunction, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(args...)
	
	// For builtin calls, we need to remove the function and all arguments from the stack
	// Calculate the target SP after removing function + numArgs arguments
	targetSP := vm.sp - numArgs - 1
	
	// Ensure we don't go below the current frame's minimum SP
	frame := vm.currentFrame()
	minSP := frame.basePointer
	if frame.cl != nil {
		minSP += frame.cl.Fn.NumLocals
	}
	
	// If the target would go below minimum, something is wrong with the stack state
	// This should not happen with proper bytecode, but let's handle it gracefully
	if targetSP < minSP {
		vm.logger.Debug("Builtin call would underflow stack (target=%d, min=%d), using minimum", targetSP, minSP)
		targetSP = minSP
	}
	
	vm.sp = targetSP

	if result != nil {
		vm.push(result)
	} else {
		vm.push(interpreter.NULL)
	}

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.constants[constIndex]
	function, ok := constant.(*interpreter.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %T", constant)
	}

	free := make([]interpreter.Value, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.safeSetSP(vm.sp - numFree)

	closure := &interpreter.Closure{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) callStringMethod(method *interpreter.StringMethod, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.safeSetSP(vm.sp - numArgs - 1)

	var result interpreter.Value
	switch method.Method {
	case "upper":
		if numArgs != 0 {
			return fmt.Errorf("upper() takes no arguments, got %d", numArgs)
		}
		result = &interpreter.String{Value: strings.ToUpper(method.String.Value)}
	case "lower":
		if numArgs != 0 {
			return fmt.Errorf("lower() takes no arguments, got %d", numArgs)
		}
		result = &interpreter.String{Value: strings.ToLower(method.String.Value)}
	case "trim":
		if numArgs != 0 {
			return fmt.Errorf("trim() takes no arguments, got %d", numArgs)
		}
		result = &interpreter.String{Value: strings.TrimSpace(method.String.Value)}
	case "contains":
		if numArgs != 1 {
			return fmt.Errorf("contains() takes 1 argument, got %d", numArgs)
		}
		searchStr, ok := args[0].(*interpreter.String)
		if !ok {
			return fmt.Errorf("contains() argument must be string")
		}
		result = &interpreter.Boolean{Value: strings.Contains(method.String.Value, searchStr.Value)}
	case "match", "matches?", "replace", "split":
		// Delegate complex methods to interpreter
		argValues := make([]interpreter.Value, numArgs)
		for i := 0; i < numArgs; i++ {
			argValues[i] = args[i]
		}
		result = interpreter.ApplyStringMethod(method, argValues, nil)
	default:
		return fmt.Errorf("unknown string method: %s", method.Method)
	}

	return vm.push(result)
}

func (vm *VM) callArrayMethod(method *interpreter.ArrayMethod, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.safeSetSP(vm.sp - numArgs - 1)

	var result interpreter.Value
	switch method.Method {
	case "push":
		if numArgs == 0 {
			return fmt.Errorf("push() requires at least 1 argument")
		}
		// Add all arguments to the array
		for _, arg := range args {
			method.Array.Elements = append(method.Array.Elements, arg)
		}
		result = method.Array  // Return the modified array, not the length
	case "pop":
		if numArgs != 0 {
			return fmt.Errorf("pop() takes no arguments, got %d", numArgs)
		}
		if len(method.Array.Elements) == 0 {
			result = interpreter.NULL
		} else {
			lastIndex := len(method.Array.Elements) - 1
			result = method.Array.Elements[lastIndex]
			method.Array.Elements = method.Array.Elements[:lastIndex]
		}
	case "join":
		separator := ", "
		if numArgs == 1 {
			sep, ok := args[0].(*interpreter.String)
			if !ok {
				return fmt.Errorf("join() separator must be string")
			}
			separator = sep.Value
		} else if numArgs > 1 {
			return fmt.Errorf("join() takes 0 or 1 arguments, got %d", numArgs)
		}
		
		parts := make([]string, len(method.Array.Elements))
		for i, elem := range method.Array.Elements {
			parts[i] = elem.Inspect()
		}
		result = &interpreter.String{Value: strings.Join(parts, separator)}
	default:
		return fmt.Errorf("unknown array method: %s", method.Method)
	}

	return vm.push(result)
}

func (vm *VM) callHashMethod(method *interpreter.HashMethod, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.safeSetSP(vm.sp - numArgs - 1)

	var result interpreter.Value
	switch method.Method {
	case "has_key":
		if numArgs != 1 {
			return fmt.Errorf("has_key() takes 1 argument, got %d", numArgs)
		}
		
		// Check if key is hashable
		key := args[0]
		switch key.(type) {
		case *interpreter.Integer, *interpreter.String, *interpreter.Boolean, *interpreter.Float:
			hashKey := interpreter.CreateHashKey(key)
			_, exists := method.Hash.Pairs[hashKey]
			result = &interpreter.Boolean{Value: exists}
		default:
			result = &interpreter.Boolean{Value: false}
		}
	default:
		return fmt.Errorf("unknown hash method: %s", method.Method)
	}

	return vm.push(result)
}

func (vm *VM) callNumberMethod(method *interpreter.NumberMethod, numArgs int) error {
	vm.safeSetSP(vm.sp - numArgs - 1)

	var result interpreter.Value
	var numValue float64
	
	// Convert number to float64 for calculations
	switch num := method.Number.(type) {
	case *interpreter.Integer:
		numValue = float64(num.Value)
	case *interpreter.Float:
		numValue = num.Value
	}

	switch method.Method {
	case "abs":
		if numArgs != 0 {
			return fmt.Errorf("abs() takes no arguments, got %d", numArgs)
		}
		if numValue < 0 {
			numValue = -numValue
		}
		// Return same type as input
		if _, ok := method.Number.(*interpreter.Integer); ok {
			result = &interpreter.Integer{Value: int64(numValue)}
		} else {
			result = &interpreter.Float{Value: numValue}
		}
	default:
		return fmt.Errorf("unknown number method: %s", method.Method)
	}

	return vm.push(result)
}

func (vm *VM) callJSONMethod(method *interpreter.JSONMethod, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.safeSetSP(vm.sp - numArgs - 1)
	
	// Convert args to slice of interpreter.Value
	argValues := make([]interpreter.Value, numArgs)
	for i := 0; i < numArgs; i++ {
		argValues[i] = args[i]
	}
	
	// Use the existing applyJSONMethod function from interpreter
	result := interpreter.ApplyJSONMethod(method, argValues, nil)
	
	return vm.push(result)
}

func (vm *VM) callRegexpMethod(method *interpreter.RegexpMethod, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.safeSetSP(vm.sp - numArgs - 1)
	
	// Convert args to slice of interpreter.Value
	argValues := make([]interpreter.Value, numArgs)
	for i := 0; i < numArgs; i++ {
		argValues[i] = args[i]
	}
	
	// Use the existing ApplyRegexpMethod function from interpreter
	result := interpreter.ApplyRegexpMethod(method, argValues, nil)
	
	return vm.push(result)
}

func (vm *VM) callClassConstructor(class *interpreter.Class, numArgs int) error {
	// Create new instance
	instance := &interpreter.Object{
		Class:        class,
		InstanceVars: make(map[string]interpreter.Value),
	}
	
	// Pop arguments from stack
	args := make([]interpreter.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}
	
	// Pop class from stack
	vm.pop()
	
	// Call initialize method if it exists
	if initMethod, ok := class.CompiledMethods["initialize"]; ok {
		// Push arguments back on stack
		for _, arg := range args {
			err := vm.push(arg)
			if err != nil {
				return err
			}
		}
		
		// Create closure and call initialize method with instance context
		closure := &interpreter.Closure{Fn: initMethod}
		err := vm.callClosureWithSelf(closure, numArgs, instance)
		if err != nil {
			return err
		}
		
		// Pop the return value from initialize method
		vm.pop()
	}
	
	// Push new instance on stack
	return vm.push(instance)
}

func (vm *VM) callObjectBoundMethod(boundMethod *ObjectBoundMethod, numArgs int) error {
	// Pop arguments from stack
	args := make([]interpreter.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}
	
	// Pop the bound method from stack
	vm.pop()
	
	// Push arguments back on stack
	for _, arg := range args {
		err := vm.push(arg)
		if err != nil {
			return err
		}
	}
	
	// Call the method with the object context
	return vm.callClosureWithSelf(boundMethod.Method, numArgs, boundMethod.Object)
}

func nativeBoolToPushBool(input bool) interpreter.Value {
	if input {
		return interpreter.TRUE
	}
	return interpreter.FALSE
}

// Helper function for debugging
func getMethodNames(methods map[string]*interpreter.CompiledFunction) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	return names
}

// LastPoppedStackElem returns the last popped element (for testing)
func (vm *VM) LastPoppedStackElem() interpreter.Value {
	return vm.StackTop()
}

// executeMethodCall handles method invocation on objects
func (vm *VM) executeMethodCall(object interpreter.Value, methodName string, args []interpreter.Value) error {
	switch obj := object.(type) {
	case *interpreter.Object:
		// Instance method call
		class := obj.Class
		method, ok := class.CompiledMethods[methodName]
		if !ok {
			// Try superclass
			if class.SuperClass != nil {
				method, ok = class.SuperClass.CompiledMethods[methodName]
			}
			if !ok {
				return fmt.Errorf("undefined method '%s' for class %s", methodName, class.Name)
			}
		}
		
		// Create closure and call it
		closure := &interpreter.Closure{Fn: method}
		
		// Push arguments back on stack in reverse order
		for i := len(args) - 1; i >= 0; i-- {
			err := vm.push(args[i])
			if err != nil {
				return err
			}
		}
		
		// Push closure on stack
		err := vm.push(closure)
		if err != nil {
			return err
		}
		
		// Call the method with object context
		return vm.callClosureWithSelf(closure, len(args), obj)
		
	case *interpreter.Class:
		// Class method call (constructor)
		if methodName == "new" {
			// Create new instance
			instance := &interpreter.Object{
				Class:        obj,
				InstanceVars: make(map[string]interpreter.Value),
			}
			
			// Call initialize method if it exists
			if initMethod, ok := obj.CompiledMethods["initialize"]; ok {
				// Set up method call context with instance as 'self'
				closure := &interpreter.Closure{Fn: initMethod}
				
				// Push arguments
				for i := len(args) - 1; i >= 0; i-- {
					err := vm.push(args[i])
					if err != nil {
						return err
					}
				}
				
				// Push closure
				err := vm.push(closure)
				if err != nil {
					return err
				}
				
				// Call initialize method with instance context
				err = vm.callClosureWithSelf(closure, len(args), instance)
				if err != nil {
					return err
				}
				
				// Pop the return value and replace with instance
				vm.pop() // Remove return value
			}
			
			// Push new instance
			return vm.push(instance)
		}
		
		return fmt.Errorf("undefined class method '%s' for class %s", methodName, obj.Name)
		
	default:
		return fmt.Errorf("cannot call method on %T", object)
	}
}

// valueToString converts any value to its string representation for coercion
func valueToString(val interpreter.Value) string {
	switch val.Type() {
	case interpreter.STRING_VALUE:
		return val.(*interpreter.String).Value
	case interpreter.INTEGER_VALUE:
		return fmt.Sprintf("%d", val.(*interpreter.Integer).Value)
	case interpreter.FLOAT_VALUE:
		return fmt.Sprintf("%g", val.(*interpreter.Float).Value)
	case interpreter.BOOLEAN_VALUE:
		if val.(*interpreter.Boolean).Value {
			return "true"
		}
		return "false"
	case interpreter.NULL_VALUE:
		return "null"
	default:
		return val.Inspect() // fallback to existing representation
	}
}

func (vm *VM) executeBinaryMixedNumberOperation(op bytecode.Opcode, left, right interpreter.Value) error {
	var leftFloat, rightFloat float64
	
	if left.Type() == interpreter.INTEGER_VALUE {
		leftFloat = float64(left.(*interpreter.Integer).Value)
	} else {
		leftFloat = left.(*interpreter.Float).Value
	}
	
	if right.Type() == interpreter.INTEGER_VALUE {
		rightFloat = float64(right.(*interpreter.Integer).Value)
	} else {
		rightFloat = right.(*interpreter.Float).Value
	}
	
	var result float64
	switch op {
	case bytecode.OpAdd:
		result = leftFloat + rightFloat
	case bytecode.OpSub:
		result = leftFloat - rightFloat
	case bytecode.OpMul:
		result = leftFloat * rightFloat
	case bytecode.OpDiv:
		if rightFloat == 0 {
			return fmt.Errorf("division by zero")
		}
		result = leftFloat / rightFloat
	case bytecode.OpMod:
		if rightFloat == 0 {
			return fmt.Errorf("modulo by zero")
		}
		result = float64(int64(leftFloat) % int64(rightFloat))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
	
	return vm.push(&interpreter.Float{Value: result})
}

func (vm *VM) executeBinaryStringCoercionOperation(op bytecode.Opcode, left, right interpreter.Value) error {
	leftStr := valueToString(left)
	rightStr := valueToString(right)
	
	switch op {
	case bytecode.OpAdd:
		return vm.push(&interpreter.String{Value: leftStr + rightStr})
	default:
		// String coercion should only be used for addition/concatenation
		// If we get here, there's likely a parsing issue causing incorrect operator application
		leftTypeName := vm.getTypeName(left.Type())
		rightTypeName := vm.getTypeName(right.Type())
		opName := vm.getOperatorName(op)
		return fmt.Errorf("invalid string coercion: %s %s %s (opcode %d)", leftTypeName, opName, rightTypeName, op)
	}
}

// Helper function to convert ValueType to test-expected type name
func (vm *VM) getTypeName(valueType interpreter.ValueType) string {
	switch valueType {
	case interpreter.INTEGER_VALUE:
		return "INTEGER"
	case interpreter.BOOLEAN_VALUE:
		return "BOOLEAN"
	case interpreter.STRING_VALUE:
		return "STRING"
	case interpreter.FLOAT_VALUE:
		return "FLOAT"
	case interpreter.ARRAY_VALUE:
		return "ARRAY"
	case interpreter.HASH_VALUE:
		return "HASH"
	case interpreter.FUNCTION_VALUE:
		return "FUNCTION"
	case interpreter.CLOSURE_VALUE:
		return "FUNCTION"
	case interpreter.NULL_VALUE:
		return "NULL"
	case interpreter.BUILTIN_VALUE:
		return "BUILTIN"
	case interpreter.COMPILED_FUNCTION_VALUE:
		return "COMPILED_FUNCTION"
	case interpreter.RETURN_VALUE:
		return "RETURN_VALUE"
	case interpreter.BREAK_VALUE:
		return "BREAK_VALUE"
	case interpreter.CONTINUE_VALUE:
		return "CONTINUE_VALUE"
	case interpreter.EXCEPTION_VALUE:
		return "EXCEPTION"
	case interpreter.CLASS_VALUE:
		return "CLASS"
	case interpreter.INSTANCE_VALUE:
		return "INSTANCE"
	case interpreter.BOUND_METHOD_VALUE:
		return "BOUND_METHOD"
	case interpreter.HASH_METHOD_VALUE:
		return "HASH_METHOD"
	case interpreter.STRING_METHOD_VALUE:
		return "STRING_METHOD"
	case interpreter.ARRAY_METHOD_VALUE:
		return "ARRAY_METHOD"
	case interpreter.NUMBER_METHOD_VALUE:
		return "NUMBER_METHOD"
	case interpreter.FILE_VALUE:
		return "FILE"
	case interpreter.DIRECTORY_VALUE:
		return "DIRECTORY"
	case interpreter.PATH_VALUE:
		return "PATH"
	case interpreter.FILE_METHOD_VALUE:
		return "FILE_METHOD"
	case interpreter.DIRECTORY_METHOD_VALUE:
		return "DIRECTORY_METHOD"
	case interpreter.PATH_METHOD_VALUE:
		return "PATH_METHOD"
	case interpreter.JSON_VALUE:
		return "JSON"
	case interpreter.JSON_METHOD_VALUE:
		return "JSON_METHOD"
	case interpreter.JSON_NAMESPACE_VALUE:
		return "JSON_NAMESPACE"
	case interpreter.TIME_VALUE:
		return "TIME"
	case interpreter.TIME_METHOD_VALUE:
		return "TIME_METHOD"
	case interpreter.TIME_NAMESPACE_VALUE:
		return "TIME_NAMESPACE"
	case interpreter.DURATION_VALUE:
		return "DURATION"
	case interpreter.DURATION_METHOD_VALUE:
		return "DURATION_METHOD"
	case interpreter.DURATION_NAMESPACE_VALUE:
		return "DURATION_NAMESPACE"
	case interpreter.TIMEZONE_VALUE:
		return "TIMEZONE"
	case interpreter.TIMEZONE_METHOD_VALUE:
		return "TIMEZONE_METHOD"
	case interpreter.TIMEZONE_NAMESPACE_VALUE:
		return "TIMEZONE_NAMESPACE"
	default:
		return "UNKNOWN"
	}
}

// Helper function to convert bytecode opcode to operator symbol
func (vm *VM) getOperatorName(op bytecode.Opcode) string {
	switch op {
	case bytecode.OpAdd:
		return "+"
	case bytecode.OpSub:
		return "-"
	case bytecode.OpMul:
		return "*"
	case bytecode.OpDiv:
		return "/"
	case bytecode.OpMod:
		return "%"
	case bytecode.OpEqual:
		return "=="
	case bytecode.OpNotEqual:
		return "!="
	case bytecode.OpGreaterThan:
		return ">"
	case bytecode.OpLessThan:
		return "<"
	default:
		return "UNKNOWN"
	}
}

// Helper function to convert bytecode opcode to instruction name for logging
func (vm *VM) getOpcodeName(op bytecode.Opcode) string {
	switch op {
	case bytecode.OpConstant:
		return "OpConstant"
	case bytecode.OpPop:
		return "OpPop"
	case bytecode.OpAdd:
		return "OpAdd"
	case bytecode.OpSub:
		return "OpSub"
	case bytecode.OpMul:
		return "OpMul"
	case bytecode.OpDiv:
		return "OpDiv"
	case bytecode.OpMod:
		return "OpMod"
	case bytecode.OpTrue:
		return "OpTrue"
	case bytecode.OpFalse:
		return "OpFalse"
	case bytecode.OpNull:
		return "OpNull"
	case bytecode.OpEqual:
		return "OpEqual"
	case bytecode.OpNotEqual:
		return "OpNotEqual"
	case bytecode.OpGreaterThan:
		return "OpGreaterThan"
	case bytecode.OpGreaterEqual:
		return "OpGreaterEqual"
	case bytecode.OpAnd:
		return "OpAnd"
	case bytecode.OpOr:
		return "OpOr"
	case bytecode.OpNot:
		return "OpNot"
	case bytecode.OpMinus:
		return "OpMinus"
	case bytecode.OpJump:
		return "OpJump"
	case bytecode.OpJumpNotTruthy:
		return "OpJumpNotTruthy"
	case bytecode.OpJumpTruthy:
		return "OpJumpTruthy"
	case bytecode.OpSetGlobal:
		return "OpSetGlobal"
	case bytecode.OpGetGlobal:
		return "OpGetGlobal"
	case bytecode.OpSetLocal:
		return "OpSetLocal"
	case bytecode.OpGetLocal:
		return "OpGetLocal"
	case bytecode.OpArray:
		return "OpArray"
	case bytecode.OpHash:
		return "OpHash"
	case bytecode.OpIndex:
		return "OpIndex"
	case bytecode.OpSetIndex:
		return "OpSetIndex"
	case bytecode.OpCall:
		return "OpCall"
	case bytecode.OpReturn:
		return "OpReturn"
	case bytecode.OpReturnVoid:
		return "OpReturnVoid"
	case bytecode.OpGetBuiltin:
		return "OpGetBuiltin"
	case bytecode.OpClosure:
		return "OpClosure"
	case bytecode.OpGetFree:
		return "OpGetFree"
	case bytecode.OpSetFree:
		return "OpSetFree"
	case bytecode.OpCurrentClosure:
		return "OpCurrentClosure"
	case bytecode.OpThrow:
		return "OpThrow"
	case bytecode.OpTryBegin:
		return "OpTryBegin"
	case bytecode.OpTryEnd:
		return "OpTryEnd"
	case bytecode.OpCatch:
		return "OpCatch"
	case bytecode.OpFinally:
		return "OpFinally"
	case bytecode.OpImport:
		return "OpImport"
	case bytecode.OpExport:
		return "OpExport"
	case bytecode.OpGetProperty:
		return "OpGetProperty"
	case bytecode.OpClass:
		return "OpClass"
	case bytecode.OpInherit:
		return "OpInherit"
	case bytecode.OpMethod:
		return "OpMethod"
	case bytecode.OpInvoke:
		return "OpInvoke"
	case bytecode.OpGetInstance:
		return "OpGetInstance"
	case bytecode.OpSetInstance:
		return "OpSetInstance"
	case bytecode.OpGetSuper:
		return "OpGetSuper"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", op)
	}
}