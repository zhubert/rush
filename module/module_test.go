package module

import (
  "os"
  "path/filepath"
  "strings"
  "testing"
)

func TestNewModuleResolver(t *testing.T) {
  resolver := NewModuleResolver()
  
  if resolver == nil {
    t.Fatal("NewModuleResolver returned nil")
  }
  
  if resolver.cache == nil {
    t.Error("Module resolver cache not initialized")
  }
  
  if resolver.loadStack == nil {
    t.Error("Module resolver loadStack not initialized")
  }
  
  if len(resolver.cache) != 0 {
    t.Error("Initial cache should be empty")
  }
  
  if len(resolver.loadStack) != 0 {
    t.Error("Initial loadStack should be empty")
  }
}

func TestModuleResolverResolvePath(t *testing.T) {
  resolver := NewModuleResolver()
  baseDir := "/tmp/test"
  
  tests := []struct {
    name         string
    modulePath   string
    expectedSuffix string
  }{
    {
      name:         "Relative path with extension",
      modulePath:   "./math_module.rush",
      expectedSuffix: "math_module.rush",
    },
    {
      name:         "Relative path without extension",
      modulePath:   "./math_module",
      expectedSuffix: "math_module.rush",
    },
    {
      name:         "Parent directory path",
      modulePath:   "../utils/helper",
      expectedSuffix: "helper.rush",
    },
    {
      name:         "Absolute path with extension",
      modulePath:   "/usr/local/lib/module.rush",
      expectedSuffix: "module.rush",
    },
    {
      name:         "Absolute path without extension",
      modulePath:   "/usr/local/lib/module",
      expectedSuffix: "module.rush",
    },
  }
  
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      resolved, err := resolver.resolvePath(tt.modulePath, baseDir)
      if err != nil {
        t.Fatalf("resolvePath failed: %v", err)
      }
      
      if !strings.HasSuffix(resolved, tt.expectedSuffix) {
        t.Errorf("Expected path to end with %q, got %q", tt.expectedSuffix, resolved)
      }
      
      if !strings.HasSuffix(resolved, ".rush") {
        t.Errorf("Expected resolved path to have .rush extension, got %q", resolved)
      }
    })
  }
}

func TestModuleResolverLoadModule(t *testing.T) {
  resolver := NewModuleResolver()
  
  // Create a temporary directory and module file
  tmpDir, err := os.MkdirTemp("", "module_test")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpDir)
  
  // Create a test module file
  moduleContent := `export add = fn(x, y) { return x + y }
export PI = 3.14159`
  
  modulePath := filepath.Join(tmpDir, "test_module.rush")
  if err := os.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
    t.Fatal(err)
  }
  
  // Test loading the module
  module, err := resolver.LoadModule("./test_module", tmpDir)
  if err != nil {
    t.Fatalf("LoadModule failed: %v", err)
  }
  
  if module == nil {
    t.Fatal("LoadModule returned nil module")
  }
  
  if module.Path == "" {
    t.Error("Module path should not be empty")
  }
  
  if module.AST == nil {
    t.Error("Module AST should not be nil")
  }
  
  if module.Exports == nil {
    t.Error("Module exports should not be nil")
  }
  
  // Test that module is cached
  if !resolver.IsLoaded(module.Path) {
    t.Error("Module should be cached after loading")
  }
  
  // Test loading the same module again (should return cached version)
  module2, err := resolver.LoadModule("./test_module", tmpDir)
  if err != nil {
    t.Fatalf("Second LoadModule failed: %v", err)
  }
  
  if module != module2 {
    t.Error("Second load should return the same cached module instance")
  }
}

func TestModuleResolverLoadModuleErrors(t *testing.T) {
  resolver := NewModuleResolver()
  
  tests := []struct {
    name           string
    modulePath     string
    baseDir        string
    expectedError  string
  }{
    {
      name:          "Non-existent module",
      modulePath:    "./nonexistent",
      baseDir:       "/tmp",
      expectedError: "failed to read module",
    },
    {
      name:          "Invalid module syntax",
      modulePath:    "./invalid_module",
      baseDir:       "",
      expectedError: "parse errors in module",
    },
  }
  
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // For invalid syntax test, create a file with invalid content
      if tt.name == "Invalid module syntax" {
        tmpDir, err := os.MkdirTemp("", "module_error_test")
        if err != nil {
          t.Fatal(err)
        }
        defer os.RemoveAll(tmpDir)
        
        tt.baseDir = tmpDir
        invalidContent := `export add = fn(x, y { return x + y` // Missing closing parenthesis
        modulePath := filepath.Join(tmpDir, "invalid_module.rush")
        if err := os.WriteFile(modulePath, []byte(invalidContent), 0644); err != nil {
          t.Fatal(err)
        }
      }
      
      module, err := resolver.LoadModule(tt.modulePath, tt.baseDir)
      if err == nil {
        t.Errorf("Expected error containing %q, but got no error", tt.expectedError)
        return
      }
      
      if !strings.Contains(err.Error(), tt.expectedError) {
        t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
      }
      
      if module != nil {
        t.Error("Expected nil module on error")
      }
    })
  }
}

