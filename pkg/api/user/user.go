// Package user contains user application services
package user

import (
	"github.com/labstack/echo"

	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/utl/query"
)

// Create creates a new user account
func (u User) Create(c echo.Context, req gorsk.User) (gorsk.User, error) {
	if err := u.rbac.AccountCreate(c, req.RoleID, req.CompanyID, req.LocationID); err != nil {
		return gorsk.User{}, err
	}
	req.Password = u.sec.Hash(req.Password)
	return u.udb.Create(u.db, req)
}

// List returns list of users
func (u User) List(c echo.Context, p gorsk.Pagination) ([]gorsk.User, error) {
	au := u.rbac.User(c)
	q, err := query.List(au)
	if err != nil {
		return nil, err
	}
	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u User) View(c echo.Context, id int) (gorsk.User, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return gorsk.User{}, err
	}
	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u User) Delete(c echo.Context, id int) error {
	user, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}
	if err := u.rbac.IsLowerRole(c, user.Role.AccessLevel); err != nil {
		return err
	}
	return u.udb.Delete(u.db, user)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName string
	LastName  string
	Mobile    string
	Phone     string
	Address   string
}

// Update updates user's contact information
func (u User) Update(c echo.Context, r Update) (gorsk.User, error) {
	if err := u.rbac.EnforceUser(c, r.ID); err != nil {
		return gorsk.User{}, err
	}

	if err := u.udb.Update(u.db, gorsk.User{
		Base:      gorsk.Base{ID: r.ID},
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Mobile:    r.Mobile,
		Address:   r.Address,
	}); err != nil {
		return gorsk.User{}, err
	}

	return u.udb.View(u.db, r.ID)
}
