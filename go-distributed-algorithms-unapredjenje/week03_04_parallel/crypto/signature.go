package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ============================================================================
// DIGITAL SIGNATURES (Simplified Model)
// ============================================================================

// KeyPair holds a public-private key pair for a process.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// KeyRegistry manages cryptographic keys for all processes.
// Thread-safe: reads happen from multiple goroutines but registration
// is done before the simulation starts, so no mutex is needed.
type KeyRegistry struct {
	keys       map[string]KeyPair // processID -> key pair
	sharedKeys map[string]string  // "p:q" -> shared secret
}

// NewKeyRegistry creates an empty key registry.
func NewKeyRegistry() *KeyRegistry {
	return &KeyRegistry{
		keys:       make(map[string]KeyPair),
		sharedKeys: make(map[string]string),
	}
}

// RegisterProcess generates and stores a key pair for a process.
func (kr *KeyRegistry) RegisterProcess(processID string) KeyPair {
	kp := KeyPair{
		PublicKey:  HashString("pub:" + processID),
		PrivateKey: HashString("priv:" + processID),
	}
	kr.keys[processID] = kp
	return kp
}

// RegisterSharedKey creates a shared secret between two processes.
func (kr *KeyRegistry) RegisterSharedKey(p, q string) string {
	key := HashString(fmt.Sprintf("shared:%s:%s", p, q))
	kr.sharedKeys[sharedKeyID(p, q)] = key
	kr.sharedKeys[sharedKeyID(q, p)] = key
	return key
}

// GetSharedKey retrieves the shared key between two processes.
func (kr *KeyRegistry) GetSharedKey(p, q string) (string, bool) {
	key, ok := kr.sharedKeys[sharedKeyID(p, q)]
	return key, ok
}

// Sign creates a digital signature for data using the process's private key.
func (kr *KeyRegistry) Sign(processID string, data []byte) string {
	kp, exists := kr.keys[processID]
	if !exists {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(kp.PrivateKey))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks a digital signature against the process's key pair.
func (kr *KeyRegistry) Verify(processID string, data []byte, sig string) bool {
	expected := kr.Sign(processID, data)
	if expected == "" {
		return false
	}
	return hmac.Equal([]byte(expected), []byte(sig))
}

// GetPublicKey returns the public key of a registered process.
func (kr *KeyRegistry) GetPublicKey(processID string) (string, bool) {
	kp, exists := kr.keys[processID]
	if !exists {
		return "", false
	}
	return kp.PublicKey, true
}

func sharedKeyID(p, q string) string {
	return p + ":" + q
}
