// Package link implements communication link abstractions from
// Cachin et al., Chapter 2, Section 2.4.
//
// PARALLEL VERSION: Links embed their own mutex-protected RNG
// instead of receiving *rand.Rand per call. This makes them safe
// for concurrent use from multiple goroutines.
//
// Link hierarchy:
//
//	Unreliable -> Fair-Loss -> Stubborn -> Perfect -> Authenticated
package link

import (
	"math/rand"
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// Link is the interface for all communication link abstractions.
//
// PARALLEL VERSION differences:
//   - Send/Receive no longer take *rand.Rand — links own their RNG.
//   - Retransmissions() replaces Tick() — returns buffered messages
//     to retransmit (called by the runtime's retransmit goroutine).
type Link interface {
	// Send processes an outgoing message. Returns messages to enqueue.
	Send(msg process.Message) []process.Message

	// Receive checks if a message should be delivered.
	// Returns (processed_message, should_deliver).
	Receive(msg process.Message) (process.Message, bool)

	// Retransmissions returns any messages that need retransmitting.
	// Called periodically by the runtime's retransmitter goroutine.
	Retransmissions() []process.Message
}

// ============================================================================
// UNRELIABLE LINK (baseline)
// ============================================================================

// UnreliableLink provides NO guarantees. Messages may be lost at any time.
type UnreliableLink struct {
	lossRate float64
	mu       sync.Mutex
	rng      *rand.Rand
}

// NewUnreliableLink creates an unreliable link with the given loss rate and seed.
func NewUnreliableLink(lossRate float64, seed int64) *UnreliableLink {
	return &UnreliableLink{
		lossRate: lossRate,
		rng:      rand.New(rand.NewSource(seed)),
	}
}

// Send may drop the message based on LossRate.
func (l *UnreliableLink) Send(msg process.Message) []process.Message {
	l.mu.Lock()
	lost := l.rng.Float64() < l.lossRate
	l.mu.Unlock()
	if lost {
		return nil
	}
	return []process.Message{msg}
}

// Receive always accepts messages (no filtering).
func (l *UnreliableLink) Receive(msg process.Message) (process.Message, bool) {
	return msg, true
}

// Retransmissions returns nil (no retransmission for unreliable links).
func (l *UnreliableLink) Retransmissions() []process.Message {
	return nil
}
