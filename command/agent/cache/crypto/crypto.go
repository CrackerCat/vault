package crypto

import (
	"context"
	"fmt"
)

const (
	KeyID = "root"
)

// KeyManager TODO
type KeyManager interface {
	Get() []byte
	Renewable() bool
	Renewer(context.Context, chan struct{}) error
	Encrypt(context.Context, []byte, []byte) ([]byte, error)
	Decrypt(context.Context, []byte, []byte) ([]byte, error)
}

// Encrypt TODO
func Encrypt(ctx context.Context, k KeyManager, plaintext []byte, aad []byte) ([]byte, error) {
	return k.Encrypt(ctx, plaintext, aad)
}

// Decrypt TODO
func Decrypt(ctx context.Context, k KeyManager, ciphertext, aad []byte) ([]byte, error) {
	return k.Decrypt(ctx, ciphertext, aad)
}

// New TODO
func New(keyType string, existingKey []byte) (KeyManager, error) {
	switch keyType {
	case "k8s":
		k8s, err := NewK8s(existingKey)
		return k8s, err
	case "response":
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("invalid key type: %s", keyType)
	}
}

// Get TODO
func Get(k KeyManager) []byte {
	return k.Get()
}

// Renewable TODO
func Renewable(k KeyManager) bool {
	return k.Renewable()
}

// Renewer TODO
func Renewer(ctx context.Context, k KeyManager, notify chan struct{}) error {
	return k.Renewer(ctx, notify)
}
