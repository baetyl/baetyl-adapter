package modbus

import dm "github.com/baetyl/baetyl-go/v2/dmcontext"

const (
	SlaveOffline = 0
	SlaveOnline  = 1
)

type Slave struct {
	client *MbClient
	dev    *dm.DeviceInfo
	ctx    dm.Context
	cfg    SlaveConfig
	fail   int
	status int
}

func NewSlave(ctx dm.Context, dev *dm.DeviceInfo, cfg SlaveConfig, client *MbClient) *Slave {
	return &Slave{
		status: SlaveOffline,
		ctx:    ctx,
		dev:    dev,
		client: client,
		cfg:    cfg,
	}
}

func (s *Slave) UpdateStatus(status int) error {
	if status == s.status {
		return nil
	}
	if status == SlaveOffline {
		s.fail++
		if s.fail == 3 {
			err := s.ctx.Offline(s.dev)
			if err != nil {
				return err
			}
			s.status = SlaveOffline
			s.fail = 0
		}
	} else if status == SlaveOnline {
		err := s.ctx.Online(s.dev)
		if err != nil {
			return err
		}
		s.status = SlaveOnline
	}
	return nil
}