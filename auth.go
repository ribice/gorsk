package gorsk

import (
	"github.com/labstack/echo"
)

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken holds authentication token details
type RefreshToken struct {
	Token string `json:"token"`
}

// RBACService represents role-based access control service interface
type RBACService interface {
	User(echo.Context) AuthUser
	EnforceRole(echo.Context, AccessRole) error
	EnforceUser(echo.Context, int) error
	EnforceCompany(echo.Context, int) error
	EnforceLocation(echo.Context, int) error
	AccountCreate(echo.Context, AccessRole, int, int) error
	IsLowerRole(echo.Context, AccessRole) error
}
