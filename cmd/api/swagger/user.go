package swagger

import (
	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/request"
)

// Account create request
// swagger:parameters accCreate
type swaggAccCreateReq struct {
	// in:body
	Body request.Register
}

// Password change request
// swagger:parameters pwChange
type swaggPwChange struct {
	// in:body
	Body request.Password
}

// User update request
// swagger:parameters userUpdate
type swaggUserUpdateReq struct {
	// in:body
	Body request.UpdateUser
}

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*model.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []model.User `json:"users"`
		Page  int          `json:"page"`
	}
}
