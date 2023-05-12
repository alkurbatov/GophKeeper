package cmdline

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/cheynewallace/tabby"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [secret id]",
	Short: "Pull full secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func doPull(cmd *cobra.Command, args []string) error {
	id, err := uuid.FromString(args[0])
	if err != nil {
		return err
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	secret, data, err := clientApp.Usecases.Secrets.Get(cmd.Context(), clientApp.AccessToken, id)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return entity.Unwrap(err)
	}

	t := tabby.New()
	header := []any{"ID", "Name", "Kind", "Description"}
	line := []any{
		secret.GetId(),
		secret.GetName(),
		secret.GetKind().String(),
		string(secret.GetMetadata()),
	}

	switch d := data.(type) {
	case *goph.Binary:
		t.AddHeader(header...)
		t.AddLine(line...)
		t.Print()

		clientApp.Log.Info().Msg(string(d.Binary))

	case *goph.Text:
		t.AddHeader(header...)
		t.AddLine(line...)
		t.Print()

		clientApp.Log.Info().Msg(d.Text)
	}

	return nil
}
