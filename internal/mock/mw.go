package mock

import (
	"time"

	"github.com/ribice/gorsk/internal"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*model.User) (string, time.Time, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *model.User) (string, time.Time, error) {
	return j.GenerateTokenFn(u)
}
