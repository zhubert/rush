//go:build llvm14 || llvm15 || llvm16 || llvm17 || llvm18 || llvm19 || llvm20

package aot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rush/compiler"
	"rush/bytecode"
	"rush/interpreter"
	
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

// GenerateExecutable creates an executable by translating bytecode to LLVM IR
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
	
	// Declare runtime functions
	if err := g.declareRuntimeFunctions(module, context); err != nil {
		return nil, fmt.Errorf("failed to declare runtime functions: %w", err)
	}
	
	// Create main function type: int main()
	mainType := llvm.FunctionType(context.Int32Type(), nil, false)
	mainFunc := llvm.AddFunction(module, "main", mainType)
	
	// Create entry basic block
	entry := context.AddBasicBlock(mainFunc, "entry")
	builder.SetInsertPointAtEnd(entry)
	
	// Translate bytecode to LLVM IR
	if err := g.translateBytecode(bytecode, module, context, builder, mainFunc); err != nil {
		return nil, fmt.Errorf("bytecode translation failed: %w", err)
	}
	
	// Return 0 from main
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

// declareRuntimeFunctions declares Rush runtime functions in the LLVM module
func (g *SimpleLLVMCodeGenerator) declareRuntimeFunctions(module llvm.Module, context llvm.Context) error {
	// void rush_print(char* str, size_t len)
	charPtrType := llvm.PointerType(context.Int8Type(), 0)
	sizeTType := context.Int64Type()
	printType := llvm.FunctionType(context.VoidType(), []llvm.Type{charPtrType, sizeTType}, false)
	llvm.AddFunction(module, "rush_print", printType)
	
	// void rush_print_line(char* str, size_t len) - print with newline
	printlnType := llvm.FunctionType(context.VoidType(), []llvm.Type{charPtrType, sizeTType}, false)
	llvm.AddFunction(module, "rush_print_line", printlnType)
	
	return nil
}

