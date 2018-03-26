package request

import (
	"net/http"
	"strconv"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/gin-gonic/gin"
)

const defaultLimit = 100
const maxLimit = 1000

// Pagination contains pagination request
type Pagination struct {
	Limit  int `form:"limit"`
	Page   int `form:"page" binding:"min=0"`
	Offset int `json:"-"`
}

// Paginate validates pagination requests
func Paginate(c *gin.Context) (*Pagination, error) {
	p := new(Pagination)
	if err := c.ShouldBindQuery(p); err != nil {
		apperr.Response(c, err)
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
// In case of conversion error to int, request will be aborted with StatusBadRequest.
func ID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return 0, apperr.BadRequest
	}
	return id, nil
}
