package main

import (
	"log"
	"os"

	"github.com/go-pg/pg/orm"

	"github.com/ribice/gorsk/cmd/api/config"
	"github.com/ribice/gorsk/cmd/migration/queries"
	"github.com/ribice/gorsk/internal"

	"github.com/go-pg/pg"

	"github.com/joho/godotenv"
)

const (
	dbCheck    = "SELECT 1"
	appEnvName = "APP_CFG_ENVIRONMENT_NAME"
)

var dbq []string
var dbm []interface{}

func init() {
	// Append any new models you create to this list
	dbm = append(dbm, &model.Company{}, &model.Location{}, &model.Role{}, &model.User{})
	// Append any new queries you create to this slice
	dbq = append(dbq, queries.DBSetupQueries()...)
	dbq = append(dbq, queries.CompanyQueries()...)
	dbq = append(dbq, queries.LocationQueries()...)
	dbq = append(dbq, queries.RoleQueries()...)
	dbq = append(dbq, queries.UserQueries()...)
}

func main() {
	err := godotenv.Load()
	checkErr(err)

	env := os.Getenv(appEnvName)
	cfg, err := config.Load(env)
	checkErr(err)

	psn := cfg.DB.PSN

	u, err := pg.ParseURL(psn)
	checkErr(err)

	db := pg.Connect(u)
	_, err = db.Exec(dbCheck)
	checkErr(err)

	log.Println("Check")
	createSchema(db, dbm...)

	for _, v := range dbq {
		_, err := db.Exec(v)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		}))
	}
}
