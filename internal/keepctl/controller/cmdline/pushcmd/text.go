package pushcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var (
	text string

	textCmd = &cobra.Command{
		Use:     "text [flags]",
		Short:   "Push arbitrary text",
		PreRunE: preRun,
		RunE:    doPushText,
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

func doPushText(cmd *cobra.Command, _args []string) error {
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
