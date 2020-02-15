package pgsql

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"

	gorsk2 "github.com/ribice/gorsk"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user on database
func (u *User) Create(db orm.DB, usr gorsk2.User) (*gorsk2.User, error) {
	var user = new(gorsk2.User)
	err := db.Model(user).Where("lower(username) = ? or lower(email) = ? and deleted_at is null",
		strings.ToLower(usr.Username), strings.ToLower(usr.Email)).Select()
	if (err == nil) || (err != nil && err != pg.ErrNoRows) {
		return nil, ErrAlreadyExists
	}

	if err := db.Insert(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*gorsk2.User, error) {
	var user = new(gorsk2.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates user's contact info
func (u *User) Update(db orm.DB, user *gorsk2.User) error {
	_, err := db.Model(user).UpdateNotNull()
	return err
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db orm.DB, qp *gorsk2.ListQuery, p *gorsk2.Pagination) ([]gorsk2.User, error) {
	var users []gorsk2.User
	q := db.Model(&users).Column("user.*", "Role").Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return users, nil
}

// Delete sets deleted_at for a user
func (u *User) Delete(db orm.DB, user *gorsk2.User) error {
	return db.Delete(user)
}
