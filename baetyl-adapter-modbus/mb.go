package main

import (
	"encoding/binary"
	"fmt"
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/protocol/mqtt"
	"time"
)

type mb struct {
	Interval time.Duration
	ms []*Map
}

func NewMb(interval time.Duration) *mb{
	return &mb {
		Interval: interval,
	}
}

func (m *mb) AddMap(ma *Map) {
	m.ms = append(m.ms, ma)
}

func Package(ma *Map) ([]byte, error){
	results, err := ma.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to collect data from slave.go: id=%d function=%d address=%d quantity=%d",
			ma.cfg.SlaveID, ma.cfg.Function, ma.cfg.Address, ma.cfg.Quantity)
	}
	pld := make([]byte, 9+ma.cfg.Quantity*2)
	pld[0] = ma.cfg.SlaveID
	binary.BigEndian.PutUint16(pld[1:], ma.cfg.Address)
	binary.BigEndian.PutUint16(pld[3:], ma.cfg.Quantity)
	binary.BigEndian.PutUint32(pld[5:], uint32(time.Now().Unix()))
	copy(pld[9:], results)
	return pld, nil
}

func (m *mb) Execute(mqttCli *mqtt.Dispatcher, publish Publish) error {
	for _, ma := range m.ms {
		payload, err := Package(ma)
		if err != nil {
			return err
		}
		pkt := packet.NewPublish()
		pkt.Message.Topic = publish.Topic
		pkt.Message.QOS = packet.QOS(publish.QOS)
		pkt.Message.Payload = payload
		err = mqttCli.Send(pkt)
		if err != nil {
			return fmt.Errorf("failed to publish: %s", err.Error())
		}
	}
	return nil
}


