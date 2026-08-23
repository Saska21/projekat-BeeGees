package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ============================================================================
// MESSAGE AUTHENTICATION CODES (MACs)
// ============================================================================

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
