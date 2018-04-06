package request

import (
	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal/errors"
)

// UpdateUser contains user update data from json request
type UpdateUser struct {
	ID        int     `json:"-"`
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2"`
	Mobile    *string `json:"mobile,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
}

// UserUpdate validates user update request
func UserUpdate(c *gin.Context) (*UpdateUser, error) {
	var u UpdateUser
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(&u); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	u.ID = id
	return &u, nil
}
