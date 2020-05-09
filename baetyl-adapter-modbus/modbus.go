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
	ctx    context.Context
	wg     *sync.WaitGroup
	cfg    Config
	ws     []*Worker
	logger *log.Logger
	mqtt   *mqtt.Client
	slaves map[byte]*Slave
}

func NewModbus(ctx context.Context, cfg Config) (*Modbus, error) {
	mqttCfg := ctx.ServiceConfig().MQTT
	if mqttCfg.MaxCacheMessages < len(cfg.Jobs)*2 {
		mqttCfg.MaxCacheMessages = len(cfg.Jobs) * 2
	}
	if mqttCfg.ClientID == "" {
		mqttCfg.ClientID = ctx.ServiceName()
	}
	option, err := mqttCfg.ToClientOptions(nil)
	if err != nil {
		return nil, err
	}
	mqtt := mqtt.NewClient(*option)
	logger := ctx.Log()

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
	mod := &Modbus{
		wg:     new(sync.WaitGroup),
		ctx:    ctx,
		cfg:    cfg,
		ws:     ws,
		mqtt:   mqtt,
		logger: logger,
		slaves: slaves,
	}
	for _, worker := range mod.ws {
		mod.wg.Add(1)
		go mod.working(worker)
	}
	return mod, nil
}

func (mod *Modbus) working(w *Worker) {
	ticker := time.NewTicker(w.job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := w.Execute(mod.cfg.Publish)
			if err != nil {
				mod.logger.Error("failed to execute job", log.Error(err))
			}
		case <-mod.ctx.WaitChan():
			mod.logger.Warn("worker stopped", log.Any("worker", w))
			mod.wg.Done()
			return
		}
	}
}

func (mod *Modbus) Close() error {
	mod.wg.Wait()
	var msgs []string
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			msgs = append(msgs, err.Error())
		}
	}
	if len(msgs) != 0 {
		return fmt.Errorf("failed to close slaves: %s", strings.Join(msgs, ";"))
	}
	return mod.mqtt.Close()
}
