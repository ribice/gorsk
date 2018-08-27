package pgsql

import (
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"
)

// NewUserDB returns a new UserDB instance
func NewUserDB(l echo.Logger) *UserDB {
	return &UserDB{l}
}

// UserDB represents the client for user table
type UserDB struct {
	log echo.Logger
}

// View returns single user by ID
func (u *UserDB) View(db orm.DB, id int) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, id)
	if err != nil {
		u.log.Warnf("AccountDB Error: %v", err)
	}
	return user, err
}

// FindByUsername queries for single user by username
func (u *UserDB) FindByUsername(db orm.DB, uname string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, uname)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}

// FindByToken queries for single user by token
func (u *UserDB) FindByToken(db orm.DB, token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}

// List returns list of all users retrievable for the current user, depending on role
func (u *UserDB) List(db orm.DB, qp *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	var users []model.User
	q := db.Model(&users).Column("user.*", "Role").Limit(p.Limit).Offset(p.Offset).Where(notDeleted).Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warnf("UserDB Error: %v", err)
		return nil, err
	}
	return users, nil
}

// Delete sets deleted_at for a user
func (u *UserDB) Delete(db orm.DB, user *model.User) error {
	err := db.Delete(user)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return err
}

// Update updates user's contact info
func (u *UserDB) Update(db orm.DB, user *model.User) (*model.User, error) {
	_, err := db.Model(user).WherePK().Update()
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}
