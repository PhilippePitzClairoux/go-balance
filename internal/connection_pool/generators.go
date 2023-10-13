package connection_pool

import (
	"io"
	"net/http"
	"strings"
)

const FailOverServersDown = "ALL_SERVERS_DOWN"

// RoundRobinNextServer iterates through the server list
func RoundRobinNextServer(config *Config) string {
	value := config.Connections[config.lastServerIndexUsed%uint64(len(config.Connections))]
	config.lastServerIndexUsed++

	return value
}

// OneConnectionPool always returns the same server
func OneConnectionPool(config *Config) string {
	return config.Connections[0]
}

// FailOverPool pings the server, if it's down goes to the next one until we find one that's up
func FailOverPool(config *Config) string {
	lastIndexUsed := config.lastServerIndexUsed % uint64(len(config.Connections))
	target := config.Connections[lastIndexUsed]

	if config.TestConnection.retry == 3 {
		return FailOverServersDown
	}

	if !config.TestConnection.isServerUp(target) {
		config.TestConnection.retry++
		config.lastServerIndexUsed++
		return FailOverPool(config)
	}

	config.TestConnection.retry = 0
	return config.Connections[lastIndexUsed]
}

func (c *TestConnection) isServerUp(target string) bool {
	request, err := http.NewRequest(c.Method, target+c.Path, nil)
	if err != nil {
		return false
	}

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		return false
	}

	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return false
	}

	return do.StatusCode == c.ExpectedStatusCode && strings.Contains(string(all), c.BodyMustContain)
}
