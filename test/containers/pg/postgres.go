// Package pg wraps-up postgres.
package pg

import (
	"context"
	"fmt"
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
	dbPassword     = "postgres"
	dbUser         = "admin"
	driverName     = "pgx"
	dbPort         = "5432/tcp"
	dbInitScript   = "containers/db/init.sh"
)

// CancelFn cancel function for stopping/clearing DB related stuff.
type CancelFn = func()

// Start starts PG in a docker.
func Start(ctx context.Context) (*postgres.PostgresContainer, CancelFn, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithSQLDriver(driverName),
		postgres.WithInitScripts(dbInitScript),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort(dbPort),
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

// Connect creates connection to PG within docker container using default URL.
func Connect(ctx context.Context, c *postgres.PostgresContainer) (*sqlx.DB, error) {
	cs, err := c.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return sqlx.ConnectContext(ctx, driverName, cs)
}

// ConnectByURL creates connection to PG within docker container using given connection string.
func ConnectByURL(ctx context.Context, connectionString string) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, driverName, connectionString)
}

// ConnectionString created connection string for given container and specified port.
func ConnectionString(ctx context.Context, c *postgres.PostgresContainer, port int) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dbUser, dbPassword,
		host, port,
		dbName,
	), nil
}

// ContainerURL create URL for given cotainer.
func ContainerURL(ctx context.Context, c *postgres.PostgresContainer) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := c.MappedPort(ctx, dbPort)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}
