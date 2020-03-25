package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/ribice/gorsk"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, gorsk.User) (gorsk.User, error)
	ViewFn           func(orm.DB, int) (gorsk.User, error)
	FindByUsernameFn func(orm.DB, string) (gorsk.User, error)
	FindByTokenFn    func(orm.DB, string) (gorsk.User, error)
	ListFn           func(orm.DB, *gorsk.ListQuery, gorsk.Pagination) ([]gorsk.User, error)
	DeleteFn         func(orm.DB, gorsk.User) error
	UpdateFn         func(orm.DB, gorsk.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr gorsk.User) (gorsk.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (gorsk.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (gorsk.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (gorsk.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *gorsk.ListQuery, p gorsk.Pagination) ([]gorsk.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr gorsk.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr gorsk.User) error {
	return u.UpdateFn(db, usr)
}
