package model

import (
	"time"

	"github.com/go-pg/pg/orm"
)

// User represents user domain model
type User struct {
	Base
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Address   string     `json:"address,omitempty"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Active    bool       `json:"active"`
	Token     string     `json:"-"`

	Role *Role `json:"role,omitempty"`

	RoleID     AccessRole `json:"-"`
	CompanyID  int        `json:"company_id"`
	LocationID int        `json:"location_id"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID         int
	CompanyID  int
	LocationID int
	Username   string
	Email      string
	Role       AccessRole
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin() {
	t := time.Now()
	u.LastLogin = &t
}

// AccountDB represents account related database interface (repository)
type AccountDB interface {
	Create(orm.DB, User) (*User, error)
	ChangePassword(orm.DB, *User) error
}

// UserDB represents user database interface (repository)
type UserDB interface {
	View(orm.DB, int) (*User, error)
	FindByUsername(orm.DB, string) (*User, error)
	FindByToken(orm.DB, string) (*User, error)
	List(orm.DB, *ListQuery, *Pagination) ([]User, error)
	Delete(orm.DB, *User) error
	Update(orm.DB, *User) (*User, error)
}
