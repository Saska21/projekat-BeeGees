// Package process defines the core abstractions for distributed processes
// and message passing, following Chapter 2 of Cachin, Guerraoui, Rodrigues:
// "Reliable and Secure Distributed Programming".
//
// In a distributed system, processes communicate EXCLUSIVELY through messages.
// There is no shared memory between processes.
package process

import "fmt"

// ============================================================================
// PROCESS IDENTIFIER
// ============================================================================

// ProcessID uniquely identifies a process in the distributed system.
// In the formal model, processes are named p, q, r, etc.
// In our simulator, we use descriptive string identifiers like "p1", "alice".
type ProcessID string

// ============================================================================
// MESSAGE
// ============================================================================

// Message represents a message exchanged between processes.
//
// From Cachin et al., Section 2.1:
//   - Processes communicate by sending and receiving messages.
//   - A message is uniquely identified by its sender and content.
//   - Messages may contain arbitrary data.
//
// The Meta field carries auxiliary information used by link abstractions
// and cryptographic primitives (e.g., MAC tags, signatures, sequence numbers).
// It is NOT part of the "logical" message — it is metadata added by
// the communication infrastructure.
type Message struct {
	From ProcessID      // Sender process identifier
	To   ProcessID      // Destination process identifier
	Type string         // Message type label (e.g., "PING", "PONG", "ACK")
	Data any            // Payload data (arbitrary)
	Meta map[string]any // Metadata for links/crypto (tags, signatures, etc.)
}

// NewMessage creates a new Message with an initialized metadata map.
func NewMessage(from, to ProcessID, msgType string, data any) Message {
	return Message{
		From: from,
		To:   to,
		Type: msgType,
		Data: data,
		Meta: make(map[string]any),
	}
}

// String returns a human-readable representation of the message.
// Useful for trace logging and debugging.
func (m Message) String() string {
	return fmt.Sprintf("[%s -> %s] %s: %v", m.From, m.To, m.Type, m.Data)
}

// ContentString returns a canonical string of the message content,
// EXCLUDING metadata. Used for MAC/signature computation.
// Two messages with the same content but different metadata produce
// the same ContentString.
func (m Message) ContentString() string {
	return fmt.Sprintf("%s:%s:%s:%v", m.From, m.To, m.Type, m.Data)
}

// Clone creates a deep copy of the message.
// Important for link abstractions that retransmit or duplicate messages —
// each copy must be independent so modifications don't leak between them.
func (m Message) Clone() Message {
	meta := make(map[string]any, len(m.Meta))
	for k, v := range m.Meta {
		meta[k] = v
	}
	return Message{
		From: m.From,
		To:   m.To,
		Type: m.Type,
		Data: m.Data,
		Meta: meta,
	}
}
