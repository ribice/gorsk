package model

import "github.com/gin-gonic/gin"

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken holds authentication token details
type RefreshToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// AuthService represents authentication service interface
type AuthService interface {
	User(*gin.Context) *AuthUser
}

// RBACService represents role-based access control service interface
type RBACService interface {
	EnforceRole(*gin.Context, AccessRole) bool
	EnforceUser(*gin.Context, int) bool
	EnforceCompany(*gin.Context, int) bool
	EnforceLocation(*gin.Context, int) bool
	AccountCreate(*gin.Context, int, int, int) bool
	IsLowerRole(*gin.Context, AccessRole) bool
}
