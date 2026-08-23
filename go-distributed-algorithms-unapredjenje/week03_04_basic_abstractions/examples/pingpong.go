package examples

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/failures"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/link"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
	simrt "github.com/danilokacanski/da/week03_04_basic_abstractions/runtime"
)

// PingProcess sends PING messages and counts PONG responses.
type PingProcess struct {
	id    process.ProcessID
	peer  process.ProcessID
	count int
	limit int
}

func (p *PingProcess) ID() process.ProcessID { return p.id }

func (p *PingProcess) Handle(msg process.Message) []process.Message {
	switch msg.Type {
	case "INIT":
		fmt.Printf("  [%s] Starting ping-pong (limit: %d)\n", p.id, p.limit)
		return []process.Message{
			process.NewMessage(p.id, p.peer, "PING", p.count),
		}
	case "PONG":
		fmt.Printf("  [%s] Got PONG(%v) — exchange #%d done\n", p.id, msg.Data, p.count)
		p.count++
		if p.count < p.limit {
			return []process.Message{
				process.NewMessage(p.id, p.peer, "PING", p.count),
			}
		}
		fmt.Printf("  [%s] All %d exchanges complete!\n", p.id, p.limit)
	}
	return nil
}

// PongProcess responds to PING messages with PONG.
type PongProcess struct {
	id process.ProcessID
}

func (p *PongProcess) ID() process.ProcessID { return p.id }

func (p *PongProcess) Handle(msg process.Message) []process.Message {
	if msg.Type == "PING" {
		fmt.Printf("  [%s] Got PING(%v) — sending PONG\n", p.id, msg.Data)
		return []process.Message{
			process.NewMessage(p.id, msg.From, "PONG", msg.Data),
		}
	}
	return nil
}

// RunPingPong demonstrates basic message passing with perfect links.
func RunPingPong() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 1: PING-PONG (Perfect Links)")
	fmt.Println("  Demonstrates: message passing, nondeterministic scheduling,")
	fmt.Println("  and perfect link guarantees (PL1, PL2, PL3).")
	fmt.Println("================================================================")
	fmt.Println()

	fll := link.NewFairLossLink(0.0)
	sl := link.NewStubbornLink(fll, 5)
	pl := link.NewPerfectLink(sl)

	fm := failures.NewNoFailure()
	rt := simrt.NewRuntime(42, pl, fm, 30)

	rt.AddChecker(simrt.NoDuplicationChecker())
	rt.AddChecker(simrt.ReliableDeliveryChecker())

	rt.Register(&PingProcess{id: "ping", peer: "pong", limit: 3})
	rt.Register(&PongProcess{id: "pong"})

	rt.Inject(process.NewMessage("system", "ping", "INIT", nil))
	rt.Run()
}
