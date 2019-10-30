package main

import (
	"github.com/baetyl/baetyl/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollect(t *testing.T) {
	server := MbSlave{}
	server.StartTCPSlave()

	cfg := SlaveConfig{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client := NewClient(cfg)
	client.Connect()
	defer client.Close()
	slave := NewSlave(cfg, client)

	mapConfig := MapConfig{
		SlaveID:  1,
		Function: 2,
		Address:  0,
		Quantity: 1,
	}

	log := logger.WithField("modbus", "test")
	ma := NewMap(mapConfig, slave, log)

	results, err := Package(ma)
	assert.NoError(t, err)
	expected := []byte{1}
	// address is 0(uint16) corresponding to []byte{0, 0}
	assert.Equal(t, results[1:3], []byte{0, 0})
	// quantity is 1(uint16) corresponding to []byte{0, 1}
	assert.Equal(t, results[3:5], []byte{0, 1})
	// result should be 1 (slave have already set it to 1)
	assert.Equal(t, results[9:10], expected)

	// invalid quantity
	mapConfig2 := MapConfig{
		SlaveID:  1,
		Function: 2,
		Address:  0,
		Quantity: 0,
	}
	ma2 := NewMap(mapConfig2, slave, log)
	_, err = Package(ma2)
	assert.Error(t, err, "failed to collect data from slave.go: id=1 function=2 address=0 quantity=0")
	server.Stop()
}
