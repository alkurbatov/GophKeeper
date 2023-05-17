package usecase_test

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
)

func newTestKey() entity.Key {
	return entity.NewKey(gophtest.Username, gophtest.Password)
}
