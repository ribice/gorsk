// Package user contains user application services
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/errors"
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
func (s *Service) List(c *gin.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.udb.List(c, q, p)
}

// View returns single user
func (s *Service) View(c *gin.Context, id int) (*model.User, error) {
	if !s.rbac.EnforceUser(c, id) {
		return nil, apperr.Forbidden
	}
	return s.udb.View(c, id)
}

// Delete deletes a user
func (s *Service) Delete(c *gin.Context, id int) error {
	u, err := s.udb.View(c, id)
	if err != nil {
		return err
	}
	if !s.rbac.IsLowerRole(c, u.Role.AccessLevel) {
		return apperr.Forbidden
	}
	u.Delete()
	return s.udb.Delete(c, u)
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
func (s *Service) Update(c *gin.Context, u *Update) (*model.User, error) {
	if !s.rbac.EnforceUser(c, u.ID) {
		return nil, apperr.Forbidden
	}
	usr, err := s.udb.View(c, u.ID)
	if err != nil {
		return nil, err
	}
	structs.Merge(usr, u)
	return s.udb.Update(c, usr)
}
