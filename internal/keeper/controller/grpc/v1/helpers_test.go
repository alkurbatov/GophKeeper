package v1_test

import (
	"context"
	"net"
	"testing"

	v1 "github.com/alkurbatov/goph-keeper/internal/keeper/controller/grpc/v1"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
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
		Auth:  &usecase.AuthUseCaseMock{},
		Users: &usecase.UsersUseCaseMock{},
	}
}

func createTestServer(
	t *testing.T,
	useCases usecase.UseCases,
) *grpc.ClientConn {
	t.Helper()
	require := require.New(t)

	srv := grpc.NewServer()
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
