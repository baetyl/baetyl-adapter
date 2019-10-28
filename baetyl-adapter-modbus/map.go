package main

type read func(address, quantity uint16) (results []byte, err error)

type Map struct {
	cfg MapItem
	r   read
	s   *Slave
}

func NewMap(item MapItem, s *Slave) *Map {
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
	}
}

func (m *Map) Read() (results []byte, err error) {
	results, err = m.r(m.cfg.Address, m.cfg.Quantity)
	if err != nil {
		// try connect again
		m.s.client.Close()
		m.s.client.Connect()
		return nil, err
	}
	return results, err
}
