package mockdb

import (
	"github.com/ribice/gorsk/internal"
)

// User database mock
type User struct {
	ViewFn           func(int) (*model.User, error)
	FindByUsernameFn func(string) (*model.User, error)
	FindByTokenFn    func(string) (*model.User, error)
	ListFn           func(*model.ListQuery, *model.Pagination) ([]model.User, error)
	DeleteFn         func(*model.User) error
	UpdateFn         func(*model.User) (*model.User, error)
}

// View mock
func (u *User) View(id int) (*model.User, error) {
	return u.ViewFn(id)
}

// FindByUsername mock
func (u *User) FindByUsername(username string) (*model.User, error) {
	return u.FindByUsernameFn(username)
}

// FindByToken mock
func (u *User) FindByToken(token string) (*model.User, error) {
	return u.FindByTokenFn(token)
}

// List mock
func (u *User) List(lq *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	return u.ListFn(lq, p)
}

// Delete mock
func (u *User) Delete(usr *model.User) error {
	return u.DeleteFn(usr)
}

// Update mock
func (u *User) Update(usr *model.User) (*model.User, error) {
	return u.UpdateFn(usr)
}
