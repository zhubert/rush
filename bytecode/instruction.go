package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Opcode represents a bytecode instruction opcode
type Opcode byte

// Bytecode instruction set for Rush language
const (
	// Stack operations
	OpConstant Opcode = iota // Push constant to stack
	OpPop                    // Pop value from stack
	OpDup                    // Duplicate top of stack
	OpSwap                   // Swap top two stack values

	// Literals
	OpTrue   // Push true to stack
	OpFalse  // Push false to stack
	OpNull   // Push null to stack

	// Arithmetic operations
	OpAdd // Pop two values, push sum
	OpSub // Pop two values, push difference  
	OpMul // Pop two values, push product
	OpDiv // Pop two values, push quotient
	OpMod // Pop two values, push remainder

	// Comparison operations
	OpEqual       // Pop two values, push equality result
	OpNotEqual    // Pop two values, push inequality result
	OpGreaterThan // Pop two values, push greater than result
	OpLessThan    // Pop two values, push less than result
	OpGreaterEqual // Pop two values, push greater equal result
	OpLessEqual    // Pop two values, push less equal result

	// Logical operations
	OpAnd // Pop two values, push logical AND
	OpOr  // Pop two values, push logical OR
	OpNot // Pop one value, push logical NOT

	// Unary operations
	OpMinus // Pop one value, push negation

	// Variable operations
	OpGetGlobal // Push global variable to stack
	OpSetGlobal // Pop value, set global variable
	OpGetLocal  // Push local variable to stack
	OpSetLocal  // Pop value, set local variable
	OpGetFree   // Push free variable (closure) to stack
	OpSetFree   // Pop value, set free variable

	// Array operations
	OpArray     // Pop n elements, create array, push to stack
	OpIndex     // Pop index and array, push element
	OpSetIndex  // Pop value, index, and array; set element

	// Hash operations
	OpHash      // Pop 2n elements, create hash, push to stack
	OpGetHash   // Pop key and hash, push value
	OpSetHash   // Pop value, key, and hash; set element

	// Property/method access
	OpGetProperty // Pop object, push property value
	OpSetProperty // Pop value and object, set property
	OpCallMethod  // Pop args, receiver; call method

	// Control flow
	OpJump          // Unconditional jump
	OpJumpNotTruthy // Jump if top of stack is not truthy
	OpJumpTruthy    // Jump if top of stack is truthy

	// Function operations
	OpCall     // Pop args and function, call function
	OpReturn   // Return with value from stack
	OpReturnVoid // Return without value

	// Closure operations
	OpClosure      // Create closure with free variables
	OpGetBuiltin   // Push builtin function to stack
	OpCurrentClosure // Push current closure to stack

	// Loop control
	OpBreak    // Break from loop
	OpContinue // Continue loop

	// Exception handling
	OpThrow     // Throw exception
	OpTryBegin  // Begin try block
	OpTryEnd    // End try block
	OpCatch     // Catch exception
	OpFinally   // Finally block

	// Class operations
	OpClass        // Create class
	OpInherit      // Set up inheritance
	OpGetSuper     // Get superclass method
	OpMethod       // Define method
	OpInvoke       // Invoke method
	OpGetInstance  // Get instance variable
	OpSetInstance  // Set instance variable

	// Module operations
	OpImport // Import module
	OpExport // Export value

	// Switch operations
	OpSwitch     // Begin switch statement
	OpCase       // Case comparison
	OpDefault    // Default case

	// Advanced operations
	OpIterator   // Create iterator
	OpIterNext   // Get next iteration value
	OpIterDone   // Check if iteration is done

	// Dot notation method operations
	OpStringMethod    // String method call
	OpArrayMethod     // Array method call  
	OpNumberMethod    // Number method call
	OpHashMethod      // Hash method call
	OpJSONMethod      // JSON method call
	OpFileMethod      // File method call
	OpDirectoryMethod // Directory method call
	OpPathMethod      // Path method call
	OpTimeMethod      // Time method call
	OpDurationMethod  // Duration method call
	OpTimezoneMethod  // Timezone method call
)

