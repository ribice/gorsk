package pgsql

import (
	"net/http"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"
)

// NewAccountDB returns a new AccountDB instance
func NewAccountDB(l echo.Logger) *AccountDB {
	return &AccountDB{l}
}

// AccountDB represents the client for user table
type AccountDB struct {
	log echo.Logger
}

// Create creates a new user on database
func (a *AccountDB) Create(db orm.DB, usr model.User) (*model.User, error) {
	var user = new(model.User)
	res, err := db.Query(user, "select id from users where username = ? or email = ? and deleted_at is null", usr.Username, usr.Email)
	if err != nil {
		a.log.Error("AccountDB Error: %v", err)
		return nil, err
	}
	if res.RowsReturned() != 0 {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
	}
	if err := db.Insert(&usr); err != nil {
		a.log.Error("AccountDB Error: %v", err)
		return nil, err
	}
	return &usr, nil
}

// ChangePassword changes user's password
func (a *AccountDB) ChangePassword(db orm.DB, usr *model.User) error {
	_, err := db.Model(usr).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warnf("AccountDB Error: %v", err)
	}
	return err
}
