package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/protocol/mqtt"
	"github.com/baetyl/baetyl/sdk/baetyl-go"
)

func main() {
	// Running module in baetyl context
	baetyl.Run(func(ctx baetyl.Context) error {
		var cfg Config
		// load custom config
		err := ctx.LoadConfig(&cfg)
		if err != nil {
			return err
		}
		log := ctx.Log()
		// create a hub client
		mqttCli, err := ctx.NewHubClient("", nil)
		if err != nil {
			return err
		}
		//start client to keep connection with hub
		mqttCli.Start(nil)
		defer mqttCli.Close()

		// modbus slave.go connection init
		clients := map[byte]handler{}
		slaves := map[byte]*Slave{}
		for _, slaveItem := range cfg.Slaves {
			client := NewClient(slaveItem)
			clients[slaveItem.ID] = client
			client.Connect()
			slaves[slaveItem.ID] = NewSlave(slaveItem, client)
		}
		defer func() {
			for _, client := range clients {
				client.Close()
			}
		}()
		var wg sync.WaitGroup
		for _, item := range cfg.Tables {
			m := NewMap(item, slaves[item.SlaveID])
			wg.Add(1)
			go func(m Map, wg *sync.WaitGroup) {
				ticker := time.NewTicker(m.s.cfg.Interval)
				defer ticker.Stop()
				defer wg.Done()
				for {
					select {
					case <-ticker.C:
						err := execute(m, mqttCli, cfg.Publish)
						if err != nil {
							log.Errorf(err.Error())
						}
					case <-ctx.WaitChan():
						return
					}
				}
			}(*m, &wg)
		}
		wg.Wait()
		return nil
	})
}

func execute(m Map, mqttCli *mqtt.Dispatcher, publish *Publish) error {
	results, err := m.Read()
	if err != nil {
		return fmt.Errorf("failed to collect data from slave.go: id=%d function=%d address=%d quantity=%d",
			m.cfg.SlaveID, m.cfg.Function, m.cfg.Address, m.cfg.Quantity)
	}
	pld := make([]byte, 9+m.cfg.Quantity*2)
	pld[0] = m.cfg.SlaveID
	binary.BigEndian.PutUint16(pld[1:], m.cfg.Address)
	binary.BigEndian.PutUint16(pld[3:], m.cfg.Quantity)
	binary.BigEndian.PutUint32(pld[5:], uint32(time.Now().Unix()))
	copy(pld[9:], results)
	pkt := packet.NewPublish()
	pkt.Message.Topic = publish.Topic
	pkt.Message.QOS = packet.QOS(publish.QOS)
	pkt.Message.Payload = pld
	err = mqttCli.Send(pkt)
	if err != nil {
		return fmt.Errorf("failed to publish: %s", err.Error())
	}
	return nil
}
