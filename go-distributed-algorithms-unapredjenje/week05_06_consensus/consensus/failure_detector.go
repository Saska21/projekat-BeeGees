package consensus

import (
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"
)

// ============================================================================
// PERFECT FAILURE DETECTOR (simplified)
// ============================================================================
//
// The book's PerfectFailureDetector (P) satisfies:
//   - PFD1 (Strong completeness): every crashed process is eventually detected.
//   - PFD2 (Strong accuracy):     no correct process is ever detected as crashed.
//
// In a real system this requires timeouts + heartbeats over a synchronous network.
//
// Simplification for this teaching implementation:
//   - We use an explicit "oracle" that knows in advance which processes will
//     crash and at what simulated time, and injects CRASH_DETECTED messages.
//   - This is conceptually equivalent to the PFD output events in the algorithm.
//   - It avoids the need for a full heartbeat protocol and keeps examples clean.

// CrashEvent describes one planned crash injection.
type CrashEvent struct {
	Crashed  process.ProcessID // which process crashes
	NotifyAt time.Duration     // how long after Run() to inject CRASH_DETECTED
}

// PerfectFailureDetector sends CRASH_DETECTED messages to all surviving processes
// after a configurable delay, simulating the output of a PFD oracle.
//
// Usage:
//
//	pfd := NewPerfectFailureDetector(rt, allIDs, []CrashEvent{
//	    {Crashed: "node-D", NotifyAt: 50 * time.Millisecond},
//	})
//	pfd.Start(crashFailure)   // also marks the process as crashed in the runtime
type PerfectFailureDetector struct {
	rt      *simrt.Runtime
	allIDs  []process.ProcessID
	crashes []CrashEvent
}

// NewPerfectFailureDetector creates a PFD oracle.
func NewPerfectFailureDetector(
	rt *simrt.Runtime,
	allIDs []process.ProcessID,
	crashes []CrashEvent,
) *PerfectFailureDetector {
	return &PerfectFailureDetector{rt: rt, allIDs: allIDs, crashes: crashes}
}

// Inject fires all crash events.  Call this in a goroutine BEFORE rt.Run().
// Each crash event:
//  1. Waits NotifyAt duration.
//  2. Crashes the process in the runtime failure model (so it stops receiving).
//  3. Injects CRASH_DETECTED messages to every surviving process.
func (pfd *PerfectFailureDetector) Inject(crashFM interface {
	CrashProcess(process.ProcessID)
}) {
	for _, ev := range pfd.crashes {
		// Fire each event in its own goroutine so events with different
		// delays don't block each other.
		ev := ev
		go func() {
			time.Sleep(ev.NotifyAt)

			// Mark process as crashed in the runtime failure model.
			crashFM.CrashProcess(ev.Crashed)
			pfd.rt.CrashProcess(ev.Crashed)

			// Notify every other process (simulate PFD Crash | p output).
			for _, to := range pfd.allIDs {
				if to == ev.Crashed {
					continue // don't notify the crashed process itself
				}
				msg := process.NewMessage(
					"pfd", // logical sender: failure detector
					to,
					ConsensusType,
					ConsensusMessage{
						Kind:    KindCrashDetected,
						Crashed: ev.Crashed,
					},
				)
				pfd.rt.Inject(msg)
			}
		}()
	}
}
