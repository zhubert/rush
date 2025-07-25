package compiler

// SymbolScope represents the scope of a symbol
type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

// Symbol represents a symbol in the symbol table
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable manages variable scoping and symbol resolution
type SymbolTable struct {
	Outer          *SymbolTable    // Parent scope
	store          map[string]Symbol
	numDefinitions int             // Number of definitions in current scope
	FreeSymbols    []Symbol        // Free variables (closures)
	isFunction     bool            // True if this is a function scope (not a block scope)
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free}
}

// NewEnclosedSymbolTable creates a new enclosed symbol table (for function scopes)
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	s.isFunction = true  // Mark as function scope
	return s
}

// NewEnclosedBlockTable creates a new enclosed symbol table (for block scopes like loops)
func NewEnclosedBlockTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	s.isFunction = false  // Mark as block scope
	return s
}

// Define adds a new symbol to the symbol table
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// DefineBuiltin adds a builtin function to the symbol table
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltinScope, Index: index}
	s.store[name] = symbol
	return symbol
}

// DefineFree adds a free variable to the symbol table (for closures)
func (s *SymbolTable) DefineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	
	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope
	
	s.store[original.Name] = symbol
	return symbol
}

// Resolve looks up a symbol in the symbol table hierarchy
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		// Convert to free variable only if it's a local variable (not global, builtin, or already free)
		if obj.Scope == LocalScope {
			return s.DefineFree(obj), true
		}
		// If it's already a free variable, pass it through directly
		if obj.Scope == FreeScope {
			return s.DefineFree(obj), true
		}
	}
	return obj, ok
}