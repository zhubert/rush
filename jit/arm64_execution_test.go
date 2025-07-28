package jit

import (
	"runtime"
	"testing"
	"unsafe"

	"rush/interpreter"
)

func TestARM64ExecutionEnvironment(t *testing.T) {
	// Skip if not on ARM64
	if runtime.GOARCH != "arm64" {
		t.Skip("ARM64 execution tests require ARM64 architecture")
	}

	err := ValidateExecutionEnvironment()
	if err != nil {
		t.Fatalf("Execution environment validation failed: %v", err)
	}
}

func TestValueMarshaling(t *testing.T) {
	code := &CompiledCode{}

	tests := []struct {
		name     string
		value    interpreter.Value
		expected uint64
	}{
		{
			name:     "Integer positive",
			value:    &interpreter.Integer{Value: 42},
			expected: 42,
		},
		{
			name:     "Integer negative",
			value:    &interpreter.Integer{Value: -10},
			expected: uint64(^uint64(10) + 1), // Two's complement representation
		},
		{
			name:     "Boolean true",
			value:    &interpreter.Boolean{Value: true},
			expected: 1,
		},
		{
			name:     "Boolean false",
			value:    &interpreter.Boolean{Value: false},
			expected: 0,
		},
		{
			name:     "Null",
			value:    &interpreter.Null{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := code.valueToUint64(tt.value)
			if err != nil {
				t.Fatalf("valueToUint64() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("valueToUint64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValueUnmarshaling(t *testing.T) {
	code := &CompiledCode{}

	tests := []struct {
		name          string
		value         uint64
		typeIndicator uint8
		expectedType  interpreter.ValueType
		expectedValue interface{}
	}{
		{
			name:          "Integer",
			value:         42,
			typeIndicator: 0,
			expectedType:  interpreter.INTEGER_VALUE,
			expectedValue: int64(42),
		},
		{
			name:          "Boolean true",
			value:         1,
			typeIndicator: 2,
			expectedType:  interpreter.BOOLEAN_VALUE,
			expectedValue: true,
		},
		{
			name:          "Boolean false",
			value:         0,
			typeIndicator: 2,
			expectedType:  interpreter.BOOLEAN_VALUE,
			expectedValue: false,
		},
		{
			name:          "Null",
			value:         0,
			typeIndicator: 4,
			expectedType:  interpreter.NULL_VALUE,
			expectedValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := code.uint64ToValue(tt.value, tt.typeIndicator)
			if err != nil {
				t.Fatalf("uint64ToValue() error = %v", err)
			}

			// Check type
			if result.Type() != tt.expectedType {
				t.Errorf("uint64ToValue() type = %v, want %v", result.Type(), tt.expectedType)
			}

			// Check value based on type
			switch v := result.(type) {
			case *interpreter.Integer:
				if v.Value != tt.expectedValue.(int64) {
					t.Errorf("Integer = %v, want %v", v.Value, tt.expectedValue)
				}
			case *interpreter.Boolean:
				if v.Value != tt.expectedValue.(bool) {
					t.Errorf("Boolean = %v, want %v", v.Value, tt.expectedValue)
				}
			case *interpreter.Null:
				// Null values don't have a meaningful value to check
			}
		})
	}
}

func TestARM64ArgumentMarshaling(t *testing.T) {
	code := &CompiledCode{}

	args := []interpreter.Value{
		&interpreter.Integer{Value: 10},
		&interpreter.Integer{Value: 20},
		&interpreter.Boolean{Value: true},
		&interpreter.Null{},
	}

	arm64Args, err := code.marshalArgsToARM64(args)
	if err != nil {
		t.Fatalf("marshalArgsToARM64() error = %v", err)
	}

	expected := [8]uint64{10, 20, 1, 0, 0, 0, 0, 0}
	if arm64Args != expected {
		t.Errorf("marshalArgsToARM64() = %v, want %v", arm64Args, expected)
	}
}

func TestARM64ArgumentMarshalingTooManyArgs(t *testing.T) {
	code := &CompiledCode{}

	// Create 9 arguments (more than ARM64 can handle)
	args := make([]interpreter.Value, 9)
	for i := 0; i < 9; i++ {
		args[i] = &interpreter.Integer{Value: int64(i)}
	}

	_, err := code.marshalArgsToARM64(args)
	if err == nil {
		t.Error("marshalArgsToARM64() should fail with too many arguments")
	}
}

func TestARM64ResultUnmarshaling(t *testing.T) {
	code := &CompiledCode{}

	tests := []struct {
		name     string
		result   *ARM64Result
		wantErr  bool
		expected interface{}
	}{
		{
			name: "Success integer",
			result: &ARM64Result{
				Value: 42,
				Type:  0, // INTEGER
				Error: 0,
			},
			wantErr:  false,
			expected: int64(42),
		},
		{
			name: "Success boolean",
			result: &ARM64Result{
				Value: 1,
				Type:  2, // BOOLEAN
				Error: 0,
			},
			wantErr:  false,
			expected: true,
		},
		{
			name: "Error result",
			result: &ARM64Result{
				Value: 0,
				Type:  0,
				Error: 1,
			},
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := code.unmarshalResultFromARM64(tt.result)
			
			if tt.wantErr {
				if err == nil {
					t.Error("unmarshalResultFromARM64() should return error")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("unmarshalResultFromARM64() error = %v", err)
			}

			switch v := result.(type) {
			case *interpreter.Integer:
				if v.Value != tt.expected.(int64) {
					t.Errorf("Integer = %v, want %v", v.Value, tt.expected)
				}
			case *interpreter.Boolean:
				if v.Value != tt.expected.(bool) {
					t.Errorf("Boolean = %v, want %v", v.Value, tt.expected)
				}
			}
		})
	}
}

func TestARM64CodeValidation(t *testing.T) {
	tests := []struct {
		name    string
		code    []byte
		wantErr bool
	}{
		{
			name:    "Empty code",
			code:    []byte{},
			wantErr: true,
		},
		{
			name:    "Unaligned code",
			code:    []byte{0x01, 0x02, 0x03},
			wantErr: true,
		},
		{
			name:    "Valid aligned code",
			code:    []byte{0x00, 0x00, 0x80, 0xD2}, // MOV X0, #0
			wantErr: false,
		},
		{
			name:    "Invalid instruction (all zeros)",
			code:    []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: true,
		},
		{
			name:    "Invalid instruction (all ones)",
			code:    []byte{0xFF, 0xFF, 0xFF, 0xFF},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateARM64Code(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateARM64Code() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestARM64InstructionValidation(t *testing.T) {
	tests := []struct {
		name        string
		instruction uint32
		expected    bool
	}{
		{
			name:        "All zeros",
			instruction: 0x00000000,
			expected:    false,
		},
		{
			name:        "All ones",
			instruction: 0xFFFFFFFF,
			expected:    false,
		},
		{
			name:        "MOV immediate",
			instruction: 0xD2800000, // MOV X0, #0
			expected:    true,
		},
		{
			name:        "ADD immediate",
			instruction: 0x91000000, // ADD X0, X0, #0
			expected:    true,
		},
		{
			name:        "RET",
			instruction: 0xD65F03C0,
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidARM64Instruction(tt.instruction)
			if result != tt.expected {
				t.Errorf("isValidARM64Instruction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExceptionHandler(t *testing.T) {
	handler := NewARM64ExceptionHandler()

	// Test installation and uninstallation
	err := handler.Install()
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if !handler.isActive {
		t.Error("Handler should be active after installation")
	}

	err = handler.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	if handler.isActive {
		t.Error("Handler should not be active after uninstallation")
	}
}

func TestExceptionClassification(t *testing.T) {
	handler := NewARM64ExceptionHandler()

	tests := []struct {
		name        string
		instruction uint32
		expected    ExceptionType
	}{
		{
			name:        "Illegal instruction (all zeros)",
			instruction: 0x00000000,
			expected:    ExceptionIllegalInstruction,
		},
		{
			name:        "Division instruction",
			instruction: 0x9AC00C00, // SDIV
			expected:    ExceptionArithmeticError,
		},
		{
			name:        "Load instruction",
			instruction: 0xF9400000, // LDR
			expected:    ExceptionSegmentationFault,
		},
		{
			name:        "Normal instruction",
			instruction: 0xD2800000, // MOV
			expected:    ExceptionUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.classifyException(tt.instruction)
			if result != tt.expected {
				t.Errorf("classifyException() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExceptionRecoverability(t *testing.T) {
	handler := NewARM64ExceptionHandler()

	tests := []struct {
		name      string
		excType   ExceptionType
		expected  bool
	}{
		{
			name:     "Arithmetic exception",
			excType:  ExceptionArithmeticError,
			expected: true,
		},
		{
			name:     "Floating point exception",
			excType:  ExceptionFloatingPointError,
			expected: true,
		},
		{
			name:     "Illegal instruction",
			excType:  ExceptionIllegalInstruction,
			expected: false,
		},
		{
			name:     "Segmentation fault",
			excType:  ExceptionSegmentationFault,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.isRecoverable(tt.excType)
			if result != tt.expected {
				t.Errorf("isRecoverable() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestARM64CallingConvention(t *testing.T) {
	conv := GetARM64CallingConvention()

	if len(conv.ArgRegs) != 8 {
		t.Errorf("Expected 8 argument registers, got %d", len(conv.ArgRegs))
	}

	if conv.ReturnReg != "X0" {
		t.Errorf("Expected return register X0, got %s", conv.ReturnReg)
	}

	if conv.FramePointer != "X29" {
		t.Errorf("Expected frame pointer X29, got %s", conv.FramePointer)
	}

	if conv.LinkRegister != "X30" {
		t.Errorf("Expected link register X30, got %s", conv.LinkRegister)
	}
}

func TestExecutionContextCreation(t *testing.T) {
	args := []interpreter.Value{
		&interpreter.Integer{Value: 10},
		&interpreter.Integer{Value: 20},
	}
	globals := []interpreter.Value{
		&interpreter.String{Value: "test"},
	}

	ctx := &ExecutionContext{
		Args:    args,
		Globals: globals,
		Stack:   make([]interpreter.Value, 1024),
	}

	if len(ctx.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(ctx.Args))
	}

	if len(ctx.Globals) != 1 {
		t.Errorf("Expected 1 global, got %d", len(ctx.Globals))
	}

	if len(ctx.Stack) != 1024 {
		t.Errorf("Expected stack size 1024, got %d", len(ctx.Stack))
	}
}

func TestARM64ExecutionContextCreation(t *testing.T) {
	args := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	
	ctx := &ARM64ExecutionContext{
		Args:       args,
		GlobalsPtr: uintptr(unsafe.Pointer(&args[0])),
		StackPtr:   uintptr(unsafe.Pointer(&args[1])),
		StackSize:  1024,
	}

	if ctx.Args != args {
		t.Errorf("ARM64 args not set correctly")
	}

	if ctx.StackSize != 1024 {
		t.Errorf("Expected stack size 1024, got %d", ctx.StackSize)
	}

	if ctx.GlobalsPtr == 0 {
		t.Error("GlobalsPtr should not be zero")
	}

	if ctx.StackPtr == 0 {
		t.Error("StackPtr should not be zero")
	}
}

// Benchmark tests for performance
func BenchmarkValueMarshaling(b *testing.B) {
	code := &CompiledCode{}
	value := &interpreter.Integer{Value: 42}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := code.valueToUint64(value)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValueUnmarshaling(b *testing.B) {
	code := &CompiledCode{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := code.uint64ToValue(42, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArgumentMarshaling(b *testing.B) {
	code := &CompiledCode{}
	args := []interpreter.Value{
		&interpreter.Integer{Value: 10},
		&interpreter.Integer{Value: 20},
		&interpreter.Boolean{Value: true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := code.marshalArgsToARM64(args)
		if err != nil {
			b.Fatal(err)
		}
	}
}