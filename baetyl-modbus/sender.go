package main

import (
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/protocol/mqtt"
	"github.com/baetyl/baetyl/utils"
	"github.com/baidubce/bce-sdk-go/util/log"
)

type sender struct {
	client *mqtt.Dispatcher
	msgChan chan *packet.Publish
	tomb utils.Tomb
}

func newSender(cli *mqtt.Dispatcher, msgChan chan *packet.Publish) *sender {
	return &sender {
		client: cli,
		msgChan: msgChan,
	}
}

func (s *sender) start() error {
	return s.tomb.Go(s.processMsg)
}

func (s *sender) processMsg() error {
	for {
		select {
		case <- s.tomb.Dying():
			return nil
		case m := <- s.msgChan:
			err := s.client.Send(m)
			if err != nil {
				log.Errorf("failed to publish: %s", err.Error())
			}
		}
	}
}

func (s *sender) close() {
	s.tomb.Kill(nil)
	s.tomb.Wait()
}