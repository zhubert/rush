package jit

import (
	"fmt"

	"rush/bytecode"
)

// ARM64CodeGen generates ARM64 assembly code from Rush bytecode
type ARM64CodeGen struct {
	code        []byte            // Generated machine code
	labels      map[int]int       // Jump target labels
	relocations []Relocation      // Addresses that need fixing up
}

// Relocation represents a jump target that needs to be resolved
type Relocation struct {
	Offset int  // Offset in code where address needs to be patched
	Target int  // Target instruction index
	Type   byte // Type of relocation (branch, call, etc.)
}

// ARM64 instruction encoding helpers
const (
	// ARM64 register numbers
	X0  = 0  // First argument register
	X1  = 1  // Second argument register
	X2  = 2  // Third argument register
	X3  = 3  // Fourth argument register
	X8  = 8  // Return value register
	X9  = 9  // Temporary register
	X10 = 10 // Temporary register
	X11 = 11 // Temporary register
	SP  = 31 // Stack pointer
	XZR = 31 // Zero register (in different contexts)
	
	// ARM64 instruction opcodes (simplified)
	ARM64_ADD_IMM = 0x91000000  // ADD (immediate)
	ARM64_SUB_IMM = 0xD1000000  // SUB (immediate)
	ARM64_ADD_REG = 0x8B000000  // ADD (register)
	ARM64_SUB_REG = 0xCB000000  // SUB (register)
	ARM64_MUL     = 0x9B000000  // MUL
	ARM64_SDIV    = 0x9AC00C00  // SDIV (signed divide)
	ARM64_CMP_REG = 0xEB000000  // CMP (register)
	ARM64_B       = 0x14000000  // B (branch)
	ARM64_BEQ     = 0x54000000  // B.EQ (branch if equal)
	ARM64_BNE     = 0x54000001  // B.NE (branch if not equal)
	ARM64_BGT     = 0x5400000C  // B.GT (branch if greater)
	ARM64_BLT     = 0x5400000B  // B.LT (branch if less)
	ARM64_LDR_IMM = 0xF9400000  // LDR (immediate)
	ARM64_STR_IMM = 0xF9000000  // STR (immediate)
	ARM64_MOV_IMM = 0xD2800000  // MOV (immediate)
	ARM64_RET     = 0xD65F03C0  // RET
	ARM64_BL      = 0x94000000  // BL (branch with link)
)

// NewARM64CodeGen creates a new ARM64 code generator
func NewARM64CodeGen() *ARM64CodeGen {
	return &ARM64CodeGen{
		code:        make([]byte, 0, 4096),
		labels:      make(map[int]int),
		relocations: make([]Relocation, 0),
	}
}

// Generate compiles bytecode instructions to ARM64 machine code
func (g *ARM64CodeGen) Generate(instructions bytecode.Instructions) ([]byte, error) {
	// Reset state
	g.code = g.code[:0]
	g.labels = make(map[int]int)
	g.relocations = g.relocations[:0]
	
	// Generate function prologue
	g.emitPrologue()
	
	// Process each bytecode instruction
	for ip := 0; ip < len(instructions); {
		// Record label position for jumps
		g.labels[ip] = len(g.code)
		
		opcode := bytecode.Opcode(instructions[ip])
		
		switch opcode {
		case bytecode.OpConstant:
			// Load constant to stack
			constIndex := int(instructions[ip+1])<<8 | int(instructions[ip+2])
			g.emitLoadConstant(constIndex)
			ip += 3
			
		case bytecode.OpPop:
			// Pop value from stack
			g.emitPop()
			ip++
			
		case bytecode.OpAdd:
			// Pop two values, add, push result
			g.emitBinaryOp(ARM64_ADD_REG)
			ip++
			
		case bytecode.OpSub:
			// Pop two values, subtract, push result
			g.emitBinaryOp(ARM64_SUB_REG)
			ip++
			
		case bytecode.OpMul:
			// Pop two values, multiply, push result
			g.emitBinaryOp(ARM64_MUL)
			ip++
			
		case bytecode.OpDiv:
			// Pop two values, divide, push result
			g.emitBinaryOp(ARM64_SDIV)
			ip++
			
		case bytecode.OpTrue:
			// Push true to stack
			g.emitPushBoolean(true)
			ip++
			
		case bytecode.OpFalse:
			// Push false to stack
			g.emitPushBoolean(false)
			ip++
			
		case bytecode.OpEqual:
			// Pop two values, compare for equality
			g.emitComparison(ARM64_BEQ)
			ip++
			
		case bytecode.OpNotEqual:
			// Pop two values, compare for inequality
			g.emitComparison(ARM64_BNE)
			ip++
			
		case bytecode.OpGreaterThan:
			// Pop two values, compare for greater than
			g.emitComparison(ARM64_BGT)
			ip++
			
		case bytecode.OpLessThan:
			// Pop two values, compare for less than
			g.emitComparison(ARM64_BLT)
			ip++
			
		case bytecode.OpJump:
			// Unconditional jump
			target := int(instructions[ip+1])<<8 | int(instructions[ip+2])
			g.emitJump(target)
			ip += 3
			
		case bytecode.OpJumpNotTruthy:
			// Conditional jump if top of stack is not truthy
			target := int(instructions[ip+1])<<8 | int(instructions[ip+2])
			g.emitConditionalJump(target, false)
			ip += 3
			
		case bytecode.OpReturn:
			// Return from function
			g.emitReturn()
			ip++
			
		default:
			return nil, fmt.Errorf("unsupported opcode for JIT compilation: %v", opcode)
		}
	}
	
	// Generate function epilogue
	g.emitEpilogue()
	
	// Resolve relocations (fix up jump targets)
	if err := g.resolveRelocations(); err != nil {
		return nil, fmt.Errorf("failed to resolve relocations: %v", err)
	}
	
	// Return generated code
	result := make([]byte, len(g.code))
	copy(result, g.code)
	return result, nil
}

