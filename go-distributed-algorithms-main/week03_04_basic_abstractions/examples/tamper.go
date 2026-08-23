package examples

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/crypto"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/failures"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/link"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
	simrt "github.com/danilokacanski/da/week03_04_basic_abstractions/runtime"
)

// SenderProcess sends a single authenticated message.
type SenderProcess struct {
	id   process.ProcessID
	peer process.ProcessID
}

func (p *SenderProcess) ID() process.ProcessID { return p.id }

func (p *SenderProcess) Handle(msg process.Message) []process.Message {
	if msg.Type == "INIT" {
		fmt.Printf("  [%s] Sending SECRET to %s\n", p.id, p.peer)
		return []process.Message{
			process.NewMessage(p.id, p.peer, "SECRET", "transfer $100 to Bob"),
		}
	}
	return nil
}

// ReceiverProcess receives authenticated messages.
type ReceiverProcess struct {
	id            process.ProcessID
	ReceivedCount int
	LastData      string
}

func (p *ReceiverProcess) ID() process.ProcessID { return p.id }

func (p *ReceiverProcess) Handle(msg process.Message) []process.Message {
	if msg.Type == "SECRET" {
		p.ReceivedCount++
		if data, ok := msg.Data.(string); ok {
			p.LastData = data
		}
		fmt.Printf("  [%s] Received authenticated SECRET: %v\n", p.id, msg.Data)
	}
	return nil
}

// RunTamper demonstrates authenticated links and message forgery detection.
func RunTamper() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 3: TAMPER (Authenticated Links + Forgery Detection)")
	fmt.Println("  Demonstrates: MACs, message authentication, and how")
	fmt.Println("  authenticated links detect network-level tampering.")
	fmt.Println("================================================================")

	// -- Phase 1: Honest communication --
	fmt.Println("\n-- Phase 1: Honest Communication (no tampering) --")
	fmt.Println()

	keys := crypto.NewKeyRegistry()
	keys.RegisterProcess("alice")
	keys.RegisterProcess("bob")

	fll := link.NewFairLossLink(0.0)
	sl := link.NewStubbornLink(fll, 5)
	pl := link.NewPerfectLink(sl)
	al := link.NewAuthenticatedLink(pl, keys)

	fm := failures.NewNoFailure()
	rt := simrt.NewRuntime(42, al, fm, 20)

	receiver := &ReceiverProcess{id: "bob"}
	rt.Register(&SenderProcess{id: "alice", peer: "bob"})
	rt.Register(receiver)
	rt.Inject(process.NewMessage("system", "alice", "INIT", nil))
	rt.Run()

	fmt.Printf("\n  Phase 1 Result: Bob received %d message(s)\n", receiver.ReceivedCount)
	if receiver.ReceivedCount > 0 {
		fmt.Printf("  Content: %q — valid MAC\n", receiver.LastData)
	}

	// -- Phase 2: Network attacker tampers with messages --
	fmt.Println("\n-- Phase 2: Network Attacker (tampering in transit) --")
	fmt.Println()

	keys2 := crypto.NewKeyRegistry()
	keys2.RegisterProcess("alice")
	keys2.RegisterProcess("bob")

	fll2 := link.NewFairLossLink(0.0)
	sl2 := link.NewStubbornLink(fll2, 5)
	pl2 := link.NewPerfectLink(sl2)
	al2 := link.NewAuthenticatedLink(pl2, keys2)

	fm2 := failures.NewNoFailure()
	rt2 := simrt.NewRuntime(99, al2, fm2, 20)

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

	fmt.Printf("\n  Phase 2 Result: Bob received %d message(s)\n", receiver2.ReceivedCount)
	if receiver2.ReceivedCount == 0 {
		fmt.Println("  -> Tampered message REJECTED by authenticated link!")
		fmt.Println("  -> MAC verification failed (content was modified).")
		fmt.Println("  -> This is the AL1 (Authenticity) property.")
	} else {
		fmt.Println("  -> WARNING: Message delivered despite tampering.")
	}
	fmt.Println()
}
