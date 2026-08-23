// Package runtime implements the CONCURRENT scheduler for the
// distributed algorithms simulator.
//
// PARALLEL VERSION: Instead of a single-threaded step loop,
// each process runs in its own goroutine. Messages flow through
// channels. The runtime manages:
//   - A buffered network channel (simulating the network)
//   - One goroutine per registered process
//   - A router goroutine that picks messages from the network and delivers them
//   - A retransmitter goroutine for stubborn link retransmissions
//   - An idle monitor that cancels the context when no activity is detected
//
// Termination: The simulation ends when either:
//   - The idle timeout expires (no messages delivered for a period)
//   - The max duration is reached
package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// RUNTIME
// ============================================================================

// procEntry holds a process and its dedicated inbox channel.
type procEntry struct {
	proc  process.Process
	inbox chan process.Message
}

// Runtime is the concurrent scheduler for the distributed system simulation.
type Runtime struct {
	// processes maps process IDs to their entries (process + inbox).
	processes map[process.ProcessID]*procEntry

	// network is the buffered channel simulating the network.
	// All outgoing messages from processes go here first.
	network chan process.Message

	// link is the communication link abstraction.
	link link.Link

	// failures is the failure model.
	failures failures.FailureInjector

	// trace records all events in the execution.
	trace *Trace

	// checkers verify properties of the execution.
	checkers []Checker

	// interceptor allows modifying messages in transit (network attacker).
	interceptor func(process.Message) process.Message

	// lastDelivery tracks the time of the last delivered message (for idle detection).
	lastDelivery atomic.Value // stores time.Time

	// options
	maxDuration        time.Duration
	idleTimeout        time.Duration
	retransmitInterval time.Duration
	verbose            bool
}

// Option configures the Runtime.
type Option func(*Runtime)

// WithMaxDuration sets the maximum wall-clock duration for the simulation.
func WithMaxDuration(d time.Duration) Option {
	return func(r *Runtime) { r.maxDuration = d }
}

// WithIdleTimeout sets how long the runtime waits with no deliveries
// before concluding the simulation is finished.
func WithIdleTimeout(d time.Duration) Option {
	return func(r *Runtime) { r.idleTimeout = d }
}

// WithRetransmitInterval sets the interval between retransmission rounds.
func WithRetransmitInterval(d time.Duration) Option {
	return func(r *Runtime) { r.retransmitInterval = d }
}

// WithVerbose enables or disables trace printing.
func WithVerbose(v bool) Option {
	return func(r *Runtime) { r.verbose = v }
}

// NewRuntime creates a new concurrent simulation runtime.
func NewRuntime(linkLayer link.Link, fm failures.FailureInjector, opts ...Option) *Runtime {
	r := &Runtime{
		processes:          make(map[process.ProcessID]*procEntry),
		network:            make(chan process.Message, 1000),
		link:               linkLayer,
		failures:           fm,
		checkers:           make([]Checker, 0),
		maxDuration:        10 * time.Second,
		idleTimeout:        500 * time.Millisecond,
		retransmitInterval: 100 * time.Millisecond,
		verbose:            true,
	}
	for _, o := range opts {
		o(r)
	}
	r.trace = NewTrace(r.verbose)
	r.lastDelivery.Store(time.Now())
	return r
}

// ============================================================================
// CONFIGURATION
// ============================================================================

// Register adds a process to the simulation.
// Must be called before Run().
func (r *Runtime) Register(p process.Process) {
	r.processes[p.ID()] = &procEntry{
		proc:  p,
		inbox: make(chan process.Message, 100),
	}
}

// AddChecker adds a property checker.
func (r *Runtime) AddChecker(c Checker) {
	r.checkers = append(r.checkers, c)
}

// SetInterceptor configures a message interceptor (network attacker).
func (r *Runtime) SetInterceptor(fn func(process.Message) process.Message) {
	r.interceptor = fn
}

// ============================================================================
// MESSAGE INJECTION
// ============================================================================

// Inject places a message directly into the recipient's inbox,
// bypassing the network channel and link layer.
// Used for INIT messages to start processes.
func (r *Runtime) Inject(msg process.Message) {
	if msg.Meta == nil {
		msg.Meta = make(map[string]any)
	}
	msg.Meta["_injected"] = true

	entry, exists := r.processes[msg.To]
	if !exists {
		r.trace.Log(Event{
			Type:   EventDrop,
			Detail: fmt.Sprintf("Inject failed — unknown recipient %s: %s", msg.To, msg.String()),
		})
		return
	}

	r.trace.Log(Event{
		Type:    EventEnqueue,
		Message: &msg,
		Detail:  fmt.Sprintf("Injected: %s", msg.String()),
	})

	entry.inbox <- msg
}

