// Package process defines the core abstractions for distributed processes
// and message passing, following Chapter 2 of Cachin, Guerraoui, Rodrigues:
// "Reliable and Secure Distributed Programming".
//
// PARALLEL VERSION: Processes run as goroutines and communicate via channels.
package process

import "fmt"

// ============================================================================
// PROCESS IDENTIFIER
// ============================================================================

// ProcessID uniquely identifies a process in the distributed system.
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
func (m Message) String() string {
	return fmt.Sprintf("[%s -> %s] %s: %v", m.From, m.To, m.Type, m.Data)
}

// ContentString returns a canonical string of the message content,
// EXCLUDING metadata. Used for MAC/signature computation.
func (m Message) ContentString() string {
	return fmt.Sprintf("%s:%s:%s:%v", m.From, m.To, m.Type, m.Data)
}

// Clone creates a deep copy of the message.
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
