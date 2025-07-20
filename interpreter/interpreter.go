package interpreter

import (
	"fmt"

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
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	
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
			if result.Type() == ERROR_VALUE {
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
	Message string
	Line    int
	Column  int
}

func (e *Error) Type() ValueType { return ERROR_VALUE }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("ERROR at line %d:%d: %s", e.Line, e.Column, e.Message)
	}
	return "ERROR: " + e.Message
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func newErrorWithPosition(line, column int, format string, a ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
		Line:    line,
		Column:  column,
	}
}

func isError(val Value) bool {
	if val != nil {
		return val.Type() == ERROR_VALUE
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
			if rt == RETURN_VALUE || rt == ERROR_VALUE {
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

func applyFunction(fn Value, args []Value) Value {
	switch fn := fn.(type) {
	case *Function:
		if len(args) != len(fn.Parameters) {
			return newError("wrong number of arguments: want=%d, got=%d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *BuiltinFunction:
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
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index Value) Value {
	stringObject := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
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
			if rt == RETURN_VALUE || rt == ERROR_VALUE {
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
			if rt == RETURN_VALUE || rt == ERROR_VALUE {
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