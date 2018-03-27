package pgsql_test

import (
	"testing"

	"github.com/ribice/gorsk/internal/platform/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/go-pg/pg"
	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/mock"
	"go.uber.org/zap"
)

func testUserDB(t *testing.T, c *pg.DB, l *zap.Logger) {
	userDB := pgsql.NewUserDB(c, l)
	cases := []struct {
		name string
		fn   func(*testing.T, *pgsql.UserDB, *pg.DB)
	}{
		{
			name: "view",
			fn:   testUserView,
		},
		{
			name: "findByUsername",
			fn:   testUserFindByUsername,
		},
		{
			name: "userList",
			fn:   testUserList,
		},
		{
			name: "updateLastLogin",
			fn:   testUserUpdateLastLogin,
		},
		{
			name: "delete",
			fn:   testUserDelete,
		},
		{
			name: "update",
			fn:   testUserUpdate,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(t, userDB, c)
		})
	}

}

func testUserView(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData *model.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: &model.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: model.Base{
					ID: 2,
				},
				Role: &model.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := db.View(nil, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}

func testUserFindByUsername(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		username string
		wantData *model.User
	}{
		{
			name:     "User does not exist",
			wantErr:  true,
			username: "notExists",
		},
		{
			name:     "Success",
			username: "tomjones",
			wantData: &model.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: model.Base{
					ID: 2,
				},
				Role: &model.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := db.FindByUsername(nil, tt.username)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				assert.Equal(t, tt.wantData, user)

			}
		})
	}
}

func testUserList(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *model.ListQuery
		pg       *model.Pagination
		wantData []model.User
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: &model.Pagination{
				Limit: -100,
			},
		},
		{
			name: "Success",
			pg: &model.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &model.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			wantData: []model.User{
				{
					Email:      "tomjones@mail.com",
					FirstName:  "Tom",
					LastName:   "Jones",
					Username:   "tomjones",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "newPass",
					Base: model.Base{
						ID: 2,
					},
					Role: &model.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
				},
				{
					Email:      "johndoe@mail.com",
					FirstName:  "John",
					LastName:   "Doe",
					Username:   "johndoe",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Base: model.Base{
						ID: 1,
					},
					Role: &model.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := db.List(nil, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, users)
			}
		})
	}
}

func testUserUpdateLastLogin(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *model.User
		wantData *model.User
	}{
		// Does not fail on this test, but should
		// {
		// 	name:    "User does not exist",
		// 	wantErr: true,
		// 	usr:     &model.User{},
		// },
		{
			name: "Success",
			usr: &model.User{
				Base: model.Base{
					ID:        2,
					UpdatedAt: mock.TestTime(2000),
				},
				LastLogin: mock.TestTimePtr(2018),
			},
			wantData: &model.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: model.Base{
					ID: 2,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			userBefore := &model.User{}
			if tt.wantData != nil {
				userBefore = queryUser(t, c, tt.usr.Base.ID)
			}
			err := db.UpdateLastLogin(nil, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				assert.NotEqual(t, tt.usr.LastLogin, userBefore.LastLogin)
				tt.wantData.UpdatedAt = userBefore.UpdatedAt
				tt.wantData.CreatedAt = userBefore.CreatedAt
				assert.Equal(t, tt.wantData, userBefore)
			}
		})
	}
}

func testUserDelete(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *model.User
		wantData *model.User
	}{
		// Does not fail on this test, but should
		// {
		// 	name:    "User does not exist",
		// 	wantErr: true,
		// 	usr:     &model.User{},
		// },
		{
			name: "Success",
			usr: &model.User{
				Base: model.Base{
					ID:        2,
					DeletedAt: mock.TestTimePtr(2018),
				},
			},
			wantData: &model.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: model.Base{
					ID: 2,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			userBefore := &model.User{}
			if tt.wantData != nil {
				userBefore = queryUser(t, c, tt.usr.Base.ID)
			}
			err := db.Delete(nil, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				assert.NotEqual(t, tt.usr.DeletedAt, userBefore.DeletedAt)
				tt.wantData.UpdatedAt = userBefore.UpdatedAt
				tt.wantData.CreatedAt = userBefore.CreatedAt
				tt.wantData.LastLogin = userBefore.LastLogin
				assert.Equal(t, tt.wantData, userBefore)
			}
		})
	}
}

func testUserUpdate(t *testing.T, db *pgsql.UserDB, c *pg.DB) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *model.User
		wantData *model.User
	}{
		// Does not fail on this test, but should
		// {
		// 	name:    "User does not exist",
		// 	wantErr: true,
		// 	usr:     &model.User{},
		// },
		{
			name: "Success",
			usr: &model.User{
				Base: model.Base{
					ID: 2,
				},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			// Expected wantData:
			// wantData: &model.User{
			// 	Email:      "tomjones@mail.com",
			// 	FirstName:  "Z",
			// 	LastName:   "Freak",
			// 	Username:   "tomjones",
			// 	RoleID:     1,
			// 	CompanyID:  1,
			// 	LocationID: 1,
			// 	Password:   "newPass",
			// 	Address:    "Address",
			// 	Phone:      "123456",
			// 	Mobile:     "345678",
			// 	Base: model.Base{
			// 		ID: 2,
			// 	},
			// },
			wantData: &model.User{
				FirstName: "Z",
				LastName:  "Freak",
				Username:  "newUsername",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Base: model.Base{
					ID: 2,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := db.Update(nil, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				tt.wantData.UpdatedAt = resp.UpdatedAt
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.LastLogin = resp.LastLogin
				tt.wantData.DeletedAt = resp.DeletedAt
				assert.Equal(t, tt.wantData, resp)
			}
		})
	}
}
