package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/internal/mock/mockdb"

	"github.com/ribice/gorsk/internal/auth"
	"github.com/ribice/gorsk/internal/mock"
)

func TestAuthenticate(t *testing.T) {
	type args struct {
		c    context.Context
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
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
					return nil, apperr.DB
				},
			},
		},
		{
			name:    "Fail on hashing",
			args:    args{user: "juzernejm", pass: "notHashedPassword"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
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
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
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
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, time.Time, error) {
					return "", mock.TestTime(1), apperr.Generic
				},
			},
		},
		{
			name:    "Fail on updating last login",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
				UpdateLastLoginFn: func(c context.Context, u *model.User) error {
					return apperr.DB
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, time.Time, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockdb.User{
				FindByUsernameFn: func(c context.Context, user string) (*model.User, error) {
					return &model.User{
						Username: user,
						Password: auth.HashPassword("pass"),
						Active:   true,
					}, nil
				},
				UpdateLastLoginFn: func(c context.Context, u *model.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *model.User) (string, time.Time, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000), nil
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
			s := auth.New(tt.udb, tt.jwt)
			token, err := s.Authenticate(tt.args.c, tt.args.user, tt.args.pass)
			assert.Equal(t, tt.wantData, token)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUser(t *testing.T) {
	ctx := mock.GinCtxWithKeys([]string{
		"id", "company_id", "location_id", "user", "email", "role"},
		9, 15, 52, "ribice", "ribice@gmail.com", int8(1))
	wantUser := &model.AuthUser{
		ID:         9,
		Username:   "ribice",
		CompanyID:  15,
		LocationID: 52,
		Email:      "ribice@gmail.com",
		Role:       model.SuperAdminRole,
	}
	rbacSvc := auth.New(nil, nil)
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
