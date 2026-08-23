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

// Digital signatures provide NON-REPUDIATION: a process can prove that
// a specific process signed a message, and the signer cannot deny it.
//
// From Cachin et al., Section 2.3:
//
//	  A digital signature scheme uses:
//	    - A PRIVATE KEY (known only to the signer) to produce signatures
//	    - A PUBLIC KEY (known to everyone) to verify signatures
//
// IMPORTANT SIMPLIFICATION:
//   In a real system, digital signatures use asymmetric cryptography
//   (e.g., RSA, Ed25519, ECDSA). In our teaching simulator, we use
//   HMAC as a stand-in. This captures the API and properties:
//     - Only the key holder can produce valid signatures
//     - Anyone with the registry can verify
//     - Tampering is detectable
//
//   This is NOT cryptographically secure for real applications.
//   The point is to teach the ABSTRACTION, not the implementation.

// KeyPair holds a public-private key pair for a process.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// KeyRegistry manages cryptographic keys for all processes.
// It simulates a Public Key Infrastructure (PKI) where:
//   - Each process has a key pair (public + private)
//   - Shared keys can be established between process pairs (for MACs)
//   - The registry provides sign/verify operations
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
// In our simplified model, keys are derived deterministically from
// the process ID. In a real system, they would be randomly generated.
func (kr *KeyRegistry) RegisterProcess(processID string) KeyPair {
	kp := KeyPair{
		PublicKey:  HashString("pub:" + processID),
		PrivateKey: HashString("priv:" + processID),
	}
	kr.keys[processID] = kp
	return kp
}

// RegisterSharedKey creates a shared secret between two processes.
// Used for MAC-based authentication between specific pairs.
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
//
// In our simplified model, this uses HMAC with the private key.
// In a real system, this would use an asymmetric algorithm (e.g., Ed25519).
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
//
// In our simplified model, this re-computes the HMAC and compares.
// In a real system, this would use the public key for verification.
//
// The key teaching point:
//   - Only the key holder (private key) can produce valid signatures
//   - Anyone (with the registry) can verify
//   - If the message was tampered with, verification FAILS
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
