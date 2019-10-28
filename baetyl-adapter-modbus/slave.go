package main

type Slave struct {
	client *MbClient
	cfg    SlaveItem
}

func NewSlave(item SlaveItem, client *MbClient) *Slave {
	return &Slave{
		client: client,
		cfg:    item,
	}
}
