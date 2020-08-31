package opc

import (
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl-go/v2/mqtt"
)

//go:generate mockgen -destination=mock/sender.go -package=mock github.com/baetyl/baetyl-adapter/opc Sender

type Sender interface {
	Send(msg []byte) error
	Close() error
}

type mqttSender struct {
	*mqtt.Client
	Publish Publish
}

func NewSender(publish Publish, client *mqtt.Client) Sender {
	return mqttSender{
		Client:  client,
		Publish: publish,
	}
}

func (s mqttSender) Send(msg []byte) error {
	pkt := packet.NewPublish()
	pkt.Message.Topic = s.Publish.Topic
	pkt.Message.QOS = packet.QOS(s.Publish.QOS)
	pkt.Message.Payload = msg
	if err := s.Client.Send(pkt); err != nil {
		return err
	}
	return nil
}

func (s mqttSender) Close() error {
	return s.Client.Close()
}
