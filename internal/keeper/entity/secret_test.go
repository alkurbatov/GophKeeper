package entity_test

import (
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestSecretToString(t *testing.T) {
	tt := []struct {
		name string
		data string
	}{
		{
			name: "Print password",
			data: "1q2w3e",
		},
		{
			name: "Print empty secret",
			data: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := entity.Secret(tc.data)
			snaps.MatchSnapshot(t, sat.String())
		})
	}
}

func TestSecretURIToString(t *testing.T) {
	tt := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name: "Print database URI",
			data: "postgres://postgres:postgres@127.0.0.1:5432/goph?sslmode=disable",
		},
		{
			name: "Print empty secret",
			data: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := entity.SecretURI(tc.data)

			snaps.MatchSnapshot(t, sat.String())
		})
	}
}
