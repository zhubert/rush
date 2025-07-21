package interpreter

import (
	"fmt"
	"path/filepath"

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
				if method, exists := obj.Class.Methods[methodName]; exists {
					// Set up method call environment with 'self' and parameters
					methodEnv := NewEnclosedEnvironment(method.Env)
					methodEnv.Set("self", obj)
					
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
		return applyFunction(function, args, node, env)
	
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	
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
			if rt == RETURN_VALUE || rt == ERROR_VALUE || rt == EXCEPTION_VALUE {
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
	// Don't unwrap exceptions - they should propagate as-is
	return val
}

func evalIndexExpression(left, index Value) Value {
	switch {
	case left.Type() == ARRAY_VALUE && index.Type() == INTEGER_VALUE:
		return evalArrayIndexExpression(left, index)
	case left.Type() == STRING_VALUE && index.Type() == INTEGER_VALUE:
		return evalStringIndexExpression(left, index)
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
		}
	}

	return result
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
	
	// Import the specified names into the current environment
	for _, name := range node.Names {
		if value, exists := module.Exports[name.Value]; exists {
			// Convert interface{} back to Value
			if val, ok := value.(Value); ok {
				env.Set(name.Value, val)
			} else {
				return newError("invalid export type for %s", name.Value)
			}
		} else {
			return newError("module %s does not export %s", modulePath, name.Value)
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
	
	// Check if it's an object instance and return method
	if obj, ok := object.(*Object); ok {
		methodName := node.Property.Value
		if method, exists := obj.Class.Methods[methodName]; exists {
			// Return a bound method (for now, just return the function)
			// In a full implementation, we'd create a bound method that includes 'self'
			return method
		}
		return newError("undefined method %s for class %s", methodName, obj.Class.Name)
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