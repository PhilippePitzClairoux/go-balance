package dispatcher

import (
	"load-balancer/internal/connection_pool"
	"net/http"
	"strings"
)

type muxServe struct {
	triggerValue             string
	target                   string
	subdomains               map[string]*connection_pool.Config
	pathConnectionPoolConfig *connection_pool.Config
	isPath                   bool
}

func (ms *muxServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.Path, ms.triggerValue) > 0 {
		http.Error(w, "The requested path does not exist", http.StatusNotFound)
		return
	}

	singleHostProxyServe(w, r, ms.pathConnectionPoolConfig, cutPathPrefix{cut: true, value: ms.triggerValue})
}
