package interpreter

// Environment represents a scope for variable storage
type Environment struct {
	store map[string]Value
	outer *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Value)
	env := &Environment{store: s, outer: nil}
	
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