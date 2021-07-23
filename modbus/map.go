package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
)

var (
	ErrClientInvalid = errors.New("device client is invalid")
)

type Map struct {
	ctx dm.Context
	cfg MapConfig
	s   *Slave
	log *log.Logger
}

func NewMap(ctx dm.Context, cfg MapConfig, s *Slave, log *log.Logger) *Map {
	return &Map{
		ctx: ctx,
		cfg: cfg,
		s:   s,
		log: log,
	}
}

func (m *Map) read() (results []byte, err error) {
	if m.s.client.Client == nil {
		return nil, errors.Trace(ErrClientInvalid)
	}
	switch m.cfg.Function {
	case Coil:
		return m.s.client.ReadCoils(m.cfg.Address, m.cfg.Quantity)
	case DiscreteInput:
		return m.s.client.ReadDiscreteInputs(m.cfg.Address, m.cfg.Quantity)
	case HoldingRegister:
		return m.s.client.ReadHoldingRegisters(m.cfg.Address, m.cfg.Quantity)
	case InputRegister:
		return m.s.client.ReadInputRegisters(m.cfg.Address, m.cfg.Quantity)
	default:
		return
	}
}

func (m *Map) Collect() ([]byte, error) {
	res, err := m.read()
	if err != nil {
		m.log.Error("failed to collect data from slave", log.Any("slave id", m.s.cfg.Id), log.Any("config", m.cfg), log.Error(err))
		if err1 := m.s.client.Reconnect(); err1 == nil {
			m.log.Info("reconnected successfully", log.Any("slave id", m.s.cfg.Id))
			if err2 := m.ctx.Online(m.s.dev); err2 != nil {
				m.log.Error("failed to report online status", log.Any("slave id", m.s.cfg.Id), log.Error(err2))
			}
		} else {
			m.log.Error("failed to reconnect", log.Any("slave id", m.s.cfg.Id), log.Error(err1))
			return nil, err
		}
		// try to read again
		res, err = m.read()
		if err != nil {
			return nil, err
		}
	}
	pld := make([]byte, 4)
	if m.cfg.SwapByte {
		binary.LittleEndian.PutUint16(pld, m.cfg.Address)
		binary.LittleEndian.PutUint16(pld[2:], m.cfg.Quantity)
	} else {
		binary.BigEndian.PutUint16(pld, m.cfg.Address)
		binary.BigEndian.PutUint16(pld[2:], m.cfg.Quantity)
	}
	pld = append(pld, res...)
	return pld, nil
}

func (m *Map) Parse(data []byte) (interface{}, error) {
	if m.cfg.Function == Coil || m.cfg.Function == DiscreteInput {
		if len(data) != 1 {
			return nil, fmt.Errorf("quantity should be 1 when parsing coil or discrete input")
		}
		data[0] = data[0] & 0x1
		if m.cfg.Type != Bool {
			return nil, fmt.Errorf("field type should be bool when parsing coil or discrete input")
		}
	} else {
		if m.cfg.SwapRegister {
			for i := 0; i < len(data)-1; i += 2 {
				data[i], data[i+1] = data[i+1], data[i]
			}
		}
	}
	p, err := parse(bytes.NewBuffer(data), m.cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parse(reader io.Reader, cfg MapConfig) (res interface{}, err error) {
	var order binary.ByteOrder = binary.BigEndian
	if cfg.SwapByte {
		order = binary.LittleEndian
	}
	switch cfg.Type {
	case Bool:
		var b bool
		err = binary.Read(reader, order, &b)
		res = b
	case Int16:
		var i16 int16
		err = binary.Read(reader, order, &i16)
		res = i16
	case UInt16:
		var u16 uint16
		err = binary.Read(reader, order, &u16)
		res = u16
	case Int32:
		var i32 int32
		err = binary.Read(reader, order, &i32)
		res = i32
	case UInt32:
		var u32 uint32
		err = binary.Read(reader, order, &u32)
		res = u32
	case Int64:
		var i64 int64
		err = binary.Read(reader, order, &i64)
		res = i64
	case UInt64:
		var u64 uint64
		err = binary.Read(reader, order, &u64)
		res = u64
	case Float32:
		var f32 float32
		err = binary.Read(reader, order, &f32)
		res = f32
	case Float64:
		var f64 float64
		err = binary.Read(reader, order, &f64)
		res = f64
	default:
		err = errors.New("unsupported field type")
	}
	return res, err
}
