package model

import (
	"errors"
	"time"

	"github.com/go-pg/pg/orm"
)

// ErrGeneric is used for testing purposes and for errors handled later in the callstack
var ErrGeneric = errors.New("generic error")

// Base contains common fields for all tables
type Base struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}

// Pagination holds paginations data
type Pagination struct {
	Limit  int
	Offset int
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    int
}

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeInsert(_ orm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(_ orm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
