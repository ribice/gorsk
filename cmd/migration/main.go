package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-pg/pg/orm"

	"github.com/ribice/gorsk/internal"

	"github.com/go-pg/pg"
	"github.com/ribice/gorsk/internal/auth"
)

func main() {
	dbInsert := `INSERT INTO public.companies VALUES (1, now(), now(), NULL, 'admin_company', true);
	INSERT INTO public.locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);
	INSERT INTO public.roles VALUES (1, 1, 'SUPER_ADMIN');
	INSERT INTO public.roles VALUES (2, 2, 'ADMIN');
	INSERT INTO public.roles VALUES (3, 3, 'COMPANY_ADMIN');
	INSERT INTO public.roles VALUES (4, 4, 'LOCATION_ADMIN');
	INSERT INTO public.roles VALUES (5, 5, 'USER');`
	var psn = ``
	queries := strings.Split(dbInsert, ";")

	u, err := pg.ParseURL(psn)
	checkErr(err)
	db := pg.Connect(u)
	_, err = db.Exec("SELECT 1")
	checkErr(err)
	createSchema(db, &model.Company{}, &model.Location{}, &model.Role{}, &model.User{})

	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		checkErr(err)
	}
	userInsert := `INSERT INTO public.users VALUES (1, now(),now(), NULL, 'Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', NULL, NULL, NULL, NULL, true, 1, 1, 1);`
	_, err = db.Exec(fmt.Sprintf(userInsert, auth.HashPassword("admin")))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		}))
	}
}
