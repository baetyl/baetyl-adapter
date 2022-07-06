package opcua

import (
	"github.com/baetyl/baetyl-adapter/v2/dmp"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	v2log "github.com/baetyl/baetyl-go/v2/log"
	mqtt2 "github.com/baetyl/baetyl-go/v2/mqtt"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gopcua/opcua/ua"
)

var (
	ErrWorkerNotExist = errors.New("worker not exist")
)

type Opcua struct {
	cfg  Config
	ctx  context.Context
	ws   map[string]*Worker
	log  *v2log.Logger
	mqtt *mqtt2.Client
}

func NewOpcua(ctx dm.Context, cfg Config) (*Opcua, error) {
	infos := make(map[string]dm.DeviceInfo)
	for _, info := range ctx.GetAllDevices() {
		infos[info.Name] = info
	}
	log := ctx.Log().With(v2log.Any("module", "opcua"))
	devs := make(map[string]*Device)
	for _, dCfg := range cfg.Devices {
		if info, ok := infos[dCfg.Device]; ok {
			dev, err := NewDevice(&info, dCfg)
			if err != nil {
				log.Error("ignore device which failed to establish connection", v2log.Any("device", dCfg.Device), v2log.Error(err))
				continue
			}
			devs[dCfg.Device] = dev
			err = ctx.Online(&info)
			if err != nil {
				log.Error("failed to report device status", v2log.Any("device", dCfg.Device))
			}
		}
	}
	ws := make(map[string]*Worker)
	for _, job := range cfg.Jobs {
		if dev := devs[job.Device]; dev != nil {
			ws[dev.info.Name] = NewWorker(ctx, job, dev, log)
		} else {
			log.Error("device of job not exist", v2log.Any("device id", job.Device))
		}
	}
	o := &Opcua{
		ctx: ctx,
		cfg: cfg,
		ws:  ws,
		log: log,
	}
	err := ctx.RegisterDeltaCallback(o.DeltaCallback)
	if err != nil {
		return nil, err
	}
	err = ctx.RegisterEventCallback(o.EventCallback)
	if err != nil {
		return nil, err
	}
	for _, worker := range o.ws {
		go o.working(worker)
	}
	return o, nil
}

func (o *Opcua) working(w *Worker) {
	ticker := time.NewTicker(w.job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := w.Execute()
			if err != nil {
				o.log.Error("failed to execute job", v2log.Error(err))
			}
		case <-o.ctx.WaitChan():
			o.log.Warn("worker stopped", v2log.Any("worker", w))
			return
		}
	}
}

func (o *Opcua) Close() error {
	return nil
}

func (o *Opcua) DeltaCallback(info *dm.DeviceInfo, delta v1.Delta) error {
	w, ok := o.ws[info.Name]
	if !ok {
		o.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
		return ErrWorkerNotExist
	}
	accessTemplate, err := w.ctx.GetAccessTemplates(info)
	if err != nil {
		o.log.Warn("get access template err", v2log.Any("device", info.Name))
		return err
	}
	for key, val := range delta {
		// support model Number-setting
		id, err := dmp.GetConfigIdByModelName(key, accessTemplate)
		if id == "" || err != nil {
			o.log.Warn("prop not exist", v2log.Any("name", key))
			continue
		}
		propName, err := dmp.GetMappingName(id, accessTemplate)
		if err != nil {
			o.log.Warn("prop name not exist", v2log.Any("id", id))
			continue
		}
		propVal, err := dmp.GetPropValueByModelName(key, val, accessTemplate)
		if err != nil {
			o.log.Warn("get prop value err", v2log.Any("name", propName))
			continue
		}

		for _, prop := range w.job.Properties {
			if propName == prop.Name {
				value, err := dmp.ParsePropertyValue(prop.Type, propVal)
				if err != nil {
					return errors.Trace(err)
				}
				variant, err := ua.NewVariant(value)
				if err != nil {
					return errors.Trace(err)
				}
				id, err := ua.ParseNodeID(prop.NodeID)
				if err != nil {
					return errors.Trace(err)
				}
				req := &ua.WriteRequest{
					NodesToWrite: []*ua.WriteValue{{
						NodeID:      id,
						AttributeID: ua.AttributeIDValue,
						Value: &ua.DataValue{
							EncodingMask: ua.DataValueValue,
							Value:        variant,
						}}},
				}
				_, err = w.device.opcuaClient.Write(req)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (o *Opcua) EventCallback(info *dm.DeviceInfo, event *dm.Event) error {
	w, ok := o.ws[info.Name]
	if !ok {
		o.log.Warn("worker not exist according to device", v2log.Any("device", info.Name))
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