// ============================================================================
// SEND FUNCTION FACTORY
// ============================================================================

// makeSendFunc creates a send function for a process.
// The send function:
//  1. Logs the send event
//  2. Applies the failure model (Byzantine alteration)
//  3. Sends each message through the link layer
//  4. Puts surviving messages on the network channel
func (r *Runtime) makeSendFunc(ctx context.Context, pid process.ProcessID) func(process.Message) {
	return func(msg process.Message) {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgCopy := msg
		r.trace.Log(Event{
			Type:    EventSend,
			Message: &msgCopy,
			Process: pid,
			Detail:  fmt.Sprintf("Process %s sends: %s", pid, msg.String()),
		})

		// Apply failure model (Byzantine alteration)
		altered := r.failures.MaybeAlter(pid, []process.Message{msg})

		// Send each message through the link layer
		for _, m := range altered {
			enqueued := r.link.Send(m)
			for _, eq := range enqueued {
				eqCopy := eq
				// Place the message on the network channel first, then log ENQUEUE.
				select {
				case <-ctx.Done():
					return
				case r.network <- eqCopy:
					// Log only after the message has been successfully enqueued
					r.trace.Log(Event{
						Type:    EventEnqueue,
						Message: &eqCopy,
						Detail:  fmt.Sprintf("Enqueued: %s", eq.String()),
					})
					// Update lastDelivery (activity) so the idle monitor sees progress.
					// This prevents the idle monitor from cancelling the simulation
					// while messages are queued but not yet delivered.
					r.lastDelivery.Store(time.Now())
				}
			}
		}
	}
}

// ============================================================================
// ROUTER GOROUTINE
// ============================================================================

// routeMessages reads from the network channel and delivers messages
// to the appropriate process inbox, applying interceptor, failure model,
// and link-level receive checks.
func (r *Runtime) routeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-r.network:
			if !ok {
				return
			}
			r.processMessage(ctx, msg)
		}
	}
}

// processMessage handles a single message from the network.
// Steps: dequeue -> intercept -> alive check -> delivery check -> link receive -> deliver
//
// Each trace event gets its own message snapshot to prevent pointer
// aliasing (where a later mutation retroactively changes earlier events).
func (r *Runtime) processMessage(ctx context.Context, msg process.Message) {
	dequeued := msg // snapshot for trace
	r.trace.Log(Event{
		Type:    EventDequeue,
		Message: &dequeued,
		Detail:  fmt.Sprintf("Dequeued: %s", msg.String()),
	})

	// --- Apply interceptor (network attacker) ---
	if r.interceptor != nil {
		original := msg.String()
		msg = r.interceptor(msg)
		if msg.String() != original {
			forged := msg // snapshot
			r.trace.Log(Event{
				Type:    EventForge,
				Message: &forged,
				Detail:  fmt.Sprintf("Interceptor modified: %s -> %s", original, msg.String()),
			})
		}
	}

	// --- Check if recipient is alive ---
	if !r.failures.IsAlive(msg.To) {
		dropped := msg
		r.trace.Log(Event{
			Type:    EventDrop,
			Message: &dropped,
			Process: msg.To,
			Detail:  fmt.Sprintf("Process %s is crashed — dropping: %s", msg.To, msg.String()),
		})
		return
	}

	// --- Check failure-based delivery (omission) ---
	if !r.failures.ShouldDeliver(msg) {
		dropped := msg
		r.trace.Log(Event{
			Type:    EventDrop,
			Message: &dropped,
			Process: msg.To,
			Detail:  fmt.Sprintf("Omission failure — dropping: %s", msg.String()),
		})
		return
	}

	// --- Link receive check (dedup, auth) ---
	_, isInjected := msg.Meta["_injected"]
	if !isInjected {
		msg, accepted := r.link.Receive(msg)
		if !accepted {
			rejected := msg
			r.trace.Log(Event{
				Type:    EventDrop,
				Message: &rejected,
				Process: msg.To,
				Detail:  fmt.Sprintf("Link rejected (dedup/auth): %s", msg.String()),
			})
			return
		}
	}

	// --- Deliver to recipient's inbox ---
	entry, exists := r.processes[msg.To]
	if !exists {
		unknown := msg
		r.trace.Log(Event{
			Type:    EventDrop,
			Message: &unknown,
			Detail:  fmt.Sprintf("Unknown recipient %s — dropping: %s", msg.To, msg.String()),
		})
		return
	}

	// Blocking send with timeout — prevents losing messages that
	// passed dedup. If we silently dropped after dedup marked it
	// delivered, the message would be permanently lost.
	select {
	case <-ctx.Done():
		return
	case entry.inbox <- msg:
		// success
	}

	delivered := msg // snapshot after successful inbox delivery
	r.trace.Log(Event{
		Type:    EventDeliver,
		Message: &delivered,
		Process: msg.To,
		Detail:  fmt.Sprintf("Delivered to %s: %s", msg.To, msg.String()),
	})

	// Update last delivery time for idle detection
	r.lastDelivery.Store(time.Now())

	// --- Run safety checkers ---
	for _, checker := range r.checkers {
		if checker.AtEnd {
			continue
		}
		if violation := checker.Check(r.trace.Events()); violation != "" {
			r.trace.Log(Event{
				Type:   EventViolation,
				Detail: fmt.Sprintf("[%s] %s", checker.Name, violation),
			})
		}
	}
}

