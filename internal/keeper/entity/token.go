package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

const _defaultTokenLifeTime = 15 * time.Minute

var _defaultSigningMethod = jwt.SigningMethodHS256

// AccessToken is JWT token used for authentication.
type AccessToken string

// Claims contain token's payload with various info about token itself and authentciated user.
type Claims struct {
	jwt.RegisteredClaims
	Username string
}

// NewAccessToken issues new access token valid for limited period of time.
func NewAccessToken(user User, secret Secret) (AccessToken, error) {
	now := time.Now()

	claims := jwt.MapClaims{}

	claims["iss"] = "Goph"
	claims["jti"] = uuid.NewV4()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(_defaultTokenLifeTime).Unix()

	// User info
	claims["sub"] = user.ID
	claims["username"] = user.Username

	rawToken := jwt.NewWithClaims(_defaultSigningMethod, claims)

	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("AccessToken - NewAccessToken - rawToken.SignedString: %w", err)
	}

	return AccessToken(signedToken), nil
}

// String converts AccessToken to string.
func (t AccessToken) String() string {
	return string(t)
}

// Decode decodes token, verifies it's signature and return claims if the token is valid.
func (t AccessToken) Decode(secret Secret) (*Claims, error) {
	claims := new(Claims)

	if _, err := jwt.ParseWithClaims(
		t.String(),
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	); err != nil {
		return nil, err
	}

	return claims, nil
}
