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

type observer struct {
	slaves map[byte]*Slave
	log    *log.Logger
}

type CtrData struct {
	SlaveID    byte                   `yaml:"slaveid" json:"slaveid"`
	Attributes map[string]interface{} `yaml:"attr" json:"attr"`
}

func NewObserver(slaves map[byte]*Slave, log *log.Logger) mqtt.Observer {
	return &observer{
		slaves: slaves,
		log:    log,
	}
}

func (o *observer) OnPublish(pkt *packet.Publish) error {
	var ctrData CtrData
	err := json.Unmarshal(pkt.Message.Payload, &ctrData)
	if err != nil {
		return errors.Trace(err)
	}
	var slave *Slave
	slave, ok := o.slaves[ctrData.SlaveID]
	if !ok {
		o.log.Error("device to write data not exist", log.Any("id", ctrData.SlaveID))
		return errors.Errorf("device to write data not exist")
	}
	if err := o.Write(slave, ctrData.Attributes); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (o *observer) Write(slave *Slave, attr map[string]interface{}) error {
	config, ok := configRecoder[slave.cfg.ID]
	if !ok {
		o.log.Error("map config of slave id not exist", log.Any("id", slave.cfg.ID))
		return errors.Errorf("map config of slave id [%d] not exist", slave.cfg.ID)
	}
	for key, val := range attr {
		cfg, ok := config[key]
		if !ok {
			o.log.Warn("ignore key whose map config not exist", log.Any("key", key))
			continue
		}
		value, err := validateAndTransform(val, cfg.Field.Type)
		if err != nil {
			o.log.Warn("ignore illegal data type of val", log.Any("value", val), log.Any("type", cfg.Field.Type))
			continue
		}
		switch cfg.Function {
		case DiscreteInput:
		case InputRegister:
			return errors.Errorf("can not write data with illegal function code: [%d]", cfg.Function)
		case Coil:
			if _, err := slave.client.WriteMultipleCoils(cfg.Address, cfg.Quantity, value); err != nil {
				return err
			}
		case HoldingRegister:
			if _, err := slave.client.WriteMultipleRegisters(cfg.Address, cfg.Quantity, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateAndTransform(source interface{}, fieldType string) ([]byte, error) {
	var value interface{}
	var ok bool
	var num float64
	if fieldType != Bool {
		num, ok = source.(float64)
	}
	switch fieldType {
	case Bool:
		var b bool
		b, ok = source.(bool)
		value = b
	case Int16:
		i16 := int16(num)
		value = i16
	case UInt16:
		u16 := uint16(num)
		value = u16
	case Int32:
		i32 := int32(num)
		value = i32
	case UInt32:
		u32 := uint32(num)
		value = u32
	case Int64:
		i64 := int64(num)
		value = i64
	case UInt64:
		u64 := uint64(num)
		value = u64
	case Float32:
		f32 := float32(num)
		value = f32
	case Float64:
		value = num
	default:
		return nil, errors.Errorf("unsupported field type [%s]", fieldType)
	}
	if !ok {
		return nil, errors.Errorf("value [%v] not compatible with type [%s] ", source, fieldType)
	}
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return buf.Bytes(), nil
}

func (o *observer) OnPuback(pkt *packet.Puback) error {
	return nil
}

func (o *observer) OnError(err error) {
	o.log.Error("receive mqtt message error", log.Error(err))
}
