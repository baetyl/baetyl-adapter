package main

import (
	"container/heap"
	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl/sdk/baetyl-go"
	"time"
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
		if cfg.Publish == nil {
			log.Errorf("topic of collected data can not be empty")
			return nil
		}
		// create a hub client
		mqttCli, err := ctx.NewHubClient("", nil)
		if err != nil {
			return err
		}
		//start client to keep connection with hub
		mqttCli.Start(nil)

		// modbus slave connection init
		clients := map[byte]*client{}
		for _, slave := range cfg.Slaves {
			client, err := newClient(slave)
			if err != nil {
				log.Errorf("failed to connect slave: id=%d, %s", slave.ID, err.Error())
			} else {
				clients[slave.ID] = client
			}
		}
		msgChan := make(chan *packet.Publish, cfg.MsgBufferSize)
		sender := newSender(mqttCli, msgChan)
		err = sender.start()
		if err != nil {
			return err
		}
		defer func() {
			for _, client := range clients {
				err := client.close()
				if err != nil {
					log.Errorf("failed to close modbus client %s", err.Error())
				}
			}
			sender.close()
			err = mqttCli.Close()
			if err != nil {
				log.Errorf("failed to close mqtt client %s", err.Error())
			}
		}()
		//create task according to parse item
		pq := make(PriorityQueue, 0)
		heap.Init(&pq)
		for _, parse := range cfg.ParseItems {
			client, ok := clients[parse.SlaveID]
			if ok {
				if parse.Interval == 0 {
					parse.Interval = client.cfg.Interval
				}
				ta := newTask(client, mqttCli, parse, msgChan, cfg.Publish, ctx.Log().WithField("modbus", "task"))
				heap.Push(&pq, ta)
			} else {
				log.Errorf("slave id=%d does not exist or disconnect", parse.SlaveID)
			}
		}
		//create a timer
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for i := 0; i < len(pq) && pq[i].timeUp(); i++ {
					pq[i].execute()
					// since runTime of task changed, index of task should be adjusted
					heap.Fix(&pq, i)
				}
			case <-ctx.WaitChan():
				// wait until service is stopped
				return nil
			}
		}
	})
}
