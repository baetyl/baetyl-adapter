package modbus

import (
	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
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
	r := make(map[string]interface{})
	for _, m := range w.maps {
		p, err := m.Collect()
		if err != nil {
			return err
		}
		pa, err := m.Parse(p[4:])
		if err != nil {
			return err
		}
		r[m.cfg.Field.Name] = pa
	}
	if err := w.ctx.ReportDeviceProperties(w.slave.dev, r); err != nil {
		return errors.Trace(err)
	}
	w.log.Debug("report properties successfully", log.Any("content", r))
	return nil
}
