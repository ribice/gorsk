package mockdb

import (
	"github.com/go-pg/pg/orm"
	"github.com/ribice/gorsk/internal"
)

// Account database mock
type Account struct {
	CreateFn         func(orm.DB, model.User) (*model.User, error)
	ChangePasswordFn func(orm.DB, *model.User) error
}

// Create mock
func (a *Account) Create(db orm.DB, usr model.User) (*model.User, error) {
	return a.CreateFn(db, usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(db orm.DB, usr *model.User) error {
	return a.ChangePasswordFn(db, usr)
}
