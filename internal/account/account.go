package account

import (
	"net/http"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/auth"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/gin-gonic/gin"
)

// New creates new user application service
func New(adb model.AccountDB, udb model.UserDB, rbac model.RBACService) *Service {
	return &Service{
		adb:  adb,
		udb:  udb,
		rbac: rbac,
	}
}

// Service represents account application service
type Service struct {
	adb  model.AccountDB
	udb  model.UserDB
	rbac model.RBACService
}

// Create creates a new user account
func (s *Service) Create(c *gin.Context, req *model.User) error {
	if !s.rbac.AccountCreate(c, req.RoleID, req.CompanyID, req.LocationID) {
		return apperr.Forbidden
	}
	req.Password = auth.HashPassword(req.Password)
	return s.adb.Create(c, req)
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(c *gin.Context, oldPass, newPass string, id int) error {
	if !s.rbac.EnforceUser(c, id) {
		return apperr.Forbidden
	}
	u, err := s.udb.View(c, id)
	if err != nil {
		return err
	}
	if !auth.HashMatchesPassword(u.Password, oldPass) {
		return apperr.New(http.StatusBadRequest, "old password is not correct")
	}
	u.Password = auth.HashPassword(newPass)
	return s.adb.ChangePassword(c, u)
}
