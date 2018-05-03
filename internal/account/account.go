package account

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/auth"
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
func (s *Service) Create(c echo.Context, req model.User) (*model.User, error) {
	if err := s.rbac.AccountCreate(c, req.RoleID, req.CompanyID, req.LocationID); err != nil {
		return nil, err
	}
	req.Password = auth.HashPassword(req.Password)
	return s.adb.Create(req)
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(c echo.Context, oldPass, newPass string, id int) error {
	if err := s.rbac.EnforceUser(c, id); err != nil {
		return err
	}
	u, err := s.udb.View(id)
	if err != nil {
		return err
	}
	if !auth.HashMatchesPassword(u.Password, oldPass) {
		return echo.NewHTTPError(http.StatusBadRequest, "old password is not correct")
	}
	u.Password = auth.HashPassword(newPass)
	return s.adb.ChangePassword(u)
}
