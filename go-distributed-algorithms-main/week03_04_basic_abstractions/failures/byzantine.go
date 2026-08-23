package failures

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ByzantineFailure models Byzantine (arbitrary) process failures.
//
// From Cachin et al., Section 2.2:
//
//	"A Byzantine process may deviate ARBITRARILY from its algorithm.
//	 It may crash, omit messages, send wrong messages, or even
//	 coordinate with other Byzantine processes to subvert the system."
//
// This is the STRONGEST failure model. A Byzantine process can:
//   - Send messages it shouldn't send
//   - Modify message content (forge messages)
//   - Impersonate other processes (without authentication)
//   - Stop processing messages (like a crash)
//   - Do all of the above non-deterministically
//
// In our simulator, Byzantine behavior is modelled through an
// "alteration function" that can modify outgoing messages from
// Byzantine processes.
type ByzantineFailure struct {
	byzantineProcesses map[process.ProcessID]bool
	alterFunc          func(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message
}

// NewByzantineFailure creates a Byzantine failure model.
func NewByzantineFailure() *ByzantineFailure {
	return &ByzantineFailure{
		byzantineProcesses: make(map[process.ProcessID]bool),
	}
}

// SetByzantine marks a process as Byzantine.
func (f *ByzantineFailure) SetByzantine(pid process.ProcessID) {
	f.byzantineProcesses[pid] = true
}

// SetAlterFunc configures what Byzantine processes do to outgoing messages.
func (f *ByzantineFailure) SetAlterFunc(fn func(process.ProcessID, []process.Message, *rand.Rand) []process.Message) {
	f.alterFunc = fn
}

func (f *ByzantineFailure) IsAlive(pid process.ProcessID) bool     { return true }
func (f *ByzantineFailure) ShouldDeliver(msg process.Message) bool { return true }
func (f *ByzantineFailure) CrashProcess(pid process.ProcessID)     {}
func (f *ByzantineFailure) RecoverProcess(pid process.ProcessID)   {}

// MaybeAlter applies the alteration function to outgoing messages
// from Byzantine processes. Non-Byzantine messages pass through unmodified.
func (f *ByzantineFailure) MaybeAlter(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message {
	if !f.byzantineProcesses[from] {
		return msgs
	}
	if f.alterFunc == nil {
		return msgs
	}
	return f.alterFunc(from, msgs, rng)
}
