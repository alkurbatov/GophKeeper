package cmdline

import (
	"context"

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
	accessToken, err := clientApp.Usecases.Users.Register(
		context.Background(),
		username,
		password,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("access-token", accessToken).Msg("New user created successfully")

	return nil
}
