package main

import (
	"encoding/binary"
	"time"

	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/logger"
	"github.com/baetyl/baetyl/protocol/mqtt"
)

type task struct {
	mbCli     *client
	mqttCli   *mqtt.Dispatcher
	runTime   time.Time
	parseItem ParseItem
	msgChan   chan *packet.Publish
	publish   *Publish
	log       logger.Logger
	index     int // index of task in priority queue
}

func newTask(mbCli *client, mqttCli *mqtt.Dispatcher, item ParseItem, msgChan chan *packet.Publish,
	publish *Publish, log logger.Logger) *task {
	return &task{
		mbCli:     mbCli,
		mqttCli:   mqttCli,
		parseItem: item,
		msgChan:   msgChan,
		publish:   publish,
		log:       log,
		runTime:   time.Now().Add(item.Interval),
	}
}

func (t *task) timeUp() bool {
	now := time.Now()
	if t.runTime.Sub(now) < 0 {
		t.runTime = t.runTime.Add(t.parseItem.Interval)
		return true
	}
	return false
}

func (t *task) execute() {
	results, err := t.mbCli.read(t.parseItem)
	if err != nil {
		logger.Errorf("failed to collect data from slave: id=%d function=%d address=%d quantity=%d",
			t.parseItem.SlaveID, t.parseItem.Function, t.parseItem.Address, t.parseItem.Quantity)
		return
	}
	pld := make([]byte, 9+t.parseItem.Quantity*2)
	pld[0] = t.parseItem.SlaveID
	binary.BigEndian.PutUint16(pld[1:], t.parseItem.Address)
	binary.BigEndian.PutUint16(pld[3:], t.parseItem.Quantity)
	binary.BigEndian.PutUint32(pld[5:], uint32(time.Now().Unix()))
	copy(pld[9:], results)
	pkt := packet.NewPublish()
	pkt.Message.Topic = t.publish.Topic
	pkt.Message.QOS = packet.QOS(t.publish.QOS)
	pkt.Message.Payload = pld
	t.msgChan <- pkt
}
