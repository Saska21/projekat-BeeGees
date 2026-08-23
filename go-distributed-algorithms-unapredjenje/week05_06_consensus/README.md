# Week 05-06: Flooding Consensus

A teaching implementation of **Algorithm 5.1 (Flooding Consensus)** from:

> Cachin, Guerraoui, Rodrigues — *Reliable and Secure Distributed Programming*, Chapter 5.1.2

---

## What is Flooding Consensus?

Flooding Consensus solves the **consensus problem** for fail-stop (crash-stop) processes. The consensus problem requires all correct processes to agree on a single value that was proposed by one of them.

### Consensus Properties

| Property | Description |
|---|---|
| **Termination** | Every correct process eventually decides a value. |
| **Validity** | The decided value was proposed by some process. |
| **Integrity** | Each process decides at most once. |
| **Agreement** | No two correct processes decide different values. |

---

## The Algorithm (Algorithm 5.1)

The key idea: processes execute sequential rounds, each time broadcasting the full set of proposals they know. A process can only decide when:

1. It has heard from every process it believes is still correct in this round.
2. The set of processes it heard from in this round **equals** the set from the previous round (no new crash detected).

If these conditions hold, the process decides `min(proposals[round])`.  
Otherwise it advances to the next round and floods again.

### Pseudocode (abbreviated)

```
Init:
  correct = all processes
  round = 1
  receivedfrom[0] = all processes

On Propose(v):
  proposals[1] += {v}
  broadcast PROPOSAL(1, proposals[1])

On receive PROPOSAL(r, ps) from p:
  receivedfrom[r] += {p}
  proposals[r]    += ps
  if correct ⊆ receivedfrom[round]:
    if receivedfrom[round] == receivedfrom[round-1]:
      decide min(proposals[round])
      broadcast DECIDED(decision)
    else:
      round++; broadcast PROPOSAL(round, proposals[round-1])

On detect crash of p:
  correct -= {p}
  re-check termination condition

On receive DECIDED(v) from correct p (and not yet decided):
  decide v; broadcast DECIDED(v)
```

---

## Abstractions and Simplifications

### BestEffortBroadcast (BEB)

The book uses BEB to flood proposals. Here we implement it as **sending one unicast message to every process** via the `send` callback. Reliability is provided by the `PerfectLink` layer from `week03_04_parallel`.

### PerfectFailureDetector (P)

The book requires a perfect failure detector satisfying:
- **Strong completeness**: every crashed process is eventually detected.
- **Strong accuracy**: no correct process is ever detected.

We simulate this with an **oracle PFD**: a goroutine that waits a configurable delay, then injects `CRASH_DETECTED` messages to all surviving processes. This is semantically equivalent to the PFD output events in the algorithm.

---

## Project Structure

```
week05_06_consensus/
├── main.go                     # Runs all three examples
├── go.mod                      # Module (depends on week03_04_parallel)
├── consensus/
│   ├── types.go                # ProposalValue, ConsensusMessage, helpers
│   ├── node.go                 # ConsensusNode (implements process.Process)
│   ├── broadcast.go            # Broadcast helper (BEB simplified)
│   ├── failure_detector.go     # PerfectFailureDetector oracle
│   └── checkers.go             # DecisionRecorder + 4 property checkers
└── examples/
    ├── no_failure.go           # 3 nodes, no crash, decide in round 1
    ├── one_crash.go            # 4 nodes, 1 crash, may need round 2
    └── delayed_decision.go     # DECIDED rebroadcast scenario
```

---

## How to Run

```bash
cd week05_06_consensus
go run .
```

---

## Examples

### Example 1: No Failure

- **3 nodes**: node-A, node-B, node-C
- **Proposals**: "A", "B", "C"
- **Expected**: all decide `"A"` in round 1

Demonstrates: when there are no failures, a single communication round suffices.

### Example 2: One Crash

- **4 nodes**: node-A, node-B, node-C, node-D
- **Proposals**: "A", "B", "C", "D"
- **node-D crashes at 60ms**
- **Expected**: correct nodes decide the minimum of proposals seen across all rounds

Demonstrates: why an additional round is required when a crash is detected mid-round.

### Example 3: Delayed Decision via DECIDED Rebroadcast

- **3 nodes**: node-A, node-B, node-C
- **Proposals**: "A", "B", "C"
- **node-C crashes at 40ms**
- **Expected**: node-A decides early and rebroadcasts DECIDED; node-B adopts the decision

Demonstrates: the DECIDED rebroadcast rule that ensures late-deciding processes still agree.

---

## Expected Output

Each example prints:
- Per-node consensus activity (PROPOSAL received, crash detected, round advance, decision)
- A decision summary table
- Property check results: **[PASS]** or **[FAIL]** for Agreement, Validity, Integrity, Termination

---

## Reuse from week03_04_parallel

| Component | Reused | Purpose |
|---|---|---|
| `runtime.Runtime` | ✅ | Execution harness (goroutines, router, idle monitor) |
| `process.Process` | ✅ | Interface implemented by `ConsensusNode` |
| `link.FairLossLink` | ✅ | Transport (0% loss for consensus examples) |
| `link.StubbornLink` | ✅ | Retransmission for reliability |
| `link.PerfectLink`  | ✅ | Deduplication |
| `failures.CrashFailure` | ✅ | Marks crashed processes in the runtime |
| `runtime.Trace` | ✅ | Event logging (verbose=false in examples) |
