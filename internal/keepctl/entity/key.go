package entity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/alkurbatov/goph-keeper/internal/libraries/creds"
)

// NB (alkurbatov): Never use more than 2^32 random nonces with a given key
// because of the risk of a repeat.
// See https://pkg.go.dev/crypto/cipher#example-NewGCM-Encrypt
const _defaultNonceLength = 12

// Key is user's encyption key.
type Key struct {
	sum [sha256.Size]byte
}

// NewKey creates new encryption key from the provided username and password.
func NewKey(username string, password creds.Password) Key {
	sum := sha256.Sum256([]byte(username + "@" + string(password)))

	return Key{sum}
}

// Hash provides hash of the encryption key.
func (k Key) Hash() string {
	return hex.EncodeToString(k.sum[:])
}

// Encrypt encrypts provided message with secret key.
func (k Key) Encrypt(data []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(k.sum[:])
	if err != nil {
		return nil, fmt.Errorf("Key - Encrypt - aes.NewCipher: %w", err)
	}

	nonce := make([]byte, _defaultNonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("Key - Encrypt - io.ReadFull: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("Key - Encrypt - aes.NewGCM: %w", err)
	}

	// NB (alkurbatov): We don't plan to store nonce somewhere,
	// thus lets put it at the beginning of the data.
	return aesgcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts provided bytes.
func (k Key) Decrypt(data []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(k.sum[:])
	if err != nil {
		return nil, fmt.Errorf("Key - Decrypt - aes.NewCipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("Key - Decrypt - aes.NewGCM: %w", err)
	}

	nonce, payload := data[:_defaultNonceLength], data[_defaultNonceLength:]

	decrypted, err := aesgcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("Key - Decrypt - aesgcm.Open: %w", err)
	}

	return decrypted, nil
}
