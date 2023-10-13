package connection_pool

import (
	"math"
	"net/http"
	"testing"
)

func TestServerPool(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked even tho it shouldn't%v", r)
		}
	}()

	conf := Config{
		Type: "ROUND_ROBIN",
		Connections: []string{
			"test1",
			"test2",
			"test3",
		},
		lastServerIndexUsed: math.MaxUint64 - 3,
	}
	expectedOutput := []string{
		"test1",
		"test2",
		"test3",
	}

	conf.SetupPoolConfig()
	for i, j := conf.lastServerIndexUsed, 0; i != 3; i++ {
		if conf.GetNextServer() != expectedOutput[j] {
			t.Log("GetNextServer did not match expected output")
			t.FailNow()
		}
		j++
	}

	t.Log("All servers matched!")
}

func TestFailOverPool(t *testing.T) {
	conf := Config{
		Type: "FAIL_OVER",
		Connections: []string{
			"test1",
			"test2",
			"http://localhost:8080",
		},
		TestConnection: &TestConnection{
			Path:               "/",
			Method:             "GET",
			BodyMustContain:    "",
			ExpectedStatusCode: 404,
			retry:              0,
		},
		lastServerIndexUsed: 0,
	}

	const expected = "http://localhost:8080"
	conf.SetupPoolConfig()
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()

	server := conf.GetNextServer()
	if server != expected {
		t.Logf("FailOver did not work properly - expected %s, got %s", expected, server)
		t.FailNow()
	}

	t.Log("FailOver worked!")
}
