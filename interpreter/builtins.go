package interpreter

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Builtins is a list of builtin function names for the compiler
var Builtins = []string{
	"JSON",
	"Time", 
	"Duration",
	"TimeZone",
	"Regexp",
	"len",
	"print",
	"puts",
	"type",
	"ord",
	"chr", 
	"substr",
	"split",
	"push",
	"pop",
	"slice",
	"Error",
	"ValidationError", 
	"TypeError",
	"IndexError",
	"ArgumentError",
	"RuntimeError",
	"to_string",
	"builtin_abs",
	"builtin_min",
	"builtin_max",
	"builtin_floor", 
	"builtin_ceil",
	"builtin_round",
	"builtin_sqrt",
	"builtin_pow",
	"builtin_random",
	"builtin_random_int",
	"builtin_sum", 
	"builtin_average",
	"builtin_hash_keys",
	"builtin_hash_values",
	"builtin_hash_has_key",
	"builtin_hash_get",
	"builtin_hash_set",
	"builtin_hash_delete",
	"builtin_hash_merge",
	"array_to_hash",
	"file",
	"directory", 
	"path",
}

// GetBuiltin returns a builtin function by name
func GetBuiltin(name string) (*BuiltinFunction, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

var builtins = map[string]*BuiltinFunction{
	"JSON": {
		Fn: func(args ...Value) Value {
			return &JSONNamespace{}
		},
	},
	"Time": {
		Fn: func(args ...Value) Value {
			return &TimeNamespace{}
		},
	},
	"Duration": {
		Fn: func(args ...Value) Value {
			return &DurationNamespace{}
		},
	},
	"TimeZone": {
		Fn: func(args ...Value) Value {
			return &TimeZoneNamespace{}
		},
	},
	"Regexp": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			pattern, ok := args[0].(*String)
			if !ok {
				return newError("argument to `Regexp` must be STRING, got %s", args[0].Type())
			}

			regex, err := regexp.Compile(pattern.Value)
			if err != nil {
				return newError("invalid regular expression: %s", err.Error())
			}

			return &Regexp{
				Pattern: pattern.Value,
				Regex:   regex,
			}
		},
	},
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
			case *Hash:
				return &Integer{Value: int64(len(arg.Keys))}
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
	"puts": {
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
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
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
					isFloat = true // Set to true whenever we encounter a float
					if val.Value > maxFloat {
						maxFloat = val.Value
						maxVal = val
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
	"builtin_is_number?": {
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
	"builtin_is_integer?": {
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
	"builtin_hash_keys": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("argument to `builtin_hash_keys` must be HASH, got %s", args[0].Type())
			}

			return &Array{Elements: hash.Keys}
		},
	},
	"builtin_hash_values": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("argument to `builtin_hash_values` must be HASH, got %s", args[0].Type())
			}

			values := []Value{}
			for _, key := range hash.Keys {
				hashKey := CreateHashKey(key)
				values = append(values, hash.Pairs[hashKey])
			}

			return &Array{Elements: values}
		},
	},
	"builtin_hash_has_key": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("first argument to `builtin_hash_has_key` must be HASH, got %s", args[0].Type())
			}

			key := args[1]
			if !isHashable(key) {
				return newError("unusable as hash key: %T", key)
			}

			hashKey := CreateHashKey(key)
			_, exists := hash.Pairs[hashKey]
			return &Boolean{Value: exists}
		},
	},
	"builtin_hash_get": {
		Fn: func(args ...Value) Value {
			if len(args) < 2 || len(args) > 3 {
				return newError("wrong number of arguments. got=%d, want=2 or 3", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("first argument to `builtin_hash_get` must be HASH, got %s", args[0].Type())
			}

			key := args[1]
			if !isHashable(key) {
				return newError("unusable as hash key: %T", key)
			}

			hashKey := CreateHashKey(key)
			value, exists := hash.Pairs[hashKey]
			if !exists {
				if len(args) == 3 {
					return args[2] // return default value
				}
				return NULL
			}

			return value
		},
	},
	"builtin_hash_set": {
		Fn: func(args ...Value) Value {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("first argument to `builtin_hash_set` must be HASH, got %s", args[0].Type())
			}

			key := args[1]
			if !isHashable(key) {
				return newError("unusable as hash key: %T", key)
			}

			value := args[2]

			// Create new hash (immutable operation)
			newPairs := make(map[HashKey]Value)
			newKeys := []Value{}

			// Copy existing pairs
			for _, k := range hash.Keys {
				hashKey := CreateHashKey(k)
				newPairs[hashKey] = hash.Pairs[hashKey]
				newKeys = append(newKeys, k)
			}

			// Add or update the key-value pair
			hashKey := CreateHashKey(key)
			if _, exists := newPairs[hashKey]; !exists {
				newKeys = append(newKeys, key)
			}
			newPairs[hashKey] = value

			return &Hash{Pairs: newPairs, Keys: newKeys}
		},
	},
	"builtin_hash_delete": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			hash, ok := args[0].(*Hash)
			if !ok {
				return newError("first argument to `builtin_hash_delete` must be HASH, got %s", args[0].Type())
			}

			key := args[1]
			if !isHashable(key) {
				return newError("unusable as hash key: %T", key)
			}

			hashKey := CreateHashKey(key)

			// Create new hash without the key (immutable operation)
			newPairs := make(map[HashKey]Value)
			newKeys := []Value{}

			for _, k := range hash.Keys {
				kHashKey := CreateHashKey(k)
				if kHashKey != hashKey {
					newPairs[kHashKey] = hash.Pairs[kHashKey]
					newKeys = append(newKeys, k)
				}
			}

			return &Hash{Pairs: newPairs, Keys: newKeys}
		},
	},
	"builtin_hash_merge": {
		Fn: func(args ...Value) Value {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			hash1, ok1 := args[0].(*Hash)
			if !ok1 {
				return newError("first argument to `builtin_hash_merge` must be HASH, got %s", args[0].Type())
			}

			hash2, ok2 := args[1].(*Hash)
			if !ok2 {
				return newError("second argument to `builtin_hash_merge` must be HASH, got %s", args[1].Type())
			}

			// Create new hash (immutable operation)
			newPairs := make(map[HashKey]Value)
			newKeys := []Value{}

			// Copy hash1
			for _, key := range hash1.Keys {
				hashKey := CreateHashKey(key)
				newPairs[hashKey] = hash1.Pairs[hashKey]
				newKeys = append(newKeys, key)
			}

			// Merge hash2 (overrides conflicts)
			for _, key := range hash2.Keys {
				hashKey := CreateHashKey(key)
				if _, exists := newPairs[hashKey]; !exists {
					newKeys = append(newKeys, key)
				}
				newPairs[hashKey] = hash2.Pairs[hashKey]
			}

			return &Hash{Pairs: newPairs, Keys: newKeys}
		},
	},
	
	"array_to_hash": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			
			array, ok := args[0].(*Array)
			if !ok {
				return newError("argument to `array_to_hash` must be ARRAY, got %s", args[0].Type())
			}
			
			pairs := make(map[HashKey]Value)
			keys := make([]Value, 0, len(array.Elements))
			
			for _, element := range array.Elements {
				pair, ok := element.(*Array)
				if !ok {
					return newError("array elements must be arrays of [key, value] pairs, got %s", element.Type())
				}
				
				if len(pair.Elements) < 2 {
					return newError("each pair must have at least 2 elements (key and value), got %d", len(pair.Elements))
				}
				
				key := pair.Elements[0]
				value := pair.Elements[1]
				hashKey := CreateHashKey(key)
				
				// Only add key if it doesn't exist yet (preserve first occurrence)
				if _, exists := pairs[hashKey]; !exists {
					keys = append(keys, key)
				}
				pairs[hashKey] = value
			}
			
			return &Hash{Pairs: pairs, Keys: keys}
		},
	},
	"file": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			path, ok := args[0].(*String)
			if !ok {
				return newError("argument to `file` must be STRING, got %s", args[0].Type())
			}

			// Validate path to prevent directory traversal attacks
			if strings.Contains(path.Value, "..") {
				return newError("invalid file path: path traversal not allowed")
			}

			return &File{
				Path:   path.Value,
				Handle: nil,
				IsOpen: false,
			}
		},
	},
	"directory": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			path, ok := args[0].(*String)
			if !ok {
				return newError("argument to `directory` must be STRING, got %s", args[0].Type())
			}

			// Validate path to prevent directory traversal attacks
			if strings.Contains(path.Value, "..") {
				return newError("invalid directory path: path traversal not allowed")
			}

			return &Directory{
				Path: path.Value,
			}
		},
	},
	"path": {
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			value, ok := args[0].(*String)
			if !ok {
				return newError("argument to `path` must be STRING, got %s", args[0].Type())
			}

			return &Path{
				Value: value.Value,
			}
		},
	},
}

