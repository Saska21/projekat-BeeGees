package process

// ============================================================================
// PROCESS INTERFACE
// ============================================================================

// Process represents a participant in a distributed computation.
//
// From Cachin et al., Section 2.1:
//
//	"We consider a distributed system as a collection of processes
//	 that communicate by exchanging messages."
//
// Key modelling rules:
//  1. Each process has a unique identifier (ProcessID).
//  2. A process is a STATE MACHINE: it has internal state that changes
//     only when it handles a message.
//  3. On each step, a process handles exactly ONE incoming message.
//  4. Handling produces zero or more outgoing messages.
//  5. Processes NEVER send messages directly to other processes;
//     all communication goes through the runtime scheduler.
//  6. Handle() performs ONE ATOMIC STEP — no partial execution.
//
// IMPORTANT: Handle must be DETERMINISTIC given the same state and message.
// All nondeterminism (delivery order, losses) comes from the runtime/links.
type Process interface {
	// ID returns the unique identifier of this process.
	ID() ProcessID

	// Handle processes a single incoming message and returns
	// zero or more outgoing messages.
	//
	// This is ONE ATOMIC STEP in the computation.
	//
	// The process may update its internal state and produce
	// response messages, but it must NOT:
	//   - Send messages directly (return them instead)
	//   - Block or wait
	//   - Access other processes' state
	//   - Use goroutines or channels
	//
	// The returned messages will be enqueued by the runtime
	// and delivered according to the configured link abstraction.
	Handle(msg Message) []Message
}
