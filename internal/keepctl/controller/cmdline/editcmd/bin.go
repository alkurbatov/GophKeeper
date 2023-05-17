package editcmd

import (
	"errors"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/spf13/cobra"
)

var (
	data []byte

	binCmd = &cobra.Command{
		Use:     "bin [secret id] [flags]",
		Short:   "Edit stored binary secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditBin,
	}
)

func init() {
	binCmd.Flags().BytesHexVarP(
		&data,
		"binary-data",
		"b",
		nil,
		"New binary data in hex format",
	)
}

func doEditBin(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription && len(data) == 0 {
		return errFlagsRequired
	}

	if err := clientApp.Usecases.Secrets.EditBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		data,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, usecase.ErrKindMismatch) {
			return usecase.ErrKindMismatch
		}

		return entity.Unwrap(err)
	}

	return nil
}
