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
// EXAMPLE 2: One Crash (round 2 required)
// ============================================================================
//
// Setup:
//   - 4 processes: node-A, node-B, node-C, node-D
//   - node-D crashes before sending PROPOSAL to anyone.
//     Its INIT_PROPOSE is not injected so it never broadcasts at all.
//
// Expected outcome:
//
//	Round 1:
//	  A, B, C broadcast PROPOSAL(1) and collect from each other.
//	  receivedfrom[1] = {A,B,C} for all — blocked waiting for D.
//	  PFD fires at 10ms → correct = {A,B,C}.
//	  receivedfrom[1] = {A,B,C} ≠ receivedfrom[0] = {A,B,C,D} → ALL advance to round 2.
//
//	Round 2:
//	  All three broadcast PROPOSAL(2, {A,B,C}) and collect from each other.
//	  receivedfrom[2] = {A,B,C} = receivedfrom[1] → ALL decide min({A,B,C}) = "A".
//
// Contrast with example 3:
//
//	Here nobody decides in round 1 — the DECIDED rebroadcast path is never taken.
//	All surviving nodes go through round 2 together because none of them
//	heard from D. Example 3 shows what happens when one node hears a partial
//	message from the crashing process and can decide early.
//
// Properties demonstrated:
//   - When a crash is detected and receivedfrom sets differ from the previous
//     round, an extra round is always required before deciding.
//   - All correct nodes advance to round 2 in lockstep.
func RunOneCrash() {
	fmt.Println("================================================================")
	fmt.Println("  EXAMPLE 2: ONE CRASH — Flooding Consensus (4 nodes, 1 crash)")
	fmt.Println("  node-D crashes before sending to anyone → all advance to round 2.")
	fmt.Println("  Reference: Cachin et al., Algorithm 5.1 / Fig. 5.1")
	fmt.Println("================================================================")

	allIDs := []process.ProcessID{"node-A", "node-B", "node-C", "node-D"}

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

	// Register all four nodes so every node starts with
	// correct = {A,B,C,D} and receivedfrom[0] = {A,B,C,D}.
	for _, pid := range allIDs {
		node := consensus.NewConsensusNode(pid, allIDs, recorder)
		rt.Register(node)
	}

	// Inject INIT_PROPOSE for A, B, C only.
	// D never receives INIT_PROPOSE and never sends PROPOSAL(1) to anyone.
	for pid, val := range map[process.ProcessID]consensus.ProposalValue{
		"node-A": "A",
		"node-B": "B",
		"node-C": "C",
	} {
		rt.Inject(process.NewMessage(
			"system", pid,
			consensus.ConsensusType,
			consensus.ConsensusMessage{Kind: consensus.KindPropose, Value: val},
		))
		recorder.RecordProposal(pid, val)
	}

	// Crash D at 10ms. By then A, B, C have exchanged PROPOSAL(1) but are
	// still blocked waiting for D. Crash fires → all advance to round 2.
	pfd := consensus.NewPerfectFailureDetector(rt, allIDs, []consensus.CrashEvent{
		{Crashed: "node-D", NotifyAt: 10 * time.Millisecond},
	})
	pfd.Inject(fm)

	rt.Run()

	snap := recorder.Snapshot()
	printDecisions(snap, allIDs)
	consensus.RunAllChecks(snap, allIDs)
}
