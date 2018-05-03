package pgsql_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/platform/postgres"

	"github.com/ribice/gorsk/cmd/api/config"

	"github.com/fortytw2/dockertest"
	"github.com/go-pg/pg"
)

func TestNew(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	_, err = pgsql.New(&config.Database{PSN: "PSN"})
	if err == nil {
		t.Error("Expected error")
	}

	_, err = pgsql.New(&config.Database{PSN: "postgres://postgres:postgres@localhost:1234/postgres?sslmode=disable"})
	if err == nil {
		t.Error("Expected error")
	}

	dbLogTest, err := pgsql.New(&config.Database{PSN: "postgres://postgres:postgres@" + container.Addr + "/postgres?sslmode=disable", Log: true})
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}
	dbLogTest.Close()

	dbCfg := &config.Database{PSN: "postgres://postgres:postgres@" + container.Addr + "/postgres?sslmode=disable", CreateSchema: true}

	db, err := pgsql.New(dbCfg)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	defer db.Close()

	cases := []struct {
		name string
		fn   func(t *testing.T, db *pg.DB, log echo.Logger)
	}{
		{
			name: "AccountDB",
			fn:   testAccountDB,
		},
		{
			name: "UserDB",
			fn:   testUserDB,
		},
	}

	seedData(t, db)

	e := echo.New()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(t, db, e.Logger)
		})
	}

}

func seedData(t *testing.T, db *pg.DB) {

	dbInsert := `INSERT INTO companies VALUES (1, now(), now(), NULL, 'admin_company', true);
INSERT INTO locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);
INSERT INTO roles VALUES (1, 1, 'SUPER_ADMIN');
INSERT INTO roles VALUES (2, 2, 'ADMIN');
INSERT INTO roles VALUES (3, 3, 'COMPANY_ADMIN');
INSERT INTO roles VALUES (4, 4, 'LOCATION_ADMIN');
INSERT INTO roles VALUES (5, 5, 'USER');
INSERT INTO users VALUES (1, now(),now(), NULL, 'John', 'Doe', 'johndoe', 'hunter2', 'johndoe@mail.com', NULL, NULL, NULL, NULL, NULL, 'loginrefresh',1, 1, 1);`

	queries := strings.Split(dbInsert, ";")
	for _, v := range queries[0 : len(queries)-1] {
		_, err := db.Exec(v)
		if err != nil {
			t.Fatalf("Fail on seeding data: %v", err)
		}

	}
}

func queryUser(t *testing.T, db *pg.DB, id int) *model.User {
	user := &model.User{
		Base: model.Base{
			ID: id,
		},
	}
	if err := db.Select(user); err != nil {
		t.Errorf("Could not get user with ID %d due to error %v", id, err)
	}
	return user
}
