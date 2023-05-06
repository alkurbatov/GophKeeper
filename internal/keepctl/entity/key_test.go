package entity_test

import (
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/libraries/creds"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestKeyToHash(t *testing.T) {
	sat := entity.NewKey(gophtest.Username, gophtest.Password)

	snaps.MatchSnapshot(t, sat.Hash())
}

func TestEncryptDecrypt(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password creds.Password
	}{
		{
			name:     "Basic key",
			username: gophtest.Username,
			password: gophtest.Password,
		},
		{
			name:     "Short key",
			username: "a",
			password: "b",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			msg := []byte("TestEncryptDecrypt")
			sat := entity.NewKey(gophtest.Username, gophtest.Password)

			encrypted, err := sat.Encrypt(msg)
			require.NoError(t, err)

			decrypted, err := sat.Decrypt(encrypted)
			require.NoError(t, err)
			require.Equal(t, msg, decrypted)
		})
	}
}
