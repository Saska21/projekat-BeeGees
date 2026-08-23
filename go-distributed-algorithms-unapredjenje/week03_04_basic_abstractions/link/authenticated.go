package link

import (
	"math/rand"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/crypto"
	"github.com/danilokacanski/da/week03_04_basic_abstractions/process"
)

// ============================================================================
// AUTHENTICATED PERFECT LINK (AL)
// ============================================================================

// AuthenticatedLink models the Authenticated Perfect Link abstraction.
//
// From Cachin et al., Section 2.5:
//
// Properties (in addition to PL1-PL3):
//
//	AL1 (Authenticity): If a correct process q delivers a message m
//	  with sender p, and p is correct, then m was previously sent
//	  by p to q.
//
// This prevents message forgery: a Byzantine process or network attacker
// cannot impersonate another process.
//
// Implementation:
//
//	Uses a Perfect Link underneath.
//	On send: attaches a MAC (Message Authentication Code) to the message.
//	On receive: verifies the MAC before delivery.
//
//	If the message was modified in transit (by a network attacker),
//	the MAC verification fails and the message is rejected.
type AuthenticatedLink struct {
	// underlying is the perfect link for reliable delivery.
	underlying *PerfectLink

	// keys is the key registry for MAC computation and verification.
	keys *crypto.KeyRegistry
}

// NewAuthenticatedLink creates an authenticated link.
func NewAuthenticatedLink(pl *PerfectLink, keys *crypto.KeyRegistry) *AuthenticatedLink {
	return &AuthenticatedLink{
		underlying: pl,
		keys:       keys,
	}
}

// Send attaches a MAC to the message, then forwards via perfect link.
//
// The MAC is computed over the message content (excluding metadata).
// This ensures that any modification to From, To, Type, or Data
// will be detected by the receiver.
func (l *AuthenticatedLink) Send(msg process.Message, rng *rand.Rand) []process.Message {
	// Clone to avoid modifying the original
	authenticated := msg.Clone()

	// Compute MAC over message content using sender's key
	content := msg.ContentString()
	tag := l.keys.Sign(string(msg.From), []byte(content))
	authenticated.Meta["mac"] = tag

	return l.underlying.Send(authenticated, rng)
}

// Tick delegates to the underlying perfect link.
func (l *AuthenticatedLink) Tick(step int, rng *rand.Rand) []process.Message {
	return l.underlying.Tick(step, rng)
}

// Receive verifies the MAC before delivering.
//
// Two checks happen in order:
//  1. Perfect link deduplication (have we seen this message before?)
//  2. MAC verification (is the message authentic?)
//
// If either check fails, the message is rejected.
func (l *AuthenticatedLink) Receive(msg process.Message) (process.Message, bool) {
	// First: deduplication via perfect link
	msg, ok := l.underlying.Receive(msg)
	if !ok {
		return msg, false // Duplicate — already delivered
	}

	// Second: verify if MAC exists
	tag, exists := msg.Meta["mac"].(string)
	if !exists {
		return msg, false // No MAC — reject (possible forgery)
	}
	// Third: verify MAC using sender's key
	content := msg.ContentString()
	valid := l.keys.Verify(string(msg.From), []byte(content), tag)
	if !valid {
		return msg, false // Invalid MAC — reject (tampered message)
	}

	return msg, true
}
