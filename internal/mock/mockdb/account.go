package mockdb

import (
	"context"

	"github.com/ribice/gorsk/internal"
)

// Account database mock
type Account struct {
	CreateFn         func(context.Context, *model.User) error
	ChangePasswordFn func(context.Context, *model.User) error
}

// Create mock
func (a *Account) Create(c context.Context, usr *model.User) error {
	return a.CreateFn(c, usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(c context.Context, usr *model.User) error {
	return a.ChangePasswordFn(c, usr)
}
