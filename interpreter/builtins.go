package interpreter

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
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
	"ord": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			str, ok := args[0].(*String)
			if !ok {
				return newError("argument to `ord` must be STRING, got %s", args[0].Type())
			}

			if len(str.Value) != 1 {
				return newError("argument to `ord` must be a single character, got length %d", len(str.Value))
			}

			return &Integer{Value: int64(str.Value[0])}
		},
	},
	"chr": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			code, ok := args[0].(*Integer)
			if !ok {
				return newError("argument to `chr` must be INTEGER, got %s", args[0].Type())
			}

			if code.Value < 0 || code.Value > 127 {
				return newError("argument to `chr` must be between 0 and 127, got %d", code.Value)
			}

			return &String{Value: string(byte(code.Value))}
		},
	},
	// String functions - will be moved to std/string
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
	// Array functions - will be moved to std/array
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
	"to_string": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			
			// Use the built-in Inspect() method which provides proper string representation
			// for all value types including integers, floats, booleans, arrays, etc.
			return &String{Value: args[0].Inspect()}
		},
	},
	// Math functions
	"builtin_abs": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				if arg.Value < 0 {
					return &Integer{Value: -arg.Value}
				}
				return arg
			case *Float:
				return &Float{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to `builtin_abs` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"builtin_min": {
		Fn: func(args ...Value) Value {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=0, want at least 1")
			}

			var minVal Value = args[0]
			var minFloat float64
			var isFloat bool

			// Check first argument and set initial values
			switch val := args[0].(type) {
			case *Integer:
				minFloat = float64(val.Value)
			case *Float:
				minFloat = val.Value
				isFloat = true
			default:
				return newError("arguments to `builtin_min` must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			// Compare with remaining arguments
			for i := 1; i < len(args); i++ {
				switch val := args[i].(type) {
				case *Integer:
					if float64(val.Value) < minFloat {
						minFloat = float64(val.Value)
						minVal = val
					}
				case *Float:
					if val.Value < minFloat {
						minFloat = val.Value
						minVal = val
						isFloat = true
					}
				default:
					return newError("arguments to `builtin_min` must be INTEGER or FLOAT, got %s", args[i].Type())
				}
			}

			// Return appropriate type
			if isFloat {
				return &Float{Value: minFloat}
			}
			return minVal
		},
	},
	"builtin_max": {
		Fn: func(args ...Value) Value {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=0, want at least 1")
			}

			var maxVal Value = args[0]
			var maxFloat float64
			var isFloat bool

			// Check first argument and set initial values
			switch val := args[0].(type) {
			case *Integer:
				maxFloat = float64(val.Value)
			case *Float:
				maxFloat = val.Value
				isFloat = true
			default:
				return newError("arguments to `builtin_max` must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			// Compare with remaining arguments
			for i := 1; i < len(args); i++ {
				switch val := args[i].(type) {
				case *Integer:
					if float64(val.Value) > maxFloat {
						maxFloat = float64(val.Value)
						maxVal = val
					}
				case *Float:
					if val.Value > maxFloat {
						maxFloat = val.Value
						maxVal = val
						isFloat = true
					}
				default:
					return newError("arguments to `builtin_max` must be INTEGER or FLOAT, got %s", args[i].Type())
				}
			}

			// Return appropriate type
			if isFloat {
				return &Float{Value: maxFloat}
			}
			return maxVal
		},
	},
	"builtin_floor": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return arg // Already an integer
			case *Float:
				return &Integer{Value: int64(math.Floor(arg.Value))}
			default:
				return newError("argument to `builtin_floor` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"builtin_ceil": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return arg // Already an integer
			case *Float:
				return &Integer{Value: int64(math.Ceil(arg.Value))}
			default:
				return newError("argument to `builtin_ceil` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"builtin_round": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return arg // Already an integer
			case *Float:
				return &Integer{Value: int64(math.Round(arg.Value))}
			default:
				return newError("argument to `builtin_round` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"builtin_sqrt": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var val float64
			switch arg := args[0].(type) {
			case *Integer:
				val = float64(arg.Value)
			case *Float:
				val = arg.Value
			default:
				return newError("argument to `builtin_sqrt` must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			if val < 0 {
				return newError("argument to `builtin_sqrt` cannot be negative, got %f", val)
			}

			return &Float{Value: math.Sqrt(val)}
		},
	},
	"builtin_pow": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			var base, exp float64
			switch arg := args[0].(type) {
			case *Integer:
				base = float64(arg.Value)
			case *Float:
				base = arg.Value
			default:
				return newError("first argument to `builtin_pow` must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			switch arg := args[1].(type) {
			case *Integer:
				exp = float64(arg.Value)
			case *Float:
				exp = arg.Value
			default:
				return newError("second argument to `builtin_pow` must be INTEGER or FLOAT, got %s", args[1].Type())
			}

			result := math.Pow(base, exp)
			return &Float{Value: result}
		},
	},
	"builtin_random": {
		Fn: func(args ...Value) Value {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0", len(args))
			}

			// Initialize random seed once per session
			rand.Seed(time.Now().UnixNano())
			return &Float{Value: rand.Float64()}
		},
	},
	"builtin_random_int": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			minVal, ok := args[0].(*Integer)
			if !ok {
				return newError("first argument to `builtin_random_int` must be INTEGER, got %s", args[0].Type())
			}

			maxVal, ok := args[1].(*Integer)
			if !ok {
				return newError("second argument to `builtin_random_int` must be INTEGER, got %s", args[1].Type())
			}

			if minVal.Value > maxVal.Value {
				return newError("first argument to `builtin_random_int` cannot be greater than second argument")
			}

			// Initialize random seed once per session
			rand.Seed(time.Now().UnixNano())
			range_ := maxVal.Value - minVal.Value + 1
			result := minVal.Value + rand.Int63n(range_)
			return &Integer{Value: result}
		},
	},
	"builtin_sum": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			arr, ok := args[0].(*Array)
			if !ok {
				return newError("argument to `builtin_sum` must be ARRAY, got %s", args[0].Type())
			}

			var sum float64
			var hasFloat bool

			for _, element := range arr.Elements {
				switch val := element.(type) {
				case *Integer:
					sum += float64(val.Value)
				case *Float:
					sum += val.Value
					hasFloat = true
				default:
					return newError("array elements for `builtin_sum` must be INTEGER or FLOAT, got %s", element.Type())
				}
			}

			if hasFloat {
				return &Float{Value: sum}
			}
			return &Integer{Value: int64(sum)}
		},
	},
	"builtin_average": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			arr, ok := args[0].(*Array)
			if !ok {
				return newError("argument to `builtin_average` must be ARRAY, got %s", args[0].Type())
			}

			if len(arr.Elements) == 0 {
				return newError("cannot calculate average of empty array")
			}

			var sum float64

			for _, element := range arr.Elements {
				switch val := element.(type) {
				case *Integer:
					sum += float64(val.Value)
				case *Float:
					sum += val.Value
				default:
					return newError("array elements for `builtin_average` must be INTEGER or FLOAT, got %s", element.Type())
				}
			}

			average := sum / float64(len(arr.Elements))
			return &Float{Value: average}
		},
	},
	"builtin_is_number": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *Integer, *Float:
				return &Boolean{Value: true}
			default:
				return &Boolean{Value: false}
			}
		},
	},
	"builtin_is_integer": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *Integer:
				return &Boolean{Value: true}
			default:
				return &Boolean{Value: false}
			}
		},
	},
}