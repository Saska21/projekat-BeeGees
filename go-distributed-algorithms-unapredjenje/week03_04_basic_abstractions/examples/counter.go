package examples

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/failures"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/link"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
	simrt "github.com/danilokacanski/da/week03_04_basic_abstractions/runtime"
)

// ============================================================================
// COUNTER EXAMPLE — Message Loss with Fair-Loss Links
// ============================================================================

// This example demonstrates:
//   - How message loss affects system behavior
//   - Why fair-loss links alone are not enough for reliable delivery
//   - Why reliable links matter for correctness
//   - Safety properties that ALWAYS hold (no creation)
//   - Liveness concerns when messages are lost
//
// Scenario:
//   Two worker processes each send 5 INCREMENT messages to a counter.
//   With a fair-loss link (50% loss), many messages are lost.
//
//   Expected: counter = 10 (if all messages arrive)
//   Actual:   counter < 10 (due to message loss)
//
// TEACHING POINT:
//   Even with message loss, the "no creation" property (PL3) holds:
//   the counter will NEVER exceed 10, because no messages are fabricated.
//   But reliable delivery (PL1) may NOT hold with a fair-loss link.

// CounterProcess maintains a running count of INCREMENT messages received.
type CounterProcess struct {
	id    process.ProcessID
	Count int // Exported so the example can inspect the final value
}

func (p *CounterProcess) ID() process.ProcessID { return p.id }

// Handle increments the counter for each INCREMENT message.
func (p *CounterProcess) Handle(msg process.Message) []process.Message {
	if msg.Type == "INCREMENT" {
		p.Count++
		fmt.Printf("  [%s] Received INCREMENT from %s — count is now %d\n",
			p.id, msg.From, p.Count)
	}
	return nil
}

// WorkerProcess sends a fixed number of INCREMENT messages.
type WorkerProcess struct {
	id      process.ProcessID
	counter process.ProcessID
	total   int
}

func (p *WorkerProcess) ID() process.ProcessID { return p.id }

// Handle: on INIT, sends all INCREMENT messages at once.
// In a real system, these would be sent over time. For simplicity,
// we enqueue them all in response to INIT.
func (p *WorkerProcess) Handle(msg process.Message) []process.Message {
	if msg.Type != "INIT" {
		return nil
	}

	fmt.Printf("  [%s] Sending %d INCREMENT messages to %s\n", p.id, p.total, p.counter)
	msgs := make([]process.Message, p.total)
	for i := 0; i < p.total; i++ {
		msgs[i] = process.NewMessage(p.id, p.counter, "INCREMENT", i+1)
	}
	return msgs
}

// RunCounter sets up and runs the counter demonstration.
func RunCounter() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        EXAMPLE 2: COUNTER (Fair-Loss Link, 50% Loss)       ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ Demonstrates: message loss with fair-loss links, safety     ║")
	fmt.Println("║ (no creation) vs liveness (reliable delivery) properties.   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Use a fair-loss link with 50% loss rate
	// This means each message has a 50% chance of being dropped
	fll := link.NewFairLossLink(0.5)

	// No process failures — the loss is at the LINK level
	fm := failures.NewNoFailure()

	// Create runtime
	rt := simrt.NewRuntime(42, fll, fm, 50)

	// Create counter process (we keep a reference to inspect the final count)
	counterProc := &CounterProcess{id: "counter"}
	rt.Register(counterProc)

	// Create two worker processes, each sending 5 increments
	rt.Register(&WorkerProcess{id: "worker-A", counter: "counter", total: 5})
	rt.Register(&WorkerProcess{id: "worker-B", counter: "counter", total: 5})

	// Add safety checker: counter should never exceed 10
	rt.AddChecker(simrt.Checker{
		Name:  "Counter bound",
		AtEnd: false,
		Check: func(events []simrt.Event) string {
			delivers := 0
			for _, e := range events {
				if e.Type == simrt.EventDeliver && e.Message != nil && e.Message.Type == "INCREMENT" {
					delivers++
				}
			}
			if delivers > 10 {
				return fmt.Sprintf("Counter bound violated: %d > 10", delivers)
			}
			return ""
		},
	})

	// Start both workers
	rt.Inject(process.NewMessage("system", "worker-A", "INIT", nil))
	rt.Inject(process.NewMessage("system", "worker-B", "INIT", nil))

	// Run
	rt.Run()

	// Report results
	fmt.Printf("\n  RESULT: Counter = %d out of 10 expected\n", counterProc.Count)
	if counterProc.Count < 10 {
		fmt.Println("  → Messages were LOST due to the fair-loss link (50% loss rate).")
		fmt.Println("  → This demonstrates why Perfect Links (with retransmission) are needed")
		fmt.Println("    for reliable delivery.")
	} else {
		fmt.Println("  → All messages arrived despite the lossy link (lucky seed!).")
	}
	fmt.Println()
}
