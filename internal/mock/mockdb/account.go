package mockdb

import (
	"github.com/ribice/gorsk/internal"
)

// Account database mock
type Account struct {
	CreateFn         func(model.User) (*model.User, error)
	ChangePasswordFn func(*model.User) error
}

// Create mock
func (a *Account) Create(usr model.User) (*model.User, error) {
	return a.CreateFn(usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(usr *model.User) error {
	return a.ChangePasswordFn(usr)
}
