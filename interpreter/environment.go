package interpreter

import (
	"fmt"
	"strings"
	
	"rush/module"
)

// CallFrame represents a single function call in the call stack
type CallFrame struct {
	FunctionName string
	Line         int
	Column       int
}

// Environment represents a scope for variable storage
type Environment struct {
	store          map[string]Value
	outer          *Environment
	moduleResolver *module.ModuleResolver
	currentDir     string // current directory for module resolution
	exports        map[string]Value // for tracking exports in modules
	callStack      []CallFrame // for tracking function calls
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Value)
	env := &Environment{
		store:          s, 
		outer:          nil,
		moduleResolver: module.NewModuleResolver(),
		currentDir:     ".",
		exports:        make(map[string]Value),
		callStack:      make([]CallFrame, 0),
	}
	
	// Add built-in functions
	for name, builtin := range builtins {
		env.store[name] = builtin
	}
	
	return env
}

// NewEnclosedEnvironment creates a new environment that encloses an outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	// Inherit module resolver and current directory from outer environment
	if outer != nil {
		env.moduleResolver = outer.moduleResolver
		env.currentDir = outer.currentDir
		// Inherit call stack from outer environment
		env.callStack = make([]CallFrame, len(outer.callStack))
		copy(env.callStack, outer.callStack)
	}
	return env
}

// NewModuleEnvironment creates a new environment specifically for module execution
func NewModuleEnvironment(outer *Environment) *Environment {
	env := NewEnclosedEnvironment(outer)
	env.exports = make(map[string]Value) // Fresh exports map for the module
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Value, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

// Set stores a value in the environment
// If the variable exists in an outer scope, it updates it there
// Otherwise, it creates a new variable in the current scope
func (e *Environment) Set(name string, val Value) Value {
	// Check if the variable exists in the current environment
	if _, exists := e.store[name]; exists {
		e.store[name] = val
		return val
	}
	
	// Check if the variable exists in outer environments
	if e.outer != nil {
		if _, exists := e.outer.Get(name); exists {
			// Update the variable in the outer environment
			return e.outer.Set(name, val)
		}
	}
	
	// Variable doesn't exist anywhere, create it in current scope
	e.store[name] = val
	return val
}

// SetLocal always creates/updates a variable in the current environment only
// This is used for variable shadowing (like catch variables)
func (e *Environment) SetLocal(name string, val Value) Value {
	e.store[name] = val
	return val
}

// SetCurrentDir sets the current directory for module resolution
func (e *Environment) SetCurrentDir(dir string) {
	e.currentDir = dir
}

// GetCurrentDir returns the current directory for module resolution
func (e *Environment) GetCurrentDir() string {
	return e.currentDir
}

// GetModuleResolver returns the module resolver
func (e *Environment) GetModuleResolver() *module.ModuleResolver {
	return e.moduleResolver
}

// AddExport adds a value to the exports map
func (e *Environment) AddExport(name string, value Value) {
	e.exports[name] = value
}

// GetExports returns the exports map
func (e *Environment) GetExports() map[string]Value {
	return e.exports
}

// PushCall adds a function call to the call stack
func (e *Environment) PushCall(functionName string, line, column int) {
	frame := CallFrame{
		FunctionName: functionName,
		Line:         line,
		Column:       column,
	}
	e.callStack = append(e.callStack, frame)
}

// PopCall removes the most recent function call from the call stack
func (e *Environment) PopCall() {
	if len(e.callStack) > 0 {
		e.callStack = e.callStack[:len(e.callStack)-1]
	}
}

// GetStackTrace returns a formatted stack trace string
func (e *Environment) GetStackTrace() string {
	if len(e.callStack) == 0 {
		return ""
	}
	
	var frames []string
	for i := len(e.callStack) - 1; i >= 0; i-- {
		frame := e.callStack[i]
		if frame.Line > 0 {
			frames = append(frames, fmt.Sprintf("  at %s (line %d:%d)", frame.FunctionName, frame.Line, frame.Column))
		} else {
			frames = append(frames, fmt.Sprintf("  at %s", frame.FunctionName))
		}
	}
	return strings.Join(frames, "\n")
}

// GetCallStack returns a copy of the current call stack
func (e *Environment) GetCallStack() []CallFrame {
	stack := make([]CallFrame, len(e.callStack))
	copy(stack, e.callStack)
	return stack
}