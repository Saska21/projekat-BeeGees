package link

import (
	"github.com/danilokacanski/da/week03_04_parallel/crypto"
	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// ============================================================================
// AUTHENTICATED PERFECT LINK (AL)
// ============================================================================

// AuthenticatedLink models the Authenticated Perfect Link abstraction.
//
// From Cachin et al., Section 2.5:
//
//	AL1 (Authenticity): If a correct process q delivers a message m
//	  with sender p, and p is correct, then m was previously sent by p to q.
//
// Implementation:
//   - On send: attach a MAC (Message Authentication Code) to the message.
//   - On receive: verify the MAC before delivery.
type AuthenticatedLink struct {
	underlying *PerfectLink
	keys       *crypto.KeyRegistry
}

// NewAuthenticatedLink creates an authenticated link.
func NewAuthenticatedLink(pl *PerfectLink, keys *crypto.KeyRegistry) *AuthenticatedLink {
	return &AuthenticatedLink{
		underlying: pl,
		keys:       keys,
	}
}

// Send attaches a MAC to the message, then forwards via perfect link.
func (l *AuthenticatedLink) Send(msg process.Message) []process.Message {
	authenticated := msg.Clone()
	content := msg.ContentString()
	tag := l.keys.Sign(string(msg.From), []byte(content))
	authenticated.Meta["mac"] = tag
	return l.underlying.Send(authenticated)
}

// Receive authenticates FIRST, then deduplicates.
//
// This order is critical: if we deduplicated first, a tampered message
// would burn the dedup slot and prevent the valid retransmission from
// ever being delivered.
func (l *AuthenticatedLink) Receive(msg process.Message) (process.Message, bool) {
	// Step 1: Verify MAC BEFORE deduplication.
	tag, exists := msg.Meta["mac"].(string)
	if !exists {
		return msg, false // No MAC — reject (possible forgery)
	}

	content := msg.ContentString()
	valid := l.keys.Verify(string(msg.From), []byte(content), tag)
	if !valid {
		return msg, false // Invalid MAC — reject (tampered message)
	}

	// Step 2: Deduplication via perfect link (only for authentic messages).
	msg, ok := l.underlying.Receive(msg)
	if !ok {
		return msg, false // Already delivered — suppress duplicate
	}

	return msg, true
}

// Retransmissions delegates to the underlying perfect link.
func (l *AuthenticatedLink) Retransmissions() []process.Message {
	return l.underlying.Retransmissions()
}
