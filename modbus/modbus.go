package modbus

import (
	"github.com/baetyl/baetyl-go/context"
	"sync"
	"time"

	"github.com/baetyl/baetyl-go/log"
)

type Modbus struct {
	ctx    context.Context
	wg     sync.WaitGroup
	sender Sender
	logger *log.Logger
	slaves map[byte]*Slave
}

func NewModbus(ctx context.Context, cfg Config, sender Sender) (*Modbus, error) {
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client := NewClient(slaveConfig)
		err := client.Connect()
		if err != nil {
			ctx.Log().Error("failed to connect slave", log.Any("id", slaveConfig.ID), log.Error(err))
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
	}
	mod := &Modbus{
		ctx:    ctx,
		sender: sender,
		logger: ctx.Log(),
		slaves: slaves,
	}
	var ws []*Worker
	for _, job := range cfg.Jobs {
		if slave := slaves[job.SlaveId]; slave != nil {
			w := NewWorker(job, slave, sender, log.With(log.Any("slaveid", job.SlaveId)))
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
	return nil
}
