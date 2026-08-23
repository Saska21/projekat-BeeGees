package link

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ============================================================================
// STUBBORN LINK (SL)
// ============================================================================

// StubbornLink models the Stubborn Link abstraction.
//
// From Cachin et al., Section 2.4.2:
//
// Properties:
//
//	SL1 (Stubborn delivery): If a correct process p sends a message m
//	  once to a correct process q, then q delivers m an infinite number
//	  of times.
//	SL2 (No creation): If some process q delivers a message m with
//	  sender p, then m was previously sent to q by process p.
//
// Implementation (Algorithm 2.1):
//
//	Uses a Fair-Loss Link underneath.
//	Stores every sent message.
//	On each tick, re-sends ALL stored messages through the fair-loss link.
//
// This is the key insight: stubborn delivery is achieved by
// REPEATED RETRANSMISSION over a fair-loss channel.
type StubbornLink struct {
	// underlying is the fair-loss link used for actual transmission.
	underlying *FairLossLink

	// sent stores all messages that have been sent (for retransmission).
	// In a real system, this would be bounded or use acknowledgments.
	sent []process.Message

	// RetransmitInterval controls how often retransmissions occur.
	// Retransmit every N steps (1 = every step).
	RetransmitInterval int
}

// NewStubbornLink creates a stubborn link wrapping a fair-loss link.
func NewStubbornLink(fll *FairLossLink, retransmitInterval int) *StubbornLink {
	if retransmitInterval <= 0 {
		retransmitInterval = 3 // Default: retransmit every 3 steps
	}
	return &StubbornLink{
		underlying:         fll,
		sent:               make([]process.Message, 0),
		RetransmitInterval: retransmitInterval,
	}
}

// Send stores the message and forwards it through the fair-loss link.
func (l *StubbornLink) Send(msg process.Message, rng *rand.Rand) []process.Message {
	// Store for retransmission (stubborn delivery guarantee)
	l.sent = append(l.sent, msg.Clone())
	// Forward through fair-loss link (may be lost)
	return l.underlying.Send(msg, rng)
}

// Tick retransmits ALL stored messages periodically.
// This is the mechanism that achieves stubborn delivery:
// even if individual sends fail (fair-loss), we keep trying.
func (l *StubbornLink) Tick(step int, rng *rand.Rand) []process.Message {
	if step%l.RetransmitInterval != 0 {
		return nil
	}
	var retransmissions []process.Message
	for _, msg := range l.sent {
		// Each retransmission still goes through fair-loss
		delivered := l.underlying.Send(msg.Clone(), rng)
		retransmissions = append(retransmissions, delivered...)
	}
	return retransmissions
}

// Receive always accepts (stubborn link has no deduplication).
// Duplicates ARE expected — that's the price of stubborn delivery.
func (l *StubbornLink) Receive(msg process.Message) (process.Message, bool) {
	return msg, true
}
