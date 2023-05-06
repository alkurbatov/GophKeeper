package v1_test

import (
	"context"
	"testing"

	v1 "github.com/alkurbatov/goph-keeper/internal/keeper/controller/grpc/v1"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const testMethod = "/goph.keeper.v1.Secrets/Create"

func newTestServerInfo() *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: testMethod}
}

func fakeHandler(_ context.Context, data any) (any, error) {
	return data, nil
}

func TestAuthOfAnonymousMethods(t *testing.T) {
	tt := []struct {
		name   string
		method string
	}{
		{
			name:   "User Register is allowed",
			method: "/goph.keeper.v1.Users/Register",
		},
		{
			name:   "Auth Login is allowed",
			method: "/goph.keeper.v1.Auth/Login",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: tc.method}

			sat := v1.AuthUnaryInterceptor(gophtest.Secret)
			_, err := sat(context.Background(), nil, info, fakeHandler)

			require.NoError(t, err)
		})
	}
}

func TestAuthIfNoMetadata(t *testing.T) {
	info := newTestServerInfo()

	sat := v1.AuthUnaryInterceptor(gophtest.Secret)
	_, err := sat(context.Background(), nil, info, fakeHandler)

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestAuthInterceptor(t *testing.T) {
	tt := []struct {
		name string
		keys map[string]string
		code codes.Code
	}{
		{
			name: "Access blocked if no authorization key",
			keys: map[string]string{},
			code: codes.Unauthenticated,
		},
		{
			name: "Access blocked if token is empty string",
			keys: map[string]string{"authorization": ""},
			code: codes.Unauthenticated,
		},
		{
			name: "Access blocked if token is invalid",
			keys: map[string]string{"authorization": "xxx"},
			code: codes.Unauthenticated,
		},
		{
			name: "Access blocked if token signed with different secret",
			keys: map[string]string{
				"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk2ODM0ODkzODUsImlhdCI6MTY4MzQ4ODQ4NSwiaXNzIjoiR29waCIsImp0aSI6ImUzZDliYzUyLTBiNGYtNDE2Yi04MGM5LTUxYjVjNDYyNGFkZCIsIm5iZiI6MTY4MzQ4ODQ4NSwic3ViIjoiOThiZDgxM2QtMmIxYS00MDkyLThhNzAtMzkzMWEzNGJiNGFkIiwidXNlcm5hbWUiOiJhbCJ9._Tq_925CZeqqEfcw-zO2W_duDdv8_oncl1hn05Lw40Y", //nolint:lll // JWT token
			},
			code: codes.Unauthenticated,
		},
		{
			name: "Access blocked if token is expired",
			keys: map[string]string{
				"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODM0ODkzODUsImlhdCI6MTY4MzQ4ODQ4NSwiaXNzIjoiR29waCIsImp0aSI6ImUzZDliYzUyLTBiNGYtNDE2Yi04MGM5LTUxYjVjNDYyNGFkZCIsIm5iZiI6MTY4MzQ4ODQ4NSwic3ViIjoiOThiZDgxM2QtMmIxYS00MDkyLThhNzAtMzkzMWEzNGJiNGFkIiwidXNlcm5hbWUiOiJhbCJ9.5aqybHcm-OEE4MV91zwDiwBXVFewADscEYY-N5Se_Sw", //nolint:lll // JWT token
			},
			code: codes.Unauthenticated,
		},
		{
			name: "Access granted for valid token",
			keys: map[string]string{
				"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk2ODM0ODkzODUsImlhdCI6MTY4MzQ4ODQ4NSwiaXNzIjoiR29waCIsImp0aSI6ImUzZDliYzUyLTBiNGYtNDE2Yi04MGM5LTUxYjVjNDYyNGFkZCIsIm5iZiI6MTY4MzQ4ODQ4NSwic3ViIjoiOThiZDgxM2QtMmIxYS00MDkyLThhNzAtMzkzMWEzNGJiNGFkIiwidXNlcm5hbWUiOiJhbCJ9.rFwa_AcWwikgUukktcu468aJdPk-xBf6ZDB93YqUxMY", //nolint:lll // JWT token
			},
			code: codes.OK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			info := newTestServerInfo()
			md := metadata.New(tc.keys)
			ctx := metadata.NewIncomingContext(context.Background(), md)

			sat := v1.AuthUnaryInterceptor(gophtest.Secret)
			_, err := sat(ctx, nil, info, fakeHandler)

			requireEqualCode(t, tc.code, err)
		})
	}
}