// translateBytecode converts Rush bytecode to LLVM IR
func (g *SimpleLLVMCodeGenerator) translateBytecode(bc *compiler.Bytecode, module llvm.Module, context llvm.Context, builder llvm.Builder, mainFunc llvm.Value) error {
	// Simple stack simulation using alloca
	stackPtr := builder.CreateAlloca(llvm.ArrayType(context.Int64Type(), 1000), "stack")
	stackTop := builder.CreateAlloca(context.Int32Type(), "stack_top")
	builder.CreateStore(llvm.ConstInt(context.Int32Type(), 0, false), stackTop)
	
	// Keep track of constants on the stack
	constantStack := make([]interface{}, 0)
	
	// Process each bytecode instruction
	for i := 0; i < len(bc.Instructions); {
		if i >= len(bc.Instructions) {
			break
		}
		opcode := bytecode.Opcode(bc.Instructions[i])
		
		switch opcode {
		case bytecode.OpConstant:
			// Extract constant index from instruction
			if i+2 >= len(bc.Instructions) {
				return fmt.Errorf("OpConstant instruction incomplete at position %d", i)
			}
			constIndex := int(bytecode.ReadUint16(bc.Instructions[i+1:]))
			
			// Bounds check for constants array
			if constIndex >= len(bc.Constants) {
				return fmt.Errorf("constant index %d out of range (constants length: %d)", constIndex, len(bc.Constants))
			}
			constant := bc.Constants[constIndex]
			
			// Store the constant on our simulated stack
			constantStack = append(constantStack, constant)
			
			// Convert Rush value to LLVM value and push to stack
			llvmValue, err := g.convertRushValueToLLVM(constant, context, builder)
			if err != nil {
				return fmt.Errorf("failed to convert constant: %w", err)
			}
			
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push constant: %w", err)
			}
			
			i += 3 // OpConstant takes 3 bytes
			
		case bytecode.OpTrue:
			// Push boolean true value onto constant stack
			trueValue := &interpreter.Boolean{Value: true}
			constantStack = append(constantStack, trueValue)
			
			// Convert to LLVM value and push to stack
			llvmValue := llvm.ConstInt(context.Int1Type(), 1, false)
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push true: %w", err)
			}
			
			i += 1 // OpTrue takes 1 byte
			
		case bytecode.OpFalse:
			// Push boolean false value onto constant stack
			falseValue := &interpreter.Boolean{Value: false}
			constantStack = append(constantStack, falseValue)
			
			// Convert to LLVM value and push to stack
			llvmValue := llvm.ConstInt(context.Int1Type(), 0, false)
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push false: %w", err)
			}
			
			i += 1 // OpFalse takes 1 byte
			
		case bytecode.OpNull:
			// Push null value onto constant stack
			nullValue := &interpreter.Null{}
			constantStack = append(constantStack, nullValue)
			
			// For LLVM, represent null as 0
			llvmValue := llvm.ConstInt(context.Int64Type(), 0, false)
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push null: %w", err)
			}
			
			i += 1 // OpNull takes 1 byte
			
		case bytecode.OpArray:
			// Extract element count from instruction
			if i+2 >= len(bc.Instructions) {
				return fmt.Errorf("OpArray instruction incomplete at position %d", i)
			}
			elementCount := int(bytecode.ReadUint16(bc.Instructions[i+1:]))
			
			// Pop elements from constant stack and create array
			if len(constantStack) < elementCount {
				return fmt.Errorf("not enough elements on stack for array of size %d", elementCount)
			}
			
			// Get the elements (they're at the end of the stack)
			elements := constantStack[len(constantStack)-elementCount:]
			constantStack = constantStack[:len(constantStack)-elementCount]
			
			// Create array object
			arrayElements := make([]interpreter.Value, elementCount)
			for i, elem := range elements {
				if val, ok := elem.(interpreter.Value); ok {
					arrayElements[i] = val
				} else {
					return fmt.Errorf("invalid array element type: %T", elem)
				}
			}
			arrayValue := &interpreter.Array{Elements: arrayElements}
			constantStack = append(constantStack, arrayValue)
			
			// For LLVM, just push a placeholder (arrays need complex runtime support)
			llvmValue := llvm.ConstInt(context.Int64Type(), uint64(elementCount), false)
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push array: %w", err)
			}
			
			i += 3 // OpArray takes 3 bytes (1 opcode + 2 count)
			
		case bytecode.OpHash:
			// Extract pair count from instruction
			if i+2 >= len(bc.Instructions) {
				return fmt.Errorf("OpHash instruction incomplete at position %d", i)
			}
			pairCount := int(bytecode.ReadUint16(bc.Instructions[i+1:]))
			elementCount := pairCount * 2 // keys + values
			
			// Pop elements from constant stack and create hash
			if len(constantStack) < elementCount {
				return fmt.Errorf("not enough elements on stack for hash with %d pairs", pairCount)
			}
			
			// Get the elements (they're at the end of the stack)
			elements := constantStack[len(constantStack)-elementCount:]
			constantStack = constantStack[:len(constantStack)-elementCount]
			
			// Create hash object (pairs are key, value, key, value, ...)
			pairs := make(map[interpreter.HashKey]interpreter.Value)
			keys := make([]interpreter.Value, 0, pairCount)
			for i := 0; i < len(elements); i += 2 {
				keyElem := elements[i]
				valueElem := elements[i+1]
				
				// Convert to Value types
				key, ok1 := keyElem.(interpreter.Value)
				value, ok2 := valueElem.(interpreter.Value)
				if !ok1 || !ok2 {
					return fmt.Errorf("invalid hash element types: key %T, value %T", keyElem, valueElem)
				}
				
				// Create hash key and store
				hashKey := interpreter.CreateHashKey(key)
				pairs[hashKey] = value
				keys = append(keys, key)
			}
			hashValue := &interpreter.Hash{Pairs: pairs, Keys: keys}
			constantStack = append(constantStack, hashValue)
			
			// For LLVM, just push a placeholder (hashes need complex runtime support)
			llvmValue := llvm.ConstInt(context.Int64Type(), uint64(pairCount), false)
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push hash: %w", err)
			}
			
			i += 3 // OpHash takes 3 bytes (1 opcode + 2 count)
			
		case bytecode.OpCall:
			// Extract argument count
			if i+1 >= len(bc.Instructions) {
				return fmt.Errorf("OpCall instruction incomplete at position %d", i)
			}
			argCount := int(bc.Instructions[i+1])
			
			// For built-in functions, check if it's a print call
			// Get the argument from the constant stack
			var currentConstant interface{}
			if len(constantStack) > 0 {
				// Pop the last constant from the stack
				currentConstant = constantStack[len(constantStack)-1]
				constantStack = constantStack[:len(constantStack)-1]
			}
			
			if err := g.handleBuiltinCall(builder, context, module, stackPtr, stackTop, argCount, currentConstant); err != nil {
				return fmt.Errorf("failed to handle builtin call: %w", err)
			}
			
			i += 2 // OpCall takes 2 bytes
			
		default:
			// For unhandled opcodes, use the definition to get the correct length
			def, err := bytecode.Lookup(opcode)
			if err != nil {
				return fmt.Errorf("unknown opcode %d at position %d: %w", int(opcode), i, err)
			}
			
			// Calculate instruction length: 1 (opcode) + sum of operand widths
			instrLength := 1
			for _, width := range def.OperandWidths {
				instrLength += width
			}
			
			// Silently skip unhandled opcodes
			i += instrLength
		}
	}
	
	return nil
}

