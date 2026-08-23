package link

import (
	"fmt"
	"sync"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// PERFECT LINK (PL)
// ============================================================================

// PerfectLink models the Perfect Link abstraction.
//
// From Cachin et al., Section 2.4.3:
//
//	PL1 (Reliable delivery): If p sends m to q, q eventually delivers m.
//	PL2 (No duplication): No message is delivered more than once.
//	PL3 (No creation): Delivered messages were previously sent.
//
// Implementation: Perfect = Stubborn + Deduplication.
//
// PARALLEL VERSION: Uses sync.Mutex to protect the delivered set.
type PerfectLink struct {
	underlying *StubbornLink
	mu         sync.Mutex
	delivered  map[string]bool
}

// NewPerfectLink creates a perfect link wrapping a stubborn link.
func NewPerfectLink(sl *StubbornLink) *PerfectLink {
	return &PerfectLink{
		underlying: sl,
		delivered:  make(map[string]bool),
	}
}

// Send forwards the message through the stubborn link.
func (l *PerfectLink) Send(msg process.Message) []process.Message {
	return l.underlying.Send(msg)
}

// Receive implements deduplication.
// Only delivers messages that haven't been seen before.
func (l *PerfectLink) Receive(msg process.Message) (process.Message, bool) {
	key := messageKey(msg)
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.delivered[key] {
		return msg, false // Already delivered — suppress duplicate
	}
	l.delivered[key] = true
	return msg, true
}

// Retransmissions delegates to the stubborn link.
func (l *PerfectLink) Retransmissions() []process.Message {
	return l.underlying.Retransmissions()
}

// messageKey creates a unique identifier for deduplication.
func messageKey(msg process.Message) string {
	return fmt.Sprintf("%s:%s:%s:%v", msg.From, msg.To, msg.Type, msg.Data)
}
