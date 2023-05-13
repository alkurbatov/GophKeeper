package editcmd

import (
	"errors"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/spf13/cobra"
)

var (
	number     string
	expiration string
	holder     string
	cvv        int32

	cardCmd = &cobra.Command{
		Use:     "card [secret id] [flags]",
		Short:   "Edit stored bank card secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditCard,
	}
)

func init() {
	cardCmd.Flags().StringVar(
		&number,
		"number",
		"",
		"Card number",
	)
	cardCmd.Flags().StringVar(
		&expiration,
		"expiration",
		"",
		"Card expiration date",
	)
	cardCmd.Flags().StringVar(
		&holder,
		"holder",
		"",
		"Card holder name and surname",
	)
	cardCmd.Flags().Int32Var(
		&cvv,
		"cvv",
		0,
		"Card verification value",
	)
}

func doEditCard(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription &&
		number == "" && expiration == "" && holder == "" && cvv == 0 {
		return errFlagsRequired
	}

	if err := clientApp.Usecases.Secrets.EditCard(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		number,
		expiration,
		holder,
		cvv,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, usecase.ErrKindMismatch) {
			return usecase.ErrKindMismatch
		}

		return entity.Unwrap(err)
	}

	return nil
}
