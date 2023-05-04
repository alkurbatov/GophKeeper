package cmdline

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

func login(*cobra.Command, []string) error {
	token, err := clientApp.Usecases.Auth.Login(
		context.Background(),
		username,
		password,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().
		Str("access-token", token).
		Msg("Login successful")

	return nil
}
