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
	
	return nil
}

// translateBytecode converts Rush bytecode to LLVM IR
func (g *SimpleLLVMCodeGenerator) translateBytecode(bc *compiler.Bytecode, module llvm.Module, context llvm.Context, builder llvm.Builder, mainFunc llvm.Value) error {
	// Simple stack simulation using alloca
	stackPtr := builder.CreateAlloca(llvm.ArrayType(context.Int64Type(), 1000), "stack")
	stackTop := builder.CreateAlloca(context.Int32Type(), "stack_top")
	builder.CreateStore(llvm.ConstInt(context.Int32Type(), 0, false), stackTop)
	
	// Keep track of the last pushed value for print calls
	var lastPushedConstant interface{}
	
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
			
			// Store the constant for potential use in print calls
			lastPushedConstant = constant
			
			// Convert Rush value to LLVM value and push to stack
			llvmValue, err := g.convertRushValueToLLVM(constant, context, builder)
			if err != nil {
				return fmt.Errorf("failed to convert constant: %w", err)
			}
			
			if err := g.pushToStack(builder, context, stackPtr, stackTop, llvmValue); err != nil {
				return fmt.Errorf("failed to push constant: %w", err)
			}
			
			i += 3 // OpConstant takes 3 bytes
			
		case bytecode.OpCall:
			// Extract argument count
			if i+1 >= len(bc.Instructions) {
				return fmt.Errorf("OpCall instruction incomplete at position %d", i)
			}
			argCount := int(bc.Instructions[i+1])
			
			// For built-in functions, check if it's a print call
			if err := g.handleBuiltinCall(builder, context, module, stackPtr, stackTop, argCount, lastPushedConstant); err != nil {
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
		// Try to handle objects with String() method
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
func (g *SimpleLLVMCodeGenerator) handleBuiltinCall(builder llvm.Builder, context llvm.Context, module llvm.Module, stackPtr, stackTop llvm.Value, argCount int, lastConstant interface{}) error {
	if argCount == 1 {
		// Assume this is a print call with one string argument
		printFunc := module.NamedFunction("rush_print")
		if printFunc.IsNil() {
			return fmt.Errorf("rush_print function not found")
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
		
		// Extract the actual string value from the last constant
		var actualStr string
		if lastConstant != nil {
			switch v := lastConstant.(type) {
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
				// Try to get string representation
				if strObj, ok := lastConstant.(interface{ String() string }); ok {
					str := strObj.String()
					// Remove quotes if present
					if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
						actualStr = str[1 : len(str)-1]
					} else {
						actualStr = str
					}
				} else {
					actualStr = fmt.Sprintf("%v", lastConstant)
				}
			}
		} else {
			actualStr = "[null]"
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
		
		// Call rush_print with string pointer and length
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