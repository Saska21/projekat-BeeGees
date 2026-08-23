package examples

import (
	"fmt"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"
	"github.com/danilokacanski/da/week05_06_consensus/consensus"
)

// ============================================================================
// EXAMPLE 3: Delayed Decision via DECIDED Rebroadcast
// ============================================================================
//
// Setup:
//   - 3 processes: node-A, node-B, node-C
//   - node-C sends its PROPOSAL(1) to node-A only, then crashes.
//     (same direct-inject technique as example 2)
//
// Why this specifically demonstrates DECIDED rebroadcast:
//
//	node-A: receives PROPOSAL from A, B, and injected C
//	  → receivedfrom[1] = {A,B,C} = receivedfrom[0] = {A,B,C}
//	  → decides round 1, broadcasts DECIDED("A")
//
//	node-B: receives PROPOSAL from A and B only (never gets C's)
//	  → receivedfrom[1] = {A,B}, blocked waiting for C
//	  → receives DECIDED("A") from node-A
//	  → adopts the decision WITHOUT completing round 2
//
// This is the second case from the book's correctness proof:
//
//	"p decides after delivering a DECIDED message from some process q."
//
// Contrast with example 2 (one_crash):
//
//	Example 2 has 4 nodes and B/C advance to round 2 before DECIDED arrives.
//	Example 3 has 3 nodes and B decides PURELY via the DECIDED relay —
//	no round 2 is needed at all.
func RunDelayedDecision() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 3: DELAYED DECISION — DECIDED rebroadcast")
	fmt.Println("  node-C sends to node-A only, then crashes.")
	fmt.Println("  node-A decides round 1; node-B decides via DECIDED relay.")
	fmt.Println("  Reference: Cachin et al., Algorithm 5.1 agreement proof")
	fmt.Println("================================================================")

	allIDs := []process.ProcessID{"node-A", "node-B", "node-C"}

	recorder := consensus.NewDecisionRecorder()

	fl := link.NewFairLossLink(0.0, 42)
	sl := link.NewStubbornLink(fl)
	pl := link.NewPerfectLink(sl)

	fm := failures.NewCrashFailure()

	rt := simrt.NewRuntime(pl, fm,
		simrt.WithMaxDuration(15*time.Second),
		simrt.WithIdleTimeout(800*time.Millisecond),
		simrt.WithRetransmitInterval(50*time.Millisecond),
		simrt.WithVerbose(false),
	)

	for _, pid := range allIDs {
		node := consensus.NewConsensusNode(pid, allIDs, recorder)
		rt.Register(node)
	}

	// Inject INIT_PROPOSE for A and B only.
	// C never receives INIT_PROPOSE so it never broadcasts via StubbornLink.
	for pid, val := range map[process.ProcessID]consensus.ProposalValue{
		"node-A": "A",
		"node-B": "B",
	} {
		rt.Inject(process.NewMessage(
			"system", pid,
			consensus.ConsensusType,
			consensus.ConsensusMessage{Kind: consensus.KindPropose, Value: val},
		))
		recorder.RecordProposal(pid, val)
	}

	// Simulate "node-C sends one message then crashes":
	// inject C's PROPOSAL directly into node-A's inbox only.
	// rt.Inject bypasses the link layer — no StubbornLink buffer, no retransmit.
	rt.Inject(process.NewMessage(
		"node-C", "node-A",
		consensus.ConsensusType,
		consensus.ConsensusMessage{
			Kind:   consensus.KindProposal,
			Round:  1,
			From:   "node-C",
			Values: []consensus.ProposalValue{"C"},
		},
	))
	recorder.RecordProposal("node-C", "C")

	// Crash C at 10ms.
	// node-A: receivedfrom[1]={A,B,C} = receivedfrom[0] → decides round 1.
	// node-B: receivedfrom[1]={A,B} ≠ {A,B,C} → blocked, then receives DECIDED from A.
	pfd := consensus.NewPerfectFailureDetector(rt, allIDs, []consensus.CrashEvent{
		{Crashed: "node-C", NotifyAt: 10 * time.Millisecond},
	})
	pfd.Inject(fm)

	rt.Run()

	snap := recorder.Snapshot()
	printDecisions(snap, allIDs)
	consensus.RunAllChecks(snap, allIDs)
}
