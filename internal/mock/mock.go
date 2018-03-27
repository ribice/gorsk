package mock

import (
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

// TestTime is used for testing time fields
func TestTime(year int) time.Time {
	return time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
}

// TestTimePtr is used for testing pointer time fields
func TestTimePtr(year int) *time.Time {
	t := time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
	return &t
}

// Str2Ptr converts string to pointer
func Str2Ptr(s string) *string {
	return &s
}

// GinCtxWithKeys returns new gin context with keys
func GinCtxWithKeys(keys []string, values ...interface{}) *gin.Context {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	for i, k := range keys {
		c.Set(k, values[i])
	}
	return c
}

// HeaderValid is used for jwt testing
func HeaderValid() string {
	return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlciI6ImpvaG5kb2UiLCJlbWFpbCI6ImpvaG5kb2VAbWFpbC5jb20iLCJyb2xlIjoxLCJjb21wYW55X2lkIjoxLCJsb2NhdGlvbl9pZCI6MSwiZXhwIjo0MTA5MzIwODk0LCJpYXQiOjE1MTYyMzkwMjJ9.BvdLN0EDQA6stRcePjYRX_ag0b_m7deKqWeBBOAoqqk"
}

// HeaderInvalid is used for jwt testing
func HeaderInvalid() string {
	return "Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlciI6ImpvaG5kb2UiLCJlbWFpbCI6ImpvaG5kb2VAbWFpbC5jb20iLCJyb2xlIjoxLCJjb21wYW55X2lkIjoxLCJsb2NhdGlvbl9pZCI6MSwiZXhwIjo0MTA5MzIwODk0LCJpYXQiOjE1MTYyMzkwMjJ9.9fvF4eR0lqcNekPVrSLhlSXe-JFvKhlR0jaSb7pRDwipkcAAt3-GiPGZOmTtmvM_"
}
