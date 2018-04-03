package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/satori/go.uuid"

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
	GenerateToken(*model.User) (string, string, error)
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
		return nil, apperr.Unauthorized
	}
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.Unauthorized
	}

	u.UpdateLastLogin()
	u.Token = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	if err := s.udb.UpdateLogin(c, u); err != nil {
		return nil, err
	}

	return &model.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (s *Service) Refresh(c context.Context, token string) (*model.RefreshToken, error) {
	user, err := s.udb.FindByToken(c, token)
	if err != nil {
		return nil, err
	}
	token, expire, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, apperr.Generic
	}
	return &model.RefreshToken{Token: token, Expires: expire}, nil
}

// User returns user data stored in jwt token
func (s *Service) User(c *gin.Context) *model.AuthUser {
	id := c.GetInt("id")
	companyID := c.GetInt("company_id")
	locationID := c.GetInt("location_id")
	user := c.GetString("username")
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
