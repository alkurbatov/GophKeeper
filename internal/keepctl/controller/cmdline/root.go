package cmdline

import (
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/app"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg       *config.Config
	clientApp *app.App

	verbose  bool
	address  string
	caPath   string
	username string
	password string

	rootCmd = &cobra.Command{
		Use:               "keepctl",
		Short:             "keepctl is an intercative commandline client for the goph-keeper service",
		PersistentPreRunE: initApp,
	}
)

// Execute executes the root command.
func Execute(buildVersion, buildDate string) error {
	rootCmd.Version = fmt.Sprintf("%s (%s)", buildVersion, buildDate)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initializeApp)
	cobra.OnFinalize(finalizeApp)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(
		&address,
		"address",
		"",
		"Address and port of the keeper service",
	)
	rootCmd.PersistentFlags().StringVar(
		&caPath,
		"ca-path",
		"",
		"Path to certificate authority to verify server certificate",
	)
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Name of a user")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Master password")

	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")

	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("ca-path", rootCmd.PersistentFlags().Lookup("ca-path"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initializeApp does initialization routine before reading commandline flags.
func initializeApp() {
	cfg = config.New()
}

func initApp(cmd *cobra.Command, _ []string) error {
	// NB (alkurbatov): Prerun is executed for EVERY command, even for noop like help.
	if cmd.Name() == "help" {
		return nil
	}

	var err error

	clientApp, err = app.New(cfg)
	if err != nil {
		return err
	}

	return nil
}

// finalizeApp does cleanup at the end of commandline application.
func finalizeApp() {
	if clientApp == nil {
		return
	}

	clientApp.Shutdown()
}
