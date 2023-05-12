package entity

import (
	"regexp"
	"strings"
)

// Secret is sensitive value (e.g. password) which shouldn't leak to logs.
type Secret string

// String converts Secret to string.
func (s Secret) String() string {
	return strings.Repeat("*", len(s))
}

// SecretURI is URI with sensitive values (e.g. login:password) which shouldn't leak to logs.
type SecretURI string

var _URISecrets = regexp.MustCompile(`(://).*:.*(@)`)

// String converts SecretURI to string.
func (u SecretURI) String() string {
	return string(_URISecrets.ReplaceAll([]byte(u), []byte("$1*****:*****$2")))
}
