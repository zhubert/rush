package aot

import (
	"fmt"
	"rush/compiler"
)

// StaticAnalyzer performs whole-program analysis and optimization
type StaticAnalyzer struct {
	passes []OptimizationPass
}

// OptimizationPass represents a single optimization transformation
type OptimizationPass interface {
	Name() string
	Apply(code *compiler.Bytecode) (*compiler.Bytecode, error)
	IsEnabled(level int) bool
}

// AnalysisResult contains the results of static analysis
type AnalysisResult struct {
	CallGraph        *CallGraph
	DeadCode         []int          // Dead instruction indices
	HotPaths         [][]int        // Hot execution paths
	ConstantValues   map[int]interface{} // Constant propagation results
	InlineCandidates []FunctionInfo // Functions suitable for inlining
}

// CallGraph represents the program's function call relationships
type CallGraph struct {
	Functions []FunctionInfo
	Calls     map[int][]int // Function index -> called function indices
}

// FunctionInfo contains metadata about a function
type FunctionInfo struct {
	StartIndex    int
	EndIndex      int
	CallCount     int64
	Size          int
	IsRecursive   bool
	CanInline     bool
}

// NewStaticAnalyzer creates a new static analyzer with optimization passes
func NewStaticAnalyzer() *StaticAnalyzer {
	return &StaticAnalyzer{
		passes: []OptimizationPass{
			NewDeadCodeEliminationPass(),
			NewConstantPropagationPass(),
			NewFunctionInliningPass(),
			NewTailCallOptimizationPass(),
			NewLoopOptimizationPass(),
		},
	}
}

// AnalyzeAndOptimize performs complete static analysis and optimization
func (a *StaticAnalyzer) AnalyzeAndOptimize(code *compiler.Bytecode) (*compiler.Bytecode, error) {
	// Perform static analysis (results currently unused in simplified implementation)
	_, err := a.analyzeProgram(code)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Apply optimization passes
	optimized := code
	for _, pass := range a.passes {
		if pass.IsEnabled(1) { // TODO: Use actual optimization level
			result, err := pass.Apply(optimized)
			if err != nil {
				return nil, fmt.Errorf("optimization pass %s failed: %w", pass.Name(), err)
			}
			optimized = result
		}
	}

	return optimized, nil
}

// analyzeProgram performs static analysis to gather optimization information
func (a *StaticAnalyzer) analyzeProgram(code *compiler.Bytecode) (*AnalysisResult, error) {
	result := &AnalysisResult{
		CallGraph:        a.buildCallGraph(code),
		DeadCode:         a.findDeadCode(code),
		HotPaths:         a.identifyHotPaths(code),
		ConstantValues:   a.findConstants(code),
		InlineCandidates: a.findInlineCandidates(code),
	}
	return result, nil
}

// buildCallGraph constructs the program's call graph
func (a *StaticAnalyzer) buildCallGraph(code *compiler.Bytecode) *CallGraph {
	graph := &CallGraph{
		Functions: make([]FunctionInfo, 0),
		Calls:     make(map[int][]int),
	}

	// Simplified implementation - treat entire program as one function
	if len(code.Instructions) > 0 {
		mainFunc := FunctionInfo{
			StartIndex: 0,
			EndIndex:   len(code.Instructions) - 1,
			Size:       len(code.Instructions),
			CallCount:  1,
			IsRecursive: false,
			CanInline:   false,
		}
		graph.Functions = append(graph.Functions, mainFunc)
	}

	return graph
}

// findDeadCode identifies unreachable instructions
func (a *StaticAnalyzer) findDeadCode(code *compiler.Bytecode) []int {
	// Simplified implementation - assume all instructions are reachable
	// In a full implementation, this would perform proper reachability analysis
	return []int{}
}

// markReachable recursively marks reachable instructions (unused in simplified version)
func (a *StaticAnalyzer) markReachable(code *compiler.Bytecode, index int, reachable []bool) {
	// Simplified implementation - not used
}

// identifyHotPaths finds frequently executed code paths
func (a *StaticAnalyzer) identifyHotPaths(code *compiler.Bytecode) [][]int {
	// Simplified implementation - return empty hot paths
	return [][]int{}
}

// findConstants identifies constant values through the program
func (a *StaticAnalyzer) findConstants(code *compiler.Bytecode) map[int]interface{} {
	// Simplified implementation - return constants from constant pool
	constants := make(map[int]interface{})
	for i, constant := range code.Constants {
		constants[i] = constant
	}
	return constants
}

// findInlineCandidates identifies functions suitable for inlining
func (a *StaticAnalyzer) findInlineCandidates(code *compiler.Bytecode) []FunctionInfo {
	candidates := make([]FunctionInfo, 0)
	
	// Simple heuristic: small functions without recursion
	callGraph := a.buildCallGraph(code)
	for i, fn := range callGraph.Functions {
		if fn.Size <= 10 && !fn.IsRecursive { // Small, non-recursive functions
			fn.CanInline = true
			candidates = append(candidates, fn)
		}
		_ = i // Suppress unused variable warning
	}
	
	return candidates
}