// Package containers wraps-up containers set-up.
package containers

import (
	"context"
	"fmt"
	"sync"

	"sqlstream/test/containers/pg"
	"sqlstream/test/containers/toxiproxy"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Service holds state of started containers.
type Service struct {
	pg        *postgres.PostgresContainer
	pgCancel  pg.CancelFn
	toxiproxy *toxiproxy.Container
}

// New start containers.
func New(ctx context.Context) (*Service, error) {
	var wg sync.WaitGroup
	errs := make([]error, 0)
	s := &Service{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		if s.pg, s.pgCancel, err = pg.Start(ctx); err != nil {
			errs = append(errs, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		if s.toxiproxy, err = toxiproxy.Start(ctx); err != nil {
			errs = append(errs, err)
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		s.Close()
		return nil, fmt.Errorf("failed to start containers: %w", joinErr(errs))
	}

	return s, nil
}

// NewDB creates direct connect to database.
func (s *Service) NewDB(ctx context.Context) (*sqlx.DB, error) {
	return pg.Connect(ctx, s.pg)
}

// NewDBProxy creates connection to database via proxy.
func (s *Service) NewDBProxy(ctx context.Context) (*DBProxy, error) {
	directURL, err := pg.ContainerURL(ctx, s.pg)
	if err != nil {
		return nil, fmt.Errorf("postgres URL: %w", err)
	}

	proxy, port, err := s.toxiproxy.NewProxy(directURL)
	if err != nil {
		return nil, err
	}

	proxyConnStr, err := pg.ConnectionString(ctx, s.pg, port)
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}

	conn, err := pg.ConnectByURL(ctx, proxyConnStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connection: %w", err)
	}

	return &DBProxy{
		DB:    conn,
		Proxy: proxy,
	}, nil
}

// Close closes all created containers gracefully.
func (s *Service) Close() {
	if s.pgCancel != nil {
		s.pgCancel()
	}
	if s.toxiproxy != nil {
		_ = s.toxiproxy.Close(context.Background())
	}
}
