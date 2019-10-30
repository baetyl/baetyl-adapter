package main

import "github.com/baetyl/baetyl/logger"

type read func(address, quantity uint16) (results []byte, err error)

type Map struct {
	cfg MapConfig
	r   read
	s   *Slave
	log logger.Logger
}

func NewMap(item MapConfig, s *Slave, log logger.Logger) *Map {
	var r read
	switch item.Function {
	case 1:
		r = s.client.ReadCoils
	case 2:
		r = s.client.ReadDiscreteInputs
	case 3:
		r = s.client.ReadHoldingRegisters
	case 4:
		r = s.client.ReadInputRegisters
	default:
	}
	return &Map{
		cfg: item,
		r:   r,
		s:   s,
		log: log,
	}
}

func (m *Map) Read() (results []byte, err error) {
	results, err = m.r(m.cfg.Address, m.cfg.Quantity)
	if err != nil {
		// reconnect slave
		m.s.client.Close()
		err2 := m.s.client.Connect()
		if err2 != nil {
			m.log.Errorf("failed to reconnect slave id=%d: %s", m.cfg.SlaveID, err2.Error())
		} else {
			return m.r(m.cfg.Address, m.cfg.Quantity)
		}
		return nil, err
	}
	return results, err
}
