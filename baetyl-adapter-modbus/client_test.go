package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	slave := MbSlave{}
	slave.StartTCPSlave()

	cfg := SlaveItem{}
	client := NewClient(cfg)
	err := client.Connect()
	defer client.Close()
	assert.Error(t, err, "failed to connect: no such file or directory")

	cfg = SlaveItem{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client = NewClient(cfg)
	err = client.Connect()
	defer client.Close()
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
	slave.Stop()

	cfg = SlaveItem{
		ID:      2,
		Address: "tcp://127.0.0.1:50201",
	}
	client = NewClient(cfg)
	err = client.Connect()
	defer client.Close()
	assert.Error(t, err, "failed to connect: dial tcp 127.0.0.1:50201: connect: connection refused")
}
