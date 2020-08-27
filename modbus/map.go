package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
)

type read func(address, quantity uint16) (results []byte, err error)

type Map struct {
	cfg    MapConfig
	r      read
	s      *Slave
	logger *log.Logger
}

func NewMap(cfg MapConfig, s *Slave, logger *log.Logger) *Map {
	var r read
	switch cfg.Function {
	case Coil:
		r = s.client.ReadCoils
	case DiscreteInput:
		r = s.client.ReadDiscreteInputs
	case HoldingRegister:
		r = s.client.ReadHoldingRegisters
	case InputRegister:
		r = s.client.ReadInputRegisters
	default:
	}
	return &Map{
		cfg:    cfg,
		r:      r,
		s:      s,
		logger: logger,
	}
}

func (m *Map) read() (results []byte, err error) {
	results, err = m.r(m.cfg.Address, m.cfg.Quantity)
	if err != nil {
		// reconnect slave
		m.s.client.Close()
		err2 := m.s.client.Connect()
		if err2 == nil {
			return m.r(m.cfg.Address, m.cfg.Quantity)
		}
		m.logger.Error("failed to reconnect: ", log.Any("slave", m.s.cfg.ID), log.Error(err2))
		return nil, err
	}
	return results, err
}

func (m *Map) write(function int, address, quantity uint16, value int) error {
	switch function {
	case DiscreteInput:
	case InputRegister:
		return errors.Errorf("illegal function code")

	}
	return nil
}

func (m *Map) Collect() ([]byte, error) {
	res, err := m.read()
	if err != nil {
		return nil, fmt.Errorf("failed to collect data from slave.go: id=%d function=%d address=%d quantity=%d",
			m.s.cfg.ID, m.cfg.Function, m.cfg.Address, m.cfg.Quantity)
	}
	pld := make([]byte, 4)
	binary.BigEndian.PutUint16(pld, m.cfg.Address)
	binary.BigEndian.PutUint16(pld[2:], m.cfg.Quantity)
	pld = append(pld, res...)
	return pld, nil
}

func (m *Map) Parse(data []byte) (interface{}, error) {
	if m.cfg.Function == Coil || m.cfg.Function == DiscreteInput {
		if len(data) != 1 {
			return nil, fmt.Errorf("quantity should be 1 when parsing coil or discrete input")
		}
		data[0] = data[0] & 0x1
		if m.cfg.Field.Type != Bool {
			return nil, fmt.Errorf("field type should be bool when parsing coil or discrete input")
		}
	}
	p, err := parse(bytes.NewBuffer(data), m.cfg.Field.Type)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parse(reader io.Reader, fieldType string) (res interface{}, err error) {
	switch fieldType {
	case Bool:
		var b bool
		err = binary.Read(reader, binary.BigEndian, &b)
		res = b
	case Int16:
		var i16 int16
		err = binary.Read(reader, binary.BigEndian, &i16)
		res = i16
	case UInt16:
		var u16 uint16
		err = binary.Read(reader, binary.BigEndian, &u16)
		res = u16
	case Int32:
		var i32 int32
		err = binary.Read(reader, binary.BigEndian, &i32)
		res = i32
	case UInt32:
		var u32 uint32
		err = binary.Read(reader, binary.BigEndian, &u32)
		res = u32
	case Int64:
		var i64 int64
		err = binary.Read(reader, binary.BigEndian, &i64)
		res = i64
	case UInt64:
		var u64 uint64
		err = binary.Read(reader, binary.BigEndian, &u64)
		res = u64
	case Float32:
		var f32 float32
		err = binary.Read(reader, binary.BigEndian, &f32)
		res = f32
	case Float64:
		var f64 float64
		err = binary.Read(reader, binary.BigEndian, &f64)
		res = f64
	default:
		err = errors.New("unsupported field type")
	}
	return res, err
}
