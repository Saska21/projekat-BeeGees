package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ============================================================================
// MESSAGE AUTHENTICATION CODES (MACs)
// ============================================================================

// A MAC is a short piece of information used to authenticate a message.
//
// From Cachin et al., Section 2.3:
//
//	A MAC takes a key k and a message m, and produces a tag t = MAC(k, m).
//	Anyone with key k can verify that the tag matches the message.
//	Without the key, it is infeasible to produce a valid tag.
//
// Properties:
//   - UNFORGEABILITY: Without the key, you cannot create a valid (message, tag) pair.
//   - VERIFICATION:   Given the key, message, and tag, you can check validity.
//
// Key difference from hashes:
//   - Hash = no key (anyone can compute)
//   - MAC  = shared key (only key holders can compute/verify)
//
// Implementation: HMAC-SHA256 (standard library).

// ComputeMAC computes an HMAC-SHA256 tag for the given data and key.
func ComputeMAC(key string, data []byte) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyMAC checks if a MAC tag is valid for the given data and key.
// Uses constant-time comparison to prevent timing attacks.
func VerifyMAC(key string, data []byte, tag string) bool {
	expected := ComputeMAC(key, data)
	return hmac.Equal([]byte(expected), []byte(tag))
}
