package main

import (
	"sync"
	"time"

	"github.com/baetyl/baetyl/sdk/baetyl-go"

	"github.com/baetyl/baetyl/logger"
	"github.com/baetyl/baetyl/protocol/mqtt"
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
		m := NewMap(mapConfig, slaves[mapConfig.SlaveID], logger.WithField("modbus", "map point"))
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

func (mod *Modbus) Start(ctx baetyl.Context) {
	var wg sync.WaitGroup
	for _, m := range mod.mbs {
		wg.Add(1)
		go func(m *mb, wg *sync.WaitGroup) {
			defer wg.Done()
			ticker := time.NewTicker(m.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := m.Execute(mod.mqttCli, mod.cfg.Publish)
					if err != nil {
						mod.log.Errorf(err.Error())
					}
				case <-ctx.WaitChan():
					return
				}
			}
		}(m, &wg)
	}
	wg.Wait()
}

func (mod *Modbus) Close() error {
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			return err
		}
	}
	return nil
}