// convertRushValueToLLVM converts a Rush object to LLVM value
func (g *SimpleLLVMCodeGenerator) convertRushValueToLLVM(obj interface{}, context llvm.Context, builder llvm.Builder) (llvm.Value, error) {
	// Handle different Rush value types
	switch v := obj.(type) {
	case *interpreter.String:
		// Rush String object
		str := v.Value
		strConst := context.ConstString(str, false)
		return strConst, nil
		
	case *interpreter.Integer:
		// Rush Integer object
		intVal := llvm.ConstInt(context.Int64Type(), uint64(v.Value), true)
		return intVal, nil
		
	case *interpreter.Float:
		// Rush Float object
		floatVal := llvm.ConstFloat(context.DoubleType(), v.Value)
		return floatVal, nil
		
	case *interpreter.Boolean:
		// Rush Boolean object
		var boolVal llvm.Value
		if v.Value {
			boolVal = llvm.ConstInt(context.Int1Type(), 1, false)
		} else {
			boolVal = llvm.ConstInt(context.Int1Type(), 0, false)
		}
		return boolVal, nil
		
	case string:
		// Direct string value
		str := v
		if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
			str = str[1 : len(str)-1]
		}
		strConst := context.ConstString(str, false)
		return strConst, nil
		
	default:
		// Try to handle objects with Inspect() method (proper Rush formatting)
		if inspectObj, ok := obj.(interface{ Inspect() string }); ok {
			str := inspectObj.Inspect()
			
			// Create global string constant with proper formatting
			strConst := context.ConstString(str, false)
			return strConst, nil
		}
		
		// Fallback: try to handle objects with String() method
		if strObj, ok := obj.(interface{ String() string }); ok {
			str := strObj.String()
			// Remove quotes if present
			if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
				str = str[1 : len(str)-1]
			}
			
			// Create global string constant
			strConst := context.ConstString(str, false)
			return strConst, nil
		}
	}
	
	return llvm.Value{}, fmt.Errorf("unsupported value type: %T", obj)
}

