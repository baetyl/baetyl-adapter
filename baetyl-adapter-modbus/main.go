package main

import (
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
		ctx.Config().Hub.BufferSize = len(cfg.Maps) * 2
		mqttCli, err := ctx.NewHubClient("", nil)
		if err != nil {
			return err
		}
		//start client to keep connection with hub
		mqttCli.Start(nil)
		defer mqttCli.Close()

		modbus := newModbus(cfg, mqttCli, log)
		modbus.Start(ctx)
		defer modbus.Close()
		//ctx.Wait()
		return nil
	})
}
