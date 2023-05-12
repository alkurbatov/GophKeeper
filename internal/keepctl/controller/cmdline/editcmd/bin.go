package editcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var (
	data string

	binCmd = &cobra.Command{
		Use:   "bin [secret id] [flags]",
		Short: "Edit stored binary secret",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doEditBin,
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

func doEditBin(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	if secretName == "" && description == "" && !noDescription && data == "" {
		return errFlagsRequired
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	if err := clientApp.Usecases.Secrets.EditBinary(
		cmd.Context(),
		clientApp.AccessToken,
		id,
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
