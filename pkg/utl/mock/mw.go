package mock

import (
	"github.com/ribice/gorsk"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(gorsk.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u gorsk.User) (string, error) {
	return j.GenerateTokenFn(u)
}
