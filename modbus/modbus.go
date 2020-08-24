package modbus

import (
	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/mqtt"
	"sync"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
)

type Modbus struct {
	ctx    context.Context
	wg     sync.WaitGroup
	mqtt   *mqtt.Client
	logger *log.Logger
	slaves map[byte]*Slave
}

func NewModbus(ctx context.Context, cfg Config) (*Modbus, error) {
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client := NewClient(slaveConfig)
		err := client.Connect()
		if err != nil {
			ctx.Log().Error("failed to connect slave", log.Any("id", slaveConfig.ID), log.Error(err))
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
	}
	mqttCfg := ctx.SystemConfig().Broker
	if mqttCfg.MaxCacheMessages < len(cfg.Jobs)*2 {
		mqttCfg.MaxCacheMessages = len(cfg.Jobs) * 2
	}
	if mqttCfg.ClientID == "" {
		mqttCfg.ClientID = ctx.ServiceName()
	}
	option, err := mqttCfg.ToClientOptions()
	if err != nil {
		return nil, err
	}
	mqtt := mqtt.NewClient(option)
	observer := NewObserver(slaves, ctx.Log())
	if err := mqtt.Start(observer); err != nil {
		return nil, err
	}
	mod := &Modbus{
		ctx:    ctx,
		mqtt:   mqtt,
		logger: ctx.Log(),
		slaves: slaves,
	}
	var ws []*Worker
	for _, job := range cfg.Jobs {
		if slave := slaves[job.SlaveID]; slave != nil {
			sender := NewMqttSender(job.Publish, mqtt)
			w := NewWorker(job, slave, sender, log.With(log.Any("slaveid", job.SlaveID)))
			ws = append(ws, w)
		} else {
			ctx.Log().Error("slave of job not exist")
		}
	}
	for _, worker := range ws {
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
			// TODO add independent topic of each job, supply default topic like <service_name>/<slave_id>
			err := w.Execute()
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
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			mod.logger.Warn("failed to close slave", log.Any("slave id", slave.cfg.ID))
		}
	}
	mod.mqtt.Close()
	return nil
}
