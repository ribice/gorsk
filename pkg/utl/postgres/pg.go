package postgres

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
)

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
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			if query, err := event.FormattedQuery(); err == nil {
				log.Printf("%s | %s", time.Since(event.StartTime), query)
			}
		})
	}

	return db, nil
}
