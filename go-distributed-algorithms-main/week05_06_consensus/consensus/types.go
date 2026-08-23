// Package consensus implements Flooding Consensus from
// Cachin, Guerraoui, Rodrigues: "Reliable and Secure Distributed Programming"
// Chapter 5.1.2, Algorithm 5.1.
package consensus

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// CORE VALUE TYPE
// ============================================================================

// ProposalValue is the type of values processes propose and decide.
// We use strings so values are easy to print and compare.
type ProposalValue string

// MinValue returns the minimum (lexicographically smallest) value from a set.
// This is the deterministic selection function used in the book algorithm.
func MinValue(vals map[ProposalValue]bool) ProposalValue {
	sorted := SortedValues(vals)
	if len(sorted) == 0 {
		return ""
	}
	return sorted[0]
}

// SortedValues returns the values in a set as a sorted slice.
// Used for deterministic trace output and message payloads.
func SortedValues(vals map[ProposalValue]bool) []ProposalValue {
	out := make([]ProposalValue, 0, len(vals))
	for v := range vals {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// ValuesString renders a proposal set as a human-readable string.
func ValuesString(vals map[ProposalValue]bool) string {
	sv := SortedValues(vals)
	strs := make([]string, len(sv))
	for i, v := range sv {
		strs[i] = string(v)
	}
	return "{" + strings.Join(strs, ",") + "}"
}

// CopyValues copies a set of ProposalValues.
func CopyValues(src map[ProposalValue]bool) map[ProposalValue]bool {
	dst := make(map[ProposalValue]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// CopyProcessSet copies a set of process IDs.
func CopyProcessSet(src map[process.ProcessID]bool) map[process.ProcessID]bool {
	dst := make(map[process.ProcessID]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// ProcessSetEqual returns true if two process-ID sets are identical.
func ProcessSetEqual(a, b map[process.ProcessID]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// ============================================================================
// CONSENSUS MESSAGE KINDS
// ============================================================================

const (
	// KindPropose is the initial INIT message sent to a node to start the protocol.
	KindPropose = "INIT_PROPOSE"

	// KindProposal carries a round number and proposal set (BestEffortBroadcast).
	KindProposal = "PROPOSAL"

	// KindDecided carries the decided value (BestEffortBroadcast after decision).
	KindDecided = "DECIDED"

	// KindCrashDetected carries the ID of a process detected as crashed.
	// This simulates output events from the PerfectFailureDetector.
	KindCrashDetected = "CRASH_DETECTED"

	// ConsensusType is the process.Message.Type value used for all consensus messages.
	ConsensusType = "CONSENSUS"
)

// ============================================================================
// CONSENSUS MESSAGE PAYLOAD
// ============================================================================

// ConsensusMessage is the structured payload carried in every consensus
// process.Message. It encodes all event types from Algorithm 5.1.
type ConsensusMessage struct {
	Kind    string            // KindPropose, KindProposal, KindDecided, KindCrashDetected
	Round   int               // Relevant for KindProposal
	From    process.ProcessID // Logical sender (useful for broadcast detection)
	Values  []ProposalValue   // Sorted proposal set (KindProposal)
	Value   ProposalValue     // Decided or proposed value (KindPropose, KindDecided)
	Crashed process.ProcessID // Crashed process ID (KindCrashDetected)
}

// String renders a ConsensusMessage for trace output.
func (m ConsensusMessage) String() string {
	switch m.Kind {
	case KindPropose:
		return fmt.Sprintf("INIT_PROPOSE(value=%s)", m.Value)
	case KindProposal:
		return fmt.Sprintf("PROPOSAL(round=%d, from=%s, values=%v)", m.Round, m.From, m.Values)
	case KindDecided:
		return fmt.Sprintf("DECIDED(from=%s, value=%s)", m.From, m.Value)
	case KindCrashDetected:
		return fmt.Sprintf("CRASH_DETECTED(crashed=%s)", m.Crashed)
	default:
		return fmt.Sprintf("UNKNOWN(%s)", m.Kind)
	}
}

// ValuesAsSet converts the Values slice back into a set for merging.
func (m ConsensusMessage) ValuesAsSet() map[ProposalValue]bool {
	s := make(map[ProposalValue]bool, len(m.Values))
	for _, v := range m.Values {
		s[v] = true
	}
	return s
}
