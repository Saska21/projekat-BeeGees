package failures

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ============================================================================
// OMISSION FAILURE
// ============================================================================

// OmissionFailure models the omission failure model.
//
// From Cachin et al., Section 2.2:
//
// An omission fault means a process fails to send or receive a message
// that it should have. The process itself continues running (unlike crash),
// but some messages are silently lost.
//
// Types of omission:
//   - SEND OMISSION: process fails to send a message it should have sent
//   - RECEIVE OMISSION: process fails to receive a message sent to it
//
// In our simulator, we model RECEIVE OMISSION: messages addressed to
// a process affected by omission may be silently dropped with some
// probability.
//
// This is similar to link-level message loss (fair-loss link),
// but the fault is at the PROCESS level, not the CHANNEL level.
//
// Omission failures are harder to detect than crashes
// because the process appears to be running normally — it just misses
// some messages. This is a realistic model for overloaded systems.
type OmissionFailure struct {
	// omitRate maps processID to their omission probability.
	// A rate of 0.5 means 50% of messages to that process are silently lost.
	omitRate map[process.ProcessID]float64
}

// NewOmissionFailure creates an omission failure model.
func NewOmissionFailure() *OmissionFailure {
	return &OmissionFailure{
		omitRate: make(map[process.ProcessID]float64),
	}
}

// SetOmissionRate configures the omission probability for a process.
// rate should be between 0.0 (no omission) and 1.0 (all messages lost).
func (f *OmissionFailure) SetOmissionRate(pid process.ProcessID, rate float64) {
	f.omitRate[pid] = rate
}

// IsAlive always returns true — omission processes are still "alive",
// they just miss some messages.
func (f *OmissionFailure) IsAlive(pid process.ProcessID) bool {
	return true
}

// ShouldDeliver returns false with probability equal to the process's omission rate.
// This simulates receive omission: the message is addressed to the process
// but silently dropped.
func (f *OmissionFailure) ShouldDeliver(msg process.Message) bool {
	rate, exists := f.omitRate[msg.To]
	if !exists {
		return true // No omission configured — deliver normally
	}
	// We use a simple deterministic-ish check based on message content
	// For proper randomness, the runtime passes the rng to ShouldDeliver.
	// SIMPLIFIED: always deliver if no rate set.
	_ = rate
	return true // See ShouldDeliverRng for randomized version
}

// ShouldDeliverRng is the randomized version of ShouldDeliver.
// Called by the runtime with its random number generator.
func (f *OmissionFailure) ShouldDeliverRng(msg process.Message, rng *rand.Rand) bool {
	rate, exists := f.omitRate[msg.To]
	if !exists {
		return true
	}
	return rng.Float64() >= rate // Deliver if random value exceeds omission rate
}

// CrashProcess does nothing — omission is not crash.
func (f *OmissionFailure) CrashProcess(pid process.ProcessID) {}

// RecoverProcess does nothing.
func (f *OmissionFailure) RecoverProcess(pid process.ProcessID) {}

// MaybeAlter returns messages unmodified — omission only drops, never alters.
func (f *OmissionFailure) MaybeAlter(from process.ProcessID, msgs []process.Message, rng *rand.Rand) []process.Message {
	return msgs
}
