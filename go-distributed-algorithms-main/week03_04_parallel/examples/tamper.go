package examples

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/crypto"
	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"
)

// ============================================================================
// TAMPER EXAMPLE — Authenticated Links + Forgery Detection
// ============================================================================

// SenderProcess sends a single authenticated message.
type SenderProcess struct {
	id   process.ProcessID
	peer process.ProcessID
}

func (p *SenderProcess) ID() process.ProcessID { return p.id }

func (p *SenderProcess) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inbox:
			if !ok {
				return
			}
			if msg.Type == "INIT" {
				fmt.Printf("  [%s] Sending SECRET to %s\n", p.id, p.peer)
				send(process.NewMessage(p.id, p.peer, "SECRET", "transfer $100 to Bob"))
			}
		}
	}
}

// ReceiverProcess receives authenticated messages.
type ReceiverProcess struct {
	id       process.ProcessID
	received int64        // atomic counter
	lastData atomic.Value // stores string
}

func (p *ReceiverProcess) ID() process.ProcessID { return p.id }

func (p *ReceiverProcess) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inbox:
			if !ok {
				return
			}
			if msg.Type == "SECRET" {
				atomic.AddInt64(&p.received, 1)
				if data, ok := msg.Data.(string); ok {
					p.lastData.Store(data)
				}
				fmt.Printf("  [%s] Received authenticated SECRET: %v\n", p.id, msg.Data)
			}
		}
	}
}

// RunTamper demonstrates authenticated links and message forgery detection.
func RunTamper() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 3: TAMPER (Authenticated Links + Forgery Detection)")
	fmt.Println("  Demonstrates: MACs, message authentication, and how")
	fmt.Println("  authenticated links detect network-level tampering.")
	fmt.Println("  Using goroutines and channels for parallel execution.")
	fmt.Println("================================================================")

	// -- Phase 1: Honest communication --
	fmt.Println("\n-- Phase 1: Honest Communication (no tampering) --")
	fmt.Println()

	keys := crypto.NewKeyRegistry()
	keys.RegisterProcess("alice")
	keys.RegisterProcess("bob")

	fll := link.NewFairLossLink(0.0, 42)
	sl := link.NewStubbornLink(fll)
	pl := link.NewPerfectLink(sl)
	al := link.NewAuthenticatedLink(pl, keys)

	fm := failures.NewNoFailure()
	rt := simrt.NewRuntime(al, fm,
		simrt.WithIdleTimeout(300*time.Millisecond),
		simrt.WithRetransmitInterval(100*time.Millisecond),
		simrt.WithMaxDuration(5*time.Second),
	)

	receiver := &ReceiverProcess{id: "bob"}
	rt.Register(&SenderProcess{id: "alice", peer: "bob"})
	rt.Register(receiver)
	rt.Inject(process.NewMessage("system", "alice", "INIT", nil))
	rt.Run()

	receivedCount := atomic.LoadInt64(&receiver.received)
	fmt.Printf("\n  Phase 1 Result: Bob received %d message(s)\n", receivedCount)
	if receivedCount > 0 {
		if data, ok := receiver.lastData.Load().(string); ok {
			fmt.Printf("  Content: %q — valid MAC\n", data)
		}
	}

	// -- Phase 2: Network attacker tampers with messages --
	fmt.Println("\n-- Phase 2: Network Attacker (tampering in transit) --")
	fmt.Println()

	keys2 := crypto.NewKeyRegistry()
	keys2.RegisterProcess("alice")
	keys2.RegisterProcess("bob")

	fll2 := link.NewFairLossLink(0.0, 99)
	sl2 := link.NewStubbornLink(fll2)
	pl2 := link.NewPerfectLink(sl2)
	al2 := link.NewAuthenticatedLink(pl2, keys2)

	fm2 := failures.NewNoFailure()
	rt2 := simrt.NewRuntime(al2, fm2,
		simrt.WithIdleTimeout(300*time.Millisecond),
		simrt.WithRetransmitInterval(100*time.Millisecond),
		simrt.WithMaxDuration(5*time.Second),
	)

	receiver2 := &ReceiverProcess{id: "bob"}
	rt2.Register(&SenderProcess{id: "alice", peer: "bob"})
	rt2.Register(receiver2)

	// Network interceptor modifies message content in transit.
	// MAC was computed over original content, so verification FAILS.
	rt2.SetInterceptor(func(msg process.Message) process.Message {
		if msg.Type == "SECRET" {
			tampered := msg.Clone()
			tampered.Data = "transfer $10000 to Eve"
			fmt.Printf("  [ATTACKER] Tampered: %q -> %q\n",
				"transfer $100 to Bob", tampered.Data)
			return tampered
		}
		return msg
	})

	rt2.Inject(process.NewMessage("system", "alice", "INIT", nil))
	rt2.Run()

	receivedCount2 := atomic.LoadInt64(&receiver2.received)
	fmt.Printf("\n  Phase 2 Result: Bob received %d message(s)\n", receivedCount2)
	if receivedCount2 == 0 {
		fmt.Println("  -> Tampered message REJECTED by authenticated link!")
		fmt.Println("  -> MAC verification failed (content was modified).")
		fmt.Println("  -> This is the AL1 (Authenticity) property.")
	} else {
		fmt.Println("  -> WARNING: Message delivered despite tampering.")
	}
	fmt.Println()
}
