package modbus

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	v2log "github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

var (
	ErrWorkerNotExist       = errors.New("worker not exist")
	ErrIllegalValueType     = errors.New("illegal value type")
	ErrUnsupportedValueType = errors.New("value type not supported")
)

type Modbus struct {
	ctx    dmcontext.Context
	log    *v2log.Logger
	slaves map[byte]*Slave
	ws     map[string]*Worker
}

func NewModbus(ctx dmcontext.Context, cfg Config) (*Modbus, error) {
	devices := ctx.GetAllDevices()
	devMap := map[string]dmcontext.DeviceInfo{}
	for _, dev := range devices {
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
			log.Warn("connect failed", v2log.Any("id", slaveConfig.ID), v2log.Error(err))
			continue
		}
		slaves[slaveConfig.ID] = NewSlave(slaveConfig, client)
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
				log.Error("can not find device according to job config")
				continue
			}
			mod.ws[dev.Name] = NewWorker(ctx, job, slave, &dev, log)
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
	for _, dev := range devices {
		if err := ctx.Online(&dev); err != nil {
			return nil, err
		}
	}
	return mod, nil
}

func (mod *Modbus) DeltaCallback(info *dmcontext.DeviceInfo, prop v1.Delta) error {
	w, ok := mod.ws[info.Name]
	if !ok {
		mod.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
		return ErrWorkerNotExist
	}
	ms := map[string]MapConfig{}
	for _, m := range w.job.Maps {
		ms[m.Field.Name] = m
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
		value, err := validateAndTransform(val, cfg.Field.Type)
		if err != nil {
			mod.log.Warn("ignore illegal data type of val", v2log.Any("value", val), v2log.Any("type", cfg.Field.Type), v2log.Error(err))
			continue
		}
		switch cfg.Function {
		case DiscreteInput:
		case InputRegister:
			return fmt.Errorf("can not write data with illegal function code: [%d]", cfg.Function)
		case Coil:
			if _, err := slave.client.WriteMultipleCoils(cfg.Address, cfg.Quantity, value); err != nil {
				return err
			}
		case HoldingRegister:
			if _, err := slave.client.WriteMultipleRegisters(cfg.Address, cfg.Quantity, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mod *Modbus) EventCallback(info *dmcontext.DeviceInfo, event *dmcontext.Event) error {
	w, ok := mod.ws[info.Name]
	if !ok {
		mod.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
		return ErrWorkerNotExist
	}
	switch event.Type {
	case dmcontext.TypeReportEvent:
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
			mod.log.Warn("failed to close slave", v2log.Any("slave id", slave.cfg.ID), v2log.Error(err))
		}
	}
	return nil
}

func validateAndTransform(value interface{}, fieldType string) ([]byte, error) {
	bs, err := json.Marshal(value)
	if err != nil {
		return nil, errors.Trace(err)
	}
	s := string(bs)
	var res interface{}
	switch fieldType {
	case Bool:
		var ok bool
		res, ok = value.(bool)
		if !ok {
			return nil, errors.Trace(ErrIllegalValueType)
		}
	case String:
		return bs, nil
	case Int16:
		i, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = int16(i)
	case UInt16:
		ui, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = uint16(ui)
	case Int32:
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = int32(i)
	case UInt32:
		ui, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = uint32(ui)
	case Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = i
	case UInt64:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = ui
	case Float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = float32(f)
	case Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		res = f
	default:
		return nil, errors.Trace(ErrUnsupportedValueType)
	}
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.BigEndian, res)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return buf.Bytes(), nil
}
