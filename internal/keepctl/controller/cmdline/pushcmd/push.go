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
	Short: "Push secret data to the keeper service",
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
		"Additional description of stored data (activation codes, bank names etc)",
	)

	PushCmd.MarkFlagRequired("name")

	PushCmd.AddCommand(textCmd)
}