// emitPrologue generates function entry code
func (g *ARM64CodeGen) emitPrologue() {
	// Standard ARM64 function prologue following AAPCS64
	// Save frame pointer (X29) and link register (X30)
	// stp x29, x30, [sp, #-16]!
	g.emit32(0xA9BF7FFF) // STP X29, X30, [SP, #-16]!
	
	// Set up frame pointer: mov x29, sp
	g.emit32(0x910003FD) // MOV X29, SP
	
	// Allocate stack space for locals (adjust as needed)
	// sub sp, sp, #64
	g.emit32(ARM64_SUB_IMM | (SP << 0) | (SP << 5) | (64 << 10))
	
	// Arguments are already in X0-X7 per ARM64 ABI
	// Globals pointer in X8, Stack pointer in X9, Stack size in X10
}

// emitEpilogue generates function exit code
func (g *ARM64CodeGen) emitEpilogue() {
	// Standard ARM64 function epilogue following AAPCS64
	// Deallocate stack space: add sp, sp, #64
	g.emit32(ARM64_ADD_IMM | (SP << 0) | (SP << 5) | (64 << 10))
	
	// Restore frame pointer and link register
	// ldp x29, x30, [sp], #16
	g.emit32(0xA8C17FFF) // LDP X29, X30, [SP], #16
	
	// Return to caller
	g.emit32(ARM64_RET) // RET
}

// emitLoadConstant loads a constant onto the stack
func (g *ARM64CodeGen) emitLoadConstant(constIndex int) {
	// This is a simplified version - in practice, we'd need to:
	// 1. Load constant from constant pool
	// 2. Convert to Rush value representation
	// 3. Push onto stack
	
	// For now, emit placeholder instructions
	g.emit32(ARM64_MOV_IMM | (X10 << 5) | (uint32(constIndex) << 5)) // MOV X10, #constIndex
}

// emitPop pops a value from the stack
func (g *ARM64CodeGen) emitPop() {
	// Decrement stack pointer and load value
	// Simplified version
	g.emit32(ARM64_SUB_IMM | (SP << 5) | (SP << 0) | (8 << 10)) // SUB SP, SP, #8
}

// emitBinaryOp emits code for binary operations
func (g *ARM64CodeGen) emitBinaryOp(opcode uint32) {
	// Pop two operands, perform operation, push result
	// This is a simplified version
	
	// Load operands (would need proper stack management)
	g.emit32(ARM64_LDR_IMM | (X10 << 0) | (SP << 5))     // LDR X10, [SP] (right operand)
	g.emit32(ARM64_LDR_IMM | (X11 << 0) | (SP << 5) | 8)  // LDR X11, [SP, #8] (left operand)
	
	// Perform operation
	g.emit32(opcode | (X9 << 0) | (X11 << 5) | (X10 << 16)) // OP X9, X11, X10
	
	// Store result back to stack
	g.emit32(ARM64_STR_IMM | (X9 << 0) | (SP << 5) | 8)   // STR X9, [SP, #8]
	g.emit32(ARM64_ADD_IMM | (SP << 0) | (SP << 5) | (8 << 10)) // ADD SP, SP, #8 (adjust stack)
}

