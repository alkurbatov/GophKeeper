package entity

import (
	"crypto/sha256"
	"encoding/hex"
)

// Key is user's encyption key.
type Key struct {
	sum [sha256.Size]byte
}

// NewKey creates new encryption key from the provided username and password.
func NewKey(username, password string) Key {
	sum := sha256.Sum256([]byte(username + "@" + password))

	return Key{sum}
}

// Hash provides hash of the encryption key.
func (k Key) Hash() string {
	return hex.EncodeToString(k.sum[:])
}
