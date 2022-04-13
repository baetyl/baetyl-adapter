package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	v2log "github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

var (
	ErrWorkerNotExist = errors.New("worker not exist")
)

type Modbus struct {
	ctx    dm.Context
	log    *v2log.Logger
	slaves map[byte]*Slave
	ws     map[string]*Worker
}

func NewModbus(ctx dm.Context, cfg Config) (*Modbus, error) {
	devMap := map[string]dm.DeviceInfo{}
	for _, dev := range ctx.GetAllDevices() {
		devMap[dev.Name] = dev
	}
	log := ctx.Log().With(v2log.Any("module", "modbus"))
	slaves := map[byte]*Slave{}
	for _, slaveConfig := range cfg.Slaves {
		client, err := NewClient(slaveConfig)
		if err != nil {
			return nil, err
		}
		err = client.Connect()
		if err != nil {
			log.Warn("connect failed", v2log.Any("slave id", slaveConfig.Id), v2log.Error(err))
		}
		dev, ok := devMap[slaveConfig.Device]
		if !ok {
			log.Error("can not find device according to job config", v2log.Any("device", slaveConfig.Device))
			continue
		}
		slave := NewSlave(ctx, &dev, slaveConfig, client)
		slaves[slaveConfig.Id] = slave
		if err = slave.UpdateStatus(SlaveOnline); err != nil {
			log.Error("failed to update status", v2log.Any("error", err))
		}
	}
	mod := &Modbus{
		ctx:    ctx,
		ws:     make(map[string]*Worker),
		log:    log,
		slaves: slaves,
	}
	for _, job := range cfg.Jobs {
		if slave := slaves[job.SlaveID]; slave != nil {
			dev, ok := devMap[slave.cfg.Device]
			if !ok {
				log.Error("can not find device according to job config", v2log.Any("device", slave.cfg.Device))
				continue
			}
			mod.ws[dev.Name] = NewWorker(ctx, job, slave, log)
		} else {
			log.Warn("slave id of job is invalid", v2log.Any("id", job.SlaveID))
		}
	}
	if err := ctx.RegisterDeltaCallback(mod.DeltaCallback); err != nil {
		return nil, err
	}
	if err := ctx.RegisterEventCallback(mod.EventCallback); err != nil {
		return nil, err
	}
	return mod, nil
}

func (mod *Modbus) DeltaCallback(info *dm.DeviceInfo, prop v1.Delta) error {
	w, ok := mod.ws[info.Name]
	if !ok {
		mod.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
		return ErrWorkerNotExist
	}
	ms := map[string]MapConfig{}
	for _, m := range w.job.Maps {
		ms[m.Name] = m
	}
	for name, val := range prop {
		slave, ok := mod.slaves[w.job.SlaveID]
		if !ok {
			mod.log.Warn("did not find slave to write", v2log.Any("slave id", w.job.SlaveID))
			continue
		}
		cfg, ok := ms[name]
		if !ok {
			mod.log.Warn("did not find prop", v2log.Any("name", name))
			continue
		}
		bs, err := transform(val, cfg)
		if err != nil {
			mod.log.Warn("ignore illegal data type of val", v2log.Any("value", val), v2log.Any("type", cfg.Type), v2log.Error(err))
			continue
		}
		switch cfg.Function {
		case DiscreteInput:
		case InputRegister:
			return fmt.Errorf("can not write data with illegal function code: [%d]", cfg.Function)
		case Coil:
			if _, err := slave.client.WriteMultipleCoils(cfg.Address, cfg.Quantity, bs); err != nil {
				return err
			}
		case HoldingRegister:
			if _, err := slave.client.WriteMultipleRegisters(cfg.Address, cfg.Quantity, bs); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mod *Modbus) EventCallback(info *dm.DeviceInfo, event *dm.Event) error {
	w, ok := mod.ws[info.Name]
	if !ok {
		mod.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
		return ErrWorkerNotExist
	}
	switch event.Type {
	case dm.TypeReportEvent:
		if err := w.Execute(); err != nil {
			return err
		}
	default:
		return errors.New("event type not supported yet")
	}
	return nil
}

func (mod *Modbus) Start() {
	for _, worker := range mod.ws {
		go mod.working(worker)
	}
}

func (mod *Modbus) working(w *Worker) {
	ticker := time.NewTicker(w.job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := w.Execute()
			if err != nil {
				mod.log.Error("failed to execute job", v2log.Error(err))
			}
		case <-mod.ctx.WaitChan():
			mod.log.Warn("worker stopped", v2log.Any("worker", w))
			return
		}
	}
}

func (mod *Modbus) Close() error {
	for _, slave := range mod.slaves {
		if err := slave.client.Close(); err != nil {
			mod.log.Warn("failed to close slave", v2log.Any("slave id", slave.cfg.Id), v2log.Error(err))
		}
	}
	return nil
}

func transform(value interface{}, cfg MapConfig) ([]byte, error) {
	var order binary.ByteOrder = binary.BigEndian
	if cfg.Function == HoldingRegister && cfg.SwapByte {
		order = binary.LittleEndian
	}
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, order, value)
	if err != nil {
		return nil, errors.Trace(err)
	}
	bs := buf.Bytes()
	if cfg.Function == HoldingRegister && cfg.SwapRegister {
		for i := 0; i < len(bs)-1; i += 2 {
			bs[i], bs[i+1] = bs[i+1], bs[i]
		}
	}
	return bs, nil
}
