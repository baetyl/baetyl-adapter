package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl-go/log"
	"github.com/baetyl/baetyl-go/mqtt"
	"time"
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
	res := map[string]interface{}{}

	now := time.Now()
	var ts int64
	if w.job.Time.Precision == SecondPrecision {
		ts = now.Unix()
	} else if w.job.Time.Precision == NanoPrecision {
		ts = now.UnixNano()
	}
	if w.job.Encoding == BinaryEncoding {
		tp := make([]byte, 8)
		binary.BigEndian.PutUint64(tp, uint64(ts))
		pld = append(pld, tp...)
		// for aligning
		sp := make([]byte, 4)
		sp[0] = w.job.SlaveId
		pld = append(pld, sp...)
	} else if w.job.Encoding == JsonEncoding {
		if w.job.Time.Type == IntegerTime {
			res[w.job.Time.Name] = ts
		} else if w.job.Time.Type == StringTime {
			res[w.job.Time.Name] = now.Format(w.job.Time.Format)
		}
		res[SlaveId] = w.job.SlaveId
		var err error
		pld, err = json.Marshal(res)
		if err != nil {
			return err
		}
	}

	for _, m := range w.maps {
		p, err := m.Collect()
		if err != nil {
			return err
		}
		if w.job.Encoding == BinaryEncoding {
			pld = append(pld, p...)
		} else if w.job.Encoding == JsonEncoding {
			pa, err := m.Parse(p[4:])
			if err != nil {
				return err
			}
			res[m.cfg.Field.Name] = pa
		}
	}

	if w.job.Encoding == JsonEncoding {
		var err error
		pld, err = json.Marshal(res)
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
