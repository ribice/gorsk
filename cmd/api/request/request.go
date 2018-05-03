package request

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

const (
	defaultLimit = 100
	maxLimit     = 1000
)

// Pagination contains pagination request
type Pagination struct {
	Limit  int `query:"limit"`
	Page   int `query:"page" validate:"min=0"`
	Offset int `json:"-"`
}

// Paginate validates pagination requests
func Paginate(c echo.Context) (*Pagination, error) {
	p := new(Pagination)
	if err := c.Bind(p); err != nil {
		return nil, err
	}
	if p.Limit < 1 {
		p.Limit = defaultLimit
	}
	if p.Limit > 1000 {
		p.Limit = maxLimit
	}
	p.Offset = p.Limit * p.Page
	return p, nil
}

// ID returns id url parameter.
// In case of conversion error to int, StatusBadRequest will be returned as err
func ID(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest)
	}
	return id, nil
}
