package opcua

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/mqtt"
)

var configRecoder = make(map[byte]map[string]Property)

type Opcua struct {
	cfg     Config
	ctx     context.Context
	devices map[byte]*Device
	logger  *log.Logger
	mqtt    *mqtt.Client
}

func NewOpcua(ctx context.Context, cfg Config) (*Opcua, error) {
	devices := map[byte]*Device{}
	for _, dCfg := range cfg.Devices {
		dev, err := NewDevice(dCfg)
		if err != nil {
			ctx.Log().Error("ignore device which failed to establish connection", log.Any("id", dCfg.ID), log.Error(err))
			continue
		}
		devices[dCfg.ID] = dev
		configRecoder[dCfg.ID] = make(map[string]Property)
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
	observer := NewObserver(devices, ctx.Log())
	if err := mqtt.Start(observer); err != nil {
		return nil, err
	}
	o := &Opcua{
		cfg:     cfg,
		ctx:     ctx,
		mqtt:    mqtt,
		logger:  ctx.Log(),
		devices: devices,
	}
	var ws []*Worker
	for _, job := range cfg.Jobs {
		if device := devices[job.DeviceID]; device != nil {
			if job.Publish.Topic == "" {
				job.Publish.Topic = fmt.Sprintf("%s/%s", ctx.ServiceName(), job.DeviceID)
			}
			sender := NewSender(job.Publish, mqtt)
			w := NewWorker(job, device, sender, log.With(log.Any("deviceid", job.DeviceID)))
			ws = append(ws, w)
			if _, ok := configRecoder[job.DeviceID]; !ok {
				continue
			}
			for _, p := range job.Properties {
				if _, ok := configRecoder[job.DeviceID][p.Name]; ok {
					return nil, errors.Errorf("one device should not have same variable definition")
				}
				configRecoder[job.DeviceID][p.Name] = p
			}
		} else {
			ctx.Log().Error("device of job not exist")
		}
	}
	for _, worker := range ws {
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
				o.logger.Error("failed to execute job", log.Error(err))
			}
		case <-o.ctx.WaitChan():
			o.logger.Warn("worker stopped", log.Any("worker", w))
			return
		}
	}
}

func (o *Opcua) Close() error {
	return nil
}
