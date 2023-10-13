package dispatcher

import (
	"net/http"
	"strings"
)

type subdomainMuxServe struct {
	muxServe
}

func (sms *subdomainMuxServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0] // Extract hostname without port
	target, ok := sms.subdomains[host]
	if !ok {
		http.Error(w, "The requested subdomain does not exist", http.StatusNotFound)
		return
	}

	singleHostProxyServe(w, r, target, cutPathPrefix{cut: false})
}
