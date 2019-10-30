package main

import (
	"github.com/baetyl/baetyl/logger"
	"github.com/baetyl/baetyl/protocol/mqtt"
	"time"
)

type Modbus struct {
	cfg     Config
	mbs     map[byte]*mb
	mqttCli *mqtt.Dispatcher
	log     logger.Logger
	slaves  map[byte]*Slave
}

func newModbus(cfg Config, mqttCli *mqtt.Dispatcher, log logger.Logger) *Modbus {
	// modbus slave.go connection init
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client := NewClient(slaveConfig)
		err := client.Connect()
		if err != nil {
			log.Errorf("failed to connect slave id=%d: %s", slaveConfig.ID, err.Error())
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
	}
	mbs := make(map[byte]*mb, 0)
	for _, mapConfig := range cfg.Maps {
		m := NewMap(mapConfig, slaves[mapConfig.SlaveID], logger.WithField("modbus", "map"))
		if _, ok := mbs[mapConfig.SlaveID]; !ok {
			mbs[mapConfig.SlaveID] = NewMb(slaves[mapConfig.SlaveID].cfg.Interval)
		}
		mbs[mapConfig.SlaveID].AddMap(m)
	}
	return &Modbus{
		cfg:     cfg,
		mbs:     mbs,
		mqttCli: mqttCli,
		log:     log,
		slaves:  slaves,
	}
}

func (mod *Modbus) Start() {
	for _, m := range mod.mbs {
		go func(m *mb) {
			ticker := time.NewTicker(m.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := m.Execute(mod.mqttCli, mod.cfg.Publish)
					if err != nil {
						mod.log.Errorf(err.Error())
					}
				}
			}
		}(m)
	}
}

func (mod *Modbus) Close() {
	for _, slave := range mod.slaves {
		slave.client.Close()
	}
	mod.mqttCli.Close()
}
