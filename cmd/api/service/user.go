package service

import (
	"net/http"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/errors"
	"github.com/ribice/gorsk/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
)

// User represents user http service
type User struct {
	svc *user.Service
}

// NewUser creates new user http service
func NewUser(svc *user.Service, r *gin.RouterGroup) {
	u := User{svc: svc}
	ur := r.Group("/users")
	// swagger:operation GET /users users listUsers
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
	// swagger:operation GET /users/{id} users getUser
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
	// swagger:operation PATCH /users/{id} users userUpdate
	// ---
	// summary: Updates user's contact information
	// description: Updates user's contact information -> first name, last name, mobile, phone, address.
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
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	ur.PATCH("/:id", u.update)
	// swagger:operation DELETE /users/{id} users userDelete
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

func (u *User) list(c *gin.Context) {
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.List(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listResponse{result, p.Page})
}

func (u *User) view(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	result, err := u.svc.View(c, id)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *User) update(c *gin.Context) {
	usr, err := request.UserUpdate(c)
	if err != nil {
		return
	}
	usrUpdate, err := u.svc.Update(c, &user.Update{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Mobile:    usr.Mobile,
		Phone:     usr.Phone,
		Address:   usr.Address,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, usrUpdate)
}

func (u *User) delete(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	if err := u.svc.Delete(c, id); err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusOK)
}
