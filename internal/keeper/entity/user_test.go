package entity_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestUserWithFromContext(t *testing.T) {
	expected := entity.User{
		ID:       uuid.NewV4(),
		Username: gophtest.Username,
	}

	ctx := expected.WithContext(context.Background())

	require.Equal(t, expected, *entity.UserFromContext(ctx))
}

func TestUserFromCleanContext(t *testing.T) {
	require.Nil(t, entity.UserFromContext(context.Background()))
}
