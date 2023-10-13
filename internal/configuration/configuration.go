package configuration

import (
	"gopkg.in/yaml.v3"
	"load-balancer/internal/dispatcher"
	"os"
)

type Global struct {
	Dispatcher *dispatcher.Config `yaml:"dispatcher"`
}

func Load(path string) (*Global, error) {
	var config = &Global{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (g *Global) Start() {
	g.Dispatcher.Start()
}
