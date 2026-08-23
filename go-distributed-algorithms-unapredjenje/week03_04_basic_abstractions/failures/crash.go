// Package failures implements process failure models from
// Cachin et al., Chapter 2, Section 2.2.
//
// In a distributed system, processes may FAIL in various ways.
// The failure model determines what kinds of faults the system
// must tolerate, and directly affects which algorithms are possible.
//
// Failure models (from weakest to strongest):
//
//  0. NO FAILURE: All processes are correct.
//  1. CRASH-STOP: A process halts and never recovers.
//  2. OMISSION: A process may silently drop messages.
//  3. CRASH-RECOVERY: A process crashes but may restart (with state loss) which is not implemented in this simulator.
//  4. BYZANTINE: A process may behave arbitrarily (maliciously).
//
// Each model is implemented as a FailureInjector that the runtime uses
// to simulate faults during execution.
package failures

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// FailureInjector controls process failures in the simulation.
// The runtime calls these methods during message processing to
// determine if processes are alive, if messages should be delivered,
// and if outgoing messages should be altered.
type FailureInjector interface {
	// IsAlive returns true if the process is currently operational.
	// A crashed process will not receive any messages.
	IsAlive(pid process.ProcessID) bool

	// ShouldDeliver returns true if a message should be delivered.
	// Omission failures may cause this to return false.
	ShouldDeliver(msg process.Message) bool

	// MaybeAlter allows failure models to modify outgoing messages.
	// Byzantine processes may produce arbitrary messages.
	// Non-Byzantine models return messages unmodified.
	MaybeAlter(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message

	// CrashProcess marks a process as crashed.
	CrashProcess(pid process.ProcessID)

	// RecoverProcess marks a process as recovered (crash-recovery model).
	RecoverProcess(pid process.ProcessID)
}

// NoFailure is the default failure model: all processes are correct.
// No crashes, no omissions, no Byzantine behavior.
// Use this for examples that focus on link abstractions or basic algorithms.
type NoFailure struct{}

// NewNoFailure creates a no-failure model.
func NewNoFailure() *NoFailure {
	return &NoFailure{}
}

func (f *NoFailure) IsAlive(pid process.ProcessID) bool     { return true }
func (f *NoFailure) ShouldDeliver(msg process.Message) bool { return true }
func (f *NoFailure) CrashProcess(pid process.ProcessID)     {}
func (f *NoFailure) RecoverProcess(pid process.ProcessID)   {}

func (f *NoFailure) MaybeAlter(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message {
	return msgs
}

// CrashFailure models the crash-stop failure model.
//
// From Cachin et al., Section 2.2.1:
//
//	"A process that crashes simply stops executing steps of its algorithm.
//	 It does not perform any other (harmful) action, such as sending
//	 corrupted messages."
//
// Properties:
//   - A crashed process STOPS permanently (no recovery).
//   - A crashed process does not send or receive messages.
//
// This is the simplest and most commonly assumed failure model.
type CrashFailure struct {
	crashed map[process.ProcessID]bool
}

// NewCrashFailure creates a crash-stop failure model.
func NewCrashFailure() *CrashFailure {
	return &CrashFailure{
		crashed: make(map[process.ProcessID]bool),
	}
}

// IsAlive returns false if the process has crashed.
func (f *CrashFailure) IsAlive(pid process.ProcessID) bool {
	return !f.crashed[pid]
}

// ShouldDeliver returns true. Crashed processes are filtered by IsAlive.
func (f *CrashFailure) ShouldDeliver(msg process.Message) bool {
	return true
}

// CrashProcess permanently stops a process.
func (f *CrashFailure) CrashProcess(pid process.ProcessID) {
	f.crashed[pid] = true
}

// RecoverProcess does nothing in crash-STOP model (no recovery).
func (f *CrashFailure) RecoverProcess(pid process.ProcessID) {
	// In crash-stop model, processes do NOT recover.
}

// MaybeAlter returns messages unmodified. Crashed processes cannot send.
func (f *CrashFailure) MaybeAlter(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message {
	return msgs
}
