package mw

import (
	"net/http"
	"strings"
	"time"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// NewJWT generates new JWT variable necessery for auth middleware
func NewJWT(c *config.JWTConfig) *JWT {
	return &JWT{
		Realm:   c.Realm,
		Key:     []byte(c.Secret),
		Timeout: time.Duration(c.Timeout) * time.Minute,
		Algo:    c.SigningAlgorithm,
	}
}

// JWT provides a Json-Web-Token authentication implementation
type JWT struct {
	// Realm name to display to the user.
	Realm string

	// Secret key used for signing.
	Key []byte

	// Duration for which the jwt token is valid.
	Timeout time.Duration

	// JWT signing algorithm
	Algo string
}

// MWFunc makes JWT implement the Middleware interface.
func (j *JWT) MWFunc() gin.HandlerFunc {

	return func(c *gin.Context) {
		token, err := j.ParseToken(c)
		if err != nil {
			c.Header("WWW-Authenticate", "JWT realm="+j.Realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id := int(claims["id"].(float64))
		companyID := int(claims["company_id"].(float64))
		locationID := int(claims["location_id"].(float64))
		username := claims["user"].(string)
		email := claims["email"].(string)
		role := int8(claims["role"].(float64))

		c.Set("id", id)
		c.Set("company_id", companyID)
		c.Set("location_id", locationID)
		c.Set("user", username)
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

// ParseToken parses token from Authorization header
func (j *JWT) ParseToken(c *gin.Context) (*jwt.Token, error) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return nil, apperr.Generic
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, apperr.Generic
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(j.Algo) != token.Method {
			return nil, apperr.Generic
		}
		return j.Key, nil
	})
}

// GenerateToken generates new JWT token and populates it with user data
func (j *JWT) GenerateToken(u *model.User) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(j.Algo))
	claims := token.Claims.(jwt.MapClaims)

	expire := time.Now().Add(j.Timeout)
	claims["id"] = u.ID
	claims["user"] = u.Username
	claims["email"] = u.Email
	claims["role"] = u.Role.AccessLevel
	claims["company_id"] = u.CompanyID
	claims["location_id"] = u.LocationID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(j.Key)
	return tokenString, expire, err
}
