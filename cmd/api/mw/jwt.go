package mw

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/config"

	jwt "github.com/dgrijalva/jwt-go"
)

// NewJWT generates new JWT variable necessery for auth middleware
func NewJWT(c *config.JWT) *JWT {
	return &JWT{
		Key:      []byte(c.Secret),
		Duration: time.Duration(c.Duration) * time.Minute,
		Algo:     c.SigningAlgorithm,
	}
}

// JWT provides a Json-Web-Token authentication implementation
type JWT struct {
	// Secret key used for signing.
	Key []byte

	// Duration for which the jwt token is valid.
	Duration time.Duration

	// JWT signing algorithm
	Algo string
}

// MWFunc makes JWT implement the Middleware interface.
func (j *JWT) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseToken(c)
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)

			id := int(claims["id"].(float64))
			companyID := int(claims["c"].(float64))
			locationID := int(claims["l"].(float64))
			username := claims["u"].(string)
			email := claims["e"].(string)
			role := int8(claims["r"].(float64))

			c.Set("id", id)
			c.Set("company_id", companyID)
			c.Set("location_id", locationID)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("role", role)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (j *JWT) ParseToken(c echo.Context) (*jwt.Token, error) {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, model.ErrGeneric
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, model.ErrGeneric
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(j.Algo) != token.Method {
			return nil, model.ErrGeneric
		}
		return j.Key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (j *JWT) GenerateToken(u *model.User) (string, string, error) {
	expire := time.Now().Add(j.Duration)

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.Algo), jwt.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   u.Role.AccessLevel,
		"c":   u.CompanyID,
		"l":   u.LocationID,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.Key)

	return tokenString, expire.Format(time.RFC3339), err
}
