package mock

import (
	"github.com/labstack/echo"

	"github.com/ribice/gorsk/internal"
)

// Auth mock
type Auth struct {
	UserFn func(echo.Context) *model.AuthUser
}

// User mock
func (a *Auth) User(c echo.Context) *model.AuthUser {
	return a.UserFn(c)
}
