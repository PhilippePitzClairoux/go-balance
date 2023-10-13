package connection_pool

import (
	"math"
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

	conf.SetupPoolConfig()
	t.Log("lastServerIndexUsed = math.MaxUint64 - 3")

	for i := conf.lastServerIndexUsed; i != 3; i++ {
		t.Log(i, conf.GetNextServer())
	}

	t.Log("lastServerIndexUsed reached the value 3!")
}
