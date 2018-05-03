package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/cmd/api/request"

	"github.com/ribice/gorsk/internal/auth"
)

// Auth represents auth http service
type Auth struct {
	svc *auth.Service
}

// NewAuth creates new auth http service
func NewAuth(svc *auth.Service, e *echo.Echo) {
	a := Auth{svc}
	// swagger:route POST /login auth login
	// Logs in user by username and password.
	// responses:
	//  200: loginResp
	//  400: errMsg
	//  401: errMsg
	// 	403: err
	//  404: errMsg
	//  500: err
	e.POST("/login", a.login)
	// swagger:operation GET /refresh/{token} auth refresh
	// ---
	// summary: Refreshes jwt token.
	// description: Refreshes jwt token by checking at database whether refresh token exists.
	// parameters:
	// - name: token
	//   in: path
	//   description: refresh token
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/refreshResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	e.GET("/refresh/:token", a.refresh)
}

func (a *Auth) login(c echo.Context) error {
	cred, err := request.Login(c)
	if err != nil {
		return err
	}
	r, err := a.svc.Authenticate(c, cred.Username, cred.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (a *Auth) refresh(c echo.Context) error {
	r, err := a.svc.Refresh(c, c.Param("token"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}
