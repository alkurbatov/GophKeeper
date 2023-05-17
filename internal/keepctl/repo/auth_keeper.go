package repo

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

var _ Auth = (*AuthRepo)(nil)

// AuthRepo is facade to operations regarding authentication in Keeper.
type AuthRepo struct {
	client goph.AuthClient
}

// NewAuthRepo creates and initializes AuthRepo object.
func NewAuthRepo(client goph.AuthClient) *AuthRepo {
	return &AuthRepo{client}
}

// Login authenticates user in the Keeper service.
func (r *AuthRepo) Login(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	req := &goph.LoginRequest{
		Username:    username,
		SecurityKey: securityKey,
	}

	resp, err := r.client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AuthRepo - Login - r.client.Login: %w", entity.NewRequestError(err))
	}

	return resp.GetAccessToken(), nil
}
