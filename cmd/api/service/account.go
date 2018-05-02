package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/account"

	"github.com/ribice/gorsk/cmd/api/request"
)

// Account represents account http
type Account struct {
	svc *account.Service
}

// NewAccount creates new account http service
func NewAccount(svc *account.Service, e *echo.Group) {
	a := Account{svc: svc}
	ar := e.Group("/users")
	// swagger:route POST /v1/users users accCreate
	// Creates new user account.
	// responses:
	//  200: userResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	ar.POST("", a.create)
	// swagger:operation PATCH /v1/users/{id}/password users pwChange
	// ---
	// summary: Changes user's password.
	// description: If user's old passowrd is correct, it will be replaced with new password.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: int
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ar.PATCH("/:id/password", a.changePassword)
}

func (a *Account) create(c echo.Context) error {
	r, err := request.AccountCreate(c)
	if err != nil {
		return err
	}
	usr, err := a.svc.Create(c, model.User{
		Username:   r.Username,
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		CompanyID:  r.CompanyID,
		LocationID: r.LocationID,
		RoleID:     r.RoleID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, usr)
}

func (a *Account) changePassword(c echo.Context) error {
	p, err := request.PasswordChange(c)
	if err != nil {
		return err
	}
	if err := a.svc.ChangePassword(c, p.OldPassword, p.NewPassword, p.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
