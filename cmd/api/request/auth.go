package request

import (
	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal/errors"
)

// Credentials contains login request
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login validates login request
func Login(c *gin.Context) (*Credentials, error) {
	cred := new(Credentials)
	if err := c.ShouldBindJSON(cred); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return cred, nil
}
