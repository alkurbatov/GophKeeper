package editcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/spf13/cobra"
)

var (
	data string

	binCmd = &cobra.Command{
		Use:     "bin [secret id] [flags]",
		Short:   "Edit stored binary secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditBin,
	}
)

func init() {
	binCmd.Flags().StringVarP(
		&data,
		"binary-data",
		"b",
		"",
		"New binary data",
	)
}

func doEditBin(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription && data == "" {
		return errFlagsRequired
	}

	if err := clientApp.Usecases.Secrets.EditBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		[]byte(data),
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
