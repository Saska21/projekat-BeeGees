package runtime

import (
	"fmt"
	"sync"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// EventType classifies trace events for analysis and display.
type EventType int

const (
	EventSend        EventType = iota // A process produces an outgoing message.
	EventEnqueue                      // A message is placed in the network channel.
	EventDequeue                      // A message is selected from the network channel.
	EventDeliver                      // A message is delivered to its recipient.
	EventDrop                         // A message is dropped.
	EventRetransmit                   // A stubborn link retransmits a message.
	EventStateChange                  // A notable state change in the system.
	EventCrash                        // A process crashes.
	EventRecover                      // A process recovers after a crash.
	EventForge                        // A message is forged or tampered with.
	EventViolation                    // A checker detected a property violation.
)

var eventLabels = map[EventType]string{
	EventSend:        "SEND",
	EventEnqueue:     "ENQUEUE",
	EventDequeue:     "DEQUEUE",
	EventDeliver:     "DELIVER",
	EventDrop:        "DROP",
	EventRetransmit:  "RETRANSMIT",
	EventStateChange: "STATE",
	EventCrash:       "CRASH",
	EventRecover:     "RECOVER",
	EventForge:       "FORGE",
	EventViolation:   "VIOLATION",
}

// EventLabel returns the human-readable label for an event type.
func EventLabel(t EventType) string {
	if label, ok := eventLabels[t]; ok {
		return label
	}
	return "UNKNOWN"
}

// Event represents a single observable action in the execution.
type Event struct {
	Seq     int64             // Monotonically increasing sequence number
	Time    time.Time         // Wall-clock time of the event
	Type    EventType         // Classification of the event
	Message *process.Message  // Associated message (nil for non-message events)
	Process process.ProcessID // Relevant process identifier
	Detail  string            // Human-readable description
}

// String returns a formatted representation of the event.
func (e Event) String() string {
	return fmt.Sprintf("[Seq %04d] %-10s %s", e.Seq, EventLabel(e.Type), e.Detail)
}

// ============================================================================
// TRACE (thread-safe)
// ============================================================================

// Trace records the complete sequence of events in an execution.
// Thread-safe via sync.Mutex.
type Trace struct {
	mu      sync.Mutex
	events  []Event
	verbose bool
	seq     int64 // atomic counter for sequence numbers
}

// NewTrace creates a new empty trace.
func NewTrace(verbose bool) *Trace {
	return &Trace{
		events:  make([]Event, 0),
		verbose: verbose,
	}
}

// Log records an event and optionally prints it.
// Sequence numbers are assigned inside the mutex so events are
// appended in strictly increasing sequence order.
func (t *Trace) Log(e Event) {
	t.mu.Lock()
	t.seq++
	e.Seq = t.seq
	e.Time = time.Now()
	t.events = append(t.events, e)
	t.mu.Unlock()
	if t.verbose {
		fmt.Println(e.String())
	}
}

// Events returns a snapshot of all recorded events.
func (t *Trace) Events() []Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	snapshot := make([]Event, len(t.events))
	copy(snapshot, t.events)
	return snapshot
}

// EventsOfType returns all events matching the given type.
func (t *Trace) EventsOfType(eventType EventType) []Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	var result []Event
	for _, e := range t.events {
		if e.Type == eventType {
			result = append(result, e)
		}
	}
	return result
}

// Summary prints a summary of the execution.
func (t *Trace) Summary() {
	t.mu.Lock()
	defer t.mu.Unlock()
	counts := make(map[EventType]int)
	for _, e := range t.events {
		counts[e.Type]++
	}
	fmt.Println("\n--- Execution Summary ---")
	fmt.Printf("Total events: %d\n", len(t.events))
	for eventType, label := range eventLabels {
		if count, ok := counts[eventType]; ok && count > 0 {
			fmt.Printf("  %-12s %d\n", label+":", count)
		}
	}
	fmt.Println("-------------------------")
}
