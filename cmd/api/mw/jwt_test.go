package mw_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/config"
	"github.com/ribice/gorsk/internal/mock"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/mw"
)

func TestAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw.Add(r, gin.Logger())
}

func hwHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"text": "Hello World.",
	})
}

func ginHandler(jwt *mw.JWT) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(jwt.MWFunc())
	r.GET("/hello", hwHandler)
	return r
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		header     string
	}{
		{
			name:       "Empty header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid header",
			header:     mock.HeaderInvalid(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Success",
			header:     mock.HeaderValid(),
			wantStatus: http.StatusOK,
		},
	}
	jwtCfg := &config.JWTConfig{Realm: "testRealm", Secret: "jwtsecret", Timeout: 60, SigningAlgorithm: "HS256"}
	ts := httptest.NewServer(ginHandler(mw.NewJWT(jwtCfg)))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Failed creating request")
			}
			defer res.Body.Close()
			if res.StatusCode != tt.wantStatus {
				t.Errorf("expected status %v; got %v", tt.wantStatus, res.StatusCode)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		name      string
		wantToken string
		req       *model.User
	}{
		{
			name: "Success",
			req: &model.User{
				Base: model.Base{
					ID: 1,
				},
				Username: "johndoe",
				Email:    "johndoe@mail.com",
				Role: &model.Role{
					AccessLevel: model.SuperAdminRole,
				},
				CompanyID:  1,
				LocationID: 1,
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}
	jwtCfg := &config.JWTConfig{Realm: "testRealm", Secret: "jwtsecret", Timeout: 60, SigningAlgorithm: "HS256"}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			jwt := mw.NewJWT(jwtCfg)
			str, _, err := jwt.GenerateToken(tt.req)
			if err != nil {
				t.Error("Didn't expect error but received one.")
			}
			if strings.Split(str, ".")[0] != tt.wantToken {
				t.Error("Expected and received token do not match.")
			}

		})
	}
}
