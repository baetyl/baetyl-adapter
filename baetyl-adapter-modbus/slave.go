package main

type Slave struct {
	client *MbClient
	cfg    SlaveConfig
}

func NewSlave(item SlaveConfig, client *MbClient) *Slave {
	return &Slave{
		client: client,
		cfg:    item,
	}
}
