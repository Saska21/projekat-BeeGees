package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"
)

// ============================================================================
// PING-PONG EXAMPLE — Goroutine-based Process Interface
// ============================================================================

// PingProcess sends PING messages and counts PONG responses.
// Uses the native concurrent Process interface (Run with select loop).
type PingProcess struct {
	id    process.ProcessID
	peer  process.ProcessID
	limit int
}

func (p *PingProcess) ID() process.ProcessID { return p.id }

func (p *PingProcess) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inbox:
			if !ok {
				return
			}
			switch msg.Type {
			case "INIT":
				fmt.Printf("  [%s] Starting ping-pong (limit: %d)\n", p.id, p.limit)
				send(process.NewMessage(p.id, p.peer, "PING", count))
			case "PONG":
				fmt.Printf("  [%s] Got PONG(%v) — exchange #%d done\n", p.id, msg.Data, count)
				count++
				if count < p.limit {
					send(process.NewMessage(p.id, p.peer, "PING", count))
				} else {
					fmt.Printf("  [%s] All %d exchanges complete!\n", p.id, p.limit)
				}
			}
		}
	}
}

// PongProcess responds to PING messages with PONG.
type PongProcess struct {
	id process.ProcessID
}

func (p *PongProcess) ID() process.ProcessID { return p.id }

func (p *PongProcess) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inbox:
			if !ok {
				return
			}
			if msg.Type == "PING" {
				fmt.Printf("  [%s] Got PING(%v) — sending PONG\n", p.id, msg.Data)
				send(process.NewMessage(p.id, msg.From, "PONG", msg.Data))
			}
		}
	}
}

// RunPingPong demonstrates basic message passing with perfect links
// using goroutines and channels.
func RunPingPong() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 1: PING-PONG (Perfect Links, Goroutines)")
	fmt.Println("  Demonstrates: concurrent message passing, goroutines,")
	fmt.Println("  channels, and perfect link guarantees (PL1, PL2, PL3).")
	fmt.Println("================================================================")
	fmt.Println()

	fll := link.NewFairLossLink(0.0, 42)
	sl := link.NewStubbornLink(fll)
	pl := link.NewPerfectLink(sl)

	fm := failures.NewNoFailure()
	rt := simrt.NewRuntime(pl, fm,
		simrt.WithIdleTimeout(300*time.Millisecond),
		simrt.WithRetransmitInterval(100*time.Millisecond),
		simrt.WithMaxDuration(5*time.Second),
	)

	rt.AddChecker(simrt.NoDuplicationChecker())
	rt.AddChecker(simrt.ReliableDeliveryChecker())

	rt.Register(&PingProcess{id: "ping", peer: "pong", limit: 3})
	rt.Register(&PongProcess{id: "pong"})

	rt.Inject(process.NewMessage("system", "ping", "INIT", nil))
	rt.Run()
}
