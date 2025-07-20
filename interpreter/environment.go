package interpreter

import (
	"rush/module"
)

// Environment represents a scope for variable storage
type Environment struct {
	store          map[string]Value
	outer          *Environment
	moduleResolver *module.ModuleResolver
	currentDir     string // current directory for module resolution
	exports        map[string]Value // for tracking exports in modules
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
func (e *Environment) Set(name string, val Value) Value {
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