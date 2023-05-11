package editcmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	errFlagsRequired = errors.New("at least one flag required")

	secretName    string
	description   string
	noDescription bool
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit secret data stored in the keeper service",
}

func init() {
	EditCmd.PersistentFlags().StringVarP(
		&secretName,
		"name",
		"n",
		"",
		"Name of the stored secret",
	)
	EditCmd.PersistentFlags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"Additional description of stored data (activation codes, bank names etc)",
	)
	EditCmd.PersistentFlags().BoolVar(
		&noDescription,
		"no-description",
		false,
		"Remove description from the secret",
	)

	EditCmd.MarkFlagsMutuallyExclusive("description", "no-description")

	EditCmd.AddCommand(textCmd)
}
