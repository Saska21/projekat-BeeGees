// Package crypto provides simplified cryptographic abstractions for
// the distributed algorithms simulator.
//
// From Cachin et al., Section 2.3:
//
//	Cryptographic abstractions provide three fundamental services:
//	  1. INTEGRITY (hash functions): detect if data was modified
//	  2. AUTHENTICATION (MACs): verify the sender of a message
//	  3. NON-REPUDIATION (signatures): prove a message was sent by a specific process
//
// This implementation uses Go's standard library (crypto/sha256, crypto/hmac).
package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ============================================================================
// HASH FUNCTIONS
// ============================================================================

// Hash computes a SHA-256 hash of raw bytes.
// Returns a hex-encoded string.
func Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// HashString computes a SHA-256 hash of a string.
func HashString(s string) string {
	return Hash([]byte(s))
}

// HashMessage computes a hash over the logical content of a message.
func HashMessage(from, to, msgType string, data any) string {
	content := fmt.Sprintf("%s:%s:%s:%v", from, to, msgType, data)
	return Hash([]byte(content))
}

// VerifyHash checks if data matches an expected hash.
func VerifyHash(data []byte, expectedHash string) bool {
	return Hash(data) == expectedHash
}
