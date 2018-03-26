package rbac

import (
	"github.com/gin-gonic/gin"
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

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c *gin.Context, r model.AccessRole) bool {
	return !(c.MustGet("role").(int8) > int8(r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c *gin.Context, ID int) bool {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	return (c.GetInt("id") == ID) || s.isAdmin(c)
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *Service) EnforceCompany(c *gin.Context, ID int) bool {
	return (c.GetInt("company_id") == ID && s.EnforceRole(c, model.CompanyAdminRole)) || s.isAdmin(c)
}

// EnforceLocation checks whether the request to change location data
// is done by the user belonging to the requested location
func (s *Service) EnforceLocation(c *gin.Context, ID int) bool {
	return ((c.GetInt("location_id") == ID) && s.EnforceRole(c, model.LocationAdminRole)) || s.isCompanyAdmin(c)
}

func (s *Service) isAdmin(c *gin.Context) bool {
	return !(c.MustGet("role").(int8) > int8(model.AdminRole))
}

func (s *Service) isCompanyAdmin(c *gin.Context) bool {
	// Must query company ID in database for the given user
	return !(c.MustGet("role").(int8) > int8(model.CompanyAdminRole))
}

// AccountCreate performs auth check when creating a new account
// Location admin cannot create accounts, needs to be fixed on EnforceLocation function
func (s *Service) AccountCreate(c *gin.Context, roleID, companyID, locationID int) bool {
	companyCheck := s.EnforceCompany(c, companyID)
	locationCheck := s.EnforceLocation(c, locationID)
	roleCheck := s.EnforceRole(c, model.AccessRole(roleID))
	return companyCheck && locationCheck && roleCheck && s.IsLowerRole(c, model.AccessRole(roleID))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c *gin.Context, r model.AccessRole) bool {
	return !(c.MustGet("role").(int8) >= int8(r))
}
