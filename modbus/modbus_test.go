package main

import (
	"testing"
)

func TestModbus(t *testing.T) {
//	confString := `
//slaves:
//  - id: 1
//    address: tcp://127.0.0.1:50200
//    interval: 3s
//maps:
//  - slaveid: 1
//    address: 0
//    quantity: 1
//    function: 3
//publish:
// topic: test
//`
//	slave := MbSlave{}
//	slave.StartTCPSlave()
//	var cfg Config
//	utils.UnmarshalYAML([]byte(confString), &cfg)
//	mqttCfg := mqtt.ClientConfig{
//		Address:  "tcp://127.0.0.1:1883",
//		ClientID: "modbus-1",
//	}
//	log := log.With(log.Any("modbus", "test"))
//	mqttOption, err := mqttCfg.ToClientOptions(nil)
//	assert.NoError(t, err)
//	mqtt := mqtt.NewClient(*mqttOption)
//	modbus := NewModbus(cfg, mqtt, log)
//	assert.NotNil(t, modbus)
//	modbus.Close()
//	slave.Stop()
}
