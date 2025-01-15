package containers

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

const (
	postgresNoSSL      = "sslmode=disable"
	postgresDriverName = "pgx"
)

// cfg represents configuration of docker compose with random ports.
type cfg struct {
	Postgres  PostgresCfg
	ToxiProxy ToxiProxyCfg
	uid       int
}

// PostgresCfg keeps docker-compose config.
type PostgresCfg struct {
	Name string
	User string
	Pass string
	Port int
}

// ToxiProxyCfg keeps docker-compose config.
type ToxiProxyCfg struct {
	Port         int
	PostgresPort int
}

func newCfg() cfg {
	return cfg{
		uid: rand.IntN(1024), // nolint
		Postgres: PostgresCfg{
			Port: randomPort(),
			Name: "sqlstream",
			User: "postgres",
			Pass: "pass",
		},
		ToxiProxy: ToxiProxyCfg{
			Port:         randomPort(),
			PostgresPort: randomPort(),
		},
	}
}

// Env returns env variables.
func (c cfg) Env() map[string]string {
	return map[string]string{
		"uid":        strconv.Itoa(c.uid),
		"pgPort":     strconv.Itoa(c.Postgres.Port),
		"pgName":     c.Postgres.Name,
		"pgUser":     c.Postgres.User,
		"pgPass":     c.Postgres.Pass,
		"toxiPort":   strconv.Itoa(c.ToxiProxy.Port),
		"toxiPgPort": strconv.Itoa(c.ToxiProxy.PostgresPort),
	}
}

// ConnectionString returns connection string for postgres container.
func (c PostgresCfg) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		c.User, c.Pass,
		GetHost(), c.Port,
		c.Name,
		postgresNoSSL,
	)
}
