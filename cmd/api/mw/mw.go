package mw

import (
	"github.com/gin-gonic/gin"
)

// Add adds middlewares to gin engine
func Add(r *gin.Engine, h ...gin.HandlerFunc) {
	for _, v := range h {
		r.Use(v)
	}
}
