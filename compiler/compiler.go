package compiler

import (
	"fmt"
	"sort"

	"rush/ast"
	"rush/bytecode"
	"rush/interpreter"
)

// EmittedInstruction represents an instruction that has been emitted
type EmittedInstruction struct {
	Opcode   bytecode.Opcode
	Position int
}

// CompilationScope represents a compilation scope (function, class, etc.)
type CompilationScope struct {
	instructions        bytecode.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// Compiler transforms AST nodes into bytecode instructions
type Compiler struct {
	constants   []interpreter.Value // Constant pool
	symbolTable *SymbolTable        // Symbol table for variables
	scopes      []CompilationScope  // Compilation scopes stack
	scopeIndex  int                 // Current scope index
}

// Bytecode represents the compilation result
type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []interpreter.Value
}

// New creates a new compiler instance
func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        bytecode.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	
	// Add builtin functions to symbol table
	for i, name := range interpreter.Builtins {
		symbolTable.DefineBuiltin(i, name)
	}

	return &Compiler{
		constants:   []interpreter.Value{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

// NewWithState creates a new compiler with existing state (for closures)
func NewWithState(s *SymbolTable, constants []interpreter.Value) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

// Compile transforms an AST node into bytecode
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpPop)

	case *ast.IntegerLiteral:
		integer := &interpreter.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))

	case *ast.FloatLiteral:
		float := &interpreter.Float{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(float))

	case *ast.StringLiteral:
		str := &interpreter.String{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(str))

	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(bytecode.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		// Sort pairs to ensure deterministic compilation
		pairs := []ast.HashPair{}
		for _, pair := range node.Pairs {
			pairs = append(pairs, pair)
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Key.String() < pairs[j].Key.String()
		})

		for _, pair := range pairs {
			err := c.Compile(pair.Key)
			if err != nil {
				return err
			}
			err = c.Compile(pair.Value)
			if err != nil {
				return err
			}
		}
		c.emit(bytecode.OpHash, len(pairs))

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(bytecode.OpNot)
		case "-":
			c.emit(bytecode.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			// Rewrite a < b as b > a
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(bytecode.OpGreaterThan)
			return nil
		}

		if node.Operator == "<=" {
			// Rewrite a <= b as b >= a  
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(bytecode.OpGreaterEqual)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		case "-":
			c.emit(bytecode.OpSub)
		case "*":
			c.emit(bytecode.OpMul)
		case "/":
			c.emit(bytecode.OpDiv)
		case "%":
			c.emit(bytecode.OpMod)
		case ">":
			c.emit(bytecode.OpGreaterThan)
		case ">=":
			c.emit(bytecode.OpGreaterEqual)
		case "==":
			c.emit(bytecode.OpEqual)
		case "!=":
			c.emit(bytecode.OpNotEqual)
		case "&&":
			c.emit(bytecode.OpAnd)
		case "||":
			c.emit(bytecode.OpOr)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit OpJumpNotTruthy with placeholder offset
		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		// Emit jump to skip alternative (if present)
		jumpPos := c.emit(bytecode.OpJump, 9999)

		jumpNotTruthyAddr := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, jumpNotTruthyAddr)

		if node.Alternative == nil {
			c.emit(bytecode.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}
		}

		jumpAddr := len(c.currentInstructions())
		c.changeOperand(jumpPos, jumpAddr)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	case *ast.AssignmentStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// Check if this is an instance variable assignment
		if len(node.Name.Value) > 0 && node.Name.Value[0] == '@' {
			// Instance variable assignment: @name = value
			varName := node.Name.Value[1:] // Remove @ prefix
			varNameStr := &interpreter.String{Value: varName}
			varNameIndex := c.addConstant(varNameStr)
			c.emit(bytecode.OpSetInstance, varNameIndex)
		} else {
			// Regular variable assignment
			symbol := c.symbolTable.Define(node.Name.Value)
			c.storeSymbol(symbol)
		}

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpIndex)

	case *ast.IndexAssignmentStatement:
		err := c.Compile(node.Left.Left) // array/hash
		if err != nil {
			return err
		}
		err = c.Compile(node.Left.Index) // index
		if err != nil {
			return err
		}
		err = c.Compile(node.Value) // value
		if err != nil {
			return err
		}
		c.emit(bytecode.OpSetIndex)

	case *ast.PropertyAccess:
		err := c.Compile(node.Object)
		if err != nil {
			return err
		}
		
		propertyName := &interpreter.String{Value: node.Property.Value}
		c.emit(bytecode.OpGetProperty, c.addConstant(propertyName))

	case *ast.FunctionLiteral:
		c.enterScope()

		// Define parameters as local variables
		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		// Functions that don't explicitly return should return null
		if c.lastInstructionIs(bytecode.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(bytecode.OpReturn) {
			c.emit(bytecode.OpReturnVoid)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &interpreter.CompiledFunction{
			Instructions:  []byte(instructions),
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(bytecode.OpClosure, fnIndex, len(freeSymbols))

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(bytecode.OpCall, len(node.Arguments))

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpReturn)

	case *ast.WhileStatement:
		loopStart := len(c.currentInstructions())

		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

		c.emit(bytecode.OpJump, loopStart)

		jumpNotTruthyAddr := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, jumpNotTruthyAddr)

	case *ast.ForStatement:
		// Compile initialization
		if node.Init != nil {
			err := c.Compile(node.Init)
			if err != nil {
				return err
			}
		}

		loopStart := len(c.currentInstructions())

		// Compile condition
		if node.Condition != nil {
			err := c.Compile(node.Condition)
			if err != nil {
				return err
			}
		} else {
			c.emit(bytecode.OpTrue) // Infinite loop if no condition
		}

		jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		// Compile body
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		// Compile update
		if node.Update != nil {
			err := c.Compile(node.Update)
			if err != nil {
				return err
			}
		}

		c.emit(bytecode.OpJump, loopStart)

		jumpNotTruthyAddr := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, jumpNotTruthyAddr)

	case *ast.BreakStatement:
		c.emit(bytecode.OpBreak)

	case *ast.ContinueStatement:
		c.emit(bytecode.OpContinue)

	case *ast.SwitchStatement:
		// Compile the switch value expression
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// Store jump positions for cases and default
		jumpTable := make([]int, len(node.Cases))
		var defaultJumpPos int
		var endJumpPos int

		// Compile case comparisons and collect jump positions
		for i, caseClause := range node.Cases {
			var skipPos int
			// For each case value, duplicate switch value and compare
			for j, caseValue := range caseClause.Values {
				if j > 0 {
					// Multiple values in one case - use OR logic
					// If any previous comparison was true, skip this one
					skipPos = c.emit(bytecode.OpJumpTruthy, 9999)
					c.emit(bytecode.OpDup) // Duplicate switch value again
				} else {
					c.emit(bytecode.OpDup) // Duplicate switch value
				}
				
				err := c.Compile(caseValue)
				if err != nil {
					return err
				}
				c.emit(bytecode.OpEqual)
				
				if j > 0 {
					// Patch the skip jump
					c.changeOperand(skipPos, len(c.currentInstructions()))
				}
			}
			
			// Jump to case body if match found
			jumpTable[i] = c.emit(bytecode.OpJumpTruthy, 9999)
		}

		// If no cases matched, jump to default (if exists) or end
		if node.Default != nil {
			defaultJumpPos = c.emit(bytecode.OpJump, 9999)
		} else {
			// No default, jump to end
			endJumpPos = c.emit(bytecode.OpJump, 9999)
		}

		// Compile case bodies
		caseEndJumps := make([]int, 0)
		for i, caseClause := range node.Cases {
			// Patch jump to this case body
			c.changeOperand(jumpTable[i], len(c.currentInstructions()))
			
			// Pop the switch value from stack
			c.emit(bytecode.OpPop)
			
			// Compile case body
			err := c.Compile(caseClause.Body)
			if err != nil {
				return err
			}
			
			// Jump to end (break behavior) unless it's the last case
			if i < len(node.Cases)-1 || node.Default != nil {
				caseEndJumps = append(caseEndJumps, c.emit(bytecode.OpJump, 9999))
			}
		}

		// Compile default case if it exists
		if node.Default != nil {
			c.changeOperand(defaultJumpPos, len(c.currentInstructions()))
			c.emit(bytecode.OpPop) // Pop switch value
			err := c.Compile(node.Default.Body)
			if err != nil {
				return err
			}
		}

		// Patch all end jumps to point to the end
		endPos := len(c.currentInstructions())
		for _, jumpPos := range caseEndJumps {
			c.changeOperand(jumpPos, endPos)
		}
		if node.Default == nil && endJumpPos != 0 {
			c.changeOperand(endJumpPos, endPos)
		}

	case *ast.ThrowStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpThrow)

	case *ast.TryStatement:
		// Begin try block with catch handler address
		catchJumpPos := c.emit(bytecode.OpTryBegin, 9999)

		// Compile try block
		err := c.Compile(node.TryBlock)
		if err != nil {
			return err
		}

		// End try block
		c.emit(bytecode.OpTryEnd)

		// Jump over catch blocks if no exception
		endJumpPos := c.emit(bytecode.OpJump, 9999)

		// Patch catch handler address
		c.changeOperand(catchJumpPos, len(c.currentInstructions()))

		// Compile catch clauses
		for _, catchClause := range node.CatchClauses {
			if catchClause.ErrorType != nil {
				// Type-specific catch
				errorTypeName := &interpreter.String{Value: catchClause.ErrorType.Value}
				c.emit(bytecode.OpCatch, c.addConstant(errorTypeName))
			} else {
				// Catch-all
				c.emit(bytecode.OpCatch, c.addConstant(&interpreter.String{Value: ""}))
			}

			// Set error variable in current scope
			if catchClause.ErrorVar != nil {
				symbol := c.symbolTable.Define(catchClause.ErrorVar.Value)
				c.storeSymbol(symbol)
			}

			// Compile catch block
			err := c.Compile(catchClause.Body)
			if err != nil {
				return err
			}

			// Jump to finally or end
			c.emit(bytecode.OpJump, 9999) // Will be patched later
		}

		// Compile finally block if it exists
		if node.FinallyBlock != nil {
			c.emit(bytecode.OpFinally)
			err := c.Compile(node.FinallyBlock)
			if err != nil {
				return err
			}
		}

		// Patch end jump
		c.changeOperand(endJumpPos, len(c.currentInstructions()))

	case *ast.ImportStatement:
		// Compile import statement
		// For now, we'll implement a simplified version
		moduleNameStr := node.Module.Value
		moduleName := &interpreter.String{Value: moduleNameStr}
		c.emit(bytecode.OpImport, c.addConstant(moduleName))

		// For each imported item, create a binding in current scope
		for _, item := range node.Items {
			// Define the imported name in symbol table
			var bindingName string
			if item.Alias != nil {
				bindingName = item.Alias.Value
			} else {
				bindingName = item.Name.Value
			}
			
			symbol := c.symbolTable.Define(bindingName)
			c.storeSymbol(symbol)
		}

	case *ast.ExportStatement:
		// Compile the export value
		if node.Value != nil {
			err := c.Compile(node.Value)
			if err != nil {
				return err
			}
		} else {
			// Export identifier directly
			symbol, ok := c.symbolTable.Resolve(node.Name.Value)
			if !ok {
				return fmt.Errorf("undefined variable %s", node.Name.Value)
			}
			c.loadSymbol(symbol)
		}

		// Emit export instruction with name
		exportName := &interpreter.String{Value: node.Name.Value}
		c.emit(bytecode.OpExport, c.addConstant(exportName))

	case *ast.ModuleAccess:
		// For module access like "module.member", we can treat it similar to property access
		// First get the module
		moduleSymbol, ok := c.symbolTable.Resolve(node.Module.Value)
		if !ok {
			return fmt.Errorf("undefined module %s", node.Module.Value)
		}
		c.loadSymbol(moduleSymbol)
		
		// Then access the member
		memberName := &interpreter.String{Value: node.Member.Value}
		c.emit(bytecode.OpGetProperty, c.addConstant(memberName))

	case *ast.ClassDeclaration:
		// Create class name constant
		className := &interpreter.String{Value: node.Name.Value}
		classNameIndex := c.addConstant(className)
		
		// Get methods from the class body
		var methods []*ast.MethodDeclaration
		if node.Body != nil {
			for _, stmt := range node.Body.Statements {
				if method, ok := stmt.(*ast.MethodDeclaration); ok {
					methods = append(methods, method)
				}
			}
		}
		
		// Handle superclass if present
		if node.SuperClass != nil {
			// Load superclass
			superSymbol, ok := c.symbolTable.Resolve(node.SuperClass.Value)
			if !ok {
				return fmt.Errorf("undefined superclass %s", node.SuperClass.Value)
			}
			c.loadSymbol(superSymbol)
			
			// Create class with inheritance
			c.emit(bytecode.OpClass, classNameIndex, len(methods))
			c.emit(bytecode.OpInherit)
		} else {
			// Create class without inheritance
			c.emit(bytecode.OpClass, classNameIndex, len(methods))
		}
		
		// Compile methods
		for _, method := range methods {
			// Enter new scope for method compilation
			c.enterScope()
			
			// Define parameters as local variables
			for _, p := range method.Parameters {
				c.symbolTable.Define(p.Value)
			}
			
			// Compile method body
			err := c.Compile(method.Body)
			if err != nil {
				return err
			}
			
			// Handle method return
			if c.lastInstructionIs(bytecode.OpPop) {
				c.replaceLastPopWithReturn()
			}
			if !c.lastInstructionIs(bytecode.OpReturn) {
				c.emit(bytecode.OpReturnVoid)
			}
			
			// Get method instructions and leave scope
			freeSymbols := c.symbolTable.FreeSymbols
			numLocals := c.symbolTable.numDefinitions
			instructions := c.leaveScope()
			
			// Load free variables
			for _, s := range freeSymbols {
				c.loadSymbol(s)
			}
			
			// Create compiled method and push to constants
			compiledMethod := &interpreter.CompiledFunction{
				Instructions:  []byte(instructions),
				NumLocals:     numLocals,
				NumParameters: len(method.Parameters),
			}
			
			// Push compiled method as closure
			methodIndex := c.addConstant(compiledMethod)
			c.emit(bytecode.OpClosure, methodIndex, len(freeSymbols))
			
			// Create method name constant and emit OpMethod
			methodName := &interpreter.String{Value: method.Name.Value}
			methodNameIndex := c.addConstant(methodName)
			c.emit(bytecode.OpMethod, methodNameIndex)
		}
		
		// Define class in symbol table
		symbol := c.symbolTable.Define(node.Name.Value)
		c.storeSymbol(symbol)

	case *ast.NewExpression:
		// Load class constructor
		classSymbol, ok := c.symbolTable.Resolve(node.ClassName.Value)
		if !ok {
			return fmt.Errorf("undefined class %s", node.ClassName.Value)
		}
		c.loadSymbol(classSymbol)
		
		// Compile constructor arguments
		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}
		
		// Call constructor (similar to function call but for class instantiation)
		c.emit(bytecode.OpCall, len(node.Arguments))

	case *ast.InstanceVariable:
		// Instance variable access using @ syntax
		varName := &interpreter.String{Value: node.Name.Value}
		varNameIndex := c.addConstant(varName)
		c.emit(bytecode.OpGetInstance, varNameIndex)

	case *ast.SuperExpression:
		// Super method call: super(args...)
		// The method name should be available from the context (current method being compiled)
		// For now, we'll use a placeholder - this might need to be enhanced with context tracking
		methodName := &interpreter.String{Value: "initialize"} // Default to constructor
		methodNameIndex := c.addConstant(methodName)
		
		// Compile arguments
		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}
		
		// Emit super call
		c.emit(bytecode.OpGetSuper, methodNameIndex)
		c.emit(bytecode.OpCall, len(node.Arguments))

	default:
		return fmt.Errorf("compilation not implemented for %T", node)
	}

	return nil
}

// Bytecode returns the compiled bytecode and constants
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// Helper methods

func (c *Compiler) addConstant(obj interpreter.Value) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	ins := bytecode.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)
	c.scopes[c.scopeIndex].instructions = updatedInstructions
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) currentInstructions() bytecode.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) lastInstructionIs(op bytecode.Opcode) bool {
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()
	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := bytecode.Opcode(c.currentInstructions()[opPos])
	newInstruction := bytecode.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) replaceLastPopWithReturn() {
	pos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(pos, bytecode.Make(bytecode.OpReturn))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = bytecode.OpReturn
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        bytecode.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() bytecode.Instructions {
	instructions := c.currentInstructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return instructions
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(bytecode.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(bytecode.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(bytecode.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(bytecode.OpGetFree, s.Index)
	}
}

func (c *Compiler) storeSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(bytecode.OpSetGlobal, s.Index)
	case LocalScope:
		c.emit(bytecode.OpSetLocal, s.Index)
	}
}