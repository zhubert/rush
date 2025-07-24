package compiler

import (
	"fmt"
	"testing"

	"rush/ast"
	"rush/bytecode"
	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []bytecode.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "2 % 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpNot),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) { 10 }; 3333;
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 11),
				// 0010
				bytecode.Make(bytecode.OpNull),
				// 0011
				bytecode.Make(bytecode.OpPop),
				// 0012
				bytecode.Make(bytecode.OpConstant, 1),
				// 0015
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			if (true) { 10 } else { 20 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 13),
				// 0010
				bytecode.Make(bytecode.OpConstant, 1),
				// 0013
				bytecode.Make(bytecode.OpPop),
				// 0014
				bytecode.Make(bytecode.OpConstant, 2),
				// 0017
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			one = 1;
			two = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
			},
		},
		{
			input: `
			one = 1;
			one;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			one = 1;
			two = one;
			two;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"rush"`,
			expectedConstants: []interface{}{"rush"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             `"ru" + "sh"`,
			expectedConstants: []interface{}{"ru", "sh"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpHash, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpHash, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpHash, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []interface{}{1, 2, 2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpHash, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `fn() { 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `fn() { 1; 2 }`,
			expectedConstants: []interface{}{
				1,
				2,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`,
			expectedConstants: []interface{}{
				24,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			noArg = fn() { 24 };
			noArg();
			`,
			expectedConstants: []interface{}{
				24,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			oneArg = fn(a) { a };
			oneArg(24);
			`,
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturn),
				},
				24,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			manyArg = fn(a, b, c) { a; b; c };
			manyArg(24, 25, 26);
			`,
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 2),
					bytecode.Make(bytecode.OpReturn),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpCall, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			num = 55;
			fn() { num }
			`,
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			fn() {
				num = 55;
				num
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			fn() {
				a = 55;
				b = 77;
				a + b
			}
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 1),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			len([]);
			push([], 1);
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpGetBuiltin, 0),
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetBuiltin, 5),
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpCall, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `fn() { len([]) }`,
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetBuiltin, 0),
					bytecode.Make(bytecode.OpArray, 0),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			fn(a) {
				fn(b) {
					a + b
				}
			}
			`,
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 0, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			fn(a) {
				fn(b) {
					fn(c) {
						a + b + c
					}
				}
			};
			`,
			expectedConstants: []interface{}{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpGetFree, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 0, 2),
					bytecode.Make(bytecode.OpReturn),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 1, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 2, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			global = 55;
			fn() {
				a = 66;
				fn() {
					b = 77;
					fn() {
						c = 88;
						global + a + b + c;
					}
				}
			}
			`,
			expectedConstants: []interface{}{
				55,
				66,
				77,
				88,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 3),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetFree, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturn),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 2),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 4, 2),
					bytecode.Make(bytecode.OpReturn),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpClosure, 5, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpClosure, 6, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			countdown = fn(x) {
				countdown(x - 1);
			};
			countdown(1);
			`,
			expectedConstants: []interface{}{
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpCurrentClosure),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSub),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturn),
				},
				1,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 1, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			wrapper = fn() {
				countdown = fn(x) {
					countdown(x - 1);
				};
				countdown(1);
			};
			wrapper();
			`,
			expectedConstants: []interface{}{
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpCurrentClosure),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSub),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturn),
				},
				1,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpClosure, 1, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 2),
					bytecode.Make(bytecode.OpCall, 1),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpClosure, 3, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestWhileLoops(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			i = 0;
			while (i < 3) {
				i = i + 1;
			}
			`,
			expectedConstants: []interface{}{0, 3, 1},
			expectedInstructions: []bytecode.Instructions{
				// i = 0
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// loop start: while (i < 3)
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpJumpNotTruthy, 22),
				// loop body: i = i + 1
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// jump back to condition
				bytecode.Make(bytecode.OpJump, 6),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestForLoops(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			for (i = 0; i < 3; i = i + 1) {
				i;
			}
			`,
			expectedConstants: []interface{}{0, 3, 1},
			expectedInstructions: []bytecode.Instructions{
				// initialization: i = 0
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// condition: i < 3
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpJumpNotTruthy, 25),
				// body: i;
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				// update: i = i + 1
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// jump back to condition
				bytecode.Make(bytecode.OpJump, 6),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestPropertyAccess(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `"hello".length`,
			expectedConstants: []interface{}{"hello", "length"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpGetProperty, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `[1, 2, 3].length`,
			expectedConstants: []interface{}{1, 2, 3, "length"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpGetProperty, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexAssignment(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			arr = [1, 2, 3];
			arr[1] = 5;
			`,
			expectedConstants: []interface{}{1, 2, 3, 1, 5},
			expectedInstructions: []bytecode.Instructions{
				// arr = [1, 2, 3]
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// arr[1] = 5
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpSetIndex),
			},
		},
		{
			input: `
			hash = {1: 2};
			hash[1] = 3;
			`,
			expectedConstants: []interface{}{1, 2, 1, 3},
			expectedInstructions: []bytecode.Instructions{
				// hash = {1: 2}
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpHash, 1),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				// hash[1] = 3
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSetIndex),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestSymbolTable(t *testing.T) {
	global := NewSymbolTable()
	
	a := global.Define("a")
	expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0}
	if a != expected {
		t.Errorf("expected a=%+v, got=%+v", expected, a)
	}

	b := global.Define("b")
	expected = Symbol{Name: "b", Scope: GlobalScope, Index: 1}
	if b != expected {
		t.Errorf("expected b=%+v, got=%+v", expected, b)
	}
}

func TestEnclosedSymbolTable(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}
}

func TestResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		{Name: "a", Scope: BuiltinScope, Index: 0},
		{Name: "c", Scope: BuiltinScope, Index: 1},
		{Name: "e", Scope: BuiltinScope, Index: 2},
		{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, sym := range expected {
		global.DefineBuiltin(i, sym.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, result)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
		expectedFree    []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
			[]Symbol{},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
			[]Symbol{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, result)
			}
		}

		if len(tt.table.FreeSymbols) != len(tt.expectedFree) {
			t.Errorf("wrong number of free symbols. got=%d, want=%d",
				len(tt.table.FreeSymbols), len(tt.expectedFree))
			continue
		}

		for i, sym := range tt.expectedFree {
			result := tt.table.FreeSymbols[i]
			if result != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v",
					result, sym)
			}
		}
	}
}

func TestResolveUnresolvable(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	expected := Symbol{Name: "a", Scope: GlobalScope, Index: 0}

	result, ok := global.Resolve("a")
	if !ok {
		t.Errorf("name a not resolvable")
	}
	if result != expected {
		t.Errorf("expected a to resolve to %+v, got=%+v", expected, result)
	}

	_, ok = global.Resolve("b")
	if ok {
		t.Errorf("name b resolved, but was expected not to")
	}
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []bytecode.Instructions, actual bytecode.Instructions) error {
	concatted := bytecode.FlattenInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

func testConstants(t *testing.T, expected []interface{}, actual []interpreter.Value) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s",
					i, err)
			}
		case []bytecode.Instructions:
			fn, ok := actual[i].(*interpreter.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T",
					i, actual[i])
			}

			err := testInstructions(constant, bytecode.Instructions(fn.Instructions))
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s",
					i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual interpreter.Value) error {
	result, ok := actual.(*interpreter.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual interpreter.Value) error {
	result, ok := actual.(*interpreter.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}

	return nil
}