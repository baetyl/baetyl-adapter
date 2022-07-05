package modbus

import (
	"time"

	"github.com/baetyl/baetyl-adapter/v2/dmp"
	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/google/uuid"
)

type Worker struct {
	ctx   dmcontext.Context
	job   Job
	maps  []*Map
	slave *Slave
	log   *log.Logger
}

func NewWorker(ctx dmcontext.Context, job Job, slave *Slave, log *log.Logger) *Worker {
	w := &Worker{
		ctx:   ctx,
		job:   job,
		slave: slave,
		log:   log,
	}
	for _, m := range job.Maps {
		m := NewMap(ctx, m, slave, log)
		w.maps = append(w.maps, m)
	}
	return w
}

func (w *Worker) Execute() error {
	r := v1.Report{}
	for _, m := range w.maps {
		p, err := m.Collect()
		if err != nil {
			if err1 := w.slave.UpdateStatus(SlaveOffline); err1 != nil {
				w.log.Error("failed to update slave status", log.Any("error", err1), log.Any("status", "offline"))
			}
			return err
		}
		pa, err := m.Parse(p[4:])
		if err != nil {
			return err
		}
		r[m.cfg.Name] = pa
	}

	// add dmp field
	reqId := uuid.New().String()
	timestamp := time.Now().UnixNano() / 1e6
	events := make(map[string]interface{})
	bie := make(map[string]interface{})
	accessTemplate, err := w.ctx.GetAccessTemplates(w.slave.dev)
	if err != nil {
		return err
	}
	for _, model := range accessTemplate.Mappings {
		args := make(map[string]interface{})
		params, err := dmcontext.ParseExpression(model.Expression)
		if err != nil {
			return err
		}
		for _, param := range params {
			id := param[1:]
			mappingName, err := dmp.GetMappingName(id, accessTemplate)
			if err != nil {
				return err
			}
			args[param] = r[mappingName]
		}
		modelValue, err := dmcontext.ExecExpression(model.Expression, args, model.Type)
		if err != nil {
			return err
		}
		bie[model.Attribute] = modelValue
	}
	events[dmp.BIEKey] = bie
	r[dmp.DMPKey] = dmp.DMP{
		ReqId:     reqId,
		Method:    dmp.Method,
		Version:   dmp.Version,
		Timestamp: timestamp,
		BindName:  dmp.BindName,
		Events:    events,
	}

	if err := w.slave.UpdateStatus(SlaveOnline); err != nil {
		w.log.Error("failed to update slave status", log.Any("error", err), log.Any("status", "online"))
	}
	if err := w.ctx.ReportDeviceProperties(w.slave.dev, r); err != nil {
		return errors.Trace(err)
	}
	w.log.Debug("report properties successfully", log.Any("content", r))
	return nil
}
