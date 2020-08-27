package modbus

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/mqtt"
)

var configRecoder = make(map[byte]map[string]MapConfig)

type Modbus struct {
	ctx    context.Context
	mqtt   *mqtt.Client
	logger *log.Logger
	slaves map[byte]*Slave
}

func NewModbus(ctx context.Context, cfg Config) (*Modbus, error) {
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client, err := NewClient(slaveConfig)
		if err != nil {
			return nil, err
		}
		err = client.Connect()
		if err != nil {
			ctx.Log().Error("ignore slave device which failed to establish connection", log.Any("id", slaveConfig.ID), log.Error(err))
			continue
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
		configRecoder[slaveConfig.ID] = make(map[string]MapConfig)
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
			if job.Publish.Topic == "" {
				job.Publish.Topic = fmt.Sprintf("%s/%d", ctx.ServiceName(), job.SlaveID)
			}
			sender := NewMqttSender(job.Publish, mqtt)
			w := NewWorker(job, slave, sender, log.With(log.Any("slaveid", job.SlaveID)))
			ws = append(ws, w)
			if job.Encoding != JsonEncoding {
				continue
			}
			if _, ok := configRecoder[job.SlaveID]; !ok {
				continue
			}
			for _, m := range job.Maps {
				if _, ok := configRecoder[job.SlaveID][m.Field.Name]; ok {
					return nil, errors.Errorf("one device should not have same variable definition")
				}
				configRecoder[job.SlaveID][m.Field.Name] = m
			}
		} else {
			ctx.Log().Error("slave of job not exist")
		}
	}
	for _, worker := range ws {
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
			err := w.Execute()
			if err != nil {
				mod.logger.Error("failed to execute job", log.Error(err))
			}
		case <-mod.ctx.WaitChan():
			mod.logger.Warn("worker stopped", log.Any("worker", w))
			return
		}
	}
}

func (mod *Modbus) Close() error {
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			mod.logger.Warn("failed to close slave", log.Any("slave id", slave.cfg.ID))
		}
	}
	mod.mqtt.Close()
	return nil
}