// parseJSON converts a JSON string to a Rush JSON object
func parseJSON(jsonStr string) Value {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return newError("invalid JSON: %s", err.Error())
	}

	rushValue := convertFromGoValue(data)
	return &JSON{Data: rushValue}
}

// stringifyValue converts a Rush value to a JSON string
func stringifyValue(value Value) (string, error) {
	goValue, err := convertToGoValue(value)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(goValue)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// convertFromGoValue converts Go interface{} from json.Unmarshal to Rush Value
func convertFromGoValue(data interface{}) Value {
	switch v := data.(type) {
	case nil:
		return NULL
	case bool:
		return &Boolean{Value: v}
	case float64:
		// JSON numbers are always float64
		if v == float64(int64(v)) {
			return &Integer{Value: int64(v)}
		}
		return &Float{Value: v}
	case string:
		return &String{Value: v}
	case []interface{}:
		elements := make([]Value, len(v))
		for i, elem := range v {
			elements[i] = convertFromGoValue(elem)
		}
		return &Array{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[HashKey]Value)
		keys := make([]Value, 0, len(v))
		for k, val := range v {
			key := &String{Value: k}
			hashKey := CreateHashKey(key)
			pairs[hashKey] = convertFromGoValue(val)
			keys = append(keys, key)
		}
		return &Hash{Pairs: pairs, Keys: keys}
	default:
		return newError("unsupported JSON type: %T", v)
	}
}

// convertToGoValue converts Rush Value to Go interface{} for json.Marshal
func convertToGoValue(value Value) (interface{}, error) {
	switch v := value.(type) {
	case *Null:
		return nil, nil
	case *Boolean:
		return v.Value, nil
	case *Integer:
		return v.Value, nil
	case *Float:
		return v.Value, nil
	case *String:
		return v.Value, nil
	case *Array:
		goArray := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			goValue, err := convertToGoValue(elem)
			if err != nil {
				return nil, err
			}
			goArray[i] = goValue
		}
		return goArray, nil
	case *Hash:
		goMap := make(map[string]interface{})
		for _, key := range v.Keys {
			hashKey := CreateHashKey(key)
			value := v.Pairs[hashKey]
			
			// Only string keys are supported in JSON
			keyStr, ok := key.(*String)
			if !ok {
				return nil, fmt.Errorf("JSON object keys must be strings, got %s", key.Type())
			}
			
			goValue, err := convertToGoValue(value)
			if err != nil {
				return nil, err
			}
			goMap[keyStr.Value] = goValue
		}
		return goMap, nil
	case *JSON:
		return convertToGoValue(v.Data)
	default:
		return nil, fmt.Errorf("unsupported value type for JSON: %s", v.Type())
	}
}