package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var textCmd = &cobra.Command{
	Use:   "text [arbitrary secret text]",
	Short: "Push arbitrary text",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doPushText,
}

func doPushText(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	id, err := clientApp.Usecases.Secrets.PushText(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		args[0],
		description,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
