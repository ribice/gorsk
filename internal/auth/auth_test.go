package auth_test

import (
	"testing"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/auth"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/ribice/gorsk/internal/mock/mockdb"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	type args struct {
		c    echo.Context
		user string
		pass string
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.AuthToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
	}{
		{
			name:    "Fail on finding user",
			args:    args{user: "juzernejm"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return nil, model.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on hashing",
			args:    args{user: "juzernejm", pass: "notHashedPassword"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: "HashedPassword",
					}, nil
				},
			},
		},
		{
			name:    "Inactive user",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   false,
					}, nil
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, string, error) {
					return "", "", model.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on updating last login",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u *model.User) (*model.User, error) {
					return nil, model.ErrGeneric
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u *model.User) (*model.User, error) {
					return u, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			wantData: &model.AuthToken{
				Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires: mock.TestTime(2000).Format(time.RFC3339),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt)
			token, err := s.Authenticate(tt.args.c, tt.args.user, tt.args.pass)
			if tt.wantData != nil {
				tt.wantData.RefreshToken = token.RefreshToken
				assert.Equal(t, tt.wantData, token)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestRefresh(t *testing.T) {
	type args struct {
		c     echo.Context
		token string
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.RefreshToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
	}{
		{
			name:    "Fail on finding token",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*model.User, error) {
					return nil, model.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*model.User, error) {
					return &model.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, string, error) {
					return "", "", model.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{token: "refreshtoken"},
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*model.User, error) {
					return &model.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			wantData: &model.RefreshToken{
				Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires: mock.TestTime(2000).Format(time.RFC3339),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt)
			token, err := s.Refresh(tt.args.c, tt.args.token)
			assert.Equal(t, tt.wantData, token)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "company_id", "location_id", "username", "email", "role"},
		9, 15, 52, "ribice", "ribice@gmail.com", int8(1))
	wantUser := &model.AuthUser{
		ID:         9,
		Username:   "ribice",
		CompanyID:  15,
		LocationID: 52,
		Email:      "ribice@gmail.com",
		Role:       model.SuperAdminRole,
	}
	rbacSvc := auth.New(nil, nil, nil)
	assert.Equal(t, wantUser, rbacSvc.User(ctx))
}

func TestHashPassowrd(t *testing.T) {
	password := "Hunter123"
	hash := auth.HashPassword(password)
	if password == hash {
		t.Error("Passsword and hash should not be equal")
	}
}

func TestHashMatchesPassword(t *testing.T) {
	password := "Hunter123"
	if !auth.HashMatchesPassword(auth.HashPassword(password), password) {
		t.Error("Passsword and hash should match")
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name     string
		ctx      echo.Context
		wantData *model.User
		udb      *mockdb.User
		wantErr  bool
	}{
		{
			name: "Success",
			ctx: mock.EchoCtxWithKeys([]string{
				"id", "company_id", "location_id", "username", "email", "role"},
				9, 15, 52, "ribice", "ribice@gmail.com", int8(1)),
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &model.Role{
							AccessLevel: model.UserRole,
						},
					}, nil
				},
			},
			wantData: &model.User{
				Base: model.Base{
					ID:        9,
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Role: &model.Role{
					AccessLevel: model.UserRole,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, nil)
			user, err := s.Me(tt.ctx)
			assert.Equal(t, tt.wantData, user)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
