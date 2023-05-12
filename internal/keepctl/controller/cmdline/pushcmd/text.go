package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var (
	text string

	textCmd = &cobra.Command{
		Use:   "text [flags]",
		Short: "Push arbitrary text",
		RunE:  doPushText,
	}
)

func init() {
	textCmd.Flags().StringVarP(
		&text,
		"text",
		"t",
		"",
		"Text to save",
	)

	textCmd.MarkFlagRequired("text")
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
		description,
		text,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
