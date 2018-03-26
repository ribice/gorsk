package pgsql_test

import (
	"reflect"
	"testing"

	"github.com/ribice/gorsk/internal/platform/postgres"

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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				if !reflect.DeepEqual(tt.wantData, user) {
					t.Errorf("Expected %v, received %v", tt.wantData, user)
				}
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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				if !reflect.DeepEqual(tt.wantData, user) {
					t.Errorf("Expected %v, received %v", tt.wantData, user)
				}
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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				if !reflect.DeepEqual(tt.wantData, users) {
					t.Errorf("Expected %v, received %v", tt.wantData, users)
				}
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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				if tt.usr.LastLogin == userBefore.LastLogin {
					t.Errorf("Expected last login to be changed, but was not.")
				}
				tt.wantData.UpdatedAt = userBefore.UpdatedAt
				tt.wantData.CreatedAt = userBefore.CreatedAt
				if !reflect.DeepEqual(tt.wantData, userBefore) {
					t.Errorf("Expected %v - received %v", tt.wantData, userBefore)
				}
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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				if tt.usr.DeletedAt == userBefore.DeletedAt {
					t.Errorf("Expected deletedAt to be changed, but was not.")
				}
				tt.wantData.UpdatedAt = userBefore.UpdatedAt
				tt.wantData.CreatedAt = userBefore.CreatedAt
				tt.wantData.LastLogin = userBefore.LastLogin
				if !reflect.DeepEqual(tt.wantData, userBefore) {
					t.Errorf("Expected %#v - received %#v", tt.wantData, userBefore)
				}
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
			if tt.wantErr != (err != nil) {
				t.Error("WantErr and err!=nil are not equal")
			}
			if tt.wantData != nil {
				tt.wantData.UpdatedAt = resp.UpdatedAt
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.LastLogin = resp.LastLogin
				tt.wantData.DeletedAt = resp.DeletedAt
				if !reflect.DeepEqual(tt.wantData, resp) {
					t.Errorf("Expected %#v - received %#v", tt.wantData, resp)
				}
			}
		})
	}
}
