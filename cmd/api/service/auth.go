package service

import (
	"net/http"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/ribice/gorsk/cmd/api/request"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal/auth"
)

// Auth represents auth http service
type Auth struct {
	svc *auth.Service
}

// NewAuth creates new auth http service
func NewAuth(svc *auth.Service, r *gin.Engine) {
	a := Auth{svc}
	// swagger:route POST /login auth login
	// Logs in user by username and password.
	// responses:
	//  200: loginResp
	//  400: errMsg
	//  401: errMsg
	//  404: errMsg
	//  500: err
	r.POST("/login", a.login)
}

func (a *Auth) login(c *gin.Context) {
	cred, err := request.Login(c)
	if err != nil {
		return
	}
	r, err := a.svc.Authenticate(c, cred.Username, cred.Password)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, r)
}
