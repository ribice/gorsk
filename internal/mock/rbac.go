package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal"
)

// RBAC Mock
type RBAC struct {
	EnforceRoleFn     func(*gin.Context, model.AccessRole) bool
	EnforceUserFn     func(*gin.Context, int) bool
	EnforceCompanyFn  func(*gin.Context, int) bool
	EnforceLocationFn func(*gin.Context, int) bool
	AccountCreateFn   func(*gin.Context, int, int, int) bool
	IsLowerRoleFn     func(*gin.Context, model.AccessRole) bool
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c *gin.Context, role model.AccessRole) bool {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c *gin.Context, id int) bool {
	return a.EnforceUserFn(c, id)
}

// EnforceCompany mock
func (a *RBAC) EnforceCompany(c *gin.Context, id int) bool {
	return a.EnforceCompanyFn(c, id)
}

// EnforceLocation mock
func (a *RBAC) EnforceLocation(c *gin.Context, id int) bool {
	return a.EnforceLocationFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c *gin.Context, roleID, companyID, locationID int) bool {
	return a.AccountCreateFn(c, roleID, companyID, locationID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c *gin.Context, role model.AccessRole) bool {
	return a.IsLowerRoleFn(c, role)
}
