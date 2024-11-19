package test

import (
	"context"
	"os"
	"testing"
	"time"

	"sqlstream/test/db"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	defaultTimeout     = 5 * time.Second
	dockerSetupTimeout = 30 * time.Second
)

var postgresContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), dockerSetupTimeout)
	defer cancel()

	pg, pgCancel, err := db.StartPostgres(ctx)
	if err != nil {
		panic(err)
	}

	defer pgCancel()
	postgresContainer = pg

	code := m.Run()
	defer os.Exit(code)
}
