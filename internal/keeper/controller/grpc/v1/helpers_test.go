package v1_test

import (
	"context"
	"net"
	"testing"

	v1 "github.com/alkurbatov/goph-keeper/internal/keeper/controller/grpc/v1"
	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func requireEqualCode(t *testing.T, expected codes.Code, err error) {
	t.Helper()

	rv, ok := status.FromError(err)

	require.True(t, ok)
	require.Equal(t, expected, rv.Code())
}

func newUseCasesMock() usecase.UseCases {
	return usecase.UseCases{
		Auth:    &usecase.AuthUseCaseMock{},
		Secrets: &usecase.SecretsUseCaseMock{},
		Users:   &usecase.UsersUseCaseMock{},
	}
}

func fakeAuthInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	user := entity.User{
		ID:       uuid.NewV4(),
		Username: gophtest.Username,
	}

	return handler(user.WithContext(ctx), req)
}

func createTestServer(
	t *testing.T,
	useCases usecase.UseCases,
	opts ...grpc.ServerOption,
) *grpc.ClientConn {
	t.Helper()
	require := require.New(t)

	srvOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(v1.DefaultMaxMessageSize),
	}
	srvOpts = append(srvOpts, opts...)

	srv := grpc.NewServer(srvOpts...)
	v1.RegisterRoutes(srv, &useCases)

	lis := bufconn.Listen(1024 * 1024)
	go func() {
		require.NoError(srv.Serve(lis))
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.Dial(
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(err)

	t.Cleanup(func() {
		require.NoError(conn.Close())
		srv.Stop()
		require.NoError(lis.Close())
	})

	return conn
}

func createTestServerWithFakeAuth(
	t *testing.T,
	useCases usecase.UseCases,
) *grpc.ClientConn {
	t.Helper()

	return createTestServer(t, useCases, grpc.ChainUnaryInterceptor(fakeAuthInterceptor))
}
