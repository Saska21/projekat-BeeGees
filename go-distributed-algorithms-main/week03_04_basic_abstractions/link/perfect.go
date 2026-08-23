package link

import (
	"fmt"
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ============================================================================
// PERFECT LINK (PL)
// ============================================================================

// PerfectLink models the Perfect Link abstraction.
//
// From Cachin et al., Section 2.4.3:
//
// Properties:
//
//	PL1 (Reliable delivery): If a correct process p sends a message m
//	  to a correct process q, then q eventually delivers m.
//	PL2 (No duplication): No message is delivered by a process more
//	  than once.
//	PL3 (No creation): If some process q delivers a message m with
//	  sender p, then m was previously sent to q by process p.
//
// Implementation (Algorithm 2.2):
//
//	Uses a Stubborn Link underneath.
//	Maintains a set of already-delivered messages.
//	On receive from stubborn link, checks if message was already delivered.
//	Only delivers NEW messages.
//
// This is composition in action:
//
//	Perfect = Stubborn + Deduplication
type PerfectLink struct {
	// underlying is the stubborn link used for reliable transmission.
	underlying *StubbornLink

	// delivered tracks which messages have already been delivered.
	// Key format: "from:to:type:data" — unique message identifier.
	delivered map[string]bool
}

// NewPerfectLink creates a perfect link wrapping a stubborn link.
func NewPerfectLink(sl *StubbornLink) *PerfectLink {
	return &PerfectLink{
		underlying: sl,
		delivered:  make(map[string]bool),
	}
}

// Send forwards the message through the stubborn link.
func (l *PerfectLink) Send(msg process.Message, rng *rand.Rand) []process.Message {
	return l.underlying.Send(msg, rng)
}

// Tick delegates to the stubborn link for retransmissions.
func (l *PerfectLink) Tick(step int, rng *rand.Rand) []process.Message {
	return l.underlying.Tick(step, rng)
}

// Receive implements deduplication.
// Only delivers messages that haven't been seen before.
func (l *PerfectLink) Receive(msg process.Message) (process.Message, bool) {
	key := messageKey(msg)
	if l.delivered[key] {
		return msg, false // Already delivered — suppress duplicate
	}
	l.delivered[key] = true
	return msg, true
}

// messageKey creates a unique identifier for deduplication.
// In a real system, this would use sequence numbers.
func messageKey(msg process.Message) string {
	return fmt.Sprintf("%s:%s:%s:%v", msg.From, msg.To, msg.Type, msg.Data)
}
