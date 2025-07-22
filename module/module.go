package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"rush/ast"
	"rush/lexer"
	"rush/parser"
)

// Module represents a loaded module with its exports
type Module struct {
	Path    string
	Exports map[string]interface{} // Will be interpreter.Value when executed
	AST     *ast.Program
}

// ModuleResolver handles module loading and resolution
type ModuleResolver struct {
	cache     map[string]*Module
	loadStack []string // for circular dependency detection
}

// NewModuleResolver creates a new module resolver
func NewModuleResolver() *ModuleResolver {
	return &ModuleResolver{
		cache:     make(map[string]*Module),
		loadStack: []string{},
	}
}

// LoadModule loads a module from the given path
func (mr *ModuleResolver) LoadModule(modulePath string, baseDir string) (*Module, error) {
	// Resolve the actual file path
	resolvedPath, err := mr.resolvePath(modulePath, baseDir)
	if err != nil {
		return nil, err
	}

	// Check if already loaded
	if module, exists := mr.cache[resolvedPath]; exists {
		return module, nil
	}

	// Check for circular dependencies
	for _, path := range mr.loadStack {
		if path == resolvedPath {
			return nil, fmt.Errorf("circular dependency detected: %s", strings.Join(mr.loadStack, " -> "))
		}
	}

	// Add to load stack
	mr.loadStack = append(mr.loadStack, resolvedPath)
	defer func() {
		// Remove from load stack when done
		if len(mr.loadStack) > 0 {
			mr.loadStack = mr.loadStack[:len(mr.loadStack)-1]
		}
	}()

	// Read the module file
	content, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module %s: %v", resolvedPath, err)
	}

	// Parse the module
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parse errors
	if errors := p.Errors(); len(errors) > 0 {
		return nil, fmt.Errorf("parse errors in module %s: %v", resolvedPath, errors)
	}

	// Create module
	module := &Module{
		Path:    resolvedPath,
		Exports: make(map[string]interface{}),
		AST:     program,
	}

	// Cache the module
	mr.cache[resolvedPath] = module

	return module, nil
}

// resolvePath resolves a module path to an actual file path
func (mr *ModuleResolver) resolvePath(modulePath string, baseDir string) (string, error) {
	// Handle standard library modules
	if strings.HasPrefix(modulePath, "std/") {
		return mr.resolveStandardLibraryPath(modulePath)
	}

	// Handle relative paths
	if strings.HasPrefix(modulePath, "./") || strings.HasPrefix(modulePath, "../") {
		absPath := filepath.Join(baseDir, modulePath)
		
		// Add .rush extension if not present
		if !strings.HasSuffix(absPath, ".rush") {
			absPath += ".rush"
		}
		
		return filepath.Abs(absPath)
	}

	// Handle absolute paths
	if strings.HasPrefix(modulePath, "/") {
		absPath := modulePath
		if !strings.HasSuffix(absPath, ".rush") {
			absPath += ".rush"
		}
		return absPath, nil
	}

	// Handle user modules in current directory
	absPath := filepath.Join(baseDir, modulePath)
	if !strings.HasSuffix(absPath, ".rush") {
		absPath += ".rush"
	}
	
	return filepath.Abs(absPath)
}

// resolveStandardLibraryPath resolves a standard library module path
func (mr *ModuleResolver) resolveStandardLibraryPath(modulePath string) (string, error) {
	// Remove the "std/" prefix
	stdModuleName := strings.TrimPrefix(modulePath, "std/")
	
	// Try multiple search paths for the standard library
	searchPaths := []string{}
	
	// First, try relative to current working directory (for development/testing)
	if cwd, err := os.Getwd(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(cwd, "std", stdModuleName))
	}
	
	// Then try relative to executable path (for deployed installations)
	if execPath, err := os.Executable(); err == nil {
		// Resolve symlinks to get the actual executable path
		if realPath, err := filepath.EvalSymlinks(execPath); err == nil {
			execPath = realPath
		}
		execDir := filepath.Dir(execPath)
		searchPaths = append(searchPaths, filepath.Join(execDir, "std", stdModuleName))
	}
	
	// Try each search path
	for _, stdLibPath := range searchPaths {
		// Add .rush extension if not present
		if !strings.HasSuffix(stdLibPath, ".rush") {
			stdLibPath += ".rush"
		}
		
		// Check if the standard library module exists
		if _, err := os.Stat(stdLibPath); err == nil {
			return stdLibPath, nil
		}
	}
	
	return "", fmt.Errorf("standard library module not found: %s", modulePath)
}

// GetExports returns the exports of a loaded module
func (mr *ModuleResolver) GetExports(modulePath string) (map[string]interface{}, error) {
	module, exists := mr.cache[modulePath]
	if !exists {
		return nil, fmt.Errorf("module %s not loaded", modulePath)
	}
	return module.Exports, nil
}

// IsLoaded checks if a module is already loaded
func (mr *ModuleResolver) IsLoaded(modulePath string) bool {
	_, exists := mr.cache[modulePath]
	return exists
}