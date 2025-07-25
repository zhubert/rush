package bytecode

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpConstant, []int{65535}, []byte{byte(OpConstant), 255, 255}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}
		
		for i, b := range tt.expected {
			if instruction[i] != b {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpClosure, 65535, 255),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpClosure 65535 255
`

	concatted := FlattenInstructions(instructions)
	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpGetLocal, []int{255}, 1},
		{OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		
		def, err := Lookup(tt.op)
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}

func TestLookup(t *testing.T) {
	tests := []struct {
		op       Opcode
		expected *Definition
	}{
		{OpConstant, &Definition{"OpConstant", []int{2}}},
		{OpAdd, &Definition{"OpAdd", []int{}}},
		{OpGetLocal, &Definition{"OpGetLocal", []int{1}}},
		{OpClosure, &Definition{"OpClosure", []int{2, 1}}},
	}

	for _, tt := range tests {
		def, err := Lookup(tt.op)
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		if def.Name != tt.expected.Name {
			t.Errorf("wrong name. want=%q, got=%q", tt.expected.Name, def.Name)
		}

		if len(def.OperandWidths) != len(tt.expected.OperandWidths) {
			t.Errorf("wrong operand widths length. want=%d, got=%d",
				len(tt.expected.OperandWidths), len(def.OperandWidths))
		}

		for i, width := range tt.expected.OperandWidths {
			if def.OperandWidths[i] != width {
				t.Errorf("wrong operand width at %d. want=%d, got=%d",
					i, width, def.OperandWidths[i])
			}
		}
	}
}

func TestLookupUndefinedOpcode(t *testing.T) {
	_, err := Lookup(Opcode(255))
	if err == nil {
		t.Error("expected error for undefined opcode")
	}
}

func TestMakeInvalidOpcode(t *testing.T) {
	instruction := Make(Opcode(255), 1, 2, 3)
	if len(instruction) != 0 {
		t.Errorf("expected empty instruction for invalid opcode, got %v", instruction)
	}
}

func TestFlattenInstructions(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 1),
		Make(OpPop),
	}

	expected := Instructions{
		byte(OpAdd),
		byte(OpConstant), 0, 1,
		byte(OpPop),
	}

	result := FlattenInstructions(instructions)
	
	if len(result) != len(expected) {
		t.Errorf("wrong length. want=%d, got=%d", len(expected), len(result))
	}

	for i, b := range expected {
		if result[i] != b {
			t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, result[i])
		}
	}
}

func TestReadUint16(t *testing.T) {
	tests := []struct {
		input    []byte
		expected uint16
	}{
		{[]byte{0, 1}, 1},
		{[]byte{255, 255}, 65535},
		{[]byte{1, 0}, 256},
	}

	for _, tt := range tests {
		result := ReadUint16(tt.input)
		if result != tt.expected {
			t.Errorf("wrong result. want=%d, got=%d", tt.expected, result)
		}
	}
}

func TestAllOpcodesHaveDefinitions(t *testing.T) {
	for op := OpConstant; op <= OpTimezoneMethod; op++ {
		_, err := Lookup(op)
		if err != nil {
			t.Errorf("opcode %d has no definition", op)
		}
	}
}

func TestDefinitionOperandWidths(t *testing.T) {
	tests := []struct {
		op              Opcode
		expectedWidths  []int
		description     string
	}{
		{OpConstant, []int{2}, "constant index"},
		{OpAdd, []int{}, "no operands"},
		{OpGetLocal, []int{1}, "local variable index"},
		{OpSetLocal, []int{1}, "local variable index"},
		{OpGetGlobal, []int{2}, "global variable index"},
		{OpSetGlobal, []int{2}, "global variable index"},
		{OpArray, []int{2}, "element count"},
		{OpHash, []int{2}, "key-value pair count"},
		{OpJump, []int{2}, "jump offset"},
		{OpJumpNotTruthy, []int{2}, "jump offset"},
		{OpCall, []int{1}, "argument count"},
		{OpClosure, []int{2, 1}, "constant index and free variable count"},
		{OpClass, []int{2, 1}, "class name index and method count"},
	}

	for _, tt := range tests {
		def, err := Lookup(tt.op)
		if err != nil {
			t.Fatalf("definition not found for %s: %q", tt.description, err)
		}

		if len(def.OperandWidths) != len(tt.expectedWidths) {
			t.Errorf("wrong operand widths length for %s. want=%d, got=%d",
				tt.description, len(tt.expectedWidths), len(def.OperandWidths))
			continue
		}

		for i, expectedWidth := range tt.expectedWidths {
			if def.OperandWidths[i] != expectedWidth {
				t.Errorf("wrong operand width for %s at position %d. want=%d, got=%d",
					tt.description, i, expectedWidth, def.OperandWidths[i])
			}
		}
	}
}