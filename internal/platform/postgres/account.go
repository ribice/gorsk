package pgsql

import (
	"context"
	"net/http"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

// NewAccountDB returns a new AccountDB instance
func NewAccountDB(c *pg.DB, l *zap.Logger) *AccountDB {
	return &AccountDB{c, l}
}

// AccountDB represents the client for user table
type AccountDB struct {
	cl  *pg.DB
	log *zap.Logger
}

// Create creates a new user on database
func (a *AccountDB) Create(c context.Context, usr *model.User) error {
	var user = new(model.User)
	res, err := a.cl.Query(user, "select id from users where username = ? or email = ? and deleted_at is null", usr.Username, usr.Email)
	if err != nil {
		a.log.Error("AccountDB Error: ", zap.Error(err))
		return apperr.DB
	}
	if res.RowsReturned() != 0 {
		return apperr.New(http.StatusBadRequest, "Username or email already exists.")
	}

	if err := a.cl.Insert(usr); err != nil {
		a.log.Warn("AccountDB Error: ", zap.Error(err))
		return apperr.DB
	}
	return nil
}

// ChangePassword changes user's password
func (a *AccountDB) ChangePassword(c context.Context, usr *model.User) error {
	_, err := a.cl.Model(usr).Column("password", "updated_at").Update()
	if err != nil {
		a.log.Warn("AccountDB Error: ", zap.Error(err))
	}
	return err
}
