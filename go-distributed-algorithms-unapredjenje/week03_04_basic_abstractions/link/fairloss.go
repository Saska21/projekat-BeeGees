package link

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// Link is the interface for all communication link abstractions.
//
// In the formal model, a link provides two events:
//   - Send: process p requests to send message m to process q
//   - Deliver: message m is delivered to process q
//
// In our simulator, the link layer sits between the runtime scheduler
// and message delivery:
//  1. Send() is called when a process produces an outgoing message
//  2. Tick() is called each step to handle retransmissions
//  3. Receive() is called before delivering to check dedup/auth
type Link interface {
	// Send processes an outgoing message. Returns messages to enqueue.
	// May return:
	//   - empty slice: message was lost/dropped
	//   - single message: normal delivery
	//   - multiple messages: duplication occurred
	Send(msg process.Message, rng *rand.Rand) []process.Message

	// Tick is called each scheduler step. Returns any retransmissions.
	// Only meaningful for stubborn/higher links.
	Tick(step int, rng *rand.Rand) []process.Message

	// Receive checks if a message should be delivered.
	// Returns (processed_message, should_deliver).
	// Used for deduplication (perfect links) and authentication checks.
	Receive(msg process.Message) (process.Message, bool)
}

// ============================================================================
// FAIR-LOSS LINK (FLL, baseline)
// ============================================================================

// FairLossLink models the Fair-Loss Link abstraction.
//
// From Cachin et al., Section 2.4.1:
//
// Properties:
//
//	FLL1 (Fair-loss delivery): If a correct process p infinitely often
//	  sends a message m to a correct process q, then q delivers m an
//	  infinite number of times.
//	FLL2 (Finite duplication): If a correct process p sends a message m
//	  a finite number of times to process q, then m cannot be delivered
//	  an infinite number of times by q.
//	FLL3 (No creation): If some process q delivers a message m with
//	  sender p, then m was previously sent to q by process p.
//
// In our simulator:
//   - Each message has a probability of being lost
//   - Lost messages are simply not enqueued
//   - No message creation (we never fabricate messages)
//   - If a message is sent enough times, it will eventually get through
type FairLossLink struct {
	// LossRate controls the probability of losing a single message.
	LossRate float64
}

// NewFairLossLink creates a fair-loss link with the given loss rate.
func NewFairLossLink(lossRate float64) *FairLossLink {
	return &FairLossLink{LossRate: lossRate}
}

// Send may drop the message based on LossRate.
// This models FLL1: individual sends may fail, but repeated sends
// will eventually succeed.
func (l *FairLossLink) Send(msg process.Message, rng *rand.Rand) []process.Message {
	if rng.Float64() < l.LossRate {
		return nil // Message lost (but fair: repeated sends work)
	}
	return []process.Message{msg}
}

// Tick does nothing for fair-loss links (no automatic retransmission).
func (l *FairLossLink) Tick(step int, rng *rand.Rand) []process.Message {
	return nil
}

// Receive always accepts (fair-loss has no deduplication).
func (l *FairLossLink) Receive(msg process.Message) (process.Message, bool) {
	return msg, true
}
