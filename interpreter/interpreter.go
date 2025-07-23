package interpreter

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"rush/ast"
)

// Eval evaluates an AST node and returns a value
func Eval(node ast.Node, env *Environment) Value {
	switch node := node.(type) {
	
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	
	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		
		// Check if this is an instance variable assignment (@variable)
		if len(node.Name.Value) > 0 && node.Name.Value[0] == '@' {
			// Look for the current object (self) in the environment
			if self, exists := env.Get("self"); exists {
				if obj, ok := self.(*Object); ok {
					// Remove @ prefix for storage in instance variables
					varName := node.Name.Value[1:]
					obj.InstanceVars[varName] = val
					return val
				}
			}
			return newError("instance variable %s used outside of object context", node.Name.Value)
		}
		
		env.Set(node.Name.Value, val)
		return val
	
	case *ast.IndexAssignmentStatement:
		return evalIndexAssignment(node, env)
	
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	
	// Expressions
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	
	case *ast.FloatLiteral:
		return &Float{Value: node.Value}
	
	case *ast.StringLiteral:
		return &String{Value: node.Value}
	
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanValue(node.Value)
	
	case *ast.Identifier:
		return evalIdentifier(node, env)
	
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		
		// Handle short-circuit evaluation for boolean operators
		if node.Operator == "&&" {
			if !IsTruthy(left) {
				return FALSE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			if !IsTruthy(right) {
				return FALSE
			}
			return TRUE
		}
		
		if node.Operator == "||" {
			if IsTruthy(left) {
				return TRUE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			if IsTruthy(right) {
				return TRUE
			}
			return FALSE
		}
		
		// For all other operators, evaluate right side normally
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}
	
	case *ast.CallExpression:
		// Check if this is a method call (object.method())
		if propAccess, ok := node.Function.(*ast.PropertyAccess); ok {
			// Evaluate the object
			object := Eval(propAccess.Object, env)
			if isError(object) {
				return object
			}
			
			// Check if object is an instance with methods
			if obj, ok := object.(*Object); ok {
				methodName := propAccess.Property.Value
				if method := resolveMethod(obj.Class, methodName); method != nil {
					// Set up method call environment with 'self' and parameters
					methodEnv := NewEnclosedEnvironment(method.Env)
					methodEnv.Set("self", obj)
					methodEnv.Set("__current_method__", &String{Value: methodName})
					
					// Evaluate arguments
					args := evalExpressions(node.Arguments, env)
					if len(args) == 1 && isError(args[0]) {
						return args[0]
					}
					
					// Check argument count
					if len(args) != len(method.Parameters) {
						return newError("wrong number of arguments: want=%d, got=%d", len(method.Parameters), len(args))
					}
					
					// Set up parameters in method environment
					for paramIdx, param := range method.Parameters {
						methodEnv.Set(param.Value, args[paramIdx])
					}
					
					// Evaluate method body with proper environment
					result := Eval(method.Body, methodEnv)
					return unwrapReturnValue(result)
				}
				return newError("undefined method %s for class %s", methodName, obj.Class.Name)
			}
		}
		
		// Regular function call
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		
		// Check if it's a hash method call
		if hashMethod, ok := function.(*HashMethod); ok {
			return applyHashMethod(hashMethod, args, env)
		}
		
		// Check if it's a string method call
		if stringMethod, ok := function.(*StringMethod); ok {
			return applyStringMethod(stringMethod, args, env)
		}
		
		// Check if it's an array method call
		if arrayMethod, ok := function.(*ArrayMethod); ok {
			return applyArrayMethod(arrayMethod, args, env)
		}
		
		// Check if it's a number method call
		if numberMethod, ok := function.(*NumberMethod); ok {
			return applyNumberMethod(numberMethod, args, env)
		}
		
		// Check if it's a file method call
		if fileMethod, ok := function.(*FileMethod); ok {
			return applyFileMethod(fileMethod, args, env)
		}
		
		// Check if it's a directory method call
		if dirMethod, ok := function.(*DirectoryMethod); ok {
			return applyDirectoryMethod(dirMethod, args, env)
		}
		
		// Check if it's a path method call
		if pathMethod, ok := function.(*PathMethod); ok {
			return applyPathMethod(pathMethod, args, env)
		}
		
		return applyFunction(function, args, node, env)
	
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	
	case *ast.BreakStatement:
		return &BreakValue{}
	
	case *ast.ContinueStatement:
		return &ContinueValue{}
	
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	
	case *ast.SwitchStatement:
		return evalSwitchStatement(node, env)
	
	case *ast.ForStatement:
		return evalForStatement(node, env)
	
	case *ast.ImportStatement:
		return evalImportStatement(node, env)
	
	case *ast.ExportStatement:
		return evalExportStatement(node, env)
	
	case *ast.PropertyAccess:
		return evalPropertyAccess(node, env)
	
	case *ast.ModuleAccess:
		return evalModuleAccess(node, env)
	
	case *ast.TryStatement:
		return evalTryStatement(node, env)
	
	case *ast.ThrowStatement:
		return evalThrowStatement(node, env)
	
	case *ast.ClassDeclaration:
		return evalClassDeclaration(node, env)
	
	case *ast.InstanceVariable:
		return evalInstanceVariable(node, env)
	
	case *ast.NewExpression:
		return evalNewExpression(node, env)
	
	case *ast.SuperExpression:
		return evalSuperExpression(node, env)
	
	default:
		return newError("unknown node type: %T", node)
	}
}

func evalProgram(stmts []ast.Statement, env *Environment) Value {
	var result Value
	
	for _, statement := range stmts {
		result = Eval(statement, env)
		
		if result != nil {
			if result.Type() == RETURN_VALUE {
				return result.(*ReturnValue).Value
			}
			if result.Type() == ERROR_VALUE || result.Type() == EXCEPTION_VALUE {
				return result
			}
		}
	}
	
	return result
}

func evalIdentifier(node *ast.Identifier, env *Environment) Value {
	val, ok := env.Get(node.Value)
	if !ok {
		return newErrorWithPosition(node.Token.Line, node.Token.Column, "identifier not found: %s", node.Value)
	}
	
	return val
}

func evalPrefixExpression(operator string, right Value) Value {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Value) Value {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right Value) Value {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: -right.Value}
	case *Float:
		return &Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right Value) Value {
	switch {
	case left.Type() == INTEGER_VALUE && right.Type() == INTEGER_VALUE:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == FLOAT_VALUE && right.Type() == FLOAT_VALUE:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == INTEGER_VALUE && right.Type() == FLOAT_VALUE:
		return evalMixedNumberInfixExpression(operator, left, right)
	case left.Type() == FLOAT_VALUE && right.Type() == INTEGER_VALUE:
		return evalMixedNumberInfixExpression(operator, left, right)
	case left.Type() == STRING_VALUE && right.Type() == STRING_VALUE:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == STRING_VALUE || right.Type() == STRING_VALUE:
		return evalStringCoercionInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanValue(left == right)
	case operator == "!=":
		return nativeBoolToBooleanValue(left != right)
	case operator == "&&":
		return evalBooleanInfixExpression(operator, left, right)
	case operator == "||":
		return evalBooleanInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Value) Value {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	
	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: float64(leftVal) / float64(rightVal)}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanValue(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanValue(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanValue(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanValue(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanValue(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanValue(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalFloatInfixExpression(operator string, left, right Value) Value {
	leftVal := left.(*Float).Value
	rightVal := right.(*Float).Value
	
	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanValue(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanValue(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanValue(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanValue(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanValue(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanValue(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalMixedNumberInfixExpression(operator string, left, right Value) Value {
	var leftVal, rightVal float64
	
	if left.Type() == INTEGER_VALUE {
		leftVal = float64(left.(*Integer).Value)
	} else {
		leftVal = left.(*Float).Value
	}
	
	if right.Type() == INTEGER_VALUE {
		rightVal = float64(right.(*Integer).Value)
	} else {
		rightVal = right.(*Float).Value
	}
	
	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanValue(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanValue(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanValue(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanValue(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanValue(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanValue(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalStringInfixExpression(operator string, left, right Value) Value {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	
	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanValue(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanValue(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanValue(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanValue(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanValue(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanValue(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalStringCoercionInfixExpression(operator string, left, right Value) Value {
	leftStr := valueToString(left)
	rightStr := valueToString(right)
	
	switch operator {
	case "+":
		return &String{Value: leftStr + rightStr}
	case "==":
		return nativeBoolToBooleanValue(leftStr == rightStr)
	case "!=":
		return nativeBoolToBooleanValue(leftStr != rightStr)
	case "<":
		return nativeBoolToBooleanValue(leftStr < rightStr)
	case ">":
		return nativeBoolToBooleanValue(leftStr > rightStr)
	case "<=":
		return nativeBoolToBooleanValue(leftStr <= rightStr)
	case ">=":
		return nativeBoolToBooleanValue(leftStr >= rightStr)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalExpressions(exps []ast.Expression, env *Environment) []Value {
	result := []Value{}
	
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Value{evaluated}
		}
		result = append(result, evaluated)
	}
	
	return result
}

func evalHashLiteral(node *ast.HashLiteral, env *Environment) Value {
	pairs := make(map[HashKey]Value)
	keys := []Value{}

	for _, pair := range node.Pairs {
		key := Eval(pair.Key, env)
		if isError(key) {
			return key
		}

		// Check if key is hashable (integer, string, boolean, float)
		if !isHashable(key) {
			return newError("unusable as hash key: %T", key)
		}

		value := Eval(pair.Value, env)
		if isError(value) {
			return value
		}

		hashKey := CreateHashKey(key)
		pairs[hashKey] = value
		keys = append(keys, key)
	}

	return &Hash{Pairs: pairs, Keys: keys}
}

func isHashable(value Value) bool {
	switch value.(type) {
	case *Integer, *String, *Boolean, *Float:
		return true
	default:
		return false
	}
}

func nativeBoolToBooleanValue(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Error handling
const ERROR_VALUE = "ERROR"

type Error struct {
	ErrorType string // e.g., "ValidationError", "RuntimeError"
	Message   string
	Stack     string
	Line      int
	Column    int
}

func (e *Error) Type() ValueType { return ERROR_VALUE }
func (e *Error) Inspect() string {
	errorType := e.ErrorType
	if errorType == "" {
		errorType = "Error"
	}
	if e.Line > 0 {
		return fmt.Sprintf("%s at line %d:%d: %s", errorType, e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("%s: %s", errorType, e.Message)
}

func newError(format string, a ...interface{}) *Error {
	return &Error{
		ErrorType: "RuntimeError",
		Message:   fmt.Sprintf(format, a...),
		Stack:     "", // Will be populated when thrown
	}
}

func newErrorWithPosition(line, column int, format string, a ...interface{}) *Error {
	return &Error{
		ErrorType: "RuntimeError",
		Message:   fmt.Sprintf(format, a...),
		Stack:     "", // Will be populated when thrown
		Line:      line,
		Column:    column,
	}
}

// newTypedError creates a new error with a specific type
func newTypedError(errorType, message string, line, column int) *Error {
	return &Error{
		ErrorType: errorType,
		Message:   message,
		Stack:     "", // Will be populated when thrown
		Line:      line,
		Column:    column,
	}
}

// Built-in error constructor functions
var ErrorConstructors = map[string]func(string, int, int) *Error{
	"Error":           func(msg string, line, col int) *Error { return newTypedError("Error", msg, line, col) },
	"ValidationError": func(msg string, line, col int) *Error { return newTypedError("ValidationError", msg, line, col) },
	"TypeError":       func(msg string, line, col int) *Error { return newTypedError("TypeError", msg, line, col) },
	"IndexError":      func(msg string, line, col int) *Error { return newTypedError("IndexError", msg, line, col) },
	"ArgumentError":   func(msg string, line, col int) *Error { return newTypedError("ArgumentError", msg, line, col) },
	"RuntimeError":    func(msg string, line, col int) *Error { return newTypedError("RuntimeError", msg, line, col) },
}

func isError(val Value) bool {
	if val != nil {
		return val.Type() == ERROR_VALUE || val.Type() == EXCEPTION_VALUE
	}
	return false
}

func evalIfExpression(ie *ast.IfExpression, env *Environment) Value {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if IsTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) Value {
	var result Value

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE || rt == BREAK_VALUE || rt == CONTINUE_VALUE {
				return result
			}
		}
	}

	return result
}

func evalBooleanInfixExpression(operator string, left, right Value) Value {
	switch operator {
	case "&&":
		// Short-circuit evaluation for AND
		if !IsTruthy(left) {
			return FALSE
		}
		if !IsTruthy(right) {
			return FALSE
		}
		return TRUE
	case "||":
		// Short-circuit evaluation for OR
		if IsTruthy(left) {
			return TRUE
		}
		if IsTruthy(right) {
			return TRUE
		}
		return FALSE
	default:
		return newError("unknown operator: %s", operator)
	}
}

func applyFunction(fn Value, args []Value, callNode *ast.CallExpression, env *Environment) Value {
	// Get function name for stack trace
	var functionName string
	if ident, ok := callNode.Function.(*ast.Identifier); ok {
		functionName = ident.Value
	} else {
		functionName = "<anonymous>"
	}
	
	switch fn := fn.(type) {
	case *BoundMethod:
		// Handle bound method calls
		if len(args) != len(fn.Method.Parameters) {
			return newError("wrong number of arguments: want=%d, got=%d", len(fn.Method.Parameters), len(args))
		}
		
		// Set up method call environment with 'self' and parameters
		methodEnv := NewEnclosedEnvironment(fn.Method.Env)
		methodEnv.Set("self", fn.Instance)
		
		// Inherit the call stack
		methodEnv.callStack = env.callStack
		
		// Push method call onto stack
		methodName := functionName
		if methodName == "<anonymous>" {
			methodName = "<method>"
		}
		env.PushCall(methodName, callNode.Token.Line, callNode.Token.Column)
		
		// Set up parameters in method environment
		for paramIdx, param := range fn.Method.Parameters {
			methodEnv.Set(param.Value, args[paramIdx])
		}
		
		// Evaluate method body with proper environment
		result := Eval(fn.Method.Body, methodEnv)
		
		// Pop method call from stack
		env.PopCall()
		
		return unwrapReturnValue(result)
	case *Function:
		if len(args) != len(fn.Parameters) {
			return newError("wrong number of arguments: want=%d, got=%d", len(fn.Parameters), len(args))
		}
		
		// Push function call onto stack
		env.PushCall(functionName, callNode.Token.Line, callNode.Token.Column)
		
		extendedEnv := extendFunctionEnv(fn, args)
		// Inherit the call stack
		extendedEnv.callStack = env.callStack
		
		evaluated := Eval(fn.Body, extendedEnv)
		
		// Pop function call from stack
		env.PopCall()
		
		return unwrapReturnValue(evaluated)
	case *BuiltinFunction:
		// Don't track built-in function calls in stack trace
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

func applyHashMethod(hashMethod *HashMethod, args []Value, env *Environment) Value {
	switch hashMethod.Method {
	case "has_key?":
		if len(args) != 1 {
			return newError("wrong number of arguments for has_key?: want=1, got=%d", len(args))
		}
		key := args[0]
		hashKey := CreateHashKey(key)
		_, exists := hashMethod.Hash.Pairs[hashKey]
		return &Boolean{Value: exists}
		
	case "has_value?":
		if len(args) != 1 {
			return newError("wrong number of arguments for has_value?: want=1, got=%d", len(args))
		}
		searchValue := args[0]
		for _, value := range hashMethod.Hash.Pairs {
			if compareValues(value, searchValue) {
				return TRUE
			}
		}
		return FALSE
		
	case "get":
		if len(args) < 1 || len(args) > 2 {
			return newError("wrong number of arguments for get: want=1 or 2, got=%d", len(args))
		}
		key := args[0]
		hashKey := CreateHashKey(key)
		if value, exists := hashMethod.Hash.Pairs[hashKey]; exists {
			return value
		}
		if len(args) == 2 {
			return args[1] // default value
		}
		return NULL
		
	case "set":
		if len(args) != 2 {
			return newError("wrong number of arguments for set: want=2, got=%d", len(args))
		}
		return hashSet(hashMethod.Hash, args[0], args[1])
		
	case "delete":
		if len(args) != 1 {
			return newError("wrong number of arguments for delete: want=1, got=%d", len(args))
		}
		return hashDelete(hashMethod.Hash, args[0])
		
	case "merge":
		if len(args) != 1 {
			return newError("wrong number of arguments for merge: want=1, got=%d", len(args))
		}
		otherHash, ok := args[0].(*Hash)
		if !ok {
			return newError("argument to merge must be HASH, got %s", args[0].Type())
		}
		return hashMerge(hashMethod.Hash, otherHash)
		
	case "filter":
		if len(args) != 1 {
			return newError("wrong number of arguments for filter: want=1, got=%d", len(args))
		}
		predicate, ok := args[0].(*Function)
		if !ok {
			return newError("argument to filter must be FUNCTION, got %s", args[0].Type())
		}
		return hashFilter(hashMethod.Hash, predicate, env)
		
	case "map_values":
		if len(args) != 1 {
			return newError("wrong number of arguments for map_values: want=1, got=%d", len(args))
		}
		transform, ok := args[0].(*Function)
		if !ok {
			return newError("argument to map_values must be FUNCTION, got %s", args[0].Type())
		}
		return hashMapValues(hashMethod.Hash, transform, env)
		
	case "each":
		if len(args) != 1 {
			return newError("wrong number of arguments for each: want=1, got=%d", len(args))
		}
		callback, ok := args[0].(*Function)
		if !ok {
			return newError("argument to each must be FUNCTION, got %s", args[0].Type())
		}
		hashEach(hashMethod.Hash, callback, env)
		return hashMethod.Hash // return original hash
		
	case "select_keys":
		if len(args) != 1 {
			return newError("wrong number of arguments for select_keys: want=1, got=%d", len(args))
		}
		keyArray, ok := args[0].(*Array)
		if !ok {
			return newError("argument to select_keys must be ARRAY, got %s", args[0].Type())
		}
		return hashSelectKeys(hashMethod.Hash, keyArray)
		
	case "reject_keys":
		if len(args) != 1 {
			return newError("wrong number of arguments for reject_keys: want=1, got=%d", len(args))
		}
		keyArray, ok := args[0].(*Array)
		if !ok {
			return newError("argument to reject_keys must be ARRAY, got %s", args[0].Type())
		}
		return hashRejectKeys(hashMethod.Hash, keyArray)
		
	case "invert":
		if len(args) != 0 {
			return newError("wrong number of arguments for invert: want=0, got=%d", len(args))
		}
		return hashInvert(hashMethod.Hash)
		
	case "to_array":
		if len(args) != 0 {
			return newError("wrong number of arguments for to_array: want=0, got=%d", len(args))
		}
		return hashToArray(hashMethod.Hash)
		
	default:
		return newError("unknown hash method: %s", hashMethod.Method)
	}
}

func applyStringMethod(stringMethod *StringMethod, args []Value, env *Environment) Value {
	str := stringMethod.String.Value
	
	switch stringMethod.Method {
	case "trim":
		if len(args) != 0 {
			return newError("wrong number of arguments for trim: want=0, got=%d", len(args))
		}
		return &String{Value: strings.TrimSpace(str)}
		
	case "ltrim":
		if len(args) != 0 {
			return newError("wrong number of arguments for ltrim: want=0, got=%d", len(args))
		}
		return &String{Value: strings.TrimLeftFunc(str, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})}
		
	case "rtrim":
		if len(args) != 0 {
			return newError("wrong number of arguments for rtrim: want=0, got=%d", len(args))
		}
		return &String{Value: strings.TrimRightFunc(str, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})}
		
	case "upper":
		if len(args) != 0 {
			return newError("wrong number of arguments for upper: want=0, got=%d", len(args))
		}
		return &String{Value: strings.ToUpper(str)}
		
	case "lower":
		if len(args) != 0 {
			return newError("wrong number of arguments for lower: want=0, got=%d", len(args))
		}
		return &String{Value: strings.ToLower(str)}
		
	case "contains?":
		if len(args) != 1 {
			return newError("wrong number of arguments for contains?: want=1, got=%d", len(args))
		}
		searchStr, ok := args[0].(*String)
		if !ok {
			return newError("argument to contains? must be STRING, got %s", args[0].Type())
		}
		return &Boolean{Value: strings.Contains(str, searchStr.Value)}
		
	case "replace":
		if len(args) != 2 {
			return newError("wrong number of arguments for replace: want=2, got=%d", len(args))
		}
		old, ok1 := args[0].(*String)
		new, ok2 := args[1].(*String)
		if !ok1 || !ok2 {
			return newError("arguments to replace must be STRING, got %s, %s", args[0].Type(), args[1].Type())
		}
		return &String{Value: strings.ReplaceAll(str, old.Value, new.Value)}
		
	case "starts_with?":
		if len(args) != 1 {
			return newError("wrong number of arguments for starts_with?: want=1, got=%d", len(args))
		}
		prefixStr, ok := args[0].(*String)
		if !ok {
			return newError("argument to starts_with? must be STRING, got %s", args[0].Type())
		}
		return &Boolean{Value: strings.HasPrefix(str, prefixStr.Value)}
		
	case "ends_with?":
		if len(args) != 1 {
			return newError("wrong number of arguments for ends_with?: want=1, got=%d", len(args))
		}
		suffixStr, ok := args[0].(*String)
		if !ok {
			return newError("argument to ends_with? must be STRING, got %s", args[0].Type())
		}
		return &Boolean{Value: strings.HasSuffix(str, suffixStr.Value)}
		
	case "substr":
		if len(args) != 2 {
			return newError("wrong number of arguments for substr: want=2, got=%d", len(args))
		}
		start, ok1 := args[0].(*Integer)
		length, ok2 := args[1].(*Integer)
		if !ok1 || !ok2 {
			return newError("arguments to substr must be INTEGER, got %s, %s", args[0].Type(), args[1].Type())
		}
		
		startIdx := int(start.Value)
		lengthVal := int(length.Value)
		
		if startIdx < 0 || startIdx >= len(str) {
			return &String{Value: ""}
		}
		
		endIdx := startIdx + lengthVal
		if endIdx > len(str) {
			endIdx = len(str)
		}
		
		if lengthVal <= 0 {
			return &String{Value: ""}
		}
		
		return &String{Value: str[startIdx:endIdx]}
		
	case "split":
		if len(args) != 1 {
			return newError("wrong number of arguments for split: want=1, got=%d", len(args))
		}
		delimiter, ok := args[0].(*String)
		if !ok {
			return newError("argument to split must be STRING, got %s", args[0].Type())
		}
		
		parts := strings.Split(str, delimiter.Value)
		elements := make([]Value, len(parts))
		for i, part := range parts {
			elements[i] = &String{Value: part}
		}
		return &Array{Elements: elements}
		
	case "join":
		if len(args) != 1 {
			return newError("wrong number of arguments for join: want=1, got=%d", len(args))
		}
		arr, ok := args[0].(*Array)
		if !ok {
			return newError("argument to join must be ARRAY, got %s", args[0].Type())
		}
		
		strs := make([]string, len(arr.Elements))
		for i, elem := range arr.Elements {
			strs[i] = valueToString(elem)
		}
		return &String{Value: strings.Join(strs, str)}
		
	default:
		return newError("unknown string method: %s", stringMethod.Method)
	}
}

func applyArrayMethod(arrayMethod *ArrayMethod, args []Value, env *Environment) Value {
	arr := arrayMethod.Array
	
	switch arrayMethod.Method {
	case "map":
		if len(args) != 1 {
			return newError("wrong number of arguments for map: want=1, got=%d", len(args))
		}
		mapFunc, ok := args[0].(*Function)
		if !ok {
			return newError("argument to map must be FUNCTION, got %s", args[0].Type())
		}
		
		result := []Value{}
		for _, elem := range arr.Elements {
			extendedEnv := extendFunctionEnv(mapFunc, []Value{elem})
			mapped := Eval(mapFunc.Body, extendedEnv)
			if isError(mapped) {
				return mapped
			}
			result = append(result, unwrapReturnValue(mapped))
		}
		return &Array{Elements: result}
		
	case "filter":
		if len(args) != 1 {
			return newError("wrong number of arguments for filter: want=1, got=%d", len(args))
		}
		filterFunc, ok := args[0].(*Function)
		if !ok {
			return newError("argument to filter must be FUNCTION, got %s", args[0].Type())
		}
		
		result := []Value{}
		for _, elem := range arr.Elements {
			extendedEnv := extendFunctionEnv(filterFunc, []Value{elem})
			filtered := Eval(filterFunc.Body, extendedEnv)
			if isError(filtered) {
				return filtered
			}
			if IsTruthy(unwrapReturnValue(filtered)) {
				result = append(result, elem)
			}
		}
		return &Array{Elements: result}
		
	case "reduce":
		if len(args) != 2 {
			return newError("wrong number of arguments for reduce: want=2, got=%d", len(args))
		}
		reduceFunc, ok := args[0].(*Function)
		if !ok {
			return newError("first argument to reduce must be FUNCTION, got %s", args[0].Type())
		}
		
		result := args[1] // initial value
		for _, elem := range arr.Elements {
			extendedEnv := extendFunctionEnv(reduceFunc, []Value{result, elem})
			result = Eval(reduceFunc.Body, extendedEnv)
			if isError(result) {
				return result
			}
			result = unwrapReturnValue(result)
		}
		return result
		
	case "find":
		if len(args) != 1 {
			return newError("wrong number of arguments for find: want=1, got=%d", len(args))
		}
		findFunc, ok := args[0].(*Function)
		if !ok {
			return newError("argument to find must be FUNCTION, got %s", args[0].Type())
		}
		
		for _, elem := range arr.Elements {
			extendedEnv := extendFunctionEnv(findFunc, []Value{elem})
			found := Eval(findFunc.Body, extendedEnv)
			if isError(found) {
				return found
			}
			if IsTruthy(unwrapReturnValue(found)) {
				return elem
			}
		}
		return NULL
		
	case "index_of":
		if len(args) != 1 {
			return newError("wrong number of arguments for index_of: want=1, got=%d", len(args))
		}
		searchElement := args[0]
		
		for i, elem := range arr.Elements {
			if compareValues(elem, searchElement) {
				return &Integer{Value: int64(i)}
			}
		}
		return &Integer{Value: -1}
		
	case "includes?":
		if len(args) != 1 {
			return newError("wrong number of arguments for includes?: want=1, got=%d", len(args))
		}
		searchElement := args[0]
		
		for _, elem := range arr.Elements {
			if compareValues(elem, searchElement) {
				return TRUE
			}
		}
		return FALSE
		
	case "reverse":
		if len(args) != 0 {
			return newError("wrong number of arguments for reverse: want=0, got=%d", len(args))
		}
		
		result := make([]Value, len(arr.Elements))
		for i, elem := range arr.Elements {
			result[len(arr.Elements)-1-i] = elem
		}
		return &Array{Elements: result}
		
	case "sort":
		if len(args) != 0 {
			return newError("wrong number of arguments for sort: want=0, got=%d", len(args))
		}
		
		// Create a copy of the array
		result := make([]Value, len(arr.Elements))
		copy(result, arr.Elements)
		
		// Simple bubble sort
		n := len(result)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if compareForSort(result[j], result[j+1]) > 0 {
					result[j], result[j+1] = result[j+1], result[j]
				}
			}
		}
		return &Array{Elements: result}
		
	case "push":
		if len(args) != 1 {
			return newError("wrong number of arguments for push: want=1, got=%d", len(args))
		}
		
		newElements := make([]Value, len(arr.Elements)+1)
		copy(newElements, arr.Elements)
		newElements[len(arr.Elements)] = args[0]
		return &Array{Elements: newElements}
		
	case "pop":
		if len(args) != 0 {
			return newError("wrong number of arguments for pop: want=0, got=%d", len(args))
		}
		
		if len(arr.Elements) == 0 {
			return newError("cannot pop from empty array")
		}
		
		return arr.Elements[len(arr.Elements)-1]
		
	case "slice":
		if len(args) != 2 {
			return newError("wrong number of arguments for slice: want=2, got=%d", len(args))
		}
		
		start, ok1 := args[0].(*Integer)
		end, ok2 := args[1].(*Integer)
		if !ok1 || !ok2 {
			return newError("arguments to slice must be INTEGER, got %s, %s", args[0].Type(), args[1].Type())
		}
		
		startIdx := int(start.Value)
		endIdx := int(end.Value)
		
		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > len(arr.Elements) {
			endIdx = len(arr.Elements)
		}
		if startIdx >= endIdx {
			return &Array{Elements: []Value{}}
		}
		
		result := make([]Value, endIdx-startIdx)
		copy(result, arr.Elements[startIdx:endIdx])
		return &Array{Elements: result}
		
	default:
		return newError("unknown array method: %s", arrayMethod.Method)
	}
}

// compareForSort compares two values for sorting purposes
func compareForSort(a, b Value) int {
	switch aVal := a.(type) {
	case *Integer:
		if bVal, ok := b.(*Integer); ok {
			if aVal.Value < bVal.Value {
				return -1
			} else if aVal.Value > bVal.Value {
				return 1
			}
			return 0
		}
	case *Float:
		if bVal, ok := b.(*Float); ok {
			if aVal.Value < bVal.Value {
				return -1
			} else if aVal.Value > bVal.Value {
				return 1
			}
			return 0
		}
	case *String:
		if bVal, ok := b.(*String); ok {
			if aVal.Value < bVal.Value {
				return -1
			} else if aVal.Value > bVal.Value {
				return 1
			}
			return 0
		}
	}
	// For mixed types or unsupported types, convert to string and compare
	aStr := valueToString(a)
	bStr := valueToString(b)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

func applyNumberMethod(numberMethod *NumberMethod, args []Value, env *Environment) Value {
	num := numberMethod.Number
	
	switch numberMethod.Method {
	case "abs":
		if len(args) != 0 {
			return newError("wrong number of arguments for abs: want=0, got=%d", len(args))
		}
		
		switch n := num.(type) {
		case *Integer:
			if n.Value < 0 {
				return &Integer{Value: -n.Value}
			}
			return n
		case *Float:
			if n.Value < 0 {
				return &Float{Value: -n.Value}
			}
			return n
		default:
			return newError("abs not supported for type %s", num.Type())
		}
		
	case "floor":
		if len(args) != 0 {
			return newError("wrong number of arguments for floor: want=0, got=%d", len(args))
		}
		
		switch n := num.(type) {
		case *Integer:
			return n // integers are already floored
		case *Float:
			return &Float{Value: math.Floor(n.Value)}
		default:
			return newError("floor not supported for type %s", num.Type())
		}
		
	case "ceil":
		if len(args) != 0 {
			return newError("wrong number of arguments for ceil: want=0, got=%d", len(args))
		}
		
		switch n := num.(type) {
		case *Integer:
			return n // integers are already ceiled
		case *Float:
			return &Float{Value: math.Ceil(n.Value)}
		default:
			return newError("ceil not supported for type %s", num.Type())
		}
		
	case "round":
		if len(args) != 0 {
			return newError("wrong number of arguments for round: want=0, got=%d", len(args))
		}
		
		switch n := num.(type) {
		case *Integer:
			return n // integers are already rounded
		case *Float:
			return &Float{Value: math.Round(n.Value)}
		default:
			return newError("round not supported for type %s", num.Type())
		}
		
	case "sqrt":
		if len(args) != 0 {
			return newError("wrong number of arguments for sqrt: want=0, got=%d", len(args))
		}
		
		var value float64
		switch n := num.(type) {
		case *Integer:
			if n.Value < 0 {
				return newError("sqrt of negative number")
			}
			value = float64(n.Value)
		case *Float:
			if n.Value < 0 {
				return newError("sqrt of negative number")
			}
			value = n.Value
		default:
			return newError("sqrt not supported for type %s", num.Type())
		}
		
		return &Float{Value: math.Sqrt(value)}
		
	case "pow":
		if len(args) != 1 {
			return newError("wrong number of arguments for pow: want=1, got=%d", len(args))
		}
		
		var base, exponent float64
		
		// Get base value
		switch n := num.(type) {
		case *Integer:
			base = float64(n.Value)
		case *Float:
			base = n.Value
		default:
			return newError("pow not supported for type %s", num.Type())
		}
		
		// Get exponent value
		switch exp := args[0].(type) {
		case *Integer:
			exponent = float64(exp.Value)
		case *Float:
			exponent = exp.Value
		default:
			return newError("exponent to pow must be number, got %s", args[0].Type())
		}
		
		result := math.Pow(base, exponent)
		
		// If both inputs were integers and result is a whole number, return integer
		if _, baseIsInt := num.(*Integer); baseIsInt {
			if _, expIsInt := args[0].(*Integer); expIsInt {
				if result == float64(int64(result)) {
					return &Integer{Value: int64(result)}
				}
			}
		}
		
		return &Float{Value: result}
		
	default:
		return newError("unknown number method: %s", numberMethod.Method)
	}
}

func extendFunctionEnv(fn *Function, args []Value) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(val Value) Value {
	if returnValue, ok := val.(*ReturnValue); ok {
		return returnValue.Value
	}
	// Convert break/continue that escape function boundaries to errors
	if _, ok := val.(*BreakValue); ok {
		return newError("break statement not in loop")
	}
	if _, ok := val.(*ContinueValue); ok {
		return newError("continue statement not in loop")
	}
	// Don't unwrap exceptions - they should propagate as-is
	return val
}

// valueToString converts any value to its string representation for coercion
func valueToString(val Value) string {
	switch val.Type() {
	case STRING_VALUE:
		return val.(*String).Value
	case INTEGER_VALUE:
		return fmt.Sprintf("%d", val.(*Integer).Value)
	case FLOAT_VALUE:
		return fmt.Sprintf("%g", val.(*Float).Value)
	case BOOLEAN_VALUE:
		if val.(*Boolean).Value {
			return "true"
		}
		return "false"
	case NULL_VALUE:
		return "null"
	default:
		return val.Inspect() // fallback to existing representation
	}
}

func evalIndexExpression(left, index Value) Value {
	switch {
	case left.Type() == ARRAY_VALUE && index.Type() == INTEGER_VALUE:
		return evalArrayIndexExpression(left, index)
	case left.Type() == ARRAY_VALUE && index.Type() == FLOAT_VALUE:
		// Allow float indices that represent whole numbers
		floatIdx := index.(*Float).Value
		if floatIdx == float64(int64(floatIdx)) {
			intIdx := &Integer{Value: int64(floatIdx)}
			return evalArrayIndexExpression(left, intIdx)
		}
		return newError("array index must be a whole number, got: %g", floatIdx)
	case left.Type() == STRING_VALUE && index.Type() == INTEGER_VALUE:
		return evalStringIndexExpression(left, index)
	case left.Type() == STRING_VALUE && index.Type() == FLOAT_VALUE:
		// Allow float indices that represent whole numbers
		floatIdx := index.(*Float).Value
		if floatIdx == float64(int64(floatIdx)) {
			intIdx := &Integer{Value: int64(floatIdx)}
			return evalStringIndexExpression(left, intIdx)
		}
		return newError("string index must be a whole number, got: %g", floatIdx)
	case left.Type() == HASH_VALUE:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index Value) Value {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		errorObj := newTypedError("IndexError", fmt.Sprintf("array index %d out of range [0:%d]", idx, max+1), 0, 0)
		return NewException(errorObj)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index Value) Value {
	stringObject := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		errorObj := newTypedError("IndexError", fmt.Sprintf("string index %d out of range [0:%d]", idx, max+1), 0, 0)
		return NewException(errorObj)
	}

	return &String{Value: string(stringObject.Value[idx])}
}

func evalHashIndexExpression(hash, index Value) Value {
	hashObject := hash.(*Hash)
	
	if !isHashable(index) {
		return newError("unusable as hash key: %T", index)
	}

	hashKey := CreateHashKey(index)
	value, exists := hashObject.Pairs[hashKey]
	if !exists {
		return NULL
	}

	return value
}

func evalWhileStatement(ws *ast.WhileStatement, env *Environment) Value {
	var result Value = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !IsTruthy(condition) {
			break
		}

		result = Eval(ws.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE {
				return result
			}
			if rt == BREAK_VALUE {
				result = NULL // Don't return the BreakValue
				break // Exit the while loop
			}
			if rt == CONTINUE_VALUE {
				continue // Skip to next iteration
			}
		}
	}

	return result
}

func evalSwitchStatement(ss *ast.SwitchStatement, env *Environment) Value {
	// Evaluate the switch value
	switchValue := Eval(ss.Value, env)
	if isError(switchValue) {
		return switchValue
	}

	// Check each case clause
	for _, caseClause := range ss.Cases {
		for _, caseValue := range caseClause.Values {
			caseVal := Eval(caseValue, env)
			if isError(caseVal) {
				return caseVal
			}
			
			// Compare values for equality
			if compareValues(switchValue, caseVal) {
				result := Eval(caseClause.Body, env)
				if result != nil {
					rt := result.Type()
					if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE {
						return result
					}
					if rt == BREAK_VALUE {
						return NULL // Go-style automatic break
					}
					if rt == CONTINUE_VALUE {
						return result // Pass continue through
					}
				}
				return result // Go-style automatic break - execute only this case
			}
		}
	}

	// If no case matched, try default clause
	if ss.Default != nil {
		result := Eval(ss.Default.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE {
				return result
			}
			if rt == BREAK_VALUE {
				return NULL
			}
			if rt == CONTINUE_VALUE {
				return result
			}
		}
		return result
	}

	return NULL // No match and no default
}

// Helper function to compare two values for equality
func compareValues(left, right Value) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch left := left.(type) {
	case *Integer:
		if right, ok := right.(*Integer); ok {
			return left.Value == right.Value
		}
	case *Float:
		if right, ok := right.(*Float); ok {
			return left.Value == right.Value
		}
	case *String:
		if right, ok := right.(*String); ok {
			return left.Value == right.Value
		}
	case *Boolean:
		if right, ok := right.(*Boolean); ok {
			return left.Value == right.Value
		}
	}
	
	return false
}

func evalForStatement(fs *ast.ForStatement, env *Environment) Value {
	var result Value = NULL
	
	// Don't create a separate scope - use the current environment
	// This way variables are accessible and modifiable

	// Execute init statement if present
	if fs.Init != nil {
		initResult := Eval(fs.Init, env)
		if isError(initResult) {
			return initResult
		}
	}

	for {
		// Check condition (if no condition, loop forever until break/return)
		if fs.Condition != nil {
			condition := Eval(fs.Condition, env)
			if isError(condition) {
				return condition
			}
			if !IsTruthy(condition) {
				break
			}
		}

		// Execute body
		result = Eval(fs.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE {
				return result
			}
			if rt == BREAK_VALUE {
				result = NULL // Don't return the BreakValue
				break // Exit the for loop
			}
			if rt == CONTINUE_VALUE {
				// Execute update statement before continuing
				if fs.Update != nil {
					updateResult := Eval(fs.Update, env)
					if isError(updateResult) {
						return updateResult
					}
				}
				continue // Skip to next iteration
			}
		}

		// Execute update statement if present
		if fs.Update != nil {
			updateResult := Eval(fs.Update, env)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

// evalImportStatement handles import statements
func evalImportStatement(node *ast.ImportStatement, env *Environment) Value {
	// Get the module path
	modulePath := node.Module.Value
	
	// Load the module
	moduleResolver := env.GetModuleResolver()
	module, err := moduleResolver.LoadModule(modulePath, env.GetCurrentDir())
	if err != nil {
		return newError("failed to import module %s: %s", modulePath, err.Error())
	}
	
	// Execute the module to populate its exports
	if module.Exports == nil || len(module.Exports) == 0 {
		// Create a new environment for the module
		moduleEnv := NewModuleEnvironment(env)
		
		// Set the current directory to the module's directory for nested imports
		moduleDir := filepath.Dir(module.Path)
		moduleEnv.SetCurrentDir(moduleDir)
		
		// Add native functions for standard library modules
		if isStandardLibraryModule(modulePath) {
			addNativeStandardLibraryFunctions(moduleEnv, modulePath)
		}
		
		// Execute the module
		result := Eval(module.AST, moduleEnv)
		if isError(result) {
			return newError("error executing module %s: %s", modulePath, result.Inspect())
		}
		
		// Extract exports from the module environment
		module.Exports = make(map[string]interface{})
		for name, value := range moduleEnv.GetExports() {
			module.Exports[name] = value
		}
	}
	
	// Import the specified items into the current environment
	for _, item := range node.Items {
		if value, exists := module.Exports[item.Name.Value]; exists {
			// Convert interface{} back to Value
			if val, ok := value.(Value); ok {
				// Use alias if provided, otherwise use original name
				importName := item.Name.Value
				if item.Alias != nil {
					importName = item.Alias.Value
				}
				env.Set(importName, val)
			} else {
				return newError("invalid export type for %s", item.Name.Value)
			}
		} else {
			return newError("module %s does not export %s", modulePath, item.Name.Value)
		}
	}
	
	return NULL
}

// evalExportStatement handles export statements
func evalExportStatement(node *ast.ExportStatement, env *Environment) Value {
	var value Value
	
	if node.Value != nil {
		// Export an assignment: export name = value
		value = Eval(node.Value, env)
		if isError(value) {
			return value
		}
		
		// Set the value in the local environment
		env.Set(node.Name.Value, value)
	} else {
		// Export a previously defined value: export name
		var exists bool
		value, exists = env.Get(node.Name.Value)
		if !exists {
			return newError("cannot export undefined variable: %s", node.Name.Value)
		}
	}
	
	// Add the export to the environment's export list
	env.AddExport(node.Name.Value, value)
	
	return value
}

// evalModuleAccess handles module member access like module.function
func evalModuleAccess(node *ast.ModuleAccess, env *Environment) Value {
	// For now, treat this as a simple property access
	// This is a placeholder implementation
	// In a full implementation, we'd need to track loaded modules
	
	return newError("module access not yet fully implemented: %s.%s", 
		node.Module.Value, node.Member.Value)
}

// evalPropertyAccess handles property access like "object.property"
func evalPropertyAccess(node *ast.PropertyAccess, env *Environment) Value {
	// Check for nil property
	if node.Property == nil {
		return newError("property name is missing in property access")
	}
	
	// Evaluate the object expression
	object := Eval(node.Object, env)
	
	// Check if it's an error object and handle property access
	if errorObj, ok := object.(*Error); ok {
		return getErrorProperty(errorObj, node.Property.Value)
	}
	
	// Check if it's a runtime error (Exception)
	if isError(object) {
		return object
	}
	
	// Check if it's an object instance and return bound method
	if obj, ok := object.(*Object); ok {
		methodName := node.Property.Value
		if method := resolveMethod(obj.Class, methodName); method != nil {
			// Return a bound method that includes 'self'
			return &BoundMethod{
				Method:   method,
				Instance: obj,
			}
		}
		return newError("undefined method %s for class %s", methodName, obj.Class.Name)
	}
	
	// Check if it's a hash and handle property access
	if hash, ok := object.(*Hash); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "keys":
			return &Array{Elements: hash.Keys}
		case "values":
			values := make([]Value, 0, len(hash.Keys))
			for _, key := range hash.Keys {
				hashKey := CreateHashKey(key)
				values = append(values, hash.Pairs[hashKey])
			}
			return &Array{Elements: values}
		case "length", "size":
			return &Integer{Value: int64(len(hash.Keys))}
		case "empty":
			return &Boolean{Value: len(hash.Keys) == 0}
		
		// Methods (with parameters) - return bound methods
		case "has_key?", "has_value?", "get", "set", "delete", "merge", 
		     "filter", "map_values", "each", "select_keys", "reject_keys",
		     "invert", "to_array":
			return &HashMethod{Hash: hash, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for hash", node.Property.Value)
		}
	}
	
	// Check if it's a string and handle property access
	if str, ok := object.(*String); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "length":
			return &Integer{Value: int64(len(str.Value))}
		case "empty":
			return &Boolean{Value: len(str.Value) == 0}
		
		// Methods (with parameters) - return bound methods
		case "trim", "ltrim", "rtrim", "upper", "lower", "contains?", "replace",
		     "starts_with?", "ends_with?", "substr", "split", "join":
			return &StringMethod{String: str, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for string", node.Property.Value)
		}
	}
	
	// Check if it's an array and handle property access
	if arr, ok := object.(*Array); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "length":
			return &Integer{Value: int64(len(arr.Elements))}
		case "empty":
			return &Boolean{Value: len(arr.Elements) == 0}
		
		// Methods (with parameters) - return bound methods
		case "map", "filter", "reduce", "find", "index_of", "includes?", "reverse", 
		     "sort", "push", "pop", "slice":
			return &ArrayMethod{Array: arr, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for array", node.Property.Value)
		}
	}
	
	// Check if it's a number (integer or float) and handle property access
	if num, ok := object.(*Integer); ok {
		switch node.Property.Value {
		// Methods (with parameters) - return bound methods
		case "abs", "floor", "ceil", "round", "sqrt", "pow":
			return &NumberMethod{Number: num, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for integer", node.Property.Value)
		}
	}
	
	if num, ok := object.(*Float); ok {
		switch node.Property.Value {
		// Methods (with parameters) - return bound methods  
		case "abs", "floor", "ceil", "round", "sqrt", "pow":
			return &NumberMethod{Number: num, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for float", node.Property.Value)
		}
	}
	
	// Check if it's a file and handle property access
	if file, ok := object.(*File); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "path":
			return &String{Value: file.Path}
		case "is_open":
			return &Boolean{Value: file.IsOpen}
		
		// Methods (with parameters) - return bound methods
		case "open", "read", "write", "close", "exists?", "size", "delete":
			return &FileMethod{File: file, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for file", node.Property.Value)
		}
	}
	
	// Check if it's a directory and handle property access
	if dir, ok := object.(*Directory); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "path":
			return &String{Value: dir.Path}
		
		// Methods (with parameters) - return bound methods
		case "create", "list", "delete", "exists?":
			return &DirectoryMethod{Directory: dir, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for directory", node.Property.Value)
		}
	}
	
	// Check if it's a path and handle property access
	if path, ok := object.(*Path); ok {
		switch node.Property.Value {
		// Simple properties (no parameters)
		case "value":
			return &String{Value: path.Value}
		
		// Methods (with parameters) - return bound methods
		case "join", "basename", "dirname", "absolute", "clean":
			return &PathMethod{Path: path, Method: node.Property.Value}
		
		default:
			return newError("unknown property %s for path", node.Property.Value)
		}
	}
	
	// For other objects, try module access (backward compatibility)
	if ident, ok := node.Object.(*ast.Identifier); ok {
		// This looks like module.member access
		return newError("module access not yet fully implemented: %s.%s", 
			ident.Value, node.Property.Value)
	}
	
	return newError("property access not supported for type %s", object.Type())
}

// getErrorProperty returns the specified property of an error object
func getErrorProperty(errorObj *Error, propertyName string) Value {
	switch propertyName {
	case "type":
		return &String{Value: errorObj.ErrorType}
	case "message":
		return &String{Value: errorObj.Message}
	case "stack":
		return &String{Value: errorObj.Stack}
	case "line":
		return &Integer{Value: int64(errorObj.Line)}
	case "column":
		return &Integer{Value: int64(errorObj.Column)}
	default:
		return newError("error object has no property '%s'", propertyName)
	}
}

// evalThrowStatement handles throw statements
func evalThrowStatement(node *ast.ThrowStatement, env *Environment) Value {
	// Evaluate the expression being thrown
	value := Eval(node.Expression, env)
	
	// Check if evaluation resulted in a runtime error (not an Error object value)
	if value.Type() == ERROR_VALUE {
		errorObj, ok := value.(*Error)
		if ok && (errorObj.ErrorType == "RuntimeError" || errorObj.ErrorType == "") {
			// This could be either a runtime error or an Error object value
			// If the expression was a function call to an error constructor, treat it as a value
			if call, ok := node.Expression.(*ast.CallExpression); ok {
				if ident, ok := call.Function.(*ast.Identifier); ok {
					// Check if it's one of our error constructors
					if _, exists := builtins[ident.Value]; exists {
						// This is an Error object value from a constructor, wrap it in Exception
						// Populate stack trace
						errorObj.Stack = env.GetStackTrace()
						return NewException(errorObj)
					}
				}
			}
			// Otherwise, it's a runtime error that should be propagated
			return value
		}
		// For other error types, wrap in Exception
		errorObj.Stack = env.GetStackTrace()
		return NewException(value.(*Error))
	}
	
	// If not an error object, throw the value as-is wrapped in a generic error
	errorObj := newTypedError("Error", value.Inspect(), node.Token.Line, node.Token.Column)
	errorObj.Stack = env.GetStackTrace()
	return NewException(errorObj)
}

// evalTryStatement handles try-catch-finally blocks
func evalTryStatement(node *ast.TryStatement, env *Environment) Value {
	var result Value = NULL
	var exception *Exception
	
	// Execute the try block
	tryResult := evalBlockStatement(node.TryBlock, env)
	
	// Check if an exception occurred
	if ex, ok := tryResult.(*Exception); ok {
		exception = ex
	} else {
		result = tryResult
	}
	
	// If there's an exception, try to handle it with catch clauses
	if exception != nil {
		handled := false
		
		for _, catchClause := range node.CatchClauses {
			// Check if this catch clause matches the error type
			if shouldCatchError(exception.Error, catchClause) {
				// Create a new environment for the catch block to shadow variables
				catchEnv := NewEnclosedEnvironment(env)
				
				// Bind the error variable in the catch environment (force local shadowing)
				catchEnv.SetLocal(catchClause.ErrorVar.Value, exception.Error)
				
				// Execute the catch block in the new environment
				catchResult := evalBlockStatement(catchClause.Body, catchEnv)
				
				// If catch block throws another exception, propagate it
				if ex, ok := catchResult.(*Exception); ok {
					exception = ex
					result = exception
				} else {
					// Exception was handled
					exception = nil
					result = catchResult
					handled = true
				}
				break
			}
		}
		
		// If no catch clause handled the exception, it will be propagated
		if !handled {
			result = exception
		}
	}
	
	// Execute finally block if present
	if node.FinallyBlock != nil {
		finallyResult := evalBlockStatement(node.FinallyBlock, env)
		
		// If finally block throws an exception, it overrides any previous exception
		if ex, ok := finallyResult.(*Exception); ok {
			return ex
		}
		
		// If finally block returns a value, it overrides the try/catch result
		if returnVal, ok := finallyResult.(*ReturnValue); ok {
			return returnVal
		}
	}
	
	return result
}

// shouldCatchError determines if a catch clause should handle an error
func shouldCatchError(error Value, catchClause *ast.CatchClause) bool {
	// If no error type specified, catch all errors
	if catchClause.ErrorType == nil {
		return true
	}
	
	// Check if the error type matches
	if errorObj, ok := error.(*Error); ok {
		return errorObj.ErrorType == catchClause.ErrorType.Value
	}
	
	return false
}

// evalClassDeclaration evaluates class declarations
func evalClassDeclaration(node *ast.ClassDeclaration, env *Environment) Value {
  class := &Class{
    Name:    node.Name.Value,
    Methods: make(map[string]*Function),
    Env:     NewEnclosedEnvironment(env),
  }

  // Handle inheritance
  if node.SuperClass != nil {
    superVal := Eval(node.SuperClass, env)
    if isError(superVal) {
      return superVal
    }
    if superClass, ok := superVal.(*Class); ok {
      class.SuperClass = superClass
    } else {
      return newError("super class must be a class, got %T", superVal)
    }
  }

  // Parse methods from the class body
  if node.Body != nil {
    for _, stmt := range node.Body.Statements {
      if method, ok := stmt.(*ast.MethodDeclaration); ok {
        methodFunc := &Function{
          Parameters: method.Parameters,
          Body:       method.Body,
          Env:        class.Env,
        }
        class.Methods[method.Name.Value] = methodFunc
      }
    }
  }

  // Store the class in the environment
  env.Set(node.Name.Value, class)
  return class
}

// evalInstanceVariable evaluates instance variable access like "@variable"
func evalInstanceVariable(node *ast.InstanceVariable, env *Environment) Value {
  // Look for the current object (self) in the environment
  if self, exists := env.Get("self"); exists {
    if obj, ok := self.(*Object); ok {
      if value, exists := obj.InstanceVars[node.Name.Value]; exists {
        return value
      }
      return NULL // Instance variable not set yet
    }
  }
  return newError("instance variable %s used outside of object context", node.Name.Value)
}

// evalNewExpression evaluates object instantiation like "ClassName.new()"
func evalNewExpression(node *ast.NewExpression, env *Environment) Value {
  // Get the class
  classVal := Eval(node.ClassName, env)
  if isError(classVal) {
    return classVal
  }

  class, ok := classVal.(*Class)
  if !ok {
    return newError("not a class: %T", classVal)
  }

  // Create new object instance
  obj := &Object{
    Class:        class,
    InstanceVars: make(map[string]Value),
    Env:          NewEnclosedEnvironment(class.Env),
  }

  // Set self in the object's environment
  obj.Env.Set("self", obj)

  // Call initialize method if it exists
  if initMethod, exists := class.Methods["initialize"]; exists {
    // Evaluate arguments
    args := []Value{}
    for _, arg := range node.Arguments {
      evaluated := Eval(arg, env)
      if isError(evaluated) {
        return evaluated
      }
      args = append(args, evaluated)
    }

    // Set up method call environment with 'self' and parameters
    initEnv := NewEnclosedEnvironment(initMethod.Env)
    initEnv.Set("self", obj)
    
    // Check argument count
    if len(args) != len(initMethod.Parameters) {
      return newError("wrong number of arguments for initialize: want=%d, got=%d", len(initMethod.Parameters), len(args))
    }
    
    // Set up parameters in method environment
    for paramIdx, param := range initMethod.Parameters {
      initEnv.Set(param.Value, args[paramIdx])
    }
    
    // Evaluate method body with proper environment
    result := Eval(initMethod.Body, initEnv)
    if isError(result) {
      return result
    }
  }

  return obj
}

// resolveMethod walks up the inheritance chain to find a method
func resolveMethod(class *Class, methodName string) *Function {
  current := class
  for current != nil {
    if method, exists := current.Methods[methodName]; exists {
      return method
    }
    current = current.SuperClass
  }
  return nil
}

// evalSuperExpression evaluates super() calls to parent methods
func evalSuperExpression(node *ast.SuperExpression, env *Environment) Value {
  // Get the current object (self)
  self, exists := env.Get("self")
  if !exists {
    return newError("super() can only be used within a method")
  }
  
  obj, ok := self.(*Object)
  if !ok {
    return newError("super() can only be used within a method")
  }

  // Get the current method name by walking up the call stack
  // For now, we'll assume it's stored in a special environment variable
  currentMethodNameVal, exists := env.Get("__current_method__")
  if !exists {
    return newError("super() can only be used within a method")
  }
  
  currentMethodName, ok := currentMethodNameVal.(*String)
  if !ok {
    return newError("invalid method context for super()")
  }

  // Find the parent class method
  if obj.Class.SuperClass == nil {
    return newError("super() called but no superclass exists")
  }
  
  method := resolveMethod(obj.Class.SuperClass, currentMethodName.Value)
  if method == nil {
    return newError("super method %s not found", currentMethodName.Value)
  }

  // Evaluate arguments
  args := evalExpressions(node.Arguments, env)
  if len(args) == 1 && isError(args[0]) {
    return args[0]
  }

  // Check argument count
  if len(args) != len(method.Parameters) {
    return newError("wrong number of arguments for super(): want=%d, got=%d", len(method.Parameters), len(args))
  }

  // Set up method call environment with 'self' and parameters
  methodEnv := NewEnclosedEnvironment(method.Env)
  methodEnv.Set("self", obj)
  methodEnv.Set("__current_method__", currentMethodName)

  // Set up parameters in method environment
  for paramIdx, param := range method.Parameters {
    methodEnv.Set(param.Value, args[paramIdx])
  }

  // Evaluate method body with proper environment
  result := Eval(method.Body, methodEnv)
  return unwrapReturnValue(result)
}

// isStandardLibraryModule checks if a module path is for a standard library module
func isStandardLibraryModule(modulePath string) bool {
	return strings.HasPrefix(modulePath, "std/")
}

// addNativeStandardLibraryFunctions adds native functions to standard library modules
func addNativeStandardLibraryFunctions(env *Environment, modulePath string) {
	switch modulePath {
	case "std/string":
		// Add native string functions
		if builtin, exists := builtins["substr"]; exists {
			env.AddExport("substr", builtin)
		}
		if builtin, exists := builtins["split"]; exists {
			env.AddExport("split", builtin)
		}
	case "std/array":
		// Add native array functions
		if builtin, exists := builtins["push"]; exists {
			env.AddExport("push", builtin)
		}
		if builtin, exists := builtins["pop"]; exists {
			env.AddExport("pop", builtin)
		}
		if builtin, exists := builtins["slice"]; exists {
			env.AddExport("slice", builtin)
		}
	}
}

// evalIndexAssignment evaluates array element assignments like "arr[0] = 5"
func evalIndexAssignment(node *ast.IndexAssignmentStatement, env *Environment) Value {
	// Evaluate the array/string being indexed
	left := Eval(node.Left.Left, env)
	if isError(left) {
		return left
	}
	
	// Evaluate the index
	index := Eval(node.Left.Index, env)
	if isError(index) {
		return index
	}
	
	// Evaluate the value to assign
	value := Eval(node.Value, env)
	if isError(value) {
		return value
	}
	
	// Handle array index assignment
	if arr, ok := left.(*Array); ok {
		return evalArrayIndexAssignment(arr, index, value)
	}
	
	// Handle hash index assignment
	if hash, ok := left.(*Hash); ok {
		return evalHashIndexAssignment(hash, index, value, env)
	}
	
	// String index assignment is not supported (strings are immutable)
	if _, ok := left.(*String); ok {
		return newError("string index assignment not supported: strings are immutable")
	}
	
	return newError("index assignment not supported on type: %s", left.Type())
}

// evalArrayIndexAssignment handles assignment to array elements
func evalArrayIndexAssignment(array *Array, index Value, value Value) Value {
	idx, ok := index.(*Integer)
	if !ok {
		return newError("array index must be an integer, got %s", index.Type())
	}
	
	arrayLen := int64(len(array.Elements))
	i := idx.Value
	
	// Check bounds
	if i < 0 || i >= arrayLen {
		return newError("array index out of bounds: index %d, length %d", i, arrayLen)
	}
	
	// Assign the value
	array.Elements[i] = value
	return value
}

// evalHashIndexAssignment handles assignment to hash elements
func evalHashIndexAssignment(hash *Hash, index Value, value Value, env *Environment) Value {
	if !isHashable(index) {
		return newError("unusable as hash key: %T", index)
	}

	hashKey := CreateHashKey(index)
	
	// If the key doesn't exist, add it to the keys slice
	if _, exists := hash.Pairs[hashKey]; !exists {
		hash.Keys = append(hash.Keys, index)
	}
	
	// Set the value (mutating the hash in place)
	hash.Pairs[hashKey] = value
	
	return value
}

// Hash method helper functions

func hashSet(hash *Hash, key, value Value) Value {
	hashKey := CreateHashKey(key)
	newPairs := make(map[HashKey]Value)
	for k, v := range hash.Pairs {
		newPairs[k] = v
	}
	
	// Check if key already exists
	if _, exists := newPairs[hashKey]; !exists {
		// New key, add to keys array
		newKeys := make([]Value, len(hash.Keys)+1)
		copy(newKeys, hash.Keys)
		newKeys[len(hash.Keys)] = key
		newPairs[hashKey] = value
		return &Hash{Pairs: newPairs, Keys: newKeys}
	} else {
		// Existing key, just update value
		newPairs[hashKey] = value
		return &Hash{Pairs: newPairs, Keys: hash.Keys}
	}
}

func hashDelete(hash *Hash, key Value) Value {
	hashKey := CreateHashKey(key)
	if _, exists := hash.Pairs[hashKey]; !exists {
		// Key doesn't exist, return original hash
		return hash
	}
	
	newPairs := make(map[HashKey]Value)
	for k, v := range hash.Pairs {
		if k != hashKey {
			newPairs[k] = v
		}
	}
	
	// Remove key from keys array
	newKeys := make([]Value, 0, len(hash.Keys)-1)
	for _, k := range hash.Keys {
		if !compareValues(k, key) {
			newKeys = append(newKeys, k)
		}
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashMerge(hash1, hash2 *Hash) Value {
	newPairs := make(map[HashKey]Value)
	
	// Copy all pairs from hash1
	for k, v := range hash1.Pairs {
		newPairs[k] = v
	}
	
	// Copy all pairs from hash2 (overwrites conflicts)
	for k, v := range hash2.Pairs {
		newPairs[k] = v
	}
	
	// Build new keys array maintaining order
	newKeys := make([]Value, 0, len(newPairs))
	keyExists := make(map[string]bool)
	
	// Add keys from hash1 first
	for _, key := range hash1.Keys {
		keyStr := key.Inspect()
		if !keyExists[keyStr] {
			newKeys = append(newKeys, key)
			keyExists[keyStr] = true
		}
	}
	
	// Add new keys from hash2
	for _, key := range hash2.Keys {
		keyStr := key.Inspect()
		if !keyExists[keyStr] {
			newKeys = append(newKeys, key)
			keyExists[keyStr] = true
		}
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashFilter(hash *Hash, predicate *Function, env *Environment) Value {
	newPairs := make(map[HashKey]Value)
	newKeys := make([]Value, 0)
	
	for _, key := range hash.Keys {
		hashKey := CreateHashKey(key)
		value := hash.Pairs[hashKey]
		
		// Call predicate with key and value
		args := []Value{key, value}
		// Create a dummy call expression for the function call
		dummyCall := &ast.CallExpression{
			Function:  &ast.Identifier{Value: "predicate"},
			Arguments: []ast.Expression{},
		}
		result := applyFunction(predicate, args, dummyCall, env)
		
		if isError(result) {
			return result
		}
		
		if IsTruthy(result) {
			newPairs[hashKey] = value
			newKeys = append(newKeys, key)
		}
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashMapValues(hash *Hash, transform *Function, env *Environment) Value {
	newPairs := make(map[HashKey]Value)
	
	for _, key := range hash.Keys {
		hashKey := CreateHashKey(key)
		value := hash.Pairs[hashKey]
		
		// Call transform with value
		args := []Value{value}
		// Create a dummy call expression for the function call
		dummyCall := &ast.CallExpression{
			Function:  &ast.Identifier{Value: "transform"},
			Arguments: []ast.Expression{},
		}
		result := applyFunction(transform, args, dummyCall, env)
		
		if isError(result) {
			return result
		}
		
		newPairs[hashKey] = result
	}
	
	return &Hash{Pairs: newPairs, Keys: hash.Keys}
}

func hashEach(hash *Hash, callback *Function, env *Environment) {
	for _, key := range hash.Keys {
		hashKey := CreateHashKey(key)
		value := hash.Pairs[hashKey]
		
		// Call callback with key and value
		args := []Value{key, value}
		// Create a dummy call expression for the function call
		dummyCall := &ast.CallExpression{
			Function:  &ast.Identifier{Value: "callback"},
			Arguments: []ast.Expression{},
		}
		applyFunction(callback, args, dummyCall, env)
	}
}

func hashSelectKeys(hash *Hash, keyArray *Array) Value {
	newPairs := make(map[HashKey]Value)
	newKeys := make([]Value, 0)
	
	for _, key := range keyArray.Elements {
		hashKey := CreateHashKey(key)
		if value, exists := hash.Pairs[hashKey]; exists {
			newPairs[hashKey] = value
			newKeys = append(newKeys, key)
		}
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashRejectKeys(hash *Hash, keyArray *Array) Value {
	// Create set of keys to reject
	rejectSet := make(map[string]bool)
	for _, key := range keyArray.Elements {
		rejectSet[key.Inspect()] = true
	}
	
	newPairs := make(map[HashKey]Value)
	newKeys := make([]Value, 0)
	
	for _, key := range hash.Keys {
		if !rejectSet[key.Inspect()] {
			hashKey := CreateHashKey(key)
			newPairs[hashKey] = hash.Pairs[hashKey]
			newKeys = append(newKeys, key)
		}
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashInvert(hash *Hash) Value {
	newPairs := make(map[HashKey]Value)
	newKeys := make([]Value, 0, len(hash.Keys))
	
	for _, key := range hash.Keys {
		hashKey := CreateHashKey(key)
		value := hash.Pairs[hashKey]
		
		// Use value as new key, key as new value
		newHashKey := CreateHashKey(value)
		newPairs[newHashKey] = key
		newKeys = append(newKeys, value)
	}
	
	return &Hash{Pairs: newPairs, Keys: newKeys}
}

func hashToArray(hash *Hash) Value {
	pairs := make([]Value, 0, len(hash.Keys))
	
	for _, key := range hash.Keys {
		hashKey := CreateHashKey(key)
		value := hash.Pairs[hashKey]
		pair := &Array{Elements: []Value{key, value}}
		pairs = append(pairs, pair)
	}
	
	return &Array{Elements: pairs}
}

// applyFileMethod handles file method calls
func applyFileMethod(fileMethod *FileMethod, args []Value, env *Environment) Value {
	file := fileMethod.File
	
	switch fileMethod.Method {
	case "open":
		if len(args) > 1 {
			return newError("wrong number of arguments for file.open: want=0 or 1, got=%d", len(args))
		}
		
		if file.IsOpen {
			return newError("file is already open: %s", file.Path)
		}
		
		// Default to read mode
		mode := "r"
		if len(args) == 1 {
			modeArg, ok := args[0].(*String)
			if !ok {
				return newError("file mode argument must be STRING")
			}
			mode = modeArg.Value
		}
		
		// Validate mode
		switch mode {
		case "r", "w", "a", "r+", "w+", "a+":
			// Valid modes
		default:
			return newError("invalid file mode: %s", mode)
		}
		
		// Try to open the file
		var handle *os.File
		var err error
		
		switch mode {
		case "r":
			handle, err = os.Open(file.Path)
		case "w":
			handle, err = os.Create(file.Path)
		case "a":
			handle, err = os.OpenFile(file.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		case "r+":
			handle, err = os.OpenFile(file.Path, os.O_RDWR, 0644)
		case "w+":
			handle, err = os.OpenFile(file.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		case "a+":
			handle, err = os.OpenFile(file.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		}
		
		if err != nil {
			return newError("failed to open file %s: %s", file.Path, err.Error())
		}
		
		file.Handle = handle
		file.IsOpen = true
		return file
		
	case "read":
		if len(args) != 0 {
			return newError("wrong number of arguments for file.read: want=0, got=%d", len(args))
		}
		
		if !file.IsOpen {
			return newError("file is not open: %s", file.Path)
		}
		
		handle, ok := file.Handle.(*os.File)
		if !ok {
			return newError("invalid file handle")
		}
		
		content, err := ioutil.ReadAll(handle)
		if err != nil {
			return newError("failed to read file %s: %s", file.Path, err.Error())
		}
		
		return &String{Value: string(content)}
		
	case "write":
		if len(args) != 1 {
			return newError("wrong number of arguments for file.write: want=1, got=%d", len(args))
		}
		
		content, ok := args[0].(*String)
		if !ok {
			return newError("file content argument must be STRING")
		}
		
		if !file.IsOpen {
			return newError("file is not open: %s", file.Path)
		}
		
		handle, ok := file.Handle.(*os.File)
		if !ok {
			return newError("invalid file handle")
		}
		
		_, err := handle.WriteString(content.Value)
		if err != nil {
			return newError("failed to write to file %s: %s", file.Path, err.Error())
		}
		
		return &Integer{Value: int64(len(content.Value))}
		
	case "close":
		if len(args) != 0 {
			return newError("wrong number of arguments for file.close: want=0, got=%d", len(args))
		}
		
		if !file.IsOpen {
			return newError("file is not open: %s", file.Path)
		}
		
		handle, ok := file.Handle.(*os.File)
		if !ok {
			return newError("invalid file handle")
		}
		
		err := handle.Close()
		if err != nil {
			return newError("failed to close file %s: %s", file.Path, err.Error())
		}
		
		file.Handle = nil
		file.IsOpen = false
		return file
		
	case "exists?":
		if len(args) != 0 {
			return newError("wrong number of arguments for file.exists?: want=0, got=%d", len(args))
		}
		
		_, err := os.Stat(file.Path)
		return &Boolean{Value: !os.IsNotExist(err)}
		
	case "size":
		if len(args) != 0 {
			return newError("wrong number of arguments for file.size: want=0, got=%d", len(args))
		}
		
		stat, err := os.Stat(file.Path)
		if os.IsNotExist(err) {
			return newError("file does not exist: %s", file.Path)
		}
		if err != nil {
			return newError("failed to get file size for %s: %s", file.Path, err.Error())
		}
		
		return &Integer{Value: stat.Size()}
		
	case "delete":
		if len(args) != 0 {
			return newError("wrong number of arguments for file.delete: want=0, got=%d", len(args))
		}
		
		if file.IsOpen {
			// Close the file first
			if handle, ok := file.Handle.(*os.File); ok {
				handle.Close()
			}
			file.Handle = nil
			file.IsOpen = false
		}
		
		err := os.Remove(file.Path)
		if err != nil {
			return newError("failed to delete file %s: %s", file.Path, err.Error())
		}
		
		return TRUE
		
	default:
		return newError("unknown file method: %s", fileMethod.Method)
	}
}

// applyDirectoryMethod handles directory method calls
func applyDirectoryMethod(dirMethod *DirectoryMethod, args []Value, env *Environment) Value {
	dir := dirMethod.Directory
	
	switch dirMethod.Method {
	case "create":
		if len(args) != 0 {
			return newError("wrong number of arguments for directory.create: want=0, got=%d", len(args))
		}
		
		err := os.MkdirAll(dir.Path, 0755)
		if err != nil {
			return newError("failed to create directory %s: %s", dir.Path, err.Error())
		}
		
		return dir
		
	case "list":
		if len(args) != 0 {
			return newError("wrong number of arguments for directory.list: want=0, got=%d", len(args))
		}
		
		entries, err := ioutil.ReadDir(dir.Path)
		if err != nil {
			return newError("failed to list directory %s: %s", dir.Path, err.Error())
		}
		
		names := make([]Value, 0, len(entries))
		for _, entry := range entries {
			names = append(names, &String{Value: entry.Name()})
		}
		
		return &Array{Elements: names}
		
	case "delete":
		if len(args) != 0 {
			return newError("wrong number of arguments for directory.delete: want=0, got=%d", len(args))
		}
		
		err := os.RemoveAll(dir.Path)
		if err != nil {
			return newError("failed to delete directory %s: %s", dir.Path, err.Error())
		}
		
		return TRUE
		
	case "exists?":
		if len(args) != 0 {
			return newError("wrong number of arguments for directory.exists?: want=0, got=%d", len(args))
		}
		
		stat, err := os.Stat(dir.Path)
		if os.IsNotExist(err) {
			return FALSE
		}
		if err != nil {
			return newError("failed to check directory existence for %s: %s", dir.Path, err.Error())
		}
		
		return &Boolean{Value: stat.IsDir()}
		
	default:
		return newError("unknown directory method: %s", dirMethod.Method)
	}
}

// applyPathMethod handles path method calls
func applyPathMethod(pathMethod *PathMethod, args []Value, env *Environment) Value {
	path := pathMethod.Path
	
	switch pathMethod.Method {
	case "join":
		if len(args) != 1 {
			return newError("wrong number of arguments for path.join: want=1, got=%d", len(args))
		}
		
		other, ok := args[0].(*String)
		if !ok {
			return newError("path join argument must be STRING")
		}
		
		joined := filepath.Join(path.Value, other.Value)
		return &Path{Value: joined}
		
	case "basename":
		if len(args) != 0 {
			return newError("wrong number of arguments for path.basename: want=0, got=%d", len(args))
		}
		
		base := filepath.Base(path.Value)
		return &String{Value: base}
		
	case "dirname":
		if len(args) != 0 {
			return newError("wrong number of arguments for path.dirname: want=0, got=%d", len(args))
		}
		
		dir := filepath.Dir(path.Value)
		return &String{Value: dir}
		
	case "absolute":
		if len(args) != 0 {
			return newError("wrong number of arguments for path.absolute: want=0, got=%d", len(args))
		}
		
		abs, err := filepath.Abs(path.Value)
		if err != nil {
			return newError("failed to get absolute path for %s: %s", path.Value, err.Error())
		}
		
		return &Path{Value: abs}
		
	case "clean":
		if len(args) != 0 {
			return newError("wrong number of arguments for path.clean: want=0, got=%d", len(args))
		}
		
		clean := filepath.Clean(path.Value)
		return &Path{Value: clean}
		
	default:
		return newError("unknown path method: %s", pathMethod.Method)
	}
}