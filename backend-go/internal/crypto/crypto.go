// Package crypto provides §1.4 CryptoProvider interface and AES-256-GCM implementation (doc v4.0).
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// CryptoProvider defines core crypto operations (doc §1.4.1).
type CryptoProvider interface {
	Encrypt(data []byte, keyID string) ([]byte, error)
	Decrypt(data []byte, keyID string) ([]byte, error)
	Hash(data []byte) ([]byte, error)
	RandomBytes(size int) ([]byte, error)
}

// AESGCMProvider implements CryptoProvider with AES-256-GCM. KeyID maps to 32-byte key (e.g. from env or KMS).
type AESGCMProvider struct {
	Keys map[string][]byte // keyID -> 32-byte key; in production use KMS
}

// NewAESGCMProvider returns a provider. Keys must be 32 bytes each for AES-256.
func NewAESGCMProvider(keys map[string][]byte) *AESGCMProvider {
	return &AESGCMProvider{Keys: keys}
}

func (p *AESGCMProvider) getKey(keyID string) ([]byte, error) {
	k, ok := p.Keys[keyID]
	if !ok || len(k) != 32 {
		return nil, errors.New("invalid key id or key size")
	}
	return k, nil
}

// Encrypt encrypts data with AES-256-GCM; output is nonce (12) + ciphertext + tag (16).
func (p *AESGCMProvider) Encrypt(data []byte, keyID string) ([]byte, error) {
	key, err := p.getKey(keyID)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts data produced by Encrypt.
func (p *AESGCMProvider) Decrypt(data []byte, keyID string) ([]byte, error) {
	key, err := p.getKey(keyID)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Hash returns SHA-256 of data.
func (p *AESGCMProvider) Hash(data []byte) ([]byte, error) {
	h := sha256.Sum256(data)
	return h[:], nil
}

// RandomBytes returns cryptographically random bytes.
func (p *AESGCMProvider) RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	return b, err
}
