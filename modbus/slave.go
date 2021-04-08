package modbus

import dm "github.com/baetyl/baetyl-go/v2/dmcontext"

type Slave struct {
	client *MbClient
	dev    *dm.DeviceInfo
	cfg    SlaveConfig
}

func NewSlave(dev *dm.DeviceInfo, cfg SlaveConfig, client *MbClient) *Slave {
	return &Slave{
		dev:    dev,
		client: client,
		cfg:    cfg,
	}
}
