package interpreter

import (
	"fmt"
	"strings"
)

var builtins = map[string]*BuiltinFunction{
	"len": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...Value) Value {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				if arg.Type() == STRING_VALUE {
					fmt.Print(arg.(*String).Value)
				} else {
					fmt.Print(arg.Inspect())
				}
			}
			fmt.Println()
			return NULL
		},
	},
	"type": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &String{Value: string(args[0].Type())}
		},
	},
	"substr": {
		Fn: func(args ...Value) Value {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}

			str, ok := args[0].(*String)
			if !ok {
				return newError("argument to `substr` must be STRING, got %s", args[0].Type())
			}

			start, ok := args[1].(*Integer)
			if !ok {
				return newError("second argument to `substr` must be INTEGER, got %s", args[1].Type())
			}

			length, ok := args[2].(*Integer)
			if !ok {
				return newError("third argument to `substr` must be INTEGER, got %s", args[2].Type())
			}

			strLen := int64(len(str.Value))
			if start.Value < 0 || start.Value >= strLen {
				return &String{Value: ""}
			}

			end := start.Value + length.Value
			if end > strLen {
				end = strLen
			}

			return &String{Value: str.Value[start.Value:end]}
		},
	},
	"split": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			str, ok := args[0].(*String)
			if !ok {
				return newError("first argument to `split` must be STRING, got %s", args[0].Type())
			}

			sep, ok := args[1].(*String)
			if !ok {
				return newError("second argument to `split` must be STRING, got %s", args[1].Type())
			}

			parts := strings.Split(str.Value, sep.Value)
			elements := make([]Value, len(parts))
			for i, part := range parts {
				elements[i] = &String{Value: part}
			}

			return &Array{Elements: elements}
		},
	},
	"push": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			arr, ok := args[0].(*Array)
			if !ok {
				return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			newElements := make([]Value, len(arr.Elements)+1)
			copy(newElements, arr.Elements)
			newElements[len(arr.Elements)] = args[1]

			return &Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			arr, ok := args[0].(*Array)
			if !ok {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}

			if len(arr.Elements) == 0 {
				errorObj := newTypedError("IndexError", "pop from empty array", 0, 0)
				return NewException(errorObj)
			}

			return arr.Elements[len(arr.Elements)-1]
		},
	},
	"slice": {
		Fn: func(args ...Value) Value {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}

			arr, ok := args[0].(*Array)
			if !ok {
				return newError("first argument to `slice` must be ARRAY, got %s", args[0].Type())
			}

			start, ok := args[1].(*Integer)
			if !ok {
				return newError("second argument to `slice` must be INTEGER, got %s", args[1].Type())
			}

			end, ok := args[2].(*Integer)
			if !ok {
				return newError("third argument to `slice` must be INTEGER, got %s", args[2].Type())
			}

			arrLen := int64(len(arr.Elements))
			if start.Value < 0 || start.Value >= arrLen {
				return &Array{Elements: []Value{}}
			}

			endIdx := end.Value
			if endIdx > arrLen {
				endIdx = arrLen
			}
			if endIdx <= start.Value {
				return &Array{Elements: []Value{}}
			}

			newElements := make([]Value, endIdx-start.Value)
			copy(newElements, arr.Elements[start.Value:endIdx])

			return &Array{Elements: newElements}
		},
	},
	// Error constructors
	"Error": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to Error constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("Error", msg.Value, 0, 0)
		},
	},
	"ValidationError": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to ValidationError constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("ValidationError", msg.Value, 0, 0)
		},
	},
	"TypeError": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to TypeError constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("TypeError", msg.Value, 0, 0)
		},
	},
	"IndexError": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to IndexError constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("IndexError", msg.Value, 0, 0)
		},
	},
	"ArgumentError": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to ArgumentError constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("ArgumentError", msg.Value, 0, 0)
		},
	},
	"RuntimeError": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			msg, ok := args[0].(*String)
			if !ok {
				return newError("argument to RuntimeError constructor must be STRING, got %s", args[0].Type())
			}
			return newTypedError("RuntimeError", msg.Value, 0, 0)
		},
	},
}