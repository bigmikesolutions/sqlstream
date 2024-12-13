// Package toxiproxy wraps-up toxi-proxy.
package toxiproxy

import (
	"context"
	"fmt"

	client "github.com/Shopify/toxiproxy/v2/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	port      = "8474"
	portRange = "5000-6000"
	startPort = 5000
)

type (
	// Container holds toxi-proxy state.
	Container struct {
		container testcontainers.Container
		client    *client.Client
		port      int
	}
)

// Start starts toxi-proxy in a docker.
func Start(ctx context.Context) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "ghcr.io/shopify/toxiproxy",
		ExposedPorts: []string{
			port,
			portRange,
		},
		WaitingFor: wait.ForHTTP("/version").WithPort(port),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("toxi-proxy start: %w", err)
	}

	c, err := newClient(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("toxi-proxy client: %w", err)
	}

	return &Container{
		container: container,
		client:    c,
		port:      startPort,
	}, nil
}

// Close shut-down current instance.
func (c *Container) Close(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

// NewProxy create new proxy.
func (c *Container) NewProxy(url string) (*client.Proxy, int, error) {
	proxy, err := c.client.CreateProxy(
		fmt.Sprintf("proxy_%d", c.port),
		fmt.Sprintf("0.0.0.0:%d", c.port),
		url,
	)

	port := c.port
	if err == nil {
		c.port++
	}

	return proxy, port, err
}

func newClient(ctx context.Context, c testcontainers.Container) (*client.Client, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := c.MappedPort(ctx, port)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s:%s", host, port.Port())
	return client.NewClient(apiURL), nil
}
