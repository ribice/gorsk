// Package user contains user application services
package user

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/platform/query"
	"github.com/ribice/gorsk/internal/platform/structs"
)

// New creates new user application service
func New(udb model.UserDB, rbac model.RBACService, auth model.AuthService) *Service {
	return &Service{udb: udb, rbac: rbac, auth: auth}
}

// Service represents user application service
type Service struct {
	udb  model.UserDB
	rbac model.RBACService
	auth model.AuthService
}

// List returns list of users
func (s *Service) List(c echo.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.udb.List(q, p)
}

// View returns single user
func (s *Service) View(c echo.Context, id int) (*model.User, error) {
	if err := s.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return s.udb.View(id)
}

// Delete deletes a user
func (s *Service) Delete(c echo.Context, id int) error {
	u, err := s.udb.View(id)
	if err != nil {
		return err
	}
	if err := s.rbac.IsLowerRole(c, u.Role.AccessLevel); err != nil {
		return err
	}
	return s.udb.Delete(u)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName *string
	LastName  *string
	Mobile    *string
	Phone     *string
	Address   *string
}

// Update updates user's contact information
func (s *Service) Update(c echo.Context, u *Update) (*model.User, error) {
	if err := s.rbac.EnforceUser(c, u.ID); err != nil {
		return nil, err
	}
	usr, err := s.udb.View(u.ID)
	if err != nil {
		return nil, err
	}
	structs.Merge(usr, u)
	return s.udb.Update(usr)
}
