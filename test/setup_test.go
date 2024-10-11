package test

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"sqlstream/test/db"
	"testing"
	"time"
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