package crypto

import (
	"context"
	"crypto/rand"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/aead"
)

var _ KeyManager = (*KubeEncryptionKey)(nil)
var _ Encrypter = (*KubeEncryptionKey)(nil)

// KubeEncryptionKey TODO
type KubeEncryptionKey struct {
	renewable bool
	wrapper   *aead.Wrapper
}

// NewK8s returns a new instance of the Kube encryption key. Kubernetes
// encryption keys aren't renewable.
func NewK8s(existingKey []byte) (*KubeEncryptionKey, error) {
	k := &KubeEncryptionKey{
		renewable: false,
		wrapper:   aead.NewWrapper(nil),
	}

	k.wrapper.SetConfig(map[string]string{"key_id": KeyID})
	rootKey := make([]byte, 32)
	if existingKey != nil {
		rootKey = existingKey
	}

	if rootKey == nil {
		_, err := rand.Read(rootKey)
		if err != nil {
			return k, err
		}
	}

	if err := k.wrapper.SetAESGCMKeyBytes(rootKey); err != nil {
		return k, err
	}

	return k, nil
}

// Get returns the encryption key in a format optimized for storage.
// In k8s we store the key as is, so just return the key stored.
func (k *KubeEncryptionKey) Get() []byte {
	return k.wrapper.GetKeyBytes()
}

// Renewable lets the caller know if this encryption key type is
// renewable. In Kubernetes the key isn't renewable.
func (k *KubeEncryptionKey) Renewable() bool {
	return k.renewable
}

// Renewer is used when the encryption key type is renewable. Since Kubernetes
// keys aren't renewable, returning nothing.
func (k *KubeEncryptionKey) Renewer(ctx context.Context, ch chan struct{}) error {
	return nil
}

// Encrypt takes plaintext values and encrypts them using the store key and additional
// data. The ciphertext and nonce are returned and should be used for decryption.
func (k *KubeEncryptionKey) Encrypt(ctx context.Context, plaintext, aad []byte) ([]byte, []byte, error) {
	blob, err := k.wrapper.Encrypt(ctx, plaintext, aad)
	if err != nil {
		return nil, nil, err
	}
	return blob.Ciphertext, blob.IV, nil
}

// Decrypt takes ciphertext and nonce values and returns the decrypted value.
func (k *KubeEncryptionKey) Decrypt(ctx context.Context, ciphertext, nonce, aad []byte) ([]byte, error) {
	blob := &wrapping.EncryptedBlobInfo{
		Ciphertext: ciphertext,
		IV:         nonce,
		KeyInfo: &wrapping.KeyInfo{
			KeyID: KeyID,
		},
	}
	return k.wrapper.Decrypt(ctx, blob, aad)
}
