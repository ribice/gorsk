package request

import (
	"github.com/labstack/echo"
)

// Credentials contains login request
type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login validates login request
func Login(c echo.Context) (*Credentials, error) {
	cred := new(Credentials)
	if err := c.Bind(cred); err != nil {
		return nil, err

	}
	return cred, nil
}
