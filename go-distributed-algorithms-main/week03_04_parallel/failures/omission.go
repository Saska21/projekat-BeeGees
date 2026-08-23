package failures

import (
	"math/rand"
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// OMISSION FAILURE
// ============================================================================

// OmissionFailure models the omission failure model.
//
// An omission fault means a process fails to receive a message it should
// have. The process itself continues running (unlike crash), but some
// messages are silently lost.
//
// PARALLEL VERSION: Embeds its own mutex-protected RNG.
type OmissionFailure struct {
	mu       sync.RWMutex
	omitRate map[process.ProcessID]float64
	rng      *rand.Rand
}

// NewOmissionFailure creates an omission failure model.
func NewOmissionFailure(seed int64) *OmissionFailure {
	return &OmissionFailure{
		omitRate: make(map[process.ProcessID]float64),
		rng:      rand.New(rand.NewSource(seed)),
	}
}

// SetOmissionRate configures the omission probability for a process.
func (f *OmissionFailure) SetOmissionRate(pid process.ProcessID, rate float64) {
	f.mu.Lock()
	f.omitRate[pid] = rate
	f.mu.Unlock()
}

func (f *OmissionFailure) IsAlive(pid process.ProcessID) bool { return true }

// ShouldDeliver returns false with probability equal to the process's omission rate.
func (f *OmissionFailure) ShouldDeliver(msg process.Message) bool {
	f.mu.RLock()
	rate, exists := f.omitRate[msg.To]
	f.mu.RUnlock()
	if !exists {
		return true
	}
	f.mu.Lock()
	val := f.rng.Float64()
	f.mu.Unlock()
	return val >= rate
}

func (f *OmissionFailure) CrashProcess(pid process.ProcessID)   {}
func (f *OmissionFailure) RecoverProcess(pid process.ProcessID) {}

func (f *OmissionFailure) MaybeAlter(_ process.ProcessID, msgs []process.Message) []process.Message {
	return msgs
}
