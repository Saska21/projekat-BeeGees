package failures

import (
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// BYZANTINE FAILURE
// ============================================================================

// ByzantineFailure models Byzantine (arbitrary) process failures.
//
// From Cachin et al., Section 2.2:
//
//	A Byzantine process may deviate ARBITRARILY from its algorithm.
//
// PARALLEL VERSION: Uses sync.RWMutex for concurrent access.
// The alterFunc no longer takes *rand.Rand.
type ByzantineFailure struct {
	mu                 sync.RWMutex
	byzantineProcesses map[process.ProcessID]bool
	alterFunc          func(from process.ProcessID, msgs []process.Message) []process.Message
}

func NewByzantineFailure() *ByzantineFailure {
	return &ByzantineFailure{
		byzantineProcesses: make(map[process.ProcessID]bool),
	}
}

// SetByzantine marks a process as Byzantine.
func (f *ByzantineFailure) SetByzantine(pid process.ProcessID) {
	f.mu.Lock()
	f.byzantineProcesses[pid] = true
	f.mu.Unlock()
}

// SetAlterFunc configures what Byzantine processes do to outgoing messages.
func (f *ByzantineFailure) SetAlterFunc(fn func(process.ProcessID, []process.Message) []process.Message) {
	f.mu.Lock()
	f.alterFunc = fn
	f.mu.Unlock()
}

func (f *ByzantineFailure) IsAlive(pid process.ProcessID) bool     { return true }
func (f *ByzantineFailure) ShouldDeliver(msg process.Message) bool { return true }
func (f *ByzantineFailure) CrashProcess(pid process.ProcessID)     {}
func (f *ByzantineFailure) RecoverProcess(pid process.ProcessID)   {}

// MaybeAlter applies the alteration function to outgoing messages
// from Byzantine processes.
func (f *ByzantineFailure) MaybeAlter(from process.ProcessID, msgs []process.Message) []process.Message {
	f.mu.RLock()
	isByz := f.byzantineProcesses[from]
	fn := f.alterFunc
	f.mu.RUnlock()
	if !isByz || fn == nil {
		return msgs
	}
	return fn(from, msgs)
}
