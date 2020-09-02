package opcua

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/gopcua/opcua/ua"
)

type Worker struct {
	job    Job
	device *Device
	sender Sender
	logger *log.Logger
}

func NewWorker(job Job, device *Device, s Sender, logger *log.Logger) *Worker {
	w := &Worker{
		device: device,
		job:    job,
		sender: s,
		logger: logger,
	}
	return w
}

func (w *Worker) Execute() error {
	res := map[string]interface{}{}
	now := time.Now()
	var ts int64
	if w.job.Time.Precision == SecondPrecision {
		ts = now.Unix()
	} else if w.job.Time.Precision == NanoPrecision {
		ts = now.UnixNano()
	}
	if w.job.Time.Type == IntegerTime {
		res[w.job.Time.Name] = ts
	} else if w.job.Time.Type == StringTime {
		res[w.job.Time.Name] = now.Format(w.job.Time.Format)
	}
	res[DeviceID] = w.job.DeviceID

	data := make(map[string]interface{})
	for _, p := range w.job.Properties {
		val, err := w.read(p)
		if err != nil {
			w.logger.Error("failed to read", log.Error(err))
			continue
		}
		value, err := variant2Value(p.Type, val)
		if err != nil {
			w.logger.Error("failed to parse", log.Error(err))
			continue
		}
		data[p.Name] = value
	}
	// TODO if length of data equal to 0, do not send msg ?
	res["attr"] = data

	pld, err := json.Marshal(res)
	if err != nil {
		return err
	}
	if w.sender != nil {
		if err := w.sender.Send(pld); err != nil {
			return fmt.Errorf("failed to publish: %s", err.Error())
		}
	}
	return nil
}

func (w *Worker) read(prop Property) (*ua.Variant, error) {
	id, err := ua.ParseNodeID(prop.NodeID)
	if err != nil {
		w.logger.Error("invalid node id", log.Any("nodeid", prop.NodeID))
		return nil, errors.Errorf("invalid node id [%s]", prop.NodeID)
	}
	req := &ua.ReadRequest{
		MaxAge:             2000,
		NodesToRead:        []*ua.ReadValueID{{NodeID: id}},
		TimestampsToReturn: ua.TimestampsToReturnNeither,
	}
	resp, err := w.device.opcuaClient.Read(req)
	if err != nil {
		w.logger.Error("failed to read", log.Any("nodeid", prop.NodeID), log.Error(err))
		return nil, err
		//w.device.opcuaClient.Close()
		//var ctx, cancel = context.WithTimeout(context.Background(), w.device.cfg.Timeout)
		//defer cancel()
		//err2 := w.device.opcuaClient.Connect(ctx)
		//if err2 != nil {
		//	w.logger.Error("failed to reconnect: ", log.Any("deviceid", w.device.cfg.ID), log.Error(err2))
		//	return nil, err2
		//} else {
		//	resp, err = w.device.opcuaClient.Read(req)
		//	if err != nil {
		//		w.logger.Error("failed to read after reconnect", log.Any("nodeid", prop.NodeID), log.Error(err))
		//		return nil, err
		//	}
		//}
	}
	if resp == nil || len(resp.Results) == 0 {
		w.logger.Error("invalid read response", log.Any("nodeid", prop.NodeID))
		return nil, errors.Errorf("invalid read response, nodeid: [%s]", prop.NodeID)
	}
	if resp.Results[0].Status != ua.StatusOK {
		w.logger.Error("Status not OK: %v", log.Any("nodeid", prop.NodeID), log.Any("status", resp.Results[0].Status))
		return nil, errors.Errorf("status [%d] not ok", resp.Results[0].Status)
	}

	return resp.Results[0].Value, nil
}
