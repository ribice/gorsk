package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// New creates new auth service
func New(udb model.UserDB, j JWT) *Service {
	return &Service{
		udb: udb,
		jwt: j,
	}
}

// Service represents auth application service
type Service struct {
	udb model.UserDB
	jwt JWT
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, time.Time, error)
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c context.Context, user, pass string) (*model.AuthToken, error) {
	u, err := s.udb.FindByUsername(c, user)
	if err != nil {
		return nil, err
	}
	if !HashMatchesPassword(u.Password, pass) {
		return nil, apperr.New(http.StatusNotFound, "Username or password does not exist")
	}

	if !u.Active {
		return nil, apperr.New(http.StatusUnauthorized, "User is not active")
	}

	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.Forbidden
	}

	u.UpdateLastLogin()
	if err := s.udb.UpdateLastLogin(c, u); err != nil {
		return nil, err
	}

	return &model.AuthToken{Token: token, Expires: expire.Format(time.RFC3339)}, nil
}

// User returns user data stored in jwt token
func (s Service) User(c *gin.Context) *model.AuthUser {
	id := c.GetInt("id")
	companyID := c.GetInt("company_id")
	locationID := c.GetInt("location_id")
	user := c.GetString("user")
	email := c.GetString("email")
	role := c.MustGet("role").(int8)
	return &model.AuthUser{
		ID:         id,
		Username:   user,
		CompanyID:  companyID,
		LocationID: locationID,
		Email:      email,
		Role:       model.AccessRole(role),
	}
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
