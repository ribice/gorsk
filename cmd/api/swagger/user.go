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
// swagger:model pwChange
type swaggPwChange struct {
	request.Password
}

// User update request
// swagger:model userUpdate
type swaggUserUpdateReq struct {
	request.UpdateUser
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
