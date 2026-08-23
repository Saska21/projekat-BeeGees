// Package examples contains self-contained runnable scenarios for the
// Flooding Consensus simulator.
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
// EXAMPLE 1: No Failures
// ============================================================================
//
// Setup:
//   - 3 processes: node-A, node-B, node-C
//   - Proposed values: "A", "B", "C"
//   - No crashes
//
// Expected outcome:
//   - All three nodes decide in round 1 (same receivedfrom[1] == receivedfrom[0]).
//   - All decide min("A","B","C") = "A".
//
// Properties demonstrated:
//   - Termination in a single round when there are no failures.
//   - Agreement and Validity hold trivially.
func RunNoFailure() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 1: NO FAILURE — Flooding Consensus (3 nodes)")
	fmt.Println("  All nodes propose; all decide min value in round 1.")
	fmt.Println("  Reference: Cachin et al., Algorithm 5.1")
	fmt.Println("================================================================")

	allIDs := []process.ProcessID{"node-A", "node-B", "node-C"}
	proposals := map[process.ProcessID]consensus.ProposalValue{
		"node-A": "A",
		"node-B": "B",
		"node-C": "C",
	}

	recorder := consensus.NewDecisionRecorder()

	// Perfect link stack over a zero-loss fair-loss link.
	// (No message loss; we want all messages to arrive.)
	fl := link.NewFairLossLink(0.0, 42)
	sl := link.NewStubbornLink(fl)
	pl := link.NewPerfectLink(sl)

	fm := failures.NewNoFailure()

	rt := simrt.NewRuntime(pl, fm,
		simrt.WithMaxDuration(10*time.Second),
		simrt.WithIdleTimeout(500*time.Millisecond),
		simrt.WithRetransmitInterval(50*time.Millisecond),
		simrt.WithVerbose(false), // keep output focused on consensus prints
	)

	// Register nodes.
	for _, pid := range allIDs {
		node := consensus.NewConsensusNode(pid, allIDs, recorder)
		rt.Register(node)
	}

	// Inject INIT_PROPOSE messages so each node starts the algorithm.
	for pid, val := range proposals {
		rt.Inject(process.NewMessage(
			"system", pid,
			consensus.ConsensusType,
			consensus.ConsensusMessage{Kind: consensus.KindPropose, Value: val},
		))
		recorder.RecordProposal(pid, val)
	}

	rt.Run()

	snap := recorder.Snapshot()
	printDecisions(snap, allIDs)
	consensus.RunAllChecks(snap, allIDs)
}

// printDecisions prints a summary of what each process decided.
func printDecisions(s consensus.Snapshot, allIDs []process.ProcessID) {
	fmt.Println("\n--- Decision Summary ---")
	for _, pid := range allIDs {
		v, ok := s.Decisions[pid]
		if s.Crashed[pid] {
			fmt.Printf("  %s: CRASHED\n", pid)
		} else if ok {
			fmt.Printf("  %s: decided %q\n", pid, v)
		} else {
			fmt.Printf("  %s: DID NOT DECIDE\n", pid)
		}
	}
	fmt.Println("------------------------")
}
