package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/gorsk/internal/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/internal/user"

	"github.com/ribice/gorsk/internal"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/service"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/ribice/gorsk/internal/mock/mockdb"
)

func TestListUsers(t *testing.T) {
	type listResponse struct {
		Users []model.User `json:"users"`
		Page  int          `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `?limit=2222&page=-1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on query list",
			req:  `?limit=100&page=1`,
			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.UserRole,
						Email:      "john@mail.com",
					}
				}},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `?limit=100&page=1`,
			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.SuperAdminRole,
						Email:      "john@mail.com",
					}
				}},
			udb: &mockdb.User{
				ListFn: func(c context.Context, q *model.ListQuery, p *model.Pagination) ([]model.User, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []model.User{
							{
								Base: model.Base{
									ID:        10,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								},
								FirstName:  "John",
								LastName:   "Doe",
								Email:      "john@mail.com",
								CompanyID:  2,
								LocationID: 3,
								Role: &model.Role{
									ID:          1,
									AccessLevel: 1,
									Name:        "SUPER_ADMIN",
								},
							},
							{
								Base: model.Base{
									ID:        11,
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
								},
								FirstName:  "Joanna",
								LastName:   "Dye",
								Email:      "joanna@mail.com",
								CompanyID:  1,
								LocationID: 2,
								Role: &model.Role{
									ID:          2,
									AccessLevel: 2,
									Name:        "ADMIN",
								},
							},
						}, nil
					}
					return nil, apperr.DB
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Users: []model.User{
					{
						Base: model.Base{
							ID:        10,
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						},
						FirstName:  "John",
						LastName:   "Doe",
						Email:      "john@mail.com",
						CompanyID:  2,
						LocationID: 3,
						Role: &model.Role{
							ID:          1,
							AccessLevel: 1,
							Name:        "SUPER_ADMIN",
						},
					},
					{
						Base: model.Base{
							ID:        11,
							CreatedAt: mock.TestTime(2004),
							UpdatedAt: mock.TestTime(2005),
						},
						FirstName:  "Joanna",
						LastName:   "Dye",
						Email:      "joanna@mail.com",
						CompanyID:  1,
						LocationID: 2,
						Role: &model.Role{
							ID:          2,
							AccessLevel: 2,
							Name:        "ADMIN",
						},
					},
				}, Page: 1},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestViewUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
			},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.req
			res, err := http.Get(path)
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

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(*gin.Context, int) bool {
					return true
				},
			},
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
						Address:   "Work",
						Phone:     "332223",
					}, nil
				},
				UpdateFn: func(c context.Context, usr *model.User) (*model.User, error) {
					usr.UpdatedAt = mock.TestTime(2010)
					usr.Mobile = "991991"
					return usr, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2010),
				},
				FirstName: "jj",
				LastName:  "okocha",
				Username:  "JohnDoe",
				Phone:     "321321",
				Address:   "home",
				Mobile:    "991991",
			},
		},
	}
	gin.SetMode(gin.TestMode)
	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			res, err := client.Do(req)
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

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.CompanyAdminRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.CompanyAdminRole,
						},
					}, nil
				},
				DeleteFn: func(context.Context, *model.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return true
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)
	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
