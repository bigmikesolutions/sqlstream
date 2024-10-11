package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"sqlstream/test/db"
)

const defaultTimeout = 5 * time.Second

var postgresContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pg, pgCancel, err := db.StartPostgres(ctx)
	if err != nil {
		panic(err)
	}

	defer pgCancel()
	postgresContainer = pg

	code := m.Run()
	os.Exit(code)
}
