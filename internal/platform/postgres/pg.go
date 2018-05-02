package pgsql

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/cmd/api/config"
)

const notDeleted = "deleted_at is null"

// New creates new database connection to a postgres database
// Function panics if it can't connect to database
func New(cfg *config.Database) (*pg.DB, error) {
	u, err := pg.ParseURL(cfg.PSN)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(u).WithTimeout(time.Second * 5)
	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}
	if cfg.Log {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			checkErr(err)
			log.Printf("%s | %s", time.Since(event.StartTime), query)
		})
	}
	if cfg.CreateSchema {
		createSchema(db, &model.Company{}, &model.Location{}, &model.Role{}, &model.User{})
	}
	return db, nil
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, nil))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
