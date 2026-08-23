package link

import (
	"math/rand"
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// FAIR-LOSS LINK (FLL)
// ============================================================================

// FairLossLink models the Fair-Loss Link abstraction.
//
// From Cachin et al., Section 2.4.1:
//
//	FLL1 (Fair-loss delivery): If a correct process p infinitely often
//	  sends a message m to a correct process q, then q delivers m an
//	  infinite number of times.
//	FLL2 (Finite duplication): If p sends m a finite number of times,
//	  m cannot be delivered infinitely often by q.
//	FLL3 (No creation): Delivered messages were previously sent.
//
// PARALLEL VERSION: Embeds its own mutex-protected RNG.
type FairLossLink struct {
	lossRate float64
	mu       sync.Mutex
	rng      *rand.Rand
}

// NewFairLossLink creates a fair-loss link with the given loss rate and seed.
func NewFairLossLink(lossRate float64, seed int64) *FairLossLink {
	return &FairLossLink{
		lossRate: lossRate,
		rng:      rand.New(rand.NewSource(seed)),
	}
}

// Send may drop the message based on LossRate.
func (l *FairLossLink) Send(msg process.Message) []process.Message {
	l.mu.Lock()
	lost := l.rng.Float64() < l.lossRate
	l.mu.Unlock()
	if lost {
		return nil
	}
	return []process.Message{msg}
}

// Receive always accepts (fair-loss has no deduplication).
func (l *FairLossLink) Receive(msg process.Message) (process.Message, bool) {
	return msg, true
}

// Retransmissions returns nil (fair-loss has no automatic retransmission).
func (l *FairLossLink) Retransmissions() []process.Message {
	return nil
}