// pushToStack simulates pushing a value to the runtime stack
func (g *SimpleLLVMCodeGenerator) pushToStack(builder llvm.Builder, context llvm.Context, stackPtr, stackTop, value llvm.Value) error {
	// Load current stack top
	currentTop := builder.CreateLoad(context.Int32Type(), stackTop, "current_top")
	
	// For now, just increment stack top
	newTop := builder.CreateAdd(currentTop, llvm.ConstInt(context.Int32Type(), 1, false), "new_top")
	builder.CreateStore(newTop, stackTop)
	
	return nil
}

// handleBuiltinCall handles calls to built-in functions like print
func (g *SimpleLLVMCodeGenerator) handleBuiltinCall(builder llvm.Builder, context llvm.Context, module llvm.Module, stackPtr, stackTop llvm.Value, argCount int, currentConstant interface{}) error {
	if argCount == 1 {
		// Assume this is a print call with one string argument
		printFunc := module.NamedFunction("rush_print_line")
		if printFunc.IsNil() {
			return fmt.Errorf("rush_print_line function not found")
		}
		
		// Get the print function type
		charPtrType := llvm.PointerType(context.Int8Type(), 0)
		sizeTType := context.Int64Type()
		printType := llvm.FunctionType(context.VoidType(), []llvm.Type{charPtrType, sizeTType}, false)
		
		// Pop the argument from the stack (simulated)
		// Load current stack top
		currentTop := builder.CreateLoad(context.Int32Type(), stackTop, "current_top")
		
		// Decrement stack top to pop the value
		newTop := builder.CreateSub(currentTop, llvm.ConstInt(context.Int32Type(), 1, false), "new_top")
		builder.CreateStore(newTop, stackTop)
		
		// Extract the actual string value from the current constant
		var actualStr string
		if currentConstant != nil {
			switch v := currentConstant.(type) {
			case *interpreter.String:
				actualStr = v.Value
			case string:
				// Remove quotes if present
				if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
					actualStr = v[1 : len(v)-1]
				} else {
					actualStr = v
				}
			default:
				// Try to get proper Rush string representation using Inspect()
				if inspectObj, ok := currentConstant.(interface{ Inspect() string }); ok {
					actualStr = inspectObj.Inspect()
				} else if strObj, ok := currentConstant.(interface{ String() string }); ok {
					str := strObj.String()
					// Remove quotes if present
					if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
						actualStr = str[1 : len(str)-1]
					} else {
						actualStr = str
					}
				} else {
					actualStr = fmt.Sprintf("%v", currentConstant)
				}
			}
		} else {
			actualStr = "null"
		}
		
		// Create LLVM constant string from the actual value
		strConst := context.ConstString(actualStr, false)
		
		// Create global for the string
		global := llvm.AddGlobal(module, strConst.Type(), "print_str")
		global.SetInitializer(strConst)
		global.SetLinkage(llvm.InternalLinkage)
		
		// Get pointer to string data
		strPtr := builder.CreateGEP(strConst.Type(), global, []llvm.Value{
			llvm.ConstInt(context.Int32Type(), 0, false),
			llvm.ConstInt(context.Int32Type(), 0, false),
		}, "str_ptr")
		
		// Call rush_print_line with string pointer and length
		strLen := llvm.ConstInt(context.Int64Type(), uint64(len(actualStr)), false)
		builder.CreateCall(printType, printFunc, []llvm.Value{strPtr, strLen}, "")
	}
	
	return nil
}

// LinkExecutable links the object file to create an executable
func (g *SimpleLLVMCodeGenerator) LinkExecutable(objectFile, outputPath string) error {
	// Find runtime.c file - assuming we're in the rush project directory
	runtimePath := "runtime/runtime.c"
	
	// Use clang to compile and link runtime.c with the object file
	cmd := exec.Command("clang", "-o", outputPath, objectFile, runtimePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("linking failed: %s\nOutput: %s", err, string(output))
	}
	return nil
}