package config_test

import (
	"os"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/config"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestDefaultConfig(t *testing.T) {
	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}

func TestConfigFromEnv(t *testing.T) {
	os.Setenv("USERNAME", gophtest.Username)
	os.Setenv("PASSWORD", string(gophtest.Password))
	os.Setenv("ADDRESS", "192.168.0.10:8080")
	os.Setenv("CA_PATH", "/etc/ssl/root.crt")
	os.Setenv("VERBOSE", "1")

	t.Cleanup(func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("CA_PATH")
		os.Unsetenv("VERBOSE")
	})

	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}
