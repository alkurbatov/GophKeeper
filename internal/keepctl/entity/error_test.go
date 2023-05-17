package entity_test

import (
	"fmt"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRequestErrorFromBasicError(t *testing.T) {
	sat := entity.NewRequestError(grpc.ErrServerStopped)

	snaps.MatchSnapshot(t, sat.Error())
}

func TestRequestErrorFromGenericGRPCError(t *testing.T) {
	gRPCErr := status.Errorf(codes.Unauthenticated, "bad login or password")
	sat := entity.NewRequestError(gRPCErr)

	snaps.MatchSnapshot(t, sat.Error())
}

func TestRequestErrorFromGRPCBadRequest(t *testing.T) {
	details := &errdetails.BadRequest{}
	details.FieldViolations = append(
		details.FieldViolations,
		&errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: "not set",
		},
		&errdetails.BadRequest_FieldViolation{
			Field:       "security_key",
			Description: "not set",
		},
	)

	st := status.New(codes.InvalidArgument, "invalid request")
	st, err := st.WithDetails(details)
	require.NoError(t, err)

	sat := entity.NewRequestError(st.Err())

	snaps.MatchSnapshot(t, sat.Error())
}

func TestUnwrap(t *testing.T) {
	tt := []struct {
		name string
		err  error
	}{
		{
			name: "Test unwrap of RequestError",
			err:  entity.NewRequestError(grpc.ErrServerStopped),
		},
		{
			name: "Test unwrap of other errors",
			err:  grpc.ErrServerStopped,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := fmt.Errorf("ErrorTest - TestUnwrap - SomeError: %w", tc.err)

			snaps.MatchSnapshot(t, entity.Unwrap(err).Error())
		})
	}
}