func TestModuleResolverCircularDependency(t *testing.T) {
  resolver := NewModuleResolver()
  
  tmpDir, err := os.MkdirTemp("", "circular_test")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpDir)
  
  // Create two modules that import each other
  module1Content := `import { funcB } from "./module2"
export funcA = fn() { return "A" }`
  
  module2Content := `import { funcA } from "./module1"
export funcB = fn() { return "B" }`
  
  module1Path := filepath.Join(tmpDir, "module1.rush")
  module2Path := filepath.Join(tmpDir, "module2.rush")
  
  if err := os.WriteFile(module1Path, []byte(module1Content), 0644); err != nil {
    t.Fatal(err)
  }
  if err := os.WriteFile(module2Path, []byte(module2Content), 0644); err != nil {
    t.Fatal(err)
  }
  
  // This test demonstrates the circular dependency detection structure
  // Note: Actual circular dependency detection happens during interpretation
  // For now, we test that the resolver can load modules with circular imports in the AST
  module1, err := resolver.LoadModule("./module1", tmpDir)
  if err != nil {
    t.Fatalf("Failed to load module1: %v", err)
  }
  
  if module1 == nil {
    t.Error("Module1 should not be nil")
  }
  
  // Test that both modules can be parsed (circular dependency detection
  // happens at interpretation time, not at module loading time)
  module2, err := resolver.LoadModule("./module2", tmpDir)
  if err != nil {
    t.Fatalf("Failed to load module2: %v", err)
  }
  
  if module2 == nil {
    t.Error("Module2 should not be nil")
  }
}

func TestModuleResolverGetExports(t *testing.T) {
  resolver := NewModuleResolver()
  
  tmpDir, err := os.MkdirTemp("", "exports_test")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpDir)
  
  moduleContent := `export testFunc = fn() { return 42 }
export testVar = "hello"`
  
  modulePath := filepath.Join(tmpDir, "exports_module.rush")
  if err := os.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
    t.Fatal(err)
  }
  
  // Load the module
  module, err := resolver.LoadModule("./exports_module", tmpDir)
  if err != nil {
    t.Fatalf("LoadModule failed: %v", err)
  }
  
  // Test getting exports
  exports, err := resolver.GetExports(module.Path)
  if err != nil {
    t.Fatalf("GetExports failed: %v", err)
  }
  
  if exports == nil {
    t.Error("Exports should not be nil")
  }
  
  // Test getting exports for non-existent module
  _, err = resolver.GetExports("/nonexistent/path")
  if err == nil {
    t.Error("Expected error for non-existent module")
  }
  
  if !strings.Contains(err.Error(), "not loaded") {
    t.Errorf("Expected 'not loaded' error, got %q", err.Error())
  }
}

func TestModuleResolverIsLoaded(t *testing.T) {
  resolver := NewModuleResolver()
  
  tmpDir, err := os.MkdirTemp("", "loaded_test")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(tmpDir)
  
  moduleContent := `export value = 42`
  modulePath := filepath.Join(tmpDir, "loaded_module.rush")
  if err := os.WriteFile(modulePath, []byte(moduleContent), 0644); err != nil {
    t.Fatal(err)
  }
  
  // Test that module is not loaded initially
  resolvedPath, _ := resolver.resolvePath("./loaded_module", tmpDir)
  if resolver.IsLoaded(resolvedPath) {
    t.Error("Module should not be loaded initially")
  }
  
  // Load the module
  module, err := resolver.LoadModule("./loaded_module", tmpDir)
  if err != nil {
    t.Fatalf("LoadModule failed: %v", err)
  }
  
  // Test that module is now loaded
  if !resolver.IsLoaded(module.Path) {
    t.Error("Module should be loaded after LoadModule")
  }
  
  // Test with non-existent path
  if resolver.IsLoaded("/nonexistent/path") {
    t.Error("Non-existent module should not be loaded")
  }
}