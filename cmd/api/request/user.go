package request

import (
	"github.com/labstack/echo"
)

// UpdateUser contains user update data from json request
type UpdateUser struct {
	ID        int     `json:"-"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Mobile    *string `json:"mobile,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
}

// UserUpdate validates user update request
func UserUpdate(c echo.Context) (*UpdateUser, error) {
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	u := new(UpdateUser)
	if err := c.Bind(u); err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}
