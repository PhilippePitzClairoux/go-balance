package connection_pool

import "log"

// DistributionType MUST contain one of the following values :
// ONE_CONNECTION_SERVER_POOL
// ROUND_ROBIN
// FAIL_OVER
type DistributionType string

type ServerGeneratorFunc func(config *Config) string

type Config struct {
	Type                DistributionType `yaml:"distribution_type"`
	Connections         []string         `yaml:"connections"`
	TestConnection      *TestConnection  `yaml:"test_connection"`
	generator           ServerGeneratorFunc
	lastServerIndexUsed uint64
}

type TestConnection struct {
	Path               string `yaml:"path"`
	Method             string `yaml:"method"`
	ExpectedStatusCode int    `yaml:"expected_status_code"`
	BodyMustContain    string `yaml:"body_must_contain"`
	retry              int
}

func (spc *Config) GetNextServer() string {
	return spc.generator(spc)
}

func NewOneConnectionPool(server string) *Config {
	return &Config{
		Type:        "ONE_CONNECTION_SERVER_POOL",
		generator:   OneConnectionPool,
		Connections: []string{server},
	}
}

func (spc *Config) SetupPoolConfig() {
	spc.lastServerIndexUsed = 0

	switch spc.Type {
	case "ROUND_ROBIN":
		spc.generator = RoundRobinNextServer
	case "ONE_CONNECTION_SERVER_POOL":
		spc.generator = OneConnectionPool
	case "FAIL_OVER":
		if spc.TestConnection == nil {
			log.Fatal("FAIL_OVER was chosen as distribution_type but doesn't define a test_connection")
		}
		spc.generator = FailOverPool
	default:
		log.Fatal("Invalid DistributionType : ", spc.Type)
	}
}
