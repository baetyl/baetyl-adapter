package main

import (
	"fmt"
	"github.com/baetyl/baetyl-go/context"
	"strings"
	"sync"
	"time"

	"github.com/baetyl/baetyl-go/log"
	"github.com/baetyl/baetyl-go/mqtt"
)

type Modbus struct {
	cfg    Config
	ws     []*Worker
	logger *log.Logger
	slaves map[byte]*Slave
}

func NewModbus(cfg Config, mqtt *mqtt.Client, logger *log.Logger) *Modbus {
	// modbus slave connection init
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client := NewClient(slaveConfig)
		err := client.Connect()
		if err != nil {
			logger.Error("failed to connect slave id=%d: %s", log.Any("id", slaveConfig.ID), log.Error(err))
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
	}
	var ws []*Worker
	for _, job := range cfg.Jobs {
		w := NewWorker(job, slaves[job.SlaveId], mqtt, logger.With(log.Any("modbus", "map point")))
		ws = append(ws, w)
	}
	return &Modbus{
		cfg:    cfg,
		ws:     ws,
		logger: logger,
		slaves: slaves,
	}
}

func (mod *Modbus) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, worker := range mod.ws {
		wg.Add(1)
		go func(w *Worker, wg *sync.WaitGroup) {
			defer wg.Done()
			ticker := time.NewTicker(w.job.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err := w.Execute(mod.cfg.Publish)
					if err != nil {
						mod.logger.Error("failed to execute job", log.Error(err))
					}
				case <-ctx.WaitChan():
					return
				}
			}
		}(worker, &wg)
	}
	wg.Wait()
}

func (mod *Modbus) Close() error {
	var msgs []string
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			msgs = append(msgs, err.Error())
		}
	}
	if len(msgs) != 0 {
		return fmt.Errorf("failed to close slaves: %s", strings.Join(msgs, ";"))
	}
	return nil
}
