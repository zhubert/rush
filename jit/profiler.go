package jit

import (
	"sync"
	"time"
)

// ExecutionProfiler tracks function execution patterns for JIT compilation decisions
type ExecutionProfiler struct {
	functions map[uint64]*FunctionProfile
	mu        sync.RWMutex
}

// FunctionProfile holds execution statistics for a single function
type FunctionProfile struct {
	Hash            uint64
	ExecutionCount  int64
	TotalTime       time.Duration
	AverageTime     time.Duration
	LastExecution   time.Time
	FirstExecution  time.Time
	IsHot           bool
}

// NewExecutionProfiler creates a new execution profiler
func NewExecutionProfiler() *ExecutionProfiler {
	return &ExecutionProfiler{
		functions: make(map[uint64]*FunctionProfile),
	}
}

// RecordExecution records a function execution for profiling
func (p *ExecutionProfiler) RecordExecution(fnHash uint64, executionTime time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	profile, exists := p.functions[fnHash]
	if !exists {
		profile = &FunctionProfile{
			Hash:           fnHash,
			ExecutionCount: 0,
			TotalTime:      0,
			FirstExecution: time.Now(),
		}
		p.functions[fnHash] = profile
	}
	
	profile.ExecutionCount++
	profile.TotalTime += executionTime
	profile.AverageTime = time.Duration(int64(profile.TotalTime) / profile.ExecutionCount)
	profile.LastExecution = time.Now()
	
	// Mark as hot if execution count exceeds threshold
	if profile.ExecutionCount >= DefaultHotThreshold {
		profile.IsHot = true
	}
}

// GetExecutionCount returns the execution count for a function
func (p *ExecutionProfiler) GetExecutionCount(fnHash uint64) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if profile, exists := p.functions[fnHash]; exists {
		return int(profile.ExecutionCount)
	}
	return 0
}

// GetProfile returns the execution profile for a function
func (p *ExecutionProfiler) GetProfile(fnHash uint64) *FunctionProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if profile, exists := p.functions[fnHash]; exists {
		// Return a copy to avoid race conditions
		return &FunctionProfile{
			Hash:           profile.Hash,
			ExecutionCount: profile.ExecutionCount,
			TotalTime:      profile.TotalTime,
			AverageTime:    profile.AverageTime,
			LastExecution:  profile.LastExecution,
			FirstExecution: profile.FirstExecution,
			IsHot:          profile.IsHot,
		}
	}
	return nil
}

// GetHotFunctions returns all functions marked as hot
func (p *ExecutionProfiler) GetHotFunctions() []*FunctionProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	var hotFunctions []*FunctionProfile
	for _, profile := range p.functions {
		if profile.IsHot {
			hotFunctions = append(hotFunctions, &FunctionProfile{
				Hash:           profile.Hash,
				ExecutionCount: profile.ExecutionCount,
				TotalTime:      profile.TotalTime,
				AverageTime:    profile.AverageTime,
				LastExecution:  profile.LastExecution,
				FirstExecution: profile.FirstExecution,
				IsHot:          profile.IsHot,
			})
		}
	}
	return hotFunctions
}

// GetAllProfiles returns all function profiles
func (p *ExecutionProfiler) GetAllProfiles() map[uint64]*FunctionProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	profiles := make(map[uint64]*FunctionProfile)
	for hash, profile := range p.functions {
		profiles[hash] = &FunctionProfile{
			Hash:           profile.Hash,
			ExecutionCount: profile.ExecutionCount,
			TotalTime:      profile.TotalTime,
			AverageTime:    profile.AverageTime,
			LastExecution:  profile.LastExecution,
			FirstExecution: profile.FirstExecution,
			IsHot:          profile.IsHot,
		}
	}
	return profiles
}

// Reset clears all profiling data
func (p *ExecutionProfiler) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.functions = make(map[uint64]*FunctionProfile)
}

// GetStats returns overall profiling statistics
func (p *ExecutionProfiler) GetStats() ProfilerStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	stats := ProfilerStats{
		TotalFunctions: len(p.functions),
		HotFunctions:   0,
		TotalExecutions: 0,
		TotalTime:      0,
	}
	
	for _, profile := range p.functions {
		stats.TotalExecutions += profile.ExecutionCount
		stats.TotalTime += profile.TotalTime
		if profile.IsHot {
			stats.HotFunctions++
		}
	}
	
	if stats.TotalExecutions > 0 {
		stats.AverageExecutionTime = time.Duration(int64(stats.TotalTime) / stats.TotalExecutions)
	}
	
	return stats
}

// ProfilerStats holds overall profiling statistics
type ProfilerStats struct {
	TotalFunctions       int
	HotFunctions         int
	TotalExecutions      int64
	TotalTime            time.Duration
	AverageExecutionTime time.Duration
}