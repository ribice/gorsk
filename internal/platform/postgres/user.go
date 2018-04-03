package pgsql

import (
	"context"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/errors"

	"go.uber.org/zap"

	"github.com/go-pg/pg"
)

// NewUserDB returns a new UserDB instance
func NewUserDB(c *pg.DB, l *zap.Logger) *UserDB {
	return &UserDB{c, l}
}

// UserDB represents the client for user table
type UserDB struct {
	cl  *pg.DB
	log *zap.Logger
}

// View returns single user by ID
func (u *UserDB) View(c context.Context, id int) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, id)
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
		return nil, apperr.NotFound
	}
	return user, nil
}

// FindByUsername queries for single user by username
func (u *UserDB) FindByUsername(c context.Context, uname string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, uname)
	if err != nil {
		u.log.Warn("UserDB Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *UserDB) FindByToken(c context.Context, token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warn("UserDB Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return user, nil
}

// List returns list of all users retreivable for the current user, depending on role
func (u *UserDB) List(c context.Context, qp *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	var users []model.User
	q := u.cl.Model(&users).Column("user.*", "Role").Limit(p.Limit).Offset(p.Offset).Where(notDeleted).Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// UpdateLogin updates last login and refresh token for user
func (u *UserDB) UpdateLogin(c context.Context, user *model.User) error {
	_, err := u.cl.Model(user).Column("last_login", "token").Update()
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
	}
	return err
}

// Delete sets deleted_at for a user
func (u *UserDB) Delete(c context.Context, user *model.User) error {
	_, err := u.cl.Model(user).Column("deleted_at").Update()
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
	}
	return err
}

// Update updates user's contact info
func (u *UserDB) Update(c context.Context, user *model.User) (*model.User, error) {
	_, err := u.cl.Model(user).Column("first_name",
		"last_name", "mobile", "phone", "address", "updated_at").Update()
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
	}
	return user, err
}
