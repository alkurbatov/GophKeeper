package editcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var (
	text string

	textCmd = &cobra.Command{
		Use:   "text [secret id]",
		Short: "Edit stored text secret",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doEditText,
	}
)

func init() {
	EditCmd.PersistentFlags().StringVarP(
		&text,
		"text",
		"t",
		"",
		"Stored secret text",
	)
}

func doEditText(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	if secretName == "" && description == "" && !noDescription && text == "" {
		return errFlagsRequired
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	if err := clientApp.Usecases.Secrets.EditText(
		cmd.Context(),
		clientApp.AccessToken,
		id,
		secretName,
		description,
		noDescription,
		text,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
