// Package proxy wraps-up toxi-proxy.
package proxy

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"

	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	port           = "8474/tcp"
	proxyPortRange = "9000-10000"
	startProxyPort = 9000
)

type (
	// Container holds proxy state.
	Container struct {
		container testcontainers.Container
		client    *toxiproxy.Client
		proxyPort int
	}
)

// Start starts proxy in a docker.
func Start(ctx context.Context) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "ghcr.io/shopify/toxiproxy",
		ExposedPorts: []string{
			port,
			proxyPortRange,
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
		proxyPort: startProxyPort,
	}, nil
}

// Close shut-down current instance.
func (c *Container) Close(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

// Upstream create new proxy.
func (c *Container) Upstream(url string) (*Upstream, error) {
	println(url)
	proxy, err := c.client.CreateProxy(
		fmt.Sprintf("proxy_%d", c.proxyPort),
		fmt.Sprintf("0.0.0.0:%d", c.proxyPort),
		url,
	)

	port, portErr := c.container.MappedPort(context.Background(), nat.Port(fmt.Sprintf("%d/tcp", c.proxyPort)))
	if portErr != nil {
		return nil, fmt.Errorf("proxy port: %w", portErr)
	}

	upstream := &Upstream{
		proxy: proxy,
		port:  port.Int(),
	}

	if err == nil {
		c.proxyPort++
	}

	return upstream, err
}

func newClient(ctx context.Context, c testcontainers.Container) (*toxiproxy.Client, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := c.MappedPort(ctx, port)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s:%s", host, port.Port())
	return toxiproxy.NewClient(apiURL), nil
}
