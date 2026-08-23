// Package runtime implements the single-threaded scheduler for the
// distributed algorithms simulator.
//
// From Cachin et al., Section 2.1:
//
//	"We consider a distributed system as a collection of processes that
//	 communicate by exchanging messages... An execution is a sequence
//	 of steps. In each step, a process may receive a message, perform
//	 local computation, and send messages."
//
// The runtime orchestrates this execution model:
//  1. Maintains a queue of in-transit messages
//  2. On each step, picks a message NON-DETERMINISTICALLY
//  3. Applies the failure model (crash, omission, byzantine)
//  4. Applies the link abstraction (loss, dedup, authentication)
//  5. Delivers the message to the recipient process
//  6. Collects outgoing messages and enqueues them
//  7. Logs everything in the execution trace
//  8. Runs safety/liveness checkers
//
// KEY DESIGN DECISION: No goroutines, no channels, no networking.
// All nondeterminism comes from math/rand with an explicit seed.
// This makes executions REPRODUCIBLE — same seed, same execution.
package runtime

import (
	"fmt"
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/failures"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/link"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ============================================================================
// RUNTIME
// ============================================================================

// Runtime is the central scheduler for the distributed system simulation.
//
// It manages:
//   - A set of registered processes
//   - A message queue (simulating the network)
//   - A link abstraction (controlling delivery guarantees)
//   - A failure model (controlling process faults)
//   - An execution trace (recording all events)
//   - Property checkers (verifying safety/liveness)
//   - A PRNG for reproducible nondeterminism
type Runtime struct {
	// processes maps process IDs to process implementations.
	processes map[process.ProcessID]process.Process

	// queue is the set of in-transit messages.
	// In a real network, messages are "in the wire."
	// In our simulator, they wait here until scheduled.
	queue []process.Message

	// link is the communication link abstraction.
	link link.Link

	// failures is the failure model.
	failures failures.FailureInjector

	// trace records all events in the execution.
	trace *Trace

	// checkers verify properties of the execution.
	checkers []Checker

	// interceptor allows modifying messages in transit (for teaching).
	// This simulates a network-level attacker.
	interceptor func(process.Message) process.Message

	// rng is the random number generator for nondeterministic scheduling.
	// Using an explicit seed makes executions REPRODUCIBLE.
	rng *rand.Rand

	// stepCount tracks the current step number.
	stepCount int

	// maxSteps limits the execution length (prevents infinite loops).
	maxSteps int
}

// NewRuntime creates a new simulation runtime.
//
// Parameters:
//   - seed: random number seed for reproducibility
//   - linkLayer: the link abstraction to use
//   - fm: the failure model to use
//   - maxSteps: maximum number of steps before termination
func NewRuntime(seed int64, linkLayer link.Link, fm failures.FailureInjector, maxSteps int) *Runtime {
	return &Runtime{
		processes: make(map[process.ProcessID]process.Process),
		queue:     make([]process.Message, 0),
		link:      linkLayer,
		failures:  fm,
		trace:     NewTrace(true), // verbose by default
		checkers:  make([]Checker, 0),
		rng:       rand.New(rand.NewSource(seed)),
		maxSteps:  maxSteps,
	}
}

// ============================================================================
// CONFIGURATION
// ============================================================================

// Register adds a process to the simulation.
func (r *Runtime) Register(p process.Process) {
	r.processes[p.ID()] = p
}

// AddChecker adds a property checker.
func (r *Runtime) AddChecker(c Checker) {
	r.checkers = append(r.checkers, c)
}

// SetInterceptor configures a message interceptor.
// The interceptor is called after a message is dequeued but BEFORE
// it passes through the link's Receive() check.
//
// This models a network-level attacker who can observe and modify
// messages in transit. Use this for teaching about authentication.
//
// Pass nil to remove the interceptor.
func (r *Runtime) SetInterceptor(fn func(process.Message) process.Message) {
	r.interceptor = fn
}

// SetVerbose controls trace output.
func (r *Runtime) SetVerbose(v bool) {
	r.trace.verbose = v
}

// ============================================================================
// MESSAGE INJECTION
// ============================================================================

// Inject places a message directly into the queue.
// This bypasses the link layer entirely (no loss, no authentication).
// Use this for:
//   - Sending INIT messages to start processes
//   - Testing with specific message sequences
//   - Injecting "system" messages
//
// Injected messages are marked with Meta["_injected"] = true so the
// runtime skips link.Receive() for them (they have no MAC to verify).
func (r *Runtime) Inject(msg process.Message) {
	if msg.Meta == nil {
		msg.Meta = make(map[string]any)
	}
	msg.Meta["_injected"] = true
	r.queue = append(r.queue, msg)
	r.trace.Log(Event{
		Step:    r.stepCount,
		Type:    EventEnqueue,
		Message: &msg,
		Detail:  fmt.Sprintf("Injected: %s", msg.String()),
	})
}

// ============================================================================
// EXECUTION
// ============================================================================

// Step executes one step of the simulation.
//
// Returns true if a step was executed, false if the queue is empty.
//
// One step of the simulation:
//  1. Process retransmissions from the link layer (stubborn links)
//  2. If queue is empty, return false
//  3. Pick a RANDOM message from the queue (nondeterministic scheduling)
//  4. Apply interceptor (network attacker simulation)
//  5. Check if recipient process is alive (failure model)
//  6. Check if message should be delivered (omission model)
//  7. Check link-level receive (deduplication, authentication)
//  8. Deliver message to recipient process
//  9. Collect outgoing messages
//  10. Apply failure model to outgoing (Byzantine alteration)
//  11. Send each outgoing through the link layer
//  12. Run safety checkers
func (r *Runtime) Step() bool {
	r.stepCount++

	// --- 1. Link retransmissions ---
	retransmissions := r.link.Tick(r.stepCount, r.rng)
	for _, msg := range retransmissions {
		r.queue = append(r.queue, msg)
		msgCopy := msg
		r.trace.Log(Event{
			Step:    r.stepCount,
			Type:    EventRetransmit,
			Message: &msgCopy,
			Detail:  fmt.Sprintf("Retransmit: %s", msg.String()),
		})
	}

	// --- 2. Check if queue is empty ---
	if len(r.queue) == 0 {
		return false
	}

	// --- 3. Nondeterministic scheduling ---
	// Pick a RANDOM message from the queue.
	// This models the fact that in a real distributed system,
	// message delivery order is NOT guaranteed.
	idx := r.rng.Intn(len(r.queue))
	msg := r.queue[idx]
	// Remove from queue (swap with last and shrink)
	r.queue[idx] = r.queue[len(r.queue)-1]
	r.queue = r.queue[:len(r.queue)-1]

	r.trace.Log(Event{
		Step:    r.stepCount,
		Type:    EventDequeue,
		Message: &msg,
		Detail:  fmt.Sprintf("Dequeued: %s", msg.String()),
	})

	// --- 4. Apply interceptor (network attacker) ---
	if r.interceptor != nil {
		original := msg.String()
		msg = r.interceptor(msg)
		if msg.String() != original {
			r.trace.Log(Event{
				Step:    r.stepCount,
				Type:    EventForge,
				Message: &msg,
				Detail:  fmt.Sprintf("Interceptor modified message: %s → %s", original, msg.String()),
			})
		}
	}

	// --- 5. Check if recipient is alive ---
	if !r.failures.IsAlive(msg.To) {
		r.trace.Log(Event{
			Step:    r.stepCount,
			Type:    EventDrop,
			Message: &msg,
			Process: msg.To,
			Detail:  fmt.Sprintf("Process %s is crashed — dropping: %s", msg.To, msg.String()),
		})
		return true
	}

	// --- 6. Check failure-based delivery (omission) ---
	if !r.failures.ShouldDeliver(msg) {
		r.trace.Log(Event{
			Step:    r.stepCount,
			Type:    EventDrop,
			Message: &msg,
			Process: msg.To,
			Detail:  fmt.Sprintf("Omission failure — dropping: %s", msg.String()),
		})
		return true
	}

	// --- 7. Link receive check (dedup, auth) ---
	// Skip for injected messages (they bypass the link layer).
	_, isInjected := msg.Meta["_injected"]
	if !isInjected {
		msg, accepted := r.link.Receive(msg)
		if !accepted {
			r.trace.Log(Event{
				Step:    r.stepCount,
				Type:    EventDrop,
				Message: &msg,
				Process: msg.To,
				Detail:  fmt.Sprintf("Link rejected (dedup/auth): %s", msg.String()),
			})
			return true
		}
	}

	// --- 8. Deliver to recipient process ---
	proc, exists := r.processes[msg.To]
	if !exists {
		r.trace.Log(Event{
			Step:    r.stepCount,
			Type:    EventDrop,
			Message: &msg,
			Detail:  fmt.Sprintf("Unknown recipient %s — dropping: %s", msg.To, msg.String()),
		})
		return true
	}

	r.trace.Log(Event{
		Step:    r.stepCount,
		Type:    EventDeliver,
		Message: &msg,
		Process: msg.To,
		Detail:  fmt.Sprintf("Delivered to %s: %s", msg.To, msg.String()),
	})

	// --- 9. Process handles the message ---
	outgoing := proc.Handle(msg)

	// --- 10. Apply failure model to outgoing (Byzantine) ---
	outgoing = r.failures.MaybeAlter(msg.To, outgoing, r.rng)

	// --- 11. Send each outgoing through the link ---
	for _, out := range outgoing {
		outCopy := out
		r.trace.Log(Event{
			Step:    r.stepCount,
			Type:    EventSend,
			Message: &outCopy,
			Process: msg.To,
			Detail:  fmt.Sprintf("Process %s sends: %s", msg.To, out.String()),
		})

		// Link may drop, duplicate, or authenticate the message
		enqueued := r.link.Send(out, r.rng)
		for _, eq := range enqueued {
			eqCopy := eq
			r.queue = append(r.queue, eqCopy)
			r.trace.Log(Event{
				Step:    r.stepCount,
				Type:    EventEnqueue,
				Message: &eqCopy,
				Detail:  fmt.Sprintf("Enqueued: %s", eq.String()),
			})
		}
	}

	// --- 12. Run safety checkers ---
	for _, checker := range r.checkers {
		if checker.AtEnd {
			continue // Skip liveness checkers during execution
		}
		if violation := checker.Check(r.trace.Events()); violation != "" {
			r.trace.Log(Event{
				Step:   r.stepCount,
				Type:   EventViolation,
				Detail: fmt.Sprintf("[%s] %s", checker.Name, violation),
			})
		}
	}

	return true
}

// Run executes the simulation until the queue is empty or maxSteps is reached.
func (r *Runtime) Run() {
	fmt.Printf("\n>>> Starting simulation (max %d steps, seed produces deterministic execution)\n\n", r.maxSteps)

	for r.stepCount < r.maxSteps {
		if !r.Step() {
			fmt.Printf("\n>>> Queue empty at step %d — execution complete.\n", r.stepCount)
			break
		}
	}

	if r.stepCount >= r.maxSteps {
		fmt.Printf("\n>>> Reached maximum steps (%d) — execution terminated.\n", r.maxSteps)
	}

	// --- Run liveness checkers at the end ---
	for _, checker := range r.checkers {
		if !checker.AtEnd {
			continue
		}
		if violation := checker.Check(r.trace.Events()); violation != "" {
			r.trace.Log(Event{
				Step:   r.stepCount,
				Type:   EventViolation,
				Detail: fmt.Sprintf("[LIVENESS: %s] %s", checker.Name, violation),
			})
		}
	}

	// Print summary
	r.trace.Summary()
}

// ============================================================================
// INSPECTION
// ============================================================================

// Trace returns the execution trace for external analysis.
func (r *Runtime) Trace() *Trace {
	return r.trace
}

// QueueSize returns the current number of in-transit messages.
func (r *Runtime) QueueSize() int {
	return len(r.queue)
}

// StepCount returns the current step number.
func (r *Runtime) StepCount() int {
	return r.stepCount
}

// CrashProcess injects a crash failure for the named process.
func (r *Runtime) CrashProcess(pid process.ProcessID) {
	r.failures.CrashProcess(pid)
	r.trace.Log(Event{
		Step:    r.stepCount,
		Type:    EventCrash,
		Process: pid,
		Detail:  fmt.Sprintf("Process %s CRASHED", pid),
	})
}
