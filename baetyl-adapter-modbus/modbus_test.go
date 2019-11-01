package main

import (
	"testing"

	"github.com/baetyl/baetyl/logger"
	"github.com/baetyl/baetyl/protocol/mqtt"
	"github.com/baetyl/baetyl/utils"
	"github.com/stretchr/testify/assert"
)

func TestModbus(t *testing.T) {
	confString := `
slaves:
  - id: 1
    address: tcp://127.0.0.1:50200
    interval: 3s
maps:
  - slaveid: 1
    address: 0
    quantity: 1
    function: 3
publish:
 topic: test
`
	slave := MbSlave{}
	slave.StartTCPSlave()
	var cfg Config
	utils.UnmarshalYAML([]byte(confString), &cfg)
	hubCfg := mqtt.ClientInfo{
		Address:  "tcp://127.0.0.1:1883",
		ClientID: "modbus1",
	}
	log := logger.WithField("modbus", "test")
	mqttCli := mqtt.NewDispatcher(hubCfg, log)
	modbus := newModbus(cfg, mqttCli, log)
	assert.NotNil(t, modbus)
	modbus.Close()
	slave.Stop()
}
