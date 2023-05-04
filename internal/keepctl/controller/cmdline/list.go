package cmdline

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List secrets for current user",
	PreRunE: login,
	RunE:    doList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func doList(cmd *cobra.Command, args []string) error {
	return nil
}
