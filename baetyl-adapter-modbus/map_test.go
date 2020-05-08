package main

import (
	"testing"

	"github.com/baetyl/baetyl-go/log"
	"github.com/stretchr/testify/assert"
)

func TestMapRead(t *testing.T) {
	server := MbSlave{}
	server.StartTCPSlave()

	slaveConfig := SlaveConfig{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client := NewClient(slaveConfig)
	client.Connect()
	defer client.Close()
	slave := NewSlave(slaveConfig, client)
	log := log.With(log.Any("modbus", "map_test"))

	cfg1 := MapConfig{
		Function: 2,
		Address:  0,
		Quantity: 1,
	}
	m := NewMap(cfg1, slave, log)

	results, err := m.read()
	assert.NoError(t, err)
	expected1 := []byte{1}
	assert.Equal(t, results, expected1)

	cfg2 := MapConfig{
		Function: 4,
		Address:  0,
		Quantity: 2,
	}
	m = NewMap(cfg2, slave, log)

	results, err = m.read()
	assert.NoError(t, err)
	expected2 := []byte{144, 144}
	assert.Equal(t, results, expected2)

	client.WriteSingleCoil(0, 0xFF00)
	cfg3 := MapConfig{
		Function: 1,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(cfg3, slave, log)
	results, err = m.read()
	assert.NoError(t, err)
	expected3 := []byte{1}
	assert.Equal(t, results, expected3)

	client.WriteSingleRegister(0, 65535)
	cfg4 := MapConfig{
		Function: 3,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(cfg4, slave, log)
	results, err = m.read()
	assert.NoError(t, err)
	expected4 := []byte{255, 255}
	assert.Equal(t, results, expected4)

	cfg5 := MapConfig{
		Function: 3,
		Address:  0,
		Quantity: 0,
	}
	m = NewMap(cfg5, slave, log)
	results, err = m.read()
	assert.Error(t, err, "modbus: quantity '0' must be between '1' and '125',")

	server.Stop()
}

func TestMapCollect(t *testing.T) {
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
		Function: 2,
		Address:  0,
		Quantity: 1,
	}

	log := log.With(log.Any("modbus", "test"))
	ma := NewMap(mapConfig, slave, log)

	results, err := ma.Collect()
	assert.NoError(t, err)
	expected := []byte{1}
	// address is 0(uint16) corresponding to []byte{0, 0}
	assert.Equal(t, results[:2], []byte{0, 0})
	// quantity is 1(uint16) corresponding to []byte{0, 1}
	assert.Equal(t, results[2:4], []byte{0, 1})
	// result should be 1 (slave have already set it to 1)
	assert.Equal(t, results[4:], expected)

	// invalid quantity
	mapConfig2 := MapConfig{
		Function: 2,
		Address:  0,
		Quantity: 0,
	}
	ma2 := NewMap(mapConfig2, slave, log)
	_, err = ma2.Collect()
	assert.Error(t, err, "failed to collect data from slave.go: id=1 function=2 address=0 quantity=0")
	server.Stop()
}
