package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var binCmd = &cobra.Command{
	Use:   "bin [data]",
	Short: "Push arbitrary binary data",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doPushBinary,
}

func doPushBinary(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	id, err := clientApp.Usecases.Secrets.PushBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		[]byte(args[0]),
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
