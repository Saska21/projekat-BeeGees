package consensus

import (
	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// BEST-EFFORT BROADCAST (simplified)
// ============================================================================
//
// In the book, Algorithm 5.1 uses BestEffortBroadcast (BEB).
// BEB guarantees:
//   - BEB1 (Validity): if a correct process broadcasts m, every correct process delivers m.
//   - BEB2 (No duplication): a message is delivered at most once.
//   - BEB3 (No creation): only broadcast messages are delivered.
//
// Simplification for this teaching implementation:
//   - We implement broadcast as sending one unicast message to EVERY process ID,
//     including the sender itself (self-messages flow through the runtime so
//     they appear in the trace like any other message).
//   - Reliability comes from the underlying PerfectLink (dedup + stubborn retransmit).
//   - This is logically equivalent to BEB over a perfect-link network.

// Broadcast sends a ConsensusMessage to every process in allIDs.
//
// Parameters:
//   - send:    the send callback provided by the runtime (from ConsensusNode.Run)
//   - from:    sender process ID
//   - allIDs:  all process IDs in the system
//   - payload: the ConsensusMessage to broadcast
//
// The message is wrapped in a process.Message with Type=ConsensusType so the
// runtime routes it correctly and logs it in the trace.
func Broadcast(
	send func(process.Message),
	from process.ProcessID,
	allIDs []process.ProcessID,
	payload ConsensusMessage,
) {
	for _, to := range allIDs {
		msg := process.NewMessage(from, to, ConsensusType, payload)
		send(msg)
	}
}
