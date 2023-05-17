package cmdline

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

func login(cmd *cobra.Command, _ []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	token, err := clientApp.Usecases.Auth.Login(
		cmd.Context(),
		cfg.Username,
		clientApp.Key,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.AccessToken = token
	clientApp.Log.Debug().
		Str("access-token", token).
		Msg("Login successful")

	return nil
}
