package service

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/user"

	"github.com/ribice/gorsk/cmd/api/request"
)

// User represents user http service
type User struct {
	svc *user.Service
}

// NewUser creates new user http service
func NewUser(svc *user.Service, ur *echo.Group) {
	u := User{svc: svc}
	// swagger:operation GET /v1/users users listUsers
	// ---
	// summary: Returns list of users.
	// description: Returns list of users. Depending on the user role requesting it, it may return all users for SuperAdmin/Admin users, all company/location users for Company/Location admins, and an error for non-admin users.
	// parameters:
	// - name: limit
	//   in: query
	//   description: number of results
	//   type: int
	//   required: false
	// - name: page
	//   in: query
	//   description: page number
	//   type: int
	//   required: false
	// responses:
	//   "200":
	//     "$ref": "#/responses/userListResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ur.GET("", u.list)
	// swagger:operation GET /v1/users/{id} users getUser
	// ---
	// summary: Returns a single user.
	// description: Returns a single user by its ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: int
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/userResp"
	//   "400":
	//     "$ref": "#/responses/err"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "404":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ur.GET("/:id", u.view)
	// swagger:operation PATCH /v1/users/{id} users userUpdate
	// ---
	// summary: Updates user's contact information
	// description: Updates user's contact information -> first name, last name, mobile, phone, address.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of user
	//   type: int
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/userUpdate"
	// responses:
	//   "200":
	//     "$ref": "#/responses/userResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ur.PATCH("/:id", u.update)
	// swagger:operation DELETE /v1/users/{id} users userDelete
	// ---
	// summary: Deletes a user
	// description: Deletes a user with requested ID.
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
	//     "$ref": "#/responses/err"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ur.DELETE("/:id", u.delete)
}

type listResponse struct {
	Users []model.User `json:"users"`
	Page  int          `json:"page"`
}

func (u *User) list(c echo.Context) error {
	p, err := request.Paginate(c)
	if err != nil {
		return err
	}
	result, err := u.svc.List(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, listResponse{result, p.Page})
}

func (u *User) view(c echo.Context) error {
	id, err := request.ID(c)
	if err != nil {
		return err
	}
	result, err := u.svc.View(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (u *User) update(c echo.Context) error {
	req, err := request.UserUpdate(c)
	if err != nil {
		return err
	}
	usr, err := u.svc.Update(c, &user.Update{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Phone:     req.Phone,
		Address:   req.Address,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, usr)
}

func (u *User) delete(c echo.Context) error {
	id, err := request.ID(c)
	if err != nil {
		return err
	}
	if err := u.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
