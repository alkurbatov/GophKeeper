package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is main configuration of client application.
type Config struct {
	Address string
	CAPath  string
	Verbose bool
}

// New create application config by reading environment variables and
// commandline flags. The flags are read inderectly through binding in cobra.
func New() *Config {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetDefault("virbose", false)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cfg := &Config{
		Address: viper.GetString("address"),
		CAPath:  viper.GetString("ca-path"),
		Verbose: viper.GetBool("verbose"),
	}

	return cfg
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tKeeper address: %s\n", c.Address))
	sb.WriteString(fmt.Sprintf("\t\tCertificate authority path: %s\n", c.CAPath))
	sb.WriteString(fmt.Sprintf("\t\tVerbose: %t", c.Verbose))

	return sb.String()
}
