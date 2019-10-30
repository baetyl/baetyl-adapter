package main

import (
	"fmt"
	"time"

	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/protocol/mqtt"
)

type mb struct {
	Interval time.Duration
	ms       []*Map
}

func NewMb(interval time.Duration) *mb {
	return &mb{
		Interval: interval,
	}
}

func (m *mb) AddMap(ma *Map) {
	m.ms = append(m.ms, ma)
}

func (m *mb) Execute(mqttCli *mqtt.Dispatcher, publish Publish) error {
	for _, ma := range m.ms {
		payload, err := ma.Package()
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
