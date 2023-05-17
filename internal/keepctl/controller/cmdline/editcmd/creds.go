package editcmd

import (
	"errors"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/spf13/cobra"
)

var (
	login    string
	password string

	credsCmd = &cobra.Command{
		Use:     "creds [secret id] [flags]",
		Short:   "Edit stored credentials secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditCreds,
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

func doEditCreds(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription && login == "" && password == "" {
		return errFlagsRequired
	}

	if err := clientApp.Usecases.Secrets.EditCreds(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		login,
		password,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, usecase.ErrKindMismatch) {
			return usecase.ErrKindMismatch
		}

		return entity.Unwrap(err)
	}

	return nil
}
