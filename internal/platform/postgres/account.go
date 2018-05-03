package pgsql

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"

	"github.com/go-pg/pg"
)

// NewAccountDB returns a new AccountDB instance
func NewAccountDB(c *pg.DB, l echo.Logger) *AccountDB {
	return &AccountDB{c, l}
}

// AccountDB represents the client for user table
type AccountDB struct {
	cl  *pg.DB
	log echo.Logger
}

// Create creates a new user on database
func (a *AccountDB) Create(usr model.User) (*model.User, error) {
	var user = new(model.User)
	res, err := a.cl.Query(user, "select id from users where username = ? or email = ? and deleted_at is null", usr.Username, usr.Email)
	if err != nil {
		a.log.Error("AccountDB Error: %v", err)
		return nil, err
	}
	if res.RowsReturned() != 0 {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
	}
	if err := a.cl.Insert(&usr); err != nil {
		a.log.Error("AccountDB Error: %v", err)
		return nil, err
	}
	return &usr, nil
}

// ChangePassword changes user's password
func (a *AccountDB) ChangePassword(usr *model.User) error {
	_, err := a.cl.Model(usr).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warnf("AccountDB Error: %v", err)
	}
	return err
}
