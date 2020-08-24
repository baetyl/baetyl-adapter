package modbus

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/mqtt"
)

type Item struct {
	Function int     `yaml:"function" json:"function"`
	Address  uint16  `yaml:"address" json:"address"`
	Quantity uint16  `yaml:"quantity" json:"quantity"`
	Value    []int16 `yaml:"value" json:"value"`
}

type Data struct {
	SlaveID byte   `yaml:"slaveid" json:"slaveid"`
	Items   []Item `yaml:"items" json:"items"`
}

type observer struct {
	slaves map[byte]*Slave
	log    *log.Logger
}

func NewObserver(slaves map[byte]*Slave, log *log.Logger) mqtt.Observer {
	return &observer{
		slaves: slaves,
		log:    log,
	}
}

func (o *observer) OnPublish(pkt *packet.Publish) error {
	var datas []Data
	err := json.Unmarshal(pkt.Message.Payload, &datas)
	if err != nil {
		return errors.Trace(err)
	}
	for _, data := range datas {
		if slave, ok := o.slaves[data.SlaveID]; !ok {
			o.log.Error("device to write data not exist", log.Any("id", data.SlaveID))
		} else {
			if err := o.Write(slave, data); err != nil {
				o.log.Error("write data failed", log.Any("data", data), log.Error(err))
			}
		}
	}
	return nil
}

func (o *observer) Write(slave *Slave, data Data) error {
	for _, item := range data.Items {
		value, err := validateAndTransform(item)
		if err != nil {
			return err
		}
		switch item.Function {
		case DiscreteInput:
		case InputRegister:
			return errors.Errorf("illegal function code: [%v]", data)
		case Coil:
			if _, err := slave.client.WriteMultipleCoils(item.Address, item.Quantity, value); err != nil {
				return err
			}
		case HoldingRegister:
			if _, err := slave.client.WriteMultipleRegisters(item.Address, item.Quantity, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateAndTransform(item Item) ([]byte, error) {
	if int(item.Quantity) != len(item.Value) {
		return nil, errors.Errorf("quantity not equal to value length")
	}
	if item.Function == Coil {
		b := make([]byte, (len(item.Value)+7)/8)
		for i, x := range item.Value {
			if x != 0 {
				b[i/8] |= 0x1 << uint(i%8)
			}
		}
		return b, nil
	} else if item.Function == HoldingRegister {
		buf := bytes.NewBuffer(nil)
		err := binary.Write(buf, binary.BigEndian, item.Value)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return nil, nil
}

func (o *observer) OnPuback(pkt *packet.Puback) error {
	return nil
}

func (o *observer) OnError(err error) {
	o.log.Error("receive mqtt message error", log.Error(err))
}
