package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk"

	"github.com/ribice/gorsk/pkg/utl/postgres"

	"github.com/fortytw2/dockertest"
)

func TestNew(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:alpine", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	}, "-e", "POSTGRES_PASSWORD=postgres", "-e", "POSTGRES_USER=postgres")
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	_, err = postgres.New("PSN", 1, false)
	if err == nil {
		t.Error("Expected error")
	}

	_, err = postgres.New("postgres://postgres:postgres@localhost:1234/postgres?sslmode=disable", 0, false)
	if err == nil {
		t.Error("Expected error")
	}

	dbLogTest, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 0, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}
	dbLogTest.Close()

	db, err := postgres.New("postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", 1, true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	var user gorsk.User
	db.Select(&user)

	assert.NotNil(t, db)

	db.Close()

}
