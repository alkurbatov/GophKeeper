package entity_test

import (
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestAccessToken(t *testing.T) {
	user := entity.User{
		ID:       uuid.NewV4(),
		Username: "root",
	}
	secret := entity.Secret("xxx")

	token, err := entity.NewAccessToken(user, secret)
	require.NoError(t, err)

	claims, err := token.Decode(secret)
	require.NoError(t, err)

	require.Equal(t, user.ID.String(), claims.Subject)
	require.Equal(t, user.Username, claims.Username)
}

func TestAccessTokenDecodeWithWrongSecret(t *testing.T) {
	user := entity.User{
		ID:       uuid.NewV4(),
		Username: "root",
	}

	token, err := entity.NewAccessToken(user, "xxx")
	require.NoError(t, err)

	_, err = token.Decode("yyy")
	require.Error(t, err)
}
