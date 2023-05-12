package editcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var (
	login    string
	password string

	credsCmd = &cobra.Command{
		Use:   "creds [secret id] [flags]",
		Short: "Edit stored credentials secret",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doEditCreds,
	}
)

func init() {
	credsCmd.Flags().StringVarP(
		&login,
		"login",
		"l",
		"",
		"New login or username",
	)
	credsCmd.Flags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"New password",
	)
}

func doEditCreds(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	if secretName == "" && description == "" && !noDescription && login == "" && password == "" {
		return errFlagsRequired
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	if err := clientApp.Usecases.Secrets.EditCreds(
		cmd.Context(),
		clientApp.AccessToken,
		id,
		secretName,
		description,
		noDescription,
		login,
		password,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
