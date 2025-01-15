package containers

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// GetHost returns docker host.
func GetHost() string {
	if dockerHost, ok := os.LookupEnv("DOCKER_HOST"); ok {
		dh, err := url.Parse(dockerHost)
		if err == nil {
			panic(fmt.Errorf("DOCKER_HOST url: %w", err))
		}
		return dh.Hostname()
	}

	return "localhost"
}

func randomPort() int {
	listen := randomListener("tcp")
	idx := strings.LastIndex(listen.Addr().String(), ":")

	p := listen.Addr().String()[idx+1:]
	if err := listen.Close(); err != nil {
		panic(fmt.Errorf("error closing random listener: %w", err))
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		panic(fmt.Errorf("error parsing random port: %w", err))
	}
	return port
}

func randomListener(network string) net.Listener {
	listen, err := net.Listen(network, "[::1]:0")
	if err != nil {
		panic(fmt.Errorf("could not listen on random port: %w", err))
	}
	return listen
}
