package apperr_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal/errors"
	validator "gopkg.in/go-playground/validator.v8"
)

func TestNewStatus(t *testing.T) {
	apr := apperr.NewStatus(http.StatusBadRequest)
	if apr.Message != "" || apr.Status != http.StatusBadRequest {
		t.Errorf("Invalid error received.")
	}
}

func TestNew(t *testing.T) {
	apr := apperr.New(http.StatusBadRequest, "Bad request")
	if apr.Message != "Bad request" || apr.Status != http.StatusBadRequest {
		t.Errorf("Invalid error received.")
	}
}

func TestError(t *testing.T) {
	apr := apperr.New(http.StatusBadRequest, "Bad request")
	if apr.Error() != "Bad request" {
		t.Errorf("Invalid error received.")
	}
}

func TestResponse(t *testing.T) {
	type args struct {
		err error
	}
	type validationStr struct {
		FirstName string `json:"first_name" validation:"required"`
		LastName  string `json:"last_name" validation:"required"`
		Email     string `json:"email" validation:"required,email"`
	}
	cases := []struct {
		name       string
		args       args
		wantStatus int
		wantResp   string
		vld        *validationStr
	}{
		{
			name:       "AppErr without message",
			args:       args{err: apperr.Forbidden},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "AppErr with message",
			args:       args{err: apperr.New(http.StatusBadRequest, "bad request")},
			wantStatus: http.StatusBadRequest,
			wantResp:   `{"message":"bad request"}`,
		},
		{
			name:       "Validator message",
			wantStatus: http.StatusBadRequest,
			wantResp:   `{"message":["Email is required, but was not received"]}`,
			vld:        &validationStr{FirstName: "Emir", LastName: "Ribic"},
		},
		{
			name:       "Validator message other",
			wantStatus: http.StatusBadRequest,
			wantResp:   `{"message":["Email failed on email validation"]}`,
			vld:        &validationStr{FirstName: "Emir", LastName: "Ribic", Email: "testing"},
		},
		{
			name:       "Other error",
			args:       args{err: fmt.Errorf("An error occurred")},
			wantStatus: http.StatusInternalServerError,
			wantResp:   `{"message":"An error occurred"}`,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.vld != nil {
				config := &validator.Config{TagName: "validation"}
				validate := validator.New(config)
				tt.args.err = validate.Struct(tt.vld)
			}
			apperr.Response(c, tt.args.err)
			if tt.wantStatus != w.Code {
				t.Errorf("Expected status %v, received %v", tt.wantStatus, w.Code)
			}
			if tt.wantResp != "" && tt.wantResp != w.Body.String() {
				t.Errorf("Expected response %v, received %v", tt.wantResp, w.Body.String())
			}
			if !c.IsAborted() {
				t.Error("Expected context to be aborted, but was not")
			}
		})
	}
}
