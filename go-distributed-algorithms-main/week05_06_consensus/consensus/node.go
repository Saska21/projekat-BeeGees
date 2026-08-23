package consensus

import (
	"context"
	"fmt"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// CONSENSUS NODE
// ============================================================================
//
// ConsensusNode implements Algorithm 5.1 (Flooding Consensus) from
// Cachin, Guerraoui, Rodrigues Chapter 5.1.2.
//
// Each node is a concurrent process that satisfies process.Process.
// All algorithm state lives inside the Run goroutine, so no internal
// mutexes are needed for algorithm logic.
//
// State mapping to book notation:
//
//	correct       map[ProcessID]bool       — set of processes not detected crashed
//	round         int                      — current round number
//	decision      ProposalValue            — ⊥ until decided ("" = undecided)
//	receivedfrom  map[int]map[ProcessID]bool — receivedfrom[r] from book
//	proposals     map[int]map[Value]bool   — proposals[r] from book
//
// Termination: the node returns from Run when ctx is cancelled.
type ConsensusNode struct {
	id       process.ProcessID   // this node's ID
	allIDs   []process.ProcessID // all process IDs in the system (sorted)
	recorder *DecisionRecorder   // shared result recorder for checkers
}

// NewConsensusNode creates a new consensus node.
//
//   - id:       this process's ID
//   - allIDs:   all process IDs in the system (including this one)
//   - recorder: shared DecisionRecorder used by checkers
func NewConsensusNode(
	id process.ProcessID,
	allIDs []process.ProcessID,
	recorder *DecisionRecorder,
) *ConsensusNode {
	return &ConsensusNode{
		id:       id,
		allIDs:   allIDs,
		recorder: recorder,
	}
}

// ID implements process.Process.
func (n *ConsensusNode) ID() process.ProcessID { return n.id }

// ============================================================================
// RUN — main goroutine (all algorithm state lives here)
// ============================================================================

// Run implements process.Process. It is launched as a goroutine by the runtime.
//
// State lifecycle matches Algorithm 5.1 Init event:
//
//	correct        := all process IDs
//	round          := 1
//	decision       := ⊥
//	receivedfrom   := [∅...], receivedfrom[0] := all process IDs
//	proposals      := [∅...]
func (n *ConsensusNode) Run(
	ctx context.Context,
	inbox <-chan process.Message,
	send func(process.Message),
) {
	// --- Algorithm 5.1 Init ---

	// correct: set of processes not yet detected as crashed
	correct := make(map[process.ProcessID]bool, len(n.allIDs))
	for _, pid := range n.allIDs {
		correct[pid] = true
	}

	// round: current round number (starts at 1)
	round := 1

	// decision: ⊥ until decided ("" == undecided)
	var decision ProposalValue

	// receivedfrom[r] = set of processes from which we received a PROPOSAL in round r.
	// receivedfrom[0] is initialised to all processes (book: receivedfrom[0] := Π)
	receivedfrom := make(map[int]map[process.ProcessID]bool)
	receivedfrom[0] = make(map[process.ProcessID]bool, len(n.allIDs))
	for _, pid := range n.allIDs {
		receivedfrom[0][pid] = true
	}

	// proposals[r] = accumulated proposal-value set for round r.
	proposals := make(map[int]map[ProposalValue]bool)
	proposals[1] = make(map[ProposalValue]bool)

	// ensureRound creates empty sets for round r if they don't exist yet.
	ensureRound := func(r int) {
		if receivedfrom[r] == nil {
			receivedfrom[r] = make(map[process.ProcessID]bool)
		}
		if proposals[r] == nil {
			proposals[r] = make(map[ProposalValue]bool)
		}
	}

	// --- maybeAdvanceOrDecide ---
	//
	// Called whenever correct or receivedfrom[round] changes.
	// Implements the book trigger:
	//
	//   upon correct ⊆ receivedfrom[round] ∧ decision = ⊥ do
	//     if receivedfrom[round] = receivedfrom[round-1] then decide
	//     else advance round
	var maybeAdvanceOrDecide func()
	maybeAdvanceOrDecide = func() {
		// Can only fire when undecided
		if decision != "" {
			return
		}

		ensureRound(round)

		// Check: correct ⊆ receivedfrom[round]
		for pid := range correct {
			if !receivedfrom[round][pid] {
				return // still waiting for at least one correct process
			}
		}

		// All correct processes have been heard from in this round.
		prevReceived := receivedfrom[round-1]
		currReceived := receivedfrom[round]

		if ProcessSetEqual(currReceived, prevReceived) {
			// Same set of senders as previous round → safe to decide.
			// decision := min(proposals[round])
			decision = MinValue(proposals[round])
			fmt.Printf("  [%s] *** DECIDED value=%q in round=%d proposals=%s ***\n",
				n.id, decision, round, ValuesString(proposals[round]))

			// Record in the shared result store for checker inspection.
			n.recorder.RecordDecision(n.id, decision)

			// Broadcast DECIDED to all processes (best-effort broadcast).
			Broadcast(send, n.id, n.allIDs, ConsensusMessage{
				Kind:  KindDecided,
				From:  n.id,
				Value: decision,
			})
		} else {
			// Different set → some process crashed in this round.
			// Advance to next round and rebroadcast accumulated proposals.
			round++
			ensureRound(round)

			// proposals[round] := copy of proposals[round-1]
			proposals[round] = CopyValues(proposals[round-1])

			fmt.Printf("  [%s] Advancing to round=%d proposals=%s\n",
				n.id, round, ValuesString(proposals[round]))

			// Broadcast our proposals for the new round.
			Broadcast(send, n.id, n.allIDs, ConsensusMessage{
				Kind:   KindProposal,
				Round:  round,
				From:   n.id,
				Values: SortedValues(proposals[round]),
			})

			// Recurse: maybe we can decide immediately in the new round.
			// (This handles the case where receivedfrom[round] is already
			// satisfied because we received messages ahead of advancing.)
			maybeAdvanceOrDecide()
		}
	}

	// --- Main message loop ---
	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-inbox:
			if !ok {
				return
			}

			cm, ok := msg.Data.(ConsensusMessage)
			if !ok {
				continue // ignore non-consensus messages
			}

			switch cm.Kind {

			// ------------------------------------------------------------------
			// upon event ⟨ c, Propose | v ⟩ do
			//   proposals[1] := proposals[1] ∪ {v}
			//   broadcast PROPOSAL(1, proposals[1])
			// ------------------------------------------------------------------
			case KindPropose:
				v := cm.Value
				fmt.Printf("  [%s] Received INIT_PROPOSE value=%q — initialising round 1\n", n.id, v)
				n.recorder.RecordProposal(n.id, v)

				proposals[1][v] = true

				Broadcast(send, n.id, n.allIDs, ConsensusMessage{
					Kind:   KindProposal,
					Round:  1,
					From:   n.id,
					Values: SortedValues(proposals[1]),
				})

			// ------------------------------------------------------------------
			// upon event ⟨ beb, Deliver | p, [PROPOSAL, r, ps] ⟩ do
			//   receivedfrom[r] := receivedfrom[r] ∪ {p}
			//   proposals[r]   := proposals[r]   ∪ ps
			//   (then check termination condition)
			// ------------------------------------------------------------------
			case KindProposal:
				r := cm.Round
				sender := cm.From
				ensureRound(r)

				receivedfrom[r][sender] = true
				for _, v := range cm.Values {
					proposals[r][v] = true
				}

				fmt.Printf("  [%s] Received PROPOSAL(round=%d, from=%s, values=%v) — known=%s\n",
					n.id, r, sender, cm.Values, ValuesString(proposals[r]))

				// Only fire termination logic for the current round.
				if r == round {
					maybeAdvanceOrDecide()
				}

			// ------------------------------------------------------------------
			// upon event ⟨ P, Crash | p ⟩ do
			//   correct := correct \ {p}
			//   (then check termination condition, because removing a crashed
			//    process might unblock the correct ⊆ receivedfrom[round] check)
			// ------------------------------------------------------------------
			case KindCrashDetected:
				crashed := cm.Crashed
				if correct[crashed] {
					delete(correct, crashed)
					fmt.Printf("  [%s] Detected crash of %s — correct=%v\n",
						n.id, crashed, sortedIDs(correct))
					n.recorder.MarkCrashed(crashed)
					maybeAdvanceOrDecide()
				}

			// ------------------------------------------------------------------
			// upon event ⟨ beb, Deliver | p, [DECIDED, v] ⟩
			//   such that p ∈ correct ∧ decision = ⊥ do
			//   decision := v
			//   broadcast DECIDED(v)
			//   trigger ⟨ c, Decide | v ⟩
			// ------------------------------------------------------------------
			case KindDecided:
				sender := cm.From
				v := cm.Value
				if correct[sender] && decision == "" {
					decision = v
					fmt.Printf("  [%s] Received DECIDED(from=%s, value=%q) — adopting decision\n",
						n.id, sender, v)
					n.recorder.RecordDecision(n.id, decision)

					// Rebroadcast DECIDED so late processes also decide.
					Broadcast(send, n.id, n.allIDs, ConsensusMessage{
						Kind:  KindDecided,
						From:  n.id,
						Value: decision,
					})
				}
			}
		}
	}
}

// sortedIDs returns a sorted slice of process IDs from a set, for printing.
func sortedIDs(m map[process.ProcessID]bool) []process.ProcessID {
	out := make([]process.ProcessID, 0, len(m))
	for pid := range m {
		out = append(out, pid)
	}
	// sort by string value
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
