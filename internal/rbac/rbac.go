package rbac

import (
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"
)

// New creates new RBAC service
func New(udb model.UserDB) *Service {
	return &Service{udb}
}

// Service is RBAC application service
type Service struct {
	udb model.UserDB
}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r model.AccessRole) error {
	return checkBool(!(c.Get("role").(int8) > int8(r)))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(int) == ID)
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *Service) EnforceCompany(c echo.Context, ID int) error {
	if s.isAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, model.CompanyAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("company_id").(int) == ID)
}

// EnforceLocation checks whether the request to change location data
// is done by the user belonging to the requested location
func (s *Service) EnforceLocation(c echo.Context, ID int) error {
	if s.isCompanyAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, model.LocationAdminRole); err != nil {
		return err
	}
	return checkBool((c.Get("location_id").(int) == ID))
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(int8) > int8(model.AdminRole))
}

func (s *Service) isCompanyAdmin(c echo.Context) bool {
	// Must query company ID in database for the given user
	return !(c.Get("role").(int8) > int8(model.CompanyAdminRole))
}

// AccountCreate performs auth check when creating a new account
// Location admin cannot create accounts, needs to be fixed on EnforceLocation function
func (s *Service) AccountCreate(c echo.Context, roleID, companyID, locationID int) error {
	if err := s.EnforceLocation(c, locationID); err != nil {
		return err
	}
	return s.IsLowerRole(c, model.AccessRole(roleID))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r model.AccessRole) error {
	return checkBool(c.Get("role").(int8) < int8(r))
}
