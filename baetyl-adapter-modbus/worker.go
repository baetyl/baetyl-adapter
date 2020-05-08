package main

import (
	"encoding/json"
	"fmt"
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl-go/log"
	"time"

	"encoding/binary"
	"github.com/baetyl/baetyl-go/mqtt"
)

type Worker struct {
	job    Job
	mqtt   *mqtt.Client
	maps   []*Map
	logger *log.Logger
}

func NewWorker(job Job, slave *Slave, mqtt *mqtt.Client, logger *log.Logger) *Worker {
	w := &Worker{
		job:    job,
		mqtt:   mqtt,
		logger: logger,
	}
	for _, m := range job.Maps {
		m := NewMap(m, slave, logger)
		w.maps = append(w.maps, m)
	}
	return w
}

func (w *Worker) Execute(publish Publish) error {
	var pld []byte
	result := map[string]interface{}{}
	now := time.Now()
	var ts int64
	if w.job.Time.Precision == SecondPrecision {
		ts = now.Unix()
	} else if w.job.Time.Precision == NanoPrecision {
		ts = now.UnixNano()
	}
	if w.job.Kind == BinaryKind {
		pld = append(pld, w.job.SlaveId)
		tp := make([]byte, 8)
		binary.BigEndian.PutUint64(tp, uint64(ts))
		pld = append(pld, tp...)
		for _, m := range w.maps {
			p, err := m.Collect()
			if err != nil {
				return err
			}
			pld = append(pld, p...)
		}
	} else if w.job.Kind == JsonKind {
		if w.job.Time.Type == LongTimeType {
			result[TimeKey] = ts
		} else if w.job.Time.Type == StringTimeType {
			if w.job.Time.Format != "" {
				result[TimeKey] = now.Format(w.job.Time.Format)
			} else {
				result[TimeKey] = now.String()
			}
		}
		for _, m := range w.maps {
			if m.cfg.Field == nil {
				w.logger.Error("field can not be null when parsing", log.Any("map", m))
				continue
			}
			p, err := m.Parse()
			if err != nil {
				return err
			}
			result[m.cfg.Field.Name] = p
		}
		var err error
		pld, err = json.Marshal(result)
		if err != nil {
			return err
		}
	}
	if w.mqtt != nil {
		pkt := packet.NewPublish()
		pkt.Message.Topic = publish.Topic
		pkt.Message.QOS = packet.QOS(publish.QOS)
		pkt.Message.Payload = pld
		err := w.mqtt.Send(pkt)
		if err != nil {
			return fmt.Errorf("failed to publish: %s", err.Error())
		}
	}
	return nil
}
