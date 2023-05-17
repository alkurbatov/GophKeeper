package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/config"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/gkampitakis/go-snaps/snaps"
)

func unsetGophEnv() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GOPH_") {
			pair := strings.SplitN(env, "=", 2)

			os.Unsetenv(pair[0])
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	unsetGophEnv()

	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}

func TestConfigFromEnv(t *testing.T) {
	os.Setenv("GOPH_USERNAME", gophtest.Username)
	os.Setenv("GOPH_PASSWORD", string(gophtest.Password))
	os.Setenv("GOPH_ADDRESS", "192.168.0.10:8080")
	os.Setenv("GOPH_CA_PATH", "/etc/ssl/root.crt")
	os.Setenv("GOPH_VERBOSE", "1")

	t.Cleanup(unsetGophEnv)

	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}
