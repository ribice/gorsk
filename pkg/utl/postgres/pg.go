package postgres

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(pg *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	log.Printf(q.FormattedQuery())
}

// New creates new database connection to a postgres database
func New(psn string, timeout int, enableLog bool) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}
