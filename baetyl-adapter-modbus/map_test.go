package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	server := MbSlave{}
	server.StartTCPSlave()

	slaveItem := SlaveItem{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client := NewClient(slaveItem)
	client.Connect()
	defer client.Close()
	slave := NewSlave(slaveItem, client)

	item1 := MapItem{
		SlaveID:  1,
		Function: 2,
		Address:  0,
		Quantity: 1,
	}
	m := NewMap(item1, slave)

	results, err := m.Read()
	assert.NoError(t, err)
	expected1 := []byte{1}
	assert.Equal(t, results, expected1)

	item2 := MapItem{
		SlaveID:  1,
		Function: 4,
		Address:  0,
		Quantity: 2,
	}
	m = NewMap(item2, slave)

	results, err = m.Read()
	assert.NoError(t, err)
	expected2 := []byte{144, 144}
	assert.Equal(t, results, expected2)
	server.Stop()

	client.WriteSingleCoil(0, 0xFF00)
	item3 := MapItem{
		SlaveID:  1,
		Function: 1,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(item3, slave)
	results, err = m.Read()
	assert.NoError(t, err)
	expected3 := []byte{1}
	assert.Equal(t, results, expected3)

	client.WriteSingleRegister(0, 65535)
	item4 := MapItem{
		SlaveID:  1,
		Function: 3,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(item4, slave)
	results, err = m.Read()
	assert.NoError(t, err)
	expected4 := []byte{255, 255}
	assert.Equal(t, results, expected4)

	item5 := MapItem{
		SlaveID:  1,
		Function: 3,
		Address:  0,
		Quantity: 0,
	}
	m = NewMap(item5, slave)
	results, err = m.Read()
	assert.Error(t, err, "modbus: quantity '0' must be between '1' and '125',")
}
