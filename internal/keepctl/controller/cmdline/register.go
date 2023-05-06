package cmdline

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	RunE:  doRegister,
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func doRegister(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	accessToken, err := clientApp.Usecases.Users.Register(
		cmd.Context(),
		cfg.Username,
		clientApp.Key,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("access-token", accessToken).Msg("New user successfully created")

	return nil
}
