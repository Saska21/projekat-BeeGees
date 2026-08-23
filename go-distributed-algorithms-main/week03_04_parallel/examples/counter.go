package examples

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"
)

// ============================================================================
// COUNTER EXAMPLE — Message Loss with Fair-Loss Links
// ============================================================================

// This example demonstrates:
//   - How message loss affects system behavior
//   - The difference between unreliable and fair-loss links
//   - Why reliable links matter for correctness
//
// Uses WrapHandler to convert synchronous Handlers into concurrent Processes.

// CounterHandler maintains a running count of INCREMENT messages received.
// Uses sync/atomic for the Count so it can be safely read from outside.
type CounterHandler struct {
	id    process.ProcessID
	Count int64 // atomic — safe for concurrent reads
}

func (p *CounterHandler) ID() process.ProcessID { return p.id }

func (p *CounterHandler) Handle(msg process.Message) []process.Message {
	if msg.Type == "INCREMENT" {
		newCount := atomic.AddInt64(&p.Count, 1)
		fmt.Printf("  [%s] Received INCREMENT from %s — count is now %d\n",
			p.id, msg.From, newCount)
	}
	return nil
}

// WorkerHandler sends a fixed number of INCREMENT messages.
type WorkerHandler struct {
	id      process.ProcessID
	counter process.ProcessID
	total   int
}

func (p *WorkerHandler) ID() process.ProcessID { return p.id }

func (p *WorkerHandler) Handle(msg process.Message) []process.Message {
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
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 2: COUNTER (Fair-Loss Link, 50% Loss, Goroutines)")
	fmt.Println("  Demonstrates: message loss with unreliable links, safety")
	fmt.Println("  (no creation) vs liveness (reliable delivery) properties.")
	fmt.Println("================================================================")
	fmt.Println()

	// Use a fair-loss link with 50% loss rate
	fll := link.NewFairLossLink(0.5, 42)

	fm := failures.NewNoFailure()
	rt := simrt.NewRuntime(fll, fm,
		simrt.WithIdleTimeout(300*time.Millisecond),
		simrt.WithRetransmitInterval(100*time.Millisecond),
		simrt.WithMaxDuration(5*time.Second),
	)

	// Create counter process (we keep a reference to inspect the final count)
	counterHandler := &CounterHandler{id: "counter"}
	rt.Register(process.WrapHandler(counterHandler))

	// Create two worker processes, each sending 5 increments
	rt.Register(process.WrapHandler(&WorkerHandler{id: "worker-A", counter: "counter", total: 5}))
	rt.Register(process.WrapHandler(&WorkerHandler{id: "worker-B", counter: "counter", total: 5}))

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
	finalCount := atomic.LoadInt64(&counterHandler.Count)
	fmt.Printf("\n  RESULT: Counter = %d out of 10 expected\n", finalCount)
	if finalCount < 10 {
		fmt.Println("  -> Messages were LOST due to the fair-loss link (50% loss rate).")
		fmt.Println("  -> This demonstrates why Perfect Links (with retransmission) are needed")
		fmt.Println("     for reliable delivery.")
	} else {
		fmt.Println("  -> All messages arrived despite the lossy link (lucky run!).")
	}
	fmt.Println()
}
