package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var credsCmd = &cobra.Command{
	Use:   "creds [login] [password]",
	Short: "Push login and password",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd // count of required args
	RunE:  doPushCreds,
}

func doPushCreds(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	id, err := clientApp.Usecases.Secrets.PushCreds(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		args[0],
		args[1],
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
