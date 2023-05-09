package cmdline

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret id]",
	Short: "Delete secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func doDelete(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	if err := clientApp.Usecases.Secrets.Delete(cmd.Context(), clientApp.AccessToken, id); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	return nil
}
