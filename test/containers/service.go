// Package containers wraps-up containers set-up.
package containers

import (
	"context"
	"fmt"
	"sync"

	"github.com/bigmikesolutions/sqlstream/test/containers/pg"
	"github.com/bigmikesolutions/sqlstream/test/containers/proxy"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// Service holds state of started containers.
type Service struct {
	pg       *postgres.PostgresContainer
	pgCancel pg.CancelFn
	proxy    *proxy.Container
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
		if s.proxy, err = proxy.Start(ctx); err != nil {
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

// DB creates direct connect to database.
func (s *Service) DB(ctx context.Context) (*sqlx.DB, error) {
	return pg.Connect(ctx, s.pg)
}

// DBProxy creates connection to database via proxy.
func (s *Service) DBProxy(ctx context.Context) (*DBProxy, error) {
	internalURL, err := pg.InternalURL(ctx, s.pg)
	if err != nil {
		return nil, fmt.Errorf("postgres URL: %w", err)
	}

	upstream, err := s.proxy.Upstream(internalURL)
	if err != nil {
		return nil, fmt.Errorf("proxy for postgres: %w", err)
	}

	proxyConnStr, err := pg.ConnectionString(ctx, s.pg, upstream.Port())
	if err != nil {
		return nil, fmt.Errorf("postgres proxy connection string: %w", err)
	}

	conn, err := pg.ConnectByURL(ctx, proxyConnStr)
	if err != nil {
		return nil, fmt.Errorf("postgres proxy connection: %w", err)
	}

	return &DBProxy{
		DB:       conn,
		Upstream: upstream,
	}, nil
}

// Close closes all created containers gracefully.
func (s *Service) Close() {
	if s.pgCancel != nil {
		s.pgCancel()
	}
	if s.proxy != nil {
		_ = s.proxy.Close(context.Background())
	}
}