// Definition holds information about an instruction
type Definition struct {
	Name          string // Human-readable name
	OperandWidths []int  // Width of each operand in bytes
}

// definitions maps opcodes to their definitions
var definitions = map[Opcode]*Definition{
	OpConstant:        {"OpConstant", []int{2}},        // 2-byte constant index
	OpPop:             {"OpPop", []int{}},
	OpDup:             {"OpDup", []int{}},
	OpSwap:            {"OpSwap", []int{}},
	OpTrue:            {"OpTrue", []int{}},
	OpFalse:           {"OpFalse", []int{}},
	OpNull:            {"OpNull", []int{}},
	OpAdd:             {"OpAdd", []int{}},
	OpSub:             {"OpSub", []int{}},
	OpMul:             {"OpMul", []int{}},
	OpDiv:             {"OpDiv", []int{}},
	OpMod:             {"OpMod", []int{}},
	OpEqual:           {"OpEqual", []int{}},
	OpNotEqual:        {"OpNotEqual", []int{}},
	OpGreaterThan:     {"OpGreaterThan", []int{}},
	OpLessThan:        {"OpLessThan", []int{}},
	OpGreaterEqual:    {"OpGreaterEqual", []int{}},
	OpLessEqual:       {"OpLessEqual", []int{}},
	OpAnd:             {"OpAnd", []int{}},
	OpOr:              {"OpOr", []int{}},
	OpNot:             {"OpNot", []int{}},
	OpMinus:           {"OpMinus", []int{}},
	OpGetGlobal:       {"OpGetGlobal", []int{2}},       // 2-byte global index
	OpSetGlobal:       {"OpSetGlobal", []int{2}},       // 2-byte global index
	OpGetLocal:        {"OpGetLocal", []int{1}},        // 1-byte local index
	OpSetLocal:        {"OpSetLocal", []int{1}},        // 1-byte local index
	OpGetFree:         {"OpGetFree", []int{1}},         // 1-byte free variable index
	OpSetFree:         {"OpSetFree", []int{1}},         // 1-byte free variable index
	OpArray:           {"OpArray", []int{2}},           // 2-byte element count
	OpIndex:           {"OpIndex", []int{}},
	OpSetIndex:        {"OpSetIndex", []int{}},
	OpHash:            {"OpHash", []int{2}},            // 2-byte key-value pair count
	OpGetHash:         {"OpGetHash", []int{}},
	OpSetHash:         {"OpSetHash", []int{}},
	OpGetProperty:     {"OpGetProperty", []int{2}},     // 2-byte property name index
	OpSetProperty:     {"OpSetProperty", []int{2}},     // 2-byte property name index
	OpCallMethod:      {"OpCallMethod", []int{2, 1}},   // 2-byte method name, 1-byte arg count
	OpJump:            {"OpJump", []int{2}},            // 2-byte jump offset
	OpJumpNotTruthy:   {"OpJumpNotTruthy", []int{2}},   // 2-byte jump offset
	OpJumpTruthy:      {"OpJumpTruthy", []int{2}},      // 2-byte jump offset
	OpCall:            {"OpCall", []int{1}},            // 1-byte argument count
	OpReturn:          {"OpReturn", []int{}},
	OpReturnVoid:      {"OpReturnVoid", []int{}},
	OpClosure:         {"OpClosure", []int{2, 1}},      // 2-byte constant index, 1-byte free var count
	OpGetBuiltin:      {"OpGetBuiltin", []int{1}},      // 1-byte builtin index
	OpCurrentClosure:  {"OpCurrentClosure", []int{}},
	OpBreak:           {"OpBreak", []int{}},
	OpContinue:        {"OpContinue", []int{}},
	OpThrow:           {"OpThrow", []int{}},
	OpTryBegin:        {"OpTryBegin", []int{2}},        // 2-byte catch handler offset
	OpTryEnd:          {"OpTryEnd", []int{}},
	OpCatch:           {"OpCatch", []int{2}},           // 2-byte exception type index
	OpFinally:         {"OpFinally", []int{}},
	OpClass:           {"OpClass", []int{2, 1}},        // 2-byte class name, 1-byte method count
	OpInherit:         {"OpInherit", []int{}},
	OpGetSuper:        {"OpGetSuper", []int{2}},        // 2-byte method name index
	OpMethod:          {"OpMethod", []int{2}},          // 2-byte method name index
	OpInvoke:          {"OpInvoke", []int{2, 1}},       // 2-byte method name, 1-byte arg count
	OpGetInstance:     {"OpGetInstance", []int{2}},     // 2-byte instance var name index
	OpSetInstance:     {"OpSetInstance", []int{2}},     // 2-byte instance var name index
	OpImport:          {"OpImport", []int{2}},          // 2-byte module name index
	OpExport:          {"OpExport", []int{2}},          // 2-byte export name index
	OpSwitch:          {"OpSwitch", []int{1}},          // 1-byte case count
	OpCase:            {"OpCase", []int{2}},            // 2-byte jump offset
	OpDefault:         {"OpDefault", []int{2}},         // 2-byte jump offset
	OpIterator:        {"OpIterator", []int{}},
	OpIterNext:        {"OpIterNext", []int{}},
	OpIterDone:        {"OpIterDone", []int{}},
	OpStringMethod:    {"OpStringMethod", []int{1, 1}}, // 1-byte method index, 1-byte arg count
	OpArrayMethod:     {"OpArrayMethod", []int{1, 1}},  // 1-byte method index, 1-byte arg count
	OpNumberMethod:    {"OpNumberMethod", []int{1, 1}}, // 1-byte method index, 1-byte arg count
	OpHashMethod:      {"OpHashMethod", []int{1, 1}},   // 1-byte method index, 1-byte arg count
	OpJSONMethod:      {"OpJSONMethod", []int{1, 1}},   // 1-byte method index, 1-byte arg count
	OpFileMethod:      {"OpFileMethod", []int{1, 1}},   // 1-byte method index, 1-byte arg count
	OpDirectoryMethod: {"OpDirectoryMethod", []int{1, 1}}, // 1-byte method index, 1-byte arg count
	OpPathMethod:      {"OpPathMethod", []int{1, 1}},   // 1-byte method index, 1-byte arg count
	OpTimeMethod:      {"OpTimeMethod", []int{1, 1}},   // 1-byte method index, 1-byte arg count
	OpDurationMethod:  {"OpDurationMethod", []int{1, 1}}, // 1-byte method index, 1-byte arg count
	OpTimezoneMethod:  {"OpTimezoneMethod", []int{1, 1}}, // 1-byte method index, 1-byte arg count
}

// Lookup returns the definition for an opcode
func Lookup(op Opcode) (*Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make creates a bytecode instruction from an opcode and operands
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 // opcode byte
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}

	return instruction
}

// ReadOperands extracts operands from an instruction
func ReadOperands(def *Definition, ins []byte) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ins[offset])
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}

// Instructions represents a sequence of bytecode instructions
type Instructions []byte

// String returns a human-readable representation of the instructions
func (ins Instructions) String() string {
	var out strings.Builder
	
	i := 0
	for i < len(ins) {
		def, err := Lookup(Opcode(ins[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s", i, def.Name)
		
		for _, operand := range operands {
			fmt.Fprintf(&out, " %d", operand)
		}
		
		out.WriteString("\n")
		i += 1 + read
	}

	return out.String()
}

// FlattenInstructions concatenates multiple instruction slices
func FlattenInstructions(s []Instructions) Instructions {
	var out Instructions
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

// ReadUint16 reads a 16-bit unsigned integer from bytecode
func ReadUint16(ins []byte) uint16 {
	return binary.BigEndian.Uint16(ins)
}