package editcmd

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var (
	binary string

	binCmd = &cobra.Command{
		Use:   "bin [secret id]",
		Short: "Edit stored binary data",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doEditBin,
	}
)

func init() {
	binCmd.PersistentFlags().StringVarP(
		&binary,
		"bin",
		"b",
		"",
		"Stored binary data",
	)
}

func doEditBin(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	if secretName == "" && description == "" && !noDescription && binary == "" {
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
		[]byte(binary),
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
