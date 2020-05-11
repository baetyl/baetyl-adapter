package main

type Slave struct {
	client *MbClient
	cfg    SlaveConfig
}

func NewSlave(cfg SlaveConfig, client *MbClient) *Slave {
	return &Slave{
		client: client,
		cfg:    cfg,
	}
}
