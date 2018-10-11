package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/server"
	"github.com/ribice/gorsk/cmd/api/service"
	"github.com/ribice/gorsk/internal/account"
	"github.com/ribice/gorsk/internal/auth"

	"github.com/ribice/gorsk/internal/mock"
	"github.com/ribice/gorsk/internal/mock/mockdb"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *model.User
		adb        *mockdb.Account
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on userSvc",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":2}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID model.AccessRole, companyID, locationID int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":2}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID model.AccessRole, companyID, locationID int) error {
					return nil
				},
			},
			adb: &mockdb.Account{
				CreateFn: func(db orm.DB, usr model.User) (*model.User, error) {
					usr.ID = 1
					usr.CreatedAt = mock.TestTime(2018)
					usr.UpdatedAt = mock.TestTime(2018)
					return &usr, nil
				},
			},
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				},
				FirstName:  "John",
				LastName:   "Doe",
				Username:   "juzernejm",
				Email:      "johndoe@gmail.com",
				CompanyID:  1,
				LocationID: 2,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewAccount(account.New(nil, tt.adb, nil, tt.rbac), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestChangePassword(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		id         string
		udb        *mockdb.User
		adb        *mockdb.Account
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return echo.ErrForbidden
				},
			},
			id:         "1",
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				},
			},
			id: "1",
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*model.User, error) {
					return &model.User{
						Password: auth.HashPassword("oldpassw"),
					}, nil
				},
			},
			adb: &mockdb.Account{
				ChangePasswordFn: func(db orm.DB, usr *model.User) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewAccount(account.New(nil, tt.adb, tt.udb, tt.rbac), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id + "/password"
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
