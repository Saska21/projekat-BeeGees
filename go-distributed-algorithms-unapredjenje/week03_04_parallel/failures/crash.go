// Package failures implements process failure models from
// Cachin et al., Chapter 2, Section 2.2.
//
// PARALLEL VERSION: All failure injectors use sync.RWMutex for
// thread-safe access from multiple goroutines. MaybeAlter no longer
// takes *rand.Rand — implementations embed their own RNG if needed.
package failures

import (
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// FailureInjector controls process failures in the simulation.
type FailureInjector interface {
	// IsAlive returns true if the process is currently operational.
	IsAlive(pid process.ProcessID) bool

	// ShouldDeliver returns true if a message should be delivered.
	ShouldDeliver(msg process.Message) bool

	// MaybeAlter allows failure models to modify outgoing messages.
	// PARALLEL VERSION: no *rand.Rand parameter.
	MaybeAlter(from process.ProcessID, msgs []process.Message) []process.Message

	// CrashProcess marks a process as crashed.
	CrashProcess(pid process.ProcessID)

	// RecoverProcess marks a process as recovered.
	RecoverProcess(pid process.ProcessID)
}

// ============================================================================
// NO FAILURE
// ============================================================================

// NoFailure is the default failure model: all processes are correct.
type NoFailure struct{}

func NewNoFailure() *NoFailure { return &NoFailure{} }

func (f *NoFailure) IsAlive(pid process.ProcessID) bool     { return true }
func (f *NoFailure) ShouldDeliver(msg process.Message) bool { return true }
func (f *NoFailure) CrashProcess(pid process.ProcessID)     {}
func (f *NoFailure) RecoverProcess(pid process.ProcessID)   {}
func (f *NoFailure) MaybeAlter(_ process.ProcessID, msgs []process.Message) []process.Message {
	return msgs
}

// ============================================================================
// CRASH-STOP FAILURE
// ============================================================================

// CrashFailure models the crash-stop failure model.
//
// From Cachin et al., Section 2.2.1:
//
//	A crashed process STOPS permanently (no recovery).
//	A crashed process does not send or receive messages.
//
// PARALLEL VERSION: Uses sync.RWMutex for concurrent reads.
type CrashFailure struct {
	mu      sync.RWMutex
	crashed map[process.ProcessID]bool
}

func NewCrashFailure() *CrashFailure {
	return &CrashFailure{
		crashed: make(map[process.ProcessID]bool),
	}
}

func (f *CrashFailure) IsAlive(pid process.ProcessID) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return !f.crashed[pid]
}

func (f *CrashFailure) ShouldDeliver(msg process.Message) bool { return true }

func (f *CrashFailure) CrashProcess(pid process.ProcessID) {
	f.mu.Lock()
	f.crashed[pid] = true
	f.mu.Unlock()
}

func (f *CrashFailure) RecoverProcess(pid process.ProcessID) {
	// In crash-STOP model, processes do NOT recover.
}

func (f *CrashFailure) MaybeAlter(_ process.ProcessID, msgs []process.Message) []process.Message {
	return msgs
}
