package cmdline

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets for current user",
	RunE:  doList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func doList(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	data, err := clientApp.Usecases.Secrets.List(cmd.Context(), clientApp.AccessToken)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	t := tabby.New()
	t.AddHeader("ID", "Name", "Kind", "Description")

	for _, secret := range data {
		t.AddLine(secret.GetId(), secret.GetName(), secret.Kind.String(), string(secret.GetMetadata()))
	}

	t.Print()

	return nil
}
