package link

import (
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// STUBBORN LINK (SL)
// ============================================================================

// StubbornLink models the Stubborn Link abstraction.
//
// From Cachin et al., Section 2.4.2:
//
//	SL1 (Stubborn delivery): If a correct process p sends a message m
//	  once to a correct process q, then q delivers m an infinite number
//	  of times.
//	SL2 (No creation): Delivered messages were previously sent.
//
// Implementation: Uses a Fair-Loss Link underneath. Stores every sent
// message. Periodically re-sends ALL stored messages through the FLL.
//
// PARALLEL VERSION: Uses sync.Mutex to protect the sent buffer.
// Retransmissions() is called by the runtime's retransmitter goroutine.
type StubbornLink struct {
	underlying *FairLossLink
	mu         sync.Mutex
	sent       []process.Message
}

// NewStubbornLink creates a stubborn link wrapping a fair-loss link.
func NewStubbornLink(fll *FairLossLink) *StubbornLink {
	return &StubbornLink{
		underlying: fll,
		sent:       make([]process.Message, 0),
	}
}

// Send stores the message and forwards it through the fair-loss link.
func (l *StubbornLink) Send(msg process.Message) []process.Message {
	l.mu.Lock()
	l.sent = append(l.sent, msg.Clone())
	l.mu.Unlock()
	return l.underlying.Send(msg)
}

// Receive always accepts (stubborn link has no deduplication).
// Duplicates ARE expected — that's the price of stubborn delivery.
func (l *StubbornLink) Receive(msg process.Message) (process.Message, bool) {
	return msg, true
}

// Retransmissions takes a snapshot of all stored messages under the lock,
// releases the lock, then sends each through the underlying fair-loss link.
func (l *StubbornLink) Retransmissions() []process.Message {
	l.mu.Lock()
	snapshot := make([]process.Message, len(l.sent))
	for i, m := range l.sent {
		snapshot[i] = m.Clone()
	}
	l.mu.Unlock()

	var retransmissions []process.Message
	for _, msg := range snapshot {
		delivered := l.underlying.Send(msg)
		retransmissions = append(retransmissions, delivered...)
	}
	return retransmissions
}
