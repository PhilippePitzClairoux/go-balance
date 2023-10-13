package main

import (
	"flag"
	"load-balancer/internal/configuration"
	"log"
)

var (
	configFile = flag.String("config-file", "example-config.yaml", "Configuration file location")
)

func main() {
	flag.Parse()

	conf, err := configuration.Load(*configFile)
	if err != nil {
		log.Println("There was an error loading global configuration file : ", err)
	}

	conf.Start()
}