// ============================================================================
// RETRANSMITTER GOROUTINE
// ============================================================================

// retransmitLoop periodically calls link.Retransmissions() and pushes
// results onto the network channel.
func (r *Runtime) retransmitLoop(ctx context.Context) {
	ticker := time.NewTicker(r.retransmitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			retransmissions := r.link.Retransmissions()
			for _, msg := range retransmissions {
				msgCopy := msg
				r.trace.Log(Event{
					Type:    EventRetransmit,
					Message: &msgCopy,
					Detail:  fmt.Sprintf("Retransmit: %s", msg.String()),
				})

				select {
				case <-ctx.Done():
					return
				case r.network <- msgCopy:
				default:
					// Network channel full — drop retransmission
				}
			}
		}
	}
}

// ============================================================================
// IDLE MONITOR GOROUTINE
// ============================================================================

// idleMonitor cancels the context when no deliveries have occurred
// for the idle timeout duration. Also exits when ctx is cancelled
// (e.g. by max-duration timeout) to avoid leaking this goroutine.
func (r *Runtime) idleMonitor(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(r.idleTimeout / 4)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			last := r.lastDelivery.Load().(time.Time)
			if time.Since(last) > r.idleTimeout {
				cancel()
				return
			}
		}
	}
}

// ============================================================================
// EXECUTION
// ============================================================================

// Run executes the simulation.
//
// 1. Creates a context with max duration timeout
// 2. Starts a goroutine for each registered process
// 3. Starts the router goroutine
// 4. Starts the retransmitter goroutine
// 5. Starts the idle monitor goroutine
// 6. Waits for all process goroutines to finish
// 7. Runs liveness checkers
func (r *Runtime) Run() {
	fmt.Printf("\n>>> Starting PARALLEL simulation (max %s, idle timeout %s)\n\n",
		r.maxDuration, r.idleTimeout)

	// Create context with max duration
	ctx, cancel := context.WithTimeout(context.Background(), r.maxDuration)
	defer cancel()

	// Reset last delivery time
	r.lastDelivery.Store(time.Now())

	// We'll coordinate goroutines so we can drain the network before exiting.
	var wgProcesses sync.WaitGroup
	var wgRouter sync.WaitGroup
	var wgRetrans sync.WaitGroup

	// Start process goroutines
	for pid, entry := range r.processes {
		wgProcesses.Add(1)
		sendFn := r.makeSendFunc(ctx, pid)
		go func(e *procEntry, send func(process.Message)) {
			defer wgProcesses.Done()
			e.proc.Run(ctx, e.inbox, send)
		}(entry, sendFn)
	}

	// Start router goroutine
	wgRouter.Add(1)
	go func() {
		defer wgRouter.Done()
		r.routeMessages(ctx)
	}()

	// Start retransmitter goroutine with cancellable child context so we can stop it
	retransCtx, retransCancel := context.WithCancel(ctx)
	wgRetrans.Add(1)
	go func() {
		defer wgRetrans.Done()
		r.retransmitLoop(retransCtx)
	}()

	// Start idle monitor (cancels context when idle)
	go r.idleMonitor(ctx, cancel)

	// Wait for all process goroutines to finish
	wgProcesses.Wait()

	// Stop retransmitter and wait for it to finish so it won't write to the network
	retransCancel()
	wgRetrans.Wait()

	// Close the network channel to signal router to drain and exit
	close(r.network)

	// Wait for router to finish draining messages
	wgRouter.Wait()

	fmt.Printf("\n>>> Simulation ended.\n")

	// --- Run liveness checkers at the end ---
	for _, checker := range r.checkers {
		if !checker.AtEnd {
			continue
		}
		if violation := checker.Check(r.trace.Events()); violation != "" {
			r.trace.Log(Event{
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

// CrashProcess injects a crash failure for the named process.
func (r *Runtime) CrashProcess(pid process.ProcessID) {
	r.failures.CrashProcess(pid)
	r.trace.Log(Event{
		Type:    EventCrash,
		Process: pid,
		Detail:  fmt.Sprintf("Process %s CRASHED", pid),
	})
}
