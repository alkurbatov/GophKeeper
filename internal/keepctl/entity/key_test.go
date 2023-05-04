package entity_test

import (
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestKeyToHash(t *testing.T) {
	sat := entity.NewKey(gophtest.Username, gophtest.Password)

	snaps.MatchSnapshot(t, sat.Hash())
}
