# Week 3-4: Basic Abstractions

A single-threaded, deterministic simulator for the foundational abstractions
of distributed computing, based on **Chapter 2** of:

> **Cachin, Guerraoui, Rodrigues** — *Reliable and Secure Distributed Programming* (2nd ed.)

## Overview

This codebase teaches the core building blocks of every distributed algorithm:

| Concept | Package | Book Reference |
|---------|---------|----------------|
| Process abstraction | `process/` | Section 2.1 |
| Message passing | `process/message.go` | Section 2.1 |
| Link abstractions | `link/` | Section 2.4 |
| Failure models | `failures/` | Section 2.2 |
| Cryptographic primitives | `crypto/` | Section 2.5 |
| Nondeterministic execution | `runtime/` | Section 2.1 |

## Quick Start

```bash
cd week03_04_basic_abstractions
go run .
```

This runs three self-contained examples:

1. **Ping-Pong** — Basic message passing with perfect links
2. **Counter** — Message loss with fair-loss links (50% loss rate)
3. **Tamper** — Authenticated links detecting message forgery

## Architecture

### No Goroutines. No Channels. No Networking.

This simulator is **intentionally single-threaded**. All nondeterminism
comes from `math/rand` with an explicit seed, making every execution
**reproducible**.

```
┌─────────────────────────────────────────────────────┐
│                    Runtime Scheduler                │
│                                                     │
│  1. Pick random message from queue                  │
│  2. Apply failure model (crash? omission?)          │
│  3. Apply link receive (dedup? auth check?)         │
│  4. Deliver to process → Handle() → outgoing msgs   │
│  5. Apply link send (loss? attach MAC?)             │
│  6. Enqueue surviving messages                      │
│  7. Log everything, check properties                │
└─────────────────────────────────────────────────────┘
```

### Package Dependency Graph

```
process  ←──  crypto  ←──  link
   ↑             ↑           ↑
   │             │           │
failures    (standalone)     │
   ↑                         │
   └────────── runtime ──────┘
                  ↑
               examples
                  ↑
                main
```

No circular imports. Each package has a focused responsibility.

## Key Concepts

### Process Abstraction (Section 2.1)

Every participant is a **state machine** that implements one method:

```go
type Process interface {
    ID() ProcessID
    Handle(msg Message) []Message
}
```

- `Handle` is **one atomic step** — receive a message, update state, send messages
- Processes NEVER send messages directly; they return them for the runtime to enqueue
- Processes NEVER access other processes' state

### Link Hierarchy (Section 2.4)

Links form a composition hierarchy, where each level adds guarantees:

```
Fair-Loss → Stubborn → Perfect → Authenticated
    ↑           ↑          ↑           ↑
   FLL1-3     SL1-2      PL1-3        AL1
```

| Link | Key Property | Mechanism |
|------|-------------|-----------|
| Fair-Loss | Eventually delivers if sent enough | Probabilistic delivery |
| Stubborn | Delivers infinitely often | Retransmission buffer |
| Perfect | Reliable + no duplication | Stubborn + dedup set |
| Authenticated | No forgery | Perfect + MAC verification |

**Key insight**: Each link is built by **composing** the one below it:
- `Perfect = Stubborn + Deduplication`
- `Stubborn = FairLoss + Retransmission`
- `Authenticated = Perfect + MAC`

### Failure Models (Section 2.2)

| Model | Behavior | Strength |
|-------|----------|----------|
| No failure | All processes correct | Weakest |
| Crash-stop | Process halts permanently | ⬆ |
| Omission | Messages silently dropped | ⬆ |
| Byzantine | Arbitrary behavior | Strongest |

### Cryptographic Primitives (Section 2.5)

| Primitive | Purpose | Implementation |
|-----------|---------|---------------|
| Hash (SHA-256) | Integrity | `crypto/sha256` |
| MAC (HMAC-SHA256) | Authentication | `crypto/hmac` |
| Signatures | Non-repudiation | Simplified HMAC model |

## Safety vs Liveness

Every property of a distributed algorithm is either:

- **Safety**: "Something bad NEVER happens" — checked after every step
- **Liveness**: "Something good EVENTUALLY happens" — checked at the end

The `runtime/checker.go` package provides both types, and the runtime
runs them automatically during simulation.

Examples:
- **Safety**: "No message is delivered more than once" (PL2)
- **Liveness**: "Every sent message is eventually delivered" (PL1)

## Nondeterminism and Reproducibility

In a real distributed system, message delivery order is unpredictable.
Our simulator models this with **random scheduling**:

- Each step picks a **random** message from the queue
- The random seed determines the entire execution
- **Same seed → same execution** (reproducible for debugging)
- **Different seed → different execution** (explore different interleavings)

## File Structure

```
week03_04_basic_abstractions/
├── main.go                    # Entry point — runs all examples
├── go.mod                     # Module definition
├── README.md                  # This file
├── process/
│   ├── message.go             # ProcessID, Message type
│   ├── process.go             # Process interface
│   └── state.go               # State helper (key-value store)
├── crypto/
│   ├── hash.go                # SHA-256 hashing
│   ├── mac.go                 # HMAC-SHA256
│   └── signature.go           # Simplified digital signatures
├── link/
│   ├── fairloss.go            # Link interface + Fair-Loss Link (FLL1-3)
│   ├── stubborn.go            # Stubborn Link (SL1-2, Algorithm 2.1)
│   ├── perfect.go             # Perfect Link (PL1-3, Algorithm 2.2)
│   └── authenticated.go       # Authenticated Perfect Link (AL1)
├── failures/
│   ├── crash.go               # FailureInjector interface + crash-stop
│   ├── omission.go            # Omission failures
│   └── byzantine.go           # Byzantine (arbitrary) failures
├── runtime/
│   ├── runtime.go             # Central scheduler
│   ├── trace.go               # Event logging
│   └── checker.go             # Safety/liveness property checkers
└── examples/
    ├── pingpong.go            # Basic message passing demo
    ├── counter.go             # Message loss demo
    └── tamper.go              # Authentication + forgery demo
```

## Exercises for Students

1. **Change the seed**: Run with different seeds and observe different delivery orders
2. **Increase loss rate**: Change the counter example's loss rate to 0.8 — what happens?
3. **Add crash failure**: Crash the pong process mid-execution and observe the trace
4. **Implement a new process**: Create a broadcast process that sends to all peers
5. **Build a reliable broadcast**: Use perfect links to implement reliable broadcast
6. **Byzantine sender**: Make a worker process Byzantine and verify the counter bound still holds

## Course Context

This simulator covers **Weeks 3-4** of the Distributed Algorithms course:
- Week 3: Process, link, and failure abstractions
- Week 4: Cryptographic primitives and authenticated communication

The abstractions implemented here are the **foundation** for all subsequent
algorithms in the course (consensus, broadcast, replication, etc.).
