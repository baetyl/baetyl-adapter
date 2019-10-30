package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/baetyl/baetyl/logger"
)

type read func(address, quantity uint16) (results []byte, err error)

type Map struct {
	cfg MapConfig
	r   read
	s   *Slave
	log logger.Logger
}

func NewMap(cfg MapConfig, s *Slave, log logger.Logger) *Map {
	var r read
	switch cfg.Function {
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
		cfg: cfg,
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
		if err2 == nil {
			return m.r(m.cfg.Address, m.cfg.Quantity)
		}
		m.log.WithField("slave", m.cfg.SlaveID).Errorf("failed to reconnect: %s", err2.Error())
		return nil, err
	}
	return results, err
}

func (m *Map) Package() ([]byte, error) {
	results, err := m.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to collect data from slave.go: id=%d function=%d address=%d quantity=%d",
			m.cfg.SlaveID, m.cfg.Function, m.cfg.Address, m.cfg.Quantity)
	}
	pld := make([]byte, 9+m.cfg.Quantity*2)
	pld[0] = m.cfg.SlaveID
	binary.BigEndian.PutUint16(pld[1:], m.cfg.Address)
	binary.BigEndian.PutUint16(pld[3:], m.cfg.Quantity)
	binary.BigEndian.PutUint32(pld[5:], uint32(time.Now().Unix()))
	copy(pld[9:], results)
	return pld, nil
}
