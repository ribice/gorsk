package model

import (
	"time"

	"github.com/go-pg/pg/orm"
)

// Base contains common fields for all tables
type Base struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
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
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(_ orm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

// Delete sets deleted_at time to current_time
func (b *Base) Delete() {
	t := time.Now()
	b.DeletedAt = &t
}
