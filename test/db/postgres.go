package db

import (
	"context"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // drivers
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // drivers
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultTimeout = 5 * time.Second
	dbName         = "postgres"
	dbUser         = "admin"
	driverName     = "pgx"
	dbInitScript   = "db/init.sh"
)

// CancelFn cancel function for stopping/clearing DB related stuff.
type CancelFn = func()

// StartPostgres starts PG in a docker.
func StartPostgres(ctx context.Context) (*postgres.PostgresContainer, CancelFn, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithSQLDriver(driverName),
		postgres.WithInitScripts(dbInitScript),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		),
	)

	cancelFn := func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate postgres: %s", err)
		}
	}

	if err != nil {
		log.Printf("failed to start postgres: %s", err)
		return nil, nil, err
	}

	return postgresContainer, cancelFn, nil
}

// ConnectToPostgres created connection to PG within docker container.
func ConnectToPostgres(ctx context.Context, c *postgres.PostgresContainer) (*sqlx.DB, error) {
	cs, err := c.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return sqlx.Connect(driverName, cs)
}
