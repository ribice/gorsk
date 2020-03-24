package transport

import (
	"github.com/ribice/gorsk"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*gorsk.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []gorsk.User `json:"users"`
		Page  int          `json:"page"`
	}
}