// emitPushBoolean pushes a boolean value onto the stack
func (g *ARM64CodeGen) emitPushBoolean(value bool) {
	var imm uint32 = 0
	if value {
		imm = 1
	}
	
	g.emit32(ARM64_MOV_IMM | (X9 << 0) | (imm << 5))      // MOV X9, #value
	g.emit32(ARM64_SUB_IMM | (SP << 0) | (SP << 5) | (8 << 10)) // SUB SP, SP, #8
	g.emit32(ARM64_STR_IMM | (X9 << 0) | (SP << 5))       // STR X9, [SP]
}

// emitComparison emits code for comparison operations
func (g *ARM64CodeGen) emitComparison(branchOp uint32) {
	// Pop two operands, compare, push boolean result
	g.emit32(ARM64_LDR_IMM | (X10 << 0) | (SP << 5))      // LDR X10, [SP]
	g.emit32(ARM64_LDR_IMM | (X11 << 0) | (SP << 5) | 8)  // LDR X11, [SP, #8]
	g.emit32(ARM64_CMP_REG | (X11 << 5) | (X10 << 16))    // CMP X11, X10
	
	// Set result based on comparison
	g.emit32(ARM64_MOV_IMM | (X9 << 0))                   // MOV X9, #0 (default false)
	// Conditional move for true case would go here
	
	g.emit32(ARM64_STR_IMM | (X9 << 0) | (SP << 5) | 8)   // STR X9, [SP, #8]
	g.emit32(ARM64_ADD_IMM | (SP << 0) | (SP << 5) | (8 << 10)) // ADD SP, SP, #8
}

// emitJump emits an unconditional jump
func (g *ARM64CodeGen) emitJump(target int) {
	// Add relocation for jump target
	g.relocations = append(g.relocations, Relocation{
		Offset: len(g.code),
		Target: target,
		Type:   0, // Branch type
	})
	
	// Emit branch instruction with placeholder offset
	g.emit32(ARM64_B) // B #0 (offset will be patched)
}

// emitConditionalJump emits a conditional jump
func (g *ARM64CodeGen) emitConditionalJump(target int, jumpIfTrue bool) {
	// Pop condition from stack
	g.emit32(ARM64_LDR_IMM | (X10 << 0) | (SP << 5))      // LDR X10, [SP]
	g.emit32(ARM64_ADD_IMM | (SP << 0) | (SP << 5) | (8 << 10)) // ADD SP, SP, #8
	g.emit32(ARM64_CMP_REG | (X10 << 5))                  // CMP X10, #0
	
	// Add relocation for jump target
	g.relocations = append(g.relocations, Relocation{
		Offset: len(g.code),
		Target: target,
		Type:   1, // Conditional branch type
	})
	
	// Emit conditional branch
	if jumpIfTrue {
		g.emit32(ARM64_BNE) // B.NE #0 (jump if not zero/true)
	} else {
		g.emit32(ARM64_BEQ) // B.EQ #0 (jump if zero/false)
	}
}

// emitReturn emits function return
func (g *ARM64CodeGen) emitReturn() {
	// Pop return value and return
	g.emit32(ARM64_LDR_IMM | (X0 << 0) | (SP << 5))       // LDR X0, [SP] (return value)
	g.emit32(ARM64_RET)                                    // RET
}

// emit32 emits a 32-bit ARM64 instruction
func (g *ARM64CodeGen) emit32(instruction uint32) {
	g.code = append(g.code, 
		byte(instruction),
		byte(instruction >> 8),
		byte(instruction >> 16),
		byte(instruction >> 24))
}

// resolveRelocations fixes up jump targets in the generated code
func (g *ARM64CodeGen) resolveRelocations() error {
	for _, reloc := range g.relocations {
		targetAddr, exists := g.labels[reloc.Target]
		if !exists {
			return fmt.Errorf("undefined jump target: %d", reloc.Target)
		}
		
		// Calculate relative offset
		currentAddr := reloc.Offset
		offset := (targetAddr - currentAddr) / 4 // ARM64 uses word-aligned offsets
		
		// Check if offset fits in instruction
		if offset < -0x2000000 || offset > 0x1FFFFFF {
			return fmt.Errorf("jump offset too large: %d", offset)
		}
		
		// Patch the instruction
		instruction := uint32(g.code[reloc.Offset]) |
			uint32(g.code[reloc.Offset+1])<<8 |
			uint32(g.code[reloc.Offset+2])<<16 |
			uint32(g.code[reloc.Offset+3])<<24
		
		// Insert offset into instruction (bits 0-25 for branches)
		instruction |= uint32(offset&0x3FFFFFF)
		
		// Write back patched instruction
		g.code[reloc.Offset] = byte(instruction)
		g.code[reloc.Offset+1] = byte(instruction >> 8)
		g.code[reloc.Offset+2] = byte(instruction >> 16)
		g.code[reloc.Offset+3] = byte(instruction >> 24)
	}
	
	return nil
}