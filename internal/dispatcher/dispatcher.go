package dispatcher

import (
	"golang.org/x/net/http2"
	"load-balancer/internal/connection_pool"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Config dispatcher configurations
type Config struct {
	// TargetServers map[key] = value ; key=trigger, value=target
	//
	//	trigger MUST be a path or a subdomain.
	//	target MUST be a valid URL
	//	Examples :
	//
	//		"/search" : "https://google.com/search"
	//
	//	In this example, our dispatcher will transform the current request to call https://google.com/search
	//	switch the host for google.com, call the path search and append any QueryParameter of the original request to the new one.
	//	----------------------------------------------
	//
	//		"hello.mydomain.com" : "https://localhost:9001/"
	//		"hello2.mydomain.com" : "https://10.35.20.8:9001/"
	//
	//	In this example, when the request host equals the trigger, the traffic will be forwarded to the target server.
	//	Can be used if you want a subdomains to call a local service that's not exposed to the internet and running on the same
	//	machine/container/network.
	//	----------------------------------------------
	//
	//		"hello.mydomain.com" : "$custom_pool_1"
	//		"/my-path" 			 : "$custom_pool_1"
	//
	//	In this example, when the trigger is found, the traffic will be forwarded to the target pool.
	//  The method GetNextConnection() will be called and depending on the chosen
	// connection_pool.DistributionType the next server will be chosen.
	TargetServers map[string]string `yaml:"target_servers"`

	// AutoDiscover scan local network to find other dispatchers
	AutoDiscover bool `yaml:"auto_discover"`

	// ListenerPort port to bind this instance to (everything under 1000 might require root)
	ListenerPort string `yaml:"listener_port"`

	// CertPath path to public Key or certificate
	CertPath string `yaml:"cert_path"`

	// KeyPath path to private key
	KeyPath string `yaml:"key_path"`

	// ServerPools map[key] = value ; key = pool_id, value = connection_pool.Config
	ServerPools map[string]*connection_pool.Config `yaml:"server_pools"`
}

func (c *Config) Start() {
	serveMux, server := c.createServer()
	c.setupServeMux(serveMux)
	var err error

	if c.KeyPath != "" && c.CertPath != "" {
		err = server.ListenAndServeTLS(c.CertPath, c.KeyPath)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		log.Printf("Could not start dispatcher server : %s\n", err)
	}
}

func (c *Config) setupServeMux(serveMux *http.ServeMux) {
	var subDomains = make(map[string]*connection_pool.Config)

	for trigger, target := range c.TargetServers {
		if trigger == "/" {
			log.Fatal("You cannot have a trigger that's equal to \"/\". This path is reserved for subdomains")
		}

		isPath := strings.Index(trigger, "/") == 0
		serverPool := c.ServerPools[target[1:]]

		if serverPool == nil && strings.HasPrefix(target, "http") {
			serverPool = connection_pool.NewOneConnectionPool(target)
		} else if serverPool != nil {
			// setup connection pool if applicable - might set up the same instance more than once
			// but the setup isn't very costly to call many times
			serverPool.SetupPoolConfig()
		} else {
			log.Fatal("Invalid configuration : ", target, " does not exist.")
		}

		addHandle(serveMux, &subDomains, trigger, target, isPath, serverPool)
	}

	registerSubdomains(&subDomains, serveMux)
}

func registerSubdomains(subdomains *map[string]*connection_pool.Config, mux *http.ServeMux) {
	if len(*subdomains) > 0 {
		mux.Handle("/", &subdomainMuxServe{
			muxServe{
				subdomains: *subdomains,
			},
		})
	}
}

func addHandle(
	serveMux *http.ServeMux,
	subdomains *map[string]*connection_pool.Config,
	trigger string, target string, isPath bool,
	pool *connection_pool.Config,
) {
	if isPath {
		serveMux.Handle(trigger, &muxServe{
			triggerValue:             trigger,
			target:                   target,
			isPath:                   isPath,
			pathConnectionPoolConfig: pool,
		})
	} else {
		(*subdomains)[trigger] = pool
	}
}

func (c *Config) createServer() (*http.ServeMux, *http.Server) {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    "0.0.0.0:" + c.ListenerPort,
		Handler: serveMux,
	}

	// Configure the server to use HTTP/2
	err := http2.ConfigureServer(server, nil)
	if err != nil {
		log.Fatal("Could not Configure server : ", err)
	}

	return serveMux, server
}

func getTargetURLAndModifyRequest(target string) *url.URL {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Println("Could not parse target server : ", err)
		return nil
	}

	return targetURL
}

type cutPathPrefix struct {
	cut   bool
	value string
}

func singleHostProxyServe(rw http.ResponseWriter, r *http.Request, target *connection_pool.Config, cpp cutPathPrefix) {
	currentTarget := target.GetNextServer()
	if currentTarget == connection_pool.FailOverServersDown {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}

	targetURL := getTargetURLAndModifyRequest(currentTarget)
	if targetURL == nil {
		log.Println("Could not prepare request properly - targetURL is nil")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// use default director but also override request.Host with target.Host
	proxy.Director = func(request *http.Request) {
		prepareRequest(request, targetURL)
		if cpp.cut {
			requestPath := request.URL.Path
			newPath, _ := strings.CutPrefix(requestPath, cpp.value)
			request.URL.Path = newPath
		}

	}

	proxy.ServeHTTP(rw, r)
	log.Println("Served : ", targetURL.String())
}

func prepareRequest(request *http.Request, target *url.URL) {
	// Handle Path
	if request.URL.Path != "/" || target.Path != "/" {
		request.URL.Path = strings.ReplaceAll(request.URL.Path+"/"+target.Path, "//", "/")
	}
	if strings.HasSuffix(request.URL.Path, "/") {
		request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
	}

	// Set Scheme, Host and Path
	request.URL.Scheme = target.Scheme
	request.URL.Host = target.Host
	request.Host = target.Host

	// Handle Query
	request.URL.RawQuery = strings.Trim(strings.Join([]string{request.URL.RawQuery, target.RawQuery}, "&"), "&")
}
