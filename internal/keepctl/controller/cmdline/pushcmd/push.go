package pushcmd

import (
	"github.com/spf13/cobra"
)

var (
	secretName  string
	description string
)

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push secret to the Keeper service",
}

func init() {
	PushCmd.PersistentFlags().StringVarP(
		&secretName,
		"name",
		"n",
		"",
		"Name of the stored secret",
	)
	PushCmd.PersistentFlags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"Additional description of stored data (activation codes, names of banks etc)",
	)

	PushCmd.MarkPersistentFlagRequired("name")

	PushCmd.AddCommand(binCmd)
	PushCmd.AddCommand(credsCmd)
	PushCmd.AddCommand(textCmd)
}
